package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantRareListToolName = "plant_resource_plant_rare_list"

type plantRareListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"적색식물의 국명 또는 학명 검색어"`
}

type plantRareListOutput struct {
	Items      []plantRareListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                 `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                 `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                 `json:"totalCount" jsonschema:"전체 결과 수"`
}

type plantRareListItem struct {
	AgpFamilyNm      string `json:"agpFamilyNm" jsonschema:"APG과국명"`
	APGFamilyKorNm   string `json:"apgFamilyKorNm" jsonschema:"APG과명"`
	ExtrmCrssScls1Yn string `json:"extrmCrssScls1Yn" jsonschema:"멸종위기종1급 여부"`
	ExtrmCrssScls2Yn string `json:"extrmCrssScls2Yn" jsonschema:"멸종위기종2급 여부"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	PlantGnrlNm      string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantSpecsScnm   string `json:"plantSpecsScnm" jsonschema:"학명"`
	RareTpcdNm       string `json:"rareTpcdNm" jsonschema:"IUCN 적색식물 등급"`
}

type plantRareListHandler struct {
	useCase inbound.PlantRareListUseCase
}

func addPlantRareListTool(server *mcp.Server, useCase inbound.PlantRareListUseCase) {
	handler := plantRareListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantRareListToolName,
		Description: "산림청 국립수목원 적색식물 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantRareListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantRareListInput) (*mcp.CallToolResult, plantRareListOutput, error) {
	result, err := h.useCase.PlantRareList(ctx, application.PlantRareListQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantRareListOutput{}, err
	}

	items := make([]plantRareListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantRareListItem{
			AgpFamilyNm:      item.AgpFamilyNm,
			APGFamilyKorNm:   item.APGFamilyKorNm,
			ExtrmCrssScls1Yn: item.ExtrmCrssScls1Yn,
			ExtrmCrssScls2Yn: item.ExtrmCrssScls2Yn,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			PlantGnrlNm:      item.PlantGnrlNm,
			PlantSpecsScnm:   item.PlantSpecsScnm,
			RareTpcdNm:       item.RareTpcdNm,
		}
	}

	return nil, plantRareListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
