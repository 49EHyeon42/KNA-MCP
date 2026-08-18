package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kini/application/port/inbound"
)

type scnmInfoInput struct {
	ReqInsctScnmID string `json:"reqInsctScnmId" jsonschema:"조회하려는 곤충 학명ID (scnmSearch 결과의 insctScnmId)"`
}

type scnmInfoOutput struct {
	Item *scnmInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type scnmInfoItem struct {
	SuperFalmNm       string `json:"superFalmNm" jsonschema:"곤충 상과명(SuperFamily Name)"`
	ClassKorNm        string `json:"classKorNm" jsonschema:"곤충 강국명"`
	ClassNm           string `json:"classNm" jsonschema:"곤충 강명(Class Name)"`
	FalmKorNm         string `json:"falmKorNm" jsonschema:"곤충 과국명"`
	FalmNm            string `json:"falmNm" jsonschema:"곤충 과명(Family Name)"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"곤충 속국명"`
	GenusNm           string `json:"genusNm" jsonschema:"곤충 속명(Genus Name)"`
	InsctGnrlNm       string `json:"insctGnrlNm" jsonschema:"곤충 추천국명(곤충명)"`
	InsctGnrlNm2      string `json:"insctGnrlNm2" jsonschema:"곤충 비추천국명"`
	InsctScnmID       string `json:"insctScnmId" jsonschema:"곤충 학명ID"`
	InsctSpecsScnm    string `json:"insctSpecsScnm" jsonschema:"곤충 학명"`
	LastUpdtDtm       string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	OrdKorNm          string `json:"ordKorNm" jsonschema:"곤충 목국명"`
	OrdNm             string `json:"ordNm" jsonschema:"곤충 목명(Order Name)"`
	StpltScnmRltnCdNm string `json:"stpltScnmRltnCdNm" jsonschema:"곤충 학명 정명/이명 구분"`
	SubFalmKorNm      string `json:"subFalmKorNm" jsonschema:"곤충 아과국명"`
	SubFalmNm         string `json:"subFalmNm" jsonschema:"곤충 아과명(SubFamily Name)"`
	SuperFalmKorNm    string `json:"superFalmKorNm" jsonschema:"곤충 상과국명"`
}

type scnmInfoHandler struct {
	useCase inbound.ScnmInfoUseCase
}

func addScnmInfoTool(server *mcp.Server, useCase inbound.ScnmInfoUseCase) {
	handler := scnmInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kini_scnm_info",
		Description: "산림청 국립수목원 국가표준곤충목록의 곤충 학명 상세 정보를 조회합니다.",
	}, handler.handle)
}

func (h scnmInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input scnmInfoInput) (*mcp.CallToolResult, scnmInfoOutput, error) {
	result, err := h.useCase.ScnmInfo(ctx, application.ScnmInfoQuery{ReqInsctScnmID: input.ReqInsctScnmID})
	if err != nil {
		return nil, scnmInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, scnmInfoOutput{}, nil
	}

	item := result.Item
	return nil, scnmInfoOutput{Item: &scnmInfoItem{
		SuperFalmNm:       item.SuperFalmNm,
		ClassKorNm:        item.ClassKorNm,
		ClassNm:           item.ClassNm,
		FalmKorNm:         item.FalmKorNm,
		FalmNm:            item.FalmNm,
		GenusKorNm:        item.GenusKorNm,
		GenusNm:           item.GenusNm,
		InsctGnrlNm:       item.InsctGnrlNm,
		InsctGnrlNm2:      item.InsctGnrlNm2,
		InsctScnmID:       item.InsctScnmID,
		InsctSpecsScnm:    item.InsctSpecsScnm,
		LastUpdtDtm:       item.LastUpdtDtm,
		OrdKorNm:          item.OrdKorNm,
		OrdNm:             item.OrdNm,
		StpltScnmRltnCdNm: item.StpltScnmRltnCdNm,
		SubFalmKorNm:      item.SubFalmKorNm,
		SubFalmNm:         item.SubFalmNm,
		SuperFalmKorNm:    item.SuperFalmKorNm,
	}}, nil
}
