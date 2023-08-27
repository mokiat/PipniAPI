package model

import (
	"errors"
	"fmt"
	"os"

	"github.com/mokiat/PipniAPI/internal/storage"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewContext(eventBus *mvc.EventBus, cfgFileName string) *Context {
	return &Context{
		eventBus:    eventBus,
		cfgFileName: cfgFileName,
	}
}

type Context struct {
	eventBus    *mvc.EventBus
	cfgFileName string

	environments []*Environment
	selectedID   string
}

func (c *Context) Load() error {
	c.reset()

	file, err := os.Open(c.cfgFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // start with clean registry
		}
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	dtoContext, err := storage.LoadContext(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	environments := gog.Map(dtoContext.Environments, func(dtoEnvironment storage.EnvironmentDTO) *Environment {
		return &Environment{
			id:   dtoEnvironment.ID,
			name: dtoEnvironment.Name,
			// TODO: more fields
		}
	})
	c.environments = environments

	c.selectedID = gog.ValueOf(dtoContext.SelectedEnvironmentID, "")

	// TODO: Notify changed

	return nil
}

func (c *Context) Save() error {
	file, err := os.Create(c.cfgFileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// TODO

	return nil
}

func (c *Context) Environments() []*Environment {
	return c.environments
}

func (c *Context) SelectedID() string {
	return c.selectedID
}

func (c *Context) SetSelectedID(selectedID string) {
	if selectedID != c.selectedID {
		c.selectedID = selectedID
		c.eventBus.Notify(ContextSelectionChangedEvent{
			SelectedID: selectedID,
		})
	}
}

func (c *Context) reset() {
	c.selectedID = ""
	c.environments = nil
}

type Environment struct {
	id   string
	name string
}

func (e *Environment) ID() string {
	return e.id
}

func (e *Environment) Name() string {
	return e.name
}

type ContextSelectionChangedEvent struct {
	SelectedID string
}
