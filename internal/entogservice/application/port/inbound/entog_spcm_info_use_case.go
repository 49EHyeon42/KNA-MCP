package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogSpcmInfoUseCase defines the inbound port for entognath specimen detail lookups.
type EntogSpcmInfoUseCase interface {
	EntogSpcmInfo(context.Context, application.EntogSpcmInfoQuery) (application.EntogSpcmInfoResult, error)
}
