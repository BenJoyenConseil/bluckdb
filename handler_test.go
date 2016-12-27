package main

import (
	"testing"
	"strconv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"gopkg.in/gavv/httpexpect.v1"
	"os"
	"net/http"
	"github.com/BenJoyenConseil/bluckdb/bluckstore"
)

const db_test_path = "/tmp/bluckdb/"

func irisTester(t *testing.T) *httpexpect.Expect {
	server := &server{
		stores: make(map[string]bluckstore.ThreadSafeStore),
		lock: &sync.RWMutex{},
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
	response := tester.GET("/v1/meta/tmp/bluckdb/meta/").Expect()

	// Then
	response.Status(http.StatusOK).JSON().Schema(schema)
	rmDBFiles()
}

func TestIrisHandler_GET(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)

	// When
	response := tester.GET("/v1/data/tmp/bluckdb/get/").WithQuery("id", "123").Expect()

	// Then
	response.Status(http.StatusOK).JSON().Object().ContainsKey("key").ContainsKey("val")
}

func TestIrisHandler_PUT(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)

	// When
	response := tester.PUT("/v1/data/tmp/bluckdb/put/").WithQuery("id", "123").WithText("yop%20yop%20yop").Expect()

	// Then
	fmt.Println(response.Text())
	response.Status(http.StatusOK)
	rmDBFiles()
}

func TestIrisHandler_PUT_RaceCondition(t *testing.T) {
	rmDBFiles()
	// Given
	tester := irisTester(t)

	// When
	for i := 0; i < 100; i++ {
		go tester.PUT("/v1/data/tmp/bluckdb/race/").WithQuery("id", "mykey").WithText("|myvalue->thread n°" + strconv.Itoa(i) + ";").Expect()
		go tester.PUT("/v1/data/tmp/bluckdb/race/").WithQuery("id", "mykey").WithText("|myvalue->thread n°" + strconv.Itoa(i) + ";").Expect()
		go tester.PUT("/v1/data/tmp/bluckdb/race/").WithQuery("id", "mykey").WithText("|myvalue->thread n°" + strconv.Itoa(i) + ";").Expect()
		go tester.PUT("/v1/data/tmp/bluckdb/race/").WithQuery("id", "mykey").WithText("|myvalue->thread n°" + strconv.Itoa(i) + ";").Expect()
	}

	// Then
	rmDBFiles()
}

func TestIrisHandler_GET_DEBUG(t *testing.T) {
	rmDBFiles()

	// Given
	tester := irisTester(t)

	// When
	response := tester.GET("/v1/debug/tmp/bluckdb/debug/").WithQuery("page_id", "0").Expect()

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
