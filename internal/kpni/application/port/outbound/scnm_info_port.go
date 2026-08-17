package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

// ScnmInfoPort defines the outbound port for scientific name details.
type ScnmInfoPort interface {
	ScnmInfo(context.Context, application.ScnmInfoQuery) (application.ScnmInfoResult, error)
}
