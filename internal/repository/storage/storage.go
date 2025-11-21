// Package storage provides abstractions for work with filesystem
package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadFile(file *os.File, data any) error {
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("can not seek to file beginning: %w", err)
	}

	dec := json.NewDecoder(file)
	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("can not unmarshal data: %w", err)
	}
	return nil
}

func UpdateFile(file *os.File, data any) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("can not truncate file: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("can not seek to file beginning: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("can not marshal data: %w", err)
	}

	return nil
}
