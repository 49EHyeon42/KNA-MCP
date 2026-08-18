package application

// AlchnSpcmInfoQuery contains the lichen specimen detail condition.
type AlchnSpcmInfoQuery struct {
	Q1 string
}

// AlchnSpcmInfoResult contains lichen specimen detail information.
type AlchnSpcmInfoResult struct {
	Item *AlchnSpcmInfoItem
}

// AlchnSpcmInfoItem contains one lichen specimen detail record.
type AlchnSpcmInfoItem struct {
	Btnc          string
	ClarDtlDscrt  string
	CllcrNm       string
	CltrNm        string
	CprtCtnt      string
	EngNm         string
	ExmneNm       string
	FamilyKorNm   string
	FamilyNm      string
	FrstRgstnDtm  string
	GenusKorNm    string
	GenusNm       string
	Grdnt         string
	HaslvVal      string
	HbttChrcrCont string
	ImgURL        string
	InsttSmplNo   string
	JapNm         string
	LastUpdtDtm   string
	LchnGnrlNm    string
	LchnScnmID    string
	LchnSmplNo    string
	OrbrnCd       string
	PrkNm         string
	SmplCllcnDt   string
}
