package deploy

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type deployConfig struct {
	Secrets  map[string]*secretConfig `yaml:"secrets"`
	Clusters []*clusterConfig         `yaml:"clusters"`
}

type clusterConfig struct {
	Name    string   `yaml:"name"`
	Context string   `yaml:"context"`
	Secrets []string `yaml:"secrets"`
}

type secretConfig struct {
	Name       string            `yaml:"name"`
	Files      map[string]string `yaml:"files"`
	Namespaces []string          `yaml:"namespaces"`
}

func loadConfig(filename string) (*deployConfig, error) {
	var cfg deployConfig
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
