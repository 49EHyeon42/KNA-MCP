package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
)

type fngsSmplSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 버섯표본의 학명 또는 국명 (대소문자를 구분하지 않는 부분 문자열 검색)"`
}

type fngsSmplSearchOutput struct {
	Items      []fngsSmplSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                  `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                  `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                  `json:"totalCount" jsonschema:"전체 결과 수"`
}

type fngsSmplSearchItem struct {
	Cnt         string `json:"cnt" jsonschema:"표본 수"`
	FamilyKorNm string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm    string `json:"familyNm" jsonschema:"과명"`
	FngsGnrlNm  string `json:"fngsGnrlNm" jsonschema:"국명(버섯명)"`
	FngsID      string `json:"fngsId" jsonschema:"버섯 종ID"`
	FngsScnm    string `json:"fngsScnm" jsonschema:"학명"`
	GenusKorNm  string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm     string `json:"genusNm" jsonschema:"속명"`
}

type fngsSmplSearchHandler struct {
	useCase inbound.FngsSmplSearchUseCase
}

func addFngsSmplSearchTool(server *mcp.Server, useCase inbound.FngsSmplSearchUseCase) {
	handler := fngsSmplSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fungi_resource_fngs_smpl_search",
		Description: "산림청 국립수목원 버섯표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h fngsSmplSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input fngsSmplSearchInput) (*mcp.CallToolResult, fngsSmplSearchOutput, error) {
	result, err := h.useCase.FngsSmplSearch(ctx, application.FngsSmplSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, fngsSmplSearchOutput{}, err
	}

	items := make([]fngsSmplSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = fngsSmplSearchItem{
			Cnt:         item.Cnt,
			FamilyKorNm: item.FamilyKorNm,
			FamilyNm:    item.FamilyNm,
			FngsGnrlNm:  item.FngsGnrlNm,
			FngsID:      item.FngsID,
			FngsScnm:    item.FngsScnm,
			GenusKorNm:  item.GenusKorNm,
			GenusNm:     item.GenusNm,
		}
	}

	return nil, fngsSmplSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
