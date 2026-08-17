package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/inbound"
)

type gnrlNmLtrtrSearchInput struct {
	PageNo         int    `json:"pageNo" jsonschema:"페이지 번호 (1 이상)"`
	NumOfRows      int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqPlantGnrlNm string `json:"reqPlantGnrlNm,omitempty" jsonschema:"검색하려는 식물 국명(식물명) (부분 문자열 검색)"`
}

type gnrlNmLtrtrSearchOutput struct {
	Items      []gnrlNmLtrtrSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                     `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                     `json:"pageNo" jsonschema:"페이지 번호"`
	TotalCount int                     `json:"totalCount" jsonschema:"전체 결과 수"`
}

type gnrlNmLtrtrSearchItem struct {
	RcmmnTpcdNm      string `json:"rcmmnTpcdNm" jsonschema:"식물 국명 추천/비추천 구분"`
	LtrtrInfrmNm     string `json:"ltrtrInfrmNm" jsonschema:"식물 국명 출전 기재문"`
	LvbngFrlngTpcdNm string `json:"lvbngFrlngTpcdNm" jsonschema:"국명 언어 분류"`
	PlantGnrlNm      string `json:"plantGnrlNm" jsonschema:"식물 국명(식물명)"`
	PlantSpecsScnm   string `json:"plantSpecsScnm" jsonschema:"식물 학명"`
}

type gnrlNmLtrtrSearchHandler struct {
	useCase inbound.GnrlNmLtrtrSearchUseCase
}

func addGnrlNmLtrtrSearchTool(server *mcp.Server, useCase inbound.GnrlNmLtrtrSearchUseCase) {
	handler := gnrlNmLtrtrSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kpni_gnrl_nm_ltrtr_search",
		Description: "산림청 국립수목원 국가표준식물목록의 식물 국명 출전 정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h gnrlNmLtrtrSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input gnrlNmLtrtrSearchInput) (*mcp.CallToolResult, gnrlNmLtrtrSearchOutput, error) {
	result, err := h.useCase.GnrlNmLtrtrSearch(ctx, application.GnrlNmLtrtrSearchQuery{
		PageNo:         input.PageNo,
		NumOfRows:      input.NumOfRows,
		ReqPlantGnrlNm: input.ReqPlantGnrlNm,
	})
	if err != nil {
		return nil, gnrlNmLtrtrSearchOutput{}, err
	}

	items := make([]gnrlNmLtrtrSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = gnrlNmLtrtrSearchItem{
			RcmmnTpcdNm:      item.RcmmnTpcdNm,
			LtrtrInfrmNm:     item.LtrtrInfrmNm,
			LvbngFrlngTpcdNm: item.LvbngFrlngTpcdNm,
			PlantGnrlNm:      item.PlantGnrlNm,
			PlantSpecsScnm:   item.PlantSpecsScnm,
		}
	}

	return nil, gnrlNmLtrtrSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
