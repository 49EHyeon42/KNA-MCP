package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
)

type entogIlstrSearchInput struct {
	St        string `json:"st" jsonschema:"검색어구분 (1: 국명 부분 검색, 2: 학명 부분 검색, 3: 국명 일치 검색, 4: 학명 일치 검색)"`
	Sw        string `json:"sw" jsonschema:"검색대상어"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	PageNo    int    `json:"pageNo" jsonschema:"페이지 번호 (1 이상)"`
}

type entogIlstrSearchOutput struct {
	Items      []entogIlstrSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                    `json:"numOfRows" jsonschema:"페이지당레코드수"`
	PageNo     int                    `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                    `json:"totalCount" jsonschema:"전체카운트"`
}

type entogIlstrSearchItem struct {
	Btnc             string `json:"btnc" jsonschema:"학명"`
	CprtCtnt         string `json:"cprtCtnt" jsonschema:"저작권"`
	DetailYn         string `json:"detailYn" jsonschema:"상세정보유무"`
	EntogOfnmKrlngNm string `json:"entogOfnmKrlngNm" jsonschema:"국명"`
	EntogOfnmScnmID  string `json:"entogOfnmScnmId" jsonschema:"학명ID"`
	EntogPilbkNo     string `json:"entogPilbkNo" jsonschema:"도감번호"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	GenusKorNm       string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm          string `json:"genusNm" jsonschema:"속명"`
	ImgURL           string `json:"imgUrl" jsonschema:"이미지URL"`
	OrdKorNm         string `json:"ordKorNm" jsonschema:"목국명"`
	OrdNm            string `json:"ordNm" jsonschema:"목명"`
}

type entogIlstrSearchHandler struct {
	useCase inbound.EntogIlstrSearchUseCase
}

func addEntogIlstrSearchTool(server *mcp.Server, useCase inbound.EntogIlstrSearchUseCase) {
	handler := entogIlstrSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "entog_service_entog_ilstr_search",
		Description: "산림청 국립수목원 내구강도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h entogIlstrSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input entogIlstrSearchInput) (*mcp.CallToolResult, entogIlstrSearchOutput, error) {
	result, err := h.useCase.EntogIlstrSearch(ctx, application.EntogIlstrSearchQuery{
		St:        input.St,
		Sw:        input.Sw,
		NumOfRows: input.NumOfRows,
		PageNo:    input.PageNo,
	})
	if err != nil {
		return nil, entogIlstrSearchOutput{}, err
	}

	items := make([]entogIlstrSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = entogIlstrSearchItem{
			Btnc:             item.Btnc,
			CprtCtnt:         item.CprtCtnt,
			DetailYn:         item.DetailYn,
			EntogOfnmKrlngNm: item.EntogOfnmKrlngNm,
			EntogOfnmScnmID:  item.EntogOfnmScnmID,
			EntogPilbkNo:     item.EntogPilbkNo,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			ImgURL:           item.ImgURL,
			OrdKorNm:         item.OrdKorNm,
			OrdNm:            item.OrdNm,
		}
	}

	return nil, entogIlstrSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
