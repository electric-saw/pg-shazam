package state

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

type CancelSession interface {
	cancel() error
}

type stateStore struct {
	client        *badger.DB
	activeSession map[uint32]CancelSession
}

func (s *stateStore) GetHashSet(database, table string) (*HashSet, error) {
	key := []byte(fmt.Sprintf("%s:%s", database, table))

	hs := &HashSet{}
	txn := s.client.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	err = item.Value(func(val []byte) error {
		return hs.Decode(val)
	})

	return hs, err
}

func (s *stateStore) SetHashSet(hash *HashSet) error {
	key := []byte(fmt.Sprintf("%s:%s", hash.Database, hash.Table))

	txn := s.client.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set(key, hash.Encode())
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *stateStore) GetSession(pid uint32) (*Session, error) {
	key := []byte(fmt.Sprintf("session:%d", pid))

	session := &Session{}
	txn := s.client.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	err = item.Value(func(val []byte) error {
		return session.Decode(val)
	})

	return session, err
}

func (s *stateStore) SetSession(session *Session) error {
	session.NodeId = ID
	key := []byte(fmt.Sprintf("session:%d", session.PID))

	txn := s.client.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set(key, session.Encode())
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *stateStore) DeleteSession(pid int64) error {
	key := []byte(fmt.Sprintf("session:%d", pid))

	txn := s.client.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Delete(key); err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *stateStore) CancelQuery(session *Session) error {
	if s, ok := s.activeSession[session.PID]; ok {
		return s.cancel()
	}
	return fmt.Errorf("Session not found")
}

func (s *stateStore) GetNodeID() string {
	// TODO:
	return ID
}
