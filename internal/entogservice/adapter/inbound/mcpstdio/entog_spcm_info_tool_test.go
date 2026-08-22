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

type entogSpcmInfoUseCaseStub struct {
	query  application.EntogSpcmInfoQuery
	result application.EntogSpcmInfoResult
	err    error
}

func (s *entogSpcmInfoUseCaseStub) EntogSpcmInfo(_ context.Context, query application.EntogSpcmInfoQuery) (application.EntogSpcmInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestEntogSpcmInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &entogSpcmInfoUseCaseStub{result: application.EntogSpcmInfoResult{Item: &application.EntogSpcmInfoItem{
		Btnc: "btnc", ChnNm: " ", ClarHaslvVal: "clarHaslvVal", ClctDyDesc: "clctDyDesc",
		CprtCtnt: "cprtCtnt", EngNm: "engNm", EntogGnrlNm: "entogGnrlNm", EntogPilbkNo: "entogPilbkNo",
		EntogSmplNo: "entogSmplNo", FamilyKorNm: "familyKorNm", FamilyNm: "familyNm", FrstRgstnDtm: "frstRgstnDtm",
		GenusKorNm: "genusKorNm", GenusNm: "genusNm", ImgURL: "NONE", JapNm: "japNm",
		LabelUsgCllcnNmplc: "labelUsgCllcnNmplc   ", LastUpdtDtm: " ", OrdKorNm: "ordKorNm", OrdNm: "ordNm",
		PrkNm: "prkNm", TorsoLngth: "torsoLngth", WingLngth: "wingLngth",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addEntogSpcmInfoTool(server, useCase)
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
	if tool.Name != "entog_service_entog_spcm_info" || tool.Description != "산림청 국립수목원 내구강표본 상세정보를 조회합니다." {
		t.Errorf("tool = %#v", tool)
	}
	checkSchemaProperties(t, tool.InputSchema, map[string]string{
		"q1": "조회키 (entogSpcmSearch 결과의 entogSmplNo)",
	})
	required := tool.InputSchema.(map[string]any)["required"].([]any)
	if len(required) != 1 || required[0] != "q1" {
		t.Errorf("required = %#v, want q1", required)
	}

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "entog_service_entog_spcm_info",
		Arguments: map[string]any{"q1": "test-specimen-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if useCase.query != (application.EntogSpcmInfoQuery{Q1: "test-specimen-number"}) {
		t.Errorf("query = %#v", useCase.query)
	}

	output := result.StructuredContent.(map[string]any)
	item := output["item"].(map[string]any)
	wantItem := map[string]string{
		"btnc": "btnc", "chnNm": " ", "clarHaslvVal": "clarHaslvVal", "clctDyDesc": "clctDyDesc",
		"cprtCtnt": "cprtCtnt", "engNm": "engNm", "entogGnrlNm": "entogGnrlNm", "entogPilbkNo": "entogPilbkNo",
		"entogSmplNo": "entogSmplNo", "familyKorNm": "familyKorNm", "familyNm": "familyNm", "frstRgstnDtm": "frstRgstnDtm",
		"genusKorNm": "genusKorNm", "genusNm": "genusNm", "imgUrl": "NONE", "japNm": "japNm",
		"labelUsgCllcnNmplc": "labelUsgCllcnNmplc   ", "lastUpdtDtm": " ", "ordKorNm": "ordKorNm", "ordNm": "ordNm",
		"prkNm": "prkNm", "torsoLngth": "torsoLngth", "wingLngth": "wingLngth",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if item[key] != want {
			t.Errorf("item %s = %#v, want %q", key, item[key], want)
		}
	}
	checkEntogSpcmInfoOutputSchema(t, tool.OutputSchema)

	useCase.result = application.EntogSpcmInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "entog_service_entog_spcm_info",
		Arguments: map[string]any{"q1": "test-missing-specimen-number"},
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
		Name:      "entog_service_entog_spcm_info",
		Arguments: map[string]any{"q1": "test-specimen-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}

func checkEntogSpcmInfoOutputSchema(t *testing.T, schema any) {
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
		"btnc": "학명", "chnNm": "중국명", "clarHaslvVal": "채집지해발고도", "clctDyDesc": "채집일",
		"cprtCtnt": "저작권", "engNm": "영문명", "entogGnrlNm": "국명", "entogPilbkNo": "도감번호",
		"entogSmplNo": "표본번호", "familyKorNm": "과국명", "familyNm": "과명", "frstRgstnDtm": "최초등록일",
		"genusKorNm": "속국명", "genusNm": "속명", "imgUrl": "이미지URL", "japNm": "일본명",
		"labelUsgCllcnNmplc": "라벨용채집지명", "lastUpdtDtm": "최종수정일", "ordKorNm": "목국명", "ordNm": "목명",
		"prkNm": "북한명", "torsoLngth": "몸통길이", "wingLngth": "날개길이",
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

func TestEntogSpcmInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.EntogSpcmInfoItem{})
	adapterFields := reflect.TypeOf(entogSpcmInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
