package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type childPilbkSearchUseCaseStub struct {
	query  application.ChildPilbkSearchQuery
	result application.ChildPilbkSearchResult
	err    error
}

func (s *childPilbkSearchUseCaseStub) ChildPilbkSearch(_ context.Context, query application.ChildPilbkSearchQuery) (application.ChildPilbkSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestChildPilbkSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &childPilbkSearchUseCaseStub{result: application.ChildPilbkSearchResult{
		Items: []application.ChildPilbkSearchItem{{
			BiogyNm:           "biology name",
			ChildLvbngPilbkNo: "child pictorial book number",
			FamilyKorNm:       "family Korean name",
			FamilyNm:          "family name",
			LvbngTpcdNm:       "living thing type code name",
			LvngKrlngNm:       "living thing Korean name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addChildPilbkSearchTool(server, useCase)
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

	checkToolInputSchema(t, ctx, clientSession, "child_service_child_pilbk_search",
		map[string]string{
			"pageNo":       "페이지 번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 생물의 국명 또는 학명",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "child_service_child_pilbk_search", "산림청 국립수목원 어린이생물도감 목록을 검색합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "child_service_child_pilbk_search",
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

	wantQuery := application.ChildPilbkSearchQuery{
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
	wantItem := map[string]string{
		"biogyNm":           "biology name",
		"childLvbngPilbkNo": "child pictorial book number",
		"familyKorNm":       "family Korean name",
		"familyNm":          "family name",
		"lvbngTpcdNm":       "living thing type code name",
		"lvngKrlngNm":       "living thing Korean name",
	}
	itemDescriptions := map[string]string{
		"biogyNm":           "생물학명",
		"childLvbngPilbkNo": "어린이생물도감번호",
		"familyKorNm":       "과국명",
		"familyNm":          "과명",
		"lvbngTpcdNm":       "생물분류",
		"lvngKrlngNm":       "생물국명",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "child_service_child_pilbk_search",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지 결과 수",
			"pageNo":     "페이지번호",
			"totalCount": "전체 결과 수",
		}, itemDescriptions)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "child_service_child_pilbk_search",
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

func checkKeys(t *testing.T, got map[string]any, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("keys = %#v, want %#v", reflect.ValueOf(got).MapKeys(), want)
	}
	for _, key := range want {
		if _, ok := got[key]; !ok {
			t.Errorf("missing key %q in %#v", key, got)
		}
	}
}

func mapKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func checkPropertyDescriptions(t *testing.T, toolName string, properties map[string]any, wantDescriptions map[string]string) {
	t.Helper()
	for name, property := range properties {
		propertySchema, ok := property.(map[string]any)
		if !ok {
			t.Fatalf("tool %s property %s schema = %#v", toolName, name, property)
		}
		description, ok := propertySchema["description"].(string)
		if !ok || description == "" {
			t.Errorf("tool %s property %s has no description", toolName, name)
		}
		if want, ok := wantDescriptions[name]; ok && description != want {
			t.Errorf("tool %s property %s description = %q, want %q", toolName, name, description, want)
		}
	}
}

func checkToolDescription(t *testing.T, ctx context.Context, session *mcp.ClientSession, toolName, want string) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range result.Tools {
		if tool.Name == toolName {
			if tool.Description != want {
				t.Errorf("tool %s description = %q, want %q", toolName, tool.Description, want)
			}
			return
		}
	}
	t.Fatalf("tool %s not found", toolName)
}

func checkToolInputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession, toolName string, wantDescriptions map[string]string, wantRequired []string) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range result.Tools {
		if tool.Name != toolName {
			continue
		}
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool %s input schema = %#v", toolName, tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool %s input schema properties = %#v", toolName, schema["properties"])
		}
		checkKeys(t, properties, mapKeys(wantDescriptions)...)
		checkPropertyDescriptions(t, toolName, properties, wantDescriptions)

		required, ok := schema["required"].([]any)
		if !ok {
			t.Fatalf("tool %s input schema required = %#v", toolName, schema["required"])
		}
		requiredKeys := make(map[string]any, len(required))
		for _, key := range required {
			name, ok := key.(string)
			if !ok {
				t.Fatalf("tool %s input schema required key = %#v", toolName, key)
			}
			requiredKeys[name] = nil
		}
		checkKeys(t, requiredKeys, wantRequired...)
		return
	}
	t.Fatalf("tool %s not found", toolName)
}

func checkToolOutputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession, toolName string, wantDescriptions, wantItemDescriptions map[string]string) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range result.Tools {
		if tool.Name != toolName {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool %s output schema = %#v", toolName, tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool %s output schema properties = %#v", toolName, schema["properties"])
		}
		checkKeys(t, properties, mapKeys(wantDescriptions)...)
		checkPropertyDescriptions(t, toolName, properties, wantDescriptions)

		itemsProperty, ok := properties["items"].(map[string]any)
		if !ok {
			t.Fatalf("tool %s items property = %#v", toolName, properties["items"])
		}
		itemSchema, ok := itemsProperty["items"].(map[string]any)
		if !ok {
			t.Fatalf("tool %s item schema = %#v", toolName, itemsProperty["items"])
		}
		itemProperties, ok := itemSchema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool %s item schema properties = %#v", toolName, itemSchema["properties"])
		}
		checkKeys(t, itemProperties, mapKeys(wantItemDescriptions)...)
		checkPropertyDescriptions(t, toolName, itemProperties, wantItemDescriptions)
		return
	}
	t.Fatalf("tool %s not found", toolName)
}
