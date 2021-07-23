package auth

import (
	"context"
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/electric-saw/pg-shazam/internal/pkg/backend"

	"github.com/electric-saw/pg-shazam/internal/pkg/config"

	"github.com/stretchr/testify/assert"
)

const (
	User = "test"
	Pass = "123"
)

func TestValidateUser(t *testing.T) {
	shazamCfg := config.NewShazam()
	err := shazamCfg.LoadFromFile("../../../conf/sample.yaml")
	assert.Nil(t, err)

	cluster, err := backend.NewCluster(shazamCfg.Clusters[0], shazamCfg.Pool)
	assert.Nil(t, err, "Shazan compose is up?")

	h := md5.New()
	_, _ = h.Write([]byte(Pass + User))

	conn, err := cluster.GetROConnection(context.Background())
	assert.Nil(t, err)
	defer conn.Release()

	ok, msg := ValidateUser(conn, User, fmt.Sprintf("md5%x", string(h.Sum(nil))))

	assert.True(t, ok, msg)

}
