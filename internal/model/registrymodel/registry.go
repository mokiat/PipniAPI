package registrymodel

import (
	"errors"
	"fmt"
	"os"

	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/lacking/ui/mvc"
)

var ErrRegistryNotFound = errors.New("registry file missing")

func NewRegistry(eventBus *mvc.EventBus, cfgFileName string) *Registry {
	return &Registry{
		eventBus:    eventBus,
		cfgFileName: cfgFileName,

		root: &standardContainer{
			id:   RootContainerID,
			name: "Root",
		},
		selectedID: "",
	}
}

type Registry struct {
	eventBus    *mvc.EventBus
	cfgFileName string

	root       Container
	selectedID string
}

func (r *Registry) Root() Container {
	return r.root
}

func (r *Registry) SelectedID() string {
	return r.selectedID
}

func (r *Registry) SetSelectedID(selectedID string) {
	if selectedID != r.selectedID {
		r.selectedID = selectedID
		r.eventBus.Notify(RegistrySelectionChangedEvent{
			Registry:   r,
			SelectedID: selectedID,
		})
	}
}

func (r *Registry) Load() error {
	file, err := os.Open(r.cfgFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrRegistryNotFound
		}
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	dtoRegistry, err := storage.LoadRegistry(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	r.root = fromDTO(dtoRegistry)
	r.selectedID = ""
	r.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: r,
	})
	return nil
}

func (r *Registry) Save() error {
	// NOTE: Do this first in order to reduce failure chances.
	dtoRegistry := toDTO(r.root)

	file, err := os.Create(r.cfgFileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	err = storage.SaveRegistry(file, dtoRegistry)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func (r *Registry) Clear() {
	// IDEA: Consider defaulting to a working example.
	r.root = &standardContainer{
		id:   RootContainerID,
		name: "Root",
	}
	r.selectedID = ""
	r.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: r,
	})
}
