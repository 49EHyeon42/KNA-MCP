package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

// ChildPilbkInfoPort defines the outbound port for child pictorial book detail queries.
type ChildPilbkInfoPort interface {
	ChildPilbkInfo(context.Context, application.ChildPilbkInfoQuery) (application.ChildPilbkInfoResult, error)
}
