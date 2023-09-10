package context

import (
	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/gog"
)

func (m *Model) loadFromDTO(dtoContext *storage.ContextDTO) {
	m.environments = gog.Map(dtoContext.Environments, func(dtoEnvironment storage.EnvironmentDTO) *Environment {
		return &Environment{
			id:   dtoEnvironment.ID,
			name: dtoEnvironment.Name,
			// TODO: more fields
		}
	})
	m.selectedID = gog.ValueOf(dtoContext.SelectedEnvironmentID, "")
}

func (m *Model) saveToDTO() *storage.ContextDTO {
	result := &storage.ContextDTO{}
	result.Environments = gog.Map(m.environments, func(environment *Environment) storage.EnvironmentDTO {
		return storage.EnvironmentDTO{
			ID:   environment.id,
			Name: environment.name,
		}
	})
	if m.selectedID != "" {
		result.SelectedEnvironmentID = &m.selectedID
	} else {
		result.SelectedEnvironmentID = nil
	}
	return result
}
