package util

import (
	"os"
	"gopkg.in/yaml.v3"
)

// GetYAML reads and parses a YAML file relative to the current file's directory.
func GetYAML(filePath string) (map[string]any, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    var result map[string]any
    err = yaml.Unmarshal(data, &result)
    return result, err
}