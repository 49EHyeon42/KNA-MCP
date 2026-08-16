package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantPilbkInfoToolName = "plant_resource_plant_pictorial_book_information"

type plantPilbkInfoInput struct {
	ReqPlantPilbkNo string `json:"requestPlantPictorialBookNumber" jsonschema:"식물도감 목록 검색 결과의 식물도감번호"`
}

type plantPilbkInfoOutput struct {
	APGFamilyKorNm string `json:"apgFamilyKoreanName"`
	APGFamilyNm    string `json:"apgFamilyName"`
	BfofMthod      string `json:"pestControlMethod"`
	BrdMthdDesc    string `json:"breedingMethodDescription"`
	BugInfo        string `json:"bugInformation"`
	Dstrb          string `json:"distribution"`
	EngNm          string `json:"englishName"`
	FamilyKorNm    string `json:"familyKoreanName"`
	FamilyNm       string `json:"familyName"`
	FarmSpftDesc   string `json:"farmFeatureDescription"`
	GenusKorNm     string `json:"genusKoreanName"`
	GenusNm        string `json:"genusName"`
	GrwEvrntDesc   string `json:"growthEnvironmentDescription"`
	InductionDesc  string `json:"inductionDescription"`
	LastUpdtDtm    string `json:"lastUpdateDateTime"`
	NotRcmmGnrlNm  string `json:"notRecommendedGeneralName"`
	Note           string `json:"note"`
	OrplcNm        string `json:"originPlaceName"`
	OsDstrb        string `json:"overseasDistribution"`
	PlantGnrlNm    string `json:"plantGeneralName"`
	PlantPilbkNo   string `json:"plantPictorialBookNumber"`
	PlantSpecsScnm string `json:"plantSpeciesScientificName"`
	PrtcPlnDesc    string `json:"protectionPlanDescription"`
	RrngGubun      string `json:"growthClassification"`
	RrngType       string `json:"growthType"`
	Shpe           string `json:"shape"`
	SmlrPlntDesc   string `json:"similarPlantDescription"`
	Spft           string `json:"feature"`
	UseMthdDesc    string `json:"useMethodDescription"`
	WoodDesc       string `json:"woodDescription"`
}

type plantPilbkInfoHandler struct {
	useCase inbound.PlantPilbkInfoUseCase
}

func addPlantPilbkInfoTool(server *mcp.Server, useCase inbound.PlantPilbkInfoUseCase) {
	handler := plantPilbkInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantPilbkInfoToolName,
		Description: "산림청 국립수목원 식물도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h plantPilbkInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantPilbkInfoInput) (*mcp.CallToolResult, plantPilbkInfoOutput, error) {
	result, err := h.useCase.PlantPilbkInfo(ctx, application.PlantPilbkInfoQuery{
		ReqPlantPilbkNo: input.ReqPlantPilbkNo,
	})
	if err != nil {
		return nil, plantPilbkInfoOutput{}, err
	}

	return nil, plantPilbkInfoOutput{
		APGFamilyKorNm: result.APGFamilyKorNm,
		APGFamilyNm:    result.APGFamilyNm,
		BfofMthod:      result.BfofMthod,
		BrdMthdDesc:    result.BrdMthdDesc,
		BugInfo:        result.BugInfo,
		Dstrb:          result.Dstrb,
		EngNm:          result.EngNm,
		FamilyKorNm:    result.FamilyKorNm,
		FamilyNm:       result.FamilyNm,
		FarmSpftDesc:   result.FarmSpftDesc,
		GenusKorNm:     result.GenusKorNm,
		GenusNm:        result.GenusNm,
		GrwEvrntDesc:   result.GrwEvrntDesc,
		InductionDesc:  result.InductionDesc,
		LastUpdtDtm:    result.LastUpdtDtm,
		NotRcmmGnrlNm:  result.NotRcmmGnrlNm,
		Note:           result.Note,
		OrplcNm:        result.OrplcNm,
		OsDstrb:        result.OsDstrb,
		PlantGnrlNm:    result.PlantGnrlNm,
		PlantPilbkNo:   result.PlantPilbkNo,
		PlantSpecsScnm: result.PlantSpecsScnm,
		PrtcPlnDesc:    result.PrtcPlnDesc,
		RrngGubun:      result.RrngGubun,
		RrngType:       result.RrngType,
		Shpe:           result.Shpe,
		SmlrPlntDesc:   result.SmlrPlntDesc,
		Spft:           result.Spft,
		UseMthdDesc:    result.UseMthdDesc,
		WoodDesc:       result.WoodDesc,
	}, nil
}
