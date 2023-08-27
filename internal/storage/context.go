package storage

import (
	"encoding/json"
	"io"
)

func SaveContext(out io.Writer, registry *ContextDTO) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(registry)
}

func LoadContext(in io.Reader) (*ContextDTO, error) {
	decoder := json.NewDecoder(in)
	var result ContextDTO
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ContextDTO struct {
	Environments          []EnvironmentDTO `json:"environments"`
	SelectedEnvironmentID *string          `json:"selected_environment_id,omitempty"`
}

type EnvironmentDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
