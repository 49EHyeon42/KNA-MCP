package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

const (
	plantPilbkInfoPath        = plantResourceBasePath + "/plantPilbkInfo"
	plantPilbkInfoSuccessCode = "00"
)

var plantPilbkInfoResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantPilbkInfoPort = (*Client)(nil)

type plantPilbkInfoItem struct {
	APGFamilyKorNm string `xml:"apgFamilyKorNm"`
	APGFamilyNm    string `xml:"apgFamilyNm"`
	BfofMthod      string `xml:"bfofMthod"`
	BrdMthdDesc    string `xml:"brdMthdDesc"`
	BugInfo        string `xml:"bugInfo"`
	Dstrb          string `xml:"dstrb"`
	EngNm          string `xml:"engNm"`
	FamilyKorNm    string `xml:"familyKorNm"`
	FamilyNm       string `xml:"familyNm"`
	FarmSpftDesc   string `xml:"farmSpftDesc"`
	GenusKorNm     string `xml:"genusKorNm"`
	GenusNm        string `xml:"genusNm"`
	GrwEvrntDesc   string `xml:"grwEvrntDesc"`
	InductionDesc  string `xml:"inductionDesc"`
	LastUpdtDtm    string `xml:"lastUpdtDtm"`
	NotRcmmGnrlNm  string `xml:"notRcmmGnrlNm"`
	Note           string `xml:"note"`
	OrplcNm        string `xml:"orplcNm"`
	OsDstrb        string `xml:"osDstrb"`
	PlantGnrlNm    string `xml:"plantGnrlNm"`
	PlantPilbkNo   string `xml:"plantPilbkNo"`
	PlantSpecsScnm string `xml:"plantSpecsScnm"`
	PrtcPlnDesc    string `xml:"prtcPlnDesc"`
	RrngGubun      string `xml:"rrngGubun"`
	RrngType       string `xml:"rrngType"`
	Shpe           string `xml:"shpe"`
	SmlrPlntDesc   string `xml:"smlrPlntDesc"`
	Spft           string `xml:"spft"`
	UseMthdDesc    string `xml:"useMthdDesc"`
	WoodDesc       string `xml:"woodDesc"`
}

type plantPilbkInfoBody struct {
	Item plantPilbkInfoItem `xml:"item"`
}

type plantPilbkInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantPilbkInfoBody `xml:"body"`
}

// PlantPilbkInfoError reports an error returned by plantPilbkInfo.
type PlantPilbkInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantPilbkInfo error message.
func (e *PlantPilbkInfoError) Error() string {
	return fmt.Sprintf("plantPilbkInfo: API error %s: %s", e.Code, e.Message)
}

// PlantPilbkInfo gets Korea National Arboretum plant pictorial book information.
func (c *Client) PlantPilbkInfo(ctx context.Context, query application.PlantPilbkInfoQuery) (application.PlantPilbkInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantPilbkInfoPath)
	if err != nil {
		return application.PlantPilbkInfoResult{}, fmt.Errorf("plantPilbkInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantPilbkInfoResult{}, fmt.Errorf("plantPilbkInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("reqPlantPilbkNo", query.ReqPlantPilbkNo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantPilbkInfoResult{}, fmt.Errorf("plantPilbkInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantPilbkInfoResult{}, fmt.Errorf("plantPilbkInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantPilbkInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantPilbkInfoResult{}, &PlantPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantPilbkInfoResult{}, fmt.Errorf("plantPilbkInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantPilbkInfoResult{}, fmt.Errorf("plantPilbkInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantPilbkInfoResult{}, errors.New("plantPilbkInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != plantPilbkInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantPilbkInfoResultMessages[payload.Header.ResultCode]
		}
		return application.PlantPilbkInfoResult{}, &PlantPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.Item.result(), nil
}

func (item plantPilbkInfoItem) result() application.PlantPilbkInfoResult {
	return application.PlantPilbkInfoResult{
		APGFamilyKorNm: item.APGFamilyKorNm,
		APGFamilyNm:    item.APGFamilyNm,
		BfofMthod:      item.BfofMthod,
		BrdMthdDesc:    item.BrdMthdDesc,
		BugInfo:        item.BugInfo,
		Dstrb:          item.Dstrb,
		EngNm:          item.EngNm,
		FamilyKorNm:    item.FamilyKorNm,
		FamilyNm:       item.FamilyNm,
		FarmSpftDesc:   item.FarmSpftDesc,
		GenusKorNm:     item.GenusKorNm,
		GenusNm:        item.GenusNm,
		GrwEvrntDesc:   item.GrwEvrntDesc,
		InductionDesc:  item.InductionDesc,
		LastUpdtDtm:    item.LastUpdtDtm,
		NotRcmmGnrlNm:  item.NotRcmmGnrlNm,
		Note:           item.Note,
		OrplcNm:        item.OrplcNm,
		OsDstrb:        item.OsDstrb,
		PlantGnrlNm:    item.PlantGnrlNm,
		PlantPilbkNo:   item.PlantPilbkNo,
		PlantSpecsScnm: item.PlantSpecsScnm,
		PrtcPlnDesc:    item.PrtcPlnDesc,
		RrngGubun:      item.RrngGubun,
		RrngType:       item.RrngType,
		Shpe:           item.Shpe,
		SmlrPlntDesc:   item.SmlrPlntDesc,
		Spft:           item.Spft,
		UseMthdDesc:    item.UseMthdDesc,
		WoodDesc:       item.WoodDesc,
	}
}
