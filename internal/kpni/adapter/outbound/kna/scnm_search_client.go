package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/outbound"
)

const scnmSearchPath = kpniBasePath + "/scnmSearch"

var _ outbound.ScnmSearchPort = (*Client)(nil)

type scnmSearchBody struct {
	Items      []scnmSearchItem `xml:"items>item"`
	NumOfRows  int              `xml:"numOfRows"`
	PageNo     int              `xml:"pageNo"`
	TotalCount int              `xml:"totalCount"`
}

type scnmSearchItem struct {
	ClassKorNm          string `xml:"classKorNm"`
	ClassNm             string `xml:"classNm"`
	FalmKorNm           string `xml:"falmKorNm"`
	FalmNm              string `xml:"falmNm"`
	GenusKorNm          string `xml:"genusKorNm"`
	GenusNm             string `xml:"genusNm"`
	LastUpdtDtm         string `xml:"lastUpdtDtm"`
	OrdKorNm            string `xml:"ordKorNm"`
	OrdNm               string `xml:"ordNm"`
	PhylumKorNm         string `xml:"phylumKorNm"`
	PhylumNm            string `xml:"phylumNm"`
	PlantGnrlNm         string `xml:"plantGnrlNm"`
	PlantScnmID         string `xml:"plantScnmId"`
	PlantSpecsClsscCdNm string `xml:"plantSpecsClsscCdNm"`
	PlantSpecsScnm      string `xml:"plantSpecsScnm"`
	StpltScnmRltnCdNm   string `xml:"stpltScnmRltnCdNm"`
	SubClassKorNm       string `xml:"subClassKorNm"`
	SubClassNm          string `xml:"subClassNm"`
}

type scnmSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body scnmSearchBody `xml:"body"`
}

// ScnmSearchError reports an error returned by scnmSearch.
type ScnmSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the scnmSearch error message.
func (e *ScnmSearchError) Error() string {
	return fmt.Sprintf("scnmSearch: API error %s: %s", e.Code, e.Message)
}

// ScnmSearch gets Korea National Arboretum scientific name information.
func (c *Client) ScnmSearch(ctx context.Context, query application.ScnmSearchQuery) (application.ScnmSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, scnmSearchPath)
	if err != nil {
		return application.ScnmSearchResult{}, fmt.Errorf("scnmSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.ScnmSearchResult{}, fmt.Errorf("scnmSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	if query.ReqGnrlNm != "" {
		values.Set("reqGnrlNm", query.ReqGnrlNm)
	}
	if query.ReqScnm != "" {
		values.Set("reqScnm", query.ReqScnm)
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
		return application.ScnmSearchResult{}, fmt.Errorf("scnmSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.ScnmSearchResult{}, fmt.Errorf("scnmSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload scnmSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.ScnmSearchResult{}, &ScnmSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.ScnmSearchResult{}, fmt.Errorf("scnmSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.ScnmSearchResult{}, fmt.Errorf("scnmSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.ScnmSearchResult{}, errors.New("scnmSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != kpniSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = kpniResultMessages[payload.Header.ResultCode]
		}
		return application.ScnmSearchResult{}, &ScnmSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body scnmSearchBody) result() application.ScnmSearchResult {
	items := make([]application.ScnmSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.ScnmSearchItem{
			ClassKorNm:          item.ClassKorNm,
			ClassNm:             item.ClassNm,
			FalmKorNm:           item.FalmKorNm,
			FalmNm:              item.FalmNm,
			GenusKorNm:          item.GenusKorNm,
			GenusNm:             item.GenusNm,
			LastUpdtDtm:         item.LastUpdtDtm,
			OrdKorNm:            item.OrdKorNm,
			OrdNm:               item.OrdNm,
			PhylumKorNm:         item.PhylumKorNm,
			PhylumNm:            item.PhylumNm,
			PlantGnrlNm:         item.PlantGnrlNm,
			PlantScnmID:         item.PlantScnmID,
			PlantSpecsClsscCdNm: item.PlantSpecsClsscCdNm,
			PlantSpecsScnm:      item.PlantSpecsScnm,
			StpltScnmRltnCdNm:   item.StpltScnmRltnCdNm,
			SubClassKorNm:       item.SubClassKorNm,
			SubClassNm:          item.SubClassNm,
		}
	}

	return application.ScnmSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
