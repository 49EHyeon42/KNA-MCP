package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type entogSpcmSearchUseCaseStub struct {
	query  application.EntogSpcmSearchQuery
	result application.EntogSpcmSearchResult
	err    error
}

func (s *entogSpcmSearchUseCaseStub) EntogSpcmSearch(_ context.Context, query application.EntogSpcmSearchQuery) (application.EntogSpcmSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestEntogSpcmSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &entogSpcmSearchUseCaseStub{result: application.EntogSpcmSearchResult{
		Items: []application.EntogSpcmSearchItem{{
			Btnc: "btnc", ClctDyDesc: "clctDyDesc", CprtCtnt: "cprtCtnt", DetailYn: "detailYn",
			EntogOfnmKrlngNm: "entogOfnmKrlngNm", EntogSmplNo: "entogSmplNo",
			FamilyKorNm: "familyKorNm", FamilyNm: "familyNm", FrstRgstnDtm: " ",
			GenusKorNm: " ", GenusNm: "genusNm", ImgURL: "NONE", LastUpdtDtm: " ",
			OrdKorNm: "ordKorNm", OrdNm: "ordNm",
		}},
		NumOfRows: 10, PageNo: 2, TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addEntogSpcmSearchTool(server, useCase)
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
	if tool.Name != "entog_service_entog_spcm_search" || tool.Description != "산림청 국립수목원 내구강표본 목록을 검색합니다." {
		t.Errorf("tool = %#v", tool)
	}

	wantInputDescriptions := map[string]string{
		"st":        "검색어구분 (1: 국명 부분 검색, 2: 학명 부분 검색, 3: 국명 일치 검색, 4: 학명 일치 검색)",
		"sw":        "검색대상어",
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
		Name: "entog_service_entog_spcm_search",
		Arguments: map[string]any{
			"st": "2", "sw": "test-search-word", "numOfRows": 10, "pageNo": 2,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	wantQuery := application.EntogSpcmSearchQuery{St: "2", Sw: "test-search-word", NumOfRows: 10, PageNo: 2}
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
		"btnc": "btnc", "clctDyDesc": "clctDyDesc", "cprtCtnt": "cprtCtnt", "detailYn": "detailYn",
		"entogOfnmKrlngNm": "entogOfnmKrlngNm", "entogSmplNo": "entogSmplNo",
		"familyKorNm": "familyKorNm", "familyNm": "familyNm", "frstRgstnDtm": " ",
		"genusKorNm": " ", "genusNm": "genusNm", "imgUrl": "NONE", "lastUpdtDtm": " ",
		"ordKorNm": "ordKorNm", "ordNm": "ordNm",
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
		"items": "조회 결과 목록", "numOfRows": "한 페이지 결과 수", "pageNo": "페이지번호", "totalCount": "전체 결과 수",
	}
	checkSchemaProperties(t, tool.OutputSchema, wantOutputDescriptions)
	wantItemDescriptions := map[string]string{
		"btnc": "학명", "clctDyDesc": "채집일", "cprtCtnt": "저작권", "detailYn": "상세정보유무",
		"entogOfnmKrlngNm": "국명", "entogSmplNo": "표본번호", "familyKorNm": "과국명", "familyNm": "과명",
		"frstRgstnDtm": "최초등록일", "genusKorNm": "속국명", "genusNm": "속명", "imgUrl": "이미지URL",
		"lastUpdtDtm": "최종수정일", "ordKorNm": "목국명", "ordNm": "목명",
	}
	outputSchema := tool.OutputSchema.(map[string]any)
	outputProperties := outputSchema["properties"].(map[string]any)
	itemsSchema := outputProperties["items"].(map[string]any)["items"].(map[string]any)
	checkSchemaProperties(t, itemsSchema, wantItemDescriptions)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"st": "1", "sw": "test-search-word", "numOfRows": 1, "pageNo": 1}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}
