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

type insectPilbkInfoUseCaseStub struct {
	query  application.InsectPilbkInfoQuery
	result application.InsectPilbkInfoResult
	err    error
}

func (s *insectPilbkInfoUseCaseStub) InsectPilbkInfo(_ context.Context, query application.InsectPilbkInfoQuery) (application.InsectPilbkInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestInsectPilbkInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &insectPilbkInfoUseCaseStub{result: application.InsectPilbkInfoResult{Item: &application.InsectPilbkInfoItem{
		EcoDsrct:         "ecology description",
		EggDsrct:         "egg description",
		EmrgcCnt:         "emergence count",
		EmrgcEraDscrt:    "emergence era description",
		FamilyKorNm:      "family Korean name",
		FamilyNm:         "family name",
		FemaleDsrct:      "female description",
		GenusKorNm:       "genus Korean name",
		GenusNm:          "genus name",
		GnrlDsrct:        "general description",
		HabitDsrct:       "habit description",
		InsctEngNm:       "insect English name",
		InsctGnrlNm:      "insect general name",
		InsctPilbkNo:     "test-insect-pictorial-book-number",
		InsctSpecsScnm:   "insect species scientific name",
		LarvaDsrct:       "larva description",
		LastUpdtDtm:      "last update date time",
		MaleDsrct:        "male description",
		MnmmOccrrCnt:     "minimum occurrence count",
		MxmmOccrrCnt:     "maximum occurrence count",
		OrdKorNm:         "order Korean name",
		OrdNm:            "order name",
		PestDsrct:        "pest control description",
		PupaDsrct:        "pupa description",
		ReferDsrct:       "reference description",
		SubFamilyKorNm:   "subfamily Korean name",
		SubFamilyNm:      "subfamily name",
		SuperFamilyKorNm: "superfamily Korean name",
		SuperFamilyNm:    "superfamily name",
		WinterDsrct:      "winter description",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addInsectPilbkInfoTool(server, useCase)
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

	checkToolInputSchema(t, ctx, clientSession, "insect_resource_insect_pilbk_info", map[string]string{
		"reqInsctPilbkNo": "조회할 곤충 도감번호 (insectPilbkSearch 결과의 insctPilbkNo)",
	}, []string{"reqInsctPilbkNo"})
	checkToolDescription(t, ctx, clientSession, "insect_resource_insect_pilbk_info", "산림청 국립수목원 곤충도감 상세정보를 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "insect_resource_insect_pilbk_info",
		Arguments: map[string]any{"reqInsctPilbkNo": "test-insect-pictorial-book-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if want := (application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "test-insect-pictorial-book-number"}); useCase.query != want {
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
		"ecoDsrct":         "ecology description",
		"eggDsrct":         "egg description",
		"emrgcCnt":         "emergence count",
		"emrgcEraDscrt":    "emergence era description",
		"familyKorNm":      "family Korean name",
		"familyNm":         "family name",
		"femaleDsrct":      "female description",
		"genusKorNm":       "genus Korean name",
		"genusNm":          "genus name",
		"gnrlDsrct":        "general description",
		"habitDsrct":       "habit description",
		"insctEngNm":       "insect English name",
		"insctGnrlNm":      "insect general name",
		"insctPilbkNo":     "test-insect-pictorial-book-number",
		"insctSpecsScnm":   "insect species scientific name",
		"larvaDsrct":       "larva description",
		"lastUpdtDtm":      "last update date time",
		"maleDsrct":        "male description",
		"mnmmOccrrCnt":     "minimum occurrence count",
		"mxmmOccrrCnt":     "maximum occurrence count",
		"ordKorNm":         "order Korean name",
		"ordNm":            "order name",
		"pestDsrct":        "pest control description",
		"pupaDsrct":        "pupa description",
		"referDsrct":       "reference description",
		"subFamilyKorNm":   "subfamily Korean name",
		"subFamilyNm":      "subfamily name",
		"superFamilyKorNm": "superfamily Korean name",
		"superFamilyNm":    "superfamily name",
		"winterDsrct":      "winter description",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	checkInsectPilbkInfoOutputSchema(t, ctx, clientSession)

	useCase.result = application.InsectPilbkInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "insect_resource_insect_pilbk_info",
		Arguments: map[string]any{"reqInsctPilbkNo": "not-found"},
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
		Name:      "insect_resource_insect_pilbk_info",
		Arguments: map[string]any{"reqInsctPilbkNo": "test-insect-pictorial-book-number"},
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

func checkInsectPilbkInfoOutputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	itemDescriptions := map[string]string{
		"ecoDsrct":         "생태",
		"eggDsrct":         "알",
		"emrgcCnt":         "출현수",
		"emrgcEraDscrt":    "출현시기설명",
		"familyKorNm":      "과국명",
		"familyNm":         "과명",
		"femaleDsrct":      "성충(암)",
		"genusKorNm":       "속국명",
		"genusNm":          "속명",
		"gnrlDsrct":        "일반 특징",
		"habitDsrct":       "습성",
		"insctEngNm":       "영문명",
		"insctGnrlNm":      "국명(곤충명)",
		"insctPilbkNo":     "곤충도감번호",
		"insctSpecsScnm":   "학명",
		"larvaDsrct":       "유충",
		"lastUpdtDtm":      "최종수정일",
		"maleDsrct":        "성충(수)",
		"mnmmOccrrCnt":     "최소발생수",
		"mxmmOccrrCnt":     "최대발생수",
		"ordKorNm":         "목국명",
		"ordNm":            "목명",
		"pestDsrct":        "방제법",
		"pupaDsrct":        "번데기",
		"referDsrct":       "참고사항",
		"subFamilyKorNm":   "아과국명",
		"subFamilyNm":      "아과명",
		"superFamilyKorNm": "상과국명",
		"superFamilyNm":    "상과명",
		"winterDsrct":      "월동",
	}

	for _, tool := range result.Tools {
		if tool.Name != "insect_resource_insect_pilbk_info" {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool insect_resource_insect_pilbk_info output schema = %#v", tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool insect_resource_insect_pilbk_info output schema properties = %#v", schema["properties"])
		}
		checkKeys(t, properties, "item")
		checkPropertyDescriptions(t, "insect_resource_insect_pilbk_info", properties, map[string]string{"item": "상세 조회 결과"})

		itemProperty, ok := properties["item"].(map[string]any)
		if !ok {
			t.Fatalf("tool insect_resource_insect_pilbk_info item property = %#v", properties["item"])
		}
		itemProperties, nullable := nullableObjectProperties(itemProperty)
		if !nullable {
			t.Errorf("tool insect_resource_insect_pilbk_info item schema is not nullable: %#v", itemProperty)
		}
		if itemProperties == nil {
			t.Fatalf("tool insect_resource_insect_pilbk_info item schema has no object properties: %#v", itemProperty)
		}
		checkKeys(t, itemProperties, mapKeys(itemDescriptions)...)
		checkPropertyDescriptions(t, "insect_resource_insect_pilbk_info item", itemProperties, itemDescriptions)
		return
	}
	t.Fatal("tool insect_resource_insect_pilbk_info not found")
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

func TestInsectPilbkInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.InsectPilbkInfoItem{})
	adapterFields := reflect.TypeOf(insectPilbkInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
