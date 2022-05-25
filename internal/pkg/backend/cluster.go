package backend

import (
	"context"

	"github.com/electric-saw/pg-shazam/internal/pkg/log"

	"github.com/electric-saw/pg-shazam/internal/pkg/config"
)

type Cluster struct {
	rw *Client
	ro []*Client
}

func NewCluster(clusterCfg *config.Cluster, pool *config.Pool) (*Cluster, error) {
	rw, err := NewClient(clusterCfg.Rw, pool)
	if err != nil {
		return nil, err
	}

	c := &Cluster{
		rw: rw,
	}

	for _, serverCfg := range clusterCfg.Ro {
		server, err := NewClient(serverCfg, pool)
		if err != nil {
			c.Close()
			return nil, err
		}

		c.ro = append(c.ro, server)
	}

	return c, nil
}

func (c *Cluster) GetROConnection(ctx context.Context) (*Conn, error) {
	s := c.byLowConnections()

	log.Debugf("GetConnection selected %d: %s", s.Id, s.Server.ConnectionStringHiddenPass())
	return s.Acquire(ctx)
}
func (c *Cluster) GetRWConnection(ctx context.Context) (*Conn, error) {
	return c.rw.Acquire(ctx)
}

func (c *Cluster) Close() {
	c.rw.Close()
	for _, ro := range c.ro {
		ro.Close()
	}
}
