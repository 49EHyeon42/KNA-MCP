package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

type plantSampleSearchUseCaseStub struct {
	query  application.PlantSampleSearchQuery
	result application.PlantSampleSearchResult
	err    error
}

func (s *plantSampleSearchUseCaseStub) PlantSampleSearch(_ context.Context, query application.PlantSampleSearchQuery) (application.PlantSampleSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSampleSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSampleSearchUseCaseStub{result: application.PlantSampleSearchResult{
		Items: []application.PlantSampleSearchItem{{
			Count:            123,
			PlantGeneralName: "plant general name",
			PlantSpeciesID:   "plant species ID",
		}},
		NumberOfRows: 10,
		PageNumber:   2,
		TotalCount:   21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newServer(UseCases{PlantSampleSearch: useCase}).Connect(ctx, serverTransport, nil)
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

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: plantResourcePlantSampleSearchToolName,
		Arguments: map[string]any{
			"pageNumber":        2,
			"numberOfRows":      10,
			"requestSearchWord": "test-search-word",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantSampleSearchQuery{
		PageNumber:        2,
		NumberOfRows:      10,
		RequestSearchWord: "test-search-word",
	}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	if output["totalCount"] != float64(21) {
		t.Errorf("totalCount = %#v, want 21", output["totalCount"])
	}
	items, ok := output["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", output["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok || item["count"] != float64(123) || item["plantSpeciesId"] != "plant species ID" {
		t.Errorf("item = %#v", items[0])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      plantResourcePlantSampleSearchToolName,
		Arguments: map[string]any{"pageNumber": 1, "numberOfRows": 1},
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
