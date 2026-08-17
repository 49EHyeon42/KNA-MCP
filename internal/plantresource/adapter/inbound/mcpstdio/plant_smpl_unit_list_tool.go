package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantSmplUnitListInput struct {
	PageNo          int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows       int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqPlantSpecsID string `json:"reqPlantSpecsId" jsonschema:"검색할 식물표본의 식물종ID (plantSmplSearch 결과의 plantSpecsId)"`
}

type plantSmplUnitListOutput struct {
	Items      []plantSmplUnitListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                     `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                     `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                     `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantSmplUnitListItem struct {
	AgpFamilyKorNm     string `json:"agpFamilyKorNm" jsonschema:"APG과국명"`
	AgpFamilyNm        string `json:"agpFamilyNm" jsonschema:"APG과명"`
	BspcsInsttNm       string `json:"bspcsInsttNm" jsonschema:"표본소장기관"`
	ClarHaslvVal       string `json:"clarHaslvVal" jsonschema:"채집지해발고도"`
	ClarNm             string `json:"clarNm" jsonschema:"채집지"`
	CllcrNm            string `json:"cllcrNm" jsonschema:"채집자"`
	FamilyKorNm        string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm           string `json:"familyNm" jsonschema:"과명"`
	HbttChrcrCont      string `json:"hbttChrcrCont" jsonschema:"서식지특성"`
	HbttTpcdNm         string `json:"hbttTpcdNm" jsonschema:"서식지구분"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm" jsonschema:"식물번식형태"`
	PlantGnrlNm        string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantPilbkNo       string `json:"plantPilbkNo" jsonschema:"식물도감번호"`
	PlantSmplNo        string `json:"plantSmplNo" jsonschema:"식물표본번호"`
	PlantSpecsID       string `json:"plantSpecsId" jsonschema:"식물종ID"`
	PlantSpecsScnm     string `json:"plantSpecsScnm" jsonschema:"학명"`
	SmplCllcnDt        string `json:"smplCllcnDt" jsonschema:"채집일"`
	SmplClnyNm         string `json:"smplClnyNm" jsonschema:"표본군락명"`
	SmplKindCdNm       string `json:"smplKindCdNm" jsonschema:"표본종류"`
	SmplWrdt           string `json:"smplWrdt" jsonschema:"표본작성일"`
	VgttnTpeCdNm       string `json:"vgttnTpeCdNm" jsonschema:"식생유형"`
}

type plantSmplUnitListHandler struct {
	useCase inbound.PlantSmplUnitListUseCase
}

func addPlantSmplUnitListTool(server *mcp.Server, useCase inbound.PlantSmplUnitListUseCase) {
	handler := plantSmplUnitListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_smpl_unit_list",
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
