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

func TestAddToolsRegistersAllKfniTools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := mcpserver.NewServer()
	if err := AddTools(server, UseCases{
		ScnmSearch: &scnmSearchUseCaseStub{},
		ScnmInfo:   &scnmInfoUseCaseStub{},
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
	want := []string{"kfni_scnm_info", "kfni_scnm_search"}
	if !slices.Equal(got, want) {
		t.Fatalf("tools = %#v, want %#v", got, want)
	}
}

func TestAddToolsRequiresScnmInfoUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{ScnmSearch: &scnmSearchUseCaseStub{}})
	if err == nil || err.Error() != "scnmInfo use case is required" {
		t.Errorf("error = %v, want scnmInfo use case is required", err)
	}
}

func TestAddToolsRequiresScnmSearchUseCase(t *testing.T) {
	err := AddTools(mcpserver.NewServer(), UseCases{})
	if err == nil || err.Error() != "scnmSearch use case is required" {
		t.Errorf("error = %v, want scnmSearch use case is required", err)
	}
}
