package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application/port/inbound"
)

type scnmInfoInput struct {
	ReqFngsScnmID string `json:"reqFngsScnmId" jsonschema:"버섯 학명ID (scnmSearch 결과의 fngsScnmId)"`
}

type scnmInfoOutput struct {
	Item *scnmInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type scnmInfoItem struct {
	StpltScnmRltnCdNm string `json:"stpltScnmRltnCdNm" jsonschema:"버섯 학명의 정명/이명 구분"`
	FalmNm            string `json:"falmNm" jsonschema:"버섯 학명의 과명(Family name)"`
	FalnKorNm         string `json:"falnKorNm" jsonschema:"버섯 학명 과명의 과국명"`
	FngsEclgTpcdNm    string `json:"fngsEclgTpcdNm" jsonschema:"버섯 생태형"`
	FngsGnrlNm        string `json:"fngsGnrlNm" jsonschema:"버섯의 추천 국명(버섯명)"`
	FngsGnrlNm2       string `json:"fngsGnrlNm2" jsonschema:"버섯의 비추천 국명(버섯명)"`
	FngsPrpseTpcdNm   string `json:"fngsPrpseTpcdNm" jsonschema:"버섯의 식독 정보"`
	FngsScnm          string `json:"fngsScnm" jsonschema:"버섯 학명"`
	FngsScnmID        string `json:"fngsScnmId" jsonschema:"버섯 학명ID"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"버섯 학명 속명의 속국명"`
	GenusNm           string `json:"genusNm" jsonschema:"버섯 학명의 속명(Genus name)"`
	LastUpdtDtm       string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	OrdscLtrtrNm      string `json:"ordscLtrtrNm" jsonschema:"버섯 학명의 기재문 정보"`
	Rmrk              string `json:"rmrk" jsonschema:"비고"`
}

type scnmInfoHandler struct {
	useCase inbound.ScnmInfoUseCase
}

func addScnmInfoTool(server *mcp.Server, useCase inbound.ScnmInfoUseCase) {
	handler := scnmInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kfni_scnm_info",
		Description: "산림청 국립수목원 국가표준버섯목록의 버섯 학명 상세 정보를 조회합니다.",
	}, handler.handle)
}

func (h scnmInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input scnmInfoInput) (*mcp.CallToolResult, scnmInfoOutput, error) {
	result, err := h.useCase.ScnmInfo(ctx, application.ScnmInfoQuery{ReqFngsScnmID: input.ReqFngsScnmID})
	if err != nil {
		return nil, scnmInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, scnmInfoOutput{}, nil
	}

	item := result.Item
	return nil, scnmInfoOutput{Item: &scnmInfoItem{
		StpltScnmRltnCdNm: item.StpltScnmRltnCdNm,
		FalmNm:            item.FalmNm,
		FalnKorNm:         item.FalnKorNm,
		FngsEclgTpcdNm:    item.FngsEclgTpcdNm,
		FngsGnrlNm:        item.FngsGnrlNm,
		FngsGnrlNm2:       item.FngsGnrlNm2,
		FngsPrpseTpcdNm:   item.FngsPrpseTpcdNm,
		FngsScnm:          item.FngsScnm,
		FngsScnmID:        item.FngsScnmID,
		GenusKorNm:        item.GenusKorNm,
		GenusNm:           item.GenusNm,
		LastUpdtDtm:       item.LastUpdtDtm,
		OrdscLtrtrNm:      item.OrdscLtrtrNm,
		Rmrk:              item.Rmrk,
	}}, nil
}
