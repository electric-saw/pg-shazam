package frontend

import (
	"fmt"
	"net"

	"github.com/electric-saw/pg-shazam/internal/pkg/backend"
	"github.com/electric-saw/pg-shazam/internal/pkg/config"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/internal/pkg/proxy"
)

type Frontend struct {
	listenAddress string
	shazam        *backend.Shazam
	conf          *config.Shazam
	proxyPool     *proxy.ProxyPool
}

func NewFrontend(configFile string) (*Frontend, error) {
	conf := config.NewShazam()
	err := conf.LoadFromFile(configFile)
	if err != nil {
		return nil, err
	}

	shazam, err := backend.NewShazam(conf)
	if err != nil {
		return nil, fmt.Errorf("Failed on create cluster %w", err)
	}

	return &Frontend{
		conf:      conf,
		shazam:    shazam,
		proxyPool: proxy.NewProxyPool(shazam),
	}, nil
}

func (s *Frontend) Close() {
	s.shazam.StateServer.Close()
}

func (s *Frontend) Run() error {
	PrintHead()

	ln, err := net.Listen("tcp", s.conf.ListenAddress)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Infof("Listening on %s", ln.Addr())
	if err != nil {
		log.Fatalf("%v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Accepted connection from %s", conn.RemoteAddr())

		err = s.proxyPool.AddJob(conn)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
}
