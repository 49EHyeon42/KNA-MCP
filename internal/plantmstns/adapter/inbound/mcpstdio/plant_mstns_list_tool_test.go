package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/plantmstns/application"
)

type plantMstnsListUseCaseStub struct {
	query  application.PlantMstnsListQuery
	result application.PlantMstnsListResult
	err    error
}

func (s *plantMstnsListUseCaseStub) PlantMstnsList(_ context.Context, query application.PlantMstnsListQuery) (application.PlantMstnsListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantMstnsListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantMstnsListUseCaseStub{result: application.PlantMstnsListResult{
		Items: []application.PlantMstnsListItem{{
			DistrAraDscrt:         "distribution area description",
			MinitrTpcdNm:          "miniature type code name",
			PlantBrdgFomTpcdNm:    "plant breeding form type code name",
			PlantGnrlNm:           "plant general name",
			PlantMinitrAthrNm:     "plant miniature author name",
			PlantMinitrMnfctMonth: "plant miniature manufacture month",
			PlantMinitrMnfctYr:    "plant miniature manufacture year",
			PlantMinitrPsinsNm:    "plant miniature possession name",
			PlantSpecsScnm:        "plant species scientific name",
			RrnssPlantYn:          "rarity plant yes or no",
			SpcltPlantYn:          "specialty plant yes or no",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantMstnsListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_mstns_plant_mstns_list",
		[]string{"pageNo", "numOfRows", "reqSearchWrd", "reqMnfctYr"},
		[]string{"pageNo", "numOfRows"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_mstns_plant_mstns_list",
		Arguments: map[string]any{
			"pageNo":       2,
			"numOfRows":    10,
			"reqSearchWrd": "test-search-word",
			"reqMnfctYr":   "2014",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantMstnsListQuery{
		PageNo:       2,
		NumOfRows:    10,
		ReqSearchWrd: "test-search-word",
		ReqMnfctYr:   "2014",
	}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	checkKeys(t, output, "items", "numOfRows", "pageNo", "totalCount")
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(21) {
		t.Errorf("pagination = %#v", output)
	}
	items, ok := output["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", output["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item = %#v", items[0])
	}
	wantItem := map[string]string{
		"distrAraDscrt":         "distribution area description",
		"minitrTpcdNm":          "miniature type code name",
		"plantBrdgFomTpcdNm":    "plant breeding form type code name",
		"plantGnrlNm":           "plant general name",
		"plantMinitrAthrNm":     "plant miniature author name",
		"plantMinitrMnfctMonth": "plant miniature manufacture month",
		"plantMinitrMnfctYr":    "plant miniature manufacture year",
		"plantMinitrPsinsNm":    "plant miniature possession name",
		"plantSpecsScnm":        "plant species scientific name",
		"rrnssPlantYn":          "rarity plant yes or no",
		"spcltPlantYn":          "specialty plant yes or no",
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
		Name:      "plant_mstns_plant_mstns_list",
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

func checkKeys(t *testing.T, got map[string]any, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("keys = %#v, want %#v", reflect.ValueOf(got).MapKeys(), want)
	}
	for _, key := range want {
		if _, ok := got[key]; !ok {
			t.Errorf("missing key %q in %#v", key, got)
		}
	}
}

func checkToolInputSchema(t *testing.T, ctx context.Context, session *mcp.ClientSession, toolName string, wantProperties, wantRequired []string) {
	t.Helper()
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range result.Tools {
		if tool.Name != toolName {
			continue
		}
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool %s input schema = %#v", toolName, tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("tool %s input schema properties = %#v", toolName, schema["properties"])
		}
		checkKeys(t, properties, wantProperties...)

		required, ok := schema["required"].([]any)
		if !ok {
			t.Fatalf("tool %s input schema required = %#v", toolName, schema["required"])
		}
		requiredKeys := make(map[string]any, len(required))
		for _, key := range required {
			name, ok := key.(string)
			if !ok {
				t.Fatalf("tool %s input schema required key = %#v", toolName, key)
			}
			requiredKeys[name] = nil
		}
		checkKeys(t, requiredKeys, wantRequired...)
		return
	}
	t.Fatalf("tool %s not found", toolName)
}
