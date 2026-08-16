package application

// PlantRareListQuery contains the rare plant list conditions.
type PlantRareListQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// PlantRareListResult contains a page of rare plant information.
type PlantRareListResult struct {
	Items      []PlantRareListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantRareListItem contains one rare plant record.
type PlantRareListItem struct {
	AgpFamilyNm      string
	APGFamilyKorNm   string
	ExtrmCrssScls1Yn string
	ExtrmCrssScls2Yn string
	FamilyKorNm      string
	FamilyNm         string
	PlantGnrlNm      string
	PlantSpecsScnm   string
	RareTpcdNm       string
}
