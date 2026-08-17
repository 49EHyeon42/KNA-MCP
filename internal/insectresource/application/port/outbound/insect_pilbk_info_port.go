package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

// InsectPilbkInfoPort defines the outbound port for insect pictorial book details.
type InsectPilbkInfoPort interface {
	InsectPilbkInfo(context.Context, application.InsectPilbkInfoQuery) (application.InsectPilbkInfoResult, error)
}
