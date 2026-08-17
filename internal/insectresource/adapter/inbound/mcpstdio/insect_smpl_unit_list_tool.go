package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
)

type insectSmplUnitListInput struct {
	PageNo          int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows       int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqInsctSpecsID string `json:"reqInsctSpecsId" jsonschema:"검색할 곤충 종ID (insectSmplSearch 결과의 insctSpecsId)"`
}

type insectSmplUnitListOutput struct {
	Items      []insectSmplUnitListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                      `json:"numOfRows" jsonschema:"한페이지 결과 수"`
	PageNo     int                      `json:"pageNo" jsonschema:"페이지 번호"`
	TotalCount int                      `json:"totalCount" jsonschema:"전체 결과  건 수"`
}

type insectSmplUnitListItem struct {
	BspcsInsttNm       string `json:"bspcsInsttNm" jsonschema:"표본소장기관"`
	ClarHaslvVal       string `json:"clarHaslvVal" jsonschema:"채집지해발고도"`
	SmplCllcnDt        string `json:"smplCllcnDt" jsonschema:"표본채집일"`
	GynndTpcd          string `json:"gynndTpcd" jsonschema:"암수구분"`
	HbttTpcd           string `json:"hbttTpcd" jsonschema:"서식지구분"`
	InsctSmplNo        string `json:"insctSmplNo" jsonschema:"곤충표본번호"`
	InsctSpecsID       string `json:"insctSpecsId" jsonschema:"곤충종ID"`
	InsctSpecsScnm     string `json:"insctSpecsScnm" jsonschema:"학명"`
	LabelUsgCllcnNmplc string `json:"labelUsgCllcnNmplc" jsonschema:"라벨용채집지명"`
	LastUpdtDtm        string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	PrsrtStcd          string `json:"prsrtStcd" jsonschema:"보존상태"`
	SlistTpcd          string `json:"slistTpcd" jsonschema:"미소곤충구분"`
	SmplKindCd         string `json:"smplKindCd" jsonschema:"표본종류"`
	TorsoLngth         string `json:"torsoLngth" jsonschema:"몸통길이"`
	WingLngth          string `json:"wingLngth" jsonschema:"날개길이"`
	InsctGnrlNm        string `json:"insctGnrlNm" jsonschema:"국명(곤충명)"`
}

type insectSmplUnitListHandler struct {
	useCase inbound.InsectSmplUnitListUseCase
}

func addInsectSmplUnitListTool(server *mcp.Server, useCase inbound.InsectSmplUnitListUseCase) {
	handler := insectSmplUnitListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insect_resource_insect_smpl_unit_list",
		Description: "산림청 국립수목원 곤충표본 상세정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h insectSmplUnitListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input insectSmplUnitListInput) (*mcp.CallToolResult, insectSmplUnitListOutput, error) {
	result, err := h.useCase.InsectSmplUnitList(ctx, application.InsectSmplUnitListQuery{
		PageNo:          input.PageNo,
		NumOfRows:       input.NumOfRows,
		ReqInsctSpecsID: input.ReqInsctSpecsID,
	})
	if err != nil {
		return nil, insectSmplUnitListOutput{}, err
	}

	items := make([]insectSmplUnitListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = insectSmplUnitListItem{
			BspcsInsttNm:       item.BspcsInsttNm,
			ClarHaslvVal:       item.ClarHaslvVal,
			SmplCllcnDt:        item.SmplCllcnDt,
			GynndTpcd:          item.GynndTpcd,
			HbttTpcd:           item.HbttTpcd,
			InsctSmplNo:        item.InsctSmplNo,
			InsctSpecsID:       item.InsctSpecsID,
			InsctSpecsScnm:     item.InsctSpecsScnm,
			LabelUsgCllcnNmplc: item.LabelUsgCllcnNmplc,
			LastUpdtDtm:        item.LastUpdtDtm,
			PrsrtStcd:          item.PrsrtStcd,
			SlistTpcd:          item.SlistTpcd,
			SmplKindCd:         item.SmplKindCd,
			TorsoLngth:         item.TorsoLngth,
			WingLngth:          item.WingLngth,
			InsctGnrlNm:        item.InsctGnrlNm,
		}
	}

	return nil, insectSmplUnitListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
