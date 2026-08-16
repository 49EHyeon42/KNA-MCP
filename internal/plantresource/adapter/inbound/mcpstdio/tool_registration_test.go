package mcpstdio

import (
	"context"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

func TestAddToolsRegistersAllPlantResourceTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, completePlantResourceUseCases()); err != nil {
		t.Fatal(err)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
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

	result, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	got := make([]string, len(result.Tools))
	validName := regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`)
	for i, tool := range result.Tools {
		got[i] = tool.Name
		if !validName.MatchString(tool.Name) {
			t.Errorf("tool name = %q, want 1 to 128 allowed characters", tool.Name)
		}
	}
	slices.Sort(got)

	want := []string{
		"plant_resource_plant_folk_area_list",
		"plant_resource_plant_folk_search",
		"plant_resource_plant_naturalized_list",
		"plant_resource_plant_pilbk_info",
		"plant_resource_plant_pilbk_search",
		"plant_resource_plant_rare_list",
		"plant_resource_plant_seed_grmnt_list",
		"plant_resource_plant_seed_search",
		"plant_resource_plant_seed_unit_list",
		"plant_resource_plant_smpl_search",
		"plant_resource_plant_smpl_unit_list",
		"plant_resource_plant_spclt_list",
		"plant_resource_plant_word_list",
	}
	if !slices.Equal(got, want) {
		t.Errorf("tool names = %#v, want %#v", got, want)
	}
}

func TestAddToolsRequiresEveryPlantResourceUseCase(t *testing.T) {
	for _, test := range []struct {
		name      string
		remove    func(*UseCases)
		wantError string
	}{
		{name: "plantPilbkSearch", remove: func(useCases *UseCases) { useCases.PlantPilbkSearch = nil }, wantError: "plantPilbkSearch use case is required"},
		{name: "plantPilbkInfo", remove: func(useCases *UseCases) { useCases.PlantPilbkInfo = nil }, wantError: "plantPilbkInfo use case is required"},
		{name: "plantSmplSearch", remove: func(useCases *UseCases) { useCases.PlantSmplSearch = nil }, wantError: "plantSmplSearch use case is required"},
		{name: "plantSmplUnitList", remove: func(useCases *UseCases) { useCases.PlantSmplUnitList = nil }, wantError: "plantSmplUnitList use case is required"},
		{name: "plantSeedSearch", remove: func(useCases *UseCases) { useCases.PlantSeedSearch = nil }, wantError: "plantSeedSearch use case is required"},
		{name: "plantSeedUnitList", remove: func(useCases *UseCases) { useCases.PlantSeedUnitList = nil }, wantError: "plantSeedUnitList use case is required"},
		{name: "plantSeedGrmntList", remove: func(useCases *UseCases) { useCases.PlantSeedGrmntList = nil }, wantError: "plantSeedGrmntList use case is required"},
		{name: "plantFolkSearch", remove: func(useCases *UseCases) { useCases.PlantFolkSearch = nil }, wantError: "plantFolkSearch use case is required"},
		{name: "plantFolkAreaList", remove: func(useCases *UseCases) { useCases.PlantFolkAreaList = nil }, wantError: "plantFolkAreaList use case is required"},
		{name: "plantNaturalizedList", remove: func(useCases *UseCases) { useCases.PlantNaturalizedList = nil }, wantError: "plantNaturalizedList use case is required"},
		{name: "plantRareList", remove: func(useCases *UseCases) { useCases.PlantRareList = nil }, wantError: "plantRareList use case is required"},
		{name: "plantSpcltList", remove: func(useCases *UseCases) { useCases.PlantSpcltList = nil }, wantError: "plantSpcltList use case is required"},
		{name: "plantWordList", remove: func(useCases *UseCases) { useCases.PlantWordList = nil }, wantError: "plantWordList use case is required"},
	} {
		t.Run(test.name, func(t *testing.T) {
			useCases := completePlantResourceUseCases()
			test.remove(&useCases)

			err := AddTools(mcpserver.NewServer(), useCases)
			if err == nil || err.Error() != test.wantError {
				t.Errorf("error = %v, want %q", err, test.wantError)
			}
		})
	}
}

func completePlantResourceUseCases() UseCases {
	return UseCases{
		PlantPilbkSearch:     &plantPilbkSearchUseCaseStub{},
		PlantPilbkInfo:       &plantPilbkInfoUseCaseStub{},
		PlantSmplSearch:      &plantSmplSearchUseCaseStub{},
		PlantSmplUnitList:    &plantSmplUnitListUseCaseStub{},
		PlantSeedSearch:      &plantSeedSearchUseCaseStub{},
		PlantSeedUnitList:    &plantSeedUnitListUseCaseStub{},
		PlantSeedGrmntList:   &plantSeedGrmntListUseCaseStub{},
		PlantFolkSearch:      &plantFolkSearchUseCaseStub{},
		PlantFolkAreaList:    &plantFolkAreaListUseCaseStub{},
		PlantNaturalizedList: &plantNaturalizedListUseCaseStub{},
		PlantRareList:        &plantRareListUseCaseStub{},
		PlantSpcltList:       &plantSpcltListUseCaseStub{},
		PlantWordList:        &plantWordListUseCaseStub{},
	}
}
