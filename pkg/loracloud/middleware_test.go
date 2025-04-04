package loracloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func startMockServer(handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

func TestPost(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := startMockServer(mux)
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
	t.Run("Successful request", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"result": {
					"deveui": "01-23-45-67-89-AB-CD-EF",
					"pending_requests": {
						"requests": [],
						"id": 1,
						"updelay": 0,
						"upcount": 0
					},
					"info_fields": {},
					"log_messages": [],
					"fports": {
						"dmport": 1,
						"gnssport": 2,
						"wifiport": 3,
						"fragport": 4,
						"streamport": 5,
						"gnssngport": 6
					},
					"operation": "other"
				}
			}`))
		})

		server := startMockServer(mux)
		middleware := NewLoracloudMiddleware("access_token")
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789ABCDEF",
		}

		response, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if response.Result.Deveui != "0123456789ABCDEF" {
			t.Errorf("expected deveui: %v, got: %v", "0123456789ABCDEF", response.Result.Deveui)
		}
	})

	t.Run("Validation error", func(t *testing.T) {
		server := startMockServer(nil)
		middleware := NewLoracloudMiddleware("access_token")
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "",
			FCount:  123,
			Port:    1,
			Payload: "0123456789ABCDEF",
		}

		_, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !strings.Contains(err.Error(), "error validating uplink message") {
			t.Errorf("expected validation error, got: %v", err)
		}
	})

	t.Run("Unexpected status code", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors": ["Invalid request"]}`))
		})

		server := startMockServer(mux)
		middleware := NewLoracloudMiddleware("access_token")
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789ABCDEF",
		}

		_, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !strings.Contains(err.Error(), "unexpected status code returned by loracloud") {
			t.Errorf("expected status code error, got: %v", err)
		}
	})

	t.Run("Error decoding response", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid-json`))
		})

		server := startMockServer(mux)
		middleware := NewLoracloudMiddleware("access_token")
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789ABCDEF",
		}

		_, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !strings.Contains(err.Error(), "error decoding loracloud response") {
			t.Errorf("expected decoding error, got: %v", err)
		}
	})
}

func TestResponseVariants(t *testing.T) {
	var tests = []struct {
		name   string
		result []byte
	}{
		{
			name: "normal response",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": 1722433373.18046
				},
				"operation": "gnss"
			}
		}`),
		},
		{
			name: "llh empty array",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"llh": [],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": 1722433373.18046
				},
				"operation": "gnss"
			}
		}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(test.result)

				server := startMockServer(mux)
				middleware := NewLoracloudMiddleware("token")
				middleware.BaseUrl = server.URL
				defer server.Close()

				devEui := "b2e6876e64be9692"
				uplinkMsg := UplinkMsg{
					MsgType: "uplink",
					FCount:  42,
					Port:    192,
					Payload: "8c9e50de366a460e8a70fe72e04445db95d1eca8dcdac252",
				}

				_, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
				if err != nil {
					t.Fatalf("error %s", err)
				}
			})
		})
	}
}
