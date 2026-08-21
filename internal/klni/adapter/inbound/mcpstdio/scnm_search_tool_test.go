package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type scnmSearchUseCaseStub struct {
	query  application.ScnmSearchQuery
	result application.ScnmSearchResult
	err    error
}

func (s *scnmSearchUseCaseStub) ScnmSearch(_ context.Context, query application.ScnmSearchQuery) (application.ScnmSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestScnmSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &scnmSearchUseCaseStub{result: application.ScnmSearchResult{
		Items: []application.ScnmSearchItem{{
			StpltScnmRltnCdNm: "scientific name relation code name",
			ClassKorNm:        "class Korean name",
			ClassNm:           "class name",
			FalmNm:            "family name",
			FalnKorNm:         "family Korean name",
			GenusKorNm:        "genus Korean name",
			GenusNm:           "genus name",
			LastUpdtDtm:       "last update date time",
			LchnGnrlNm:        "lichen general name",
			LchnScnm:          "lichen scientific name",
			LchnScnmID:        "lichen scientific name ID",
			OrdKorNm:          "order Korean name",
			OrdNm:             "order name",
			PhylumKorNm:       "phylum Korean name",
			PhylumNm:          "phylum name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addScnmSearchTool(server, useCase)
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
		"pageNo":    "페이지번호 (1 이상)",
		"numOfRows": "한 페이지 결과 수 (1 이상)",
		"reqGnrlNm": "검색하려는 지의류 국명 (부분 문자열 검색)",
		"reqScnm":   "검색하려는 지의류 학명 (대소문자를 구분하지 않는 부분 문자열 검색)",
		"dateFrom":  "최종수정일 이후 정보 (yyyyMMdd)",
		"dateTo":    "최종수정일 이전 정보 (yyyyMMdd)",
	}
	checkToolInputSchema(t, ctx, clientSession, "klni_scnm_search", inputDescriptions, []string{"pageNo", "numOfRows"})
	checkToolDescription(t, ctx, clientSession, "klni_scnm_search", "산림청 국립수목원 국가표준지의류목록의 지의류 학명 목록을 검색합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "klni_scnm_search",
		Arguments: map[string]any{
			"pageNo":    2,
			"numOfRows": 10,
			"reqGnrlNm": "general name",
			"reqScnm":   "scientific name",
			"dateFrom":  "20240101",
			"dateTo":    "20241231",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.ScnmSearchQuery{
		PageNo:    2,
		NumOfRows: 10,
		ReqGnrlNm: "general name",
		ReqScnm:   "scientific name",
		DateFrom:  "20240101",
		DateTo:    "20241231",
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
		"stpltScnmRltnCdNm": "scientific name relation code name",
		"classKorNm":        "class Korean name",
		"classNm":           "class name",
		"falmNm":            "family name",
		"falnKorNm":         "family Korean name",
		"genusKorNm":        "genus Korean name",
		"genusNm":           "genus name",
		"lastUpdtDtm":       "last update date time",
		"lchnGnrlNm":        "lichen general name",
		"lchnScnm":          "lichen scientific name",
		"lchnScnmId":        "lichen scientific name ID",
		"ordKorNm":          "order Korean name",
		"ordNm":             "order name",
		"phylumKorNm":       "phylum Korean name",
		"phylumNm":          "phylum name",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	outputDescriptions := map[string]string{
		"items":      "조회 결과 목록",
		"numOfRows":  "한 페이지 결과 수",
		"pageNo":     "페이지 번호",
		"totalCount": "전체 결과 수",
	}
	itemDescriptions := map[string]string{
		"stpltScnmRltnCdNm": "지의류 학명의 정명/이명 구분",
		"classKorNm":        "지의류 학명 분류군의 강국명",
		"classNm":           "지의류 학명 분류군의 강명(Class Name)",
		"falmNm":            "지의류 학명 분류군의 과명(Family Name)",
		"falnKorNm":         "지의류 학명 분류군의 과국명",
		"genusKorNm":        "지의류 학명 분류군의 속국명",
		"genusNm":           "지의류 학명 분류군의 속명(Genus Name)",
		"lastUpdtDtm":       "최종수정일",
		"lchnGnrlNm":        "지의류 국명(지의류명)",
		"lchnScnm":          "지의류 학명",
		"lchnScnmId":        "지의류 학명ID",
		"ordKorNm":          "지의류 학명 분류군의 목국명",
		"ordNm":             "지의류 학명 분류군의 목명(Order Name)",
		"phylumKorNm":       "지의류 학명 분류군의 문국명",
		"phylumNm":          "지의류 학명 분류군의 문명(Phylum Name)",
	}
	checkToolOutputSchema(t, ctx, clientSession, "klni_scnm_search", outputDescriptions, itemDescriptions)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "klni_scnm_search",
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
