package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantPilbkInfoInput struct {
	ReqPlantPilbkNo string `json:"reqPlantPilbkNo" jsonschema:"검색할 식물도감번호 (plantPilbkSearch 결과의 plantPilbkNo)"`
}

type plantPilbkInfoOutput struct {
	APGFamilyKorNm string `json:"apgFamilyKorNm" jsonschema:"APG과국명"`
	APGFamilyNm    string `json:"apgFamilyNm" jsonschema:"APG과명"`
	BfofMthod      string `json:"bfofMthod" jsonschema:"방제방법"`
	BrdMthdDesc    string `json:"brdMthdDesc" jsonschema:"번식방법"`
	BugInfo        string `json:"bugInfo" jsonschema:"병충해정보"`
	Dstrb          string `json:"dstrb" jsonschema:"분포"`
	EngNm          string `json:"engNm" jsonschema:"영문명"`
	FamilyKorNm    string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm       string `json:"familyNm" jsonschema:"과명"`
	FarmSpftDesc   string `json:"farmSpftDesc" jsonschema:"재배특성"`
	GenusKorNm     string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm        string `json:"genusNm" jsonschema:"속명"`
	GrwEvrntDesc   string `json:"grwEvrntDesc" jsonschema:"생육환경"`
	InductionDesc  string `json:"inductionDesc" jsonschema:"도입여부"`
	LastUpdtDtm    string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	NotRcmmGnrlNm  string `json:"notRcmmGnrlNm" jsonschema:"비추천국명"`
	Note           string `json:"note" jsonschema:"비고"`
	OrplcNm        string `json:"orplcNm" jsonschema:"원산지"`
	OsDstrb        string `json:"osDstrb" jsonschema:"해외분포"`
	PlantGnrlNm    string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantPilbkNo   string `json:"plantPilbkNo" jsonschema:"식물도감번호"`
	PlantSpecsScnm string `json:"plantSpecsScnm" jsonschema:"학명"`
	PrtcPlnDesc    string `json:"prtcPlnDesc" jsonschema:"보호방안"`
	RrngGubun      string `json:"rrngGubun" jsonschema:"생육상 구분"`
	RrngType       string `json:"rrngType" jsonschema:"생육형"`
	Shpe           string `json:"shpe" jsonschema:"형태"`
	SmlrPlntDesc   string `json:"smlrPlntDesc" jsonschema:"유사종"`
	Spft           string `json:"spft" jsonschema:"특징"`
	UseMthdDesc    string `json:"useMthdDesc" jsonschema:"이용방안"`
	WoodDesc       string `json:"woodDesc" jsonschema:"목재"`
}

type plantPilbkInfoHandler struct {
	useCase inbound.PlantPilbkInfoUseCase
}

func addPlantPilbkInfoTool(server *mcp.Server, useCase inbound.PlantPilbkInfoUseCase) {
	handler := plantPilbkInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_pilbk_info",
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
