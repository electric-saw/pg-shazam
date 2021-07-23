package backend

import (
	"fmt"
	"sort"

	"github.com/electric-saw/pg-shazam/internal/pkg/log"
	"github.com/electric-saw/pg-shazam/internal/pkg/parser"

	"github.com/cespare/xxhash/v2"
)

const MaxPartitions = 64

type Partition struct {
	Id    uint64
	Value *Cluster
}

type Ring struct {
	size       uint64
	partitions []*Partition
}

func NewRing(nodes ...*Cluster) (*Ring, error) {
	s := len(nodes)
	if s <= 0 || s > MaxPartitions {
		return nil, fmt.Errorf("invalid number of paritions: %d, keep on 1 and %d", s, MaxPartitions)
	}

	ring := &Ring{
		partitions: make([]*Partition, s),
		size:       uint64(s),
	}

	for id, node := range nodes {
		ring.partitions[id] = &Partition{
			Id:    uint64(id),
			Value: node,
		}
	}

	return ring, nil
}

func (r *Ring) GetPartition(qry *parser.Query, shardColumns []string) *Partition {
	key := getShardingKeyValue(qry, shardColumns)
	idx := xxhash.Sum64(key) % r.size
	log.Tracef("GetPartition of %s is %d", key, idx)
	return r.partitions[idx]
}

func getShardingKeyValue(qry *parser.Query, shardColumns []string) []byte {
	sort.Strings(shardColumns)

	keyShard := ""

	for _, v := range qry.Conditions {
		for _, col := range shardColumns {
			if v.Field == col {
				keyShard = keyShard + v.Value
			}
		}
	}

	return []byte(keyShard)
}
