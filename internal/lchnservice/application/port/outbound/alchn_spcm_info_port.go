package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnSpcmInfoPort defines the outbound port for lichen specimen detail information.
type AlchnSpcmInfoPort interface {
	AlchnSpcmInfo(context.Context, application.AlchnSpcmInfoQuery) (application.AlchnSpcmInfoResult, error)
}
