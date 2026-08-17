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

type insectSmplSearchUseCaseStub struct {
	query  application.InsectSmplSearchQuery
	result application.InsectSmplSearchResult
	err    error
}

func (s *insectSmplSearchUseCaseStub) InsectSmplSearch(_ context.Context, query application.InsectSmplSearchQuery) (application.InsectSmplSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestInsectSmplSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &insectSmplSearchUseCaseStub{result: application.InsectSmplSearchResult{
		Items: []application.InsectSmplSearchItem{{
			Cnt:              "sample count",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			GenusKorNm:       "genus Korean name",
			GenusNm:          "genus name",
			InsctGnrlNm:      "insect general name",
			InsctSpecsID:     "insect species ID",
			InsctSpecsScnm:   "insect species scientific name",
			SubFamilyKorNm:   "subfamily Korean name",
			SubFamilyNm:      "subfamily name",
			SuperFamilyKorNm: "superfamily Korean name",
			SuperFamilyNm:    "superfamily name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addInsectSmplSearchTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "insect_resource_insect_smpl_search",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 곤충의 학명 또는 국명 (대소문자를 구분하지 않는 부분 문자열 검색)",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "insect_resource_insect_smpl_search", "산림청 국립수목원 곤충표본 목록을 검색합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "insect_resource_insect_smpl_search",
		Arguments: map[string]any{
			"pageNo":       2,
			"numOfRows":    10,
			"reqSearchWrd": "test-search-word",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.InsectSmplSearchQuery{
		PageNo:       2,
		NumOfRows:    10,
		ReqSearchWrd: "test-search-word",
	}
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
	wantItem := map[string]any{
		"cnt":              "sample count",
		"familyKorNm":      "family Korean name",
		"familyNm":         "family name",
		"genusKorNm":       "genus Korean name",
		"genusNm":          "genus name",
		"insctGnrlNm":      "insect general name",
		"insctSpecsId":     "insect species ID",
		"insctSpecsScnm":   "insect species scientific name",
		"subFamilyKorNm":   "subfamily Korean name",
		"subFamilyNm":      "subfamily name",
		"superFamilyKorNm": "superfamily Korean name",
		"superFamilyNm":    "superfamily name",
	}
	if !reflect.DeepEqual(item, wantItem) {
		t.Errorf("item = %#v, want %#v", item, wantItem)
	}
	checkToolOutputSchema(t, ctx, clientSession, "insect_resource_insect_smpl_search",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지 당 건 수",
			"pageNo":     "페이지 번호",
			"totalCount": "전체 건 수",
		}, map[string]string{
			"cnt":              "표본수",
			"familyKorNm":      "과국명",
			"familyNm":         "과명",
			"genusKorNm":       "속국명",
			"genusNm":          "속명",
			"insctGnrlNm":      "국명(곤충명)",
			"insctSpecsId":     "곤충종ID",
			"insctSpecsScnm":   "학명",
			"subFamilyKorNm":   "아과국명",
			"subFamilyNm":      "아과명",
			"superFamilyKorNm": "상과국명",
			"superFamilyNm":    "상과명",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "insect_resource_insect_smpl_search",
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
