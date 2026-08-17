package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
)

type insectSmplSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 곤충의 학명 또는 국명 (대소문자를 구분하지 않는 부분 문자열 검색)"`
}

type insectSmplSearchOutput struct {
	Items      []insectSmplSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                    `json:"numOfRows" jsonschema:"한 페이지 당 건 수"`
	PageNo     int                    `json:"pageNo" jsonschema:"페이지 번호"`
	TotalCount int                    `json:"totalCount" jsonschema:"전체 건 수"`
}

type insectSmplSearchItem struct {
	Cnt              string `json:"cnt" jsonschema:"표본수"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	GenusKorNm       string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm          string `json:"genusNm" jsonschema:"속명"`
	InsctGnrlNm      string `json:"insctGnrlNm" jsonschema:"국명(곤충명)"`
	InsctSpecsID     string `json:"insctSpecsId" jsonschema:"곤충종ID"`
	InsctSpecsScnm   string `json:"insctSpecsScnm" jsonschema:"학명"`
	SubFamilyKorNm   string `json:"subFamilyKorNm" jsonschema:"아과국명"`
	SubFamilyNm      string `json:"subFamilyNm" jsonschema:"아과명"`
	SuperFamilyKorNm string `json:"superFamilyKorNm" jsonschema:"상과국명"`
	SuperFamilyNm    string `json:"superFamilyNm" jsonschema:"상과명"`
}

type insectSmplSearchHandler struct {
	useCase inbound.InsectSmplSearchUseCase
}

func addInsectSmplSearchTool(server *mcp.Server, useCase inbound.InsectSmplSearchUseCase) {
	handler := insectSmplSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insect_resource_insect_smpl_search",
		Description: "산림청 국립수목원 곤충표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h insectSmplSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input insectSmplSearchInput) (*mcp.CallToolResult, insectSmplSearchOutput, error) {
	result, err := h.useCase.InsectSmplSearch(ctx, application.InsectSmplSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, insectSmplSearchOutput{}, err
	}

	items := make([]insectSmplSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = insectSmplSearchItem{
			Cnt:              item.Cnt,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			InsctGnrlNm:      item.InsctGnrlNm,
			InsctSpecsID:     item.InsctSpecsID,
			InsctSpecsScnm:   item.InsctSpecsScnm,
			SubFamilyKorNm:   item.SubFamilyKorNm,
			SubFamilyNm:      item.SubFamilyNm,
			SuperFamilyKorNm: item.SuperFamilyKorNm,
			SuperFamilyNm:    item.SuperFamilyNm,
		}
	}

	return nil, insectSmplSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
