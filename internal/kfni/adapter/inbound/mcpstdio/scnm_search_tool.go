package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application/port/inbound"
)

type scnmSearchInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqGnrlNm string `json:"reqGnrlNm,omitempty" jsonschema:"검색하려는 버섯 국명(버섯명) (부분 문자열 검색)"`
	ReqScnm   string `json:"reqScnm,omitempty" jsonschema:"검색하려는 버섯 학명 (대소문자를 구분하지 않는 부분 문자열 검색)"`
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
	StpltScnmRltnCdNm string `json:"stpltScnmRltnCdNm" jsonschema:"버섯 학명의 정/이명 구분"`
	ClassKorNm        string `json:"classKorNm" jsonschema:"버섯 분류군  강명의 한국명칭"`
	ClassNm           string `json:"classNm" jsonschema:"버섯 분류군의 강명(Class Name)"`
	FalmNm            string `json:"falmNm" jsonschema:"버섯 분류군의  과명(family Name)"`
	FalnKorNm         string `json:"falnKorNm" jsonschema:"버섯 분류군  과명의 한국명칭"`
	FngsGnrlNm        string `json:"fngsGnrlNm" jsonschema:"버섯 국명(버섯명)"`
	FngsScnm          string `json:"fngsScnm" jsonschema:"버섯 학명"`
	FngsScnmID        string `json:"fngsScnmId" jsonschema:"버섯 학명ID"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"버섯 분류군  속명의 한국명칭"`
	GenusNm           string `json:"genusNm" jsonschema:"버섯 분류군의  속명(Genus Name)"`
	LastUpdtDtm       string `json:"lastUpdtDtm" jsonschema:"최종 수정일"`
	OrdKorNm          string `json:"ordKorNm" jsonschema:"버섯 분류군  목명의 한국명칭"`
	OrdNm             string `json:"ordNm" jsonschema:"버섯 분류군의  목명(Order Name)"`
	PhylumKorNm       string `json:"phylumKorNm" jsonschema:"버섯 분류군  문명의 한국명칭"`
	PhylumNm          string `json:"phylumNm" jsonschema:"버섯 분류군의  문명(Phylum Name)"`
}

type scnmSearchHandler struct {
	useCase inbound.ScnmSearchUseCase
}

func addScnmSearchTool(server *mcp.Server, useCase inbound.ScnmSearchUseCase) {
	handler := scnmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kfni_scnm_search",
		Description: "산림청 국립수목원 국가표준버섯목록의 버섯 학명 목록을 검색합니다.",
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
			StpltScnmRltnCdNm: item.StpltScnmRltnCdNm,
			ClassKorNm:        item.ClassKorNm,
			ClassNm:           item.ClassNm,
			FalmNm:            item.FalmNm,
			FalnKorNm:         item.FalnKorNm,
			FngsGnrlNm:        item.FngsGnrlNm,
			FngsScnm:          item.FngsScnm,
			FngsScnmID:        item.FngsScnmID,
			GenusKorNm:        item.GenusKorNm,
			GenusNm:           item.GenusNm,
			LastUpdtDtm:       item.LastUpdtDtm,
			OrdKorNm:          item.OrdKorNm,
			OrdNm:             item.OrdNm,
			PhylumKorNm:       item.PhylumKorNm,
			PhylumNm:          item.PhylumNm,
		}
	}

	return nil, scnmSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
