package outbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPictorialBookInformationPort defines the outbound port for plant pictorial book information.
type PlantPictorialBookInformationPort interface {
	PlantPictorialBookInformation(context.Context, application.PlantPictorialBookInformationQuery) (application.PlantPictorialBookInformationResult, error)
}
