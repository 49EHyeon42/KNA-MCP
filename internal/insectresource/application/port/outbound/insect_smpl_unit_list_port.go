package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectSmplUnitListPort defines the outbound port for insect specimen detail lists.
type InsectSmplUnitListPort interface {
	InsectSmplUnitList(context.Context, application.InsectSmplUnitListQuery) (application.InsectSmplUnitListResult, error)
}
