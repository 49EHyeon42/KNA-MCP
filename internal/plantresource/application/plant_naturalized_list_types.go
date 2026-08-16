package application

// PlantNaturalizedListQuery contains the naturalized plant list conditions.
type PlantNaturalizedListQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
	DateFrom     string
	DateTo       string
}

// PlantNaturalizedListResult contains a page of naturalized plant information.
type PlantNaturalizedListResult struct {
	Items      []PlantNaturalizedListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantNaturalizedListItem contains one naturalized plant record.
type PlantNaturalizedListItem struct {
	AgpFamilyNm        string
	APGFamilyKorNm     string
	BlprdEnmnt         string
	BlprdStmnt         string
	DistrAraDscrt      string
	EclgDstrbYn        string
	ExtcPlantCdNm      string
	FamilyKorNm        string
	FamilyNm           string
	FrtTpcdNm          string
	LastUpdtDtm        string
	NtldgTpcdNm        string
	NtrlzEraTpcdNm     string
	OrplcNm            string
	PlantBrdgFomTpcdNm string
	PlantDistrGrcd     string
	PlantDistrQntt     string
	PlantDistrQnttGrcd string
	PlantEngNm         string
	PlantGnrlNm        string
	PlantJpnNm         string
	PlantLfcclTpcdNm   string
	PlantSpecsScnm     string
}
