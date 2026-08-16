package application

// PlantFolkAreaListQuery contains the folk plant area list conditions.
type PlantFolkAreaListQuery struct {
	PageNo    int
	NumOfRows int
	FlpltID   string
}

// PlantFolkAreaListResult contains a page of folk plant area information.
type PlantFolkAreaListResult struct {
	Items      []PlantFolkAreaListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantFolkAreaListItem contains one folk plant area record.
type PlantFolkAreaListItem struct {
	FlcstPlantExmnnAraTpcdNm string
	FlcstPlantLcltDscrt      string
	FlcstPlantPrpseDscrt     string
	FlpltID                  string
	PlantBrdgFomTpcdNm       string
	PlantGnrlNm              string
	PlantSpecsScnm           string
}
