package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type entogIlstrInfoUseCaseStub struct {
	query  application.EntogIlstrInfoQuery
	result application.EntogIlstrInfoResult
	err    error
}

func (s *entogIlstrInfoUseCaseStub) EntogIlstrInfo(_ context.Context, query application.EntogIlstrInfoQuery) (application.EntogIlstrInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestEntogIlstrInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &entogIlstrInfoUseCaseStub{result: application.EntogIlstrInfoResult{Item: &application.EntogIlstrInfoItem{
		Btnc: "btnc", Cont1: "cont1", Cont2: "cont2", Cont3: "cont3", Cont4: "cont4", Cont5: "cont5",
		Cont6: "cont6", Cont7: "cont7", Cont8: "cont8", Cont9: "cont9", Cont10: "cont10", Cont11: "cont11",
		CprtCtnt: "cprtCtnt", EmrgcCnt: "emrgcCnt", EmrgcEraDscrt: "emrgcEraDscrt", EntogAthrNm: "entogAthrNm",
		EntogEngNm: "entogEngNm", EntogOfnmKrlngNm: "entogOfnmKrlngNm", EntogPilbkNo: "entogPilbkNo",
		EntogSpecsNm: "entogSpecsNm", FamilyKorNm: "familyKorNm", FamilyNm: "familyNm", GenusKorNm: "genusKorNm",
		GenusNm: "genusNm", ImgURL: "imgUrl", MnmmOccrrCnt: "mnmmOccrrCnt", MxmmOccrrCnt: "mxmmOccrrCnt",
		OrdKorNm: "ordKorNm", OrdNm: "ordNm",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addEntogIlstrInfoTool(server, useCase)
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
	if tool.Name != "entog_service_entog_ilstr_info" || tool.Description != "산림청 국립수목원 내구강도감 상세정보를 조회합니다." {
		t.Errorf("tool = %#v", tool)
	}
	checkSchemaProperties(t, tool.InputSchema, map[string]string{
		"q1": "조회키 (entogIlstrSearch 결과의 entogPilbkNo)",
	})
	required := tool.InputSchema.(map[string]any)["required"].([]any)
	if len(required) != 1 || required[0] != "q1" {
		t.Errorf("required = %#v, want q1", required)
	}

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "entog_service_entog_ilstr_info",
		Arguments: map[string]any{"q1": "test-entognath-pictorial-book-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if useCase.query != (application.EntogIlstrInfoQuery{Q1: "test-entognath-pictorial-book-number"}) {
		t.Errorf("query = %#v", useCase.query)
	}

	output := result.StructuredContent.(map[string]any)
	item := output["item"].(map[string]any)
	wantItem := map[string]string{
		"btnc": "btnc", "cont1": "cont1", "cont2": "cont2", "cont3": "cont3", "cont4": "cont4",
		"cont5": "cont5", "cont6": "cont6", "cont7": "cont7", "cont8": "cont8", "cont9": "cont9",
		"cont10": "cont10", "cont11": "cont11", "cprtCtnt": "cprtCtnt", "emrgcCnt": "emrgcCnt",
		"emrgcEraDscrt": "emrgcEraDscrt", "entogAthrNm": "entogAthrNm", "entogEngNm": "entogEngNm",
		"entogOfnmKrlngNm": "entogOfnmKrlngNm", "entogPilbkNo": "entogPilbkNo", "entogSpecsNm": "entogSpecsNm",
		"familyKorNm": "familyKorNm", "familyNm": "familyNm", "genusKorNm": "genusKorNm", "genusNm": "genusNm",
		"imgUrl": "imgUrl", "mnmmOccrrCnt": "mnmmOccrrCnt", "mxmmOccrrCnt": "mxmmOccrrCnt",
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
	checkEntogIlstrInfoOutputSchema(t, tool.OutputSchema)

	useCase.result = application.EntogIlstrInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "entog_service_entog_ilstr_info",
		Arguments: map[string]any{"q1": "test-missing-entognath-pictorial-book-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	output = result.StructuredContent.(map[string]any)
	if output["item"] != nil {
		t.Errorf("item = %#v, want nil", output["item"])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "entog_service_entog_ilstr_info",
		Arguments: map[string]any{"q1": "test-entognath-pictorial-book-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}

func checkEntogIlstrInfoOutputSchema(t *testing.T, schema any) {
	t.Helper()
	properties := schema.(map[string]any)["properties"].(map[string]any)
	if len(properties) != 1 || properties["item"] == nil {
		t.Fatalf("properties = %#v, want item", properties)
	}
	itemProperty := properties["item"].(map[string]any)
	if itemProperty["description"] != "상세 조회 결과" {
		t.Errorf("item description = %#v", itemProperty["description"])
	}
	itemProperties, nullable := entogNullableObjectProperties(itemProperty)
	if !nullable || itemProperties == nil {
		t.Fatalf("item schema = %#v, want nullable object", itemProperty)
	}
	wantDescriptions := map[string]string{
		"btnc": "학명", "cont1": "일반특징", "cont2": "성충(수)", "cont3": "성충(암)", "cont4": "번데기",
		"cont5": "유충", "cont6": "참고사항", "cont7": "생태", "cont8": "습성", "cont9": "월동",
		"cont10": "방제법", "cont11": "알", "cprtCtnt": "저작권", "emrgcCnt": "출현수",
		"emrgcEraDscrt": "출현시기설명", "entogAthrNm": "명명자명", "entogEngNm": "영문명",
		"entogOfnmKrlngNm": "국명", "entogPilbkNo": "도감번호", "entogSpecsNm": "종소명",
		"familyKorNm": "과국명", "familyNm": "과명", "genusKorNm": "속국명", "genusNm": "속명",
		"imgUrl": "이미지URL", "mnmmOccrrCnt": "최소발생횟수", "mxmmOccrrCnt": "최대발생횟수",
		"ordKorNm": "목국명", "ordNm": "목명",
	}
	if len(itemProperties) != len(wantDescriptions) {
		t.Errorf("item property count = %d, want %d", len(itemProperties), len(wantDescriptions))
	}
	for name, want := range wantDescriptions {
		property, ok := itemProperties[name].(map[string]any)
		if !ok {
			t.Errorf("missing property %q", name)
			continue
		}
		if property["description"] != want {
			t.Errorf("property %s description = %#v, want %q", name, property["description"], want)
		}
	}
}

func entogNullableObjectProperties(schema map[string]any) (map[string]any, bool) {
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

func TestEntogIlstrInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.EntogIlstrInfoItem{})
	adapterFields := reflect.TypeOf(entogIlstrInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
