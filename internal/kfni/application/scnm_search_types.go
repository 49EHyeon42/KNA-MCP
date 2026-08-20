package application

// ScnmSearchQuery contains the fungi scientific name list conditions.
type ScnmSearchQuery struct {
	PageNo    int
	NumOfRows int
	ReqGnrlNm string
	ReqScnm   string
	DateFrom  string
	DateTo    string
}

// ScnmSearchResult contains a page of fungi scientific name information.
type ScnmSearchResult struct {
	Items      []ScnmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// ScnmSearchItem contains one fungi scientific name record.
type ScnmSearchItem struct {
	StpltScnmRltnCdNm string
	ClassKorNm        string
	ClassNm           string
	FalmNm            string
	FalnKorNm         string
	FngsGnrlNm        string
	FngsScnm          string
	FngsScnmID        string
	GenusKorNm        string
	GenusNm           string
	LastUpdtDtm       string
	OrdKorNm          string
	OrdNm             string
	PhylumKorNm       string
	PhylumNm          string
}
