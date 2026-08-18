package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

// FngsPilbkInfoUseCase defines the fungi pictorial book detail use case.
type FngsPilbkInfoUseCase interface {
	FngsPilbkInfo(context.Context, application.FngsPilbkInfoQuery) (application.FngsPilbkInfoResult, error)
}
