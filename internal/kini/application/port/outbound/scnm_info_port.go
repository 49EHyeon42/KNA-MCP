package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
)

// ScnmInfoPort defines the outbound port for insect scientific name details.
type ScnmInfoPort interface {
	ScnmInfo(context.Context, application.ScnmInfoQuery) (application.ScnmInfoResult, error)
}
