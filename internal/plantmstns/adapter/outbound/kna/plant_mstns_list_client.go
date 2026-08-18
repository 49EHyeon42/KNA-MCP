package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/port/outbound"
)

const (
	plantMstnsListPath        = "/1400119/PlantMstnsService/plantMstnsList"
	plantMstnsListSuccessCode = "00"
)

var plantMstnsListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantMstnsListPort = (*Client)(nil)

type plantMstnsListBody struct {
	Items      []plantMstnsListItem `xml:"items>item"`
	NumOfRows  int                  `xml:"numOfRows"`
	PageNo     int                  `xml:"pageNo"`
	TotalCount int                  `xml:"totalCount"`
}

type plantMstnsListItem struct {
	DistrAraDscrt         string `xml:"distrAraDscrt"`
	MinitrTpcdNm          string `xml:"minitrTpcdNm"`
	PlantBrdgFomTpcdNm    string `xml:"plantBrdgFomTpcdNm"`
	PlantGnrlNm           string `xml:"plantGnrlNm"`
	PlantMinitrAthrNm     string `xml:"plantMinitrAthrNm"`
	PlantMinitrMnfctMonth string `xml:"plantMinitrMnfctMonth"`
	PlantMinitrMnfctYr    string `xml:"plantMinitrMnfctYr"`
	PlantMinitrPsinsNm    string `xml:"plantMinitrPsinsNm"`
	PlantSpecsScnm        string `xml:"plantSpecsScnm"`
	RrnssPlantYn          string `xml:"rrnssPlantYn"`
	SpcltPlantYn          string `xml:"spcltPlantYn"`
}

type plantMstnsListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantMstnsListBody `xml:"body"`
}

// PlantMstnsListError reports an error returned by plantMstnsList.
type PlantMstnsListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantMstnsList error message.
func (e *PlantMstnsListError) Error() string {
	return fmt.Sprintf("plantMstnsList: API error %s: %s", e.Code, e.Message)
}

// PlantMstnsList gets Korea National Arboretum plant miniature information.
func (c *Client) PlantMstnsList(ctx context.Context, query application.PlantMstnsListQuery) (application.PlantMstnsListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantMstnsListPath)
	if err != nil {
		return application.PlantMstnsListResult{}, fmt.Errorf("plantMstnsList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantMstnsListResult{}, fmt.Errorf("plantMstnsList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	if query.ReqSearchWrd != "" {
		values.Set("reqSearchWrd", query.ReqSearchWrd)
	}
	if query.ReqMnfctYr != "" {
		values.Set("reqMnfctYr", query.ReqMnfctYr)
	}
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantMstnsListResult{}, fmt.Errorf("plantMstnsList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantMstnsListResult{}, fmt.Errorf("plantMstnsList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantMstnsListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantMstnsListResult{}, &PlantMstnsListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantMstnsListResult{}, fmt.Errorf("plantMstnsList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantMstnsListResult{}, fmt.Errorf("plantMstnsList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantMstnsListResult{}, errors.New("plantMstnsList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantMstnsListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantMstnsListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantMstnsListResult{}, &PlantMstnsListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantMstnsListBody) result() application.PlantMstnsListResult {
	items := make([]application.PlantMstnsListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantMstnsListItem{
			DistrAraDscrt:         item.DistrAraDscrt,
			MinitrTpcdNm:          item.MinitrTpcdNm,
			PlantBrdgFomTpcdNm:    item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:           item.PlantGnrlNm,
			PlantMinitrAthrNm:     item.PlantMinitrAthrNm,
			PlantMinitrMnfctMonth: item.PlantMinitrMnfctMonth,
			PlantMinitrMnfctYr:    item.PlantMinitrMnfctYr,
			PlantMinitrPsinsNm:    item.PlantMinitrPsinsNm,
			PlantSpecsScnm:        item.PlantSpecsScnm,
			RrnssPlantYn:          item.RrnssPlantYn,
			SpcltPlantYn:          item.SpcltPlantYn,
		}
	}

	return application.PlantMstnsListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
