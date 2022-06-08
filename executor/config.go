package executor

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

func NewConfig(configFilePath string) (*Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type Operation string

const (
	AnyOperation    Operation = "ANY"
	CreateOperation Operation = "CREATE"
	WriteOperation  Operation = "WRITE"
	RemoveOperation Operation = "REMOVE"
	RenameOperation Operation = "RENAME"
	ChmodOperation  Operation = "CHMOD"
)

func (o Operation) String() string {
	return string(o)
}

type Rule struct {
	Operation Operation `yaml:"operation"`
	Path      string    `yaml:"path"`
	Cmd       string    `yaml:"cmd"`
}
