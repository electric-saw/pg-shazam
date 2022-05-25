package auth

import (
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	User = "test"
	Pass = "123"
)

func TestValidateUser(t *testing.T) {
	h := md5.New()
	_, _ = h.Write([]byte(Pass + User))

	pass := fmt.Sprintf("md5%x", string(h.Sum(nil)))

	ok, err := passwordCheck(User, Pass, pass)

	assert.True(t, ok)
	assert.NoError(t, err)

}
