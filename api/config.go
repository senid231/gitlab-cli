package api

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

const CONFIG_NAME = ".gitlab_cli.yml"

type Config struct {
	Token       string `yaml:"token"`
	Url         string `yaml:"url"`
	ProjectName string `yaml:"project_name"`
}

func NewConfig(path *string) (config *Config, err error) {
	fmt.Println("NewConfig")
	config = &Config{Url: "https://gitlab.com"}
	yamlPath := *path
	if yamlPath == "" {
		yamlPath = "./"
	}
	if !strings.HasSuffix(yamlPath, "/") {
		yamlPath += "/"
	}
	yamlPath += CONFIG_NAME
	fmt.Printf("gonna read yaml from %s", yamlPath)
	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Printf("error read yaml\n%v\n", err)
		return
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("error unmarshal yaml\n%v\n", err)
	}
	return
}
