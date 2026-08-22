package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

// EntogIlstrInfoPort defines the outbound port for entognath pictorial book details.
type EntogIlstrInfoPort interface {
	EntogIlstrInfo(context.Context, application.EntogIlstrInfoQuery) (application.EntogIlstrInfoResult, error)
}
