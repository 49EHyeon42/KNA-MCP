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

type plantSeedUnitListUseCaseStub struct {
	query  application.PlantSeedUnitListQuery
	result application.PlantSeedUnitListResult
	err    error
}

func (s *plantSeedUnitListUseCaseStub) PlantSeedUnitList(_ context.Context, query application.PlantSeedUnitListQuery) (application.PlantSeedUnitListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSeedUnitListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSeedUnitListUseCaseStub{result: application.PlantSeedUnitListResult{
		Items: []application.PlantSeedUnitListItem{{
			CllcnDate:        "collection date",
			PlantGnrlNm:      "plant general name",
			QualtFllnsRt:     "quality fullness rate",
			SdwghWeght:       "thousand seed weight",
			SeedAdmcn:        "seed air-dry moisture content",
			SeedCllctPlace:   "seed collection place",
			SeedHoldGrainCnt: "seed holding grain count",
			SeedHoldQntt:     "seed holding quantity",
			SeedNo:           "seed number",
			SeedSpecsID:      "seed species ID",
			StoreChrcrTpcdNm: "storage characteristic type",
			Vtlfct:           "vitality rate",
			VtlfctTestYr:     "vitality test year",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	AddTools(server, UseCases{PlantSeedUnitList: useCase})
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_seed_unit_list",
		[]string{"pageNo", "numOfRows", "reqSeedSpecsId"},
		[]string{"pageNo", "numOfRows", "reqSeedSpecsId"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_seed_unit_list",
		Arguments: map[string]any{
			"pageNo":         2,
			"numOfRows":      10,
			"reqSeedSpecsId": "test-seed-species-id",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantSeedUnitListQuery{
		PageNo:         2,
		NumOfRows:      10,
		ReqSeedSpecsID: "test-seed-species-id",
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
		"cllcnDate":        "collection date",
		"plantGnrlNm":      "plant general name",
		"qualtFllnsRt":     "quality fullness rate",
		"sdwghWeght":       "thousand seed weight",
		"seedAdmcn":        "seed air-dry moisture content",
		"seedCllctPlace":   "seed collection place",
		"seedHoldGrainCnt": "seed holding grain count",
		"seedHoldQntt":     "seed holding quantity",
		"seedNo":           "seed number",
		"seedSpecsId":      "seed species ID",
		"storeChrcrTpcdNm": "storage characteristic type",
		"vtlfct":           "vitality rate",
		"vtlfctTestYr":     "vitality test year",
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
		Name: "plant_resource_plant_seed_unit_list",
		Arguments: map[string]any{
			"pageNo":         1,
			"numOfRows":      1,
			"reqSeedSpecsId": "test-seed-species-id",
		},
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
