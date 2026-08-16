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

type plantPictorialBookInformationUseCaseStub struct {
	query  application.PlantPictorialBookInformationQuery
	result application.PlantPictorialBookInformationResult
	err    error
}

func (s *plantPictorialBookInformationUseCaseStub) PlantPictorialBookInformation(_ context.Context, query application.PlantPictorialBookInformationQuery) (application.PlantPictorialBookInformationResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantPictorialBookInformationTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantPictorialBookInformationUseCaseStub{result: application.PlantPictorialBookInformationResult{
		APGFamilyKoreanName:           "apg family Korean name",
		APGFamilyName:                 "apg family name",
		BfofMethod:                    "bfof method",
		BreedingMethodDescription:     "breeding method description",
		BugInformation:                "bug information",
		Distribution:                  "distribution",
		EnglishName:                   "English name",
		FamilyKoreanName:              "family Korean name",
		FamilyName:                    "family name",
		FarmSpecialFeatureDescription: "farm special feature description",
		GenusKoreanName:               "genus Korean name",
		GenusName:                     "genus name",
		GrowthEnvironmentDescription:  "growth environment description",
		InductionDescription:          "induction description",
		LastUpdateDateTime:            "last update date time",
		NotRecommendedGeneralName:     "not recommended general name",
		Note:                          "note",
		OriginPlaceName:               "origin place name",
		OverseasDistribution:          "overseas distribution",
		PlantGeneralName:              "plant general name",
		PlantPictorialBookNumber:      "plant pictorial book number",
		PlantSpeciesScientificName:    "plant species scientific name",
		ProtectionPlanDescription:     "protection plan description",
		RearingGubun:                  "rearing gubun",
		RearingType:                   "rearing type",
		Shape:                         "shape",
		SimilarPlantDescription:       "similar plant description",
		SpecialFeature:                "special feature",
		UseMethodDescription:          "use method description",
		WoodDescription:               "wood description",
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newServer(UseCases{PlantPictorialBookInformation: useCase}).Connect(ctx, serverTransport, nil)
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
		Name: plantResourcePlantPictorialBookInformationToolName,
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

	wantQuery := application.PlantPictorialBookInformationQuery{RequestPlantPictorialBookNumber: "test-book-number"}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	wantOutput := map[string]string{
		"apgFamilyKoreanName":           "apg family Korean name",
		"apgFamilyName":                 "apg family name",
		"bfofMethod":                    "bfof method",
		"breedingMethodDescription":     "breeding method description",
		"bugInformation":                "bug information",
		"distribution":                  "distribution",
		"englishName":                   "English name",
		"familyKoreanName":              "family Korean name",
		"familyName":                    "family name",
		"farmSpecialFeatureDescription": "farm special feature description",
		"genusKoreanName":               "genus Korean name",
		"genusName":                     "genus name",
		"growthEnvironmentDescription":  "growth environment description",
		"inductionDescription":          "induction description",
		"lastUpdateDateTime":            "last update date time",
		"notRecommendedGeneralName":     "not recommended general name",
		"note":                          "note",
		"originPlaceName":               "origin place name",
		"overseasDistribution":          "overseas distribution",
		"plantGeneralName":              "plant general name",
		"plantPictorialBookNumber":      "plant pictorial book number",
		"plantSpeciesScientificName":    "plant species scientific name",
		"protectionPlanDescription":     "protection plan description",
		"rearingGubun":                  "rearing gubun",
		"rearingType":                   "rearing type",
		"shape":                         "shape",
		"similarPlantDescription":       "similar plant description",
		"specialFeature":                "special feature",
		"useMethodDescription":          "use method description",
		"woodDescription":               "wood description",
	}
	for key, want := range wantOutput {
		if got := output[key]; got != want {
			t.Errorf("output %s = %#v, want %q", key, got, want)
		}
	}

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      plantResourcePlantPictorialBookInformationToolName,
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
