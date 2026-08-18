package application

// FngsSmplUnitListQuery contains the fungi specimen detail list conditions.
type FngsSmplUnitListQuery struct {
	PageNo    int
	NumOfRows int
	ReqFngsID string
}

// FngsSmplUnitListResult contains a page of fungi specimen details.
type FngsSmplUnitListResult struct {
	Items      []FngsSmplUnitListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// FngsSmplUnitListItem contains one fungi specimen detail.
type FngsSmplUnitListItem struct {
	ClarDtlDscrt     string
	ClarHaslvVal     string
	CllcrNm          string
	FamilyKorNm      string
	FamilyNm         string
	FngsEclgTpcdNm   string
	FngsGnrlNm       string
	FngsID           string
	FngsScnm         string
	FngsSmplKindCdNm string
	FngsSmplNo       string
	GenusKorNm       string
	GenusNm          string
	HbttChrcrCont    string
	HstCont          string
	LastUpdtDtm      string
	SmplCllcnDt      string
}
