package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnSpcmInfoUseCase defines the lichen specimen detail use case.
type AlchnSpcmInfoUseCase interface {
	AlchnSpcmInfo(context.Context, application.AlchnSpcmInfoQuery) (application.AlchnSpcmInfoResult, error)
}
