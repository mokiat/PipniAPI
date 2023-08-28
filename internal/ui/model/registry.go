package model

import (
	"errors"
	"fmt"
	"os"

	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui/mvc"
)

type Resource interface {
	ID() string
}

func NewRegistry(eventBus *mvc.EventBus, cfgFileName string) *Registry {
	return &Registry{
		eventBus:    eventBus,
		cfgFileName: cfgFileName,
	}
}

type Registry struct {
	eventBus    *mvc.EventBus
	cfgFileName string

	root       *Folder
	selectedID string
}

func (r *Registry) Load() error {
	r.reset()

	file, err := os.Open(r.cfgFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // start with clean registry
		}
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	dtoRegistry, err := storage.LoadRegistry(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	endpoints := gog.Map(dtoRegistry.Endpoints, func(dtoEndpoint storage.EndpointDTO) *Endpoint {
		return &Endpoint{
			id:   dtoEndpoint.ID,
			name: dtoEndpoint.Name,
			// TODO: more fields
		}
	})
	r.root.endpoints = endpoints

	workflows := gog.Map(dtoRegistry.Workflows, func(dtoWorkflow storage.WorkflowDTO) *Workflow {
		return &Workflow{
			id:   dtoWorkflow.ID,
			name: dtoWorkflow.Name,
			// TODO: more fields
		}
	})
	r.root.workflows = workflows

	// TODO: Notify changed

	return nil
}

func (r *Registry) Save() error {
	file, err := os.Create(r.cfgFileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// TODO

	return nil
}

func (r *Registry) SelectedResource() Resource {
	for _, endpoint := range r.root.endpoints {
		if endpoint.id == r.selectedID {
			return endpoint
		}
	}
	for _, workflow := range r.root.workflows {
		if workflow.id == r.selectedID {
			return workflow
		}
	}
	return nil
}

func (r *Registry) SelectedID() string {
	return r.selectedID
}

func (r *Registry) SetSelectedID(selectedID string) {
	if selectedID != r.selectedID {
		r.selectedID = selectedID
		r.eventBus.Notify(RegistrySelectionChangedEvent{
			SelectedID: selectedID,
		})
	}
}

func (r *Registry) Root() *Folder {
	return r.root
}

func (r *Registry) reset() {
	r.selectedID = ""
	r.root = &Folder{}
	// TODO: Notify so that all subscribers can nuke state
}

type Folder struct {
	id string

	endpoints []*Endpoint
	workflows []*Workflow
}

func (f *Folder) ID() string {
	return f.id
}

func (f *Folder) Endpoints() []*Endpoint {
	return f.endpoints
}

func (f *Folder) Workflows() []*Workflow {
	return f.workflows
}

type Endpoint struct {
	id   string
	name string
}

func (e *Endpoint) ID() string {
	return e.id
}

func (e *Endpoint) Name() string {
	return e.name
}

func (e *Endpoint) SetName(name string) {
	e.name = name
}

type Workflow struct {
	id   string
	name string
}

func (w *Workflow) ID() string {
	return w.id
}

func (w *Workflow) Name() string {
	return w.name
}

type RegistrySelectionChangedEvent struct {
	SelectedID string
}
