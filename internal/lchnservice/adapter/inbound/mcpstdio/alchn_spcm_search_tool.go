package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
)

type alchnSpcmSearchInput struct {
	St        string `json:"st" jsonschema:"검색어구분 (1: 국명 부분 검색, 2: 학명 부분 검색, 3: 국명 일치 검색, 4: 학명 일치 검색)"`
	Sw        string `json:"sw" jsonschema:"검색대상어"`
	DateGbn   string `json:"dateGbn,omitempty" jsonschema:"날짜검색 구분 (1: 등록일, 2: 수정일)"`
	DateFrom  string `json:"dateFrom,omitempty" jsonschema:"검색 시작일 (dateGbn 입력 시 필수, yyyyMMdd)"`
	DateTo    string `json:"dateTo,omitempty" jsonschema:"검색 종료일 (dateGbn 입력 시 필수, yyyyMMdd)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	PageNo    int    `json:"pageNo" jsonschema:"페이지 번호 (1 이상)"`
}

type alchnSpcmSearchOutput struct {
	Items      []alchnSpcmSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"페이지당레코드수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체카운트"`
}

type alchnSpcmSearchItem struct {
	Btnc         string `json:"btnc" jsonschema:"학명"`
	CltrNm       string `json:"cltrNm" jsonschema:"채집자명"`
	CprtCtnt     string `json:"cprtCtnt" jsonschema:"저작권"`
	DetailYn     string `json:"detailYn" jsonschema:"상세유무"`
	EngNm        string `json:"engNm" jsonschema:"영문명"`
	FamilyKorNm  string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm     string `json:"familyNm" jsonschema:"과명"`
	FrstRgstnDtm string `json:"frstRgstnDtm" jsonschema:"최초등록일시"`
	GenusKorNm   string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm      string `json:"genusNm" jsonschema:"속명"`
	ImgURL       string `json:"imgUrl" jsonschema:"이미지URL"`
	JapNm        string `json:"japNm" jsonschema:"일어명"`
	LastUpdtDtm  string `json:"lastUpdtDtm" jsonschema:"최종수정일시"`
	LchnGnrlNm   string `json:"lchnGnrlNm" jsonschema:"국명"`
	LchnScnmID   string `json:"lchnScnmId" jsonschema:"학명ID"`
	LchnSmplNo   string `json:"lchnSmplNo" jsonschema:"표본번호"`
	PrkNm        string `json:"prkNm" jsonschema:"북한명"`
}

type alchnSpcmSearchHandler struct {
	useCase inbound.AlchnSpcmSearchUseCase
}

func addAlchnSpcmSearchTool(server *mcp.Server, useCase inbound.AlchnSpcmSearchUseCase) {
	handler := alchnSpcmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "lchn_service_alchn_spcm_search",
		Description: "산림청 국립수목원 지의류표본 목록을 검색합니다.",
	}, handler.handle)
}

func (h alchnSpcmSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input alchnSpcmSearchInput) (*mcp.CallToolResult, alchnSpcmSearchOutput, error) {
	result, err := h.useCase.AlchnSpcmSearch(ctx, application.AlchnSpcmSearchQuery{
		St:        input.St,
		Sw:        input.Sw,
		DateGbn:   input.DateGbn,
		DateFrom:  input.DateFrom,
		DateTo:    input.DateTo,
		NumOfRows: input.NumOfRows,
		PageNo:    input.PageNo,
	})
	if err != nil {
		return nil, alchnSpcmSearchOutput{}, err
	}

	items := make([]alchnSpcmSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = alchnSpcmSearchItem{
			Btnc:         item.Btnc,
			CltrNm:       item.CltrNm,
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
			LchnScnmID:   item.LchnScnmID,
			LchnSmplNo:   item.LchnSmplNo,
			PrkNm:        item.PrkNm,
		}
	}

	return nil, alchnSpcmSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
