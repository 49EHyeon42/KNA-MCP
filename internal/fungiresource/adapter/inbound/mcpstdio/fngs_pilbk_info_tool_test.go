package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type fngsPilbkInfoUseCaseStub struct {
	query  application.FngsPilbkInfoQuery
	result application.FngsPilbkInfoResult
	err    error
}

func (s *fngsPilbkInfoUseCaseStub) FngsPilbkInfo(_ context.Context, query application.FngsPilbkInfoQuery) (application.FngsPilbkInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestFngsPilbkInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &fngsPilbkInfoUseCaseStub{result: application.FngsPilbkInfoResult{Item: &application.FngsPilbkInfoItem{
		MshrmColorCdNm:      "mushroom color code name",
		CrpphFomTpcdNm:      "carpophore form type code name",
		FamilyKorNm:         "family Korean name",
		FamilyNm:            "family name",
		FngsEclgTpcdNm:      "fungi ecology type code name",
		FngsGnrlNm:          "fungi general name",
		FngsPilbkNo:         "test-fungi-pictorial-book-number",
		FngsPrpseTpcdNm:     "fungi purpose type code name",
		FngsScnm:            "fungi scientific name",
		GenusKorNm:          "genus Korean name",
		GenusNm:             "genus name",
		GrwEvrntDesc:        "growth environment description",
		LastUpdtDtm:         "last update date time",
		MicroShpe:           "microscopic shape",
		MshrmTpcdNm:         "mushroom type code name",
		OccrrSsnNm:          "occurrence season name",
		RsrcActoClsscTpcdNm: "resource classification type code name",
		Shpe:                "shape",
	}}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addFngsPilbkInfoTool(server, useCase)
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

	checkToolInputSchema(t, ctx, clientSession, "fungi_resource_fngs_pilbk_info", map[string]string{
		"reqFngsPilbkNo": "검색할 버섯도감의 버섯도감번호 (fngsPilbkSearch 결과의 fngsPilbkNo)",
	}, []string{"reqFngsPilbkNo"})
	checkToolDescription(t, ctx, clientSession, "fungi_resource_fngs_pilbk_info", "산림청 국립수목원 버섯도감 상세정보를 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "fungi_resource_fngs_pilbk_info",
		Arguments: map[string]any{"reqFngsPilbkNo": "test-fungi-pictorial-book-number"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if want := (application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "test-fungi-pictorial-book-number"}); useCase.query != want {
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
		"mshrmColorCdNm":      "mushroom color code name",
		"crpphFomTpcdNm":      "carpophore form type code name",
		"familyKorNm":         "family Korean name",
		"familyNm":            "family name",
		"fngsEclgTpcdNm":      "fungi ecology type code name",
		"fngsGnrlNm":          "fungi general name",
		"fngsPilbkNo":         "test-fungi-pictorial-book-number",
		"fngsPrpseTpcdNm":     "fungi purpose type code name",
		"fngsScnm":            "fungi scientific name",
		"genusKorNm":          "genus Korean name",
		"genusNm":             "genus name",
		"grwEvrntDesc":        "growth environment description",
		"lastUpdtDtm":         "last update date time",
		"microShpe":           "microscopic shape",
		"mshrmTpcdNm":         "mushroom type code name",
		"occrrSsnNm":          "occurrence season name",
		"rsrcActoClsscTpcdNm": "resource classification type code name",
		"shpe":                "shape",
	}
	checkKeys(t, item, mapKeys(wantItem)...)
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	checkFngsPilbkInfoOutputSchema(t, ctx, clientSession)

	useCase.result = application.FngsPilbkInfoResult{}
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "fungi_resource_fngs_pilbk_info",
		Arguments: map[string]any{"reqFngsPilbkNo": "not-found"},
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
		Name:      "fungi_resource_fngs_pilbk_info",
		Arguments: map[string]any{"reqFngsPilbkNo": "test-fungi-pictorial-book-number"},
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

func checkFngsPilbkInfoOutputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	itemDescriptions := map[string]string{
		"mshrmColorCdNm":      "버섯색상",
		"crpphFomTpcdNm":      "자실체형태",
		"familyKorNm":         "과국명",
		"familyNm":            "과명",
		"fngsEclgTpcdNm":      "버섯생태형",
		"fngsGnrlNm":          "국명(버섯명)",
		"fngsPilbkNo":         "버섯도감번호",
		"fngsPrpseTpcdNm":     "버섯용도",
		"fngsScnm":            "학명",
		"genusKorNm":          "속국명",
		"genusNm":             "속명",
		"grwEvrntDesc":        "발생장소",
		"lastUpdtDtm":         "최종수정일",
		"microShpe":           "현미경적 특징",
		"mshrmTpcdNm":         "버섯구분",
		"occrrSsnNm":          "발생계절",
		"rsrcActoClsscTpcdNm": "자원분류",
		"shpe":                "외부 형태적 특징",
	}

	for _, tool := range result.Tools {
		if tool.Name != "fungi_resource_fngs_pilbk_info" {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool fungi_resource_fngs_pilbk_info output schema = %#v", tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool fungi_resource_fngs_pilbk_info output schema properties = %#v", schema["properties"])
		}
		checkKeys(t, properties, "item")
		checkPropertyDescriptions(t, "fungi_resource_fngs_pilbk_info", properties, map[string]string{"item": "상세 조회 결과"})

		itemProperty, ok := properties["item"].(map[string]any)
		if !ok {
			t.Fatalf("tool fungi_resource_fngs_pilbk_info item property = %#v", properties["item"])
		}
		itemProperties, nullable := nullableObjectProperties(itemProperty)
		if !nullable {
			t.Errorf("tool fungi_resource_fngs_pilbk_info item schema is not nullable: %#v", itemProperty)
		}
		if itemProperties == nil {
			t.Fatalf("tool fungi_resource_fngs_pilbk_info item schema has no object properties: %#v", itemProperty)
		}
		checkKeys(t, itemProperties, mapKeys(itemDescriptions)...)
		checkPropertyDescriptions(t, "fungi_resource_fngs_pilbk_info item", itemProperties, itemDescriptions)
		return
	}
	t.Fatal("tool fungi_resource_fngs_pilbk_info not found")
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

func TestFngsPilbkInfoOutputFieldNamesMatchApplication(t *testing.T) {
	applicationFields := reflect.TypeOf(application.FngsPilbkInfoItem{})
	adapterFields := reflect.TypeOf(fngsPilbkInfoItem{})
	if applicationFields.NumField() != adapterFields.NumField() {
		t.Fatalf("application fields = %d, adapter fields = %d", applicationFields.NumField(), adapterFields.NumField())
	}
	for index := range applicationFields.NumField() {
		if applicationFields.Field(index).Name != adapterFields.Field(index).Name {
			t.Errorf("field %d = %s, want %s", index, adapterFields.Field(index).Name, applicationFields.Field(index).Name)
		}
	}
}
