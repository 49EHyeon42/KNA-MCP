package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

const insectPilbkInfoPath = insectResourceBasePath + "/insectPilbkInfo"

var _ outbound.InsectPilbkInfoPort = (*Client)(nil)

type insectPilbkInfoBody struct {
	Item *insectPilbkInfoItem `xml:"item"`
}

type insectPilbkInfoItem struct {
	EcoDsrct         string `xml:"ecoDsrct"`
	EggDsrct         string `xml:"eggDsrct"`
	EmrgcCnt         string `xml:"emrgcCnt"`
	EmrgcEraDscrt    string `xml:"emrgcEraDscrt"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	FemaleDsrct      string `xml:"femaleDsrct"`
	GenusKorNm       string `xml:"genusKorNm"`
	GenusNm          string `xml:"genusNm"`
	GnrlDsrct        string `xml:"gnrlDsrct"`
	HabitDsrct       string `xml:"habitDsrct"`
	InsctEngNm       string `xml:"insctEngNm"`
	InsctGnrlNm      string `xml:"insctGnrlNm"`
	InsctPilbkNo     string `xml:"insctPilbkNo"`
	InsctSpecsScnm   string `xml:"insctSpecsScnm"`
	LarvaDsrct       string `xml:"larvaDsrct"`
	LastUpdtDtm      string `xml:"lastUpdtDtm"`
	MaleDsrct        string `xml:"maleDsrct"`
	MnmmOccrrCnt     string `xml:"mnmmOccrrCnt"`
	MxmmOccrrCnt     string `xml:"mxmmOccrrCnt"`
	OrdKorNm         string `xml:"ordKorNm"`
	OrdNm            string `xml:"ordNm"`
	PestDsrct        string `xml:"pestDsrct"`
	PupaDsrct        string `xml:"pupaDsrct"`
	ReferDsrct       string `xml:"referDsrct"`
	SubFamilyKorNm   string `xml:"subFamilyKorNm"`
	SubFamilyNm      string `xml:"subFamilyNm"`
	SuperFamilyKorNm string `xml:"superFamilyKorNm"`
	SuperFamilyNm    string `xml:"superFamilyNm"`
	WinterDsrct      string `xml:"winterDsrct"`
}

type insectPilbkInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body insectPilbkInfoBody `xml:"body"`
}

// InsectPilbkInfoError reports an error returned by insectPilbkInfo.
type InsectPilbkInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the insectPilbkInfo error message.
func (e *InsectPilbkInfoError) Error() string {
	return fmt.Sprintf("insectPilbkInfo: API error %s: %s", e.Code, e.Message)
}

// InsectPilbkInfo gets Korea National Arboretum insect pictorial book detail information.
func (c *Client) InsectPilbkInfo(ctx context.Context, query application.InsectPilbkInfoQuery) (application.InsectPilbkInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, insectPilbkInfoPath)
	if err != nil {
		return application.InsectPilbkInfoResult{}, fmt.Errorf("insectPilbkInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.InsectPilbkInfoResult{}, fmt.Errorf("insectPilbkInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("reqInsctPilbkNo", query.ReqInsctPilbkNo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.InsectPilbkInfoResult{}, fmt.Errorf("insectPilbkInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.InsectPilbkInfoResult{}, fmt.Errorf("insectPilbkInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload insectPilbkInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.InsectPilbkInfoResult{}, &InsectPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.InsectPilbkInfoResult{}, fmt.Errorf("insectPilbkInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.InsectPilbkInfoResult{}, fmt.Errorf("insectPilbkInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.InsectPilbkInfoResult{}, errors.New("insectPilbkInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != insectResourceSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = insectResourceResultMessages[payload.Header.ResultCode]
		}
		return application.InsectPilbkInfoResult{}, &InsectPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body insectPilbkInfoBody) result() application.InsectPilbkInfoResult {
	if body.Item == nil {
		return application.InsectPilbkInfoResult{}
	}

	return application.InsectPilbkInfoResult{Item: &application.InsectPilbkInfoItem{
		EcoDsrct:         body.Item.EcoDsrct,
		EggDsrct:         body.Item.EggDsrct,
		EmrgcCnt:         body.Item.EmrgcCnt,
		EmrgcEraDscrt:    body.Item.EmrgcEraDscrt,
		FamilyKorNm:      body.Item.FamilyKorNm,
		FamilyNm:         body.Item.FamilyNm,
		FemaleDsrct:      body.Item.FemaleDsrct,
		GenusKorNm:       body.Item.GenusKorNm,
		GenusNm:          body.Item.GenusNm,
		GnrlDsrct:        body.Item.GnrlDsrct,
		HabitDsrct:       body.Item.HabitDsrct,
		InsctEngNm:       body.Item.InsctEngNm,
		InsctGnrlNm:      body.Item.InsctGnrlNm,
		InsctPilbkNo:     body.Item.InsctPilbkNo,
		InsctSpecsScnm:   body.Item.InsctSpecsScnm,
		LarvaDsrct:       body.Item.LarvaDsrct,
		LastUpdtDtm:      body.Item.LastUpdtDtm,
		MaleDsrct:        body.Item.MaleDsrct,
		MnmmOccrrCnt:     body.Item.MnmmOccrrCnt,
		MxmmOccrrCnt:     body.Item.MxmmOccrrCnt,
		OrdKorNm:         body.Item.OrdKorNm,
		OrdNm:            body.Item.OrdNm,
		PestDsrct:        body.Item.PestDsrct,
		PupaDsrct:        body.Item.PupaDsrct,
		ReferDsrct:       body.Item.ReferDsrct,
		SubFamilyKorNm:   body.Item.SubFamilyKorNm,
		SubFamilyNm:      body.Item.SubFamilyNm,
		SuperFamilyKorNm: body.Item.SuperFamilyKorNm,
		SuperFamilyNm:    body.Item.SuperFamilyNm,
		WinterDsrct:      body.Item.WinterDsrct,
	}}
}
