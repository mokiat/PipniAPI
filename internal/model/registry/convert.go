package registry

import (
	"net/http"

	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/gog"
	"golang.org/x/exp/slices"
)

func (m *Model) loadFromDTO(dtoRegistry *storage.RegistryDTO) {
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

			method: dtoEndpoint.Method,
			uri:    dtoEndpoint.URI,
			headers: func() http.Header {
				result := make(http.Header)
				for _, header := range dtoEndpoint.Headers {
					result.Add(header.Name, header.Value)
				}
				return result
			}(),
			body: gog.ValueOf(dtoEndpoint.Body, ""),
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

	m.root = root
	m.selectedID = ""
}

func (m *Model) saveToDTO() *storage.RegistryDTO {
	result := &storage.RegistryDTO{}
	for i, resource := range m.root.Resources() {
		switch resource := resource.(type) {
		case *Endpoint:
			result.Endpoints = append(result.Endpoints, storage.EndpointDTO{
				ID:       resource.id,
				FolderID: nil, // Currently only root
				Name:     resource.name,
				Position: i,

				Method: resource.method,
				URI:    resource.uri,
				Headers: func() []storage.HeaderDTO {
					var result []storage.HeaderDTO
					for name, values := range resource.headers {
						for _, value := range values {
							result = append(result, storage.HeaderDTO{
								Name:  name,
								Value: value,
							})
						}
					}
					return result
				}(),
				Body: &resource.body,
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
