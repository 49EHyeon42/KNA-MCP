package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
)

type entogSpcmSearchInput struct {
	St string `json:"st" jsonschema:"검색어구분 (1: 국명 부분 검색, 2: 학명 부분 검색, 3: 국명 일치 검색, 4: 학명 일치 검색)"`
	Sw string `json:"sw" jsonschema:"검색대상어"`
	// dateGbn, dateFrom, and dateTo are not exposed because the upstream API ignores them.
	NumOfRows int `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	PageNo    int `json:"pageNo" jsonschema:"페이지 번호 (1 이상)"`
}

type entogSpcmSearchOutput struct {
	Items      []entogSpcmSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 결과 수"`
}

type entogSpcmSearchItem struct {
	Btnc             string `json:"btnc" jsonschema:"학명"`
	ClctDyDesc       string `json:"clctDyDesc" jsonschema:"채집일"`
	CprtCtnt         string `json:"cprtCtnt" jsonschema:"저작권"`
	DetailYn         string `json:"detailYn" jsonschema:"상세정보유무"`
	EntogOfnmKrlngNm string `json:"entogOfnmKrlngNm" jsonschema:"국명"`
	EntogSmplNo      string `json:"entogSmplNo" jsonschema:"표본번호"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	FrstRgstnDtm     string `json:"frstRgstnDtm" jsonschema:"최초등록일"`
	GenusKorNm       string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm          string `json:"genusNm" jsonschema:"속명"`
	ImgURL           string `json:"imgUrl" jsonschema:"이미지URL"`
	LastUpdtDtm      string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	OrdKorNm         string `json:"ordKorNm" jsonschema:"목국명"`
	OrdNm            string `json:"ordNm" jsonschema:"목명"`
}

type entogSpcmSearchHandler struct {
	useCase inbound.EntogSpcmSearchUseCase
}

func addEntogSpcmSearchTool(server *mcp.Server, useCase inbound.EntogSpcmSearchUseCase) {
	handler := entogSpcmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "entog_service_entog_spcm_search",
		Description: "산림청 국립수목원 내구강표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h entogSpcmSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input entogSpcmSearchInput) (*mcp.CallToolResult, entogSpcmSearchOutput, error) {
	result, err := h.useCase.EntogSpcmSearch(ctx, application.EntogSpcmSearchQuery{
		St:        input.St,
		Sw:        input.Sw,
		NumOfRows: input.NumOfRows,
		PageNo:    input.PageNo,
	})
	if err != nil {
		return nil, entogSpcmSearchOutput{}, err
	}

	items := make([]entogSpcmSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = entogSpcmSearchItem{
			Btnc:             item.Btnc,
			ClctDyDesc:       item.ClctDyDesc,
			CprtCtnt:         item.CprtCtnt,
			DetailYn:         item.DetailYn,
			EntogOfnmKrlngNm: item.EntogOfnmKrlngNm,
			EntogSmplNo:      item.EntogSmplNo,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			FrstRgstnDtm:     item.FrstRgstnDtm,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			ImgURL:           item.ImgURL,
			LastUpdtDtm:      item.LastUpdtDtm,
			OrdKorNm:         item.OrdKorNm,
			OrdNm:            item.OrdNm,
		}
	}

	return nil, entogSpcmSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
