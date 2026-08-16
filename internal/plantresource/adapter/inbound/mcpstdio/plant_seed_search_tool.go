package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

type plantSeedSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"식물종자의 국명 또는 학명 검색어"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type plantSeedSearchOutput struct {
	Items      []plantSeedSearchItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                   `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                   `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                   `json:"totalCount" jsonschema:"전체 검색 결과 수"`
}

type plantSeedSearchItem struct {
	APGFamilyKorNm   string `json:"apgFamilyKorNm" jsonschema:"APG과국명"`
	APGFamilyNm      string `json:"apgFamilyNm" jsonschema:"APG과명"`
	BlprdEnmnt       string `json:"blprdEnmnt" jsonschema:"개화기종료일"`
	BlprdStmnt       string `json:"blprdStmnt" jsonschema:"개화기시작일"`
	ClrngMthodCdNm   string `json:"clrngMthodCdNm" jsonschema:"정선방법"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	FritCdNm         string `json:"fritCdNm" jsonschema:"열매형태"`
	FrssnEnmnt       string `json:"frssnEnmnt" jsonschema:"결실기종료일"`
	FrssnStmnt       string `json:"frssnStmnt" jsonschema:"결실기시작일"`
	LastUpdtDtm      string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	PlantGnrlNm      string `json:"plantGnrlNm" jsonschema:"국명(식물명)"`
	PlantSpecsScnm   string `json:"plantSpecsScnm" jsonschema:"학명"`
	RfrncLtrtrCont   string `json:"rfrncLtrtrCont" jsonschema:"참고문헌"`
	SeedCtsrfcDesc   string `json:"seedCtsrfcDesc" jsonschema:"종자표면형태설명"`
	SeedCtsrfcTpcdNm string `json:"seedCtsrfcTpcdNm" jsonschema:"종자표면형태"`
	SeedEmbrTpcdNm   string `json:"seedEmbrTpcdNm" jsonschema:"배아형태"`
	SeedMnmmBrdth    string `json:"seedMnmmBrdth" jsonschema:"종자최소너비"`
	SeedMnmmLngth    string `json:"seedMnmmLngth" jsonschema:"종자최소길이"`
	SeedMxmmBrdth    string `json:"seedMxmmBrdth" jsonschema:"종자최대너비"`
	SeedMxmmLngth    string `json:"seedMxmmLngth" jsonschema:"종자최대길이"`
	SeedShpDesc      string `json:"seedShpDesc" jsonschema:"종자형태설명"`
	SeedShpTpcdNm    string `json:"seedShpTpcdNm" jsonschema:"종자형태"`
	SeedSpecsID      string `json:"seedSpecsId" jsonschema:"종자종ID"`
	SeedTpcdNm       string `json:"seedTpcdNm" jsonschema:"종자구분"`
}

type plantSeedSearchHandler struct {
	useCase inbound.PlantSeedSearchUseCase
}

func addPlantSeedSearchTool(server *mcp.Server, useCase inbound.PlantSeedSearchUseCase) {
	handler := plantSeedSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "plant_resource_plant_seed_search",
		Description: "산림청 국립수목원 식물종자 기본정보 목록을 검색합니다.",
	}, handler.handle)
}

func (h plantSeedSearchHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input plantSeedSearchInput) (*mcp.CallToolResult, plantSeedSearchOutput, error) {
	result, err := h.useCase.PlantSeedSearch(ctx, application.PlantSeedSearchQuery{
		PageNo:       input.PageNo,
		NumOfRows:    input.NumOfRows,
		ReqSearchWrd: input.ReqSearchWrd,
	})
	if err != nil {
		return nil, plantSeedSearchOutput{}, err
	}

	items := make([]plantSeedSearchItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = plantSeedSearchItem{
			APGFamilyKorNm:   item.APGFamilyKorNm,
			APGFamilyNm:      item.APGFamilyNm,
			BlprdEnmnt:       item.BlprdEnmnt,
			BlprdStmnt:       item.BlprdStmnt,
			ClrngMthodCdNm:   item.ClrngMthodCdNm,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			FritCdNm:         item.FritCdNm,
			FrssnEnmnt:       item.FrssnEnmnt,
			FrssnStmnt:       item.FrssnStmnt,
			LastUpdtDtm:      item.LastUpdtDtm,
			PlantGnrlNm:      item.PlantGnrlNm,
			PlantSpecsScnm:   item.PlantSpecsScnm,
			RfrncLtrtrCont:   item.RfrncLtrtrCont,
			SeedCtsrfcDesc:   item.SeedCtsrfcDesc,
			SeedCtsrfcTpcdNm: item.SeedCtsrfcTpcdNm,
			SeedEmbrTpcdNm:   item.SeedEmbrTpcdNm,
			SeedMnmmBrdth:    item.SeedMnmmBrdth,
			SeedMnmmLngth:    item.SeedMnmmLngth,
			SeedMxmmBrdth:    item.SeedMxmmBrdth,
			SeedMxmmLngth:    item.SeedMxmmLngth,
			SeedShpDesc:      item.SeedShpDesc,
			SeedShpTpcdNm:    item.SeedShpTpcdNm,
			SeedSpecsID:      item.SeedSpecsID,
			SeedTpcdNm:       item.SeedTpcdNm,
		}
	}

	return nil, plantSeedSearchOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
