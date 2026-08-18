package application

// ChildPilbkSearchQuery contains the child pictorial book search parameters.
type ChildPilbkSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// ChildPilbkSearchResult contains a page of child pictorial book search results.
type ChildPilbkSearchResult struct {
	Items      []ChildPilbkSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// ChildPilbkSearchItem contains one child pictorial book search result.
type ChildPilbkSearchItem struct {
	BiogyNm           string
	ChildLvbngPilbkNo string
	FamilyKorNm       string
	FamilyNm          string
	LvbngTpcdNm       string
	LvngKrlngNm       string
}
