package application

// PlantPilbkSearchQuery contains the plant pictorial book search conditions.
type PlantPilbkSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
	DateFrom     string
	DateTo       string
}

// PlantPilbkSearchResult contains a page of plant pictorial book search results.
type PlantPilbkSearchResult struct {
	Items      []PlantPilbkSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantPilbkSearchItem contains one plant pictorial book search result.
type PlantPilbkSearchItem struct {
	APGFamilyKorNm string
	APGFamilyNm    string
	FamilyKorNm    string
	FamilyNm       string
	GenusKorNm     string
	GenusNm        string
	LastUpdtDtm    string
	NotRcmmGnrlNm  string
	PlantGnrlNm    string
	PlantPilbkNo   string
	PlantSpecsScnm string
}
