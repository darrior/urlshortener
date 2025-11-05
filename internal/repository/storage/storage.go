// Package storage provides abstractions for work with filesystem
package storage

import (
	"encoding/json"
	"os"
)

func ReadFile(file *os.File, data any) error {
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	dec := json.NewDecoder(file)
	if err := dec.Decode(data); err != nil {
		return err
	}
	return nil
}

func UpdateFile(file *os.File, data any) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}
