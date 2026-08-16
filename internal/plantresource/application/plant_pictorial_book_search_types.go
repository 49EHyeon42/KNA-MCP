package application

// PlantPictorialBookSearchQuery contains the plant pictorial book search conditions.
type PlantPictorialBookSearchQuery struct {
	PageNumber        int
	NumberOfRows      int
	RequestSearchWord string
	DateFrom          string
	DateTo            string
}

// PlantPictorialBookSearchResult contains a page of plant pictorial book search results.
type PlantPictorialBookSearchResult struct {
	Items        []PlantPictorialBookSearchItem
	NumberOfRows int
	PageNumber   int
	TotalCount   int
}

// PlantPictorialBookSearchItem contains one plant pictorial book search result.
type PlantPictorialBookSearchItem struct {
	APGFamilyKoreanName        string
	APGFamilyName              string
	FamilyKoreanName           string
	FamilyName                 string
	GenusKoreanName            string
	GenusName                  string
	LastUpdateDateTime         string
	NotRecommendedGeneralName  string
	PlantGeneralName           string
	PlantPictorialBookNumber   string
	PlantSpeciesScientificName string
}
