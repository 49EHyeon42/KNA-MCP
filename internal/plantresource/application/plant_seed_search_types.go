package application

// PlantSeedSearchQuery contains the plant seed search conditions.
type PlantSeedSearchQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
	DateFrom     string
	DateTo       string
}

// PlantSeedSearchResult contains a page of plant seed search results.
type PlantSeedSearchResult struct {
	Items      []PlantSeedSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantSeedSearchItem contains one plant seed search result.
type PlantSeedSearchItem struct {
	APGFamilyKorNm   string
	APGFamilyNm      string
	BlprdEnmnt       string
	BlprdStmnt       string
	ClrngMthodCdNm   string
	FamilyKorNm      string
	FamilyNm         string
	FritCdNm         string
	FrssnEnmnt       string
	FrssnStmnt       string
	LastUpdtDtm      string
	PlantGnrlNm      string
	PlantSpecsScnm   string
	RfrncLtrtrCont   string
	SeedCtsrfcDesc   string
	SeedCtsrfcTpcdNm string
	SeedEmbrTpcdNm   string
	SeedMnmmBrdth    string
	SeedMnmmLngth    string
	SeedMxmmBrdth    string
	SeedMxmmLngth    string
	SeedShpDesc      string
	SeedShpTpcdNm    string
	SeedSpecsID      string
	SeedTpcdNm       string
}
