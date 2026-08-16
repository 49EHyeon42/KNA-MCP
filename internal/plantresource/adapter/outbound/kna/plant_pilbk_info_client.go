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

var _ outbound.PlantPictorialBookInformationPort = (*Client)(nil)

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
	Body struct {
		Item plantPilbkInfoItem `xml:"item"`
	} `xml:"body"`
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

// PlantPictorialBookInformation gets Korea National Arboretum plant pictorial book information.
func (c *Client) PlantPictorialBookInformation(ctx context.Context, query application.PlantPictorialBookInformationQuery) (application.PlantPictorialBookInformationResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantPilbkInfoPath)
	if err != nil {
		return application.PlantPictorialBookInformationResult{}, fmt.Errorf("plantPilbkInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantPictorialBookInformationResult{}, fmt.Errorf("plantPilbkInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("reqPlantPilbkNo", query.RequestPlantPictorialBookNumber)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantPictorialBookInformationResult{}, fmt.Errorf("plantPilbkInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantPictorialBookInformationResult{}, fmt.Errorf("plantPilbkInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantPilbkInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantPictorialBookInformationResult{}, &PlantPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantPictorialBookInformationResult{}, fmt.Errorf("plantPilbkInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantPictorialBookInformationResult{}, fmt.Errorf("plantPilbkInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantPictorialBookInformationResult{}, errors.New("plantPilbkInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != plantPilbkInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantPilbkInfoResultMessages[payload.Header.ResultCode]
		}
		return application.PlantPictorialBookInformationResult{}, &PlantPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.Item.result(), nil
}

func (item plantPilbkInfoItem) result() application.PlantPictorialBookInformationResult {
	return application.PlantPictorialBookInformationResult{
		APGFamilyKoreanName:           item.APGFamilyKorNm,
		APGFamilyName:                 item.APGFamilyNm,
		BfofMethod:                    item.BfofMthod,
		BreedingMethodDescription:     item.BrdMthdDesc,
		BugInformation:                item.BugInfo,
		Distribution:                  item.Dstrb,
		EnglishName:                   item.EngNm,
		FamilyKoreanName:              item.FamilyKorNm,
		FamilyName:                    item.FamilyNm,
		FarmSpecialFeatureDescription: item.FarmSpftDesc,
		GenusKoreanName:               item.GenusKorNm,
		GenusName:                     item.GenusNm,
		GrowthEnvironmentDescription:  item.GrwEvrntDesc,
		InductionDescription:          item.InductionDesc,
		LastUpdateDateTime:            item.LastUpdtDtm,
		NotRecommendedGeneralName:     item.NotRcmmGnrlNm,
		Note:                          item.Note,
		OriginPlaceName:               item.OrplcNm,
		OverseasDistribution:          item.OsDstrb,
		PlantGeneralName:              item.PlantGnrlNm,
		PlantPictorialBookNumber:      item.PlantPilbkNo,
		PlantSpeciesScientificName:    item.PlantSpecsScnm,
		ProtectionPlanDescription:     item.PrtcPlnDesc,
		RearingGubun:                  item.RrngGubun,
		RearingType:                   item.RrngType,
		Shape:                         item.Shpe,
		SimilarPlantDescription:       item.SmlrPlntDesc,
		SpecialFeature:                item.Spft,
		UseMethodDescription:          item.UseMthdDesc,
		WoodDescription:               item.WoodDesc,
	}
}
