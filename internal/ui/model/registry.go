package model

import (
	"github.com/mokiat/lacking/ui/mvc"
)

func NewRegistry(eventBus *mvc.EventBus) *Registry {
	return &Registry{
		eventBus: eventBus,

		environments: []*Environment{
			{
				name: "Staging",
			},
			{
				name: "Production",
			},
		},
	}
}

type Registry struct {
	eventBus *mvc.EventBus

	environments      []*Environment
	activeEnvironment *Environment
}

func (r *Registry) Environments() []*Environment {
	return r.environments
}

func (r *Registry) SetActiveEnvironment(env *Environment) {
	if env != r.activeEnvironment {
		r.activeEnvironment = env
		r.eventBus.Notify(ActiveEnvironmentChangedEvent{
			ActiveEnvironment: env,
		})
	}
}

func (r *Registry) ActiveEnvironment() *Environment {
	return r.activeEnvironment
}

func (r *Registry) RootFolders() []*Folder {
	return nil
}

type Environment struct {
	name string
	// settings map[string]string
}

func (e *Environment) Name() string {
	return e.name
}

type Folder struct {
}

func (f *Folder) Name() string {
	return ""
}

func (f *Folder) SubFolders() *Folder {
	return nil
}

func (f *Folder) Resources() []*Resource {
	return nil
}

type Resource struct {
}

func (r *Resource) Name() string {
	return ""
}

type ActiveEnvironmentChangedEvent struct {
	ActiveEnvironment *Environment
}
