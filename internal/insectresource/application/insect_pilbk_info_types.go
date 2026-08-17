package application

// InsectPilbkInfoQuery contains the insect pictorial book detail condition.
type InsectPilbkInfoQuery struct {
	ReqInsctPilbkNo string
}

// InsectPilbkInfoResult contains insect pictorial book detail information.
type InsectPilbkInfoResult struct {
	Item *InsectPilbkInfoItem
}

// InsectPilbkInfoItem contains one insect pictorial book detail record.
type InsectPilbkInfoItem struct {
	EcoDsrct         string
	EggDsrct         string
	EmrgcCnt         string
	EmrgcEraDscrt    string
	FamilyKorNm      string
	FamilyNm         string
	FemaleDsrct      string
	GenusKorNm       string
	GenusNm          string
	GnrlDsrct        string
	HabitDsrct       string
	InsctEngNm       string
	InsctGnrlNm      string
	InsctPilbkNo     string
	InsctSpecsScnm   string
	LarvaDsrct       string
	LastUpdtDtm      string
	MaleDsrct        string
	MnmmOccrrCnt     string
	MxmmOccrrCnt     string
	OrdKorNm         string
	OrdNm            string
	PestDsrct        string
	PupaDsrct        string
	ReferDsrct       string
	SubFamilyKorNm   string
	SubFamilyNm      string
	SuperFamilyKorNm string
	SuperFamilyNm    string
	WinterDsrct      string
}
