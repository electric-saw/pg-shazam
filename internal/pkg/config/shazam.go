package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/google/uuid"
)

type Shazam struct {
	ListenAddress string       `yaml:"address"`
	Clusters      []*Cluster   `yaml:"clusters"`
	Health        *Health      `yaml:"health"`
	Replication   *Replication `yaml:"replication"`
	Pool          *Pool        `yaml:"pool"`
	Sync          *Sync        `yaml:"sync"`
	NodeID        string       `yaml:"nodeId"`
}

func NewShazam() *Shazam {
	return &Shazam{
		ListenAddress: "0.0.0.0:5432",
		Health:        NewHealth(),
		Replication:   NewReplication(),
		Pool:          NewPoolConfig(),
		NodeID:        uuid.New().URN(),
		Sync:          NewSync(),
	}
}

func (c *Shazam) LoadFromFile(file string) error {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(raw, c)
}
