package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantSmplSearchToolName = "plant_resource_plant_specimen_search"

type plantSmplSearchInput struct {
	PageNo       int    `json:"pageNumber" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numberOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"requestSearchWord,omitempty" jsonschema:"식물표본의 국명 또는 학명 검색어"`
}

type plantSmplSearchOutput struct {
	Items      []plantSmplSearchItem `json:"items"`
	NumOfRows  int                   `json:"numberOfRows"`
	PageNo     int                   `json:"pageNumber"`
	TotalCount int                   `json:"totalCount"`
}

type plantSmplSearchItem struct {
	Cnt            int    `json:"count"`
	FamilyKorNm    string `json:"familyKoreanName"`
	FamilyNm       string `json:"familyName"`
	PlantGnrlNm    string `json:"plantGeneralName"`
	PlantSpecsID   string `json:"plantSpeciesId"`
	PlantSpecsScnm string `json:"plantSpeciesScientificName"`
}

type plantSmplSearchHandler struct {
	useCase inbound.PlantSmplSearchUseCase
}

func addPlantSmplSearchTool(server *mcp.Server, useCase inbound.PlantSmplSearchUseCase) {
	handler := plantSmplSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSmplSearchToolName,
		Description: "산림청 국립수목원 식물표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantSmplSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSmplSearchInput) (*mcp.CallToolResult, plantSmplSearchOutput, error) {
	result, err := h.useCase.PlantSmplSearch(ctx, application.PlantSmplSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantSmplSearchOutput{}, err
	}

	items := make([]plantSmplSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSmplSearchItem{
			Cnt:            item.Cnt,
			FamilyKorNm:    item.FamilyKorNm,
			FamilyNm:       item.FamilyNm,
			PlantGnrlNm:    item.PlantGnrlNm,
			PlantSpecsID:   item.PlantSpecsID,
			PlantSpecsScnm: item.PlantSpecsScnm,
		}
	}

	return nil, plantSmplSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
