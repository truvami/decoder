package loracloud

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func startMockServer() *httptest.Server {
	server := httptest.NewServer(nil)
	return server
}

func TestPost(t *testing.T) {
	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := startMockServer()
	middleware := NewLoracloudMiddleware("access_token")
	middleware.BaseUrl = server.URL
	defer server.Close()

	// Test case 1: Successful request
	url := fmt.Sprintf("%v/success", middleware.BaseUrl)
	body := []byte(`{"key": "value"}`)

	response, err := middleware.post(url, body)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected response: %v, got: %v", http.StatusOK, response)
	}

	// Test case 2: Request with error
	url = fmt.Sprintf("%v/error", middleware.BaseUrl)
	body = []byte(`{"key": "value}`)

	response, err = middleware.post(url, body)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected response: %v, got: %v", http.StatusInternalServerError, response)
	}
}
func TestDeliverUplinkMessage(t *testing.T) {
	http.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		// check if request body contains 0123456789ABCDEC
		bodyString, _ := io.ReadAll(r.Body)
		if strings.Contains(string(bodyString), "01-23-45-67-89-AB-CD-EC") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"status\": failed}"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\": \"success\"}"))
	})

	server := startMockServer()
	middleware := NewLoracloudMiddleware("access_token")
	middleware.BaseUrl = server.URL

	defer server.Close()

	// Test case 1: Successful request
	devEui := "0123456789ABCDEF"
	uplinkMsg := UplinkMsg{
		MsgType:                 "uplink",
		FCount:                  123,
		Port:                    1,
		Payload:                 "0123456789ABCDEF",
		DR:                      nil,
		Frequency:               nil,
		Timestamp:               nil,
		DNMTU:                   nil,
		GNSSCaptureTime:         nil,
		GNSSCaptureTimeAccuracy: nil,
		GNSSAssistPosition:      nil,
		GNSSAssistAltitude:      nil,
		GNSSUse2DSolver:         nil,
	}

	_, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	_, err = middleware.DeliverUplinkMessage("0123456789ABCDEC", uplinkMsg)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
