package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
)

type insectPilbkSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 곤충의 국명 또는 학명"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type insectPilbkSearchOutput struct {
	Items      []insectPilbkSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                     `json:"numOfRows" jsonschema:"한 페이지당 건 수"`
	PageNo     int                     `json:"pageNo" jsonschema:"페이지 번호"`
	TotalCount int                     `json:"totalCount" jsonschema:"전체 건 수"`
}

type insectPilbkSearchItem struct {
	FamilyKorNm    string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm       string `json:"familyNm" jsonschema:"과명"`
	GenusKorNm     string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm        string `json:"genusNm" jsonschema:"속명"`
	InsctGnrlNm    string `json:"insctGnrlNm" jsonschema:"국명(곤충명)"`
	InsctPilbkNo   string `json:"insctPilbkNo" jsonschema:"곤충도감번호"`
	InsctSpecsScnm string `json:"insctSpecsScnm" jsonschema:"학명"`
	LastUpdtDtm    string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
}

type insectPilbkSearchHandler struct {
	useCase inbound.InsectPilbkSearchUseCase
}

func addInsectPilbkSearchTool(server *mcp.Server, useCase inbound.InsectPilbkSearchUseCase) {
	handler := insectPilbkSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insect_resource_insect_pilbk_search",
		Description: "산림청 국립수목원 곤충도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h insectPilbkSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input insectPilbkSearchInput) (*mcp.CallToolResult, insectPilbkSearchOutput, error) {
	result, err := h.useCase.InsectPilbkSearch(ctx, application.InsectPilbkSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, insectPilbkSearchOutput{}, err
	}

	items := make([]insectPilbkSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = insectPilbkSearchItem{
			FamilyKorNm:    item.FamilyKorNm,
			FamilyNm:       item.FamilyNm,
			GenusKorNm:     item.GenusKorNm,
			GenusNm:        item.GenusNm,
			InsctGnrlNm:    item.InsctGnrlNm,
			InsctPilbkNo:   item.InsctPilbkNo,
			InsctSpecsScnm: item.InsctSpecsScnm,
			LastUpdtDtm:    item.LastUpdtDtm,
		}
	}

	return nil, insectPilbkSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
