package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectPilbkInfoUseCase defines the insect pictorial book detail use case.
type InsectPilbkInfoUseCase interface {
	InsectPilbkInfo(context.Context, application.InsectPilbkInfoQuery) (application.InsectPilbkInfoResult, error)
}
