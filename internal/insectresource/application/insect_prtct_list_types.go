package application

// InsectPrtctListQuery contains the endangered insect list pagination.
type InsectPrtctListQuery struct {
	PageNo    int
	NumOfRows int
}

// InsectPrtctListResult contains a page of endangered insects.
type InsectPrtctListResult struct {
	Items      []InsectPrtctListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// InsectPrtctListItem contains one endangered insect.
type InsectPrtctListItem struct {
	FamilyKorNm    string
	FamilyNm       string
	InsctGnrlNm    string
	InsctPcmtt     string
	InsctPilbkNo   string
	InsctSpecsScnm string
}
