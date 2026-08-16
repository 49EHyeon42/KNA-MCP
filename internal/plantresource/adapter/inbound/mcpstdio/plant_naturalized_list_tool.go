package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/inbound"
)

const plantResourcePlantNaturalizedListToolName = "plant_resource_plant_naturalized_list"

type plantNaturalizedListInput struct {
	PageNo       int    `json:"pageNo" jsonschema:"페이지 번호(1 이상)"`
	NumOfRows    int    `json:"numOfRows" jsonschema:"페이지당 결과 수(1 이상)"`
	ReqSearchWrd string `json:"reqSearchWrd,omitempty" jsonschema:"외래식물의 국명 또는 학명 검색어"`
	// dateFrom and dateTo are disabled because the upstream API returns ORA-00908.
}

type plantNaturalizedListOutput struct {
	Items      []plantNaturalizedListItem `json:"items"`
	NumOfRows  int                        `json:"numOfRows"`
	PageNo     int                        `json:"pageNo"`
	TotalCount int                        `json:"totalCount"`
}

type plantNaturalizedListItem struct {
	AgpFamilyNm        string `json:"agpFamilyNm"`
	APGFamilyKorNm     string `json:"apgFamilyKorNm"`
	BlprdEnmnt         string `json:"blprdEnmnt"`
	BlprdStmnt         string `json:"blprdStmnt"`
	DistrAraDscrt      string `json:"distrAraDscrt"`
	EclgDstrbYn        string `json:"eclgDstrbYn"`
	ExtcPlantCdNm      string `json:"extcPlantCdNm"`
	FamilyKorNm        string `json:"familyKorNm"`
	FamilyNm           string `json:"familyNm"`
	FrtTpcdNm          string `json:"frtTpcdNm"`
	LastUpdtDtm        string `json:"lastUpdtDtm"`
	NtldgTpcdNm        string `json:"ntldgTpcdNm"`
	NtrlzEraTpcdNm     string `json:"ntrlzEraTpcdNm"`
	OrplcNm            string `json:"orplcNm"`
	PlantBrdgFomTpcdNm string `json:"plantBrdgFomTpcdNm"`
	PlantDistrGrcd     string `json:"plantDistrGrcd"`
	PlantDistrQntt     string `json:"plantDistrQntt"`
	PlantDistrQnttGrcd string `json:"plantDistrQnttGrcd"`
	PlantEngNm         string `json:"plantEngNm"`
	PlantGnrlNm        string `json:"plantGnrlNm"`
	PlantJpnNm         string `json:"plantJpnNm"`
	PlantLfcclTpcdNm   string `json:"plantLfcclTpcdNm"`
	PlantSpecsScnm     string `json:"plantSpecsScnm"`
}

type plantNaturalizedListHandler struct {
	useCase inbound.PlantNaturalizedListUseCase
}

func addPlantNaturalizedListTool(server *mcp.Server, useCase inbound.PlantNaturalizedListUseCase) {
	handler := plantNaturalizedListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        plantResourcePlantNaturalizedListToolName,
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
