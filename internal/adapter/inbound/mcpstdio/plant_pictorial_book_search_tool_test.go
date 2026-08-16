package mcpstdio

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"kna-mcp/internal/application"
)

type plantPictorialBookSearchUseCaseStub struct {
	query  application.PlantPictorialBookSearchQuery
	result application.PlantPictorialBookSearchResult
}

func (s *plantPictorialBookSearchUseCaseStub) PlantPictorialBookSearch(_ context.Context, query application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error) {
	s.query = query
	return s.result, nil
}

func TestPlantPictorialBookSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantPictorialBookSearchUseCaseStub{result: application.PlantPictorialBookSearchResult{
		Items:        []application.PlantPictorialBookSearchItem{{PlantGeneralName: "소나무"}},
		NumberOfRows: 10,
		PageNumber:   2,
		TotalCount:   21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newServer(useCase).Connect(ctx, serverTransport, nil)
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
			"requestSearchWord": "소나무",
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
	if !ok || item["plantGeneralName"] != "소나무" {
		t.Errorf("item = %#v", items[0])
	}
}
