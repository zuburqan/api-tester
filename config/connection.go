package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type ConnectionConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func LoadConnection(path string) *ConnectionConfig {
	cfg := &ConnectionConfig{}

	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
