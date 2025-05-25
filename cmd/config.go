package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of config.yaml
type Config struct {
    Headless         bool     `yaml:"headless"`
    SearchLevel      int      `yaml:"search_level"`
    ConnectionLevel  int      `yaml:"connection_level"`
    PerCompanyLimit  int      `yaml:"per_company_limit"`
    SearchList       []string `yaml:"search_list"`
    JobTitles        []string `yaml:"job_titles"`
    Blacklist        []string `yaml:"blacklist"`
}

func loadConfig(path string) (*Config, error) {
    // Read the file
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %v", err)
    }

    // Unmarshal into the Config struct
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %v", err)
    }

    return &config, nil
}

func appendToConfigList(path, listName, name string) error {
    // Read the current config
    config, err := loadConfig(path)
    if err != nil {
        return err
    }

    // Append the name to the specified list
    switch listName {
    case "blacklist":
        config.Blacklist = append(config.Blacklist, name)
    case "search_list":
        config.SearchList = append(config.SearchList, name)
    case "job_titles":
        config.JobTitles = append(config.JobTitles, name)
    default:
        return fmt.Errorf("unsupported list name: %s", listName)
    }

    // Write the updated config back to the file
    data, err := yaml.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %v", err)
    }
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        return fmt.Errorf("failed to write config: %v", err)
    }

    return nil
}