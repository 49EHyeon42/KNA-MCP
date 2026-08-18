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

type childPilbkInfoUseCaseStub struct {
	query  application.ChildPilbkInfoQuery
	result application.ChildPilbkInfoResult
	err    error
}

func (s *childPilbkInfoUseCaseStub) ChildPilbkInfo(_ context.Context, query application.ChildPilbkInfoQuery) (application.ChildPilbkInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestChildPilbkInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &childPilbkInfoUseCaseStub{result: application.ChildPilbkInfoResult{Item: &application.ChildPilbkInfoItem{
		BiogyNm:           "biology name",
		ChildLvbngPilbkNo: "child pictorial book number",
		ExtrmCrss:         "extinction crisis",
		FamilyKorNm:       "family Korean name",
		FamilyNm:          "family name",
		GenusKorNm:        "genus Korean name",
		GenusNm:           "genus name",
		HbttFieldYn:       "field habitat flag",
		HbttFrestYn:       "forest habitat flag",
		HbttRiverYn:       "river habitat flag",
		LvbngDscrt:        "living thing description <br/> next",
		LvbngTpcdNm:       "living thing type code name",
		LvngKrlngNm:       "living thing Korean name",
		PrtctSpecsTpcdNm:  " ",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addChildPilbkInfoTool(server, useCase)
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

	checkToolInputSchema(t, ctx, clientSession, "child_service_child_pilbk_info", map[string]string{
		"reqChildLvbngPilbkNo": "도감번호(childLvbngPilbkNo) (childPilbkSearch 결과의 childLvbngPilbkNo)",
	}, []string{"reqChildLvbngPilbkNo"})
	checkToolDescription(t, ctx, clientSession, "child_service_child_pilbk_info", "산림청 국립수목원 어린이생물도감 상세정보를 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "child_service_child_pilbk_info",
		Arguments: map[string]any{"reqChildLvbngPilbkNo": "test-child-pictorial-book-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if want := (application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "test-child-pictorial-book-number"}); useCase.query != want {
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
		"biogyNm":           "biology name",
		"childLvbngPilbkNo": "child pictorial book number",
		"extrmCrss":         "extinction crisis",
		"familyKorNm":       "family Korean name",
		"familyNm":          "family name",
		"genusKorNm":        "genus Korean name",
		"genusNm":           "genus name",
		"hbttFieldYn":       "field habitat flag",
		"hbttFrestYn":       "forest habitat flag",
		"hbttRiverYn":       "river habitat flag",
		"lvbngDscrt":        "living thing description <br/> next",
		"lvbngTpcdNm":       "living thing type code name",
		"lvngKrlngNm":       "living thing Korean name",
		"prtctSpecsTpcdNm":  " ",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	checkChildPilbkInfoOutputSchema(t, ctx, clientSession)

	useCase.result = application.ChildPilbkInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "child_service_child_pilbk_info",
		Arguments: map[string]any{"reqChildLvbngPilbkNo": "999999999999"},
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
		Name:      "child_service_child_pilbk_info",
		Arguments: map[string]any{"reqChildLvbngPilbkNo": "test-child-pictorial-book-number"},
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

func checkChildPilbkInfoOutputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	itemDescriptions := map[string]string{
		"biogyNm":           "생물학명",
		"childLvbngPilbkNo": "어린이생물도감번호",
		"extrmCrss":         "멸종위기종구분",
		"familyKorNm":       "과국명",
		"familyNm":          "과명",
		"genusKorNm":        "속국명",
		"genusNm":           "속명",
		"hbttFieldYn":       "서식지 들 여부",
		"hbttFrestYn":       "서식지 숲 여부",
		"hbttRiverYn":       "서식지 강 여부",
		"lvbngDscrt":        "생물설명",
		"lvbngTpcdNm":       "생물분류",
		"lvngKrlngNm":       "생물국명",
		"prtctSpecsTpcdNm":  "보호종구분",
	}

	for _, tool := range result.Tools {
		if tool.Name != "child_service_child_pilbk_info" {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool child_service_child_pilbk_info output schema = %#v", tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool child_service_child_pilbk_info output schema properties = %#v", schema["properties"])
		}
		checkKeys(t, properties, "item")
		checkPropertyDescriptions(t, "child_service_child_pilbk_info", properties, map[string]string{"item": "상세 조회 결과"})

		itemProperty, ok := properties["item"].(map[string]any)
		if !ok {
			t.Fatalf("tool child_service_child_pilbk_info item property = %#v", properties["item"])
		}
		itemProperties, nullable := nullableObjectProperties(itemProperty)
		if !nullable {
			t.Errorf("tool child_service_child_pilbk_info item schema is not nullable: %#v", itemProperty)
		}
		if itemProperties == nil {
			t.Fatalf("tool child_service_child_pilbk_info item schema has no object properties: %#v", itemProperty)
		}
		checkKeys(t, itemProperties, mapKeys(itemDescriptions)...)
		checkPropertyDescriptions(t, "child_service_child_pilbk_info item", itemProperties, itemDescriptions)
		return
	}
	t.Fatal("tool child_service_child_pilbk_info not found")
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

func TestChildPilbkInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.ChildPilbkInfoItem{})
	adapterFields := reflect.TypeOf(childPilbkInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
