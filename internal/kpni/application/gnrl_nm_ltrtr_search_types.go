package application

// GnrlNmLtrtrSearchQuery contains the plant general name literature list conditions.
type GnrlNmLtrtrSearchQuery struct {
	PageNo         int
	NumOfRows      int
	ReqPlantGnrlNm string
}

// GnrlNmLtrtrSearchResult contains a page of plant general name literature information.
type GnrlNmLtrtrSearchResult struct {
	Items      []GnrlNmLtrtrSearchItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// GnrlNmLtrtrSearchItem contains one plant general name literature record.
type GnrlNmLtrtrSearchItem struct {
	RcmmnTpcdNm      string
	LtrtrInfrmNm     string
	LvbngFrlngTpcdNm string
	PlantGnrlNm      string
	PlantSpecsScnm   string
}
