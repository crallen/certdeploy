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
	Name       string        `yaml:"name"`
	Files      []*secretFile `yaml:"files"`
	Namespaces []string      `yaml:"namespaces"`
}

type secretFile struct {
	Key      string `yaml:"key"`
	Filename string `yaml:"filename"`
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
