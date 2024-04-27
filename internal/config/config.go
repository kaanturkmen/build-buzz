package config

import (
	"encoding/json"
	"os"
)

type EmailOverrides map[string]string

func LoadEmailOverrides(path string) (EmailOverrides, error) {
	var overrides EmailOverrides
	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &overrides)

	return overrides, err
}
