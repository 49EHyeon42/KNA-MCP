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

type alchnSpcmSearchUseCaseStub struct {
	query  application.AlchnSpcmSearchQuery
	result application.AlchnSpcmSearchResult
	err    error
}

func (s *alchnSpcmSearchUseCaseStub) AlchnSpcmSearch(_ context.Context, query application.AlchnSpcmSearchQuery) (application.AlchnSpcmSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestAlchnSpcmSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &alchnSpcmSearchUseCaseStub{result: application.AlchnSpcmSearchResult{
		Items: []application.AlchnSpcmSearchItem{{
			Btnc: "btnc", CltrNm: " ", CprtCtnt: "cprtCtnt", DetailYn: "detailYn", EngNm: "engNm",
			FamilyKorNm: "familyKorNm", FamilyNm: "familyNm", FrstRgstnDtm: "frstRgstnDtm",
			GenusKorNm: "genusKorNm", GenusNm: "genusNm", ImgURL: "imgUrl\tvalue", JapNm: "japNm",
			LastUpdtDtm: "lastUpdtDtm", LchnGnrlNm: "lchnGnrlNm", LchnScnmID: "lchnScnmId",
			LchnSmplNo: "lchnSmplNo", PrkNm: "prkNm",
		}},
		NumOfRows: 10, PageNo: 2, TotalCount: 1068,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addAlchnSpcmSearchTool(server, useCase)
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
	if tool.Name != "lchn_service_alchn_spcm_search" || tool.Description != "산림청 국립수목원 지의류표본 목록을 검색합니다." {
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
	required := tool.InputSchema.(map[string]any)["required"].([]any)
	requiredNames := make([]string, len(required))
	for i, name := range required {
		requiredNames[i] = name.(string)
	}
	slices.Sort(requiredNames)
	if !slices.Equal(requiredNames, []string{"numOfRows", "pageNo", "st", "sw"}) {
		t.Errorf("required = %#v", requiredNames)
	}

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: tool.Name,
		Arguments: map[string]any{
			"st": "2", "sw": "Cladonia", "dateGbn": "1", "dateFrom": "20240101", "dateTo": "20241231", "numOfRows": 10, "pageNo": 2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	wantQuery := application.AlchnSpcmSearchQuery{St: "2", Sw: "Cladonia", DateGbn: "1", DateFrom: "20240101", DateTo: "20241231", NumOfRows: 10, PageNo: 2}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output := result.StructuredContent.(map[string]any)
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(1068) {
		t.Errorf("pagination = %#v", output)
	}
	item := output["items"].([]any)[0].(map[string]any)
	wantItem := map[string]string{
		"btnc": "btnc", "cltrNm": " ", "cprtCtnt": "cprtCtnt", "detailYn": "detailYn", "engNm": "engNm",
		"familyKorNm": "familyKorNm", "familyNm": "familyNm", "frstRgstnDtm": "frstRgstnDtm",
		"genusKorNm": "genusKorNm", "genusNm": "genusNm", "imgUrl": "imgUrl\tvalue", "japNm": "japNm",
		"lastUpdtDtm": "lastUpdtDtm", "lchnGnrlNm": "lchnGnrlNm", "lchnScnmId": "lchnScnmId",
		"lchnSmplNo": "lchnSmplNo", "prkNm": "prkNm",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if item[key] != want {
			t.Errorf("item %s = %#v, want %q", key, item[key], want)
		}
	}
	checkSchemaProperties(t, tool.OutputSchema, map[string]string{
		"items": "조회 결과 목록", "numOfRows": "페이지당레코드수", "pageNo": "페이지번호", "totalCount": "전체카운트",
	})
	outputProperties := tool.OutputSchema.(map[string]any)["properties"].(map[string]any)
	itemSchema := outputProperties["items"].(map[string]any)["items"].(map[string]any)
	checkSchemaProperties(t, itemSchema, map[string]string{
		"btnc": "학명", "cltrNm": "채집자명", "cprtCtnt": "저작권", "detailYn": "상세유무", "engNm": "영문명",
		"familyKorNm": "과국명", "familyNm": "과명", "frstRgstnDtm": "최초등록일시", "genusKorNm": "속국명",
		"genusNm": "속명", "imgUrl": "이미지URL", "japNm": "일어명", "lastUpdtDtm": "최종수정일시",
		"lchnGnrlNm": "국명", "lchnScnmId": "학명ID", "lchnSmplNo": "표본번호", "prkNm": "북한명",
	})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"st": "2", "sw": "Cladonia", "numOfRows": 1, "pageNo": 1}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}

func TestAlchnSpcmSearchOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.AlchnSpcmSearchItem{})
	adapterFields := reflect.TypeOf(alchnSpcmSearchItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
