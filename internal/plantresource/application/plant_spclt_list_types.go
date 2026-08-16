package application

// PlantSpcltListQuery contains the endemic plant list conditions.
type PlantSpcltListQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// PlantSpcltListResult contains a page of endemic plant information.
type PlantSpcltListResult struct {
	Items      []PlantSpcltListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantSpcltListItem contains one endemic plant record.
type PlantSpcltListItem struct {
	AgpFamilyKorNm     string
	AgpFamilyNm        string
	ExtrmCrssScls1Yn   string
	ExtrmCrssScls2Yn   string
	FamilyKorNm        string
	FamilyNm           string
	PlantBrdgFomTpcdNm string
	PlantGnrlNm        string
	PlantSpecsScnm     string
}
