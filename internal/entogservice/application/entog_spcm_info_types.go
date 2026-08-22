package application

// EntogSpcmInfoQuery contains the entognath specimen detail lookup key.
type EntogSpcmInfoQuery struct {
	Q1 string
}

// EntogSpcmInfoResult contains entognath specimen detail information.
type EntogSpcmInfoResult struct {
	Item *EntogSpcmInfoItem
}

// EntogSpcmInfoItem contains one entognath specimen detail result.
type EntogSpcmInfoItem struct {
	Btnc               string
	ChnNm              string
	ClarHaslvVal       string
	ClctDyDesc         string
	CprtCtnt           string
	EngNm              string
	EntogGnrlNm        string
	EntogPilbkNo       string
	EntogSmplNo        string
	FamilyKorNm        string
	FamilyNm           string
	FrstRgstnDtm       string
	GenusKorNm         string
	GenusNm            string
	ImgURL             string
	JapNm              string
	LabelUsgCllcnNmplc string
	LastUpdtDtm        string
	OrdKorNm           string
	OrdNm              string
	PrkNm              string
	TorsoLngth         string
	WingLngth          string
}
