package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/outbound"
)

const (
	childPilbkInfoPath        = childServiceBasePath + "/childPilbkInfo"
	childPilbkInfoSuccessCode = "00"
)

var childPilbkInfoResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.ChildPilbkInfoPort = (*Client)(nil)

type childPilbkInfoBody struct {
	Item *childPilbkInfoItem `xml:"item"`
}

type childPilbkInfoItem struct {
	BiogyNm           string `xml:"biogyNm"`
	ChildLvbngPilbkNo string `xml:"childLvbngPilbkNo"`
	ExtrmCrss         string `xml:"extrmCrss"`
	FamilyKorNm       string `xml:"familyKorNm"`
	FamilyNm          string `xml:"familyNm"`
	GenusKorNm        string `xml:"genusKorNm"`
	GenusNm           string `xml:"genusNm"`
	HbttFieldYn       string `xml:"hbttFieldYn"`
	HbttFrestYn       string `xml:"hbttFrestYn"`
	HbttRiverYn       string `xml:"hbttRiverYn"`
	LvbngDscrt        string `xml:"lvbngDscrt"`
	LvbngTpcdNm       string `xml:"lvbngTpcdNm"`
	LvngKrlngNm       string `xml:"lvngKrlngNm"`
	PrtctSpecsTpcdNm  string `xml:"prtctSpecsTpcdNm"`
}

type childPilbkInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body childPilbkInfoBody `xml:"body"`
}

// ChildPilbkInfoError reports an error returned by childPilbkInfo.
type ChildPilbkInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the childPilbkInfo error message.
func (e *ChildPilbkInfoError) Error() string {
	return fmt.Sprintf("childPilbkInfo: API error %s: %s", e.Code, e.Message)
}

// ChildPilbkInfo gets Korea National Arboretum child pictorial book detail information.
func (c *Client) ChildPilbkInfo(ctx context.Context, query application.ChildPilbkInfoQuery) (application.ChildPilbkInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, childPilbkInfoPath)
	if err != nil {
		return application.ChildPilbkInfoResult{}, fmt.Errorf("childPilbkInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.ChildPilbkInfoResult{}, fmt.Errorf("childPilbkInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("reqChildLvbngPilbkNo", query.ReqChildLvbngPilbkNo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.ChildPilbkInfoResult{}, fmt.Errorf("childPilbkInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.ChildPilbkInfoResult{}, fmt.Errorf("childPilbkInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload childPilbkInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.ChildPilbkInfoResult{}, &ChildPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.ChildPilbkInfoResult{}, fmt.Errorf("childPilbkInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.ChildPilbkInfoResult{}, fmt.Errorf("childPilbkInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.ChildPilbkInfoResult{}, errors.New("childPilbkInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != childPilbkInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = childPilbkInfoResultMessages[payload.Header.ResultCode]
		}
		return application.ChildPilbkInfoResult{}, &ChildPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body childPilbkInfoBody) result() application.ChildPilbkInfoResult {
	if body.Item == nil {
		return application.ChildPilbkInfoResult{}
	}

	return application.ChildPilbkInfoResult{Item: &application.ChildPilbkInfoItem{
		BiogyNm:           body.Item.BiogyNm,
		ChildLvbngPilbkNo: body.Item.ChildLvbngPilbkNo,
		ExtrmCrss:         body.Item.ExtrmCrss,
		FamilyKorNm:       body.Item.FamilyKorNm,
		FamilyNm:          body.Item.FamilyNm,
		GenusKorNm:        body.Item.GenusKorNm,
		GenusNm:           body.Item.GenusNm,
		HbttFieldYn:       body.Item.HbttFieldYn,
		HbttFrestYn:       body.Item.HbttFrestYn,
		HbttRiverYn:       body.Item.HbttRiverYn,
		LvbngDscrt:        body.Item.LvbngDscrt,
		LvbngTpcdNm:       body.Item.LvbngTpcdNm,
		LvngKrlngNm:       body.Item.LvngKrlngNm,
		PrtctSpecsTpcdNm:  body.Item.PrtctSpecsTpcdNm,
	}}
}
