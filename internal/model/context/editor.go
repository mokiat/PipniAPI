package context

import (
	"fmt"
	"slices"

	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, reg *registry.Model, cont *registry.Context) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		reg:      reg,
		cont:     cont,

		properties: slices.Clone(cont.Properties()),
	}
}

type Editor struct {
	eventBus *mvc.EventBus
	reg      *registry.Model
	cont     *registry.Context

	properties []gog.KV[string, string]
}

func (e *Editor) ID() string {
	return e.cont.ID()
}

func (e *Editor) Name() string {
	return e.cont.Name()
}

func (e *Editor) CanSave() bool {
	return !slices.Equal(e.properties, e.cont.Properties())
}

func (e *Editor) Save() error {
	e.cont.SetProperties(e.properties)
	if err := e.reg.Save(); err != nil {
		return fmt.Errorf("error saving registry: %w", err)
	}
	e.notifyModified()
	return nil
}

func (e *Editor) Properties() []gog.KV[string, string] {
	return slices.Clone(e.properties)
}

func (e *Editor) AddProperty() {
	e.properties = append(e.properties, gog.KV[string, string]{
		Key:   "",
		Value: "",
	})
	e.eventBus.Notify(PropertiesChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) SetPropertyName(index int, name string) {
	e.properties[index].Key = name
	e.eventBus.Notify(PropertiesChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) SetPropertyValue(index int, value string) {
	e.properties[index].Value = value
	e.eventBus.Notify(PropertiesChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) DeleteProperty(index int) {
	e.properties = slices.Delete(e.properties, index, index+1)
	e.eventBus.Notify(PropertiesChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) notifyModified() {
	e.eventBus.Notify(workspace.EditorModifiedEvent{
		Editor: e,
	})
}
