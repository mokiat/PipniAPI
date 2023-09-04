package context

import (
	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/gog"
	"golang.org/x/exp/slices"
)

func (m *Model) loadFromDTO(dtoContext *storage.ContextDTO) {
	m.environments = gog.Map(dtoContext.Environments, func(dtoEnvironment storage.EnvironmentDTO) *Environment {
		return &Environment{
			id:   dtoEnvironment.ID,
			name: dtoEnvironment.Name,
			// TODO: more fields
		}
	})
	selectedIndex := slices.IndexFunc(m.environments, func(env *Environment) bool {
		return env.id == *dtoContext.SelectedEnvironmentID
	})
	if selectedIndex >= 0 {
		m.selectedEnvironment = m.environments[0]
	} else {
		m.selectedEnvironment = nil
	}
}

func (m *Model) saveToDTO() *storage.ContextDTO {
	result := &storage.ContextDTO{}
	result.Environments = gog.Map(m.environments, func(environment *Environment) storage.EnvironmentDTO {
		return storage.EnvironmentDTO{
			ID:   environment.id,
			Name: environment.name,
		}
	})
	if m.selectedEnvironment != nil {
		result.SelectedEnvironmentID = gog.PtrOf(m.selectedEnvironment.ID())
	}
	return result
}
