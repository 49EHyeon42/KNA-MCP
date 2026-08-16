package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantSeedUnitListToolName = "plant_resource_plant_seed_unit_list"

type plantSeedUnitListInput struct {
	PageNo         int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows      int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSeedSpecsID string `json:"reqSeedSpecsId" jsonschema:"식물종자 기본정보 목록 검색 결과의 종자종 ID"`
}

type plantSeedUnitListOutput struct {
	Items      []plantSeedUnitListItem `json:"items"`
	NumOfRows  int                     `json:"numOfRows"`
	PageNo     int                     `json:"pageNo"`
	TotalCount int                     `json:"totalCount"`
}

type plantSeedUnitListItem struct {
	CllcnDate        string `json:"cllcnDate"`
	PlantGnrlNm      string `json:"plantGnrlNm"`
	QualtFllnsRt     string `json:"qualtFllnsRt"`
	SdwghWeght       string `json:"sdwghWeght"`
	SeedAdmcn        string `json:"seedAdmcn"`
	SeedCllctPlace   string `json:"seedCllctPlace"`
	SeedHoldGrainCnt string `json:"seedHoldGrainCnt"`
	SeedHoldQntt     string `json:"seedHoldQntt"`
	SeedNo           string `json:"seedNo"`
	SeedSpecsID      string `json:"seedSpecsId"`
	StoreChrcrTpcdNm string `json:"storeChrcrTpcdNm"`
	Vtlfct           string `json:"vtlfct"`
	VtlfctTestYr     string `json:"vtlfctTestYr"`
}

type plantSeedUnitListHandler struct {
	useCase inbound.PlantSeedUnitListUseCase
}

func addPlantSeedUnitListTool(server *mcp.Server, useCase inbound.PlantSeedUnitListUseCase) {
	handler := plantSeedUnitListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSeedUnitListToolName,
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
