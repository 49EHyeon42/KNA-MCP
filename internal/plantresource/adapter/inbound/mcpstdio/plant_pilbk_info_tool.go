package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantPilbkInfoToolName = "plant_resource_plant_pilbk_info"

type plantPilbkInfoInput struct {
	ReqPlantPilbkNo string `json:"reqPlantPilbkNo" jsonschema:"식물도감 목록 검색 결과의 식물도감번호"`
}

type plantPilbkInfoOutput struct {
	APGFamilyKorNm string `json:"apgFamilyKorNm"`
	APGFamilyNm    string `json:"apgFamilyNm"`
	BfofMthod      string `json:"bfofMthod"`
	BrdMthdDesc    string `json:"brdMthdDesc"`
	BugInfo        string `json:"bugInfo"`
	Dstrb          string `json:"dstrb"`
	EngNm          string `json:"engNm"`
	FamilyKorNm    string `json:"familyKorNm"`
	FamilyNm       string `json:"familyNm"`
	FarmSpftDesc   string `json:"farmSpftDesc"`
	GenusKorNm     string `json:"genusKorNm"`
	GenusNm        string `json:"genusNm"`
	GrwEvrntDesc   string `json:"grwEvrntDesc"`
	InductionDesc  string `json:"inductionDesc"`
	LastUpdtDtm    string `json:"lastUpdtDtm"`
	NotRcmmGnrlNm  string `json:"notRcmmGnrlNm"`
	Note           string `json:"note"`
	OrplcNm        string `json:"orplcNm"`
	OsDstrb        string `json:"osDstrb"`
	PlantGnrlNm    string `json:"plantGnrlNm"`
	PlantPilbkNo   string `json:"plantPilbkNo"`
	PlantSpecsScnm string `json:"plantSpecsScnm"`
	PrtcPlnDesc    string `json:"prtcPlnDesc"`
	RrngGubun      string `json:"rrngGubun"`
	RrngType       string `json:"rrngType"`
	Shpe           string `json:"shpe"`
	SmlrPlntDesc   string `json:"smlrPlntDesc"`
	Spft           string `json:"spft"`
	UseMthdDesc    string `json:"useMthdDesc"`
	WoodDesc       string `json:"woodDesc"`
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
