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
	entogIlstrSearchPath        = entogServiceBasePath + "/entogIlstrSearch"
	entogIlstrSearchSuccessCode = "00"
)

var entogIlstrSearchResultMessages = map[string]string{
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

var _ outbound.EntogIlstrSearchPort = (*Client)(nil)

type entogIlstrSearchBody struct {
	Items      []entogIlstrSearchItem `xml:"items>item"`
	NumOfRows  int                    `xml:"numOfRows"`
	PageNo     int                    `xml:"pageNo"`
	TotalCount int                    `xml:"totalCount"`
}

type entogIlstrSearchItem struct {
	Btnc             string `xml:"btnc"`
	CprtCtnt         string `xml:"cprtCtnt"`
	DetailYn         string `xml:"detailYn"`
	EntogOfnmKrlngNm string `xml:"entogOfnmKrlngNm"`
	EntogOfnmScnmID  string `xml:"entogOfnmScnmId"`
	EntogPilbkNo     string `xml:"entogPilbkNo"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	GenusKorNm       string `xml:"genusKorNm"`
	GenusNm          string `xml:"genusNm"`
	ImgURL           string `xml:"imgUrl"`
	OrdKorNm         string `xml:"ordKorNm"`
	OrdNm            string `xml:"ordNm"`
}

type entogIlstrSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body entogIlstrSearchBody `xml:"body"`
}

// EntogIlstrSearchError reports an error returned by entogIlstrSearch.
type EntogIlstrSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the entogIlstrSearch error message.
func (e *EntogIlstrSearchError) Error() string {
	return fmt.Sprintf("entogIlstrSearch: API error %s: %s", e.Code, e.Message)
}

// EntogIlstrSearch searches Korea National Arboretum entognath pictorial book entries.
func (c *Client) EntogIlstrSearch(ctx context.Context, query application.EntogIlstrSearchQuery) (application.EntogIlstrSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, entogIlstrSearchPath)
	if err != nil {
		return application.EntogIlstrSearchResult{}, fmt.Errorf("entogIlstrSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.EntogIlstrSearchResult{}, fmt.Errorf("entogIlstrSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("st", query.St)
	values.Set("sw", query.Sw)
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.EntogIlstrSearchResult{}, fmt.Errorf("entogIlstrSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.EntogIlstrSearchResult{}, fmt.Errorf("entogIlstrSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload entogIlstrSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.EntogIlstrSearchResult{}, &EntogIlstrSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.EntogIlstrSearchResult{}, fmt.Errorf("entogIlstrSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.EntogIlstrSearchResult{}, fmt.Errorf("entogIlstrSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.EntogIlstrSearchResult{}, errors.New("entogIlstrSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != entogIlstrSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = entogIlstrSearchResultMessages[payload.Header.ResultCode]
		}
		return application.EntogIlstrSearchResult{}, &EntogIlstrSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body entogIlstrSearchBody) result() application.EntogIlstrSearchResult {
	items := make([]application.EntogIlstrSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.EntogIlstrSearchItem{
			Btnc:             item.Btnc,
			CprtCtnt:         item.CprtCtnt,
			DetailYn:         item.DetailYn,
			EntogOfnmKrlngNm: item.EntogOfnmKrlngNm,
			EntogOfnmScnmID:  item.EntogOfnmScnmID,
			EntogPilbkNo:     item.EntogPilbkNo,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			ImgURL:           item.ImgURL,
			OrdKorNm:         item.OrdKorNm,
			OrdNm:            item.OrdNm,
		}
	}

	return application.EntogIlstrSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
