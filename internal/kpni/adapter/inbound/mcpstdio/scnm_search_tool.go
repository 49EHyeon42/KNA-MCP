package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/inbound"
)

type scnmSearchInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqGnrlNm string `json:"reqGnrlNm,omitempty" jsonschema:"검색하려는 식물 국명 (부분 문자열 검색)"`
	ReqScnm   string `json:"reqScnm,omitempty" jsonschema:"검색하려는 식물 학명 (대소문자를 구분하지 않는 부분 문자열 검색)"`
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
	ClassKorNm          string `json:"classKorNm" jsonschema:"식물 학명 분류군 강국명"`
	ClassNm             string `json:"classNm" jsonschema:"식물 학명 분류군 강명(Class Name)"`
	FalmKorNm           string `json:"falmKorNm" jsonschema:"식물 학명 분류군 과국명"`
	FalmNm              string `json:"falmNm" jsonschema:"식물 학명 분류군 과명(Family Name)"`
	GenusKorNm          string `json:"genusKorNm" jsonschema:"식물 학명 분류군 속국명"`
	GenusNm             string `json:"genusNm" jsonschema:"식물 학명 분류군 속명(Genus Name)"`
	LastUpdtDtm         string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	OrdKorNm            string `json:"ordKorNm" jsonschema:"식물 학명 분류군 목국명"`
	OrdNm               string `json:"ordNm" jsonschema:"식물 학명 분류군 목명(Order Name)"`
	PhylumKorNm         string `json:"phylumKorNm" jsonschema:"식물 학명 분류군 문국명"`
	PhylumNm            string `json:"phylumNm" jsonschema:"식물 학명 분류군 문명(Phylum Name)"`
	PlantGnrlNm         string `json:"plantGnrlNm" jsonschema:"식물 국명"`
	PlantScnmID         string `json:"plantScnmId" jsonschema:"식물 학명ID"`
	PlantSpecsClsscCdNm string `json:"plantSpecsClsscCdNm" jsonschema:"식물 분류명(자생, 재배, 외래)"`
	PlantSpecsScnm      string `json:"plantSpecsScnm" jsonschema:"식물 학명"`
	StpltScnmRltnCdNm   string `json:"stpltScnmRltnCdNm" jsonschema:"식물 학명 정명/이명 구분"`
	SubClassKorNm       string `json:"subClassKorNm" jsonschema:"식물 학명 분류군 아강국명"`
	SubClassNm          string `json:"subClassNm" jsonschema:"식물 학명 분류군 아강명(SubClass Name)"`
}

type scnmSearchHandler struct {
	useCase inbound.ScnmSearchUseCase
}

func addScnmSearchTool(server *mcp.Server, useCase inbound.ScnmSearchUseCase) {
	handler := scnmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kpni_scnm_search",
		Description: "산림청 국립수목원 국가표준식물목록의 식물 학명 목록을 조회합니다.",
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
			ClassKorNm:          item.ClassKorNm,
			ClassNm:             item.ClassNm,
			FalmKorNm:           item.FalmKorNm,
			FalmNm:              item.FalmNm,
			GenusKorNm:          item.GenusKorNm,
			GenusNm:             item.GenusNm,
			LastUpdtDtm:         item.LastUpdtDtm,
			OrdKorNm:            item.OrdKorNm,
			OrdNm:               item.OrdNm,
			PhylumKorNm:         item.PhylumKorNm,
			PhylumNm:            item.PhylumNm,
			PlantGnrlNm:         item.PlantGnrlNm,
			PlantScnmID:         item.PlantScnmID,
			PlantSpecsClsscCdNm: item.PlantSpecsClsscCdNm,
			PlantSpecsScnm:      item.PlantSpecsScnm,
			StpltScnmRltnCdNm:   item.StpltScnmRltnCdNm,
			SubClassKorNm:       item.SubClassKorNm,
			SubClassNm:          item.SubClassNm,
		}
	}

	return nil, scnmSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
