package application

// PlantSmplUnitListQuery contains the plant specimen detail list conditions.
type PlantSmplUnitListQuery struct {
	PageNo          int
	NumOfRows       int
	ReqPlantSpecsID string
}

// PlantSmplUnitListResult contains a page of plant specimen details.
type PlantSmplUnitListResult struct {
	Items      []PlantSmplUnitListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantSmplUnitListItem contains one plant specimen detail.
type PlantSmplUnitListItem struct {
	AgpFamilyKorNm     string
	AgpFamilyNm        string
	BspcsInsttNm       string
	ClarHaslvVal       string
	ClarNm             string
	CllcrNm            string
	FamilyKorNm        string
	FamilyNm           string
	HbttChrcrCont      string
	HbttTpcdNm         string
	PlantBrdgFomTpcdNm string
	PlantGnrlNm        string
	PlantPilbkNo       string
	PlantSmplNo        string
	PlantSpecsID       string
	PlantSpecsScnm     string
	SmplCllcnDt        string
	SmplClnyNm         string
	SmplKindCdNm       string
	SmplWrdt           string
	VgttnTpeCdNm       string
}
