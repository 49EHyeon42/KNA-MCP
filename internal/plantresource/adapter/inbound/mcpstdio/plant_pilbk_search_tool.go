package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantPilbkSearchToolName = "plant_resource_plant_pictorial_book_search"

type plantPilbkSearchInput struct {
	PageNo       int    `json:"pageNumber" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numberOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"requestSearchWord,omitempty" jsonschema:"식물 검색어"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type plantPilbkSearchOutput struct {
	Items      []plantPilbkSearchItem `json:"items"`
	NumOfRows  int                    `json:"numberOfRows"`
	PageNo     int                    `json:"pageNumber"`
	TotalCount int                    `json:"totalCount"`
}

type plantPilbkSearchItem struct {
	APGFamilyKorNm string `json:"apgFamilyKoreanName"`
	APGFamilyNm    string `json:"apgFamilyName"`
	FamilyKorNm    string `json:"familyKoreanName"`
	FamilyNm       string `json:"familyName"`
	GenusKorNm     string `json:"genusKoreanName"`
	GenusNm        string `json:"genusName"`
	LastUpdtDtm    string `json:"lastUpdateDateTime"`
	NotRcmmGnrlNm  string `json:"notRecommendedGeneralName"`
	PlantGnrlNm    string `json:"plantGeneralName"`
	PlantPilbkNo   string `json:"plantPictorialBookNumber"`
	PlantSpecsScnm string `json:"plantSpeciesScientificName"`
}

type plantPilbkSearchHandler struct {
	useCase inbound.PlantPilbkSearchUseCase
}

func addPlantPilbkSearchTool(server *mcp.Server, useCase inbound.PlantPilbkSearchUseCase) {
	handler := plantPilbkSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantPilbkSearchToolName,
		Description: "산림청 국립수목원 식물도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantPilbkSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantPilbkSearchInput) (*mcp.CallToolResult, plantPilbkSearchOutput, error) {
	result, err := h.useCase.PlantPilbkSearch(ctx, application.PlantPilbkSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantPilbkSearchOutput{}, err
	}

	items := make([]plantPilbkSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantPilbkSearchItem{
			APGFamilyKorNm: item.APGFamilyKorNm,
			APGFamilyNm:    item.APGFamilyNm,
			FamilyKorNm:    item.FamilyKorNm,
			FamilyNm:       item.FamilyNm,
			GenusKorNm:     item.GenusKorNm,
			GenusNm:        item.GenusNm,
			LastUpdtDtm:    item.LastUpdtDtm,
			NotRcmmGnrlNm:  item.NotRcmmGnrlNm,
			PlantGnrlNm:    item.PlantGnrlNm,
			PlantPilbkNo:   item.PlantPilbkNo,
			PlantSpecsScnm: item.PlantSpecsScnm,
		}
	}

	return nil, plantPilbkSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
