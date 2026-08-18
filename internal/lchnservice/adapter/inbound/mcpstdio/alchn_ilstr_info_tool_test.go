package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type alchnIlstrInfoUseCaseStub struct {
	query  application.AlchnIlstrInfoQuery
	result application.AlchnIlstrInfoResult
	err    error
}

func (s *alchnIlstrInfoUseCaseStub) AlchnIlstrInfo(_ context.Context, query application.AlchnIlstrInfoQuery) (application.AlchnIlstrInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestAlchnIlstrInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &alchnIlstrInfoUseCaseStub{result: application.AlchnIlstrInfoResult{Item: &application.AlchnIlstrInfoItem{
		Btnc: "btnc", Cont1: " ", Cont2: "cont2", Cont3: "cont3", Cont4: "cont4", Cont5: "cont5", Cont6: "cont6",
		Cont7: "cont7", Cont8: "cont8", Cont9: "cont9", Cont10: "cont10", Cont11: "cont11", Cont12: "cont12",
		CprtCtnt: "cprtCtnt", EngNm: "engNm", FamilyKorNm: "familyKorNm", FamilyNm: "familyNm",
		FrstRgstnDtm: "frstRgstnDtm", GenusKorNm: "genusKorNm", GenusNm: "genusNm", ImgURL: "imgUrl",
		JapNm: "japNm", LastUpdtDtm: "lastUpdtDtm", LchnGnrlNm: "lchnGnrlNm", LchnInfrpNm: "lchnInfrpNm",
		LchnPilbkNo: "lchnPilbkNo", LchnScnmID: "lchnScnmId", LchnTtnm: "lchnTtnm", PrkNm: "prkNm",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addAlchnIlstrInfoTool(server, useCase)
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
	if tool.Name != "lchn_service_alchn_ilstr_info" || tool.Description != "산림청 국립수목원 지의류도감 상세정보를 조회합니다." {
		t.Errorf("tool = %#v", tool)
	}
	checkSchemaProperties(t, tool.InputSchema, map[string]string{"q1": "조회키 (alchnIlstrSearch 결과의 lchnPilbkNo)"})
	required := tool.InputSchema.(map[string]any)["required"].([]any)
	if len(required) != 1 || required[0] != "q1" {
		t.Errorf("required = %#v, want q1", required)
	}

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"q1": "test-lichen-pictorial-book-number"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if useCase.query != (application.AlchnIlstrInfoQuery{Q1: "test-lichen-pictorial-book-number"}) {
		t.Errorf("query = %#v", useCase.query)
	}

	output := result.StructuredContent.(map[string]any)
	item := output["item"].(map[string]any)
	wantItem := map[string]string{
		"btnc": "btnc", "cont1": " ", "cont2": "cont2", "cont3": "cont3", "cont4": "cont4", "cont5": "cont5",
		"cont6": "cont6", "cont7": "cont7", "cont8": "cont8", "cont9": "cont9", "cont10": "cont10",
		"cont11": "cont11", "cont12": "cont12", "cprtCtnt": "cprtCtnt", "engNm": "engNm",
		"familyKorNm": "familyKorNm", "familyNm": "familyNm", "frstRgstnDtm": "frstRgstnDtm",
		"genusKorNm": "genusKorNm", "genusNm": "genusNm", "imgUrl": "imgUrl", "japNm": "japNm",
		"lastUpdtDtm": "lastUpdtDtm", "lchnGnrlNm": "lchnGnrlNm", "lchnInfrpNm": "lchnInfrpNm",
		"lchnPilbkNo": "lchnPilbkNo", "lchnScnmId": "lchnScnmId", "lchnTtnm": "lchnTtnm", "prkNm": "prkNm",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if item[key] != want {
			t.Errorf("item %s = %#v, want %q", key, item[key], want)
		}
	}
	checkAlchnIlstrInfoOutputSchema(t, tool.OutputSchema)

	useCase.result = application.AlchnIlstrInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"q1": "test-missing-lichen-pictorial-book-number"}})
	if err != nil {
		t.Fatal(err)
	}
	output = result.StructuredContent.(map[string]any)
	if output["item"] != nil {
		t.Errorf("item = %#v, want nil", output["item"])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"q1": "test-lichen-pictorial-book-number"}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}

func checkAlchnIlstrInfoOutputSchema(t *testing.T, schema any) {
	t.Helper()
	properties := schema.(map[string]any)["properties"].(map[string]any)
	if len(properties) != 1 || properties["item"] == nil {
		t.Fatalf("properties = %#v, want item", properties)
	}
	itemProperty := properties["item"].(map[string]any)
	if itemProperty["description"] != "상세 조회 결과" {
		t.Errorf("item description = %#v", itemProperty["description"])
	}
	itemProperties, nullable := nullableObjectProperties(itemProperty)
	if !nullable || itemProperties == nil {
		t.Fatalf("item schema = %#v, want nullable object", itemProperty)
	}
	wantDescriptions := map[string]string{
		"btnc": "학명", "cont1": "영문설명", "cont2": "미사용", "cont3": "형태에 의한 분류", "cont4": "지의류형태(국문 종기술)",
		"cont5": "미사용", "cont6": "미사용", "cont7": "미사용", "cont8": "미사용", "cont9": "지의물질",
		"cont10": "분포", "cont11": "미사용", "cont12": "비고", "cprtCtnt": "저작권", "engNm": "영문명",
		"familyKorNm": "과국명", "familyNm": "과명", "frstRgstnDtm": "최초등록일시", "genusKorNm": "속국명",
		"genusNm": "속명", "imgUrl": "이미지URL", "japNm": "일어명", "lastUpdtDtm": "최종수정일시",
		"lchnGnrlNm": "국명", "lchnInfrpNm": "종하명", "lchnPilbkNo": "도감번호", "lchnScnmId": "학명ID",
		"lchnTtnm": "종소명", "prkNm": "북한명",
	}
	if len(itemProperties) != len(wantDescriptions) {
		t.Errorf("item property count = %d, want %d", len(itemProperties), len(wantDescriptions))
	}
	for name, want := range wantDescriptions {
		property := itemProperties[name].(map[string]any)
		if property["description"] != want {
			t.Errorf("property %s description = %#v, want %q", name, property["description"], want)
		}
	}
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

func TestAlchnIlstrInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.AlchnIlstrInfoItem{})
	adapterFields := reflect.TypeOf(alchnIlstrInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
