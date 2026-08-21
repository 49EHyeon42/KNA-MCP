package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/inbound"
)

type scnmInfoInput struct {
	ReqLchnScnmID string `json:"reqLchnScnmId" jsonschema:"지의류 학명ID (scnmSearch 결과의 lchnScnmId)"`
}

type scnmInfoOutput struct {
	Item *scnmInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type scnmInfoItem struct {
	StpltScnmRltnCdNm string `json:"stpltScnmRltnCdNm" jsonschema:"지의류 학명의 정명/이명 구분"`
	FalmNm            string `json:"falmNm" jsonschema:"지의류 학명 분류군의 과명(Family Name)"`
	FalnKorNm         string `json:"falnKorNm" jsonschema:"지의류 학명 분류군의 과국명"`
	GenusKorNm        string `json:"genusKorNm" jsonschema:"지의류 학명 분류군의 속국명"`
	GenusNm           string `json:"genusNm" jsonschema:"지의류 학명 분류군의 속명(Genus Name)"`
	LastUpdtDtm       string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	LchnGnrlNm        string `json:"lchnGnrlNm" jsonschema:"지의류 추천 국명(지의류명)"`
	LchnGnrlNm2       string `json:"lchnGnrlNm2" jsonschema:"지의류 비추천 국명"`
	LchnScnm          string `json:"lchnScnm" jsonschema:"지의류 학명"`
	LchnScnmID        string `json:"lchnScnmId" jsonschema:"지의류 학명ID"`
	OrdscLtrtrNm      string `json:"ordscLtrtrNm" jsonschema:"지의류 학명 기재문"`
	Rmrk              string `json:"rmrk" jsonschema:"비고"`
}

type scnmInfoHandler struct {
	useCase inbound.ScnmInfoUseCase
}

func addScnmInfoTool(server *mcp.Server, useCase inbound.ScnmInfoUseCase) {
	handler := scnmInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "klni_scnm_info",
		Description: "산림청 국립수목원 국가표준지의류목록의 지의류 학명 상세 정보를 조회합니다.",
	}, handler.handle)
}

func (h scnmInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input scnmInfoInput) (*mcp.CallToolResult, scnmInfoOutput, error) {
	result, err := h.useCase.ScnmInfo(ctx, application.ScnmInfoQuery{ReqLchnScnmID: input.ReqLchnScnmID})
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
		GenusKorNm:        item.GenusKorNm,
		GenusNm:           item.GenusNm,
		LastUpdtDtm:       item.LastUpdtDtm,
		LchnGnrlNm:        item.LchnGnrlNm,
		LchnGnrlNm2:       item.LchnGnrlNm2,
		LchnScnm:          item.LchnScnm,
		LchnScnmID:        item.LchnScnmID,
		OrdscLtrtrNm:      item.OrdscLtrtrNm,
		Rmrk:              item.Rmrk,
	}}, nil
}
