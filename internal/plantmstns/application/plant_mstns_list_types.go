package application

// PlantMstnsListQuery contains the plant miniature list conditions.
type PlantMstnsListQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
	ReqMnfctYr   string
}

// PlantMstnsListResult contains a page of plant miniature information.
type PlantMstnsListResult struct {
	Items      []PlantMstnsListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantMstnsListItem contains one plant miniature record.
type PlantMstnsListItem struct {
	DistrAraDscrt         string
	MinitrTpcdNm          string
	PlantBrdgFomTpcdNm    string
	PlantGnrlNm           string
	PlantMinitrAthrNm     string
	PlantMinitrMnfctMonth string
	PlantMinitrMnfctYr    string
	PlantMinitrPsinsNm    string
	PlantSpecsScnm        string
	RrnssPlantYn          string
	SpcltPlantYn          string
}
