package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"kna-mcp/internal/application"
	"kna-mcp/internal/application/port/inbound"
)

const plantResourcePlantSpecimenSearchToolName = "plant_resource_plant_specimen_search"

type plantSpecimenSearchInput struct {
	PageNumber        int    `json:"pageNumber" jsonschema:"페이지 번호(1 이상)"`
	NumberOfRows      int    `json:"numberOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	RequestSearchWord string `json:"requestSearchWord,omitempty" jsonschema:"식물표본의 국명 또는 학명 검색어"`
}

type plantSpecimenSearchOutput struct {
	Items        []plantSpecimenSearchItem `json:"items"`
	NumberOfRows int                       `json:"numberOfRows"`
	PageNumber   int                       `json:"pageNumber"`
	TotalCount   int                       `json:"totalCount"`
}

type plantSpecimenSearchItem struct {
	Count                      int    `json:"count"`
	FamilyKoreanName           string `json:"familyKoreanName"`
	FamilyName                 string `json:"familyName"`
	PlantGeneralName           string `json:"plantGeneralName"`
	PlantSpeciesID             string `json:"plantSpeciesId"`
	PlantSpeciesScientificName string `json:"plantSpeciesScientificName"`
}

type plantSpecimenSearchHandler struct {
	useCase inbound.PlantSpecimenSearchUseCase
}

func addPlantSpecimenSearchTool(server *mcp.Server, useCase inbound.PlantSpecimenSearchUseCase) {
	handler := plantSpecimenSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSpecimenSearchToolName,
		Description: "산림청 국립수목원 식물표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantSpecimenSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSpecimenSearchInput) (*mcp.CallToolResult, plantSpecimenSearchOutput, error) {
	result, err := h.useCase.PlantSpecimenSearch(ctx, application.PlantSpecimenSearchQuery{
		PageNumber:        input.PageNumber,
		NumberOfRows:      input.NumberOfRows,
		RequestSearchWord: input.RequestSearchWord,
	})
	if err != nil {
		return nil, plantSpecimenSearchOutput{}, err
	}

	items := make([]plantSpecimenSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSpecimenSearchItem{
			Count:                      item.Count,
			FamilyKoreanName:           item.FamilyKoreanName,
			FamilyName:                 item.FamilyName,
			PlantGeneralName:           item.PlantGeneralName,
			PlantSpeciesID:             item.PlantSpeciesID,
			PlantSpeciesScientificName: item.PlantSpeciesScientificName,
		}
	}

	return nil, plantSpecimenSearchOutput{
		Items:        items,
		NumberOfRows: result.NumberOfRows,
		PageNumber:   result.PageNumber,
		TotalCount:   result.TotalCount,
	}, nil
}
