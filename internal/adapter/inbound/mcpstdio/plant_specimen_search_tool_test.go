package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"kna-mcp/internal/application"
)

type plantSpecimenSearchUseCaseStub struct {
	query  application.PlantSpecimenSearchQuery
	result application.PlantSpecimenSearchResult
	err    error
}

func (s *plantSpecimenSearchUseCaseStub) PlantSpecimenSearch(_ context.Context, query application.PlantSpecimenSearchQuery) (application.PlantSpecimenSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSpecimenSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSpecimenSearchUseCaseStub{result: application.PlantSpecimenSearchResult{
		Items: []application.PlantSpecimenSearchItem{{
			Count:            436,
			PlantGeneralName: "리기다소나무",
			PlantSpeciesID:   "P000004951",
		}},
		NumberOfRows: 10,
		PageNumber:   2,
		TotalCount:   21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newServer(UseCases{PlantSpecimenSearch: useCase}).Connect(ctx, serverTransport, nil)
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
		Name: plantResourcePlantSpecimenSearchToolName,
		Arguments: map[string]any{
			"pageNumber":        2,
			"numberOfRows":      10,
			"requestSearchWord": "소나무",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantSpecimenSearchQuery{
		PageNumber:        2,
		NumberOfRows:      10,
		RequestSearchWord: "소나무",
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
	if !ok || item["count"] != float64(436) || item["plantSpeciesId"] != "P000004951" {
		t.Errorf("item = %#v", items[0])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      plantResourcePlantSpecimenSearchToolName,
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
