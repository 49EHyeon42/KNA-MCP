package application

// EntogSpcmSearchQuery contains the entognath specimen search conditions.
type EntogSpcmSearchQuery struct {
	St        string
	Sw        string
	DateGbn   string
	DateFrom  string
	DateTo    string
	NumOfRows int
	PageNo    int
}

// EntogSpcmSearchResult contains a page of entognath specimen search results.
type EntogSpcmSearchResult struct {
	Items      []EntogSpcmSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// EntogSpcmSearchItem contains one entognath specimen search result.
type EntogSpcmSearchItem struct {
	Btnc             string
	ClctDyDesc       string
	CprtCtnt         string
	DetailYn         string
	EntogOfnmKrlngNm string
	EntogSmplNo      string
	FamilyKorNm      string
	FamilyNm         string
	FrstRgstnDtm     string
	GenusKorNm       string
	GenusNm          string
	ImgURL           string
	LastUpdtDtm      string
	OrdKorNm         string
	OrdNm            string
}
