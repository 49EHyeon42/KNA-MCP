package application

// PlantWordListQuery contains the plant word list conditions.
type PlantWordListQuery struct {
	PageNo       int
	NumOfRows    int
	ReqSearchWrd string
}

// PlantWordListResult contains a page of plant word information.
type PlantWordListResult struct {
	Items      []PlantWordListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// PlantWordListItem contains one plant word record.
type PlantWordListItem struct {
	EnglsWrdNm string
	KrnWrdNm   string
	PrfcnWrdNm string
	Wrddscrt   string
}
