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

type alchnSpcmInfoUseCaseStub struct {
	query  application.AlchnSpcmInfoQuery
	result application.AlchnSpcmInfoResult
	err    error
}

func (s *alchnSpcmInfoUseCaseStub) AlchnSpcmInfo(_ context.Context, query application.AlchnSpcmInfoQuery) (application.AlchnSpcmInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestAlchnSpcmInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &alchnSpcmInfoUseCaseStub{result: application.AlchnSpcmInfoResult{Item: &application.AlchnSpcmInfoItem{
		Btnc:          "btnc",
		ClarDtlDscrt:  "clarDtlDscrt",
		CllcrNm:       "cllcrNm",
		CltrNm:        "cltrNm",
		CprtCtnt:      "cprtCtnt",
		EngNm:         "engNm",
		ExmneNm:       "exmneNm",
		FamilyKorNm:   "familyKorNm",
		FamilyNm:      "familyNm",
		FrstRgstnDtm:  "frstRgstnDtm",
		GenusKorNm:    "genusKorNm",
		GenusNm:       "genusNm",
		Grdnt:         "grdnt",
		HaslvVal:      "haslvVal",
		HbttChrcrCont: "hbttChrcrCont",
		ImgURL:        "imgUrl",
		InsttSmplNo:   "insttSmplNo",
		JapNm:         "japNm",
		LastUpdtDtm:   "lastUpdtDtm",
		LchnGnrlNm:    "lchnGnrlNm",
		LchnScnmID:    "lchnScnmId",
		LchnSmplNo:    "lchnSmplNo",
		OrbrnCd:       "orbrnCd",
		PrkNm:         "prkNm",
		SmplCllcnDt:   "smplCllcnDt",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addAlchnSpcmInfoTool(server, useCase)
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
	if tool.Name != "lchn_service_alchn_spcm_info" || tool.Description != "산림청 국립수목원 지의류표본 상세정보를 조회합니다." {
		t.Errorf("tool = %#v", tool)
	}
	checkSchemaProperties(t, tool.InputSchema, map[string]string{"q1": "조회키 (alchnSpcmSearch 결과의 lchnSmplNo)"})
	required := tool.InputSchema.(map[string]any)["required"].([]any)
	if len(required) != 1 || required[0] != "q1" {
		t.Errorf("required = %#v, want q1", required)
	}

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"q1": "TEST-SAMPLE-001"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if useCase.query != (application.AlchnSpcmInfoQuery{Q1: "TEST-SAMPLE-001"}) {
		t.Errorf("query = %#v", useCase.query)
	}

	output := result.StructuredContent.(map[string]any)
	item := output["item"].(map[string]any)
	wantItem := map[string]string{
		"btnc":          "btnc",
		"clarDtlDscrt":  "clarDtlDscrt",
		"cllcrNm":       "cllcrNm",
		"cltrNm":        "cltrNm",
		"cprtCtnt":      "cprtCtnt",
		"engNm":         "engNm",
		"exmneNm":       "exmneNm",
		"familyKorNm":   "familyKorNm",
		"familyNm":      "familyNm",
		"frstRgstnDtm":  "frstRgstnDtm",
		"genusKorNm":    "genusKorNm",
		"genusNm":       "genusNm",
		"grdnt":         "grdnt",
		"haslvVal":      "haslvVal",
		"hbttChrcrCont": "hbttChrcrCont",
		"imgUrl":        "imgUrl",
		"insttSmplNo":   "insttSmplNo",
		"japNm":         "japNm",
		"lastUpdtDtm":   "lastUpdtDtm",
		"lchnGnrlNm":    "lchnGnrlNm",
		"lchnScnmId":    "lchnScnmId",
		"lchnSmplNo":    "lchnSmplNo",
		"orbrnCd":       "orbrnCd",
		"prkNm":         "prkNm",
		"smplCllcnDt":   "smplCllcnDt",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if item[key] != want {
			t.Errorf("item %s = %#v, want %q", key, item[key], want)
		}
	}
	checkAlchnSpcmInfoOutputSchema(t, tool.OutputSchema)

	useCase.result = application.AlchnSpcmInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"q1": "KNA-MCP-NO-RESULT"}})
	if err != nil {
		t.Fatal(err)
	}
	output = result.StructuredContent.(map[string]any)
	if output["item"] != nil {
		t.Errorf("item = %#v, want nil", output["item"])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: tool.Name, Arguments: map[string]any{"q1": "TEST-SAMPLE-001"}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("tool result is not an error")
	}
}

func checkAlchnSpcmInfoOutputSchema(t *testing.T, schema any) {
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
		"btnc":          "학명",
		"clarDtlDscrt":  "채집지상세설명",
		"cllcrNm":       "채집자명",
		"cltrNm":        "채집자그룹멤버명",
		"cprtCtnt":      "저작권",
		"engNm":         "영문명",
		"exmneNm":       "조사자명",
		"familyKorNm":   "과국명",
		"familyNm":      "과명",
		"frstRgstnDtm":  "최초등록일시",
		"genusKorNm":    "속국명",
		"genusNm":       "속명",
		"grdnt":         "경사도",
		"haslvVal":      "해발고도값",
		"hbttChrcrCont": "기물설명",
		"imgUrl":        "이미지URL",
		"insttSmplNo":   "기관표본번호",
		"japNm":         "일어명",
		"lastUpdtDtm":   "최종수정일시",
		"lchnGnrlNm":    "국명",
		"lchnScnmId":    "학명ID",
		"lchnSmplNo":    "표본번호",
		"orbrnCd":       "방위코드",
		"prkNm":         "북한명",
		"smplCllcnDt":   "표본채집일자",
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

func TestAlchnSpcmInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.AlchnSpcmInfoItem{})
	adapterFields := reflect.TypeOf(alchnSpcmInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
