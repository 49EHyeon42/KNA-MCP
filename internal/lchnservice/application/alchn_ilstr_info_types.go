package application

// AlchnIlstrInfoQuery contains the lichen pictorial book detail condition.
type AlchnIlstrInfoQuery struct {
	Q1 string
}

// AlchnIlstrInfoResult contains lichen pictorial book detail information.
type AlchnIlstrInfoResult struct {
	Item *AlchnIlstrInfoItem
}

// AlchnIlstrInfoItem contains one lichen pictorial book detail record.
type AlchnIlstrInfoItem struct {
	Btnc         string
	Cont1        string
	Cont2        string
	Cont3        string
	Cont4        string
	Cont5        string
	Cont6        string
	Cont7        string
	Cont8        string
	Cont9        string
	Cont10       string
	Cont11       string
	Cont12       string
	CprtCtnt     string
	EngNm        string
	FamilyKorNm  string
	FamilyNm     string
	FrstRgstnDtm string
	GenusKorNm   string
	GenusNm      string
	ImgURL       string
	JapNm        string
	LastUpdtDtm  string
	LchnGnrlNm   string
	LchnInfrpNm  string
	LchnPilbkNo  string
	LchnScnmID   string
	LchnTtnm     string
	PrkNm        string
}
