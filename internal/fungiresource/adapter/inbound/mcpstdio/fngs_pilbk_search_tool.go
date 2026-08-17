package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
)

type fngsPilbkSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 버섯의 국명 또는 학명"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type fngsPilbkSearchOutput struct {
	Items      []fngsPilbkSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 결과 수"`
}

type fngsPilbkSearchItem struct {
	FamilyKorNm string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm    string `json:"familyNm" jsonschema:"과명"`
	FngsGnrlNm  string `json:"fngsGnrlNm" jsonschema:"국명(버섯명)"`
	FngsPilbkNo string `json:"fngsPilbkNo" jsonschema:"버섯도감번호"`
	FngsScnm    string `json:"fngsScnm" jsonschema:"학명"`
	GenusKorNm  string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm     string `json:"genusNm" jsonschema:"속명"`
	LastUpdtDtm string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
}

type fngsPilbkSearchHandler struct {
	useCase inbound.FngsPilbkSearchUseCase
}

func addFngsPilbkSearchTool(server *mcp.Server, useCase inbound.FngsPilbkSearchUseCase) {
	handler := fngsPilbkSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fungi_resource_fngs_pilbk_search",
		Description: "산림청 국립수목원 버섯도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h fngsPilbkSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input fngsPilbkSearchInput) (*mcp.CallToolResult, fngsPilbkSearchOutput, error) {
	result, err := h.useCase.FngsPilbkSearch(ctx, application.FngsPilbkSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, fngsPilbkSearchOutput{}, err
	}

	items := make([]fngsPilbkSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = fngsPilbkSearchItem{
			FamilyKorNm: item.FamilyKorNm,
			FamilyNm:    item.FamilyNm,
			FngsGnrlNm:  item.FngsGnrlNm,
			FngsPilbkNo: item.FngsPilbkNo,
			FngsScnm:    item.FngsScnm,
			GenusKorNm:  item.GenusKorNm,
			GenusNm:     item.GenusNm,
			LastUpdtDtm: item.LastUpdtDtm,
		}
	}

	return nil, fngsPilbkSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
