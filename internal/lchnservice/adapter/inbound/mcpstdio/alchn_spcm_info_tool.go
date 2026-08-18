package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
)

type alchnSpcmInfoInput struct {
	Q1 string `json:"q1" jsonschema:"조회키 (alchnSpcmSearch 결과의 lchnSmplNo)"`
}

type alchnSpcmInfoOutput struct {
	Item *alchnSpcmInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type alchnSpcmInfoItem struct {
	Btnc          string `json:"btnc" jsonschema:"학명"`
	ClarDtlDscrt  string `json:"clarDtlDscrt" jsonschema:"채집지상세설명"`
	CllcrNm       string `json:"cllcrNm" jsonschema:"채집자명"`
	CltrNm        string `json:"cltrNm" jsonschema:"채집자그룹멤버명"`
	CprtCtnt      string `json:"cprtCtnt" jsonschema:"저작권"`
	EngNm         string `json:"engNm" jsonschema:"영문명"`
	ExmneNm       string `json:"exmneNm" jsonschema:"조사자명"`
	FamilyKorNm   string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm      string `json:"familyNm" jsonschema:"과명"`
	FrstRgstnDtm  string `json:"frstRgstnDtm" jsonschema:"최초등록일시"`
	GenusKorNm    string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm       string `json:"genusNm" jsonschema:"속명"`
	Grdnt         string `json:"grdnt" jsonschema:"경사도"`
	HaslvVal      string `json:"haslvVal" jsonschema:"해발고도값"`
	HbttChrcrCont string `json:"hbttChrcrCont" jsonschema:"기물설명"`
	ImgURL        string `json:"imgUrl" jsonschema:"이미지URL"`
	InsttSmplNo   string `json:"insttSmplNo" jsonschema:"기관표본번호"`
	JapNm         string `json:"japNm" jsonschema:"일어명"`
	LastUpdtDtm   string `json:"lastUpdtDtm" jsonschema:"최종수정일시"`
	LchnGnrlNm    string `json:"lchnGnrlNm" jsonschema:"국명"`
	LchnScnmID    string `json:"lchnScnmId" jsonschema:"학명ID"`
	LchnSmplNo    string `json:"lchnSmplNo" jsonschema:"표본번호"`
	OrbrnCd       string `json:"orbrnCd" jsonschema:"방위코드"`
	PrkNm         string `json:"prkNm" jsonschema:"북한명"`
	SmplCllcnDt   string `json:"smplCllcnDt" jsonschema:"표본채집일자"`
}

type alchnSpcmInfoHandler struct {
	useCase inbound.AlchnSpcmInfoUseCase
}

func addAlchnSpcmInfoTool(server *mcp.Server, useCase inbound.AlchnSpcmInfoUseCase) {
	handler := alchnSpcmInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "lchn_service_alchn_spcm_info",
		Description: "산림청 국립수목원 지의류표본 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h alchnSpcmInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input alchnSpcmInfoInput) (*mcp.CallToolResult, alchnSpcmInfoOutput, error) {
	result, err := h.useCase.AlchnSpcmInfo(ctx, application.AlchnSpcmInfoQuery{Q1: input.Q1})
	if err != nil {
		return nil, alchnSpcmInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, alchnSpcmInfoOutput{}, nil
	}

	item := result.Item
	return nil, alchnSpcmInfoOutput{Item: &alchnSpcmInfoItem{
		Btnc:          item.Btnc,
		ClarDtlDscrt:  item.ClarDtlDscrt,
		CllcrNm:       item.CllcrNm,
		CltrNm:        item.CltrNm,
		CprtCtnt:      item.CprtCtnt,
		EngNm:         item.EngNm,
		ExmneNm:       item.ExmneNm,
		FamilyKorNm:   item.FamilyKorNm,
		FamilyNm:      item.FamilyNm,
		FrstRgstnDtm:  item.FrstRgstnDtm,
		GenusKorNm:    item.GenusKorNm,
		GenusNm:       item.GenusNm,
		Grdnt:         item.Grdnt,
		HaslvVal:      item.HaslvVal,
		HbttChrcrCont: item.HbttChrcrCont,
		ImgURL:        item.ImgURL,
		InsttSmplNo:   item.InsttSmplNo,
		JapNm:         item.JapNm,
		LastUpdtDtm:   item.LastUpdtDtm,
		LchnGnrlNm:    item.LchnGnrlNm,
		LchnScnmID:    item.LchnScnmID,
		LchnSmplNo:    item.LchnSmplNo,
		OrbrnCd:       item.OrbrnCd,
		PrkNm:         item.PrkNm,
		SmplCllcnDt:   item.SmplCllcnDt,
	}}, nil
}
