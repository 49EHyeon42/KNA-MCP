package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantSeedGrmntListToolName = "plant_resource_plant_seed_grmnt_list"

type plantSeedGrmntListInput struct {
	PageNo         int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows      int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSeedSpecsID string `json:"reqSeedSpecsId" jsonschema:"식물종자 기본정보 목록 검색 결과의 종자종 ID"`
}

type plantSeedGrmntListOutput struct {
	Items      []plantSeedGrmntListItem `json:"items"`
	NumOfRows  int                      `json:"numOfRows"`
	PageNo     int                      `json:"pageNo"`
	TotalCount int                      `json:"totalCount"`
}

type plantSeedGrmntListItem struct {
	AvrgGrmntDcnt     string `json:"avrgGrmntDcnt"`
	GrmntBfrPrcesCont string `json:"grmntBfrPrcesCont"`
	GrmntClmdmCont    string `json:"grmntClmdmCont"`
	GrmntDscrt        string `json:"grmntDscrt"`
	GrmntExprmNo      string `json:"grmntExprmNo"`
	GrmntExprmSeq     string `json:"grmntExprmSeq"`
	GrmntLightCndtn   string `json:"grmntLightCndtn"`
	GrmntRt           string `json:"grmntRt"`
	GrmntTmpCndtn     string `json:"grmntTmpCndtn"`
	PlantGnrlNm       string `json:"plantGnrlNm"`
	SeedNo            string `json:"seedNo"`
	SeedSpecsID       string `json:"seedSpecsId"`
}

type plantSeedGrmntListHandler struct {
	useCase inbound.PlantSeedGrmntListUseCase
}

func addPlantSeedGrmntListTool(server *mcp.Server, useCase inbound.PlantSeedGrmntListUseCase) {
	handler := plantSeedGrmntListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSeedGrmntListToolName,
		Description: "산림청 국립수목원 종자 발아율정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantSeedGrmntListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSeedGrmntListInput) (*mcp.CallToolResult, plantSeedGrmntListOutput, error) {
	result, err := h.useCase.PlantSeedGrmntList(ctx, application.PlantSeedGrmntListQuery{
		PageNo:         input.PageNo,
		NumOfRows:      input.NumOfRows,
		ReqSeedSpecsID: input.ReqSeedSpecsID,
	})
	if err != nil {
		return nil, plantSeedGrmntListOutput{}, err
	}

	items := make([]plantSeedGrmntListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSeedGrmntListItem{
			AvrgGrmntDcnt:     item.AvrgGrmntDcnt,
			GrmntBfrPrcesCont: item.GrmntBfrPrcesCont,
			GrmntClmdmCont:    item.GrmntClmdmCont,
			GrmntDscrt:        item.GrmntDscrt,
			GrmntExprmNo:      item.GrmntExprmNo,
			GrmntExprmSeq:     item.GrmntExprmSeq,
			GrmntLightCndtn:   item.GrmntLightCndtn,
			GrmntRt:           item.GrmntRt,
			GrmntTmpCndtn:     item.GrmntTmpCndtn,
			PlantGnrlNm:       item.PlantGnrlNm,
			SeedNo:            item.SeedNo,
			SeedSpecsID:       item.SeedSpecsID,
		}
	}

	return nil, plantSeedGrmntListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
