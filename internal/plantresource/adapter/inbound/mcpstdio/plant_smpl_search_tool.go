package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantSmplSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 식물표본의 국명 또는 학명"`
}

type plantSmplSearchOutput struct {
	Items      []plantSmplSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantSmplSearchItem struct {
	Cnt            int    `json:"cnt" jsonschema:"표본수"`
	FamilyKorNm    string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm       string `json:"familyNm" jsonschema:"과명"`
	PlantGnrlNm    string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantSpecsID   string `json:"plantSpecsId" jsonschema:"식물 종ID"`
	PlantSpecsScnm string `json:"plantSpecsScnm" jsonschema:"학명"`
}

type plantSmplSearchHandler struct {
	useCase inbound.PlantSmplSearchUseCase
}

func addPlantSmplSearchTool(server *mcp.Server, useCase inbound.PlantSmplSearchUseCase) {
	handler := plantSmplSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_smpl_search",
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
