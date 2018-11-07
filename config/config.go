package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type APIConfig struct {
	Destination   string    `yaml:"destination"`
	Auth          string    `yaml:"auth"`
	ClientTimeout int       `yaml:"client_timeout"`
	Sleep         int32     `yaml:"sleep"`
	Users         int       `yaml:"users"`
	LogLevel      string    `yaml:"log_level"`
	StatsHost     string    `yaml:"stats_host"`
	StatsPort     string    `yaml:"stats_port"`
	Journeys      []Journey `yaml:"journeys"`
}

type Journey struct {
	Name     string    `yaml:"name"`
	Setup    []Request `yaml:"setup"`
	Requests []Request `yaml:"requests"`
	Cleanup  []Request `yaml:"cleanup"`
}

type Request struct {
	Method         string `yaml:"method"`
	Endpoint       string `yaml:"endpoint"`
	Payload        string `yaml:"payload"`
	ExpectedStatus int    `yaml:"expected_status"`
}

func Load(path string) *APIConfig {
	cfg := &APIConfig{}

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
