package config

type Sync struct {
	ListenAddress string `yaml:"address"`
	DataDir       string `yaml:"data_path"`
}

func NewSync() *Sync {
	return &Sync{
		ListenAddress: "127.0.0.1:5333",
		DataDir:       "/tmp/data",
	}
}
