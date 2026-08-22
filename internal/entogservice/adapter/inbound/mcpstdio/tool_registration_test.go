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

func TestAddToolsRegistersAllEntogServiceTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, UseCases{
		EntogIlstrSearch: &entogIlstrSearchUseCaseStub{},
		EntogIlstrInfo:   &entogIlstrInfoUseCaseStub{},
		EntogSpcmSearch:  &entogSpcmSearchUseCaseStub{},
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
	want := []string{"entog_service_entog_ilstr_info", "entog_service_entog_ilstr_search", "entog_service_entog_spcm_search"}
	if !slices.Equal(names, want) {
		t.Fatalf("tools = %#v, want %#v", names, want)
	}
}

func TestAddToolsRequiresEntogIlstrInfoUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		EntogIlstrSearch: &entogIlstrSearchUseCaseStub{},
		EntogSpcmSearch:  &entogSpcmSearchUseCaseStub{},
	})
	if err == nil || err.Error() != "entogIlstrInfo use case is required" {
		t.Errorf("error = %v, want entogIlstrInfo use case is required", err)
	}
}

func TestAddToolsRequiresEntogSpcmSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{
		EntogIlstrSearch: &entogIlstrSearchUseCaseStub{},
		EntogIlstrInfo:   &entogIlstrInfoUseCaseStub{},
	})
	if err == nil || err.Error() != "entogSpcmSearch use case is required" {
		t.Errorf("error = %v, want entogSpcmSearch use case is required", err)
	}
}

func TestAddToolsRequiresEntogIlstrSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{})
	if err == nil || err.Error() != "entogIlstrSearch use case is required" {
		t.Errorf("error = %v, want entogIlstrSearch use case is required", err)
	}
}
