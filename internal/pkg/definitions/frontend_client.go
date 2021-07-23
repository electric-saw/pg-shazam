package definitions

import (
	"context"
	"encoding/json"
	"net"

	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/jackc/pgproto3/v2"
)

type FrontendClient struct {
	PID          uint32
	SecretKey    uint32
	SetCommands  []string
	CurrDatabase string //TODO: feed this
	Backend      *pgproto3.Backend
	Conn         net.Conn
	MsgChan      chan pgproto3.FrontendMessage
	NextChan     chan struct{}
}

func NewFrontendClient(conn net.Conn) *FrontendClient {
	return &FrontendClient{
		Backend:  pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn),
		Conn:     conn,
		MsgChan:  make(chan pgproto3.FrontendMessage),
		NextChan: make(chan struct{}),
	}
}

func (f *FrontendClient) ReadClient(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-f.NextChan:
			msg, err := f.Backend.Receive()
			if err != nil {
				f.MsgChan <- &Error{err}
				return
			}
			f.logMessage(msg)
			f.MsgChan <- msg
		}
	}
}

func (f *FrontendClient) logMessage(msg pgproto3.Message) {
	if log.IsLevel(log.TraceLevel) {
		buf, err := json.Marshal(msg)
		if err != nil {
			return
		}
		log.Tracef("F-> %s", string(buf))
	}
}

func (f *FrontendClient) ReadNext() {
	f.NextChan <- struct{}{}
}
