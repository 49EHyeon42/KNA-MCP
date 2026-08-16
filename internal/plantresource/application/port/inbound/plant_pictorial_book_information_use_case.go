package inbound

import (
	"context"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

// PlantPictorialBookInformationUseCase defines the plant pictorial book information use case.
type PlantPictorialBookInformationUseCase interface {
	PlantPictorialBookInformation(context.Context, application.PlantPictorialBookInformationQuery) (application.PlantPictorialBookInformationResult, error)
}
