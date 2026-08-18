package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/inbound"
)

type childPilbkInfoInput struct {
	ReqChildLvbngPilbkNo string `json:"reqChildLvbngPilbkNo" jsonschema:"도감번호(childLvbngPilbkNo) (childPilbkSearch 결과의 childLvbngPilbkNo)"`
}

type childPilbkInfoOutput struct {
	Item *childPilbkInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type childPilbkInfoItem struct {
	BiogyNm           string `json:"biogyNm" jsonschema:"생물학명"`
	ChildLvbngPilbkNo string `json:"childLvbngPilbkNo" jsonschema:"어린이생물도감번호"`
	ExtrmCrss         string `json:"extrmCrss" jsonschema:"멸종위기종구분"`
	FamilyKorNm       string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm          string `json:"familyNm" jsonschema:"과명"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm           string `json:"genusNm" jsonschema:"속명"`
	HbttFieldYn       string `json:"hbttFieldYn" jsonschema:"서식지 들 여부"`
	HbttFrestYn       string `json:"hbttFrestYn" jsonschema:"서식지 숲 여부"`
	HbttRiverYn       string `json:"hbttRiverYn" jsonschema:"서식지 강 여부"`
	LvbngDscrt        string `json:"lvbngDscrt" jsonschema:"생물설명"`
	LvbngTpcdNm       string `json:"lvbngTpcdNm" jsonschema:"생물분류"`
	LvngKrlngNm       string `json:"lvngKrlngNm" jsonschema:"생물국명"`
	PrtctSpecsTpcdNm  string `json:"prtctSpecsTpcdNm" jsonschema:"보호종구분"`
}

type childPilbkInfoHandler struct {
	useCase inbound.ChildPilbkInfoUseCase
}

func addChildPilbkInfoTool(server *mcp.Server, useCase inbound.ChildPilbkInfoUseCase) {
	handler := childPilbkInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "child_service_child_pilbk_info",
		Description: "산림청 국립수목원 어린이생물도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h childPilbkInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input childPilbkInfoInput) (*mcp.CallToolResult, childPilbkInfoOutput, error) {
	result, err := h.useCase.ChildPilbkInfo(ctx, application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: input.ReqChildLvbngPilbkNo})
	if err != nil {
		return nil, childPilbkInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, childPilbkInfoOutput{}, nil
	}

	item := result.Item
	return nil, childPilbkInfoOutput{Item: &childPilbkInfoItem{
		BiogyNm:           item.BiogyNm,
		ChildLvbngPilbkNo: item.ChildLvbngPilbkNo,
		ExtrmCrss:         item.ExtrmCrss,
		FamilyKorNm:       item.FamilyKorNm,
		FamilyNm:          item.FamilyNm,
		GenusKorNm:        item.GenusKorNm,
		GenusNm:           item.GenusNm,
		HbttFieldYn:       item.HbttFieldYn,
		HbttFrestYn:       item.HbttFrestYn,
		HbttRiverYn:       item.HbttRiverYn,
		LvbngDscrt:        item.LvbngDscrt,
		LvbngTpcdNm:       item.LvbngTpcdNm,
		LvngKrlngNm:       item.LvngKrlngNm,
		PrtctSpecsTpcdNm:  item.PrtctSpecsTpcdNm,
	}}, nil
}
