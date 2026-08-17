package mcpstdio

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/inbound"
)

type fngsSmplUnitListInput struct {
	PageNo    int    `json:"pageNo" jsonschema:"페이지번호 (1 이상)"`
	NumOfRows int    `json:"numOfRows" jsonschema:"한 페이지 결과 수 (1 이상)"`
	ReqFngsID string `json:"reqFngsId" jsonschema:"검색할 버섯표본의 버섯 종ID (fngsSmplSearch 결과의 fngsId)"`
}

type fngsSmplUnitListOutput struct {
	Items      []fngsSmplUnitListItem `json:"items" jsonschema:"조회 결과 목록"`
	NumOfRows  int                    `json:"numOfRows" jsonschema:"한 페이지 결과 수"`
	PageNo     int                    `json:"pageNo" jsonschema:"페이지번호"`
	TotalCount int                    `json:"totalCount" jsonschema:"전체 결과 수"`
}

type fngsSmplUnitListItem struct {
	ClarDtlDscrt     string `json:"clarDtlDscrt" jsonschema:"채집지 상세"`
	ClarHaslvVal     string `json:"clarHaslvVal" jsonschema:"채집지 해발 고도"`
	CllcrNm          string `json:"cllcrNm" jsonschema:"채집자"`
	FamilyKorNm      string `json:"familyKorNm" jsonschema:"과국명"`
	FamilyNm         string `json:"familyNm" jsonschema:"과명"`
	FngsEclgTpcdNm   string `json:"fngsEclgTpcdNm" jsonschema:"버섯 생태형"`
	FngsGnrlNm       string `json:"fngsGnrlNm" jsonschema:"국명(버섯명)"`
	FngsID           string `json:"fngsId" jsonschema:"버섯종ID"`
	FngsScnm         string `json:"fngsScnm" jsonschema:"학명"`
	FngsSmplKindCdNm string `json:"fngsSmplKindCdNm" jsonschema:"버섯표본종류"`
	FngsSmplNo       string `json:"fngsSmplNo" jsonschema:"버섯표본번호"`
	GenusKorNm       string `json:"genusKorNm" jsonschema:"속국명"`
	GenusNm          string `json:"genusNm" jsonschema:"속명"`
	HbttChrcrCont    string `json:"hbttChrcrCont" jsonschema:"서식지 특성"`
	HstCont          string `json:"hstCont" jsonschema:"기주 정보"`
	LastUpdtDtm      string `json:"lastUpdtDtm" jsonschema:"최종수정일"`
	SmplCllcnDt      string `json:"smplCllcnDt" jsonschema:"표본 채집일"`
}

type fngsSmplUnitListHandler struct {
	useCase inbound.FngsSmplUnitListUseCase
}

func addFngsSmplUnitListTool(server *mcp.Server, useCase inbound.FngsSmplUnitListUseCase) {
	handler := fngsSmplUnitListHandler{useCase: useCase}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fungi_resource_fngs_smpl_unit_list",
		Description: "산림청 국립수목원 버섯표본 상세정보 목록을 조회합니다.",
	}, handler.handle)
}

func (h fngsSmplUnitListHandler) handle(ctx context.Context, _ *mcp.CallToolRequest, input fngsSmplUnitListInput) (*mcp.CallToolResult, fngsSmplUnitListOutput, error) {
	result, err := h.useCase.FngsSmplUnitList(ctx, application.FngsSmplUnitListQuery{
		PageNo:    input.PageNo,
		NumOfRows: input.NumOfRows,
		ReqFngsID: input.ReqFngsID,
	})
	if err != nil {
		return nil, fngsSmplUnitListOutput{}, err
	}

	items := make([]fngsSmplUnitListItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = fngsSmplUnitListItem{
			ClarDtlDscrt:     item.ClarDtlDscrt,
			ClarHaslvVal:     item.ClarHaslvVal,
			CllcrNm:          item.CllcrNm,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			FngsEclgTpcdNm:   item.FngsEclgTpcdNm,
			FngsGnrlNm:       item.FngsGnrlNm,
			FngsID:           item.FngsID,
			FngsScnm:         item.FngsScnm,
			FngsSmplKindCdNm: item.FngsSmplKindCdNm,
			FngsSmplNo:       item.FngsSmplNo,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			HbttChrcrCont:    item.HbttChrcrCont,
			HstCont:          item.HstCont,
			LastUpdtDtm:      item.LastUpdtDtm,
			SmplCllcnDt:      item.SmplCllcnDt,
		}
	}

	return nil, fngsSmplUnitListOutput{
		Items:      items,
		NumOfRows:  result.NumOfRows,
		PageNo:     result.PageNo,
		TotalCount: result.TotalCount,
	}, nil
}
