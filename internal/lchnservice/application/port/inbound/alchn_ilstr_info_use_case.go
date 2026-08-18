package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnIlstrInfoUseCase defines the lichen pictorial book detail use case.
type AlchnIlstrInfoUseCase interface {
	AlchnIlstrInfo(context.Context, application.AlchnIlstrInfoQuery) (application.AlchnIlstrInfoResult, error)
}
