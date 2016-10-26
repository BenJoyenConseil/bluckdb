package main

import (
	"net/http"
	"testing"

	"fmt"
	"github.com/BenJoyenConseil/bluckdb/bluckstore"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"
	"os"
)

const db_test_path = "/tmp/bluckdb/"

func irisTester(t *testing.T) *httpexpect.Expect {
	server := &server{
		stores: make(map[string]bluckstore.KVStore),
	}
	handler := IrisHandler(server)
	handler.Build()

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://localhost:2233",
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(handler.Router),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
		},
	})
}

func rmDBFiles() {
	os.RemoveAll(db_test_path)
}

func TestIrisHandler_GET_META(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)
	schema := `{
			"type": "object",
			"properties": {
				"LastPageId": {
					"type": "integer"
				}
			},
			"properties": {
				"globalDepth": {
					"type": "integer"
				}
			},
			"properties": {
				"table": {
					"type": "array",
					"items": {
						"type": "integer"
					}
				}
			}

		    }`

	// When
	response := tester.GET("/v1/meta/tmp/bluckdb/").Expect()

	// Then
	response.Status(http.StatusOK).JSON().Schema(schema)
	rmDBFiles()
}

func TestIrisHandler_GET(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)

	// When
	response := tester.GET("/v1/data/tmp/bluckdb/").WithQuery("id", "123").Expect()

	// Then
	response.Status(http.StatusOK).JSON().Object().ContainsKey("key").ContainsKey("val")
}

func TestIrisHandler_PUT(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)

	// When
	response := tester.PUT("/v1/data/tmp/bluckdb/").WithQuery("id", "123").WithText("yop%20yop%20yop").Expect()

	// Then
	fmt.Println(response.Text())
	response.Status(http.StatusOK)
	rmDBFiles()
}

func TestIrisHandler_GET_DEBUG(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)

	// When
	response := tester.GET("/v1/debug/tmp/bluckdb/").WithQuery("page_id", "0").Expect()

	// Then
	fmt.Println(response.Text())
	response.Status(http.StatusOK)
	rmDBFiles()
}

func TestExtractDynamicPath(t *testing.T) {
	// Given
	fixedPath := "/v1/meta"
	fullPath := "/v1/meta/path/to/the/table/"

	// When
	result := extractDynamicPath(fixedPath, fullPath)

	// Then
	assert.Equal(t, "/path/to/the/table/", result)
}

func TestServerGetStore_StoreExists(t *testing.T) {
	// Given
	store1 := &mockStore{}
	s := &server{
		stores: map[string]bluckstore.KVStore{
			"/bla/bla/bla/": store1,
		},
	}

	// When
	result := s.getStore("/bla/bla/bla/")

	// Then
	assert.Equal(t, store1, result)
	rmDBFiles()
}

func TestServerClose(t *testing.T) {
	// Given
	store1 := &mockStore{hasCalledClose: false}
	store2 := &mockStore{hasCalledClose: false}
	s := &server{
		stores: map[string]bluckstore.KVStore{
			"/bla/bla/bla/": store1,
			"/blu/blu/blu/": store2,
		},
	}

	// When
	s.close()

	// Then
	assert.True(t, store1.hasCalledClose)
	assert.True(t, store2.hasCalledClose)
	rmDBFiles()
}

type mockStore struct {
	hasCalledClose bool
}

func (s *mockStore) Get(k string) string   { return "" }
func (s *mockStore) Put(k, v string) error { return nil }
func (s *mockStore) DumpPage(i int) string { return "" }
func (s *mockStore) Open(p string)         {}
func (s *mockStore) Close()                { s.hasCalledClose = true }
func (s *mockStore) Meta() *mmap.Directory { return nil }

func TestServerGetStore_StoreDoesNotExist_ShouldCreateAndOpen_WithPath_AndReturnIt(t *testing.T) {
	// Given
	s := &server{
		stores: make(map[string]bluckstore.KVStore),
	}

	// When
	result := s.getStore(db_test_path)

	// Then
	assert.NotNil(t, result)
	assert.Equal(t, db_test_path, result.(*mmap.MmapKVStore).Path)
	rmDBFiles()
}
