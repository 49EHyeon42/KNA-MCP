package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/inbound"
)

type alchnIlstrInfoInput struct {
	Q1 string `json:"q1" jsonschema:"조회키 (alchnIlstrSearch 결과의 lchnPilbkNo)"`
}

type alchnIlstrInfoOutput struct {
	Item *alchnIlstrInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type alchnIlstrInfoItem struct {
	Btnc         string `json:"btnc" jsonschema:"학명"`
	Cont1        string `json:"cont1" jsonschema:"영문설명"`
	Cont2        string `json:"cont2" jsonschema:"미사용"`
	Cont3        string `json:"cont3" jsonschema:"형태에 의한 분류"`
	Cont4        string `json:"cont4" jsonschema:"지의류형태(국문 종기술)"`
	Cont5        string `json:"cont5" jsonschema:"미사용"`
	Cont6        string `json:"cont6" jsonschema:"미사용"`
	Cont7        string `json:"cont7" jsonschema:"미사용"`
	Cont8        string `json:"cont8" jsonschema:"미사용"`
	Cont9        string `json:"cont9" jsonschema:"지의물질"`
	Cont10       string `json:"cont10" jsonschema:"분포"`
	Cont11       string `json:"cont11" jsonschema:"미사용"`
	Cont12       string `json:"cont12" jsonschema:"비고"`
	CprtCtnt     string `json:"cprtCtnt" jsonschema:"저작권"`
	EngNm        string `json:"engNm" jsonschema:"영문명"`
	FamilyKorNm  string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm     string `json:"familyNm" jsonschema:"과명"`
	FrstRgstnDtm string `json:"frstRgstnDtm" jsonschema:"최초등록일시"`
	GenusKorNm   string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm      string `json:"genusNm" jsonschema:"속명"`
	ImgURL       string `json:"imgUrl" jsonschema:"이미지URL"`
	JapNm        string `json:"japNm" jsonschema:"일어명"`
	LastUpdtDtm  string `json:"lastUpdtDtm" jsonschema:"최종수정일시"`
	LchnGnrlNm   string `json:"lchnGnrlNm" jsonschema:"국명"`
	LchnInfrpNm  string `json:"lchnInfrpNm" jsonschema:"종하명"`
	LchnPilbkNo  string `json:"lchnPilbkNo" jsonschema:"도감번호"`
	LchnScnmID   string `json:"lchnScnmId" jsonschema:"학명ID"`
	LchnTtnm     string `json:"lchnTtnm" jsonschema:"종소명"`
	PrkNm        string `json:"prkNm" jsonschema:"북한명"`
}

type alchnIlstrInfoHandler struct {
	useCase inbound.AlchnIlstrInfoUseCase
}

func addAlchnIlstrInfoTool(server *mcp.Server, useCase inbound.AlchnIlstrInfoUseCase) {
	handler := alchnIlstrInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "lchn_service_alchn_ilstr_info",
		Description: "산림청 국립수목원 지의류도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h alchnIlstrInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input alchnIlstrInfoInput) (*mcp.CallToolResult, alchnIlstrInfoOutput, error) {
	result, err := h.useCase.AlchnIlstrInfo(ctx, application.AlchnIlstrInfoQuery{Q1: input.Q1})
	if err != nil {
		return nil, alchnIlstrInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, alchnIlstrInfoOutput{}, nil
	}

	item := result.Item
	return nil, alchnIlstrInfoOutput{Item: &alchnIlstrInfoItem{
		Btnc:         item.Btnc,
		Cont1:        item.Cont1,
		Cont2:        item.Cont2,
		Cont3:        item.Cont3,
		Cont4:        item.Cont4,
		Cont5:        item.Cont5,
		Cont6:        item.Cont6,
		Cont7:        item.Cont7,
		Cont8:        item.Cont8,
		Cont9:        item.Cont9,
		Cont10:       item.Cont10,
		Cont11:       item.Cont11,
		Cont12:       item.Cont12,
		CprtCtnt:     item.CprtCtnt,
		EngNm:        item.EngNm,
		FamilyKorNm:  item.FamilyKorNm,
		FamilyNm:     item.FamilyNm,
		FrstRgstnDtm: item.FrstRgstnDtm,
		GenusKorNm:   item.GenusKorNm,
		GenusNm:      item.GenusNm,
		ImgURL:       item.ImgURL,
		JapNm:        item.JapNm,
		LastUpdtDtm:  item.LastUpdtDtm,
		LchnGnrlNm:   item.LchnGnrlNm,
		LchnInfrpNm:  item.LchnInfrpNm,
		LchnPilbkNo:  item.LchnPilbkNo,
		LchnScnmID:   item.LchnScnmID,
		LchnTtnm:     item.LchnTtnm,
		PrkNm:        item.PrkNm,
	}}, nil
}
