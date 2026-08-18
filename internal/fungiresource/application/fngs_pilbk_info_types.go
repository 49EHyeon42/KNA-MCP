package application

// FngsPilbkInfoQuery contains the fungi pictorial book detail condition.
type FngsPilbkInfoQuery struct {
	ReqFngsPilbkNo string
}

// FngsPilbkInfoResult contains fungi pictorial book detail information.
type FngsPilbkInfoResult struct {
	Item *FngsPilbkInfoItem
}

// FngsPilbkInfoItem contains one fungi pictorial book detail record.
type FngsPilbkInfoItem struct {
	MshrmColorCdNm      string
	CrpphFomTpcdNm      string
	FamilyKorNm         string
	FamilyNm            string
	FngsEclgTpcdNm      string
	FngsGnrlNm          string
	FngsPilbkNo         string
	FngsPrpseTpcdNm     string
	FngsScnm            string
	GenusKorNm          string
	GenusNm             string
	GrwEvrntDesc        string
	LastUpdtDtm         string
	MicroShpe           string
	MshrmTpcdNm         string
	OccrrSsnNm          string
	RsrcActoClsscTpcdNm string
	Shpe                string
}
