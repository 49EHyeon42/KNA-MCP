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

type plantNaturalizedListUseCaseStub struct {
	query  application.PlantNaturalizedListQuery
	result application.PlantNaturalizedListResult
	err    error
}

func (s *plantNaturalizedListUseCaseStub) PlantNaturalizedList(_ context.Context, query application.PlantNaturalizedListQuery) (application.PlantNaturalizedListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantNaturalizedListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantNaturalizedListUseCaseStub{result: application.PlantNaturalizedListResult{
		Items: []application.PlantNaturalizedListItem{{
			AgpFamilyNm:        "agp family name",
			APGFamilyKorNm:     "apg family Korean name",
			BlprdEnmnt:         "bloom period end month",
			BlprdStmnt:         "bloom period start month",
			DistrAraDscrt:      "distribution area description",
			EclgDstrbYn:        "ecological disturbance yes or no",
			ExtcPlantCdNm:      "exotic plant code name",
			FamilyKorNm:        "family Korean name",
			FamilyNm:           "family name",
			FrtTpcdNm:          "fruit type code name",
			LastUpdtDtm:        "last update date time",
			NtldgTpcdNm:        "naturalization degree type name",
			NtrlzEraTpcdNm:     "naturalization era type name",
			OrplcNm:            "original place name",
			PlantBrdgFomTpcdNm: "plant breeding form type name",
			PlantDistrGrcd:     "plant distribution grade code",
			PlantDistrQntt:     "plant distribution quantity",
			PlantDistrQnttGrcd: "plant distribution quantity grade code",
			PlantEngNm:         "plant English name",
			PlantGnrlNm:        "plant general name",
			PlantJpnNm:         "plant Japanese name",
			PlantLfcclTpcdNm:   "plant life cycle type name",
			PlantSpecsScnm:     "plant species scientific name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantNaturalizedListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_naturalized_list",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 식물의 국명 또는 학명",
		},
		[]string{"pageNo", "numOfRows"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_naturalized_list",
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

	wantQuery := application.PlantNaturalizedListQuery{
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
		"agpFamilyNm":        "agp family name",
		"apgFamilyKorNm":     "apg family Korean name",
		"blprdEnmnt":         "bloom period end month",
		"blprdStmnt":         "bloom period start month",
		"distrAraDscrt":      "distribution area description",
		"eclgDstrbYn":        "ecological disturbance yes or no",
		"extcPlantCdNm":      "exotic plant code name",
		"familyKorNm":        "family Korean name",
		"familyNm":           "family name",
		"frtTpcdNm":          "fruit type code name",
		"lastUpdtDtm":        "last update date time",
		"ntldgTpcdNm":        "naturalization degree type name",
		"ntrlzEraTpcdNm":     "naturalization era type name",
		"orplcNm":            "original place name",
		"plantBrdgFomTpcdNm": "plant breeding form type name",
		"plantDistrGrcd":     "plant distribution grade code",
		"plantDistrQntt":     "plant distribution quantity",
		"plantDistrQnttGrcd": "plant distribution quantity grade code",
		"plantEngNm":         "plant English name",
		"plantGnrlNm":        "plant general name",
		"plantJpnNm":         "plant Japanese name",
		"plantLfcclTpcdNm":   "plant life cycle type name",
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
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_naturalized_list",
		[]string{"items", "numOfRows", "pageNo", "totalCount"}, mapKeys(wantItem),
		map[string]string{"agpFamilyNm": "APG과명", "apgFamilyKorNm": "APG과국명"})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_naturalized_list",
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
