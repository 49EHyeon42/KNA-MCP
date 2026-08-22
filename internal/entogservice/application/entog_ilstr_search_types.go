package application

// EntogIlstrSearchQuery contains the entognath pictorial book search conditions.
type EntogIlstrSearchQuery struct {
	St        string
	Sw        string
	NumOfRows int
	PageNo    int
}

// EntogIlstrSearchResult contains a page of entognath pictorial book search results.
type EntogIlstrSearchResult struct {
	Items      []EntogIlstrSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// EntogIlstrSearchItem contains one entognath pictorial book search result.
type EntogIlstrSearchItem struct {
	Btnc             string
	CprtCtnt         string
	DetailYn         string
	EntogOfnmKrlngNm string
	EntogOfnmScnmID  string
	EntogPilbkNo     string
	FamilyKorNm      string
	FamilyNm         string
	GenusKorNm       string
	GenusNm          string
	ImgURL           string
	OrdKorNm         string
	OrdNm            string
}
