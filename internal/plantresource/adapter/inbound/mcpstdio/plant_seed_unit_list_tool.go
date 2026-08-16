package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantSeedUnitListInput struct {
	PageNo         int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows      int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSeedSpecsID string `json:"reqSeedSpecsId" jsonschema:"식물종자 기본정보 목록 검색 결과의 종자종 ID"`
}

type plantSeedUnitListOutput struct {
	Items      []plantSeedUnitListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                     `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                     `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                     `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantSeedUnitListItem struct {
	CllcnDate        string `json:"cllcnDate" jsonschema:"종자수집일"`
	PlantGnrlNm      string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	QualtFllnsRt     string `json:"qualtFllnsRt" jsonschema:"품질충실율"`
	SdwghWeght       string `json:"sdwghWeght" jsonschema:"천립중무게"`
	SeedAdmcn        string `json:"seedAdmcn" jsonschema:"종자기건함수율"`
	SeedCllctPlace   string `json:"seedCllctPlace" jsonschema:"종자수집장소"`
	SeedHoldGrainCnt string `json:"seedHoldGrainCnt" jsonschema:"종자보유립수"`
	SeedHoldQntt     string `json:"seedHoldQntt" jsonschema:"종자보유량"`
	SeedNo           string `json:"seedNo" jsonschema:"종자번호"`
	SeedSpecsID      string `json:"seedSpecsId" jsonschema:"종자종ID"`
	StoreChrcrTpcdNm string `json:"storeChrcrTpcdNm" jsonschema:"저장특성"`
	Vtlfct           string `json:"vtlfct" jsonschema:"활력률"`
	VtlfctTestYr     string `json:"vtlfctTestYr" jsonschema:"활력률테스트년도"`
}

type plantSeedUnitListHandler struct {
	useCase inbound.PlantSeedUnitListUseCase
}

func addPlantSeedUnitListTool(server *mcp.Server, useCase inbound.PlantSeedUnitListUseCase) {
	handler := plantSeedUnitListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_seed_unit_list",
		Description: "산림청 국립수목원 종자 점정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantSeedUnitListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSeedUnitListInput) (*mcp.CallToolResult, plantSeedUnitListOutput, error) {
	result, err := h.useCase.PlantSeedUnitList(ctx, application.PlantSeedUnitListQuery{
		PageNo:         input.PageNo,
		NumOfRows:      input.NumOfRows,
		ReqSeedSpecsID: input.ReqSeedSpecsID,
	})
	if err != nil {
		return nil, plantSeedUnitListOutput{}, err
	}

	items := make([]plantSeedUnitListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSeedUnitListItem{
			CllcnDate:        item.CllcnDate,
			PlantGnrlNm:      item.PlantGnrlNm,
			QualtFllnsRt:     item.QualtFllnsRt,
			SdwghWeght:       item.SdwghWeght,
			SeedAdmcn:        item.SeedAdmcn,
			SeedCllctPlace:   item.SeedCllctPlace,
			SeedHoldGrainCnt: item.SeedHoldGrainCnt,
			SeedHoldQntt:     item.SeedHoldQntt,
			SeedNo:           item.SeedNo,
			SeedSpecsID:      item.SeedSpecsID,
			StoreChrcrTpcdNm: item.StoreChrcrTpcdNm,
			Vtlfct:           item.Vtlfct,
			VtlfctTestYr:     item.VtlfctTestYr,
		}
	}

	return nil, plantSeedUnitListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
