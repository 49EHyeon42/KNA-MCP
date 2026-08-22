package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogIlstrInfoUseCase defines the entognath pictorial book detail use case.
type EntogIlstrInfoUseCase interface {
	EntogIlstrInfo(context.Context, application.EntogIlstrInfoQuery) (application.EntogIlstrInfoResult, error)
}
