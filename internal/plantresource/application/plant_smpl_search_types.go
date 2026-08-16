package application

// PlantSmplSearchQuery contains the plant sample search conditions.
type PlantSmplSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// PlantSmplSearchResult contains a page of plant sample search results.
type PlantSmplSearchResult struct {
	Items      []PlantSmplSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantSmplSearchItem contains one plant sample search result.
type PlantSmplSearchItem struct {
	Cnt            int
	FamilyKorNm    string
	FamilyNm       string
	PlantGnrlNm    string
	PlantSpecsID   string
	PlantSpecsScnm string
}
