package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kini/application/port/inbound"
)

type scnmSearchInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqGnrlNm string `json:"reqGnrlNm,omitempty" jsonschema:"검색하려는 곤충국명(곤충명) (부분 문자열 검색)"`
	ReqScnm   string `json:"reqScnm,omitempty" jsonschema:"검색하려는 곤충학명 (대소문자를 구분하지 않는 부분 문자열 검색)"`
	DateFrom  string `json:"dateFrom,omitempty" jsonschema:"최종수정일 이후 정보 (yyyyMMdd)"`
	DateTo    string `json:"dateTo,omitempty" jsonschema:"최종수정일 이전 정보 (yyyyMMdd)"`
}

type scnmSearchOutput struct {
	Items      []scnmSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int              `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int              `json:"pageNo" jsonschema:"페이지 번호"`
	TotalCount int              `json:"totalCount" jsonschema:"전체 결과 수"`
}

type scnmSearchItem struct {
	SuperFalmNm       string `json:"superFalmNm" jsonschema:"곤충 상과명(SuperFamily Name)"`
	ClassKorNm        string `json:"classKorNm" jsonschema:"곤충 강국명"`
	ClassNm           string `json:"classNm" jsonschema:"곤충 강명(Class Name)"`
	FalmKorNm         string `json:"falmKorNm" jsonschema:"곤충 과국명"`
	FalmNm            string `json:"falmNm" jsonschema:"곤충 과명(Family Name)"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"곤충 속국명"`
	GenusNm           string `json:"genusNm" jsonschema:"곤충 속명(Genus Name)"`
	InsctGnrlNm       string `json:"insctGnrlNm" jsonschema:"곤충 추천국명(곤충명)"`
	InsctScnmID       string `json:"insctScnmId" jsonschema:"곤충 학명ID"`
	InsctSpecsScnm    string `json:"insctSpecsScnm" jsonschema:"곤충 학명"`
	LastUpdtDtm       string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	OrdKorNm          string `json:"ordKorNm" jsonschema:"곤충 목국명"`
	OrdNm             string `json:"ordNm" jsonschema:"곤충 목명(Order Name)"`
	StpltScnmRltnCdNm string `json:"stpltScnmRltnCdNm" jsonschema:"곤충 학명 정명/이명 구분"`
	SubFalmKorNm      string `json:"subFalmKorNm" jsonschema:"곤충 아과국명"`
	SubFalmNm         string `json:"subFalmNm" jsonschema:"곤충 아과명(SubFamily Name)"`
	SuperFalmKorNm    string `json:"superFalmKorNm" jsonschema:"곤충 상과국명"`
}

type scnmSearchHandler struct {
	useCase inbound.ScnmSearchUseCase
}

func addScnmSearchTool(server *mcp.Server, useCase inbound.ScnmSearchUseCase) {
	handler := scnmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kini_scnm_search",
		Description: "산림청 국립수목원 국가표준곤충목록의 곤충 학명 목록을 조회합니다.",
	}, handler.handle)
}

func (h scnmSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input scnmSearchInput) (*mcp.CallToolResult, scnmSearchOutput, error) {
	result, err := h.useCase.ScnmSearch(ctx, application.ScnmSearchQuery{
		PageNo:    input.PageNo,
		NumOfRows: input.NumOfRows,
		ReqGnrlNm: input.ReqGnrlNm,
		ReqScnm:   input.ReqScnm,
		DateFrom:  input.DateFrom,
		DateTo:    input.DateTo,
	})
	if err != nil {
		return nil, scnmSearchOutput{}, err
	}

	items := make([]scnmSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = scnmSearchItem{
			SuperFalmNm:       item.SuperFalmNm,
			ClassKorNm:        item.ClassKorNm,
			ClassNm:           item.ClassNm,
			FalmKorNm:         item.FalmKorNm,
			FalmNm:            item.FalmNm,
			GenusKorNm:        item.GenusKorNm,
			GenusNm:           item.GenusNm,
			InsctGnrlNm:       item.InsctGnrlNm,
			InsctScnmID:       item.InsctScnmID,
			InsctSpecsScnm:    item.InsctSpecsScnm,
			LastUpdtDtm:       item.LastUpdtDtm,
			OrdKorNm:          item.OrdKorNm,
			OrdNm:             item.OrdNm,
			StpltScnmRltnCdNm: item.StpltScnmRltnCdNm,
			SubFalmKorNm:      item.SubFalmKorNm,
			SubFalmNm:         item.SubFalmNm,
			SuperFalmKorNm:    item.SuperFalmKorNm,
		}
	}

	return nil, scnmSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
