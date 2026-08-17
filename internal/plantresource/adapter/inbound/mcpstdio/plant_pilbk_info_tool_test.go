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
	addPlantPilbkInfoTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_pilbk_info",
		map[string]string{
			"reqPlantPilbkNo": "검색할 식물도감번호 (plantPilbkSearch 결과의 plantPilbkNo)",
		},
		[]string{"reqPlantPilbkNo"},
	)

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_pilbk_info",
		Arguments: map[string]any{
			"reqPlantPilbkNo": "test-book-number",
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
		"apgFamilyKorNm": "apg family Korean name",
		"apgFamilyNm":    "apg family name",
		"bfofMthod":      "pest control method",
		"brdMthdDesc":    "breeding method description",
		"bugInfo":        "bug information",
		"dstrb":          "distribution",
		"engNm":          "English name",
		"familyKorNm":    "family Korean name",
		"familyNm":       "family name",
		"farmSpftDesc":   "farm feature description",
		"genusKorNm":     "genus Korean name",
		"genusNm":        "genus name",
		"grwEvrntDesc":   "growth environment description",
		"inductionDesc":  "induction description",
		"lastUpdtDtm":    "last update date time",
		"notRcmmGnrlNm":  "not recommended general name",
		"note":           "note",
		"orplcNm":        "origin place name",
		"osDstrb":        "overseas distribution",
		"plantGnrlNm":    "plant general name",
		"plantPilbkNo":   "plant pictorial book number",
		"plantSpecsScnm": "plant species scientific name",
		"prtcPlnDesc":    "protection plan description",
		"rrngGubun":      "growth classification",
		"rrngType":       "growth type",
		"shpe":           "shape",
		"smlrPlntDesc":   "similar plant description",
		"spft":           "feature",
		"useMthdDesc":    "use method description",
		"woodDesc":       "wood description",
	}
	if len(output) != len(wantOutput) {
		t.Errorf("output key count = %d, want %d", len(output), len(wantOutput))
	}
	for key, want := range wantOutput {
		if got := output[key]; got != want {
			t.Errorf("output %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_pilbk_info", mapKeys(wantOutput), nil, nil)

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_pilbk_info",
		Arguments: map[string]any{"reqPlantPilbkNo": "test-book-number"},
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
