package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func TestAddDecoder(t *testing.T) {
	router := http.NewServeMux()
	path := "test/path"
	decoder := tagsl.NewTagSLv1Decoder()

	addDecoder(router, path, decoder)

	handler, pattern := router.Handler(&http.Request{Method: "POST", URL: &url.URL{Path: "/test/path"}})
	if handler == nil {
		t.Errorf("expected handler to be set")
	}
	if pattern != "POST /test/path" {
		t.Errorf("expected pattern to be 'POST /test/path', got '%s'", pattern)
	}
}
func TestGetHandler(t *testing.T) {
	decoder := tagsl.NewTagSLv1Decoder()
	handler := getHandler(decoder)

	reqBody := `{"port": 1, "payload": "8002cdcd1300744f5e166018040b14341a", "devEui": "example devEUI"}`
	req, err := http.NewRequest("POST", "/test/path", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	handler(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("expected Content-Type header to be %q, got %q", expectedContentType, actualContentType)
	}

	// test with invalid JSON
	reqBody = `{"port": 1, "payload": "8002cdcd1300744f5e166018040b14341a", "devEui": "example devEUI"`
	req, err = http.NewRequest("POST", "/test/path", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	recorder = httptest.NewRecorder()
	handler(recorder, req)

	resp = recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// test with invalid payload
	reqBody = `{"port": 1, "payload": "invalid", "devEui": "example devEUI"}`
	req, err = http.NewRequest("POST", "/test/path", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	recorder = httptest.NewRecorder()
	handler(recorder, req)

	resp = recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestSetHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	status := http.StatusOK

	setHeaders(recorder, status)

	expectedContentType := "application/json"
	actualContentType := recorder.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("expected Content-Type header to be %q, got %q", expectedContentType, actualContentType)
	}

	expectedAllowOrigin := "*"
	actualAllowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	if actualAllowOrigin != expectedAllowOrigin {
		t.Errorf("expected Access-Control-Allow-Origin header to be %q, got %q", expectedAllowOrigin, actualAllowOrigin)
	}

	expectedAllowMethods := "POST"
	actualAllowMethods := recorder.Header().Get("Access-Control-Allow-Methods")
	if actualAllowMethods != expectedAllowMethods {
		t.Errorf("expected Access-Control-Allow-Methods header to be %q, got %q", expectedAllowMethods, actualAllowMethods)
	}

	expectedStatus := http.StatusOK
	actualStatus := recorder.Result().StatusCode
	if actualStatus != expectedStatus {
		t.Errorf("expected status code %d, got %d", expectedStatus, actualStatus)
	}
}

func TestHTTPCmd(t *testing.T) {
	go func() {
		// call the command handler function
		httpCmd.Run(nil, []string{"--host", "::1"})
	}()

	// create a new HTTP request to simulate the command execution
	reqBody := `{"port": 1, "payload": "8002cdcd1300744f5e166018040b14341a", "devEui": "example devEUI"}`
	req, err := http.NewRequest("POST", "http://localhost:8080/tagsl/v1", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// check the Content-Type header
	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("expected Content-Type header to be %q, got %q", expectedContentType, actualContentType)
	}
}
func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	healthHandler(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("expected Content-Type header to be %q, got %q", expectedContentType, actualContentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expectedBody := "OK"
	actualBody := string(body)
	if actualBody != expectedBody {
		t.Errorf("expected response body to be %q, got %q", expectedBody, actualBody)
	}
}
