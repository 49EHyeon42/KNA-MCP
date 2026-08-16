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
	AddTools(server, UseCases{PlantSmplUnitList: useCase})
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

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: plantResourcePlantSmplUnitListToolName,
		Arguments: map[string]any{
			"pageNumber":            2,
			"numberOfRows":          10,
			"requestPlantSpeciesId": "test-plant-species-id",
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
	if output["totalCount"] != float64(21) {
		t.Errorf("totalCount = %#v, want 21", output["totalCount"])
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
		"apgFamilyKoreanName":            "agp family Korean name",
		"apgFamilyName":                  "agp family name",
		"specimenHoldingInstitutionName": "specimen holding institution",
		"collectionSiteElevation":        "collection site elevation",
		"collectionSite":                 "collection site",
		"collectorName":                  "collector name",
		"familyKoreanName":               "family Korean name",
		"familyName":                     "family name",
		"habitatCharacteristics":         "habitat characteristics",
		"habitatTypeName":                "habitat type",
		"plantReproductiveForm":          "plant reproductive form",
		"plantGeneralName":               "plant general name",
		"plantPictorialBookNumber":       "plant pictorial book number",
		"plantSpecimenNumber":            "plant specimen number",
		"plantSpeciesId":                 "plant species ID",
		"plantSpeciesScientificName":     "plant species scientific name",
		"specimenCollectionDate":         "specimen collection date",
		"specimenCommunityName":          "specimen community name",
		"specimenTypeName":               "specimen type",
		"specimenPreparationDate":        "specimen preparation date",
		"vegetationTypeName":             "vegetation type",
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: plantResourcePlantSmplUnitListToolName,
		Arguments: map[string]any{
			"pageNumber":            1,
			"numberOfRows":          1,
			"requestPlantSpeciesId": "test-plant-species-id",
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
