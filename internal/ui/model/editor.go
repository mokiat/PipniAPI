package model

type Editor interface {
	ID() string
	Title() string
	// IsDirty() bool
}
