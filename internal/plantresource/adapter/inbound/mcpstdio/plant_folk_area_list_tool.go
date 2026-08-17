package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantFolkAreaListInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	FlpltID   string `json:"flpltId" jsonschema:"검색할 민속식물ID (plantFolkSearch 결과의 flpltId)"`
}

type plantFolkAreaListOutput struct {
	Items      []plantFolkAreaListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                     `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                     `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                     `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantFolkAreaListItem struct {
	FlcstPlantExmnnAraTpcdNm string `json:"flcstPlantExmnnAraTpcdNm" jsonschema:"지역명"`
	FlcstPlantLcltDscrt      string `json:"flcstPlantLcltDscrt" jsonschema:"지방특성설명"`
	FlcstPlantPrpseDscrt     string `json:"flcstPlantPrpseDscrt" jsonschema:"지방별 용도설명"`
	FlpltID                  string `json:"flpltId" jsonschema:"민속식물ID"`
	PlantBrdgFomTpcdNm       string `json:"plantBrdgFomTpcdNm" jsonschema:"식물번식형태"`
	PlantGnrlNm              string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantSpecsScnm           string `json:"plantSpecsScnm" jsonschema:"학명"`
}

type plantFolkAreaListHandler struct {
	useCase inbound.PlantFolkAreaListUseCase
}

func addPlantFolkAreaListTool(server *mcp.Server, useCase inbound.PlantFolkAreaListUseCase) {
	handler := plantFolkAreaListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_folk_area_list",
		Description: "산림청 국립수목원 민속식물 지방별 이용정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantFolkAreaListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantFolkAreaListInput) (*mcp.CallToolResult, plantFolkAreaListOutput, error) {
	result, err := h.useCase.PlantFolkAreaList(ctx, application.PlantFolkAreaListQuery{
		PageNo:    input.PageNo,
		NumOfRows: input.NumOfRows,
		FlpltID:   input.FlpltID,
	})
	if err != nil {
		return nil, plantFolkAreaListOutput{}, err
	}

	items := make([]plantFolkAreaListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantFolkAreaListItem{
			FlcstPlantExmnnAraTpcdNm: item.FlcstPlantExmnnAraTpcdNm,
			FlcstPlantLcltDscrt:      item.FlcstPlantLcltDscrt,
			FlcstPlantPrpseDscrt:     item.FlcstPlantPrpseDscrt,
			FlpltID:                  item.FlpltID,
			PlantBrdgFomTpcdNm:       item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:              item.PlantGnrlNm,
			PlantSpecsScnm:           item.PlantSpecsScnm,
		}
	}

	return nil, plantFolkAreaListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
