package config

// TODO: patroni?
type Replication struct {
	Factor    int `yaml:"factor"`
	MinInsync int `yaml:"min.insync"`
	Strategy  int `yaml:"strategy"`
}

func NewReplication() *Replication {
	return &Replication{
		Factor:    3,
		MinInsync: 2,
	}
}
