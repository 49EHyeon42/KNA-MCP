package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type fngsSmplUnitListUseCaseStub struct {
	query  application.FngsSmplUnitListQuery
	result application.FngsSmplUnitListResult
	err    error
}

func (s *fngsSmplUnitListUseCaseStub) FngsSmplUnitList(_ context.Context, query application.FngsSmplUnitListQuery) (application.FngsSmplUnitListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestFngsSmplUnitListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &fngsSmplUnitListUseCaseStub{result: application.FngsSmplUnitListResult{
		Items: []application.FngsSmplUnitListItem{{
			ClarDtlDscrt:     "collection site detail",
			ClarHaslvVal:     "collection site elevation",
			CllcrNm:          "collector name",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			FngsEclgTpcdNm:   "fungi ecology type code name",
			FngsGnrlNm:       "fungi general name",
			FngsID:           "fungi ID",
			FngsScnm:         "fungi scientific name",
			FngsSmplKindCdNm: "fungi sample kind code name",
			FngsSmplNo:       "fungi sample number",
			GenusKorNm:       "genus Korean name",
			GenusNm:          "genus name",
			HbttChrcrCont:    "habitat characteristic content",
			HstCont:          "host content",
			LastUpdtDtm:      "last update date time",
			SmplCllcnDt:      "sample collection date",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addFngsSmplUnitListTool(server, useCase)
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Wait()

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()

	checkToolInputSchema(t, ctx, clientSession, "fungi_resource_fngs_smpl_unit_list",
		map[string]string{
			"pageNo":    "페이지번호 (1 이상)",
			"numOfRows": "한 페이지 결과 수 (1 이상)",
			"reqFngsId": "검색할 버섯표본의 버섯 종ID (fngsSmplSearch 결과의 fngsId)",
		},
		[]string{"pageNo", "numOfRows", "reqFngsId"},
	)
	checkToolDescription(t, ctx, clientSession, "fungi_resource_fngs_smpl_unit_list", "산림청 국립수목원 버섯표본 상세정보 목록을 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "fungi_resource_fngs_smpl_unit_list",
		Arguments: map[string]any{
			"pageNo":    2,
			"numOfRows": 10,
			"reqFngsId": "test-fungi-id",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.FngsSmplUnitListQuery{
		PageNo:    2,
		NumOfRows: 10,
		ReqFngsID: "test-fungi-id",
	}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	checkKeys(t, output, "items", "numOfRows", "pageNo", "totalCount")
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(21) {
		t.Errorf("pagination = %#v", output)
	}
	items, ok := output["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", output["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item = %#v", items[0])
	}
	wantItem := map[string]any{
		"clarDtlDscrt":     "collection site detail",
		"clarHaslvVal":     "collection site elevation",
		"cllcrNm":          "collector name",
		"familyKorNm":      "family Korean name",
		"familyNm":         "family name",
		"fngsEclgTpcdNm":   "fungi ecology type code name",
		"fngsGnrlNm":       "fungi general name",
		"fngsId":           "fungi ID",
		"fngsScnm":         "fungi scientific name",
		"fngsSmplKindCdNm": "fungi sample kind code name",
		"fngsSmplNo":       "fungi sample number",
		"genusKorNm":       "genus Korean name",
		"genusNm":          "genus name",
		"hbttChrcrCont":    "habitat characteristic content",
		"hstCont":          "host content",
		"lastUpdtDtm":      "last update date time",
		"smplCllcnDt":      "sample collection date",
	}
	if !reflect.DeepEqual(item, wantItem) {
		t.Errorf("item = %#v, want %#v", item, wantItem)
	}
	checkToolOutputSchema(t, ctx, clientSession, "fungi_resource_fngs_smpl_unit_list",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지 결과 수",
			"pageNo":     "페이지번호",
			"totalCount": "전체 결과 수",
		}, map[string]string{
			"clarDtlDscrt":     "채집지 상세",
			"clarHaslvVal":     "채집지 해발 고도",
			"cllcrNm":          "채집자",
			"familyKorNm":      "과국명",
			"familyNm":         "과명",
			"fngsEclgTpcdNm":   "버섯 생태형",
			"fngsGnrlNm":       "국명(버섯명)",
			"fngsId":           "버섯종ID",
			"fngsScnm":         "학명",
			"fngsSmplKindCdNm": "버섯표본종류",
			"fngsSmplNo":       "버섯표본번호",
			"genusKorNm":       "속국명",
			"genusNm":          "속명",
			"hbttChrcrCont":    "서식지 특성",
			"hstCont":          "기주 정보",
			"lastUpdtDtm":      "최종수정일",
			"smplCllcnDt":      "표본 채집일",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "fungi_resource_fngs_smpl_unit_list",
		Arguments: map[string]any{
			"pageNo":    1,
			"numOfRows": 1,
			"reqFngsId": "test-fungi-id",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
	content, ok := result.Content[0].(*mcp.TextContent)
	if !ok || content.Text != useCase.err.Error() {
		t.Errorf("error content = %#v, want %q", result.Content, useCase.err)
	}
}

func TestFngsSmplUnitListOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.FngsSmplUnitListItem{})
	adapterFields := reflect.TypeOf(fngsSmplUnitListItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for i := range applicationFields.NumField() {
		if applicationFields.Field(i).Name != adapterFields.Field(i).Name {
			t.Errorf("field %d = %s, want %s", i, adapterFields.Field(i).Name, applicationFields.Field(i).Name)
		}
	}
}
