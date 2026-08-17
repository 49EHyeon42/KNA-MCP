package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
)

type fngsPilbkInfoInput struct {
	ReqFngsPilbkNo string `json:"reqFngsPilbkNo" jsonschema:"검색할 버섯도감의 버섯도감번호 (fngsPilbkSearch 결과의 fngsPilbkNo)"`
}

type fngsPilbkInfoOutput struct {
	Item *fngsPilbkInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type fngsPilbkInfoItem struct {
	MshrmColorCdNm      string `json:"mshrmColorCdNm" jsonschema:"버섯색상"`
	CrpphFomTpcdNm      string `json:"crpphFomTpcdNm" jsonschema:"자실체형태"`
	FamilyKorNm         string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm            string `json:"familyNm" jsonschema:"과명"`
	FngsEclgTpcdNm      string `json:"fngsEclgTpcdNm" jsonschema:"버섯생태형"`
	FngsGnrlNm          string `json:"fngsGnrlNm" jsonschema:"국명(버섯명)"`
	FngsPilbkNo         string `json:"fngsPilbkNo" jsonschema:"버섯도감번호"`
	FngsPrpseTpcdNm     string `json:"fngsPrpseTpcdNm" jsonschema:"버섯용도"`
	FngsScnm            string `json:"fngsScnm" jsonschema:"학명"`
	GenusKorNm          string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm             string `json:"genusNm" jsonschema:"속명"`
	GrwEvrntDesc        string `json:"grwEvrntDesc" jsonschema:"발생장소"`
	LastUpdtDtm         string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	MicroShpe           string `json:"microShpe" jsonschema:"현미경적 특징"`
	MshrmTpcdNm         string `json:"mshrmTpcdNm" jsonschema:"버섯구분"`
	OccrrSsnNm          string `json:"occrrSsnNm" jsonschema:"발생계절"`
	RsrcActoClsscTpcdNm string `json:"rsrcActoClsscTpcdNm" jsonschema:"자원분류"`
	Shpe                string `json:"shpe" jsonschema:"외부 형태적 특징"`
}

type fngsPilbkInfoHandler struct {
	useCase inbound.FngsPilbkInfoUseCase
}

func addFngsPilbkInfoTool(server *mcp.Server, useCase inbound.FngsPilbkInfoUseCase) {
	handler := fngsPilbkInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fungi_resource_fngs_pilbk_info",
		Description: "산림청 국립수목원 버섯도감 상세정보를 조회합니다.",
	}, handler.handle)
}

func (h fngsPilbkInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input fngsPilbkInfoInput) (*mcp.CallToolResult, fngsPilbkInfoOutput, error) {
	result, err := h.useCase.FngsPilbkInfo(ctx, application.FngsPilbkInfoQuery{ReqFngsPilbkNo: input.ReqFngsPilbkNo})
	if err != nil {
		return nil, fngsPilbkInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, fngsPilbkInfoOutput{}, nil
	}

	item := result.Item
	return nil, fngsPilbkInfoOutput{Item: &fngsPilbkInfoItem{
		MshrmColorCdNm:      item.MshrmColorCdNm,
		CrpphFomTpcdNm:      item.CrpphFomTpcdNm,
		FamilyKorNm:         item.FamilyKorNm,
		FamilyNm:            item.FamilyNm,
		FngsEclgTpcdNm:      item.FngsEclgTpcdNm,
		FngsGnrlNm:          item.FngsGnrlNm,
		FngsPilbkNo:         item.FngsPilbkNo,
		FngsPrpseTpcdNm:     item.FngsPrpseTpcdNm,
		FngsScnm:            item.FngsScnm,
		GenusKorNm:          item.GenusKorNm,
		GenusNm:             item.GenusNm,
		GrwEvrntDesc:        item.GrwEvrntDesc,
		LastUpdtDtm:         item.LastUpdtDtm,
		MicroShpe:           item.MicroShpe,
		MshrmTpcdNm:         item.MshrmTpcdNm,
		OccrrSsnNm:          item.OccrrSsnNm,
		RsrcActoClsscTpcdNm: item.RsrcActoClsscTpcdNm,
		Shpe:                item.Shpe,
	}}, nil
}
