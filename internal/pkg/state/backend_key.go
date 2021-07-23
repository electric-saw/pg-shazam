package state

import (
	"net"

	"github.com/cespare/xxhash/v2"
)

func NewBackendKey(client net.Conn) (pid uint32, secretKey uint32) {
	digest := xxhash.New()
	_, _ = digest.WriteString(client.RemoteAddr().String())

	pid = uint32(digest.Sum64())

	_, _ = digest.WriteString(ID)

	secretKey = uint32(digest.Sum64())

	return
}
