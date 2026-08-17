package application

// InsectSmplSearchQuery contains the insect sample search conditions.
type InsectSmplSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// InsectSmplSearchResult contains a page of insect sample search results.
type InsectSmplSearchResult struct {
	Items      []InsectSmplSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// InsectSmplSearchItem contains one insect sample search result.
type InsectSmplSearchItem struct {
	Cnt              string
	FamilyKorNm      string
	FamilyNm         string
	GenusKorNm       string
	GenusNm          string
	InsctGnrlNm      string
	InsctSpecsID     string
	InsctSpecsScnm   string
	SubFamilyKorNm   string
	SubFamilyNm      string
	SuperFamilyKorNm string
	SuperFamilyNm    string
}
