package application

// ScnmInfoQuery contains the fungi scientific name detail condition.
type ScnmInfoQuery struct {
	ReqFngsScnmID string
}

// ScnmInfoResult contains fungi scientific name detail information.
type ScnmInfoResult struct {
	Item *ScnmInfoItem
}

// ScnmInfoItem contains one fungi scientific name detail record.
type ScnmInfoItem struct {
	StpltScnmRltnCdNm string
	FalmNm            string
	FalnKorNm         string
	FngsEclgTpcdNm    string
	FngsGnrlNm        string
	FngsGnrlNm2       string
	FngsPrpseTpcdNm   string
	FngsScnm          string
	FngsScnmID        string
	GenusKorNm        string
	GenusNm           string
	LastUpdtDtm       string
	OrdscLtrtrNm      string
	Rmrk              string
}
