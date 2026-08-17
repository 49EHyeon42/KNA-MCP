package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
)

type insectPrtctListInput struct {
	PageNo    int `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
}

type insectPrtctListOutput struct {
	Items      []insectPrtctListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한페이지 결과수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지 번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 건수"`
}

type insectPrtctListItem struct {
	FamilyKorNm    string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm       string `json:"familyNm" jsonschema:"과명"`
	InsctGnrlNm    string `json:"insctGnrlNm" jsonschema:"국명(곤충명)"`
	InsctPcmtt     string `json:"insctPcmtt" jsonschema:"멸종위기구분"`
	InsctPilbkNo   string `json:"insctPilbkNo" jsonschema:"곤충도감번호"`
	InsctSpecsScnm string `json:"insctSpecsScnm" jsonschema:"학명"`
}

type insectPrtctListHandler struct {
	useCase inbound.InsectPrtctListUseCase
}

func addInsectPrtctListTool(server *mcp.Server, useCase inbound.InsectPrtctListUseCase) {
	handler := insectPrtctListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insect_resource_insect_prtct_list",
		Description: "산림청 국립수목원 멸종위기곤충 목록을 조회합니다.",
	}, handler.handle)
}

func (h insectPrtctListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input insectPrtctListInput) (*mcp.CallToolResult, insectPrtctListOutput, error) {
	result, err := h.useCase.InsectPrtctList(ctx, application.InsectPrtctListQuery{
		PageNo:    input.PageNo,
		NumOfRows: input.NumOfRows,
	})
	if err != nil {
		return nil, insectPrtctListOutput{}, err
	}

	items := make([]insectPrtctListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = insectPrtctListItem{
			FamilyKorNm:    item.FamilyKorNm,
			FamilyNm:       item.FamilyNm,
			InsctGnrlNm:    item.InsctGnrlNm,
			InsctPcmtt:     item.InsctPcmtt,
			InsctPilbkNo:   item.InsctPilbkNo,
			InsctSpecsScnm: item.InsctSpecsScnm,
		}
	}

	return nil, insectPrtctListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
