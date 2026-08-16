package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
	"github.com/49EHyeon42/KNA-MCP/internal/application/port/inbound"
)

const plantResourcePlantSampleSearchToolName = "plant_resource_plant_sample_search"

type plantSampleSearchInput struct {
	PageNumber        int    `json:"pageNumber" jsonschema:"페이지 번호(1 이상)"`
	NumberOfRows      int    `json:"numberOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	RequestSearchWord string `json:"requestSearchWord,omitempty" jsonschema:"식물표본의 국명 또는 학명 검색어"`
}

type plantSampleSearchOutput struct {
	Items        []plantSampleSearchItem `json:"items"`
	NumberOfRows int                     `json:"numberOfRows"`
	PageNumber   int                     `json:"pageNumber"`
	TotalCount   int                     `json:"totalCount"`
}

type plantSampleSearchItem struct {
	Count                      int    `json:"count"`
	FamilyKoreanName           string `json:"familyKoreanName"`
	FamilyName                 string `json:"familyName"`
	PlantGeneralName           string `json:"plantGeneralName"`
	PlantSpeciesID             string `json:"plantSpeciesId"`
	PlantSpeciesScientificName string `json:"plantSpeciesScientificName"`
}

type plantSampleSearchHandler struct {
	useCase inbound.PlantSampleSearchUseCase
}

func addPlantSampleSearchTool(server *mcp.Server, useCase inbound.PlantSampleSearchUseCase) {
	handler := plantSampleSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSampleSearchToolName,
		Description: "산림청 국립수목원 식물표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantSampleSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSampleSearchInput) (*mcp.CallToolResult, plantSampleSearchOutput, error) {
	result, err := h.useCase.PlantSampleSearch(ctx, application.PlantSampleSearchQuery{
		PageNumber:        input.PageNumber,
		NumberOfRows:      input.NumberOfRows,
		RequestSearchWord: input.RequestSearchWord,
	})
	if err != nil {
		return nil, plantSampleSearchOutput{}, err
	}

	items := make([]plantSampleSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSampleSearchItem{
			Count:                      item.Count,
			FamilyKoreanName:           item.FamilyKoreanName,
			FamilyName:                 item.FamilyName,
			PlantGeneralName:           item.PlantGeneralName,
			PlantSpeciesID:             item.PlantSpeciesID,
			PlantSpeciesScientificName: item.PlantSpeciesScientificName,
		}
	}

	return nil, plantSampleSearchOutput{
		Items:        items,
		NumberOfRows: result.NumberOfRows,
		PageNumber:   result.PageNumber,
		TotalCount:   result.TotalCount,
	}, nil
}
