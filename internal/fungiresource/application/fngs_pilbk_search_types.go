package application

// FngsPilbkSearchQuery contains the fungi pictorial book search conditions.
type FngsPilbkSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
	DateFrom     string
	DateTo       string
}

// FngsPilbkSearchResult contains a page of fungi pictorial book search results.
type FngsPilbkSearchResult struct {
	Items      []FngsPilbkSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// FngsPilbkSearchItem contains one fungi pictorial book search result.
type FngsPilbkSearchItem struct {
	FamilyKorNm string
	FamilyNm    string
	FngsGnrlNm  string
	FngsPilbkNo string
	FngsScnm    string
	GenusKorNm  string
	GenusNm     string
	LastUpdtDtm string
}
