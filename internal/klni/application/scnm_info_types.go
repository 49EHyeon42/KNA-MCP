package application

// ScnmInfoQuery contains the lichen scientific name detail condition.
type ScnmInfoQuery struct {
	ReqLchnScnmID string
}

// ScnmInfoResult contains lichen scientific name detail information.
type ScnmInfoResult struct {
	Item *ScnmInfoItem
}

// ScnmInfoItem contains one lichen scientific name detail record.
type ScnmInfoItem struct {
	StpltScnmRltnCdNm string
	FalmNm            string
	FalnKorNm         string
	GenusKorNm        string
	GenusNm           string
	LastUpdtDtm       string
	LchnGnrlNm        string
	LchnGnrlNm2       string
	LchnScnm          string
	LchnScnmID        string
	OrdscLtrtrNm      string
	Rmrk              string
}
