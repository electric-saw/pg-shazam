package state

import (
	badger "github.com/dgraph-io/badger/v2"
	"github.com/electric-saw/pg-shazam/internal/pkg/config"

	log "github.com/sirupsen/logrus"
)

var ID string

func NewStateServer(shazam *config.Shazam) (StateServer, error) {
	ID = shazam.NodeID
	log.Info("Initializing kv server...")
	opt := badger.DefaultOptions("").WithInMemory(true).WithLogger(log.New())
	if server, err := badger.Open(opt); err != nil {
		return nil, err
	} else {
		return &stateServer{server: server}, nil
	}
}

type stateServer struct {
	server *badger.DB
}

func (s *stateServer) GetClient() StateStore {
	return &stateStore{
		client: s.server,
	}
}

func (s *stateServer) Close() {
	s.server.Close()
}
