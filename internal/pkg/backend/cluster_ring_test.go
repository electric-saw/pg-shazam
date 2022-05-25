package backend

import (
	"testing"

	"github.com/electric-saw/pg-shazam/internal/pkg/parser"
	"github.com/stretchr/testify/assert"
)

//TODO: Create new tests

func TestRing(t *testing.T) {
	r, _ := NewRing(&Cluster{}, &Cluster{}, &Cluster{})

	values := map[string]uint64{
		"rodrigo":  0,
		"lucas":    1,
		"anderson": 2,
		"joanio":   1,
		"glauco":   2,
		"paulo":    2,
		"matheus":  0,
	}

	for value, nodeId := range values {
		assert.Equalf(t, r.GetPartition(&parser.Query{
			Conditions: []parser.Condition{{
				Field: "name",
				Value: value}},
		}, []string{"name", "id"}).Id, nodeId, value)
	}
}
