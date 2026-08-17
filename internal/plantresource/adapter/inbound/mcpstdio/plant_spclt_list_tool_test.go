package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantSpcltListUseCaseStub struct {
	query  application.PlantSpcltListQuery
	result application.PlantSpcltListResult
	err    error
}

func (s *plantSpcltListUseCaseStub) PlantSpcltList(_ context.Context, query application.PlantSpcltListQuery) (application.PlantSpcltListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSpcltListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSpcltListUseCaseStub{result: application.PlantSpcltListResult{
		Items: []application.PlantSpcltListItem{{
			AgpFamilyKorNm:     "agp family Korean name",
			AgpFamilyNm:        "agp family name",
			ExtrmCrssScls1Yn:   "endangered class one yes or no",
			ExtrmCrssScls2Yn:   "endangered class two yes or no",
			FamilyKorNm:        "family Korean name",
			FamilyNm:           "family name",
			PlantBrdgFomTpcdNm: "plant breeding form type code name",
			PlantGnrlNm:        "plant general name",
			PlantSpecsScnm:     "plant species scientific name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantSpcltListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_spclt_list",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 특산식물 학명 또는 국명",
		},
		[]string{"pageNo", "numOfRows"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_spclt_list",
		Arguments: map[string]any{
			"pageNo":       2,
			"numOfRows":    10,
			"reqSearchWrd": "test-search-word",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantSpcltListQuery{
		PageNo:       2,
		NumOfRows:    10,
		ReqSearchWrd: "test-search-word",
	}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(21) {
		t.Errorf("pagination = %#v", output)
	}
	checkKeys(t, output, "items", "numOfRows", "pageNo", "totalCount")
	items, ok := output["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", output["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item = %#v", items[0])
	}
	wantItem := map[string]string{
		"agpFamilyKorNm":     "agp family Korean name",
		"agpFamilyNm":        "agp family name",
		"extrmCrssScls1Yn":   "endangered class one yes or no",
		"extrmCrssScls2Yn":   "endangered class two yes or no",
		"familyKorNm":        "family Korean name",
		"familyNm":           "family name",
		"plantBrdgFomTpcdNm": "plant breeding form type code name",
		"plantGnrlNm":        "plant general name",
		"plantSpecsScnm":     "plant species scientific name",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_spclt_list",
		[]string{"items", "numOfRows", "pageNo", "totalCount"}, mapKeys(wantItem), nil)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_spclt_list",
		Arguments: map[string]any{"pageNo": 1, "numOfRows": 1},
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
