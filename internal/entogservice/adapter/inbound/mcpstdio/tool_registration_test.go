package mcpstdio

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

func TestAddToolsRegistersAllEntogServiceTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, UseCases{EntogIlstrSearch: &entogIlstrSearchUseCaseStub{}}); err != nil {
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
	if len(result.Tools) != 1 {
		t.Fatalf("tools = %#v, want one tool", result.Tools)
	}
	if name := result.Tools[0].Name; name != "entog_service_entog_ilstr_search" {
		t.Errorf("tool name = %q", name)
	} else if !regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`).MatchString(name) {
		t.Errorf("tool name = %q, want 1 to 128 allowed characters", name)
	}
}

func TestAddToolsRequiresEntogIlstrSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{})
	if err == nil || err.Error() != "entogIlstrSearch use case is required" {
		t.Errorf("error = %v, want entogIlstrSearch use case is required", err)
	}
}
