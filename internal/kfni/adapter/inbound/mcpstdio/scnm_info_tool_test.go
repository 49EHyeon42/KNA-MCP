package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type scnmInfoUseCaseStub struct {
	query  application.ScnmInfoQuery
	result application.ScnmInfoResult
	err    error
}

func (s *scnmInfoUseCaseStub) ScnmInfo(_ context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestScnmInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &scnmInfoUseCaseStub{result: application.ScnmInfoResult{Item: &application.ScnmInfoItem{
		StpltScnmRltnCdNm: "scientific name relation code name",
		FalmNm:            "family name",
		FalnKorNm:         "family Korean name",
		FngsEclgTpcdNm:    "fungi ecology type code name",
		FngsGnrlNm:        "fungi general name",
		FngsGnrlNm2:       " ",
		FngsPrpseTpcdNm:   "fungi purpose type code name",
		FngsScnm:          "fungi scientific name",
		FngsScnmID:        "test-fungi-scientific-name-id",
		GenusKorNm:        "genus Korean name",
		GenusNm:           "genus name",
		LastUpdtDtm:       "last update date time",
		OrdscLtrtrNm:      "original description literature name",
		Rmrk:              " ",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addScnmInfoTool(server, useCase)
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

	checkToolInputSchema(t, ctx, clientSession, "kfni_scnm_info", map[string]string{
		"reqFngsScnmId": "버섯 학명ID (scnmSearch 결과의 fngsScnmId)",
	}, []string{"reqFngsScnmId"})
	checkToolDescription(t, ctx, clientSession, "kfni_scnm_info", "산림청 국립수목원 국가표준버섯목록의 버섯 학명 상세 정보를 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "kfni_scnm_info",
		Arguments: map[string]any{"reqFngsScnmId": "test-fungi-scientific-name-id"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if want := (application.ScnmInfoQuery{ReqFngsScnmID: "test-fungi-scientific-name-id"}); useCase.query != want {
		t.Errorf("query = %#v, want %#v", useCase.query, want)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	checkKeys(t, output, "item")
	item, ok := output["item"].(map[string]any)
	if !ok {
		t.Fatalf("item = %#v", output["item"])
	}
	wantItem := map[string]string{
		"stpltScnmRltnCdNm": "scientific name relation code name",
		"falmNm":            "family name",
		"falnKorNm":         "family Korean name",
		"fngsEclgTpcdNm":    "fungi ecology type code name",
		"fngsGnrlNm":        "fungi general name",
		"fngsGnrlNm2":       " ",
		"fngsPrpseTpcdNm":   "fungi purpose type code name",
		"fngsScnm":          "fungi scientific name",
		"fngsScnmId":        "test-fungi-scientific-name-id",
		"genusKorNm":        "genus Korean name",
		"genusNm":           "genus name",
		"lastUpdtDtm":       "last update date time",
		"ordscLtrtrNm":      "original description literature name",
		"rmrk":              " ",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	checkScnmInfoOutputSchema(t, ctx, clientSession)

	useCase.result = application.ScnmInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "kfni_scnm_info",
		Arguments: map[string]any{"reqFngsScnmId": "not-found"},
	})
	if err != nil {
		t.Fatal(err)
	}
	output, ok = result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	checkKeys(t, output, "item")
	if output["item"] != nil {
		t.Errorf("item = %#v, want nil", output["item"])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "kfni_scnm_info",
		Arguments: map[string]any{"reqFngsScnmId": "test-fungi-scientific-name-id"},
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

func checkScnmInfoOutputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	itemDescriptions := map[string]string{
		"stpltScnmRltnCdNm": "버섯 학명의 정명/이명 구분",
		"falmNm":            "버섯 학명의 과명(Family name)",
		"falnKorNm":         "버섯 학명 과명의 과국명",
		"fngsEclgTpcdNm":    "버섯 생태형",
		"fngsGnrlNm":        "버섯의 추천 국명(버섯명)",
		"fngsGnrlNm2":       "버섯의 비추천 국명(버섯명)",
		"fngsPrpseTpcdNm":   "버섯의 식독 정보",
		"fngsScnm":          "버섯 학명",
		"fngsScnmId":        "버섯 학명ID",
		"genusKorNm":        "버섯 학명 속명의 속국명",
		"genusNm":           "버섯 학명의 속명(Genus name)",
		"lastUpdtDtm":       "최종수정일",
		"ordscLtrtrNm":      "버섯 학명의 기재문 정보",
		"rmrk":              "비고",
	}

	for _, tool := range result.Tools {
		if tool.Name != "kfni_scnm_info" {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool kfni_scnm_info output schema = %#v", tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool kfni_scnm_info output schema properties = %#v", schema["properties"])
		}
		checkKeys(t, properties, "item")
		checkPropertyDescriptions(t, "kfni_scnm_info", properties, map[string]string{"item": "상세 조회 결과"})

		itemProperty, ok := properties["item"].(map[string]any)
		if !ok {
			t.Fatalf("tool kfni_scnm_info item property = %#v", properties["item"])
		}
		itemProperties, nullable := nullableObjectProperties(itemProperty)
		if !nullable {
			t.Errorf("tool kfni_scnm_info item schema is not nullable: %#v", itemProperty)
		}
		if itemProperties == nil {
			t.Fatalf("tool kfni_scnm_info item schema has no object properties: %#v", itemProperty)
		}
		checkKeys(t, itemProperties, mapKeys(itemDescriptions)...)
		checkPropertyDescriptions(t, "kfni_scnm_info item", itemProperties, itemDescriptions)
		return
	}
	t.Fatal("tool kfni_scnm_info not found")
}

func nullableObjectProperties(schema map[string]any) (map[string]any, bool) {
	properties, _ := schema["properties"].(map[string]any)
	types, _ := schema["type"].([]any)
	nullable := false
	for _, schemaType := range types {
		if schemaType == "null" {
			nullable = true
		}
	}
	for _, key := range []string{"anyOf", "oneOf"} {
		alternatives, _ := schema[key].([]any)
		for _, alternative := range alternatives {
			candidate, _ := alternative.(map[string]any)
			if candidate["type"] == "null" {
				nullable = true
			}
			if candidateProperties, ok := candidate["properties"].(map[string]any); ok {
				properties = candidateProperties
			}
		}
	}
	return properties, nullable
}

func TestScnmInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.ScnmInfoItem{})
	adapterFields := reflect.TypeOf(scnmInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
