package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
)

type entogSpcmInfoInput struct {
	Q1 string `json:"q1" jsonschema:"조회키 (entogSpcmSearch 결과의 entogSmplNo)"`
}

type entogSpcmInfoOutput struct {
	Item *entogSpcmInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type entogSpcmInfoItem struct {
	Btnc               string `json:"btnc" jsonschema:"학명"`
	ChnNm              string `json:"chnNm" jsonschema:"중국명"`
	ClarHaslvVal       string `json:"clarHaslvVal" jsonschema:"채집지해발고도"`
	ClctDyDesc         string `json:"clctDyDesc" jsonschema:"채집일"`
	CprtCtnt           string `json:"cprtCtnt" jsonschema:"저작권"`
	EngNm              string `json:"engNm" jsonschema:"영문명"`
	EntogGnrlNm        string `json:"entogGnrlNm" jsonschema:"국명"`
	EntogPilbkNo       string `json:"entogPilbkNo" jsonschema:"도감번호"`
	EntogSmplNo        string `json:"entogSmplNo" jsonschema:"표본번호"`
	FamilyKorNm        string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm           string `json:"familyNm" jsonschema:"과명"`
	FrstRgstnDtm       string `json:"frstRgstnDtm" jsonschema:"최초등록일"`
	GenusKorNm         string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm            string `json:"genusNm" jsonschema:"속명"`
	ImgURL             string `json:"imgUrl" jsonschema:"이미지URL"`
	JapNm              string `json:"japNm" jsonschema:"일본명"`
	LabelUsgCllcnNmplc string `json:"labelUsgCllcnNmplc" jsonschema:"라벨용채집지명"`
	LastUpdtDtm        string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	OrdKorNm           string `json:"ordKorNm" jsonschema:"목국명"`
	OrdNm              string `json:"ordNm" jsonschema:"목명"`
	PrkNm              string `json:"prkNm" jsonschema:"북한명"`
	TorsoLngth         string `json:"torsoLngth" jsonschema:"몸통길이"`
	WingLngth          string `json:"wingLngth" jsonschema:"날개길이"`
}

type entogSpcmInfoHandler struct {
	useCase inbound.EntogSpcmInfoUseCase
}

func addEntogSpcmInfoTool(server *mcp.Server, useCase inbound.EntogSpcmInfoUseCase) {
	handler := entogSpcmInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "entog_service_entog_spcm_info",
		Description: "산림청 국립수목원 내구강표본 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h entogSpcmInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input entogSpcmInfoInput) (*mcp.CallToolResult, entogSpcmInfoOutput, error) {
	result, err := h.useCase.EntogSpcmInfo(ctx, application.EntogSpcmInfoQuery{Q1: input.Q1})
	if err != nil {
		return nil, entogSpcmInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, entogSpcmInfoOutput{}, nil
	}

	item := result.Item
	return nil, entogSpcmInfoOutput{Item: &entogSpcmInfoItem{
		Btnc:               item.Btnc,
		ChnNm:              item.ChnNm,
		ClarHaslvVal:       item.ClarHaslvVal,
		ClctDyDesc:         item.ClctDyDesc,
		CprtCtnt:           item.CprtCtnt,
		EngNm:              item.EngNm,
		EntogGnrlNm:        item.EntogGnrlNm,
		EntogPilbkNo:       item.EntogPilbkNo,
		EntogSmplNo:        item.EntogSmplNo,
		FamilyKorNm:        item.FamilyKorNm,
		FamilyNm:           item.FamilyNm,
		FrstRgstnDtm:       item.FrstRgstnDtm,
		GenusKorNm:         item.GenusKorNm,
		GenusNm:            item.GenusNm,
		ImgURL:             item.ImgURL,
		JapNm:              item.JapNm,
		LabelUsgCllcnNmplc: item.LabelUsgCllcnNmplc,
		LastUpdtDtm:        item.LastUpdtDtm,
		OrdKorNm:           item.OrdKorNm,
		OrdNm:              item.OrdNm,
		PrkNm:              item.PrkNm,
		TorsoLngth:         item.TorsoLngth,
		WingLngth:          item.WingLngth,
	}}, nil
}
