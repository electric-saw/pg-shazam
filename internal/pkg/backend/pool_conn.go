package backend

import (
	"context"
	"encoding/json"
	"net"

	"github.com/electric-saw/pg-shazam/internal/pkg/definitions"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/pkg/util"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Conn struct {
	conn     *pgxpool.Conn
	MsgChan  chan pgproto3.BackendMessage
	NextChan chan struct{}
	front    *pgproto3.Frontend
}

func NewConn(rawConn *pgxpool.Conn) *Conn {
	sock := rawConn.Conn().PgConn().Conn()
	conn := &Conn{
		conn:     rawConn,
		front:    pgproto3.NewFrontend(pgproto3.NewChunkReader(sock), sock),
		MsgChan:  make(chan pgproto3.BackendMessage),
		NextChan: make(chan struct{}),
	}

	return conn
}

func (c *Conn) Release() {
	c.conn.Release()
}

func (c *Conn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return c.conn.Exec(ctx, sql, arguments)
}

func (c *Conn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.conn.Query(ctx, sql, args...)
}

func (c *Conn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.conn.QueryRow(ctx, sql, args...)
}

func (c *Conn) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.conn.SendBatch(ctx, b)
}

func (c *Conn) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return c.conn.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (c *Conn) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.conn.Begin(ctx)
}

func (c *Conn) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.conn.BeginTx(ctx, txOptions)
}

func (c *Conn) Conn() *pgx.Conn {
	return c.conn.Conn()
}

func (c *Conn) Socket() net.Conn {
	return c.Conn().PgConn().Conn()
}

func (c *Conn) ConnStr() string {
	return c.conn.Conn().Config().ConnString()
}
func (c *Conn) ConnStrHidden() string {
	return util.HiddePass(c.conn.Conn().Config().ConnString())
}

func (c *Conn) AssumeClient(ctx context.Context, client *definitions.FrontendClient, lastMsg pgproto3.FrontendMessage) error {
	for _, set := range client.SetCommands {
		_, err := c.conn.Exec(context.Background(), set)
		if err != nil {
			return err
		}
	}

	go c.readServerConn(context.Background())
	_, err := c.Socket().Write(lastMsg.Encode(nil))
	if err != nil {
		return err
	}
	c.ReadNext()
	for {
		select {
		case rawMsg := <-client.MsgChan: // From frontend
			if log.IsLevel(log.TraceLevel) {
				buf, err := json.Marshal(rawMsg)
				if err != nil {
					return err
				}
				log.Tracef("F -> %s", string(buf))
			}

			switch msg := rawMsg.(type) {
			case *definitions.Error:
				return msg.Err

			}
			client.ReadNext()
		case rawMsg := <-c.MsgChan: // From backend
			if log.IsLevel(log.TraceLevel) {
				buf, err := json.Marshal(rawMsg)
				if err != nil {
					return err
				}
				log.Tracef("B-> %s", string(buf))
			}

			switch msg := rawMsg.(type) {
			case *definitions.Error:
				return msg.Err
			case *pgproto3.ReadyForQuery:
				err := client.Backend.Send(msg)
				if err != nil {
					return err
				}
				return nil
			default:
				err := client.Backend.Send(msg)
				if err != nil {
					return err
				}
			}
			c.ReadNext()
		}
	}
}

func (c *Conn) readServerConn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.NextChan:

			msg, err := c.front.Receive()
			if err != nil {
				c.MsgChan <- &definitions.Error{Err: err}
				return
			}
			c.MsgChan <- msg
		}
	}
}

func (c *Conn) Cancel() error {
	return c.conn.Conn().PgConn().CancelRequest(context.TODO())
}

func (c *Conn) ReadNext() {
	c.NextChan <- struct{}{}
}
