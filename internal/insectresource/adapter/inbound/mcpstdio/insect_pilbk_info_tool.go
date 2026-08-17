package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/inbound"
)

type insectPilbkInfoInput struct {
	ReqInsctPilbkNo string `json:"reqInsctPilbkNo" jsonschema:"조회할 곤충 도감번호 (insectPilbkSearch 결과의 insctPilbkNo)"`
}

type insectPilbkInfoOutput struct {
	Item *insectPilbkInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type insectPilbkInfoItem struct {
	EcoDsrct         string `json:"ecoDsrct" jsonschema:"생태"`
	EggDsrct         string `json:"eggDsrct" jsonschema:"알"`
	EmrgcCnt         string `json:"emrgcCnt" jsonschema:"출현수"`
	EmrgcEraDscrt    string `json:"emrgcEraDscrt" jsonschema:"출현시기설명"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	FemaleDsrct      string `json:"femaleDsrct" jsonschema:"성충(암)"`
	GenusKorNm       string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm          string `json:"genusNm" jsonschema:"속명"`
	GnrlDsrct        string `json:"gnrlDsrct" jsonschema:"일반 특징"`
	HabitDsrct       string `json:"habitDsrct" jsonschema:"습성"`
	InsctEngNm       string `json:"insctEngNm" jsonschema:"영문명"`
	InsctGnrlNm      string `json:"insctGnrlNm" jsonschema:"국명(곤충명)"`
	InsctPilbkNo     string `json:"insctPilbkNo" jsonschema:"곤충도감번호"`
	InsctSpecsScnm   string `json:"insctSpecsScnm" jsonschema:"학명"`
	LarvaDsrct       string `json:"larvaDsrct" jsonschema:"유충"`
	LastUpdtDtm      string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	MaleDsrct        string `json:"maleDsrct" jsonschema:"성충(수)"`
	MnmmOccrrCnt     string `json:"mnmmOccrrCnt" jsonschema:"최소발생수"`
	MxmmOccrrCnt     string `json:"mxmmOccrrCnt" jsonschema:"최대발생수"`
	OrdKorNm         string `json:"ordKorNm" jsonschema:"목국명"`
	OrdNm            string `json:"ordNm" jsonschema:"목명"`
	PestDsrct        string `json:"pestDsrct" jsonschema:"방제법"`
	PupaDsrct        string `json:"pupaDsrct" jsonschema:"번데기"`
	ReferDsrct       string `json:"referDsrct" jsonschema:"참고사항"`
	SubFamilyKorNm   string `json:"subFamilyKorNm" jsonschema:"아과국명"`
	SubFamilyNm      string `json:"subFamilyNm" jsonschema:"아과명"`
	SuperFamilyKorNm string `json:"superFamilyKorNm" jsonschema:"상과국명"`
	SuperFamilyNm    string `json:"superFamilyNm" jsonschema:"상과명"`
	WinterDsrct      string `json:"winterDsrct" jsonschema:"월동"`
}

type insectPilbkInfoHandler struct {
	useCase inbound.InsectPilbkInfoUseCase
}

func addInsectPilbkInfoTool(server *mcp.Server, useCase inbound.InsectPilbkInfoUseCase) {
	handler := insectPilbkInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insect_resource_insect_pilbk_info",
		Description: "산림청 국립수목원 곤충도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h insectPilbkInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input insectPilbkInfoInput) (*mcp.CallToolResult, insectPilbkInfoOutput, error) {
	result, err := h.useCase.InsectPilbkInfo(ctx, application.InsectPilbkInfoQuery{ReqInsctPilbkNo: input.ReqInsctPilbkNo})
	if err != nil {
		return nil, insectPilbkInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, insectPilbkInfoOutput{}, nil
	}

	item := result.Item
	return nil, insectPilbkInfoOutput{Item: &insectPilbkInfoItem{
		EcoDsrct:         item.EcoDsrct,
		EggDsrct:         item.EggDsrct,
		EmrgcCnt:         item.EmrgcCnt,
		EmrgcEraDscrt:    item.EmrgcEraDscrt,
		FamilyKorNm:      item.FamilyKorNm,
		FamilyNm:         item.FamilyNm,
		FemaleDsrct:      item.FemaleDsrct,
		GenusKorNm:       item.GenusKorNm,
		GenusNm:          item.GenusNm,
		GnrlDsrct:        item.GnrlDsrct,
		HabitDsrct:       item.HabitDsrct,
		InsctEngNm:       item.InsctEngNm,
		InsctGnrlNm:      item.InsctGnrlNm,
		InsctPilbkNo:     item.InsctPilbkNo,
		InsctSpecsScnm:   item.InsctSpecsScnm,
		LarvaDsrct:       item.LarvaDsrct,
		LastUpdtDtm:      item.LastUpdtDtm,
		MaleDsrct:        item.MaleDsrct,
		MnmmOccrrCnt:     item.MnmmOccrrCnt,
		MxmmOccrrCnt:     item.MxmmOccrrCnt,
		OrdKorNm:         item.OrdKorNm,
		OrdNm:            item.OrdNm,
		PestDsrct:        item.PestDsrct,
		PupaDsrct:        item.PupaDsrct,
		ReferDsrct:       item.ReferDsrct,
		SubFamilyKorNm:   item.SubFamilyKorNm,
		SubFamilyNm:      item.SubFamilyNm,
		SuperFamilyKorNm: item.SuperFamilyKorNm,
		SuperFamilyNm:    item.SuperFamilyNm,
		WinterDsrct:      item.WinterDsrct,
	}}, nil
}
