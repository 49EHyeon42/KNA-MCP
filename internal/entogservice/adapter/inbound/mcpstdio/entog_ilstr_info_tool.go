package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/inbound"
)

type entogIlstrInfoInput struct {
	Q1 string `json:"q1" jsonschema:"조회키 (entogIlstrSearch 결과의 entogPilbkNo)"`
}

type entogIlstrInfoOutput struct {
	Item *entogIlstrInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type entogIlstrInfoItem struct {
	Btnc             string `json:"btnc" jsonschema:"학명"`
	Cont1            string `json:"cont1" jsonschema:"일반특징"`
	Cont2            string `json:"cont2" jsonschema:"성충(수)"`
	Cont3            string `json:"cont3" jsonschema:"성충(암)"`
	Cont4            string `json:"cont4" jsonschema:"번데기"`
	Cont5            string `json:"cont5" jsonschema:"유충"`
	Cont6            string `json:"cont6" jsonschema:"참고사항"`
	Cont7            string `json:"cont7" jsonschema:"생태"`
	Cont8            string `json:"cont8" jsonschema:"습성"`
	Cont9            string `json:"cont9" jsonschema:"월동"`
	Cont10           string `json:"cont10" jsonschema:"방제법"`
	Cont11           string `json:"cont11" jsonschema:"알"`
	CprtCtnt         string `json:"cprtCtnt" jsonschema:"저작권"`
	EmrgcCnt         string `json:"emrgcCnt" jsonschema:"출현수"`
	EmrgcEraDscrt    string `json:"emrgcEraDscrt" jsonschema:"출현시기설명"`
	EntogAthrNm      string `json:"entogAthrNm" jsonschema:"명명자명"`
	EntogEngNm       string `json:"entogEngNm" jsonschema:"영문명"`
	EntogOfnmKrlngNm string `json:"entogOfnmKrlngNm" jsonschema:"국명"`
	EntogPilbkNo     string `json:"entogPilbkNo" jsonschema:"도감번호"`
	EntogSpecsNm     string `json:"entogSpecsNm" jsonschema:"종소명"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	GenusKorNm       string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm          string `json:"genusNm" jsonschema:"속명"`
	ImgURL           string `json:"imgUrl" jsonschema:"이미지URL"`
	MnmmOccrrCnt     string `json:"mnmmOccrrCnt" jsonschema:"최소발생횟수"`
	MxmmOccrrCnt     string `json:"mxmmOccrrCnt" jsonschema:"최대발생횟수"`
	OrdKorNm         string `json:"ordKorNm" jsonschema:"목국명"`
	OrdNm            string `json:"ordNm" jsonschema:"목명"`
}

type entogIlstrInfoHandler struct {
	useCase inbound.EntogIlstrInfoUseCase
}

func addEntogIlstrInfoTool(server *mcp.Server, useCase inbound.EntogIlstrInfoUseCase) {
	handler := entogIlstrInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "entog_service_entog_ilstr_info",
		Description: "산림청 국립수목원 내구강도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h entogIlstrInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input entogIlstrInfoInput) (*mcp.CallToolResult, entogIlstrInfoOutput, error) {
	result, err := h.useCase.EntogIlstrInfo(ctx, application.EntogIlstrInfoQuery{Q1: input.Q1})
	if err != nil {
		return nil, entogIlstrInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, entogIlstrInfoOutput{}, nil
	}

	item := result.Item
	return nil, entogIlstrInfoOutput{Item: &entogIlstrInfoItem{
		Btnc:             item.Btnc,
		Cont1:            item.Cont1,
		Cont2:            item.Cont2,
		Cont3:            item.Cont3,
		Cont4:            item.Cont4,
		Cont5:            item.Cont5,
		Cont6:            item.Cont6,
		Cont7:            item.Cont7,
		Cont8:            item.Cont8,
		Cont9:            item.Cont9,
		Cont10:           item.Cont10,
		Cont11:           item.Cont11,
		CprtCtnt:         item.CprtCtnt,
		EmrgcCnt:         item.EmrgcCnt,
		EmrgcEraDscrt:    item.EmrgcEraDscrt,
		EntogAthrNm:      item.EntogAthrNm,
		EntogEngNm:       item.EntogEngNm,
		EntogOfnmKrlngNm: item.EntogOfnmKrlngNm,
		EntogPilbkNo:     item.EntogPilbkNo,
		EntogSpecsNm:     item.EntogSpecsNm,
		FamilyKorNm:      item.FamilyKorNm,
		FamilyNm:         item.FamilyNm,
		GenusKorNm:       item.GenusKorNm,
		GenusNm:          item.GenusNm,
		ImgURL:           item.ImgURL,
		MnmmOccrrCnt:     item.MnmmOccrrCnt,
		MxmmOccrrCnt:     item.MxmmOccrrCnt,
		OrdKorNm:         item.OrdKorNm,
		OrdNm:            item.OrdNm,
	}}, nil
}
