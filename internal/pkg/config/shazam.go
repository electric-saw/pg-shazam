package config

import (
	"io/ioutil"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const defaultAddress = "0.0.0.0:5432"

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
		ListenAddress: defaultAddress,
		Clusters:      []*Cluster{},
		Health:        NewHealth(),
		Replication:   NewReplication(),
		Pool:          NewPoolConfig(),
		Sync:          NewSync(),
		NodeID:        uuid.New().URN(),
	}
}

func (c *Shazam) LoadFromFile(file string) error {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(raw, c)
}
