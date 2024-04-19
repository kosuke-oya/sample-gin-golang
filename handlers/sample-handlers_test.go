package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSampleHandler(t *testing.T) {
	// Create a new Gin context for testing
	router := gin.Default()
	router.GET("/sample", func(c *gin.Context) {
		SampleHandler(c, "test_value")
	})

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/sample?key=ss", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	recorder := httptest.NewRecorder()

	// Serve the HTTP request to the recorder
	router.ServeHTTP(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", recorder.Code)
	}

	// Parse the response body
	var res Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	// Check the response data
	expectedData := map[string]interface{}{
		"request_key": "ss",
	}

	if res.RequestKey != expectedData["request_key"] {
		t.Errorf("expected %v; got %v", expectedData, res)
	}
}
