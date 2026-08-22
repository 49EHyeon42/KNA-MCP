package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application/port/inbound"
)

type oldSpcmSearchInput struct {
	St        string `json:"st" jsonschema:"검색어 구분 (1: 학명 부분 검색, 2: 학명 일치 검색)"`
	Sw        string `json:"sw" jsonschema:"검색대상어"`
	DateGbn   string `json:"dateGbn,omitempty" jsonschema:"날짜검색구분 (1: 등록일, 2: 수정일)"`
	DateFrom  string `json:"dateFrom,omitempty" jsonschema:"검색시작일 (dateGbn 입력 시 필수, yyyyMMdd)"`
	DateTo    string `json:"dateTo,omitempty" jsonschema:"검색종료일 (dateGbn 입력 시 필수, yyyyMMdd)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
}

type oldSpcmSearchOutput struct {
	Items      []oldSpcmSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                 `json:"numOfRows" jsonschema:"페이지당레코드수"`
	PageNo     int                 `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                 `json:"totalCount" jsonschema:"전체카운트"`
}

type oldSpcmSearchItem struct {
	CprtCtnt       string `json:"cprtCtnt" jsonschema:"저작권"`
	FamlKorNm      string `json:"famlKorNm" jsonschema:"과국명"`
	FamlNm         string `json:"famlNm" jsonschema:"과명"`
	FrstRgstnDtm   string `json:"frstRgstnDtm" jsonschema:"최초등록일"`
	ImgURL         string `json:"imgUrl" jsonschema:"이미지URL"`
	LastUpdtDtm    string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	PlantGnrlNm    string `json:"plantGnrlNm" jsonschema:"국명"`
	PlantOldSmplNo string `json:"plantOldSmplNo" jsonschema:"고표본번호"`
	PlantSpecsScnm string `json:"plantSpecsScnm" jsonschema:"학명"`
}

type oldSpcmSearchHandler struct {
	useCase inbound.OldSpcmSearchUseCase
}

func addOldSpcmSearchTool(server *mcp.Server, useCase inbound.OldSpcmSearchUseCase) {
	handler := oldSpcmSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "old_plant_service_old_spcm_search",
		Description: "산림청 국립수목원 한반도고표본 정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h oldSpcmSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input oldSpcmSearchInput) (*mcp.CallToolResult, oldSpcmSearchOutput, error) {
	result, err := h.useCase.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{
		St:        input.St,
		Sw:        input.Sw,
		DateGbn:   input.DateGbn,
		DateFrom:  input.DateFrom,
		DateTo:    input.DateTo,
		NumOfRows: input.NumOfRows,
		PageNo:    input.PageNo,
	})
	if err != nil {
		return nil, oldSpcmSearchOutput{}, err
	}

	items := make([]oldSpcmSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = oldSpcmSearchItem{
			CprtCtnt:       item.CprtCtnt,
			FamlKorNm:      item.FamlKorNm,
			FamlNm:         item.FamlNm,
			FrstRgstnDtm:   item.FrstRgstnDtm,
			ImgURL:         item.ImgURL,
			LastUpdtDtm:    item.LastUpdtDtm,
			PlantGnrlNm:    item.PlantGnrlNm,
			PlantOldSmplNo: item.PlantOldSmplNo,
			PlantSpecsScnm: item.PlantSpecsScnm,
		}
	}

	return nil, oldSpcmSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
