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

func TestAddToolsRegistersAllLchnServiceTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, UseCases{
		AlchnIlstrSearch: &alchnIlstrSearchUseCaseStub{},
		AlchnIlstrInfo:   &alchnIlstrInfoUseCaseStub{},
		AlchnSpcmSearch:  &alchnSpcmSearchUseCaseStub{},
		AlchnSpcmInfo:    &alchnSpcmInfoUseCaseStub{},
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
	names := make([]string, len(result.Tools))
	for i, tool := range result.Tools {
		names[i] = tool.Name
		if !regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`).MatchString(tool.Name) {
			t.Errorf("tool name = %q, want 1 to 128 allowed characters", tool.Name)
		}
	}
	slices.Sort(names)
	want := []string{"lchn_service_alchn_ilstr_info", "lchn_service_alchn_ilstr_search", "lchn_service_alchn_spcm_info", "lchn_service_alchn_spcm_search"}
	if !slices.Equal(names, want) {
		t.Fatalf("tools = %#v, want %#v", names, want)
	}
}

func TestAddToolsRequiresAlchnIlstrSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		AlchnIlstrInfo:  &alchnIlstrInfoUseCaseStub{},
		AlchnSpcmSearch: &alchnSpcmSearchUseCaseStub{},
		AlchnSpcmInfo:   &alchnSpcmInfoUseCaseStub{},
	})
	if err == nil || err.Error() != "alchnIlstrSearch use case is required" {
		t.Errorf("error = %v, want alchnIlstrSearch use case is required", err)
	}
}

func TestAddToolsRequiresAlchnIlstrInfoUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		AlchnIlstrSearch: &alchnIlstrSearchUseCaseStub{},
		AlchnSpcmSearch:  &alchnSpcmSearchUseCaseStub{},
		AlchnSpcmInfo:    &alchnSpcmInfoUseCaseStub{},
	})
	if err == nil || err.Error() != "alchnIlstrInfo use case is required" {
		t.Errorf("error = %v, want alchnIlstrInfo use case is required", err)
	}
}

func TestAddToolsRequiresAlchnSpcmSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		AlchnIlstrSearch: &alchnIlstrSearchUseCaseStub{},
		AlchnIlstrInfo:   &alchnIlstrInfoUseCaseStub{},
		AlchnSpcmInfo:    &alchnSpcmInfoUseCaseStub{},
	})
	if err == nil || err.Error() != "alchnSpcmSearch use case is required" {
		t.Errorf("error = %v, want alchnSpcmSearch use case is required", err)
	}
}

func TestAddToolsRequiresAlchnSpcmInfoUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		AlchnIlstrSearch: &alchnIlstrSearchUseCaseStub{},
		AlchnIlstrInfo:   &alchnIlstrInfoUseCaseStub{},
		AlchnSpcmSearch:  &alchnSpcmSearchUseCaseStub{},
	})
	if err == nil || err.Error() != "alchnSpcmInfo use case is required" {
		t.Errorf("error = %v, want alchnSpcmInfo use case is required", err)
	}
}
