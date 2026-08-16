package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

const (
	plantSeedSearchPath        = plantResourceBasePath + "/plantSeedSearch"
	plantSeedSearchSuccessCode = "00"
)

var plantSeedSearchResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantSeedSearchPort = (*Client)(nil)

type plantSeedSearchBody struct {
	Items      []plantSeedSearchItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type plantSeedSearchItem struct {
	APGFamilyKorNm   string `xml:"apgFamilyKorNm"`
	APGFamilyNm      string `xml:"apgFamilyNm"`
	BlprdEnmnt       string `xml:"blprdEnmnt"`
	BlprdStmnt       string `xml:"blprdStmnt"`
	ClrngMthodCdNm   string `xml:"clrngMthodCdNm"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	FritCdNm         string `xml:"fritCdNm"`
	FrssnEnmnt       string `xml:"frssnEnmnt"`
	FrssnStmnt       string `xml:"frssnStmnt"`
	LastUpdtDtm      string `xml:"lastUpdtDtm"`
	PlantGnrlNm      string `xml:"plantGnrlNm"`
	PlantSpecsScnm   string `xml:"plantSpecsScnm"`
	RfrncLtrtrCont   string `xml:"rfrncLtrtrCont"`
	SeedCtsrfcDesc   string `xml:"seedCtsrfcDesc"`
	SeedCtsrfcTpcdNm string `xml:"seedCtsrfcTpcdNm"`
	SeedEmbrTpcdNm   string `xml:"seedEmbrTpcdNm"`
	SeedMnmmBrdth    string `xml:"seedMnmmBrdth"`
	SeedMnmmLngth    string `xml:"seedMnmmLngth"`
	SeedMxmmBrdth    string `xml:"seedMxmmBrdth"`
	SeedMxmmLngth    string `xml:"seedMxmmLngth"`
	SeedShpDesc      string `xml:"seedShpDesc"`
	SeedShpTpcdNm    string `xml:"seedShpTpcdNm"`
	SeedSpecsID      string `xml:"seedSpecsId"`
	SeedTpcdNm       string `xml:"seedTpcdNm"`
}

type plantSeedSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantSeedSearchBody `xml:"body"`
}

// PlantSeedSearchError reports an error returned by plantSeedSearch.
type PlantSeedSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantSeedSearch error message.
func (e *PlantSeedSearchError) Error() string {
	return fmt.Sprintf("plantSeedSearch: API error %s: %s", e.Code, e.Message)
}

// PlantSeedSearch searches Korea National Arboretum plant seed information.
func (c *Client) PlantSeedSearch(ctx context.Context, query application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantSeedSearchPath)
	if err != nil {
		return application.PlantSeedSearchResult{}, fmt.Errorf("plantSeedSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantSeedSearchResult{}, fmt.Errorf("plantSeedSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	setQueryValue(values, "dateFrom", query.DateFrom)
	setQueryValue(values, "dateTo", query.DateTo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantSeedSearchResult{}, fmt.Errorf("plantSeedSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantSeedSearchResult{}, fmt.Errorf("plantSeedSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantSeedSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantSeedSearchResult{}, &PlantSeedSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantSeedSearchResult{}, fmt.Errorf("plantSeedSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantSeedSearchResult{}, fmt.Errorf("plantSeedSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantSeedSearchResult{}, errors.New("plantSeedSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != plantSeedSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantSeedSearchResultMessages[payload.Header.ResultCode]
		}
		return application.PlantSeedSearchResult{}, &PlantSeedSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantSeedSearchBody) result() application.PlantSeedSearchResult {
	items := make([]application.PlantSeedSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantSeedSearchItem{
			APGFamilyKorNm:   item.APGFamilyKorNm,
			APGFamilyNm:      item.APGFamilyNm,
			BlprdEnmnt:       item.BlprdEnmnt,
			BlprdStmnt:       item.BlprdStmnt,
			ClrngMthodCdNm:   item.ClrngMthodCdNm,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			FritCdNm:         item.FritCdNm,
			FrssnEnmnt:       item.FrssnEnmnt,
			FrssnStmnt:       item.FrssnStmnt,
			LastUpdtDtm:      item.LastUpdtDtm,
			PlantGnrlNm:      item.PlantGnrlNm,
			PlantSpecsScnm:   item.PlantSpecsScnm,
			RfrncLtrtrCont:   item.RfrncLtrtrCont,
			SeedCtsrfcDesc:   item.SeedCtsrfcDesc,
			SeedCtsrfcTpcdNm: item.SeedCtsrfcTpcdNm,
			SeedEmbrTpcdNm:   item.SeedEmbrTpcdNm,
			SeedMnmmBrdth:    item.SeedMnmmBrdth,
			SeedMnmmLngth:    item.SeedMnmmLngth,
			SeedMxmmBrdth:    item.SeedMxmmBrdth,
			SeedMxmmLngth:    item.SeedMxmmLngth,
			SeedShpDesc:      item.SeedShpDesc,
			SeedShpTpcdNm:    item.SeedShpTpcdNm,
			SeedSpecsID:      item.SeedSpecsID,
			SeedTpcdNm:       item.SeedTpcdNm,
		}
	}

	return application.PlantSeedSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
