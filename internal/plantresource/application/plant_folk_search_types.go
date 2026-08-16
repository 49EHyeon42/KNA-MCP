package application

// PlantFolkSearchQuery contains the folk plant search conditions.
type PlantFolkSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// PlantFolkSearchResult contains a page of folk plant search results.
type PlantFolkSearchResult struct {
	Items      []PlantFolkSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantFolkSearchItem contains one folk plant search result.
type PlantFolkSearchItem struct {
	FlcstPlantIdntfDscrt string
	FlpltID              string
	PlantBrdgFomTpcdNm   string
	PlantGnrlNm          string
	PlantSpecsScnm       string
	Ptnt                 string
}
