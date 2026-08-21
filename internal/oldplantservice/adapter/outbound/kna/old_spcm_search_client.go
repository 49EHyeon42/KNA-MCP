package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application/port/outbound"
)

const (
	oldSpcmSearchPath        = "/1400119/OldPlantService/oldSpcmSearch"
	oldSpcmSearchSuccessCode = "00"
)

var oldSpcmSearchResultMessages = map[string]string{
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

var _ outbound.OldSpcmSearchPort = (*Client)(nil)

type oldSpcmSearchBody struct {
	Items      []oldSpcmSearchItem `xml:"items>item"`
	NumOfRows  int                 `xml:"numOfRows"`
	PageNo     int                 `xml:"pageNo"`
	TotalCount int                 `xml:"totalCount"`
}

type oldSpcmSearchItem struct {
	CprtCtnt       string `xml:"cprtCtnt"`
	FamlKorNm      string `xml:"famlKorNm"`
	FamlNm         string `xml:"famlNm"`
	FrstRgstnDtm   string `xml:"frstRgstnDtm"`
	ImgURL         string `xml:"imgUrl"`
	LastUpdtDtm    string `xml:"lastUpdtDtm"`
	PlantGnrlNm    string `xml:"plantGnrlNm"`
	PlantOldSmplNo string `xml:"plantOldSmplNo"`
	PlantSpecsScnm string `xml:"plantSpecsScnm"`
}

type oldSpcmSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body oldSpcmSearchBody `xml:"body"`
}

// OldSpcmSearchError reports an error returned by oldSpcmSearch.
type OldSpcmSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the oldSpcmSearch error message.
func (e *OldSpcmSearchError) Error() string {
	return fmt.Sprintf("oldSpcmSearch: API error %s: %s", e.Code, e.Message)
}

// OldSpcmSearch searches Korea National Arboretum old plant specimens.
func (c *Client) OldSpcmSearch(ctx context.Context, query application.OldSpcmSearchQuery) (application.OldSpcmSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, oldSpcmSearchPath)
	if err != nil {
		return application.OldSpcmSearchResult{}, fmt.Errorf("oldSpcmSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.OldSpcmSearchResult{}, fmt.Errorf("oldSpcmSearch: parse URL: %w", err)
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
		return application.OldSpcmSearchResult{}, fmt.Errorf("oldSpcmSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.OldSpcmSearchResult{}, fmt.Errorf("oldSpcmSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload oldSpcmSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.OldSpcmSearchResult{}, &OldSpcmSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.OldSpcmSearchResult{}, fmt.Errorf("oldSpcmSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.OldSpcmSearchResult{}, fmt.Errorf("oldSpcmSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.OldSpcmSearchResult{}, errors.New("oldSpcmSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != oldSpcmSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = oldSpcmSearchResultMessages[payload.Header.ResultCode]
		}
		return application.OldSpcmSearchResult{}, &OldSpcmSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body oldSpcmSearchBody) result() application.OldSpcmSearchResult {
	items := make([]application.OldSpcmSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.OldSpcmSearchItem{
			CprtCtnt:       item.CprtCtnt,
			FamlKorNm:      item.FamlKorNm,
			FamlNm:         item.FamlNm,
			FrstRgstnDtm:   item.FrstRgstnDtm,
			ImgURL:         item.ImgURL,
			LastUpdtDtm:    item.LastUpdtDtm,
			PlantGnrlNm:    item.PlantGnrlNm,
			PlantOldSmplNo: item.PlantOldSmplNo,
			PlantSpecsScnm: item.PlantSpecsScnm,
		}
	}

	return application.OldSpcmSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
