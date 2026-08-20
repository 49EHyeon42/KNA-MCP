package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/kfni/application"
)

// ScnmInfoPort loads fungi scientific name detail information.
type ScnmInfoPort interface {
	ScnmInfo(context.Context, application.ScnmInfoQuery) (application.ScnmInfoResult, error)
}
