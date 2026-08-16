package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantSeedGrmntListInput struct {
	PageNo         int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows      int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSeedSpecsID string `json:"reqSeedSpecsId" jsonschema:"식물종자 기본정보 목록 검색 결과의 종자종 ID"`
}

type plantSeedGrmntListOutput struct {
	Items      []plantSeedGrmntListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                      `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                      `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                      `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantSeedGrmntListItem struct {
	AvrgGrmntDcnt     string `json:"avrgGrmntDcnt" jsonschema:"평균 발아 일수"`
	GrmntBfrPrcesCont string `json:"grmntBfrPrcesCont" jsonschema:"발아 전처리 내용"`
	GrmntClmdmCont    string `json:"grmntClmdmCont" jsonschema:"발아배지내용"`
	GrmntDscrt        string `json:"grmntDscrt" jsonschema:"발아설명"`
	GrmntExprmNo      string `json:"grmntExprmNo" jsonschema:"실험번호"`
	GrmntExprmSeq     string `json:"grmntExprmSeq" jsonschema:"실험순번"`
	GrmntLightCndtn   string `json:"grmntLightCndtn" jsonschema:"광조건"`
	GrmntRt           string `json:"grmntRt" jsonschema:"발아율"`
	GrmntTmpCndtn     string `json:"grmntTmpCndtn" jsonschema:"온도조건"`
	PlantGnrlNm       string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	SeedNo            string `json:"seedNo" jsonschema:"종자번호"`
	SeedSpecsID       string `json:"seedSpecsId" jsonschema:"종자종ID"`
}

type plantSeedGrmntListHandler struct {
	useCase inbound.PlantSeedGrmntListUseCase
}

func addPlantSeedGrmntListTool(server *mcp.Server, useCase inbound.PlantSeedGrmntListUseCase) {
	handler := plantSeedGrmntListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_seed_grmnt_list",
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
