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

type plantSeedGrmntListUseCaseStub struct {
	query  application.PlantSeedGrmntListQuery
	result application.PlantSeedGrmntListResult
	err    error
}

func (s *plantSeedGrmntListUseCaseStub) PlantSeedGrmntList(_ context.Context, query application.PlantSeedGrmntListQuery) (application.PlantSeedGrmntListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSeedGrmntListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSeedGrmntListUseCaseStub{result: application.PlantSeedGrmntListResult{
		Items: []application.PlantSeedGrmntListItem{{
			AvrgGrmntDcnt:     "average germination day count",
			GrmntBfrPrcesCont: "germination before processing content",
			GrmntClmdmCont:    "germination culture medium content",
			GrmntDscrt:        "germination description",
			GrmntExprmNo:      "germination experiment number",
			GrmntExprmSeq:     "germination experiment sequence",
			GrmntLightCndtn:   "germination light condition",
			GrmntRt:           "germination rate",
			GrmntTmpCndtn:     "germination temperature condition",
			PlantGnrlNm:       "plant general name",
			SeedNo:            "seed number",
			SeedSpecsID:       "seed species ID",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantSeedGrmntListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_seed_grmnt_list",
		map[string]string{
			"pageNo":         "페이지번호 (1 이상)",
			"numOfRows":      "한 페이지 결과 수 (1 이상)",
			"reqSeedSpecsId": "검색할 종자종ID (plantSeedSearch 결과의 seedSpecsId)",
		},
		[]string{"pageNo", "numOfRows", "reqSeedSpecsId"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_seed_grmnt_list",
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

	wantQuery := application.PlantSeedGrmntListQuery{
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
		"avrgGrmntDcnt":     "average germination day count",
		"grmntBfrPrcesCont": "germination before processing content",
		"grmntClmdmCont":    "germination culture medium content",
		"grmntDscrt":        "germination description",
		"grmntExprmNo":      "germination experiment number",
		"grmntExprmSeq":     "germination experiment sequence",
		"grmntLightCndtn":   "germination light condition",
		"grmntRt":           "germination rate",
		"grmntTmpCndtn":     "germination temperature condition",
		"plantGnrlNm":       "plant general name",
		"seedNo":            "seed number",
		"seedSpecsId":       "seed species ID",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_seed_grmnt_list",
		[]string{"items", "numOfRows", "pageNo", "totalCount"}, mapKeys(wantItem), nil)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_seed_grmnt_list",
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
