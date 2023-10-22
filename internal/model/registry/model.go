package registry

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui/mvc"
)

var ErrRegistryNotFound = errors.New("registry file missing")

func NewModel(eventBus *mvc.EventBus, cfgFileName string) *Model {
	return &Model{
		eventBus:    eventBus,
		cfgFileName: cfgFileName,

		root: &standardContainer{
			id:   RootContainerID,
			name: "Root",
		},
		selectedID: "",
	}
}

type Model struct {
	eventBus    *mvc.EventBus
	cfgFileName string

	root            Container
	selectedID      string
	activeContextID string
}

func (m *Model) Root() Container {
	return m.root
}

func (m *Model) AllResources() []Resource {
	return m.root.AllResources()
}

func (m *Model) ActiveContext() *Context {
	contextResources := gog.Select(m.AllResources(), func(resource Resource) bool {
		_, ok := resource.(*Context)
		return ok
	})
	contexts := gog.Map(contextResources, func(resource Resource) *Context {
		return resource.(*Context)
	})
	result, ok := gog.FindFunc(contexts, func(context *Context) bool {
		return context.ID() == m.activeContextID
	})
	if !ok {
		return nil
	}
	return result
}

func (m *Model) ActiveContextID() string {
	return m.activeContextID
}

func (m *Model) SetActiveContextID(activeID string) {
	if activeID != m.activeContextID {
		m.activeContextID = activeID
		m.eventBus.Notify(RegistryActiveContextChangedEvent{
			Registry: m,
		})
	}
}

func (m *Model) SelectedID() string {
	return m.selectedID
}

func (m *Model) SetSelectedID(selectedID string) {
	if selectedID != m.selectedID {
		m.selectedID = selectedID
		m.eventBus.Notify(RegistrySelectionChangedEvent{
			Registry:   m,
			SelectedID: selectedID,
		})
	}
}

func (m *Model) FindSelectedResource() Resource {
	return m.root.FindResource(m.selectedID)
}

func (m *Model) CanMoveUp(resource Resource) bool {
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

func (m *Model) MoveUp(resource Resource) {
	if resource == nil {
		return
	}
	container := resource.Container()
	container.MoveResourceUp(resource)
	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
}

func (m *Model) CanMoveDown(resource Resource) bool {
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

func (m *Model) MoveDown(resource Resource) {
	if resource == nil {
		return
	}
	container := resource.Container()
	container.MoveResourceDown(resource)
	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
}

func (m *Model) CreateResource(parent Container, name string, kind ResourceKind) {
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

	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
	m.SetSelectedID(resource.ID())
}

func (m *Model) RenameResource(resource Resource, name string) {
	resource.SetName(name)
	m.eventBus.Notify(RegistryResourceNameChangedEvent{
		Registry: m,
		Resource: resource,
	})
}

func (m *Model) CloneResource(resource Resource) {
	parent := resource.Container()
	newResource := resource.Clone()
	parent.AppendResource(newResource)
	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
}

func (m *Model) DeleteResource(resource Resource) {
	parent := resource.Container()
	parent.RemoveResource(resource)
	m.eventBus.Notify(RegistryResourceRemovedEvent{
		Registry: m,
		Resource: resource,
	})
	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
}

func (m *Model) Load() error {
	file, err := os.Open(m.cfgFileName)
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

	m.loadFromDTO(dtoRegistry)

	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
	return nil
}

func (m *Model) Save() error {
	dtoRegistry := m.saveToDTO()

	file, err := os.Create(m.cfgFileName)
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

func (m *Model) Clear() {
	m.root = &standardContainer{
		id:   RootContainerID,
		name: "Root",
	}

	contextID := uuid.NewString()
	m.root.AppendResource(&Context{
		id:        contextID,
		name:      "REQRES",
		container: m.root,
		properties: []gog.KV[string, string]{
			{
				Key:   "host",
				Value: "reqres.in",
			},
		},
	})

	m.root.AppendResource(&Endpoint{
		id:     uuid.NewString(),
		name:   "Create User",
		method: http.MethodPost,
		uri:    "https://{{.host}}/api/users",
		headers: []gog.KV[string, string]{
			{
				Key:   "Content-Type",
				Value: "application/json",
			},
		},
		container: m.root,
		body:      "{\"name\":\"morpheus\",\"job\":\"leader\"}",
	})

	m.root.AppendResource(&Endpoint{
		id:        uuid.NewString(),
		name:      "Get User",
		method:    http.MethodGet,
		uri:       "https://{{.host}}/api/users/2",
		headers:   []gog.KV[string, string]{},
		container: m.root,
		body:      "",
	})

	m.root.AppendResource(&Endpoint{
		id:        uuid.NewString(),
		name:      "List Users",
		method:    http.MethodGet,
		uri:       "https://{{.host}}/api/users",
		headers:   []gog.KV[string, string]{},
		container: m.root,
		body:      "",
	})

	m.root.AppendResource(&Endpoint{
		id:        uuid.NewString(),
		name:      "Delete User",
		method:    http.MethodDelete,
		uri:       "https://{{.host}}/api/users/2",
		headers:   []gog.KV[string, string]{},
		container: m.root,
		body:      "",
	})

	m.activeContextID = contextID
	m.selectedID = ""

	m.eventBus.Notify(RegistryStructureChangedEvent{
		Registry: m,
	})
}
