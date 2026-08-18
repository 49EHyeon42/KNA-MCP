package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsPilbkInfoPort defines the outbound port for fungi pictorial book detail queries.
type FngsPilbkInfoPort interface {
	FngsPilbkInfo(context.Context, application.FngsPilbkInfoQuery) (application.FngsPilbkInfoResult, error)
}
