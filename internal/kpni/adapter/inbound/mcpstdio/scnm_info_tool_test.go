package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
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
		APGFalmKorNm:       "APG family Korean name",
		APGFalmNm:          "APG family name",
		BiogyNmTpcdNm:      "biology name type code name",
		CltvaYn:            "cultivation yes or no",
		EclgDstrbYn:        "ecological disturbance yes or no",
		ExtcCncrnsYn:       "exotic concern yes or no",
		ExtcPlantCdNm:      "exotic plant code name",
		ExtcPlantYn:        "exotic plant yes or no",
		FalmKorNm:          "family Korean name",
		FalmNm:             "family name",
		GenusKorNm:         "genus Korean name",
		GenusNm:            "genus name",
		LtrtrInfrmNm:       "literature information name",
		PlantBrdgFomTpcdNm: "plant breeding form type code name",
		PlantChnNm:         "plant Chinese name",
		PlantEngNm:         "plant English name",
		PlantGnrlNm:        "plant general name",
		PlantGnrlNm2:       "plant general name 2",
		PlantJpnNm:         "plant Japanese name",
		PlantScnmID:        "1004701",
		PlantSpecsScnm:     "plant species scientific name",
		RareTpcdNm:         "rare type code name",
		RelPlantSpecsScnm:  "related plant species scientific name",
		RelScnmTpcdNm:      "related scientific name type code name",
		Rmrk:               "remark",
		RrnssPlantYn:       "rareness plant yes or no",
		SpcltPlantCdNm:     "specialty plant code name",
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

	checkToolInputSchema(t, ctx, clientSession, "kpni_scnm_info", map[string]string{
		"reqPlantScnmId": "검색하려는 식물 학명ID (scnmSearch 결과의 plantScnmId)",
	}, []string{"reqPlantScnmId"})
	checkToolDescription(t, ctx, clientSession, "kpni_scnm_info", "산림청 국립수목원 국가표준식물목록의 식물 학명 상세 정보를 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "kpni_scnm_info",
		Arguments: map[string]any{"reqPlantScnmId": "1004701"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}
	if want := (application.ScnmInfoQuery{ReqPlantScnmID: "1004701"}); useCase.query != want {
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
		"apgFalmKorNm":       "APG family Korean name",
		"apgFalmNm":          "APG family name",
		"biogyNmTpcdNm":      "biology name type code name",
		"cltvaYn":            "cultivation yes or no",
		"eclgDstrbYn":        "ecological disturbance yes or no",
		"extcCncrnsYn":       "exotic concern yes or no",
		"extcPlantCdNm":      "exotic plant code name",
		"extcPlantYn":        "exotic plant yes or no",
		"falmKorNm":          "family Korean name",
		"falmNm":             "family name",
		"genusKorNm":         "genus Korean name",
		"genusNm":            "genus name",
		"ltrtrInfrmNm":       "literature information name",
		"plantBrdgFomTpcdNm": "plant breeding form type code name",
		"plantChnNm":         "plant Chinese name",
		"plantEngNm":         "plant English name",
		"plantGnrlNm":        "plant general name",
		"plantGnrlNm2":       "plant general name 2",
		"plantJpnNm":         "plant Japanese name",
		"plantScnmId":        "1004701",
		"plantSpecsScnm":     "plant species scientific name",
		"rareTpcdNm":         "rare type code name",
		"relPlantSpecsScnm":  "related plant species scientific name",
		"relScnmTpcdNm":      "related scientific name type code name",
		"rmrk":               "remark",
		"rrnssPlantYn":       "rareness plant yes or no",
		"spcltPlantCdNm":     "specialty plant code name",
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
		Name:      "kpni_scnm_info",
		Arguments: map[string]any{"reqPlantScnmId": "not-found"},
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
		Name:      "kpni_scnm_info",
		Arguments: map[string]any{"reqPlantScnmId": "1004701"},
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
		"apgFalmKorNm":       "식물 학명 APG 분류군 과국명",
		"apgFalmNm":          "식물 학명 APG 분류군 과명(Family Name)",
		"biogyNmTpcdNm":      "식물 학명 정명/이명/서명 등의 구분명",
		"cltvaYn":            "재배여부",
		"eclgDstrbYn":        "생태계 교란종 여부",
		"extcCncrnsYn":       "외래화 우려 여부",
		"extcPlantCdNm":      "외래 식물 구분명",
		"extcPlantYn":        "침입 외래 식물 여부",
		"falmKorNm":          "식물 학명 분류군 과국명",
		"falmNm":             "식물 학명 분류군 과명(Family Name)",
		"genusKorNm":         "식물 학명 분류군 속국명",
		"genusNm":            "식물 학명 분류군 속명(Genus Name)",
		"ltrtrInfrmNm":       "학명 기재문",
		"plantBrdgFomTpcdNm": "식물 번식 구분 형태",
		"plantChnNm":         "식물 중국명",
		"plantEngNm":         "식물 영문명",
		"plantGnrlNm":        "식물 추천 국명",
		"plantGnrlNm2":       "식물 비추천 국명",
		"plantJpnNm":         "식물 일본명",
		"plantScnmId":        "식물 학명ID",
		"plantSpecsScnm":     "식물 학명",
		"rareTpcdNm":         "희귀식물 분류명",
		"relPlantSpecsScnm":  "연관 학명",
		"relScnmTpcdNm":      "연관 학명 구분명",
		"rmrk":               "비고",
		"rrnssPlantYn":       "희귀식물 여부",
		"spcltPlantCdNm":     "특산식물 분류명",
	}

	for _, tool := range result.Tools {
		if tool.Name != "kpni_scnm_info" {
			continue
		}
		schema, ok := tool.OutputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool kpni_scnm_info output schema = %#v", tool.OutputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool kpni_scnm_info output schema properties = %#v", schema["properties"])
		}
		checkKeys(t, properties, "item")
		checkPropertyDescriptions(t, "kpni_scnm_info", properties, map[string]string{"item": "상세 조회 결과"})

		itemProperty, ok := properties["item"].(map[string]any)
		if !ok {
			t.Fatalf("tool kpni_scnm_info item property = %#v", properties["item"])
		}
		itemProperties, nullable := nullableObjectProperties(itemProperty)
		if !nullable {
			t.Errorf("tool kpni_scnm_info item schema is not nullable: %#v", itemProperty)
		}
		if itemProperties == nil {
			t.Fatalf("tool kpni_scnm_info item schema has no object properties: %#v", itemProperty)
		}
		checkKeys(t, itemProperties, mapKeys(itemDescriptions)...)
		checkPropertyDescriptions(t, "kpni_scnm_info item", itemProperties, itemDescriptions)
		return
	}
	t.Fatal("tool kpni_scnm_info not found")
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
