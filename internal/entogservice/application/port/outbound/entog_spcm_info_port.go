package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogSpcmInfoPort defines the outbound port for entognath specimen detail lookups.
type EntogSpcmInfoPort interface {
	EntogSpcmInfo(context.Context, application.EntogSpcmInfoQuery) (application.EntogSpcmInfoResult, error)
}
