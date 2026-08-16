package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantFolkSearchToolName = "plant_resource_plant_folk_search"

type plantFolkSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"민속식물의 국명 또는 학명 검색어"`
}

type plantFolkSearchOutput struct {
	Items      []plantFolkSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 검색결과 수"`
}

type plantFolkSearchItem struct {
	FlcstPlantIdntfDscrt string `json:"flcstPlantIdntfDscrt" jsonschema:"식별설명"`
	FlpltID              string `json:"flpltId" jsonschema:"민속식물ID"`
	PlantBrdgFomTpcdNm   string `json:"plantBrdgFomTpcdNm" jsonschema:"식물번식형태"`
	PlantGnrlNm          string `json:"plantGnrlNm" jsonschema:"국명(식물명"`
	PlantSpecsScnm       string `json:"plantSpecsScnm" jsonschema:"학명"`
	Ptnt                 string `json:"ptnt" jsonschema:"특허정보"`
}

type plantFolkSearchHandler struct {
	useCase inbound.PlantFolkSearchUseCase
}

func addPlantFolkSearchTool(server *mcp.Server, useCase inbound.PlantFolkSearchUseCase) {
	handler := plantFolkSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantFolkSearchToolName,
		Description: "산림청 국립수목원 민속식물 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantFolkSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantFolkSearchInput) (*mcp.CallToolResult, plantFolkSearchOutput, error) {
	result, err := h.useCase.PlantFolkSearch(ctx, application.PlantFolkSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantFolkSearchOutput{}, err
	}

	items := make([]plantFolkSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantFolkSearchItem{
			FlcstPlantIdntfDscrt: item.FlcstPlantIdntfDscrt,
			FlpltID:              item.FlpltID,
			PlantBrdgFomTpcdNm:   item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:          item.PlantGnrlNm,
			PlantSpecsScnm:       item.PlantSpecsScnm,
			Ptnt:                 item.Ptnt,
		}
	}

	return nil, plantFolkSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
