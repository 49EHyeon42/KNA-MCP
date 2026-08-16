package application

// PlantPilbkInfoQuery contains the plant pictorial book information request.
type PlantPilbkInfoQuery struct {
	ReqPlantPilbkNo string
}

// PlantPilbkInfoResult contains plant pictorial book information.
type PlantPilbkInfoResult struct {
	APGFamilyKorNm string
	APGFamilyNm    string
	BfofMthod      string
	BrdMthdDesc    string
	BugInfo        string
	Dstrb          string
	EngNm          string
	FamilyKorNm    string
	FamilyNm       string
	FarmSpftDesc   string
	GenusKorNm     string
	GenusNm        string
	GrwEvrntDesc   string
	InductionDesc  string
	LastUpdtDtm    string
	NotRcmmGnrlNm  string
	Note           string
	OrplcNm        string
	OsDstrb        string
	PlantGnrlNm    string
	PlantPilbkNo   string
	PlantSpecsScnm string
	PrtcPlnDesc    string
	RrngGubun      string
	RrngType       string
	Shpe           string
	SmlrPlntDesc   string
	Spft           string
	UseMthdDesc    string
	WoodDesc       string
}
