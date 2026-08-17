package mcpstdio

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/49EHyeon42/KNA-MCP/internal/mcpstdio"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

type plantSeedSearchUseCaseStub struct {
	query  application.PlantSeedSearchQuery
	result application.PlantSeedSearchResult
	err    error
}

func (s *plantSeedSearchUseCaseStub) PlantSeedSearch(_ context.Context, query application.PlantSeedSearchQuery) (application.PlantSeedSearchResult, error) {
	s.query = query
	return s.result, s.err
}

func TestPlantSeedSearchTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &plantSeedSearchUseCaseStub{result: application.PlantSeedSearchResult{
		Items: []application.PlantSeedSearchItem{{
			APGFamilyKorNm:   "apg family Korean name",
			APGFamilyNm:      "apg family name",
			BlprdEnmnt:       "bloom period end month",
			BlprdStmnt:       "bloom period start month",
			ClrngMthodCdNm:   "cleaning method",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			FritCdNm:         "fruit type",
			FrssnEnmnt:       "fruit season end month",
			FrssnStmnt:       "fruit season start month",
			LastUpdtDtm:      "last update date time",
			PlantGnrlNm:      "plant general name",
			PlantSpecsScnm:   "plant species scientific name",
			RfrncLtrtrCont:   "reference literature",
			SeedCtsrfcDesc:   "seed surface description",
			SeedCtsrfcTpcdNm: "seed surface type",
			SeedEmbrTpcdNm:   "seed embryo type",
			SeedMnmmBrdth:    "seed minimum breadth",
			SeedMnmmLngth:    "seed minimum length",
			SeedMxmmBrdth:    "seed maximum breadth",
			SeedMxmmLngth:    "seed maximum length",
			SeedShpDesc:      "seed shape description",
			SeedShpTpcdNm:    "seed shape type",
			SeedSpecsID:      "seed species ID",
			SeedTpcdNm:       "seed type",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addPlantSeedSearchTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "plant_resource_plant_seed_search",
		map[string]string{
			"pageNo":       "페이지번호 (1 이상)",
			"numOfRows":    "한 페이지 결과 수 (1 이상)",
			"reqSearchWrd": "검색할 식물종자의 국명 또는 학명",
		},
		[]string{"pageNo", "numOfRows"},
	)
	checkToolDescription(t, ctx, clientSession, "plant_resource_plant_seed_search", "산림청 국립수목원 식물종자 기본정보 목록을 검색합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "plant_resource_plant_seed_search",
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

	wantQuery := application.PlantSeedSearchQuery{
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
	wantItem := map[string]string{
		"apgFamilyKorNm":   "apg family Korean name",
		"apgFamilyNm":      "apg family name",
		"blprdEnmnt":       "bloom period end month",
		"blprdStmnt":       "bloom period start month",
		"clrngMthodCdNm":   "cleaning method",
		"familyKorNm":      "family Korean name",
		"familyNm":         "family name",
		"fritCdNm":         "fruit type",
		"frssnEnmnt":       "fruit season end month",
		"frssnStmnt":       "fruit season start month",
		"lastUpdtDtm":      "last update date time",
		"plantGnrlNm":      "plant general name",
		"plantSpecsScnm":   "plant species scientific name",
		"rfrncLtrtrCont":   "reference literature",
		"seedCtsrfcDesc":   "seed surface description",
		"seedCtsrfcTpcdNm": "seed surface type",
		"seedEmbrTpcdNm":   "seed embryo type",
		"seedMnmmBrdth":    "seed minimum breadth",
		"seedMnmmLngth":    "seed minimum length",
		"seedMxmmBrdth":    "seed maximum breadth",
		"seedMxmmLngth":    "seed maximum length",
		"seedShpDesc":      "seed shape description",
		"seedShpTpcdNm":    "seed shape type",
		"seedSpecsId":      "seed species ID",
		"seedTpcdNm":       "seed type",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "plant_resource_plant_seed_search",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한 페이지 결과 수",
			"pageNo":     "페이지번호",
			"totalCount": "전체 검색 결과 수",
		}, map[string]string{
			"apgFamilyKorNm":   "APG과국명",
			"apgFamilyNm":      "APG과명",
			"blprdEnmnt":       "개화기종료일",
			"blprdStmnt":       "개화기시작일",
			"clrngMthodCdNm":   "정선방법",
			"familyKorNm":      "과국명",
			"familyNm":         "과명",
			"fritCdNm":         "열매형태",
			"frssnEnmnt":       "결실기종료일",
			"frssnStmnt":       "결실기시작일",
			"lastUpdtDtm":      "최종수정일",
			"plantGnrlNm":      "국명(식물명)",
			"plantSpecsScnm":   "학명",
			"rfrncLtrtrCont":   "참고문헌",
			"seedCtsrfcDesc":   "종자표면형태설명",
			"seedCtsrfcTpcdNm": "종자표면형태",
			"seedEmbrTpcdNm":   "배아형태",
			"seedMnmmBrdth":    "종자최소너비",
			"seedMnmmLngth":    "종자최소길이",
			"seedMxmmBrdth":    "종자최대너비",
			"seedMxmmLngth":    "종자최대길이",
			"seedShpDesc":      "종자형태설명",
			"seedShpTpcdNm":    "종자형태",
			"seedSpecsId":      "종자종ID",
			"seedTpcdNm":       "종자구분",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "plant_resource_plant_seed_search",
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
