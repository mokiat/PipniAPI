package fyneutil

import (
	"fmt"
	"io"
	"io/fs"
	"path"

	"fyne.io/fyne/v2"
)

func Must(res fyne.Resource, err error) fyne.Resource {
	if err != nil {
		panic(err)
	}
	return res
}

func LoadResourceFromFS(f fs.FS, resource string) (fyne.Resource, error) {
	file, err := f.Open(resource)
	if err != nil {
		return nil, fmt.Errorf("error opening fs resource: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading fs resource: %w", err)
	}

	return &fyne.StaticResource{
		StaticName:    path.Base(resource),
		StaticContent: content,
	}, nil
}
