package application

// ScnmSearchQuery contains the insect scientific name list conditions.
type ScnmSearchQuery struct {
	PageNo    int
	NumOfRows int
	ReqGnrlNm string
	ReqScnm   string
	DateFrom  string
	DateTo    string
}

// ScnmSearchResult contains a page of insect scientific name information.
type ScnmSearchResult struct {
	Items      []ScnmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// ScnmSearchItem contains one insect scientific name record.
type ScnmSearchItem struct {
	SuperFalmNm       string
	ClassKorNm        string
	ClassNm           string
	FalmKorNm         string
	FalmNm            string
	GenusKorNm        string
	GenusNm           string
	InsctGnrlNm       string
	InsctScnmID       string
	InsctSpecsScnm    string
	LastUpdtDtm       string
	OrdKorNm          string
	OrdNm             string
	StpltScnmRltnCdNm string
	SubFalmKorNm      string
	SubFalmNm         string
	SuperFalmKorNm    string
}
