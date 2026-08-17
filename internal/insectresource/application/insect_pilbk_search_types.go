package application

// InsectPilbkSearchQuery contains the insect pictorial book search conditions.
type InsectPilbkSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
	DateFrom     string
	DateTo       string
}

// InsectPilbkSearchResult contains a page of insect pictorial book search results.
type InsectPilbkSearchResult struct {
	Items      []InsectPilbkSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// InsectPilbkSearchItem contains one insect pictorial book search result.
type InsectPilbkSearchItem struct {
	FamilyKorNm    string
	FamilyNm       string
	GenusKorNm     string
	GenusNm        string
	InsctGnrlNm    string
	InsctPilbkNo   string
	InsctSpecsScnm string
	LastUpdtDtm    string
}
