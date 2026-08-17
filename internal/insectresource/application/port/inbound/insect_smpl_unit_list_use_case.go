package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectSmplUnitListUseCase defines the insect specimen detail list use case.
type InsectSmplUnitListUseCase interface {
	InsectSmplUnitList(context.Context, application.InsectSmplUnitListQuery) (application.InsectSmplUnitListResult, error)
}
