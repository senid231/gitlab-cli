// Package api provides access to gitlab HTTP API, YAML config and GIT repo info
package api

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConfigName is a name of yaml config file
const ConfigName = ".gitlab_cli.yml"

// Config holds options from config file
type Config struct {
	Token           string `yaml:"token"`
	URL             string `yaml:"url"`
	ProjectName     string `yaml:"project_name"`
	ForkProjectName string `yaml:"fork_project_name"`
}

// NewConfig reads config file from path and return Config object
func NewConfig(path *string) (*Config, error) {
	config := &Config{URL: "https://gitlab.com"}
	yamlPath := *path
	if yamlPath == "" {
		yamlPath = "./"
	}
	if !strings.HasSuffix(yamlPath, "/") {
		yamlPath += "/"
	}
	yamlPath += ConfigName
	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Printf("error read yaml\n%v\n", err)
		return config, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("error unmarshal yaml\n%v\n", err)
	}
	return config, err
}
