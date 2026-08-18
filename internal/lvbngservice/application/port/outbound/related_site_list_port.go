package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
)

// RelatedSiteListPort defines the outbound port for related site lists.
type RelatedSiteListPort interface {
	RelatedSiteList(context.Context, application.RelatedSiteListQuery) (application.RelatedSiteListResult, error)
}
