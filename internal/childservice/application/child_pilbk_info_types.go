package application

// ChildPilbkInfoQuery contains the child pictorial book detail condition.
type ChildPilbkInfoQuery struct {
	ReqChildLvbngPilbkNo string
}

// ChildPilbkInfoResult contains child pictorial book detail information.
type ChildPilbkInfoResult struct {
	Item *ChildPilbkInfoItem
}

// ChildPilbkInfoItem contains one child pictorial book detail record.
type ChildPilbkInfoItem struct {
	BiogyNm           string
	ChildLvbngPilbkNo string
	ExtrmCrss         string
	FamilyKorNm       string
	FamilyNm          string
	GenusKorNm        string
	GenusNm           string
	HbttFieldYn       string
	HbttFrestYn       string
	HbttRiverYn       string
	LvbngDscrt        string
	LvbngTpcdNm       string
	LvngKrlngNm       string
	PrtctSpecsTpcdNm  string
}
