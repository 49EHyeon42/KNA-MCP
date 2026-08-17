package application

// ScnmInfoQuery contains the scientific name detail condition.
type ScnmInfoQuery struct {
	ReqPlantScnmID string
}

// ScnmInfoResult contains scientific name detail information.
type ScnmInfoResult struct {
	Item *ScnmInfoItem
}

// ScnmInfoItem contains one scientific name detail record.
type ScnmInfoItem struct {
	APGFalmKorNm       string
	APGFalmNm          string
	BiogyNmTpcdNm      string
	CltvaYn            string
	EclgDstrbYn        string
	ExtcCncrnsYn       string
	ExtcPlantCdNm      string
	ExtcPlantYn        string
	FalmKorNm          string
	FalmNm             string
	GenusKorNm         string
	GenusNm            string
	LtrtrInfrmNm       string
	PlantBrdgFomTpcdNm string
	PlantChnNm         string
	PlantEngNm         string
	PlantGnrlNm        string
	PlantGnrlNm2       string
	PlantJpnNm         string
	PlantScnmID        string
	PlantSpecsScnm     string
	RareTpcdNm         string
	RelPlantSpecsScnm  string
	RelScnmTpcdNm      string
	Rmrk               string
	RrnssPlantYn       string
	SpcltPlantCdNm     string
}
