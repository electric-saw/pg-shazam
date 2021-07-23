package proxy

import (
	"context"
	"net"

	"github.com/electric-saw/pg-shazam/internal/pkg/backend"
	"github.com/electric-saw/pg-shazam/internal/pkg/definitions"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/internal/pkg/state"
)

type Proxy struct {
	ctx        context.Context
	stateStore state.StateStore
	shazam     *backend.Shazam
}

func NewProxy(ctx context.Context, shazam *backend.Shazam) *Proxy {
	return &Proxy{
		ctx:        ctx,
		shazam:     shazam,
		stateStore: shazam.GetStateManager(),
	}
}

func (p *Proxy) Run(conn net.Conn) {
	defer p.closeConn(conn)

	client := definitions.NewFrontendClient(conn)
	if close, err := p.startup(client); err != nil {
		if !close {
			p.handeError(client, err)
		} else {
			return
		}
	}
	defer p.DeleteSession(int64(client.PID))
	err := p.handleMessages(client)
	p.handeError(client, err)
}

func (p *Proxy) closeConn(conn net.Conn) {
	log.Infof("Closing connection %s", conn.RemoteAddr())
	conn.Close()
}

func (p *Proxy) handeError(client *definitions.FrontendClient, err error) {
	if err != nil && err.Error() != "EOF" {
		log.Errorf("[%d]- %s", client.PID, err.Error())
	}

}
