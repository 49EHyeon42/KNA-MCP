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

type plantSmplSearchUseCaseStub struct {
	query  application.PlantSmplSearchQuery
	result application.PlantSmplSearchResult
	err    error
}

func (s *plantSmplSearchUseCaseStub) PlantSmplSearch(_ context.Context, query application.PlantSmplSearchQuery) (application.PlantSmplSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSmplSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSmplSearchUseCaseStub{result: application.PlantSmplSearchResult{
		Items: []application.PlantSmplSearchItem{{
			Cnt:            123,
			FamilyKorNm:    "family Korean name",
			FamilyNm:       "family name",
			PlantGnrlNm:    "plant general name",
			PlantSpecsID:   "plant species ID",
			PlantSpecsScnm: "plant species scientific name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantSmplSearchTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_smpl_search",
		[]string{"pageNo", "numOfRows", "reqSearchWrd"},
		[]string{"pageNo", "numOfRows"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_smpl_search",
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

	wantQuery := application.PlantSmplSearchQuery{
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
	wantItem := map[string]any{
		"cnt":            float64(123),
		"familyKorNm":    "family Korean name",
		"familyNm":       "family name",
		"plantGnrlNm":    "plant general name",
		"plantSpecsId":   "plant species ID",
		"plantSpecsScnm": "plant species scientific name",
	}
	if !reflect.DeepEqual(item, wantItem) {
		t.Errorf("item = %#v, want %#v", item, wantItem)
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_smpl_search",
		[]string{"items", "numOfRows", "pageNo", "totalCount"}, mapKeys(wantItem))

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_smpl_search",
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
