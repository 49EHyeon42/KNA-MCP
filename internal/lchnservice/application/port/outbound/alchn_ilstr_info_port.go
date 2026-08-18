package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

// AlchnIlstrInfoPort defines the outbound port for lichen pictorial book detail information.
type AlchnIlstrInfoPort interface {
	AlchnIlstrInfo(context.Context, application.AlchnIlstrInfoQuery) (application.AlchnIlstrInfoResult, error)
}
