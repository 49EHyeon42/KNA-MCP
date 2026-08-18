package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
)

type alchnIlstrSearchInput struct {
	St        string `json:"st" jsonschema:"검색어구분 (1: 국명 부분 검색, 2: 학명 부분 검색, 3: 국명 일치 검색, 4: 학명 일치 검색)"`
	Sw        string `json:"sw" jsonschema:"검색대상어"`
	DateGbn   string `json:"dateGbn,omitempty" jsonschema:"날짜검색 구분 (1: 등록일, 2: 수정일)"`
	DateFrom  string `json:"dateFrom,omitempty" jsonschema:"검색 시작일 (dateGbn 입력 시 필수, yyyyMMdd)"`
	DateTo    string `json:"dateTo,omitempty" jsonschema:"검색 종료일 (dateGbn 입력 시 필수, yyyyMMdd)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	PageNo    int    `json:"pageNo" jsonschema:"페이지 번호 (1 이상)"`
}

type alchnIlstrSearchOutput struct {
	Items      []alchnIlstrSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                    `json:"numOfRows" jsonschema:"페이지당레코드수"`
	PageNo     int                    `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                    `json:"totalCount" jsonschema:"전체카운트"`
}

type alchnIlstrSearchItem struct {
	Btnc         string `json:"btnc" jsonschema:"학명"`
	CprtCtnt     string `json:"cprtCtnt" jsonschema:"저작권"`
	DetailYn     string `json:"detailYn" jsonschema:"상세정보유무"`
	EngNm        string `json:"engNm" jsonschema:"영문명"`
	FamilyKorNm  string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm     string `json:"familyNm" jsonschema:"과명"`
	FrstRgstnDtm string `json:"frstRgstnDtm" jsonschema:"최초등록일시"`
	GenusKorNm   string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm      string `json:"genusNm" jsonschema:"속명"`
	ImgURL       string `json:"imgUrl" jsonschema:"이미지URL"`
	JapNm        string `json:"japNm" jsonschema:"일본명"`
	LastUpdtDtm  string `json:"lastUpdtDtm" jsonschema:"최종수정일시"`
	LchnGnrlNm   string `json:"lchnGnrlNm" jsonschema:"국명"`
	LchnInfrpNm  string `json:"lchnInfrpNm" jsonschema:"종하명"`
	LchnPilbkNo  string `json:"lchnPilbkNo" jsonschema:"도감번호"`
	LchnScnmID   string `json:"lchnScnmId" jsonschema:"학명ID"`
	LchnTtnm     string `json:"lchnTtnm" jsonschema:"종소명"`
	PrkNm        string `json:"prkNm" jsonschema:"북한명"`
}

type alchnIlstrSearchHandler struct {
	useCase inbound.AlchnIlstrSearchUseCase
}

func addAlchnIlstrSearchTool(server *mcp.Server, useCase inbound.AlchnIlstrSearchUseCase) {
	handler := alchnIlstrSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "lchn_service_alchn_ilstr_search",
		Description: "산림청 국립수목원 지의류도감 목록을 검색합니다.",
	}, handler.handle)
}

func (h alchnIlstrSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input alchnIlstrSearchInput) (*mcp.CallToolResult, alchnIlstrSearchOutput, error) {
	result, err := h.useCase.AlchnIlstrSearch(ctx, application.AlchnIlstrSearchQuery{
		St:        input.St,
		Sw:        input.Sw,
		DateGbn:   input.DateGbn,
		DateFrom:  input.DateFrom,
		DateTo:    input.DateTo,
		NumOfRows: input.NumOfRows,
		PageNo:    input.PageNo,
	})
	if err != nil {
		return nil, alchnIlstrSearchOutput{}, err
	}

	items := make([]alchnIlstrSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = alchnIlstrSearchItem{
			Btnc:         item.Btnc,
			CprtCtnt:     item.CprtCtnt,
			DetailYn:     item.DetailYn,
			EngNm:        item.EngNm,
			FamilyKorNm:  item.FamilyKorNm,
			FamilyNm:     item.FamilyNm,
			FrstRgstnDtm: item.FrstRgstnDtm,
			GenusKorNm:   item.GenusKorNm,
			GenusNm:      item.GenusNm,
			ImgURL:       item.ImgURL,
			JapNm:        item.JapNm,
			LastUpdtDtm:  item.LastUpdtDtm,
			LchnGnrlNm:   item.LchnGnrlNm,
			LchnInfrpNm:  item.LchnInfrpNm,
			LchnPilbkNo:  item.LchnPilbkNo,
			LchnScnmID:   item.LchnScnmID,
			LchnTtnm:     item.LchnTtnm,
			PrkNm:        item.PrkNm,
		}
	}

	return nil, alchnIlstrSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
