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

type insectSmplUnitListUseCaseStub struct {
	query  application.InsectSmplUnitListQuery
	result application.InsectSmplUnitListResult
	err    error
}

func (s *insectSmplUnitListUseCaseStub) InsectSmplUnitList(_ context.Context, query application.InsectSmplUnitListQuery) (application.InsectSmplUnitListResult, error) {
	s.query = query
	return s.result, s.err
}

func TestInsectSmplUnitListTool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	useCase := &insectSmplUnitListUseCaseStub{result: application.InsectSmplUnitListResult{
		Items: []application.InsectSmplUnitListItem{{
			BspcsInsttNm:       "specimen holding institution",
			ClarHaslvVal:       "collection site elevation",
			SmplCllcnDt:        "specimen collection date",
			GynndTpcd:          "sex type",
			HbttTpcd:           "habitat type",
			InsctSmplNo:        "insect specimen number",
			InsctSpecsID:       "insect species ID",
			InsctSpecsScnm:     "insect species scientific name",
			LabelUsgCllcnNmplc: "label collection place name",
			LastUpdtDtm:        "last update date time",
			PrsrtStcd:          "preservation status",
			SlistTpcd:          "minute insect type",
			SmplKindCd:         "specimen type",
			TorsoLngth:         "torso length",
			WingLngth:          "wing length",
			InsctGnrlNm:        "insect general name",
		}},
		NumOfRows:  10,
		PageNo:     2,
		TotalCount: 21,
	}}
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcpserver.NewServer()
	addInsectSmplUnitListTool(server, useCase)
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
	checkToolInputSchema(t, ctx, clientSession, "insect_resource_insect_smpl_unit_list",
		map[string]string{
			"pageNo":          "페이지번호 (1 이상)",
			"numOfRows":       "한 페이지 결과 수 (1 이상)",
			"reqInsctSpecsId": "검색할 곤충 종ID (insectSmplSearch 결과의 insctSpecsId)",
		},
		[]string{"pageNo", "numOfRows", "reqInsctSpecsId"},
	)
	checkToolDescription(t, ctx, clientSession, "insect_resource_insect_smpl_unit_list", "산림청 국립수목원 곤충표본 상세정보 목록을 조회합니다.")

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "insect_resource_insect_smpl_unit_list",
		Arguments: map[string]any{
			"pageNo":          2,
			"numOfRows":       10,
			"reqInsctSpecsId": "test-insect-species-id",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatalf("tool error: %#v", result.Content)
	}

	wantQuery := application.InsectSmplUnitListQuery{
		PageNo:          2,
		NumOfRows:       10,
		ReqInsctSpecsID: "test-insect-species-id",
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
		"bspcsInsttNm":       "specimen holding institution",
		"clarHaslvVal":       "collection site elevation",
		"smplCllcnDt":        "specimen collection date",
		"gynndTpcd":          "sex type",
		"hbttTpcd":           "habitat type",
		"insctSmplNo":        "insect specimen number",
		"insctSpecsId":       "insect species ID",
		"insctSpecsScnm":     "insect species scientific name",
		"labelUsgCllcnNmplc": "label collection place name",
		"lastUpdtDtm":        "last update date time",
		"prsrtStcd":          "preservation status",
		"slistTpcd":          "minute insect type",
		"smplKindCd":         "specimen type",
		"torsoLngth":         "torso length",
		"wingLngth":          "wing length",
		"insctGnrlNm":        "insect general name",
	}
	if len(item) != len(wantItem) {
		t.Errorf("item key count = %d, want %d", len(item), len(wantItem))
	}
	for key, want := range wantItem {
		if got := item[key]; got != want {
			t.Errorf("item %s = %#v, want %q", key, got, want)
		}
	}
	checkToolOutputSchema(t, ctx, clientSession, "insect_resource_insect_smpl_unit_list",
		map[string]string{
			"items":      "조회 결과 목록",
			"numOfRows":  "한페이지 결과 수",
			"pageNo":     "페이지 번호",
			"totalCount": "전체 결과  건 수",
		}, map[string]string{
			"bspcsInsttNm":       "표본소장기관",
			"clarHaslvVal":       "채집지해발고도",
			"smplCllcnDt":        "표본채집일",
			"gynndTpcd":          "암수구분",
			"hbttTpcd":           "서식지구분",
			"insctSmplNo":        "곤충표본번호",
			"insctSpecsId":       "곤충종ID",
			"insctSpecsScnm":     "학명",
			"labelUsgCllcnNmplc": "라벨용채집지명",
			"lastUpdtDtm":        "최종수정일",
			"prsrtStcd":          "보존상태",
			"slistTpcd":          "미소곤충구분",
			"smplKindCd":         "표본종류",
			"torsoLngth":         "몸통길이",
			"wingLngth":          "날개길이",
			"insctGnrlNm":        "국명(곤충명)",
		})

	useCase.err = errors.New("upstream unavailable")
	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "insect_resource_insect_smpl_unit_list",
		Arguments: map[string]any{
			"pageNo":          1,
			"numOfRows":       1,
			"reqInsctSpecsId": "test-insect-species-id",
		},
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
