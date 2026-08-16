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

type plantPictorialBookSearchUseCaseStub struct {
	query  application.PlantPictorialBookSearchQuery
	result application.PlantPictorialBookSearchResult
	err    error
}

func (s *plantPictorialBookSearchUseCaseStub) PlantPictorialBookSearch(_ context.Context, query application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantPictorialBookSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantPictorialBookSearchUseCaseStub{result: application.PlantPictorialBookSearchResult{
		Items:        []application.PlantPictorialBookSearchItem{{PlantGeneralName: "plant general name"}},
		NumberOfRows: 10,
		PageNumber:   2,
		TotalCount:   21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newServer(UseCases{PlantPictorialBookSearch: useCase}).Connect(ctx, serverTransport, nil)
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
		Name: plantResourcePlantPictorialBookSearchToolName,
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

	wantQuery := application.PlantPictorialBookSearchQuery{
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
	if !ok || item["plantGeneralName"] != "plant general name" {
		t.Errorf("item = %#v", items[0])
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      plantResourcePlantPictorialBookSearchToolName,
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
