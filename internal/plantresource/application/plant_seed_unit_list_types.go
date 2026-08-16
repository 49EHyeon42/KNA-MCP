package application

// PlantSeedUnitListQuery contains the plant seed unit list conditions.
type PlantSeedUnitListQuery struct {
	PageNo         int
	NumOfRows      int
	ReqSeedSpecsID string
}

// PlantSeedUnitListResult contains a page of plant seed unit information.
type PlantSeedUnitListResult struct {
	Items      []PlantSeedUnitListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantSeedUnitListItem contains one plant seed unit.
type PlantSeedUnitListItem struct {
	CllcnDate        string
	PlantGnrlNm      string
	QualtFllnsRt     string
	SdwghWeght       string
	SeedAdmcn        string
	SeedCllctPlace   string
	SeedHoldGrainCnt string
	SeedHoldQntt     string
	SeedNo           string
	SeedSpecsID      string
	StoreChrcrTpcdNm string
	Vtlfct           string
	VtlfctTestYr     string
}
