package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantSmplUnitListToolName = "plant_resource_plant_smpl_unit_list"

type plantSmplUnitListInput struct {
	PageNo          int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows       int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqPlantSpecsID string `json:"reqPlantSpecsId" jsonschema:"식물표본 목록 검색 결과의 식물종 ID"`
}

type plantSmplUnitListOutput struct {
	Items      []plantSmplUnitListItem `json:"items"`
	NumOfRows  int                     `json:"numOfRows"`
	PageNo     int                     `json:"pageNo"`
	TotalCount int                     `json:"totalCount"`
}

type plantSmplUnitListItem struct {
	AgpFamilyKorNm     string `json:"agpFamilyKorNm"`
	AgpFamilyNm        string `json:"agpFamilyNm"`
	BspcsInsttNm       string `json:"bspcsInsttNm"`
	ClarHaslvVal       string `json:"clarHaslvVal"`
	ClarNm             string `json:"clarNm"`
	CllcrNm            string `json:"cllcrNm"`
	FamilyKorNm        string `json:"familyKorNm"`
	FamilyNm           string `json:"familyNm"`
	HbttChrcrCont      string `json:"hbttChrcrCont"`
	HbttTpcdNm         string `json:"hbttTpcdNm"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm"`
	PlantGnrlNm        string `json:"plantGnrlNm"`
	PlantPilbkNo       string `json:"plantPilbkNo"`
	PlantSmplNo        string `json:"plantSmplNo"`
	PlantSpecsID       string `json:"plantSpecsId"`
	PlantSpecsScnm     string `json:"plantSpecsScnm"`
	SmplCllcnDt        string `json:"smplCllcnDt"`
	SmplClnyNm         string `json:"smplClnyNm"`
	SmplKindCdNm       string `json:"smplKindCdNm"`
	SmplWrdt           string `json:"smplWrdt"`
	VgttnTpeCdNm       string `json:"vgttnTpeCdNm"`
}

type plantSmplUnitListHandler struct {
	useCase inbound.PlantSmplUnitListUseCase
}

func addPlantSmplUnitListTool(server *mcp.Server, useCase inbound.PlantSmplUnitListUseCase) {
	handler := plantSmplUnitListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSmplUnitListToolName,
		Description: "산림청 국립수목원 식물표본 상세정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantSmplUnitListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSmplUnitListInput) (*mcp.CallToolResult, plantSmplUnitListOutput, error) {
	result, err := h.useCase.PlantSmplUnitList(ctx, application.PlantSmplUnitListQuery{
		PageNo:          input.PageNo,
		NumOfRows:       input.NumOfRows,
		ReqPlantSpecsID: input.ReqPlantSpecsID,
	})
	if err != nil {
		return nil, plantSmplUnitListOutput{}, err
	}

	items := make([]plantSmplUnitListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSmplUnitListItem{
			AgpFamilyKorNm:     item.AgpFamilyKorNm,
			AgpFamilyNm:        item.AgpFamilyNm,
			BspcsInsttNm:       item.BspcsInsttNm,
			ClarHaslvVal:       item.ClarHaslvVal,
			ClarNm:             item.ClarNm,
			CllcrNm:            item.CllcrNm,
			FamilyKorNm:        item.FamilyKorNm,
			FamilyNm:           item.FamilyNm,
			HbttChrcrCont:      item.HbttChrcrCont,
			HbttTpcdNm:         item.HbttTpcdNm,
			PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:        item.PlantGnrlNm,
			PlantPilbkNo:       item.PlantPilbkNo,
			PlantSmplNo:        item.PlantSmplNo,
			PlantSpecsID:       item.PlantSpecsID,
			PlantSpecsScnm:     item.PlantSpecsScnm,
			SmplCllcnDt:        item.SmplCllcnDt,
			SmplClnyNm:         item.SmplClnyNm,
			SmplKindCdNm:       item.SmplKindCdNm,
			SmplWrdt:           item.SmplWrdt,
			VgttnTpeCdNm:       item.VgttnTpeCdNm,
		}
	}

	return nil, plantSmplUnitListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
