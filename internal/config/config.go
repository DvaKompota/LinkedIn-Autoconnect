package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of config.yaml
type Config struct {
    path             string   `yaml:"-"` // not marshaled to YAML
    Headless         bool     `yaml:"headless"`
    SearchLevel      int      `yaml:"search_level"`
    ConnectionLevel  int      `yaml:"connection_level"`
    PerCompanyLimit  int      `yaml:"per_company_limit"`
    SearchList       []string `yaml:"search_list"`
    JobTitles        []string `yaml:"job_titles"`
    Blacklist        []string `yaml:"blacklist"`
}

func LoadConfig(path string) (*Config, error) {
    // Read the file
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %v", err)
    }

    // Unmarshal into the Config struct
    var cfg Config
    err = yaml.Unmarshal(data, &cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %v", err)
    }
    cfg.path = path
    return &cfg, nil
}

func (cfg *Config) AppendToList(listName, name string) error {
    // Append the name to the specified list
    switch listName {
    case "blacklist":
        cfg.Blacklist = append(cfg.Blacklist, name)
    case "search_list":
        cfg.SearchList = append(cfg.SearchList, name)
    case "job_titles":
        cfg.JobTitles = append(cfg.JobTitles, name)
    default:
        return fmt.Errorf("unsupported list name: %s", listName)
    }

    // Write the updated config back to the file
    data, err := yaml.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %v", err)
    }
    err = os.WriteFile(cfg.path, data, 0644)
    if err != nil {
        return fmt.Errorf("failed to write config: %v", err)
    }
    return nil
}