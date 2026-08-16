package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantSeedSearchToolName = "plant_resource_plant_seed_search"

type plantSeedSearchInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"식물종자의 국명 또는 학명 검색어"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type plantSeedSearchOutput struct {
	Items      []plantSeedSearchItem `json:"items"`
	NumOfRows  int                   `json:"numOfRows"`
	PageNo     int                   `json:"pageNo"`
	TotalCount int                   `json:"totalCount"`
}

type plantSeedSearchItem struct {
	APGFamilyKorNm   string `json:"apgFamilyKorNm"`
	APGFamilyNm      string `json:"apgFamilyNm"`
	BlprdEnmnt       string `json:"blprdEnmnt"`
	BlprdStmnt       string `json:"blprdStmnt"`
	ClrngMthodCdNm   string `json:"clrngMthodCdNm"`
	FamilyKorNm      string `json:"familyKorNm"`
	FamilyNm         string `json:"familyNm"`
	FritCdNm         string `json:"fritCdNm"`
	FrssnEnmnt       string `json:"frssnEnmnt"`
	FrssnStmnt       string `json:"frssnStmnt"`
	LastUpdtDtm      string `json:"lastUpdtDtm"`
	PlantGnrlNm      string `json:"plantGnrlNm"`
	PlantSpecsScnm   string `json:"plantSpecsScnm"`
	RfrncLtrtrCont   string `json:"rfrncLtrtrCont"`
	SeedCtsrfcDesc   string `json:"seedCtsrfcDesc"`
	SeedCtsrfcTpcdNm string `json:"seedCtsrfcTpcdNm"`
	SeedEmbrTpcdNm   string `json:"seedEmbrTpcdNm"`
	SeedMnmmBrdth    string `json:"seedMnmmBrdth"`
	SeedMnmmLngth    string `json:"seedMnmmLngth"`
	SeedMxmmBrdth    string `json:"seedMxmmBrdth"`
	SeedMxmmLngth    string `json:"seedMxmmLngth"`
	SeedShpDesc      string `json:"seedShpDesc"`
	SeedShpTpcdNm    string `json:"seedShpTpcdNm"`
	SeedSpecsID      string `json:"seedSpecsId"`
	SeedTpcdNm       string `json:"seedTpcdNm"`
}

type plantSeedSearchHandler struct {
	useCase inbound.PlantSeedSearchUseCase
}

func addPlantSeedSearchTool(server *mcp.Server, useCase inbound.PlantSeedSearchUseCase) {
	handler := plantSeedSearchHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantSeedSearchToolName,
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
