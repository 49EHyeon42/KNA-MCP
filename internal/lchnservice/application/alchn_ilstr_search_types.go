package application

// AlchnIlstrSearchQuery contains the lichen pictorial book search conditions.
type AlchnIlstrSearchQuery struct {
	St        string
	Sw        string
	DateGbn   string
	DateFrom  string
	DateTo    string
	NumOfRows int
	PageNo    int
}

// AlchnIlstrSearchResult contains a page of lichen pictorial book search results.
type AlchnIlstrSearchResult struct {
	Items      []AlchnIlstrSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// AlchnIlstrSearchItem contains one lichen pictorial book search result.
type AlchnIlstrSearchItem struct {
	Btnc         string
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
	LchnInfrpNm  string
	LchnPilbkNo  string
	LchnScnmID   string
	LchnTtnm     string
	PrkNm        string
}
