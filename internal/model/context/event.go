package context

type StructureChangedEvent struct {
	Model *Model
}

type EnvironmentSelectedEvent struct {
	Model       *Model
	Environment *Environment
}
