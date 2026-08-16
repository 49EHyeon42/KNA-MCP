package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantPictorialBookInformationToolName = "plant_resource_plant_pictorial_book_information"

type plantPictorialBookInformationInput struct {
	RequestPlantPictorialBookNumber string `json:"requestPlantPictorialBookNumber" jsonschema:"식물도감 목록 검색 결과의 식물도감번호"`
}

type plantPictorialBookInformationOutput struct {
	APGFamilyKoreanName           string `json:"apgFamilyKoreanName"`
	APGFamilyName                 string `json:"apgFamilyName"`
	BfofMethod                    string `json:"bfofMethod"`
	BreedingMethodDescription     string `json:"breedingMethodDescription"`
	BugInformation                string `json:"bugInformation"`
	Distribution                  string `json:"distribution"`
	EnglishName                   string `json:"englishName"`
	FamilyKoreanName              string `json:"familyKoreanName"`
	FamilyName                    string `json:"familyName"`
	FarmSpecialFeatureDescription string `json:"farmSpecialFeatureDescription"`
	GenusKoreanName               string `json:"genusKoreanName"`
	GenusName                     string `json:"genusName"`
	GrowthEnvironmentDescription  string `json:"growthEnvironmentDescription"`
	InductionDescription          string `json:"inductionDescription"`
	LastUpdateDateTime            string `json:"lastUpdateDateTime"`
	NotRecommendedGeneralName     string `json:"notRecommendedGeneralName"`
	Note                          string `json:"note"`
	OriginPlaceName               string `json:"originPlaceName"`
	OverseasDistribution          string `json:"overseasDistribution"`
	PlantGeneralName              string `json:"plantGeneralName"`
	PlantPictorialBookNumber      string `json:"plantPictorialBookNumber"`
	PlantSpeciesScientificName    string `json:"plantSpeciesScientificName"`
	ProtectionPlanDescription     string `json:"protectionPlanDescription"`
	RearingGubun                  string `json:"rearingGubun"`
	RearingType                   string `json:"rearingType"`
	Shape                         string `json:"shape"`
	SimilarPlantDescription       string `json:"similarPlantDescription"`
	SpecialFeature                string `json:"specialFeature"`
	UseMethodDescription          string `json:"useMethodDescription"`
	WoodDescription               string `json:"woodDescription"`
}

type plantPictorialBookInformationHandler struct {
	useCase inbound.PlantPictorialBookInformationUseCase
}

func addPlantPictorialBookInformationTool(server *mcp.Server, useCase inbound.PlantPictorialBookInformationUseCase) {
	handler := plantPictorialBookInformationHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantPictorialBookInformationToolName,
		Description: "산림청 국립수목원 식물도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h plantPictorialBookInformationHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantPictorialBookInformationInput) (*mcp.CallToolResult, plantPictorialBookInformationOutput, error) {
	result, err := h.useCase.PlantPictorialBookInformation(ctx, application.PlantPictorialBookInformationQuery{
		RequestPlantPictorialBookNumber: input.RequestPlantPictorialBookNumber,
	})
	if err != nil {
		return nil, plantPictorialBookInformationOutput{}, err
	}

	return nil, plantPictorialBookInformationOutput{
		APGFamilyKoreanName:           result.APGFamilyKoreanName,
		APGFamilyName:                 result.APGFamilyName,
		BfofMethod:                    result.BfofMethod,
		BreedingMethodDescription:     result.BreedingMethodDescription,
		BugInformation:                result.BugInformation,
		Distribution:                  result.Distribution,
		EnglishName:                   result.EnglishName,
		FamilyKoreanName:              result.FamilyKoreanName,
		FamilyName:                    result.FamilyName,
		FarmSpecialFeatureDescription: result.FarmSpecialFeatureDescription,
		GenusKoreanName:               result.GenusKoreanName,
		GenusName:                     result.GenusName,
		GrowthEnvironmentDescription:  result.GrowthEnvironmentDescription,
		InductionDescription:          result.InductionDescription,
		LastUpdateDateTime:            result.LastUpdateDateTime,
		NotRecommendedGeneralName:     result.NotRecommendedGeneralName,
		Note:                          result.Note,
		OriginPlaceName:               result.OriginPlaceName,
		OverseasDistribution:          result.OverseasDistribution,
		PlantGeneralName:              result.PlantGeneralName,
		PlantPictorialBookNumber:      result.PlantPictorialBookNumber,
		PlantSpeciesScientificName:    result.PlantSpeciesScientificName,
		ProtectionPlanDescription:     result.ProtectionPlanDescription,
		RearingGubun:                  result.RearingGubun,
		RearingType:                   result.RearingType,
		Shape:                         result.Shape,
		SimilarPlantDescription:       result.SimilarPlantDescription,
		SpecialFeature:                result.SpecialFeature,
		UseMethodDescription:          result.UseMethodDescription,
		WoodDescription:               result.WoodDescription,
	}, nil
}
