package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

// ChildPilbkInfoUseCase defines the child pictorial book detail use case.
type ChildPilbkInfoUseCase interface {
	ChildPilbkInfo(context.Context, application.ChildPilbkInfoQuery) (application.ChildPilbkInfoResult, error)
}
