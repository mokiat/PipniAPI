package registrymodel

type RegistrySelectionChangedEvent struct {
	Registry   *Registry
	SelectedID string
}

type RegistryStructureChangedEvent struct {
	Registry *Registry
}
