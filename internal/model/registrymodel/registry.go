package registrymodel

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/lacking/log"
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

func (r *Registry) FindSelectedResource() Resource {
	return r.root.FindResource(r.selectedID)
}

func (r *Registry) CanMoveUp(resource Resource) bool {
	if resource == nil {
		return false
	}
	// TODO: Once there is nesting of containers, this should be more
	// complicated and consider the whole tree and not just the container
	// (e.g. can the resource move to an upper sibling container)
	container := resource.Container()
	position := container.ResourcePosition(resource)
	return position > 0
}

func (r *Registry) MoveUp(resource Resource) {
	if resource == nil {
		return
	}
	container := resource.Container()
	container.MoveResourceUp(resource)
	r.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: r,
	})
}

func (r *Registry) CanMoveDown(resource Resource) bool {
	if resource == nil {
		return false
	}
	// TODO: Once there is nesting of containers, this should be more
	// complicated and consider the whole tree and not just the container
	// (e.g. can the resource move to a lower sibling container)
	container := resource.Container()
	position := container.ResourcePosition(resource)
	return position < len(container.Resources())-1
}

func (r *Registry) MoveDown(resource Resource) {
	if resource == nil {
		return
	}
	container := resource.Container()
	container.MoveResourceDown(resource)
	r.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: r,
	})
}

func (r *Registry) CreateResource(parent Container, name string, kind ResourceKind) {
	var resource Resource
	switch kind {
	case ResourceKindEndpoint:
		resource = &Endpoint{
			container: parent,
			id:        uuid.Must(uuid.NewRandom()).String(),
			name:      name,
			// TODO: Initialize more props. Or maybe consider a newEndpoint function
		}
	case ResourceKindWorkflow:
		resource = &Workflow{
			container: parent,
			id:        uuid.Must(uuid.NewRandom()).String(),
			name:      name,
			// TODO: Initialize more props. Or maybe consider a newWorkflow function
		}
	default:
		log.Warn("Unknown resource kind %q", kind)
		return
	}
	parent.AppendResource(resource)

	r.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: r,
	})
	r.SetSelectedID(resource.ID())
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
