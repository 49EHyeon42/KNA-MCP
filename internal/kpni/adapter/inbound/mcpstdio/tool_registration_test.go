package mcpstdio

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

func TestAddToolsRegistersAllKpniTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, UseCases{ScnmSearch: &scnmSearchUseCaseStub{}}); err != nil {
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
	if len(result.Tools) != 1 || result.Tools[0].Name != "kpni_scnm_search" {
		t.Fatalf("tools = %#v, want kpni_scnm_search", result.Tools)
	}
	if !regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`).MatchString(result.Tools[0].Name) {
		t.Errorf("tool name = %q, want 1 to 128 allowed characters", result.Tools[0].Name)
	}
}

func TestAddToolsRequiresScnmSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{})
	if err == nil || err.Error() != "scnmSearch use case is required" {
		t.Errorf("error = %v, want scnmSearch use case is required", err)
	}
}
