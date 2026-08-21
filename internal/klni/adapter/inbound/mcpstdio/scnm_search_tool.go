package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/inbound"
)

type scnmSearchInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqGnrlNm string `json:"reqGnrlNm,omitempty" jsonschema:"검색하려는 지의류 국명 (부분 문자열 검색)"`
	ReqScnm   string `json:"reqScnm,omitempty" jsonschema:"검색하려는 지의류 학명 (대소문자를 구분하지 않는 부분 문자열 검색)"`
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
	StpltScnmRltnCdNm string `json:"stpltScnmRltnCdNm" jsonschema:"지의류 학명의 정명/이명 구분"`
	ClassKorNm        string `json:"classKorNm" jsonschema:"지의류 학명 분류군의 강국명"`
	ClassNm           string `json:"classNm" jsonschema:"지의류 학명 분류군의 강명(Class Name)"`
	FalmNm            string `json:"falmNm" jsonschema:"지의류 학명 분류군의 과명(Family Name)"`
	FalnKorNm         string `json:"falnKorNm" jsonschema:"지의류 학명 분류군의 과국명"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"지의류 학명 분류군의 속국명"`
	GenusNm           string `json:"genusNm" jsonschema:"지의류 학명 분류군의 속명(Genus Name)"`
	LastUpdtDtm       string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	LchnGnrlNm        string `json:"lchnGnrlNm" jsonschema:"지의류 국명(지의류명)"`
	LchnScnm          string `json:"lchnScnm" jsonschema:"지의류 학명"`
	LchnScnmID        string `json:"lchnScnmId" jsonschema:"지의류 학명ID"`
	OrdKorNm          string `json:"ordKorNm" jsonschema:"지의류 학명 분류군의 목국명"`
	OrdNm             string `json:"ordNm" jsonschema:"지의류 학명 분류군의 목명(Order Name)"`
	PhylumKorNm       string `json:"phylumKorNm" jsonschema:"지의류 학명 분류군의 문국명"`
	PhylumNm          string `json:"phylumNm" jsonschema:"지의류 학명 분류군의 문명(Phylum Name)"`
}

type scnmSearchHandler struct {
	useCase inbound.ScnmSearchUseCase
}

func addScnmSearchTool(server *mcp.Server, useCase inbound.ScnmSearchUseCase) {
	handler := scnmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "klni_scnm_search",
		Description: "산림청 국립수목원 국가표준지의류목록의 지의류 학명 목록을 검색합니다.",
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
			GenusKorNm:        item.GenusKorNm,
			GenusNm:           item.GenusNm,
			LastUpdtDtm:       item.LastUpdtDtm,
			LchnGnrlNm:        item.LchnGnrlNm,
			LchnScnm:          item.LchnScnm,
			LchnScnmID:        item.LchnScnmID,
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
