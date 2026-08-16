package application

// PlantSampleSearchQuery contains the plant sample search conditions.
type PlantSampleSearchQuery struct {
	PageNumber        int
	NumberOfRows      int
	RequestSearchWord string
}

// PlantSampleSearchResult contains a page of plant sample search results.
type PlantSampleSearchResult struct {
	Items        []PlantSampleSearchItem
	NumberOfRows int
	PageNumber   int
	TotalCount   int
}

// PlantSampleSearchItem contains one plant sample search result.
type PlantSampleSearchItem struct {
	Count                      int
	FamilyKoreanName           string
	FamilyName                 string
	PlantGeneralName           string
	PlantSpeciesID             string
	PlantSpeciesScientificName string
}
