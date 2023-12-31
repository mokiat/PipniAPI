package storage

import (
	"encoding/json"
	"io"
)

func SaveRegistry(out io.Writer, registry *RegistryDTO) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(registry)
}

func LoadRegistry(in io.Reader) (*RegistryDTO, error) {
	decoder := json.NewDecoder(in)
	var result RegistryDTO
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

type RegistryDTO struct {
	Folders         []FolderDTO   `json:"folders"`
	Contexts        []ContextDTO  `json:"contexts"`
	Endpoints       []EndpointDTO `json:"endpoints"`
	Workflows       []WorkflowDTO `json:"workflows"`
	ActiveContextID string        `json:"active_context_id"`
}

type FolderDTO struct {
	ID       string  `json:"id"`
	ParentID *string `json:"parent_id,omitempty"`
	Name     string  `json:"name"`
	Position int     `json:"position"`
}

type ContextDTO struct {
	ID         string        `json:"id"`
	FolderID   *string       `json:"folder_id,omitempty"`
	Name       string        `json:"name"`
	Position   int           `json:"position"`
	Properties []PropertyDTO `json:"properties"`
}

type PropertyDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EndpointDTO struct {
	ID       string      `json:"id"`
	FolderID *string     `json:"folder_id,omitempty"`
	Name     string      `json:"name"`
	Position int         `json:"position"`
	Method   string      `json:"method"`
	URI      string      `json:"uri"`
	Headers  []HeaderDTO `json:"headers"`
	Body     *string     `json:"body,omitempty"`
}

type HeaderDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WorkflowDTO struct {
	ID       string  `json:"id"`
	FolderID *string `json:"folder_id,omitempty"`
	Name     string  `json:"name"`
	Position int     `json:"position"`
}
