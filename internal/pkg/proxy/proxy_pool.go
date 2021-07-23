package proxy

import (
	"context"
	"log"
	"net"

	"github.com/electric-saw/pg-shazam/internal/pkg/backend"

	"github.com/jackc/puddle"
)

type ProxyPool struct {
	shazam *backend.Shazam
	pool   *puddle.Pool
}

func NewProxyPool(shazam *backend.Shazam) *ProxyPool {
	p := &ProxyPool{
		shazam: shazam,
	}

	p.pool = puddle.NewPool(p.createProxy, p.destroyProxy, shazam.Cfg.Pool.MaxConns)
	err := p.pool.CreateResource(context.Background())
	if err != nil {
		log.Fatalf("Failed to create pool resource %s", err.Error())
	}
	return p
}

func (p *ProxyPool) AddJob(conn net.Conn) error {
	r, err := p.pool.Acquire(context.Background())
	if err != nil {
		return err
	} else {
		go r.Value().(*Proxy).Run(conn)
	}
	return nil
}

func (p *ProxyPool) createProxy(ctx context.Context) (interface{}, error) {
	proxy := NewProxy(ctx, p.shazam)
	return proxy, nil
}
func (p *ProxyPool) destroyProxy(value interface{}) {
}
