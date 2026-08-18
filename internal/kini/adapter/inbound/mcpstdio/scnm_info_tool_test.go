package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
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
		SuperFalmNm:       "superfamily name",
		ClassKorNm:        "class Korean name",
		ClassNm:           "class name",
		FalmKorNm:         "family Korean name",
		FalmNm:            "family name",
		GenusKorNm:        "genus Korean name",
		GenusNm:           "genus name",
		InsctGnrlNm:       "insect general name",
		InsctGnrlNm2:      " ",
		InsctScnmID:       "test-insect-scientific-name-id",
		InsctSpecsScnm:    "<em>insect</em> scientific name",
		LastUpdtDtm:       "last update date time",
		OrdKorNm:          "order Korean name",
		OrdNm:             "order name",
		StpltScnmRltnCdNm: "standard insect scientific name relation code name",
		SubFalmKorNm:      "subfamily Korean name",
		SubFalmNm:         "subfamily name",
		SuperFalmKorNm:    "superfamily Korean name",
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

	checkToolInputSchema(t, ctx, clientSession, "kini_scnm_info", map[string]string{
		"reqInsctScnmId": "조회하려는 곤충 학명ID (scnmSearch 결과의 insctScnmId)",
	}, []string{"reqInsctScnmId"})
	checkToolDescription(t, ctx, clientSession, "kini_scnm_info", "산림청 국립수목원 국가표준곤충목록의 곤충 학명 상세 정보를 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "kini_scnm_info",
		Arguments: map[string]any{"reqInsctScnmId": "test-insect-scientific-name-id"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if want := (application.ScnmInfoQuery{ReqInsctScnmID: "test-insect-scientific-name-id"}); useCase.query != want {
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
		"superFalmNm":       "superfamily name",
		"classKorNm":        "class Korean name",
		"classNm":           "class name",
		"falmKorNm":         "family Korean name",
		"falmNm":            "family name",
		"genusKorNm":        "genus Korean name",
		"genusNm":           "genus name",
		"insctGnrlNm":       "insect general name",
		"insctGnrlNm2":      " ",
		"insctScnmId":       "test-insect-scientific-name-id",
		"insctSpecsScnm":    "<em>insect</em> scientific name",
		"lastUpdtDtm":       "last update date time",
		"ordKorNm":          "order Korean name",
		"ordNm":             "order name",
		"stpltScnmRltnCdNm": "standard insect scientific name relation code name",
		"subFalmKorNm":      "subfamily Korean name",
		"subFalmNm":         "subfamily name",
		"superFalmKorNm":    "superfamily Korean name",
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
		Name:      "kini_scnm_info",
		Arguments: map[string]any{"reqInsctScnmId": "not-found"},
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
		Name:      "kini_scnm_info",
		Arguments: map[string]any{"reqInsctScnmId": "test-insect-scientific-name-id"},
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
		"superFalmNm":       "곤충 상과명(SuperFamily Name)",
		"classKorNm":        "곤충 강국명",
		"classNm":           "곤충 강명(Class Name)",
		"falmKorNm":         "곤충 과국명",
		"falmNm":            "곤충 과명(Family Name)",
		"genusKorNm":        "곤충 속국명",
		"genusNm":           "곤충 속명(Genus Name)",
		"insctGnrlNm":       "곤충 추천국명(곤충명)",
		"insctGnrlNm2":      "곤충 비추천국명",
		"insctScnmId":       "곤충 학명ID",
		"insctSpecsScnm":    "곤충 학명",
		"lastUpdtDtm":       "최종수정일",
		"ordKorNm":          "곤충 목국명",
		"ordNm":             "곤충 목명(Order Name)",
		"stpltScnmRltnCdNm": "곤충 학명 정명/이명 구분",
		"subFalmKorNm":      "곤충 아과국명",
		"subFalmNm":         "곤충 아과명(SubFamily Name)",
		"superFalmKorNm":    "곤충 상과국명",
	}

	for _, tool := range result.Tools {
		if tool.Name != "kini_scnm_info" {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool kini_scnm_info output schema = %#v", tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool kini_scnm_info output schema properties = %#v", schema["properties"])
		}
		checkKeys(t, properties, "item")
		checkPropertyDescriptions(t, "kini_scnm_info", properties, map[string]string{"item": "상세 조회 결과"})

		itemProperty, ok := properties["item"].(map[string]any)
		if !ok {
			t.Fatalf("tool kini_scnm_info item property = %#v", properties["item"])
		}
		itemProperties, nullable := nullableObjectProperties(itemProperty)
		if !nullable {
			t.Errorf("tool kini_scnm_info item schema is not nullable: %#v", itemProperty)
		}
		if itemProperties == nil {
			t.Fatalf("tool kini_scnm_info item schema has no object properties: %#v", itemProperty)
		}
		checkKeys(t, itemProperties, mapKeys(itemDescriptions)...)
		checkPropertyDescriptions(t, "kini_scnm_info item", itemProperties, itemDescriptions)
		return
	}
	t.Fatal("tool kini_scnm_info not found")
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
