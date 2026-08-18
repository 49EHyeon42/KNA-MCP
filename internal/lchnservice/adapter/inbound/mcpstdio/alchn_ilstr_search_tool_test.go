package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type alchnIlstrSearchUseCaseStub struct {
	query  application.AlchnIlstrSearchQuery
	result application.AlchnIlstrSearchResult
	err    error
}

func (s *alchnIlstrSearchUseCaseStub) AlchnIlstrSearch(_ context.Context, query application.AlchnIlstrSearchQuery) (application.AlchnIlstrSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestAlchnIlstrSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &alchnIlstrSearchUseCaseStub{result: application.AlchnIlstrSearchResult{
		Items: []application.AlchnIlstrSearchItem{{
			Btnc: "btnc", CprtCtnt: "cprtCtnt", DetailYn: "detailYn", EngNm: " ",
			FamilyKorNm: "familyKorNm", FamilyNm: "familyNm", FrstRgstnDtm: "frstRgstnDtm",
			GenusKorNm: "genusKorNm", GenusNm: "genusNm", ImgURL: "imgUrl", JapNm: "japNm",
			LastUpdtDtm: "lastUpdtDtm", LchnGnrlNm: "lchnGnrlNm", LchnInfrpNm: "lchnInfrpNm",
			LchnPilbkNo: "lchnPilbkNo", LchnScnmID: "lchnScnmId", LchnTtnm: "lchnTtnm", PrkNm: "prkNm",
		}},
		NumOfRows: 10, PageNo: 2, TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addAlchnIlstrSearchTool(server, useCase)
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

	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 1 {
		t.Fatalf("tools = %#v, want one tool", tools.Tools)
	}
	tool := tools.Tools[0]
	if tool.Name != "lchn_service_alchn_ilstr_search" || tool.Description != "산림청 국립수목원 지의류도감 목록을 검색합니다." {
		t.Errorf("tool = %#v", tool)
	}

	wantInputDescriptions := map[string]string{
		"st":        "검색어구분 (1: 국명 부분 검색, 2: 학명 부분 검색, 3: 국명 일치 검색, 4: 학명 일치 검색)",
		"sw":        "검색대상어",
		"dateGbn":   "날짜검색 구분 (1: 등록일, 2: 수정일)",
		"dateFrom":  "검색 시작일 (dateGbn 입력 시 필수, yyyyMMdd)",
		"dateTo":    "검색 종료일 (dateGbn 입력 시 필수, yyyyMMdd)",
		"numOfRows": "한 페이지 결과 수 (1 이상)",
		"pageNo":    "페이지 번호 (1 이상)",
	}
	checkSchemaProperties(t, tool.InputSchema, wantInputDescriptions)
	inputSchema := tool.InputSchema.(map[string]any)
	required := inputSchema["required"].([]any)
	requiredNames := make([]string, len(required))
	for i, name := range required {
		requiredNames[i] = name.(string)
	}
	slices.Sort(requiredNames)
	if !slices.Equal(requiredNames, []string{"numOfRows", "pageNo", "st", "sw"}) {
		t.Errorf("required = %#v", requiredNames)
	}

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "lchn_service_alchn_ilstr_search",
		Arguments: map[string]any{
			"st": "2", "sw": "test-search-word", "dateGbn": "1", "dateFrom": "20240101", "dateTo": "20241231", "numOfRows": 10, "pageNo": 2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	wantQuery := application.AlchnIlstrSearchQuery{St: "2", Sw: "test-search-word", DateGbn: "1", DateFrom: "20240101", DateTo: "20241231", NumOfRows: 10, PageNo: 2}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output := result.StructuredContent.(map[string]any)
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(21) {
		t.Errorf("pagination = %#v", output)
	}
	items := output["items"].([]any)
	item := items[0].(map[string]any)
	wantItem := map[string]string{
		"btnc": "btnc", "cprtCtnt": "cprtCtnt", "detailYn": "detailYn", "engNm": " ",
		"familyKorNm": "familyKorNm", "familyNm": "familyNm", "frstRgstnDtm": "frstRgstnDtm",
		"genusKorNm": "genusKorNm", "genusNm": "genusNm", "imgUrl": "imgUrl", "japNm": "japNm",
		"lastUpdtDtm": "lastUpdtDtm", "lchnGnrlNm": "lchnGnrlNm", "lchnInfrpNm": "lchnInfrpNm",
		"lchnPilbkNo": "lchnPilbkNo", "lchnScnmId": "lchnScnmId", "lchnTtnm": "lchnTtnm", "prkNm": "prkNm",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if item[key] != want {
			t.Errorf("item %s = %#v, want %q", key, item[key], want)
		}
	}

	wantOutputDescriptions := map[string]string{
		"items": "조회 결과 목록", "numOfRows": "페이지당레코드수", "pageNo": "페이지번호", "totalCount": "전체카운트",
	}
	checkSchemaProperties(t, tool.OutputSchema, wantOutputDescriptions)
	wantItemDescriptions := map[string]string{
		"btnc": "학명", "cprtCtnt": "저작권", "detailYn": "상세정보유무", "engNm": "영문명",
		"familyKorNm": "과국명", "familyNm": "과명", "frstRgstnDtm": "최초등록일시", "genusKorNm": "속국명",
		"genusNm": "속명", "imgUrl": "이미지URL", "japNm": "일본명", "lastUpdtDtm": "최종수정일시",
		"lchnGnrlNm": "국명", "lchnInfrpNm": "종하명", "lchnPilbkNo": "도감번호", "lchnScnmId": "학명ID",
		"lchnTtnm": "종소명", "prkNm": "북한명",
	}
	outputSchema := tool.OutputSchema.(map[string]any)
	outputProperties := outputSchema["properties"].(map[string]any)
	itemsSchema := outputProperties["items"].(map[string]any)["items"].(map[string]any)
	checkSchemaProperties(t, itemsSchema, wantItemDescriptions)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"st": "2", "sw": "test-search-word", "numOfRows": 1, "pageNo": 1}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}

func checkSchemaProperties(t *testing.T, schema any, wantDescriptions map[string]string) {
	t.Helper()
	properties := schema.(map[string]any)["properties"].(map[string]any)
	if len(properties) != len(wantDescriptions) {
		t.Errorf("property count = %d, want %d", len(properties), len(wantDescriptions))
	}
	for name, want := range wantDescriptions {
		property, ok := properties[name].(map[string]any)
		if !ok {
			t.Errorf("missing property %q", name)
			continue
		}
		if property["description"] != want {
			t.Errorf("property %s description = %#v, want %q", name, property["description"], want)
		}
	}
}
