package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectPrtctListUseCase defines the endangered insect list use case.
type InsectPrtctListUseCase interface {
	InsectPrtctList(context.Context, application.InsectPrtctListQuery) (application.InsectPrtctListResult, error)
}
