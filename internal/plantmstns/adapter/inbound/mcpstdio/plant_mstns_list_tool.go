package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/port/inbound"
)

type plantMstnsListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"세밀화의 식물 국명 또는 학명"`
	ReqMnfctYr   string `json:"reqMnfctYr,omitempty" jsonschema:"세밀화 제작년도"`
}

type plantMstnsListOutput struct {
	Items      []plantMstnsListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                  `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                  `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                  `json:"totalCount" jsonschema:"전체 결과 수"`
}

type plantMstnsListItem struct {
	DistrAraDscrt         string `json:"distrAraDscrt" jsonschema:"분포정보"`
	MinitrTpcdNm          string `json:"minitrTpcdNm" jsonschema:"세밀화구분"`
	PlantBrdgFomTpcdNm    string `json:"plantBrdgFomTpcdNm" jsonschema:"식물번식형태"`
	PlantGnrlNm           string `json:"plantGnrlNm" jsonschema:"국명"`
	PlantMinitrAthrNm     string `json:"plantMinitrAthrNm" jsonschema:"작가명"`
	PlantMinitrMnfctMonth string `json:"plantMinitrMnfctMonth" jsonschema:"제작월"`
	PlantMinitrMnfctYr    string `json:"plantMinitrMnfctYr" jsonschema:"제작년도"`
	PlantMinitrPsinsNm    string `json:"plantMinitrPsinsNm" jsonschema:"보유기관"`
	PlantSpecsScnm        string `json:"plantSpecsScnm" jsonschema:"학명"`
	RrnssPlantYn          string `json:"rrnssPlantYn" jsonschema:"희귀식물여부"`
	SpcltPlantYn          string `json:"spcltPlantYn" jsonschema:"특산식물여부"`
}

type plantMstnsListHandler struct {
	useCase inbound.PlantMstnsListUseCase
}

func addPlantMstnsListTool(server *mcp.Server, useCase inbound.PlantMstnsListUseCase) {
	handler := plantMstnsListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_mstns_plant_mstns_list",
		Description: "산림청 국립수목원 식물세밀화 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantMstnsListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantMstnsListInput) (*mcp.CallToolResult, plantMstnsListOutput, error) {
	result, err := h.useCase.PlantMstnsList(ctx, application.PlantMstnsListQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
		ReqMnfctYr:   input.ReqMnfctYr,
	})
	if err != nil {
		return nil, plantMstnsListOutput{}, err
	}

	items := make([]plantMstnsListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantMstnsListItem{
			DistrAraDscrt:         item.DistrAraDscrt,
			MinitrTpcdNm:          item.MinitrTpcdNm,
			PlantBrdgFomTpcdNm:    item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:           item.PlantGnrlNm,
			PlantMinitrAthrNm:     item.PlantMinitrAthrNm,
			PlantMinitrMnfctMonth: item.PlantMinitrMnfctMonth,
			PlantMinitrMnfctYr:    item.PlantMinitrMnfctYr,
			PlantMinitrPsinsNm:    item.PlantMinitrPsinsNm,
			PlantSpecsScnm:        item.PlantSpecsScnm,
			RrnssPlantYn:          item.RrnssPlantYn,
			SpcltPlantYn:          item.SpcltPlantYn,
		}
	}

	return nil, plantMstnsListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
