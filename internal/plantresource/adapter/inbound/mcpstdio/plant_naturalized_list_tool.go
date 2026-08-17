package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantNaturalizedListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"검색할 식물의 국명 또는 학명"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type plantNaturalizedListOutput struct {
	Items      []plantNaturalizedListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                        `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                        `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                        `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantNaturalizedListItem struct {
	AgpFamilyNm        string `json:"agpFamilyNm" jsonschema:"APG과명"`
	APGFamilyKorNm     string `json:"apgFamilyKorNm" jsonschema:"APG과국명"`
	BlprdEnmnt         string `json:"blprdEnmnt" jsonschema:"개화기종료일"`
	BlprdStmnt         string `json:"blprdStmnt" jsonschema:"개화기시작일"`
	DistrAraDscrt      string `json:"distrAraDscrt" jsonschema:"분포지역"`
	EclgDstrbYn        string `json:"eclgDstrbYn" jsonschema:"생태계교란종여부"`
	ExtcPlantCdNm      string `json:"extcPlantCdNm" jsonschema:"외래식물구분"`
	FamilyKorNm        string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm           string `json:"familyNm" jsonschema:"과명"`
	FrtTpcdNm          string `json:"frtTpcdNm" jsonschema:"열매구분"`
	LastUpdtDtm        string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	NtldgTpcdNm        string `json:"ntldgTpcdNm" jsonschema:"귀화도구분"`
	NtrlzEraTpcdNm     string `json:"ntrlzEraTpcdNm" jsonschema:"유입시기구분"`
	OrplcNm            string `json:"orplcNm" jsonschema:"원산지"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm" jsonschema:"식물번식형태"`
	PlantDistrGrcd     string `json:"plantDistrGrcd" jsonschema:"외래식물확산등급"`
	PlantDistrQntt     string `json:"plantDistrQntt" jsonschema:"식물분포지역수량"`
	PlantDistrQnttGrcd string `json:"plantDistrQnttGrcd" jsonschema:"식물분포지역수량등급"`
	PlantEngNm         string `json:"plantEngNm" jsonschema:"영문명"`
	PlantGnrlNm        string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantJpnNm         string `json:"plantJpnNm" jsonschema:"일본명"`
	PlantLfcclTpcdNm   string `json:"plantLfcclTpcdNm" jsonschema:"식물생활사"`
	PlantSpecsScnm     string `json:"plantSpecsScnm" jsonschema:"학명"`
}

type plantNaturalizedListHandler struct {
	useCase inbound.PlantNaturalizedListUseCase
}

func addPlantNaturalizedListTool(server *mcp.Server, useCase inbound.PlantNaturalizedListUseCase) {
	handler := plantNaturalizedListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_naturalized_list",
		Description: "산림청 국립수목원 외래식물정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h plantNaturalizedListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantNaturalizedListInput) (*mcp.CallToolResult, plantNaturalizedListOutput, error) {
	result, err := h.useCase.PlantNaturalizedList(ctx, application.PlantNaturalizedListQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantNaturalizedListOutput{}, err
	}

	items := make([]plantNaturalizedListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantNaturalizedListItem{
			AgpFamilyNm:        item.AgpFamilyNm,
			APGFamilyKorNm:     item.APGFamilyKorNm,
			BlprdEnmnt:         item.BlprdEnmnt,
			BlprdStmnt:         item.BlprdStmnt,
			DistrAraDscrt:      item.DistrAraDscrt,
			EclgDstrbYn:        item.EclgDstrbYn,
			ExtcPlantCdNm:      item.ExtcPlantCdNm,
			FamilyKorNm:        item.FamilyKorNm,
			FamilyNm:           item.FamilyNm,
			FrtTpcdNm:          item.FrtTpcdNm,
			LastUpdtDtm:        item.LastUpdtDtm,
			NtldgTpcdNm:        item.NtldgTpcdNm,
			NtrlzEraTpcdNm:     item.NtrlzEraTpcdNm,
			OrplcNm:            item.OrplcNm,
			PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
			PlantDistrGrcd:     item.PlantDistrGrcd,
			PlantDistrQntt:     item.PlantDistrQntt,
			PlantDistrQnttGrcd: item.PlantDistrQnttGrcd,
			PlantEngNm:         item.PlantEngNm,
			PlantGnrlNm:        item.PlantGnrlNm,
			PlantJpnNm:         item.PlantJpnNm,
			PlantLfcclTpcdNm:   item.PlantLfcclTpcdNm,
			PlantSpecsScnm:     item.PlantSpecsScnm,
		}
	}

	return nil, plantNaturalizedListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
