package application

// AlchnSpcmSearchQuery contains the lichen specimen search conditions.
type AlchnSpcmSearchQuery struct {
	St        string
	Sw        string
	DateGbn   string
	DateFrom  string
	DateTo    string
	NumOfRows int
	PageNo    int
}

// AlchnSpcmSearchResult contains a page of lichen specimen search results.
type AlchnSpcmSearchResult struct {
	Items      []AlchnSpcmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// AlchnSpcmSearchItem contains one lichen specimen search result.
type AlchnSpcmSearchItem struct {
	Btnc         string
	CltrNm       string
	CprtCtnt     string
	DetailYn     string
	EngNm        string
	FamilyKorNm  string
	FamilyNm     string
	FrstRgstnDtm string
	GenusKorNm   string
	GenusNm      string
	ImgURL       string
	JapNm        string
	LastUpdtDtm  string
	LchnGnrlNm   string
	LchnScnmID   string
	LchnSmplNo   string
	PrkNm        string
}
