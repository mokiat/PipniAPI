package registry

import (
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

	for _, dtoContext := range dtoRegistry.Contexts {
		positions[dtoContext.ID] = dtoContext.Position
		root.resources = append(root.resources, &Context{
			id:        dtoContext.ID,
			name:      dtoContext.Name,
			container: root,

			properties: gog.Map(dtoContext.Properties, func(prop storage.PropertyDTO) gog.KV[string, string] {
				return gog.KV[string, string]{
					Key:   prop.Name,
					Value: prop.Value,
				}
			}),
		})
	}

	for _, dtoEndpoint := range dtoRegistry.Endpoints {
		positions[dtoEndpoint.ID] = dtoEndpoint.Position
		root.resources = append(root.resources, &Endpoint{
			id:        dtoEndpoint.ID,
			name:      dtoEndpoint.Name,
			container: root,

			method: dtoEndpoint.Method,
			uri:    dtoEndpoint.URI,
			headers: gog.Map(dtoEndpoint.Headers, func(header storage.HeaderDTO) gog.KV[string, string] {
				return gog.KV[string, string]{
					Key:   header.Name,
					Value: header.Value,
				}
			}),
			body: gog.ValueOf(dtoEndpoint.Body, ""),
		})
	}

	for _, dtoWorkflow := range dtoRegistry.Workflows {
		positions[dtoWorkflow.ID] = dtoWorkflow.Position
		root.resources = append(root.resources, &Workflow{
			id:        dtoWorkflow.ID,
			name:      dtoWorkflow.Name,
			container: root,
		})
	}

	slices.SortFunc(root.resources, func(a, b Resource) int {
		return positions[a.ID()] - positions[b.ID()]
	})

	m.root = root
	m.selectedID = ""
	m.activeContextID = dtoRegistry.ActiveContextID
}

func (m *Model) saveToDTO() *storage.RegistryDTO {
	result := &storage.RegistryDTO{
		ActiveContextID: m.activeContextID,
	}
	for i, resource := range m.root.Resources() {
		switch resource := resource.(type) {
		case *Context:
			result.Contexts = append(result.Contexts, storage.ContextDTO{
				ID:       resource.id,
				FolderID: nil, // Currently only root
				Name:     resource.name,
				Position: i,

				Properties: gog.Map(resource.properties, func(kv gog.KV[string, string]) storage.PropertyDTO {
					return storage.PropertyDTO{
						Name:  kv.Key,
						Value: kv.Value,
					}
				}),
			})
		case *Endpoint:
			result.Endpoints = append(result.Endpoints, storage.EndpointDTO{
				ID:       resource.id,
				FolderID: nil, // Currently only root
				Name:     resource.name,
				Position: i,

				Method: resource.method,
				URI:    resource.uri,
				Headers: gog.Map(resource.headers, func(kv gog.KV[string, string]) storage.HeaderDTO {
					return storage.HeaderDTO{
						Name:  kv.Key,
						Value: kv.Value,
					}
				}),
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
