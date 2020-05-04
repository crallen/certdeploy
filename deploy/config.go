package deploy

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type deployConfig struct {
	Clusters []*clusterConfig `yaml:"clusters"`
}

type clusterConfig struct {
	Name    string          `yaml:"name"`
	Context string          `yaml:"context"`
	Secrets []*secretConfig `yaml:"secrets"`
}

type secretConfig struct {
	Name       string   `yaml:"name"`
	Files      []string `yaml:"files"`
	Namespaces []string `yaml:"namespaces"`
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
