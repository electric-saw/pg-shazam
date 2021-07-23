package config

import (
	"time"

	"github.com/electric-saw/pg-shazam/internal/pkg/log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Pool struct {
	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	MaxConnLifetime string

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime string

	// MaxConns is the maximum size of the pool.
	MaxConns int32

	// HealthCheckPeriod is the duration between checks of the health of idle connections.
	HealthCheckPeriod string

	// If set to true, pool doesn't do any I/O operation on initialization.
	// And connects to the server only when the pool starts to be used.
	// The default is false.
	LazyConnect bool
}

func NewPoolConfig() *Pool {
	return &Pool{
		MaxConnLifetime:   "1h",
		MaxConnIdleTime:   "20m",
		MaxConns:          500,
		HealthCheckPeriod: "20s",
		LazyConnect:       true,
	}
}

func parseDurationDef(val string, def time.Duration) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Warnf("Value %s is not duration supported format! check your config file \n %s", val, err.Error())
		return def
	} else {
		return d
	}
}

func (p *Pool) EnsureParams(conf *pgxpool.Config) {
	conf.MaxConnLifetime = parseDurationDef(p.MaxConnLifetime, 1*time.Hour)
	conf.MaxConnIdleTime = parseDurationDef(p.MaxConnIdleTime, 20*time.Minute)
	conf.MaxConns = p.MaxConns
	conf.HealthCheckPeriod = parseDurationDef(p.HealthCheckPeriod, 20*time.Second)
	conf.LazyConnect = p.LazyConnect
}
