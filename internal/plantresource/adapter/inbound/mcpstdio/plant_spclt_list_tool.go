package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantSpcltListToolName = "plant_resource_plant_spclt_list"

type plantSpcltListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"특산식물의 국명 또는 학명 검색어"`
}

type plantSpcltListOutput struct {
	Items      []plantSpcltListItem `json:"items"`
	NumOfRows  int                  `json:"numOfRows"`
	PageNo     int                  `json:"pageNo"`
	TotalCount int                  `json:"totalCount"`
}

type plantSpcltListItem struct {
	AgpFamilyKorNm     string `json:"agpFamilyKorNm"`
	AgpFamilyNm        string `json:"agpFamilyNm"`
	ExtrmCrssScls1Yn   string `json:"extrmCrssScls1Yn"`
	ExtrmCrssScls2Yn   string `json:"extrmCrssScls2Yn"`
	FamilyKorNm        string `json:"familyKorNm"`
	FamilyNm           string `json:"familyNm"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm"`
	PlantGnrlNm        string `json:"plantGnrlNm"`
	PlantSpecsScnm     string `json:"plantSpecsScnm"`
}

type plantSpcltListHandler struct {
	useCase inbound.PlantSpcltListUseCase
}

func addPlantSpcltListTool(server *mcp.Server, useCase inbound.PlantSpcltListUseCase) {
	handler := plantSpcltListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSpcltListToolName,
		Description: "산림청 국립수목원 특산식물 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantSpcltListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSpcltListInput) (*mcp.CallToolResult, plantSpcltListOutput, error) {
	result, err := h.useCase.PlantSpcltList(ctx, application.PlantSpcltListQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantSpcltListOutput{}, err
	}

	items := make([]plantSpcltListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSpcltListItem{
			AgpFamilyKorNm:     item.AgpFamilyKorNm,
			AgpFamilyNm:        item.AgpFamilyNm,
			ExtrmCrssScls1Yn:   item.ExtrmCrssScls1Yn,
			ExtrmCrssScls2Yn:   item.ExtrmCrssScls2Yn,
			FamilyKorNm:        item.FamilyKorNm,
			FamilyNm:           item.FamilyNm,
			PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:        item.PlantGnrlNm,
			PlantSpecsScnm:     item.PlantSpecsScnm,
		}
	}

	return nil, plantSpcltListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
