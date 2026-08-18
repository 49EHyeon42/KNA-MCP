package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type gnrlNmLtrtrSearchUseCaseStub struct {
	query  application.GnrlNmLtrtrSearchQuery
	result application.GnrlNmLtrtrSearchResult
	err    error
}

func (s *gnrlNmLtrtrSearchUseCaseStub) GnrlNmLtrtrSearch(_ context.Context, query application.GnrlNmLtrtrSearchQuery) (application.GnrlNmLtrtrSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestGnrlNmLtrtrSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &gnrlNmLtrtrSearchUseCaseStub{result: application.GnrlNmLtrtrSearchResult{
		Items: []application.GnrlNmLtrtrSearchItem{{
			RcmmnTpcdNm:      "recommendation type code name",
			LtrtrInfrmNm:     "literature information name",
			LvbngFrlngTpcdNm: "living foreign language type code name",
			PlantGnrlNm:      "plant general name",
			PlantSpecsScnm:   "plant species scientific name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addGnrlNmLtrtrSearchTool(server, useCase)
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

	inputDescriptions := map[string]string{
		"pageNo":         "페이지 번호 (1 이상)",
		"numOfRows":      "한 페이지 결과 수 (1 이상)",
		"reqPlantGnrlNm": "검색하려는 식물 국명(식물명) (부분 문자열 검색)",
	}
	checkToolInputSchema(t, ctx, clientSession, "kpni_gnrl_nm_ltrtr_search", inputDescriptions, []string{"pageNo", "numOfRows"})
	checkToolDescription(t, ctx, clientSession, "kpni_gnrl_nm_ltrtr_search", "산림청 국립수목원 국가표준식물목록의 식물 국명 출전 정보 목록을 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "kpni_gnrl_nm_ltrtr_search",
		Arguments: map[string]any{
			"pageNo":         2,
			"numOfRows":      10,
			"reqPlantGnrlNm": "test-search-word",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.GnrlNmLtrtrSearchQuery{PageNo: 2, NumOfRows: 10, ReqPlantGnrlNm: "test-search-word"}
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
	wantItem := map[string]string{
		"rcmmnTpcdNm":      "recommendation type code name",
		"ltrtrInfrmNm":     "literature information name",
		"lvbngFrlngTpcdNm": "living foreign language type code name",
		"plantGnrlNm":      "plant general name",
		"plantSpecsScnm":   "plant species scientific name",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	checkToolOutputSchema(t, ctx, clientSession, "kpni_gnrl_nm_ltrtr_search", map[string]string{
		"items":      "조회 결과 목록",
		"numOfRows":  "한 페이지 결과 수",
		"pageNo":     "페이지 번호",
		"totalCount": "전체 결과 수",
	}, map[string]string{
		"rcmmnTpcdNm":      "식물 국명 추천/비추천 구분",
		"ltrtrInfrmNm":     "식물 국명 출전 기재문",
		"lvbngFrlngTpcdNm": "국명 언어 분류",
		"plantGnrlNm":      "식물 국명(식물명)",
		"plantSpecsScnm":   "식물 학명",
	})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "kpni_gnrl_nm_ltrtr_search",
		Arguments: map[string]any{"pageNo": 1, "numOfRows": 1},
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

func TestGnrlNmLtrtrSearchOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.GnrlNmLtrtrSearchItem{})
	adapterFields := reflect.TypeOf(gnrlNmLtrtrSearchItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
