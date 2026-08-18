package main

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

func TestCompositionRegistersAllTools(t *testing.T) {
	server := mcpserver.NewServer()
	if err := addPlantResourceTools(server, "test-key"); err != nil {
		t.Fatal(err)
	}
	if err := addPlantMstnsTools(server, "test-key"); err != nil {
		t.Fatal(err)
	}
	if err := addKpniTools(server, "test-key"); err != nil {
		t.Fatal(err)
	}
	if err := addInsectResourceTools(server, "test-key"); err != nil {
		t.Fatal(err)
	}
	if err := addFungiResourceTools(server, "test-key"); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	for i, tool := range result.Tools {
		got[i] = tool.Name
	}
	slices.Sort(got)

	want := []string{
		"fungi_resource_fngs_pilbk_info",
		"fungi_resource_fngs_pilbk_search",
		"fungi_resource_fngs_smpl_search",
		"fungi_resource_fngs_smpl_unit_list",
		"insect_resource_insect_pilbk_info",
		"insect_resource_insect_pilbk_search",
		"insect_resource_insect_prtct_list",
		"insect_resource_insect_smpl_search",
		"insect_resource_insect_smpl_unit_list",
		"kpni_gnrl_nm_ltrtr_search",
		"kpni_scnm_info",
		"kpni_scnm_search",
		"plant_mstns_plant_mstns_list",
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
