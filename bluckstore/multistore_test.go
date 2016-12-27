package bluckstore

import (
	"testing"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
)

const db_test_path = "/tmp/bluckdb/"

func rmDBFiles() {
	os.RemoveAll(db_test_path)
}

func TestServerGetStore_StoreExists(t *testing.T) {
	// Given
	store1 := &mockStore{}
	s := &MultiStore{
		stores: map[string]ThreadSafeStore{
			"/bla/bla/bla/": store1,
		},
		lock: &sync.RWMutex{},
	}

	// When
	result := s.GetStore("/bla/bla/bla/")

	// Then
	assert.Equal(t, store1, result)
	rmDBFiles()
}

func TestServerClose(t *testing.T) {
	// Given
	store1 := &mockStore{hasCalledClose: false}
	store2 := &mockStore{hasCalledClose: false}
	s := &MultiStore{
		stores: map[string]ThreadSafeStore{
			"/bla/bla/bla/": store1,
			"/blu/blu/blu/": store2,
		},
		lock: &sync.RWMutex{},
	}

	// When
	s.Close()

	// Then
	assert.True(t, store1.hasCalledClose)
	assert.True(t, store2.hasCalledClose)
	rmDBFiles()
}

type mockStore struct {
	hasCalledClose bool
}

func (s *mockStore) Get(k string) string   { panic("not implemented") }
func (s *mockStore) Put(k, v string) error { panic("not implemented") }
func (s *mockStore) DumpPage(i int) string { panic("not implemented") }
func (s *mockStore) Open(p string)         { panic("not implemented") }
func (s *mockStore) Close()                { s.hasCalledClose = true }
func (s *mockStore) Meta() *mmap.Directory { panic("not implemented") }
func (s *mockStore) Lock()                 { panic("not implemented") }
func (s *mockStore) Unlock()               { panic("not implemented") }

func TestServerGetStore_StoreDoesNotExist_ShouldCreateAndOpen_WithPath_AndReturnIt(t *testing.T) {
	// Given
	rmDBFiles()
	s := &MultiStore{
		stores: make(map[string]ThreadSafeStore),
		lock:   &sync.RWMutex{},
	}

	// When
	result := s.GetStore(db_test_path)

	// Then
	f, err := os.Open(db_test_path + "bluck.data")
	defer f.Close()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	rmDBFiles()
}
