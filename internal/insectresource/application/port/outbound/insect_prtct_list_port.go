package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectPrtctListPort loads endangered insects.
type InsectPrtctListPort interface {
	InsectPrtctList(context.Context, application.InsectPrtctListQuery) (application.InsectPrtctListResult, error)
}
