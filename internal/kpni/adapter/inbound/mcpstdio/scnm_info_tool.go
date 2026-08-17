package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/inbound"
)

type scnmInfoInput struct {
	ReqPlantScnmID string `json:"reqPlantScnmId" jsonschema:"검색하려는 식물 학명ID (scnmSearch 결과의 plantScnmId)"`
}

type scnmInfoOutput struct {
	Item *scnmInfoItem `json:"item" jsonschema:"상세 조회 결과"`
}

type scnmInfoItem struct {
	APGFalmKorNm       string `json:"apgFalmKorNm" jsonschema:"식물 학명 APG 분류군 과국명"`
	APGFalmNm          string `json:"apgFalmNm" jsonschema:"식물 학명 APG 분류군 과명(Family Name)"`
	BiogyNmTpcdNm      string `json:"biogyNmTpcdNm" jsonschema:"식물 학명 정명/이명/서명 등의 구분명"`
	CltvaYn            string `json:"cltvaYn" jsonschema:"재배여부"`
	EclgDstrbYn        string `json:"eclgDstrbYn" jsonschema:"생태계 교란종 여부"`
	ExtcCncrnsYn       string `json:"extcCncrnsYn" jsonschema:"외래화 우려 여부"`
	ExtcPlantCdNm      string `json:"extcPlantCdNm" jsonschema:"외래 식물 구분명"`
	ExtcPlantYn        string `json:"extcPlantYn" jsonschema:"침입 외래 식물 여부"`
	FalmKorNm          string `json:"falmKorNm" jsonschema:"식물 학명 분류군 과국명"`
	FalmNm             string `json:"falmNm" jsonschema:"식물 학명 분류군 과명(Family Name)"`
	GenusKorNm         string `json:"genusKorNm" jsonschema:"식물 학명 분류군 속국명"`
	GenusNm            string `json:"genusNm" jsonschema:"식물 학명 분류군 속명(Genus Name)"`
	LtrtrInfrmNm       string `json:"ltrtrInfrmNm" jsonschema:"학명 기재문"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm" jsonschema:"식물 번식 구분 형태"`
	PlantChnNm         string `json:"plantChnNm" jsonschema:"식물 중국명"`
	PlantEngNm         string `json:"plantEngNm" jsonschema:"식물 영문명"`
	PlantGnrlNm        string `json:"plantGnrlNm" jsonschema:"식물 추천 국명"`
	PlantGnrlNm2       string `json:"plantGnrlNm2" jsonschema:"식물 비추천 국명"`
	PlantJpnNm         string `json:"plantJpnNm" jsonschema:"식물 일본명"`
	PlantScnmID        string `json:"plantScnmId" jsonschema:"식물 학명ID"`
	PlantSpecsScnm     string `json:"plantSpecsScnm" jsonschema:"식물 학명"`
	RareTpcdNm         string `json:"rareTpcdNm" jsonschema:"희귀식물 분류명"`
	RelPlantSpecsScnm  string `json:"relPlantSpecsScnm" jsonschema:"연관 학명"`
	RelScnmTpcdNm      string `json:"relScnmTpcdNm" jsonschema:"연관 학명 구분명"`
	Rmrk               string `json:"rmrk" jsonschema:"비고"`
	RrnssPlantYn       string `json:"rrnssPlantYn" jsonschema:"희귀식물 여부"`
	SpcltPlantCdNm     string `json:"spcltPlantCdNm" jsonschema:"특산식물 분류명"`
}

type scnmInfoHandler struct {
	useCase inbound.ScnmInfoUseCase
}

func addScnmInfoTool(server *mcp.Server, useCase inbound.ScnmInfoUseCase) {
	handler := scnmInfoHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "kpni_scnm_info",
		Description: "산림청 국립수목원 국가표준식물목록의 식물 학명 상세 정보를 조회합니다.",
	}, handler.handle)
}

func (h scnmInfoHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input scnmInfoInput) (*mcp.CallToolResult, scnmInfoOutput, error) {
	result, err := h.useCase.ScnmInfo(ctx, application.ScnmInfoQuery{ReqPlantScnmID: input.ReqPlantScnmID})
	if err != nil {
		return nil, scnmInfoOutput{}, err
	}
	if result.Item == nil {
		return nil, scnmInfoOutput{}, nil
	}

	item := result.Item
	return nil, scnmInfoOutput{Item: &scnmInfoItem{
		APGFalmKorNm:       item.APGFalmKorNm,
		APGFalmNm:          item.APGFalmNm,
		BiogyNmTpcdNm:      item.BiogyNmTpcdNm,
		CltvaYn:            item.CltvaYn,
		EclgDstrbYn:        item.EclgDstrbYn,
		ExtcCncrnsYn:       item.ExtcCncrnsYn,
		ExtcPlantCdNm:      item.ExtcPlantCdNm,
		ExtcPlantYn:        item.ExtcPlantYn,
		FalmKorNm:          item.FalmKorNm,
		FalmNm:             item.FalmNm,
		GenusKorNm:         item.GenusKorNm,
		GenusNm:            item.GenusNm,
		LtrtrInfrmNm:       item.LtrtrInfrmNm,
		PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
		PlantChnNm:         item.PlantChnNm,
		PlantEngNm:         item.PlantEngNm,
		PlantGnrlNm:        item.PlantGnrlNm,
		PlantGnrlNm2:       item.PlantGnrlNm2,
		PlantJpnNm:         item.PlantJpnNm,
		PlantScnmID:        item.PlantScnmID,
		PlantSpecsScnm:     item.PlantSpecsScnm,
		RareTpcdNm:         item.RareTpcdNm,
		RelPlantSpecsScnm:  item.RelPlantSpecsScnm,
		RelScnmTpcdNm:      item.RelScnmTpcdNm,
		Rmrk:               item.Rmrk,
		RrnssPlantYn:       item.RrnssPlantYn,
		SpcltPlantCdNm:     item.SpcltPlantCdNm,
	}}, nil
}
