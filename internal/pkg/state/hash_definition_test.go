package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	h := HashSet{
		Database: "shazam",
		Table:    "test",
		Fields: []string{
			"id",
			"name",
		},
	}
	buf := h.Encode()
	assert.Equal(t, 57, len(buf))

	hh := HashSet{}
	err := hh.Decode(buf)
	assert.Nil(t, err)
	assert.EqualValues(t, h, hh)
}
