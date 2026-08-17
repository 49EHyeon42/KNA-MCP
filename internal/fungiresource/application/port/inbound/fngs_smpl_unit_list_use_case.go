package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsSmplUnitListUseCase defines the fungi specimen detail list use case.
type FngsSmplUnitListUseCase interface {
	FngsSmplUnitList(context.Context, application.FngsSmplUnitListQuery) (application.FngsSmplUnitListResult, error)
}
