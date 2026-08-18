package application

// RelatedSiteListQuery contains the related site list conditions.
type RelatedSiteListQuery struct {
	PageNo    int
	NumOfRows int
}

// RelatedSiteListResult contains a page of related site information.
type RelatedSiteListResult struct {
	Items      []RelatedSiteListItem
	NumOfRows  int
	PageNo     int
	TotalCount int
}

// RelatedSiteListItem contains one related site record.
type RelatedSiteListItem struct {
	LvbngTpcdNm string
	SiteCtgryNm string
	SiteNm      string
	SiteURL     string
}
