package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantSpcltListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"특산식물의 국명 또는 학명 검색어"`
}

type plantSpcltListOutput struct {
	Items      []plantSpcltListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                  `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                  `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                  `json:"totalCount" jsonschema:"전체 결과 수"`
}

type plantSpcltListItem struct {
	AgpFamilyKorNm     string `json:"agpFamilyKorNm" jsonschema:"APG과국명"`
	AgpFamilyNm        string `json:"agpFamilyNm" jsonschema:"APG과명"`
	ExtrmCrssScls1Yn   string `json:"extrmCrssScls1Yn" jsonschema:"멸종위기식물 1급 여부"`
	ExtrmCrssScls2Yn   string `json:"extrmCrssScls2Yn" jsonschema:"멸종위기식물 2급 여부"`
	FamilyKorNm        string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm           string `json:"familyNm" jsonschema:"과명"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm" jsonschema:"식물번식형태"`
	PlantGnrlNm        string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantSpecsScnm     string `json:"plantSpecsScnm" jsonschema:"학명"`
}

type plantSpcltListHandler struct {
	useCase inbound.PlantSpcltListUseCase
}

func addPlantSpcltListTool(server *mcp.Server, useCase inbound.PlantSpcltListUseCase) {
	handler := plantSpcltListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_spclt_list",
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
