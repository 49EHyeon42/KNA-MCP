package application

// ScnmInfoQuery contains the insect scientific name detail condition.
type ScnmInfoQuery struct {
	ReqInsctScnmID string
}

// ScnmInfoResult contains insect scientific name detail information.
type ScnmInfoResult struct {
	Item *ScnmInfoItem
}

// ScnmInfoItem contains one insect scientific name detail record.
type ScnmInfoItem struct {
	SuperFalmNm       string
	ClassKorNm        string
	ClassNm           string
	FalmKorNm         string
	FalmNm            string
	GenusKorNm        string
	GenusNm           string
	InsctGnrlNm       string
	InsctGnrlNm2      string
	InsctScnmID       string
	InsctSpecsScnm    string
	LastUpdtDtm       string
	OrdKorNm          string
	OrdNm             string
	StpltScnmRltnCdNm string
	SubFalmKorNm      string
	SubFalmNm         string
	SuperFalmKorNm    string
}
