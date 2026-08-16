package application

// PlantPictorialBookInformationQuery contains the plant pictorial book information request.
type PlantPictorialBookInformationQuery struct {
	RequestPlantPictorialBookNumber string
}

// PlantPictorialBookInformationResult contains plant pictorial book information.
type PlantPictorialBookInformationResult struct {
	APGFamilyKoreanName           string
	APGFamilyName                 string
	BfofMethod                    string
	BreedingMethodDescription     string
	BugInformation                string
	Distribution                  string
	EnglishName                   string
	FamilyKoreanName              string
	FamilyName                    string
	FarmSpecialFeatureDescription string
	GenusKoreanName               string
	GenusName                     string
	GrowthEnvironmentDescription  string
	InductionDescription          string
	LastUpdateDateTime            string
	NotRecommendedGeneralName     string
	Note                          string
	OriginPlaceName               string
	OverseasDistribution          string
	PlantGeneralName              string
	PlantPictorialBookNumber      string
	PlantSpeciesScientificName    string
	ProtectionPlanDescription     string
	RearingClassification         string
	RearingType                   string
	Shape                         string
	SimilarPlantDescription       string
	SpecialFeature                string
	UseMethodDescription          string
	WoodDescription               string
}
