package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantFolkAreaListToolName = "plant_resource_plant_folk_area_list"

type plantFolkAreaListInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	FlpltID   string `json:"flpltId" jsonschema:"민속식물 목록 검색 결과의 민속식물 ID"`
}

type plantFolkAreaListOutput struct {
	Items      []plantFolkAreaListItem `json:"items"`
	NumOfRows  int                     `json:"numOfRows"`
	PageNo     int                     `json:"pageNo"`
	TotalCount int                     `json:"totalCount"`
}

type plantFolkAreaListItem struct {
	FlcstPlantExmnnAraTpcdNm string `json:"flcstPlantExmnnAraTpcdNm"`
	FlcstPlantLcltDscrt      string `json:"flcstPlantLcltDscrt"`
	FlcstPlantPrpseDscrt     string `json:"flcstPlantPrpseDscrt"`
	FlpltID                  string `json:"flpltId"`
	PlantBrdgFomTpcdNm       string `json:"plantBrdgFomTpcdNm"`
	PlantGnrlNm              string `json:"plantGnrlNm"`
	PlantSpecsScnm           string `json:"plantSpecsScnm"`
}

type plantFolkAreaListHandler struct {
	useCase inbound.PlantFolkAreaListUseCase
}

func addPlantFolkAreaListTool(server *mcp.Server, useCase inbound.PlantFolkAreaListUseCase) {
	handler := plantFolkAreaListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantFolkAreaListToolName,
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
