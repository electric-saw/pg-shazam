package config

type Cluster struct {
	Rw *Server   `yaml:"rw"`
	Ro []*Server `yaml:"ro"`
}
