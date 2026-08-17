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

type fngsSmplSearchUseCaseStub struct {
	query  application.FngsSmplSearchQuery
	result application.FngsSmplSearchResult
	err    error
}

func (s *fngsSmplSearchUseCaseStub) FngsSmplSearch(_ context.Context, query application.FngsSmplSearchQuery) (application.FngsSmplSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestFngsSmplSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &fngsSmplSearchUseCaseStub{result: application.FngsSmplSearchResult{
		Items: []application.FngsSmplSearchItem{{
			Cnt:         "sample count",
			FamilyKorNm: "family Korean name",
			FamilyNm:    "family name",
			FngsGnrlNm:  "fungi general name",
			FngsID:      "fungi ID",
			FngsScnm:    "fungi scientific name",
			GenusKorNm:  "genus Korean name",
			GenusNm:     "genus name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addFngsSmplSearchTool(server, useCase)
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

	checkToolInputSchema(t, ctx, clientSession, "fungi_resource_fngs_smpl_search",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 버섯표본의 학명 또는 국명 (대소문자를 구분하지 않는 부분 문자열 검색)",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "fungi_resource_fngs_smpl_search", "산림청 국립수목원 버섯표본 목록을 검색합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "fungi_resource_fngs_smpl_search",
		Arguments: map[string]any{
			"pageNo":       2,
			"numOfRows":    10,
			"reqSearchWrd": "test-search-word",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.FngsSmplSearchQuery{
		PageNo:       2,
		NumOfRows:    10,
		ReqSearchWrd: "test-search-word",
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
		"cnt":         "sample count",
		"familyKorNm": "family Korean name",
		"familyNm":    "family name",
		"fngsGnrlNm":  "fungi general name",
		"fngsId":      "fungi ID",
		"fngsScnm":    "fungi scientific name",
		"genusKorNm":  "genus Korean name",
		"genusNm":     "genus name",
	}
	if !reflect.DeepEqual(item, wantItem) {
		t.Errorf("item = %#v, want %#v", item, wantItem)
	}
	checkToolOutputSchema(t, ctx, clientSession, "fungi_resource_fngs_smpl_search",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지 결과 수",
			"pageNo":     "페이지번호",
			"totalCount": "전체 결과 수",
		}, map[string]string{
			"cnt":         "표본 수",
			"familyKorNm": "과국명",
			"familyNm":    "과명",
			"fngsGnrlNm":  "국명(버섯명)",
			"fngsId":      "버섯 종ID",
			"fngsScnm":    "학명",
			"genusKorNm":  "속국명",
			"genusNm":     "속명",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "fungi_resource_fngs_smpl_search",
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

func TestFngsSmplSearchOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.FngsSmplSearchItem{})
	adapterFields := reflect.TypeOf(fngsSmplSearchItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for i := range applicationFields.NumField() {
		if applicationFields.Field(i).Name != adapterFields.Field(i).Name {
			t.Errorf("field %d = %s, want %s", i, adapterFields.Field(i).Name, applicationFields.Field(i).Name)
		}
	}
}
