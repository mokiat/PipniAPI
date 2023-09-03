package registrymodel

import (
	"github.com/mokiat/PipniAPI/internal/storage"
	"golang.org/x/exp/slices"
)

func fromDTO(dtoRegistry *storage.RegistryDTO) Container {
	positions := make(map[string]int)

	root := &standardContainer{
		id:        RootContainerID,
		name:      "Root",
		children:  nil,
		resources: make([]Resource, 0, len(dtoRegistry.Endpoints)+len(dtoRegistry.Workflows)),
	}

	for _, dtoEndpoint := range dtoRegistry.Endpoints {
		positions[dtoEndpoint.ID] = dtoEndpoint.Position
		root.resources = append(root.resources, &Endpoint{
			id:        dtoEndpoint.ID,
			name:      dtoEndpoint.Name,
			container: root,
			// TODO: more fields
		})
	}

	for _, dtoWorkflow := range dtoRegistry.Workflows {
		positions[dtoWorkflow.ID] = dtoWorkflow.Position
		root.resources = append(root.resources, &Workflow{
			id:        dtoWorkflow.ID,
			name:      dtoWorkflow.Name,
			container: root,
			// TODO: more fields
		})
	}

	slices.SortFunc(root.resources, func(a, b Resource) int {
		return positions[a.ID()] - positions[b.ID()]
	})

	return root
}

func toDTO(root Container) *storage.RegistryDTO {
	result := &storage.RegistryDTO{}
	for i, resource := range root.Resources() {
		switch resource := resource.(type) {
		case *Endpoint:
			result.Endpoints = append(result.Endpoints, storage.EndpointDTO{
				ID:       resource.id,
				FolderID: nil, // Currently only root
				Name:     resource.name,
				Position: i,
				// TODO: More params
			})
		case *Workflow:
			result.Workflows = append(result.Workflows, storage.WorkflowDTO{
				ID:       resource.id,
				FolderID: nil, // Currently only root
				Name:     resource.name,
				Position: i,
			})
		}
	}
	return result
}
