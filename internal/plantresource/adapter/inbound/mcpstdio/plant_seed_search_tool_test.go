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

type plantSeedSearchUseCaseStub struct {
	query  application.PlantSeedSearchQuery
	result application.PlantSeedSearchResult
	err    error
}

func (s *plantSeedSearchUseCaseStub) PlantSeedSearch(_ context.Context, query application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSeedSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSeedSearchUseCaseStub{result: application.PlantSeedSearchResult{
		Items: []application.PlantSeedSearchItem{{
			APGFamilyKorNm:   "apg family Korean name",
			APGFamilyNm:      "apg family name",
			BlprdEnmnt:       "bloom period end month",
			BlprdStmnt:       "bloom period start month",
			ClrngMthodCdNm:   "cleaning method",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			FritCdNm:         "fruit type",
			FrssnEnmnt:       "fruit season end month",
			FrssnStmnt:       "fruit season start month",
			LastUpdtDtm:      "last update date time",
			PlantGnrlNm:      "plant general name",
			PlantSpecsScnm:   "plant species scientific name",
			RfrncLtrtrCont:   "reference literature",
			SeedCtsrfcDesc:   "seed surface description",
			SeedCtsrfcTpcdNm: "seed surface type",
			SeedEmbrTpcdNm:   "seed embryo type",
			SeedMnmmBrdth:    "seed minimum breadth",
			SeedMnmmLngth:    "seed minimum length",
			SeedMxmmBrdth:    "seed maximum breadth",
			SeedMxmmLngth:    "seed maximum length",
			SeedShpDesc:      "seed shape description",
			SeedShpTpcdNm:    "seed shape type",
			SeedSpecsID:      "seed species ID",
			SeedTpcdNm:       "seed type",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	AddTools(server, UseCases{PlantSeedSearch: useCase})
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_seed_search",
		[]string{"pageNo", "numOfRows", "reqSearchWrd"},
		[]string{"pageNo", "numOfRows"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_seed_search",
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

	wantQuery := application.PlantSeedSearchQuery{
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
		"apgFamilyKorNm":   "apg family Korean name",
		"apgFamilyNm":      "apg family name",
		"blprdEnmnt":       "bloom period end month",
		"blprdStmnt":       "bloom period start month",
		"clrngMthodCdNm":   "cleaning method",
		"familyKorNm":      "family Korean name",
		"familyNm":         "family name",
		"fritCdNm":         "fruit type",
		"frssnEnmnt":       "fruit season end month",
		"frssnStmnt":       "fruit season start month",
		"lastUpdtDtm":      "last update date time",
		"plantGnrlNm":      "plant general name",
		"plantSpecsScnm":   "plant species scientific name",
		"rfrncLtrtrCont":   "reference literature",
		"seedCtsrfcDesc":   "seed surface description",
		"seedCtsrfcTpcdNm": "seed surface type",
		"seedEmbrTpcdNm":   "seed embryo type",
		"seedMnmmBrdth":    "seed minimum breadth",
		"seedMnmmLngth":    "seed minimum length",
		"seedMxmmBrdth":    "seed maximum breadth",
		"seedMxmmLngth":    "seed maximum length",
		"seedShpDesc":      "seed shape description",
		"seedShpTpcdNm":    "seed shape type",
		"seedSpecsId":      "seed species ID",
		"seedTpcdNm":       "seed type",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_seed_search",
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
