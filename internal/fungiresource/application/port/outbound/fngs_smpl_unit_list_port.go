package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsSmplUnitListPort defines the outbound port for fungi specimen detail lists.
type FngsSmplUnitListPort interface {
	FngsSmplUnitList(context.Context, application.FngsSmplUnitListQuery) (application.FngsSmplUnitListResult, error)
}
