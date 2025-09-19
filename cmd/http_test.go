package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/common"
	tagslDecoder "github.com/truvami/decoder/pkg/decoder/tagsl/v1"
	tagslEncoder "github.com/truvami/decoder/pkg/encoder/tagsl/v1"
)

func TestAddDecoder(t *testing.T) {
	logger.NewLogger()
	defer logger.Sync()

	router := http.NewServeMux()
	path := "test/path"
	decoder := tagslDecoder.NewTagSLv1Decoder()

	addDecoder(context.TODO(), router, path, decoder)

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

	decoder := tagslDecoder.NewTagSLv1Decoder()
	handler := getHandler(context.TODO(), decoder)

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
	// Enable health endpoint so we can wait for readiness
	if httpCmd.Flags().Set("health", "true") != nil {
		t.Fatalf("failed to set health flag")
	}

	go func() {
		// call the command handler function
		httpCmd.Run(nil, []string{})
	}()

	// Wait for server readiness using health endpoint
	ready := false
	for i := 0; i < 100; i++ { // up to ~2s with 20ms sleep
		resp, err := http.Get("http://127.0.0.1:38888/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				ready = true
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !ready {
		t.Fatalf("server not ready on /health")
	}

	// create a new HTTP request to simulate the command execution
	reqBody := `{"port": 105, "payload": "0028672658500172a741b1e238b572a741b1e08bb03498b5c583e2b172a741b1e0cda772a741beed4cc472a741beef53b7"}`
	req, err := http.NewRequest("POST", "http://127.0.0.1:38888/tagsl/v1", strings.NewReader(reqBody))
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
		Data tagslDecoder.Port105Payload `json:"data"`
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// check the response body
	expectedData := tagslDecoder.Port105Payload{
		BufferLevel:  40,
		Timestamp:    time.Date(2024, 11, 2, 16, 50, 24, 0, time.UTC),
		DutyCycle:    false,
		ConfigId:     0,
		ConfigChange: false,
		Moving:       true,
		Mac1:         "72a741b1e238",
		Rssi1:        -75,
		Mac2:         common.StringPtr("72a741b1e08b"),
		Rssi2:        common.Int8Ptr(-80),
		Mac3:         common.StringPtr("3498b5c583e2"),
		Rssi3:        common.Int8Ptr(-79),
		Mac4:         common.StringPtr("72a741b1e0cd"),
		Rssi4:        common.Int8Ptr(-89),
		Mac5:         common.StringPtr("72a741beed4c"),
		Rssi5:        common.Int8Ptr(-60),
		Mac6:         common.StringPtr("72a741beef53"),
		Rssi6:        common.Int8Ptr(-73),
	}

	if !reflect.DeepEqual(response.Data, expectedData) {
		t.Errorf("expected response data to be %v, got %v", expectedData, response.Data)
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

func TestAddEncoder(t *testing.T) {
	logger.NewLogger()
	defer logger.Sync()

	router := http.NewServeMux()
	path := "encode/test/path"
	encoder := tagslEncoder.NewTagSLv1Encoder()

	addEncoder(router, path, encoder)

	handler, pattern := router.Handler(&http.Request{Method: "POST", URL: &url.URL{Path: "/encode/test/path"}})
	if handler == nil {
		t.Errorf("expected handler to be set")
	}
	if pattern != "POST /encode/test/path" {
		t.Errorf("expected pattern to be 'POST /encode/test/path', got '%s'", pattern)
	}
}

func TestGetEncoderHandler(t *testing.T) {
	encoder := tagslEncoder.NewTagSLv1Encoder()
	handler := getEncoderHandler(encoder)

	// Test with Port128Payload
	payload := tagslEncoder.Port128Payload{
		Ble:                    false,
		Gnss:                   true,
		Wifi:                   true,
		MovingInterval:         3600,
		SteadyInterval:         7200,
		ConfigInterval:         86400,
		GnssTimeout:            120,
		AccelerometerThreshold: 300,
		AccelerometerDelay:     1500,
		BatteryInterval:        21600,
		BatchSize:              10,
		BufferSize:             4096,
	}

	reqBody, err := json.Marshal(map[string]any{
		"port":    128,
		"payload": payload,
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/encode/tagsl/v1", bytes.NewReader(reqBody))
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

	// Test with invalid JSON
	reqBody = []byte(`{"port": 128, "payload": {`)
	req, err = http.NewRequest("POST", "/encode/tagsl/v1", bytes.NewReader(reqBody))
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

	// Test with invalid port
	reqBody, err = json.Marshal(map[string]any{
		"port":    0, // Invalid port
		"payload": payload,
		"devEui":  "",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err = http.NewRequest("POST", "/encode/tagsl/v1", bytes.NewReader(reqBody))
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

func TestMetricsEndpoint(t *testing.T) {
	logger.NewLogger()
	defer logger.Sync()

	// Enable metrics endpoint
	metrics = true
	defer func() { metrics = false }()

	router := http.NewServeMux()
	if metrics {
		router.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("metrics ok"))
			if err != nil {
				t.Fatalf("failed to write response: %v", err)
				return
			}
		}))
	}

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("failed to GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if !strings.Contains(string(body), "metrics ok") {
		t.Errorf("expected response body to contain 'metrics ok', got %q", string(body))
	}
}
