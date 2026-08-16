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
	Items      []plantRareListItem `json:"items"`
	NumOfRows  int                 `json:"numOfRows"`
	PageNo     int                 `json:"pageNo"`
	TotalCount int                 `json:"totalCount"`
}

type plantRareListItem struct {
	AgpFamilyNm      string `json:"agpFamilyNm"`
	APGFamilyKorNm   string `json:"apgFamilyKorNm"`
	ExtrmCrssScls1Yn string `json:"extrmCrssScls1Yn"`
	ExtrmCrssScls2Yn string `json:"extrmCrssScls2Yn"`
	FamilyKorNm      string `json:"familyKorNm"`
	FamilyNm         string `json:"familyNm"`
	PlantGnrlNm      string `json:"plantGnrlNm"`
	PlantSpecsScnm   string `json:"plantSpecsScnm"`
	RareTpcdNm       string `json:"rareTpcdNm"`
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
