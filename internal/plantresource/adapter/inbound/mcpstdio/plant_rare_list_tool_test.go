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

type plantRareListUseCaseStub struct {
	query  application.PlantRareListQuery
	result application.PlantRareListResult
	err    error
}

func (s *plantRareListUseCaseStub) PlantRareList(_ context.Context, query application.PlantRareListQuery) (application.PlantRareListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantRareListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantRareListUseCaseStub{result: application.PlantRareListResult{
		Items: []application.PlantRareListItem{{
			AgpFamilyNm:      "agp family name",
			APGFamilyKorNm:   "apg family Korean name",
			ExtrmCrssScls1Yn: "endangered class one yes or no",
			ExtrmCrssScls2Yn: "endangered class two yes or no",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			PlantGnrlNm:      "plant general name",
			PlantSpecsScnm:   "plant species scientific name",
			RareTpcdNm:       "rare type code name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantRareListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_rare_list",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 식물의 학명 또는 국명",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "plant_resource_plant_rare_list", "산림청 국립수목원 적색식물 목록을 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_rare_list",
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

	wantQuery := application.PlantRareListQuery{
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
		"agpFamilyNm":      "agp family name",
		"apgFamilyKorNm":   "apg family Korean name",
		"extrmCrssScls1Yn": "endangered class one yes or no",
		"extrmCrssScls2Yn": "endangered class two yes or no",
		"familyKorNm":      "family Korean name",
		"familyNm":         "family name",
		"plantGnrlNm":      "plant general name",
		"plantSpecsScnm":   "plant species scientific name",
		"rareTpcdNm":       "rare type code name",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_rare_list",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지 결과 수",
			"pageNo":     "페이지번호",
			"totalCount": "전체 결과 수",
		}, map[string]string{
			"agpFamilyNm":      "APG과명",
			"apgFamilyKorNm":   "APG과국명",
			"extrmCrssScls1Yn": "멸종위기종1급 여부",
			"extrmCrssScls2Yn": "멸종위기종2급 여부",
			"familyKorNm":      "과국명",
			"familyNm":         "과명",
			"plantGnrlNm":      "국명(식물명)",
			"plantSpecsScnm":   "학명",
			"rareTpcdNm":       "IUCN 적색식물 등급",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_rare_list",
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
