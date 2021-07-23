package state

type StateServer interface {
	GetClient() StateStore
	Close()
}

type StateStore interface {
	GetHashSet(database, table string) (*HashSet, error)
	SetHashSet(hash *HashSet) error
	GetSession(pid uint32) (*Session, error)
	SetSession(session *Session) error
	DeleteSession(pid int64) error
	CancelQuery(session *Session) error
	GetNodeID() string
}
