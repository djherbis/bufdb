// bufdb is used to persist buffer.Buffer's to a bolt.DB database.
package bufdb

import (
	"github.com/boltdb/bolt"
	"github.com/djherbis/buffer"
	"github.com/djherbis/stow"
)

// BufferStore manages buffer.Buffer persistance.
type BufferStore struct {
	buffers *stow.Store
}

// NewBufferStore creates a new BufferStore, using the underlying
// bolt.DB "bucket" to persist buffers.
func NewBufferStore(db *bolt.DB, bucket []byte) *BufferStore {
	return &BufferStore{buffers: stow.NewStore(db, bucket)}
}

// Put persists buffer.Buffer b, modifying b after Put will not affect the stored buffer.
// This will fail if the buffer is not gob-encodable.
func (s *BufferStore) Put(key []byte, b buffer.Buffer) error {
	return s.buffers.Put(key, b)
}

// Pull will retrive (and removes) the buffer stored under "key".
// ErrNotFound indicates that the buffer was not found in the database,
// other errors indicate gob-decoding failures.
func (s *BufferStore) Pull(key []byte) (b buffer.Buffer, err error) {
	err = s.buffers.Pull(key, &b)
	return b, err
}

// DeleteAll will call Reset() on every buffer managed by this store,
// and then delete the buffer from the store.
func (s *BufferStore) DeleteAll() error {
	s.buffers.ForEach(func(b buffer.Buffer) {
		b.Reset()
	})
	return s.buffers.DeleteAll()
}

type BufferPoolStore struct {
	pools *stow.Store
}

func NewBufferPoolStore(db *bolt.DB, bucket []byte) *BufferPoolStore {
	return &BufferPoolStore{pools: stow.NewStore(db, bucket)}
}

func (s *BufferPoolStore) Put(key []byte, p buffer.Pool) error {
	return s.pools.Put(key, p)
}

func (s *BufferPoolStore) Pull(key []byte) (p buffer.Pool, err error) {
	err = s.pools.Pull(key, &p)
	return p, err
}

func (s *BufferPoolStore) DeleteAll() error {
	return s.pools.DeleteAll()
}
