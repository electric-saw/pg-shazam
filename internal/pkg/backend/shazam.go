package backend

import (
	"context"

	"github.com/electric-saw/pg-shazam/internal/pkg/config"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/internal/pkg/parser"
	"github.com/electric-saw/pg-shazam/internal/pkg/state"
	"github.com/jackc/pgx/v4"
)

type Shazam struct {
	Cfg         *config.Shazam
	StateServer state.StateServer
	Clusters    []*Cluster
	ring        *Ring
}

func NewShazam(shazamCfg *config.Shazam) (*Shazam, error) {
	stateServer, err := state.NewStateServer(shazamCfg)
	if err != nil {
		return nil, err
	}
	log.Infof("Initializing shazam backend...")

	c := &Shazam{
		Cfg:         shazamCfg,
		StateServer: stateServer,
	}

	for _, clusterCfg := range shazamCfg.Clusters {
		cluster, err := NewCluster(clusterCfg, shazamCfg.Pool)
		if err != nil {
			return nil, err
		}

		c.Clusters = append(c.Clusters, cluster)
	}
	r, err := NewRing(c.Clusters...)
	if err != nil {
		return nil, err
	}

	c.ring = r
	InitShazamCatalog(c)

	return c, nil
}

func (s *Shazam) GetStateManager() state.StateStore {
	return s.StateServer.GetClient()
}

func (s *Shazam) ClusterByHash(qry *parser.Query, shardColumns []string) *Cluster {
	return s.ring.GetPartition(qry, shardColumns).Value
}

func (s *Shazam) RunAllPrimaryHosts(qry string) []error {
	var txs []pgx.Tx
	var errs []error

	for _, cluster := range s.Clusters {
		conn, err := cluster.rw.Acquire(context.Background())
		if err != nil {
			return []error{err}
		}
		defer conn.Release()

		tx, err := conn.Begin(context.Background())
		if err != nil {
			return []error{err}
		}
		_, err = tx.Exec(context.Background(), qry)
		if err != nil {
			errs = append(errs, err)
		}

		txs = append(txs, tx)
	}

	for _, tx := range txs {
		if len(errs) > 0 {
			err := tx.Rollback(context.Background())
			if err != nil {
				errs = append(errs, err)
			}
		}
		err := tx.Commit(context.Background())
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (s *Shazam) GetROConnection(ctx context.Context) (*Conn, error) {
	var client *Client
	var count float32 = float32(9999999)

	for _, cluster := range s.Clusters {
		c, clientCount := cluster.byLowConnectionsWithMetrics()
		if clientCount < count {
			client = c
			count = clientCount
		}
	}

	return client.Acquire(ctx)
}

func (s *Shazam) GetRandomCluster() *Cluster {
	// TODO: adjust with metrics
	return s.Clusters[0]
}
