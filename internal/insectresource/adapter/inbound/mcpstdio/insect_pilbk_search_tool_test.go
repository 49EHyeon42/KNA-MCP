package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type insectPilbkSearchUseCaseStub struct {
	query  application.InsectPilbkSearchQuery
	result application.InsectPilbkSearchResult
	err    error
}

func (s *insectPilbkSearchUseCaseStub) InsectPilbkSearch(_ context.Context, query application.InsectPilbkSearchQuery) (application.InsectPilbkSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestInsectPilbkSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &insectPilbkSearchUseCaseStub{result: application.InsectPilbkSearchResult{
		Items: []application.InsectPilbkSearchItem{{
			FamilyKorNm:    "family Korean name",
			FamilyNm:       "family name",
			GenusKorNm:     "genus Korean name",
			GenusNm:        "genus name",
			InsctGnrlNm:    "insect general name",
			InsctPilbkNo:   "insect pictorial book number",
			InsctSpecsScnm: "insect species scientific name",
			LastUpdtDtm:    "last update date time",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addInsectPilbkSearchTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "insect_resource_insect_pilbk_search",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 곤충의 국명 또는 학명",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "insect_resource_insect_pilbk_search", "산림청 국립수목원 곤충도감 목록을 검색합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "insect_resource_insect_pilbk_search",
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

	wantQuery := application.InsectPilbkSearchQuery{
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
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(21) {
		t.Errorf("pagination = %#v", output)
	}
	checkKeys(t, output, "items", "numOfRows", "pageNo", "totalCount")
	items, ok := output["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", output["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item = %#v", items[0])
	}
	wantItem := map[string]string{
		"familyKorNm":    "family Korean name",
		"familyNm":       "family name",
		"genusKorNm":     "genus Korean name",
		"genusNm":        "genus name",
		"insctGnrlNm":    "insect general name",
		"insctPilbkNo":   "insect pictorial book number",
		"insctSpecsScnm": "insect species scientific name",
		"lastUpdtDtm":    "last update date time",
	}
	itemDescriptions := map[string]string{
		"familyKorNm":    "과국명",
		"familyNm":       "과명",
		"genusKorNm":     "속국명",
		"genusNm":        "속명",
		"insctGnrlNm":    "국명(곤충명)",
		"insctPilbkNo":   "곤충도감번호",
		"insctSpecsScnm": "학명",
		"lastUpdtDtm":    "최종수정일",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "insect_resource_insect_pilbk_search",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지당 건 수",
			"pageNo":     "페이지 번호",
			"totalCount": "전체 건 수",
		}, itemDescriptions)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "insect_resource_insect_pilbk_search",
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
		if wantItemDescriptions == nil {
			return
		}

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
