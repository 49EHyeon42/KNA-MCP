package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/outbound"
)

const (
	alchnIlstrSearchPath        = "/1400119/LchnService/alchnIlstrSearch"
	alchnIlstrSearchSuccessCode = "00"
)

var alchnIlstrSearchResultMessages = map[string]string{
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

var _ outbound.AlchnIlstrSearchPort = (*Client)(nil)

type alchnIlstrSearchBody struct {
	Items      []alchnIlstrSearchItem `xml:"items>item"`
	NumOfRows  int                    `xml:"numOfRows"`
	PageNo     int                    `xml:"pageNo"`
	TotalCount int                    `xml:"totalCount"`
}

type alchnIlstrSearchItem struct {
	Btnc         string `xml:"btnc"`
	CprtCtnt     string `xml:"cprtCtnt"`
	DetailYn     string `xml:"detailYn"`
	EngNm        string `xml:"engNm"`
	FamilyKorNm  string `xml:"familyKorNm"`
	FamilyNm     string `xml:"familyNm"`
	FrstRgstnDtm string `xml:"frstRgstnDtm"`
	GenusKorNm   string `xml:"genusKorNm"`
	GenusNm      string `xml:"genusNm"`
	ImgURL       string `xml:"imgUrl"`
	JapNm        string `xml:"japNm"`
	LastUpdtDtm  string `xml:"lastUpdtDtm"`
	LchnGnrlNm   string `xml:"lchnGnrlNm"`
	LchnInfrpNm  string `xml:"lchnInfrpNm"`
	LchnPilbkNo  string `xml:"lchnPilbkNo"`
	LchnScnmID   string `xml:"lchnScnmId"`
	LchnTtnm     string `xml:"lchnTtnm"`
	PrkNm        string `xml:"prkNm"`
}

type alchnIlstrSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body alchnIlstrSearchBody `xml:"body"`
}

// AlchnIlstrSearchError reports an error returned by alchnIlstrSearch.
type AlchnIlstrSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the alchnIlstrSearch error message.
func (e *AlchnIlstrSearchError) Error() string {
	return fmt.Sprintf("alchnIlstrSearch: API error %s: %s", e.Code, e.Message)
}

// AlchnIlstrSearch searches the Korea National Arboretum lichen pictorial book.
func (c *Client) AlchnIlstrSearch(ctx context.Context, query application.AlchnIlstrSearchQuery) (application.AlchnIlstrSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, alchnIlstrSearchPath)
	if err != nil {
		return application.AlchnIlstrSearchResult{}, fmt.Errorf("alchnIlstrSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.AlchnIlstrSearchResult{}, fmt.Errorf("alchnIlstrSearch: parse URL: %w", err)
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
		return application.AlchnIlstrSearchResult{}, fmt.Errorf("alchnIlstrSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.AlchnIlstrSearchResult{}, fmt.Errorf("alchnIlstrSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload alchnIlstrSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.AlchnIlstrSearchResult{}, &AlchnIlstrSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.AlchnIlstrSearchResult{}, fmt.Errorf("alchnIlstrSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.AlchnIlstrSearchResult{}, fmt.Errorf("alchnIlstrSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.AlchnIlstrSearchResult{}, errors.New("alchnIlstrSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != alchnIlstrSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = alchnIlstrSearchResultMessages[payload.Header.ResultCode]
		}
		return application.AlchnIlstrSearchResult{}, &AlchnIlstrSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body alchnIlstrSearchBody) result() application.AlchnIlstrSearchResult {
	items := make([]application.AlchnIlstrSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.AlchnIlstrSearchItem{
			Btnc:         item.Btnc,
			CprtCtnt:     item.CprtCtnt,
			DetailYn:     item.DetailYn,
			EngNm:        item.EngNm,
			FamilyKorNm:  item.FamilyKorNm,
			FamilyNm:     item.FamilyNm,
			FrstRgstnDtm: item.FrstRgstnDtm,
			GenusKorNm:   item.GenusKorNm,
			GenusNm:      item.GenusNm,
			ImgURL:       item.ImgURL,
			JapNm:        item.JapNm,
			LastUpdtDtm:  item.LastUpdtDtm,
			LchnGnrlNm:   item.LchnGnrlNm,
			LchnInfrpNm:  item.LchnInfrpNm,
			LchnPilbkNo:  item.LchnPilbkNo,
			LchnScnmID:   item.LchnScnmID,
			LchnTtnm:     item.LchnTtnm,
			PrkNm:        item.PrkNm,
		}
	}

	return application.AlchnIlstrSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
