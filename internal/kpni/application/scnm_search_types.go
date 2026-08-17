package application

// ScnmSearchQuery contains the scientific name list conditions.
type ScnmSearchQuery struct {
	PageNo    int
	NumOfRows int
	ReqGnrlNm string
	ReqScnm   string
	DateFrom  string
	DateTo    string
}

// ScnmSearchResult contains a page of scientific name information.
type ScnmSearchResult struct {
	Items      []ScnmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// ScnmSearchItem contains one scientific name record.
type ScnmSearchItem struct {
	ClassKorNm          string
	ClassNm             string
	FalmKorNm           string
	FalmNm              string
	GenusKorNm          string
	GenusNm             string
	LastUpdtDtm         string
	OrdKorNm            string
	OrdNm               string
	PhylumKorNm         string
	PhylumNm            string
	PlantGnrlNm         string
	PlantScnmID         string
	PlantSpecsClsscCdNm string
	PlantSpecsScnm      string
	StpltScnmRltnCdNm   string
	SubClassKorNm       string
	SubClassNm          string
}
