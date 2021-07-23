package backend

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (p *Pool) Acquire(ctx context.Context) (*Conn, error) {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	} else {
		return NewConn(conn), nil
	}
}

func (p *Pool) AcquireAllIdle(ctx context.Context) []*Conn {
	rawConns := p.pool.AcquireAllIdle(ctx)
	conns := make([]*Conn, len(rawConns))
	for idx, raw := range rawConns {
		conns[idx] = NewConn(raw)
	}
	return conns
}

func (p *Pool) Config() *pgxpool.Config {
	return p.pool.Config()
}

func (p *Pool) Stat() *pgxpool.Stat {
	return p.pool.Stat()
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

func (p *Pool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *Pool) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return p.pool.SendBatch(ctx, b)
}

func (p *Pool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.BeginTx(ctx, pgx.TxOptions{})
}
func (p *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}

func (p *Pool) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return p.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
