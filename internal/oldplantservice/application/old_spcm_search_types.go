package application

// OldSpcmSearchQuery contains the old plant specimen search conditions.
type OldSpcmSearchQuery struct {
	St        string
	Sw        string
	DateGbn   string
	DateFrom  string
	DateTo    string
	NumOfRows int
	PageNo    int
}

// OldSpcmSearchResult contains a page of old plant specimen search results.
type OldSpcmSearchResult struct {
	Items      []OldSpcmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// OldSpcmSearchItem contains one old plant specimen search result.
type OldSpcmSearchItem struct {
	CprtCtnt       string
	FamlKorNm      string
	FamlNm         string
	FrstRgstnDtm   string
	ImgURL         string
	LastUpdtDtm    string
	PlantGnrlNm    string
	PlantOldSmplNo string
	PlantSpecsScnm string
}
