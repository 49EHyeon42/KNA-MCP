package application

// ScnmSearchQuery contains the lichen scientific name list conditions.
type ScnmSearchQuery struct {
	PageNo    int
	NumOfRows int
	ReqGnrlNm string
	ReqScnm   string
	DateFrom  string
	DateTo    string
}

// ScnmSearchResult contains a page of lichen scientific name information.
type ScnmSearchResult struct {
	Items      []ScnmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// ScnmSearchItem contains one lichen scientific name record.
type ScnmSearchItem struct {
	StpltScnmRltnCdNm string
	ClassKorNm        string
	ClassNm           string
	FalmNm            string
	FalnKorNm         string
	GenusKorNm        string
	GenusNm           string
	LastUpdtDtm       string
	LchnGnrlNm        string
	LchnScnm          string
	LchnScnmID        string
	OrdKorNm          string
	OrdNm             string
	PhylumKorNm       string
	PhylumNm          string
}
