package backend

import (
	"context"
	"time"

	"github.com/electric-saw/pg-shazam/internal/pkg/config"
	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client struct {
	Id       int64
	Server   *config.Server
	pool     *Pool
	health   bool
	stopChan chan interface{}
}

func NewClient(server *config.Server, poolConf *config.Pool) (*Client, error) {
	client := &Client{
		Server:   server,
		health:   true,
		stopChan: make(chan interface{}),
	}

	connConfig, err := pgxpool.ParseConfig(server.ConnectionString())
	if err != nil {
		return nil, err
	}

	poolConf.EnsureParams(connConfig)

	pool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}

	client.pool = &Pool{
		pool: pool,
	}

	if err = client.connCheck(); err != nil {
		return nil, err
	} else {
		go client.healthCheck()
		return client, nil
	}
}

func (c *Client) connCheck() error {
	conn, err := c.pool.pool.Acquire(context.TODO())
	if err != nil {
		return err
	}

	defer conn.Release()

	ctx, fn := context.WithTimeout(context.Background(), 10*time.Second)
	defer fn()

	res, err := conn.Query(ctx, "select 1;")
	res.Close()
	return err
}

func (c *Client) healthCheck() {
	for {
		select {
		case <-c.stopChan:
			return
		default:
			if err := c.connCheck(); err != nil {
				c.health = false
				log.Warnf("Fail to connect to %s\n%s", c.Server.ConnectionStringHiddenPass(), err)
			}
			// TODO: make config
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) Stat() *pgxpool.Stat {
	return c.pool.Stat()
}

func (c *Client) Acquire(ctx context.Context) (*Conn, error) {
	return c.pool.Acquire(ctx)
}

func (c *Client) Pool() *Pool {
	return c.pool
}
