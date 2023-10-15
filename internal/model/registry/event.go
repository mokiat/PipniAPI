package registry

type RegistryActiveContextChangedEvent struct {
	Registry *Model
}

type RegistrySelectionChangedEvent struct {
	Registry   *Model
	SelectedID string
}

type RegistryStructureChangedEvent struct {
	Registry *Model
}

type RegistryResourceNameChangedEvent struct {
	Registry *Model
	Resource Resource
}

type RegistryResourceRemovedEvent struct {
	Registry *Model
	Resource Resource
}
