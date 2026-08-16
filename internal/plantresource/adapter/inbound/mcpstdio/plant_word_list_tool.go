package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantWordListToolName = "plant_resource_plant_word_list"

type plantWordListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 한글 식물 용어명"`
}

type plantWordListOutput struct {
	Items      []plantWordListItem `json:"items"`
	NumOfRows  int                 `json:"numOfRows"`
	PageNo     int                 `json:"pageNo"`
	TotalCount int                 `json:"totalCount"`
}

type plantWordListItem struct {
	EnglsWrdNm string `json:"englsWrdNm"`
	KrnWrdNm   string `json:"krnWrdNm"`
	PrfcnWrdNm string `json:"prfcnWrdNm"`
	Wrddscrt   string `json:"wrddscrt"`
}

type plantWordListHandler struct {
	useCase inbound.PlantWordListUseCase
}

func addPlantWordListTool(server *mcp.Server, useCase inbound.PlantWordListUseCase) {
	handler := plantWordListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantWordListToolName,
		Description: "산림청 국립수목원 식물 용어사전 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantWordListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantWordListInput) (*mcp.CallToolResult, plantWordListOutput, error) {
	result, err := h.useCase.PlantWordList(ctx, application.PlantWordListQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantWordListOutput{}, err
	}

	items := make([]plantWordListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantWordListItem{
			EnglsWrdNm: item.EnglsWrdNm,
			KrnWrdNm:   item.KrnWrdNm,
			PrfcnWrdNm: item.PrfcnWrdNm,
			Wrddscrt:   item.Wrddscrt,
		}
	}

	return nil, plantWordListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
