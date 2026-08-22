package application

// EntogIlstrInfoQuery contains the entognath pictorial book detail lookup key.
type EntogIlstrInfoQuery struct {
	Q1 string
}

// EntogIlstrInfoResult contains entognath pictorial book detail information.
type EntogIlstrInfoResult struct {
	Item *EntogIlstrInfoItem
}

// EntogIlstrInfoItem contains one entognath pictorial book detail result.
type EntogIlstrInfoItem struct {
	Btnc             string
	Cont1            string
	Cont2            string
	Cont3            string
	Cont4            string
	Cont5            string
	Cont6            string
	Cont7            string
	Cont8            string
	Cont9            string
	Cont10           string
	Cont11           string
	CprtCtnt         string
	EmrgcCnt         string
	EmrgcEraDscrt    string
	EntogAthrNm      string
	EntogEngNm       string
	EntogOfnmKrlngNm string
	EntogPilbkNo     string
	EntogSpecsNm     string
	FamilyKorNm      string
	FamilyNm         string
	GenusKorNm       string
	GenusNm          string
	ImgURL           string
	MnmmOccrrCnt     string
	MxmmOccrrCnt     string
	OrdKorNm         string
	OrdNm            string
}
