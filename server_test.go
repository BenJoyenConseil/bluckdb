package main


import (
	"net/http"
	"testing"

	"gopkg.in/gavv/httpexpect.v1"
	"os"
	"fmt"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
)



func irisTester(t *testing.T) *httpexpect.Expect {
	store := &mmap.MmapKVStore{}
	store.Open()
	handler := IrisHandler(store)

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://localhost:2233",
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
		},
	})
}

func TestIrisHandler_GET_META(t *testing.T) {
	// Given
	os.Remove("/tmp/bluck.data")
	os.Remove("/tmp/bluck.meta")
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
	response := tester.GET("/meta").Expect()

	// Then
	response.Status(http.StatusOK).JSON().Schema(schema)
}

func TestIrisHandler_GET(t *testing.T) {
	// Given
	os.Remove("/tmp/bluck.data")
	os.Remove("/tmp/bluck.meta")
	tester := irisTester(t)

	// When
	response := tester.GET("/").WithQuery("id", "123").Expect()

	// Then
	response.Status(http.StatusOK).JSON().Object().ContainsKey("key").ContainsKey("val")
}

func TestIrisHandler_PUT(t *testing.T) {
	// Given
	os.Remove("/tmp/bluck.data")
	os.Remove("/tmp/bluck.meta")
	tester := irisTester(t)

	// When
	response := tester.PUT("/", ).WithQuery("id", "123").WithText("yop%20yop%20yop").Expect()

	// Then
	fmt.Println(response.Text())
	response.Status(http.StatusOK)
}