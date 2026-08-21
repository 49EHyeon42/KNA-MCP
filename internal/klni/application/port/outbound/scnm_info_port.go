package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
)

// ScnmInfoPort loads lichen scientific name detail information.
type ScnmInfoPort interface {
	ScnmInfo(context.Context, application.ScnmInfoQuery) (application.ScnmInfoResult, error)
}
