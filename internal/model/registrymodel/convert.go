package registrymodel

import (
	"github.com/mokiat/PipniAPI/internal/storage"
	"golang.org/x/exp/slices"
)

func fromDTO(dtoRegistry *storage.RegistryDTO) Container {
	positions := make(map[string]int)
	resources := make([]Resource, 0, len(dtoRegistry.Endpoints)+len(dtoRegistry.Workflows))

	for _, dtoEndpoint := range dtoRegistry.Endpoints {
		positions[dtoEndpoint.ID] = dtoEndpoint.Position
		resources = append(resources, &Endpoint{
			id:   dtoEndpoint.ID,
			name: dtoEndpoint.Name,
			// TODO: more fields
		})
	}

	for _, dtoWorkflow := range dtoRegistry.Workflows {
		positions[dtoWorkflow.ID] = dtoWorkflow.Position
		resources = append(resources, &Workflow{
			id:   dtoWorkflow.ID,
			name: dtoWorkflow.Name,
			// TODO: more fields
		})
	}

	slices.SortFunc(resources, func(a, b Resource) int {
		return positions[a.ID()] - positions[b.ID()]
	})

	return &standardContainer{
		id:        RootContainerID,
		name:      "Root",
		children:  nil,
		resources: resources,
	}
}

func toDTO(root Container) *storage.RegistryDTO {
	panic("TODO")
}
