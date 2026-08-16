package application

// PlantSpecimenSearchQuery contains the plant specimen search conditions.
type PlantSpecimenSearchQuery struct {
	PageNumber        int
	NumberOfRows      int
	RequestSearchWord string
}

// PlantSpecimenSearchResult contains a page of plant specimen search results.
type PlantSpecimenSearchResult struct {
	Items        []PlantSpecimenSearchItem
	NumberOfRows int
	PageNumber   int
	TotalCount   int
}

// PlantSpecimenSearchItem contains one plant specimen search result.
type PlantSpecimenSearchItem struct {
	Count                      int
	FamilyKoreanName           string
	FamilyName                 string
	PlantGeneralName           string
	PlantSpeciesID             string
	PlantSpeciesScientificName string
}
