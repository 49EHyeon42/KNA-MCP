package application

// FngsSmplSearchQuery contains the fungi sample search conditions.
type FngsSmplSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// FngsSmplSearchResult contains a page of fungi sample search results.
type FngsSmplSearchResult struct {
	Items      []FngsSmplSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// FngsSmplSearchItem contains one fungi sample search result.
type FngsSmplSearchItem struct {
	Cnt         string
	FamilyKorNm string
	FamilyNm    string
	FngsGnrlNm  string
	FngsID      string
	FngsScnm    string
	GenusKorNm  string
	GenusNm     string
}
