package context

import (
	"errors"
	"fmt"
	"os"

	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/lacking/ui/mvc"
)

var ErrContextNotFound = errors.New("context file missing")

func NewModel(eventBus *mvc.EventBus, cfgFileName string) *Model {
	return &Model{
		eventBus:    eventBus,
		cfgFileName: cfgFileName,
	}
}

type Model struct {
	eventBus    *mvc.EventBus
	cfgFileName string

	environments []*Environment
	selectedID   string
}

func (m *Model) Environments() []*Environment {
	return m.environments
}

func (m *Model) FindEnvironment(id string) *Environment {
	for _, env := range m.environments {
		if env.id == id {
			return env
		}
	}
	return nil
}

func (m *Model) SelectedEnvironment() *Environment {
	return m.FindEnvironment(m.selectedID)
}

func (m *Model) SelectedID() string {
	return m.selectedID
}

func (m *Model) SetSelectedID(id string) {
	if id != m.selectedID {
		m.selectedID = id
		m.eventBus.Notify(EnvironmentSelectedEvent{
			Model:       m,
			Environment: m.FindEnvironment(id),
		})
	}
}

func (m *Model) Load() error {
	file, err := os.Open(m.cfgFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrContextNotFound
		}
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	dtoContext, err := storage.LoadContext(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	m.loadFromDTO(dtoContext)

	m.eventBus.Notify(StructureChangedEvent{
		Model: m,
	})
	return nil
}

func (m *Model) Save() error {
	dtoContext := m.saveToDTO()

	file, err := os.Create(m.cfgFileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	err = storage.SaveContext(file, dtoContext)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func (m *Model) Clear() {
	m.environments = nil
	m.selectedID = ""
	m.eventBus.Notify(StructureChangedEvent{
		Model: m,
	})
}
