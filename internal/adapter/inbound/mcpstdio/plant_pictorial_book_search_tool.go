package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
	"github.com/49EHyeon42/KNA-MCP/internal/application/port/inbound"
)

const plantResourcePlantPictorialBookSearchToolName = "plant_resource_plant_pictorial_book_search"

type plantPictorialBookSearchInput struct {
	PageNumber        int    `json:"pageNumber" jsonschema:"페이지 번호(1 이상)"`
	NumberOfRows      int    `json:"numberOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	RequestSearchWord string `json:"requestSearchWord,omitempty" jsonschema:"식물 검색어"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type plantPictorialBookSearchOutput struct {
	Items        []plantPictorialBookSearchItem `json:"items"`
	NumberOfRows int                            `json:"numberOfRows"`
	PageNumber   int                            `json:"pageNumber"`
	TotalCount   int                            `json:"totalCount"`
}

type plantPictorialBookSearchItem struct {
	APGFamilyKoreanName        string `json:"apgFamilyKoreanName"`
	APGFamilyName              string `json:"apgFamilyName"`
	FamilyKoreanName           string `json:"familyKoreanName"`
	FamilyName                 string `json:"familyName"`
	GenusKoreanName            string `json:"genusKoreanName"`
	GenusName                  string `json:"genusName"`
	LastUpdateDateTime         string `json:"lastUpdateDateTime"`
	NotRecommendedGeneralName  string `json:"notRecommendedGeneralName"`
	PlantGeneralName           string `json:"plantGeneralName"`
	PlantPictorialBookNumber   string `json:"plantPictorialBookNumber"`
	PlantSpeciesScientificName string `json:"plantSpeciesScientificName"`
}

type plantPictorialBookSearchHandler struct {
	useCase inbound.PlantPictorialBookSearchUseCase
}

func addPlantPictorialBookSearchTool(server *mcp.Server, useCase inbound.PlantPictorialBookSearchUseCase) {
	handler := plantPictorialBookSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantPictorialBookSearchToolName,
		Description: "산림청 국립수목원 식물도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantPictorialBookSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantPictorialBookSearchInput) (*mcp.CallToolResult, plantPictorialBookSearchOutput, error) {
	result, err := h.useCase.PlantPictorialBookSearch(ctx, application.PlantPictorialBookSearchQuery{
		PageNumber:        input.PageNumber,
		NumberOfRows:      input.NumberOfRows,
		RequestSearchWord: input.RequestSearchWord,
	})
	if err != nil {
		return nil, plantPictorialBookSearchOutput{}, err
	}

	items := make([]plantPictorialBookSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantPictorialBookSearchItem{
			APGFamilyKoreanName:        item.APGFamilyKoreanName,
			APGFamilyName:              item.APGFamilyName,
			FamilyKoreanName:           item.FamilyKoreanName,
			FamilyName:                 item.FamilyName,
			GenusKoreanName:            item.GenusKoreanName,
			GenusName:                  item.GenusName,
			LastUpdateDateTime:         item.LastUpdateDateTime,
			NotRecommendedGeneralName:  item.NotRecommendedGeneralName,
			PlantGeneralName:           item.PlantGeneralName,
			PlantPictorialBookNumber:   item.PlantPictorialBookNumber,
			PlantSpeciesScientificName: item.PlantSpeciesScientificName,
		}
	}

	return nil, plantPictorialBookSearchOutput{
		Items:        items,
		NumberOfRows: result.NumberOfRows,
		PageNumber:   result.PageNumber,
		TotalCount:   result.TotalCount,
	}, nil
}
