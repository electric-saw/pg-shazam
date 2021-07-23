package state

import (
	"testing"

	"github.com/electric-saw/pg-shazam/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestSessionEncode(t *testing.T) {
	shazam := config.NewShazam()
	ID = shazam.NodeID
	s := Session{
		NodeId: ID,
		PID:    123456789,
	}

	buf := s.Encode()
	assert.Equal(t, 70, len(buf))

	ss := Session{}
	err := ss.Decode(buf)
	assert.Nil(t, err)
	assert.EqualValues(t, s, ss)
}
