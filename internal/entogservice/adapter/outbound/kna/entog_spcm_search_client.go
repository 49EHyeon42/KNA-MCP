package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/outbound"
)

const (
	entogSpcmSearchPath        = entogServiceBasePath + "/entogSpcmSearch"
	entogSpcmSearchSuccessCode = "00"
)

var entogSpcmSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"01": "APPLICATION_ERROR",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"04": "HTTP_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"12": "NO_OPENAPI_SERVICE_ERROR",
	"20": "SERVICE_ACCESS_DENIED_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"22": "LIMITED_NUMBER_OF_SERVICE_REQUESTS_EXCEEDS_ERROR",
	"30": "SERVICE_KEY_IS_NOT_REGISTERED_ERROR",
	"31": "DEADLINE_HAS_EXPIRED_ERROR",
	"32": "UNREGISTERED_IP_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
	"99": "UNKNOWN_ERROR",
}

var _ outbound.EntogSpcmSearchPort = (*Client)(nil)

type entogSpcmSearchBody struct {
	Items      []entogSpcmSearchItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type entogSpcmSearchItem struct {
	Btnc             string `xml:"btnc"`
	ClctDyDesc       string `xml:"clctDyDesc"`
	CprtCtnt         string `xml:"cprtCtnt"`
	DetailYn         string `xml:"detailYn"`
	EntogOfnmKrlngNm string `xml:"entogOfnmKrlngNm"`
	EntogSmplNo      string `xml:"entogSmplNo"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	FrstRgstnDtm     string `xml:"frstRgstnDtm"`
	GenusKorNm       string `xml:"genusKorNm"`
	GenusNm          string `xml:"genusNm"`
	ImgURL           string `xml:"imgUrl"`
	LastUpdtDtm      string `xml:"lastUpdtDtm"`
	OrdKorNm         string `xml:"ordKorNm"`
	OrdNm            string `xml:"ordNm"`
}

type entogSpcmSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body entogSpcmSearchBody `xml:"body"`
}

// EntogSpcmSearchError reports an error returned by entogSpcmSearch.
type EntogSpcmSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the entogSpcmSearch error message.
func (e *EntogSpcmSearchError) Error() string {
	return fmt.Sprintf("entogSpcmSearch: API error %s: %s", e.Code, e.Message)
}

// EntogSpcmSearch searches Korea National Arboretum entognath specimens.
func (c *Client) EntogSpcmSearch(ctx context.Context, query application.EntogSpcmSearchQuery) (application.EntogSpcmSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, entogSpcmSearchPath)
	if err != nil {
		return application.EntogSpcmSearchResult{}, fmt.Errorf("entogSpcmSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.EntogSpcmSearchResult{}, fmt.Errorf("entogSpcmSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("st", query.St)
	values.Set("sw", query.Sw)
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	if query.DateGbn != "" {
		values.Set("dateGbn", query.DateGbn)
	}
	if query.DateFrom != "" {
		values.Set("dateFrom", query.DateFrom)
	}
	if query.DateTo != "" {
		values.Set("dateTo", query.DateTo)
	}
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.EntogSpcmSearchResult{}, fmt.Errorf("entogSpcmSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.EntogSpcmSearchResult{}, fmt.Errorf("entogSpcmSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload entogSpcmSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.EntogSpcmSearchResult{}, &EntogSpcmSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.EntogSpcmSearchResult{}, fmt.Errorf("entogSpcmSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.EntogSpcmSearchResult{}, fmt.Errorf("entogSpcmSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.EntogSpcmSearchResult{}, errors.New("entogSpcmSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != entogSpcmSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = entogSpcmSearchResultMessages[payload.Header.ResultCode]
		}
		return application.EntogSpcmSearchResult{}, &EntogSpcmSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body entogSpcmSearchBody) result() application.EntogSpcmSearchResult {
	items := make([]application.EntogSpcmSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.EntogSpcmSearchItem{
			Btnc:             item.Btnc,
			ClctDyDesc:       item.ClctDyDesc,
			CprtCtnt:         item.CprtCtnt,
			DetailYn:         item.DetailYn,
			EntogOfnmKrlngNm: item.EntogOfnmKrlngNm,
			EntogSmplNo:      item.EntogSmplNo,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			FrstRgstnDtm:     item.FrstRgstnDtm,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			ImgURL:           item.ImgURL,
			LastUpdtDtm:      item.LastUpdtDtm,
			OrdKorNm:         item.OrdKorNm,
			OrdNm:            item.OrdNm,
		}
	}

	return application.EntogSpcmSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
