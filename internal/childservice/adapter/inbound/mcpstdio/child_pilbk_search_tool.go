package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/inbound"
)

type childPilbkSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 생물의 국명 또는 학명"`
}

type childPilbkSearchOutput struct {
	Items      []childPilbkSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                    `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                    `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                    `json:"totalCount" jsonschema:"전체 결과 수"`
}

type childPilbkSearchItem struct {
	BiogyNm           string `json:"biogyNm" jsonschema:"생물학명"`
	ChildLvbngPilbkNo string `json:"childLvbngPilbkNo" jsonschema:"어린이생물도감번호"`
	FamilyKorNm       string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm          string `json:"familyNm" jsonschema:"과명"`
	LvbngTpcdNm       string `json:"lvbngTpcdNm" jsonschema:"생물분류"`
	LvngKrlngNm       string `json:"lvngKrlngNm" jsonschema:"생물국명"`
}

type childPilbkSearchHandler struct {
	useCase inbound.ChildPilbkSearchUseCase
}

func addChildPilbkSearchTool(server *mcp.Server, useCase inbound.ChildPilbkSearchUseCase) {
	handler := childPilbkSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "child_service_child_pilbk_search",
		Description: "산림청 국립수목원 어린이생물도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h childPilbkSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input childPilbkSearchInput) (*mcp.CallToolResult, childPilbkSearchOutput, error) {
	result, err := h.useCase.ChildPilbkSearch(ctx, application.ChildPilbkSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, childPilbkSearchOutput{}, err
	}

	items := make([]childPilbkSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = childPilbkSearchItem{
			BiogyNm:           item.BiogyNm,
			ChildLvbngPilbkNo: item.ChildLvbngPilbkNo,
			FamilyKorNm:       item.FamilyKorNm,
			FamilyNm:          item.FamilyNm,
			LvbngTpcdNm:       item.LvbngTpcdNm,
			LvngKrlngNm:       item.LvngKrlngNm,
		}
	}

	return nil, childPilbkSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
