package application

// PlantSeedGrmntListQuery contains the plant seed germination list conditions.
type PlantSeedGrmntListQuery struct {
	PageNo         int
	NumOfRows      int
	ReqSeedSpecsID string
}

// PlantSeedGrmntListResult contains a page of plant seed germination information.
type PlantSeedGrmntListResult struct {
	Items      []PlantSeedGrmntListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantSeedGrmntListItem contains one plant seed germination record.
type PlantSeedGrmntListItem struct {
	AvrgGrmntDcnt     string
	GrmntBfrPrcesCont string
	GrmntClmdmCont    string
	GrmntDscrt        string
	GrmntExprmNo      string
	GrmntExprmSeq     string
	GrmntLightCndtn   string
	GrmntRt           string
	GrmntTmpCndtn     string
	PlantGnrlNm       string
	SeedNo            string
	SeedSpecsID       string
}
