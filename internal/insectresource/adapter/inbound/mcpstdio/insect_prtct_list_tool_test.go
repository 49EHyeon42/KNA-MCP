package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
)

type insectPrtctListUseCaseStub struct {
	query  application.InsectPrtctListQuery
	result application.InsectPrtctListResult
	err    error
}

func (s *insectPrtctListUseCaseStub) InsectPrtctList(_ context.Context, query application.InsectPrtctListQuery) (application.InsectPrtctListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestInsectPrtctListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &insectPrtctListUseCaseStub{result: application.InsectPrtctListResult{
		Items: []application.InsectPrtctListItem{{
			FamilyKorNm:    "family Korean name",
			FamilyNm:       "family name",
			InsctGnrlNm:    "insect general name",
			InsctPcmtt:     "endangered classification",
			InsctPilbkNo:   "insect pictorial book number",
			InsctSpecsScnm: "insect species scientific name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addInsectPrtctListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "insect_resource_insect_prtct_list",
		map[string]string{
			"pageNo":    "페이지번호 (1 이상)",
			"numOfRows": "한 페이지 결과 수 (1 이상)",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "insect_resource_insect_prtct_list", "산림청 국립수목원 멸종위기곤충 목록을 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "insect_resource_insect_prtct_list",
		Arguments: map[string]any{
			"pageNo":    2,
			"numOfRows": 10,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.InsectPrtctListQuery{PageNo: 2, NumOfRows: 10}
	if !reflect.DeepEqual(useCase.query, wantQuery) {
		t.Errorf("query = %#v, want %#v", useCase.query, wantQuery)
	}

	output, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content = %#v", result.StructuredContent)
	}
	if output["numOfRows"] != float64(10) || output["pageNo"] != float64(2) || output["totalCount"] != float64(21) {
		t.Errorf("pagination = %#v", output)
	}
	checkKeys(t, output, "items", "numOfRows", "pageNo", "totalCount")
	items, ok := output["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v", output["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item = %#v", items[0])
	}
	wantItem := map[string]string{
		"familyKorNm":    "family Korean name",
		"familyNm":       "family name",
		"insctGnrlNm":    "insect general name",
		"insctPcmtt":     "endangered classification",
		"insctPilbkNo":   "insect pictorial book number",
		"insctSpecsScnm": "insect species scientific name",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "insect_resource_insect_prtct_list",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한페이지 결과수",
			"pageNo":     "페이지 번호",
			"totalCount": "전체 건수",
		}, map[string]string{
			"familyKorNm":    "과국명",
			"familyNm":       "과명",
			"insctGnrlNm":    "국명(곤충명)",
			"insctPcmtt":     "멸종위기구분",
			"insctPilbkNo":   "곤충도감번호",
			"insctSpecsScnm": "학명",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "insect_resource_insect_prtct_list",
		Arguments: map[string]any{"pageNo": 1, "numOfRows": 1},
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
