package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/truvami/decoder/internal/logger"
	tagsl "github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func TestAddDecoder(t *testing.T) {
	logger.NewLogger()
	defer logger.Sync()

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
	logger.NewLogger()
	defer logger.Sync()

	decoder := tagsl.NewTagSLv1Decoder()
	handler := getHandler(decoder)

	reqBody := `{"port": 1, "payload": "8002cdcd1300744f5e166018040b14341a", "devEui": ""}`
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expectedBody := `47.041811`
	actualBody := string(responseBody)

	if strings.Contains(expectedBody, actualBody) {
		t.Errorf("expected response body to contain %q, got %q", expectedBody, actualBody)
	}

	// test with invalid JSON
	reqBody = `{"port": 1, "payload": "8002cdcd1300744f5e166018040b14341a", "devEui": ""`
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
	reqBody = `{"port": 1, "payload": "invalid", "devEui": ""}`
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
	logger.NewLogger()
	defer logger.Sync()

	if httpCmd.Flags().Set("port", "38888") != nil {
		t.Fatalf("failed to set port flag")
	}
	if httpCmd.Flags().Set("host", "127.0.0.1") != nil {
		t.Fatalf("failed to set host flag")
	}

	go func() {
		// call the command handler function
		httpCmd.Run(nil, []string{})
	}()

	// create a new HTTP request to simulate the command execution
	reqBody := `{"port": 105, "payload": "0028672658500172a741b1e238b572a741b1e08bb03498b5c583e2b172a741b1e0cda772a741beed4cc472a741beef53b7"}`
	req, err := http.NewRequest("POST", "http://localhost:38888/tagsl/v1", strings.NewReader(reqBody))
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

	// parse the response body
	type Response struct {
		Data     tagsl.Port105Payload `json:"data"`
		Metadata tagsl.Status         `json:"metadata"`
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// check the response body
	expectedData := tagsl.Port105Payload{
		Moving:      true,
		DutyCycle:   false,
		BufferLevel: 40,
		Timestamp:   time.Date(2024, 11, 2, 16, 50, 24, 0, time.UTC),
		Mac1:        "72a741b1e238",
		Rssi1:       -75,
		Mac2:        "72a741b1e08b",
		Rssi2:       -80,
		Mac3:        "3498b5c583e2",
		Rssi3:       -79,
		Mac4:        "72a741b1e0cd",
		Rssi4:       -89,
		Mac5:        "72a741beed4c",
		Rssi5:       -60,
		Mac6:        "72a741beef53",
		Rssi6:       -73,
	}

	if response.Data != expectedData {
		t.Errorf("expected response data to be %v, got %v", expectedData, response.Data)
	}

	expectedMetadata := tagsl.Status{
		DutyCycle:           false,
		ConfigChangeId:      0,
		ConfigChangeSuccess: false,
		Moving:              true,
	}
	if response.Metadata != expectedMetadata {
		t.Errorf("expected response metadata to be %v, got %v", expectedMetadata, response.Metadata)
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
