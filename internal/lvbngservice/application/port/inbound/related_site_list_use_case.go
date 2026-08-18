package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
)

// RelatedSiteListUseCase defines the related site list use case.
type RelatedSiteListUseCase interface {
	RelatedSiteList(context.Context, application.RelatedSiteListQuery) (application.RelatedSiteListResult, error)
}
