package application

// InsectSmplUnitListQuery contains the insect specimen detail list conditions.
type InsectSmplUnitListQuery struct {
	PageNo          int
	NumOfRows       int
	ReqInsctSpecsID string
}

// InsectSmplUnitListResult contains a page of insect specimen details.
type InsectSmplUnitListResult struct {
	Items      []InsectSmplUnitListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// InsectSmplUnitListItem contains one insect specimen detail.
type InsectSmplUnitListItem struct {
	BspcsInsttNm       string
	ClarHaslvVal       string
	SmplCllcnDt        string
	GynndTpcd          string
	HbttTpcd           string
	InsctSmplNo        string
	InsctSpecsID       string
	InsctSpecsScnm     string
	LabelUsgCllcnNmplc string
	LastUpdtDtm        string
	PrsrtStcd          string
	SlistTpcd          string
	SmplKindCd         string
	TorsoLngth         string
	WingLngth          string
	InsctGnrlNm        string
}
