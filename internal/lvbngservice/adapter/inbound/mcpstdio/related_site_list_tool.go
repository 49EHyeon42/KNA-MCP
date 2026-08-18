package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application/port/inbound"
)

type relatedSiteListInput struct {
	PageNo    int `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
}

type relatedSiteListOutput struct {
	Items      []relatedSiteListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 결과 수"`
}

type relatedSiteListItem struct {
	LvbngTpcdNm string `json:"lvbngTpcdNm" jsonschema:"생물분류"`
	SiteCtgryNm string `json:"siteCtgryNm" jsonschema:"카테고리분류"`
	SiteNm      string `json:"siteNm" jsonschema:"사이트명"`
	SiteURL     string `json:"siteUrl" jsonschema:"사이트URL"`
}

type relatedSiteListHandler struct {
	useCase inbound.RelatedSiteListUseCase
}

func addRelatedSiteListTool(server *mcp.Server, useCase inbound.RelatedSiteListUseCase) {
	handler := relatedSiteListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "lvbng_service_related_site_list",
		Description: "산림청 국립수목원 생물관련사이트 목록을 조회합니다.",
	}, handler.handle)
}

func (h relatedSiteListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input relatedSiteListInput) (*mcp.CallToolResult, relatedSiteListOutput, error) {
	result, err := h.useCase.RelatedSiteList(ctx, application.RelatedSiteListQuery{
		PageNo:    input.PageNo,
		NumOfRows: input.NumOfRows,
	})
	if err != nil {
		return nil, relatedSiteListOutput{}, err
	}

	items := make([]relatedSiteListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = relatedSiteListItem{
			LvbngTpcdNm: item.LvbngTpcdNm,
			SiteCtgryNm: item.SiteCtgryNm,
			SiteNm:      item.SiteNm,
			SiteURL:     item.SiteURL,
		}
	}

	return nil, relatedSiteListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
