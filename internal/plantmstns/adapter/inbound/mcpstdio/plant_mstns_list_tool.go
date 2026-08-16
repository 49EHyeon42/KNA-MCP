package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application/port/inbound"
)

const plantMstnsPlantMstnsListToolName = "plant_mstns_plant_mstns_list"

type plantMstnsListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"식물세밀화의 식물 국명 또는 학명 검색어"`
	ReqMnfctYr   string `json:"reqMnfctYr,omitempty" jsonschema:"식물세밀화 제작년도"`
}

type plantMstnsListOutput struct {
	Items      []plantMstnsListItem `json:"items"`
	NumOfRows  int                  `json:"numOfRows"`
	PageNo     int                  `json:"pageNo"`
	TotalCount int                  `json:"totalCount"`
}

type plantMstnsListItem struct {
	DistrAraDscrt         string `json:"distrAraDscrt"`
	MinitrTpcdNm          string `json:"minitrTpcdNm"`
	PlantBrdgFomTpcdNm    string `json:"plantBrdgFomTpcdNm"`
	PlantGnrlNm           string `json:"plantGnrlNm"`
	PlantMinitrAthrNm     string `json:"plantMinitrAthrNm"`
	PlantMinitrMnfctMonth string `json:"plantMinitrMnfctMonth"`
	PlantMinitrMnfctYr    string `json:"plantMinitrMnfctYr"`
	PlantMinitrPsinsNm    string `json:"plantMinitrPsinsNm"`
	PlantSpecsScnm        string `json:"plantSpecsScnm"`
	RrnssPlantYn          string `json:"rrnssPlantYn"`
	SpcltPlantYn          string `json:"spcltPlantYn"`
}

type plantMstnsListHandler struct {
	useCase inbound.PlantMstnsListUseCase
}

func addPlantMstnsListTool(server *mcp.Server, useCase inbound.PlantMstnsListUseCase) {
	handler := plantMstnsListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantMstnsPlantMstnsListToolName,
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
