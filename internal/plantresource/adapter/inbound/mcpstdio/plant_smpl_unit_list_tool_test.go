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

type plantSmplUnitListUseCaseStub struct {
	query  application.PlantSmplUnitListQuery
	result application.PlantSmplUnitListResult
	err    error
}

func (s *plantSmplUnitListUseCaseStub) PlantSmplUnitList(_ context.Context, query application.PlantSmplUnitListQuery) (application.PlantSmplUnitListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSmplUnitListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSmplUnitListUseCaseStub{result: application.PlantSmplUnitListResult{
		Items: []application.PlantSmplUnitListItem{{
			AgpFamilyKorNm:     "agp family Korean name",
			AgpFamilyNm:        "agp family name",
			BspcsInsttNm:       "specimen holding institution",
			ClarHaslvVal:       "collection site elevation",
			ClarNm:             "collection site",
			CllcrNm:            "collector name",
			FamilyKorNm:        "family Korean name",
			FamilyNm:           "family name",
			HbttChrcrCont:      "habitat characteristics",
			HbttTpcdNm:         "habitat type",
			PlantBrdgFomTpcdNm: "plant reproductive form",
			PlantGnrlNm:        "plant general name",
			PlantPilbkNo:       "plant pictorial book number",
			PlantSmplNo:        "plant specimen number",
			PlantSpecsID:       "plant species ID",
			PlantSpecsScnm:     "plant species scientific name",
			SmplCllcnDt:        "specimen collection date",
			SmplClnyNm:         "specimen community name",
			SmplKindCdNm:       "specimen type",
			SmplWrdt:           "specimen preparation date",
			VgttnTpeCdNm:       "vegetation type",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantSmplUnitListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_smpl_unit_list",
		map[string]string{
			"pageNo":          "페이지번호 (1 이상)",
			"numOfRows":       "한 페이지 결과 수 (1 이상)",
			"reqPlantSpecsId": "검색할 식물표본의 식물종ID (plantSmplSearch 결과의 plantSpecsId)",
		},
		[]string{"pageNo", "numOfRows", "reqPlantSpecsId"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_smpl_unit_list",
		Arguments: map[string]any{
			"pageNo":          2,
			"numOfRows":       10,
			"reqPlantSpecsId": "test-plant-species-id",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantSmplUnitListQuery{
		PageNo:          2,
		NumOfRows:       10,
		ReqPlantSpecsID: "test-plant-species-id",
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
		"agpFamilyKorNm":     "agp family Korean name",
		"agpFamilyNm":        "agp family name",
		"bspcsInsttNm":       "specimen holding institution",
		"clarHaslvVal":       "collection site elevation",
		"clarNm":             "collection site",
		"cllcrNm":            "collector name",
		"familyKorNm":        "family Korean name",
		"familyNm":           "family name",
		"hbttChrcrCont":      "habitat characteristics",
		"hbttTpcdNm":         "habitat type",
		"plantBrdgFomTpcdNm": "plant reproductive form",
		"plantGnrlNm":        "plant general name",
		"plantPilbkNo":       "plant pictorial book number",
		"plantSmplNo":        "plant specimen number",
		"plantSpecsId":       "plant species ID",
		"plantSpecsScnm":     "plant species scientific name",
		"smplCllcnDt":        "specimen collection date",
		"smplClnyNm":         "specimen community name",
		"smplKindCdNm":       "specimen type",
		"smplWrdt":           "specimen preparation date",
		"vgttnTpeCdNm":       "vegetation type",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_smpl_unit_list",
		[]string{"items", "numOfRows", "pageNo", "totalCount"}, mapKeys(wantItem), nil)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_smpl_unit_list",
		Arguments: map[string]any{
			"pageNo":          1,
			"numOfRows":       1,
			"reqPlantSpecsId": "test-plant-species-id",
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
