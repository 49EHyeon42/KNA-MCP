package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
)

// ScnmInfoUseCase defines the insect scientific name detail use case.
type ScnmInfoUseCase interface {
	ScnmInfo(context.Context, application.ScnmInfoQuery) (application.ScnmInfoResult, error)
}
