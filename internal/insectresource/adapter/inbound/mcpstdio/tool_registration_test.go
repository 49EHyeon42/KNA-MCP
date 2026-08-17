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

func TestAddToolsRegistersAllInsectResourceTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, UseCases{
		InsectPilbkSearch:  &insectPilbkSearchUseCaseStub{},
		InsectPilbkInfo:    &insectPilbkInfoUseCaseStub{},
		InsectPrtctList:    &insectPrtctListUseCaseStub{},
		InsectSmplSearch:   &insectSmplSearchUseCaseStub{},
		InsectSmplUnitList: &insectSmplUnitListUseCaseStub{},
	}); err != nil {
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
	got := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		got = append(got, tool.Name)
		if !regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`).MatchString(tool.Name) {
			t.Errorf("tool name = %q, want 1 to 128 allowed characters", tool.Name)
		}
	}
	slices.Sort(got)
	want := []string{"insect_resource_insect_pilbk_info", "insect_resource_insect_pilbk_search", "insect_resource_insect_prtct_list", "insect_resource_insect_smpl_search", "insect_resource_insect_smpl_unit_list"}
	if !slices.Equal(got, want) {
		t.Fatalf("tools = %#v, want %#v", got, want)
	}
}

func TestAddToolsRequiresInsectPilbkSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{})
	if err == nil || err.Error() != "insectPilbkSearch use case is required" {
		t.Errorf("error = %v, want insectPilbkSearch use case is required", err)
	}
}

func TestAddToolsRequiresInsectPilbkInfoUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{InsectPilbkSearch: &insectPilbkSearchUseCaseStub{}})
	if err == nil || err.Error() != "insectPilbkInfo use case is required" {
		t.Errorf("error = %v, want insectPilbkInfo use case is required", err)
	}
}

func TestAddToolsRequiresInsectPrtctListUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		InsectPilbkSearch: &insectPilbkSearchUseCaseStub{},
		InsectPilbkInfo:   &insectPilbkInfoUseCaseStub{},
	})
	if err == nil || err.Error() != "insectPrtctList use case is required" {
		t.Errorf("error = %v, want insectPrtctList use case is required", err)
	}
}

func TestAddToolsRequiresInsectSmplSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		InsectPilbkSearch: &insectPilbkSearchUseCaseStub{},
		InsectPilbkInfo:   &insectPilbkInfoUseCaseStub{},
		InsectPrtctList:   &insectPrtctListUseCaseStub{},
	})
	if err == nil || err.Error() != "insectSmplSearch use case is required" {
		t.Errorf("error = %v, want insectSmplSearch use case is required", err)
	}
}

func TestAddToolsRequiresInsectSmplUnitListUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		InsectPilbkSearch: &insectPilbkSearchUseCaseStub{},
		InsectPilbkInfo:   &insectPilbkInfoUseCaseStub{},
		InsectPrtctList:   &insectPrtctListUseCaseStub{},
		InsectSmplSearch:  &insectSmplSearchUseCaseStub{},
	})
	if err == nil || err.Error() != "insectSmplUnitList use case is required" {
		t.Errorf("error = %v, want insectSmplUnitList use case is required", err)
	}
}
