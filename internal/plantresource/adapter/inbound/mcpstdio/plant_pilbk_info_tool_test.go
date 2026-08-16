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

type plantPilbkInfoUseCaseStub struct {
	query  application.PlantPilbkInfoQuery
	result application.PlantPilbkInfoResult
	err    error
}

func (s *plantPilbkInfoUseCaseStub) PlantPilbkInfo(_ context.Context, query application.PlantPilbkInfoQuery) (application.PlantPilbkInfoResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantPilbkInfoTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantPilbkInfoUseCaseStub{result: application.PlantPilbkInfoResult{
		APGFamilyKorNm: "apg family Korean name",
		APGFamilyNm:    "apg family name",
		BfofMthod:      "pest control method",
		BrdMthdDesc:    "breeding method description",
		BugInfo:        "bug information",
		Dstrb:          "distribution",
		EngNm:          "English name",
		FamilyKorNm:    "family Korean name",
		FamilyNm:       "family name",
		FarmSpftDesc:   "farm feature description",
		GenusKorNm:     "genus Korean name",
		GenusNm:        "genus name",
		GrwEvrntDesc:   "growth environment description",
		InductionDesc:  "induction description",
		LastUpdtDtm:    "last update date time",
		NotRcmmGnrlNm:  "not recommended general name",
		Note:           "note",
		OrplcNm:        "origin place name",
		OsDstrb:        "overseas distribution",
		PlantGnrlNm:    "plant general name",
		PlantPilbkNo:   "plant pictorial book number",
		PlantSpecsScnm: "plant species scientific name",
		PrtcPlnDesc:    "protection plan description",
		RrngGubun:      "growth classification",
		RrngType:       "growth type",
		Shpe:           "shape",
		SmlrPlntDesc:   "similar plant description",
		Spft:           "feature",
		UseMthdDesc:    "use method description",
		WoodDesc:       "wood description",
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	AddTools(server, UseCases{PlantPilbkInfo: useCase})
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
		Name: plantResourcePlantPilbkInfoToolName,
		Arguments: map[string]any{
			"requestPlantPictorialBookNumber": "test-book-number",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "test-book-number"}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	wantOutput := map[string]string{
		"apgFamilyKoreanName":          "apg family Korean name",
		"apgFamilyName":                "apg family name",
		"pestControlMethod":            "pest control method",
		"breedingMethodDescription":    "breeding method description",
		"bugInformation":               "bug information",
		"distribution":                 "distribution",
		"englishName":                  "English name",
		"familyKoreanName":             "family Korean name",
		"familyName":                   "family name",
		"farmFeatureDescription":       "farm feature description",
		"genusKoreanName":              "genus Korean name",
		"genusName":                    "genus name",
		"growthEnvironmentDescription": "growth environment description",
		"inductionDescription":         "induction description",
		"lastUpdateDateTime":           "last update date time",
		"notRecommendedGeneralName":    "not recommended general name",
		"note":                         "note",
		"originPlaceName":              "origin place name",
		"overseasDistribution":         "overseas distribution",
		"plantGeneralName":             "plant general name",
		"plantPictorialBookNumber":     "plant pictorial book number",
		"plantSpeciesScientificName":   "plant species scientific name",
		"protectionPlanDescription":    "protection plan description",
		"growthClassification":         "growth classification",
		"growthType":                   "growth type",
		"shape":                        "shape",
		"similarPlantDescription":      "similar plant description",
		"feature":                      "feature",
		"useMethodDescription":         "use method description",
		"woodDescription":              "wood description",
	}
	for key, want := range wantOutput {
		if got := output[key]; got != want {
			t.Errorf("output %s = %#v, want %q", key, got, want)
		}
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      plantResourcePlantPilbkInfoToolName,
		Arguments: map[string]any{"requestPlantPictorialBookNumber": "test-book-number"},
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
