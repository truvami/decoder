package loracloud

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	"go.uber.org/zap"
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
	middleware := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
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
		middleware := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
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
		middleware := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
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
		middleware := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
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
		middleware := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
		middleware.BaseUrl = server.URL
		defer server.Close()

		devEui := "0123456789ABCDEF"
		uplinkMsg := UplinkMsg{
			MsgType: "uplink",
			FCount:  123,
			Port:    1,
			Payload: "0123456789abcdef",
		}

		_, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
		if err == nil || !strings.Contains(err.Error(), "error decoding loracloud response") {
			t.Errorf("expected decoding error, got: %v", err)
		}
	})
}

func TestResponseVariants(t *testing.T) {
	type Expected = struct {
		timestamp *time.Time
		latitude  float64
		longitude float64
		altitude  float64
	}
	var tests = []struct {
		name     string
		result   []byte
		expected Expected
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
			expected: Expected{
				timestamp: common.TimePointer(1722433373.18046),
				latitude:  51.49278,
				longitude: 0.0212,
				altitude:  83.93,
			},
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
			expected: Expected{
				timestamp: common.TimePointer(1722433373.18046),
				latitude:  0,
				longitude: 0,
				altitude:  0,
			},
		},
		{
			name: "captured at null and no algorithm type",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": null
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{
				timestamp: nil,
				latitude:  51.49278,
				longitude: 0.0212,
				altitude:  83.93,
			},
		},
		{
			name: "captured at null and gnss ng algorithm type",
			result: []byte(`{
			"result": {
				"deveui": "927da4b72110927d",
				"position_solution": {
						"llh": [51.49278, 0.0212, 83.93],
						"accuracy": 20.7,
						"gdop": 2.48,
						"capture_time_utc": null,
						"capture_times_utc": [1722433364.06164, 1722433373.18046],
						"algorithm_type": "gnssng"
				},
				"operation": "gnss"
			}
		}`),
			expected: Expected{
				timestamp: common.TimePointer(1722433373.18046),
				latitude:  51.49278,
				longitude: 0.0212,
				altitude:  83.93,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v1/device/send", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(test.result)
			})

			server := startMockServer(mux)
			middleware := NewLoracloudClient(context.TODO(), "access_token", zap.NewExample())
			middleware.BaseUrl = server.URL
			defer server.Close()

			devEui := "b2e6876e64be9692"
			uplinkMsg := UplinkMsg{
				MsgType: "uplink",
				FCount:  42,
				Port:    192,
				Payload: "8c9e50de366a460e8a70fe72e04445db95d1eca8dcdac252",
			}

			response, err := middleware.DeliverUplinkMessage(devEui, uplinkMsg)
			if err != nil {
				t.Fatalf("error %s", err)
			}

			if !common.TimePointerCompare(response.GetTimestamp(), test.expected.timestamp) {
				t.Fatalf("expected timestamp %s got %s", test.expected.timestamp, response.GetTimestamp())
			}
			if response.GetLatitude() != test.expected.latitude {
				t.Fatalf("expected latitude %f got %f", test.expected.latitude, response.GetLatitude())
			}
			if response.GetLongitude() != test.expected.longitude {
				t.Fatalf("expected longitude %f got %f", test.expected.longitude, response.GetLongitude())
			}
			if response.GetAltitude() != test.expected.altitude {
				t.Fatalf("expected altitude %f got %f", test.expected.altitude, response.GetAltitude())
			}
		})
	}
}

func TestValidateContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
	}{
		{
			name:    "missing port",
			ctx:     context.WithValue(context.Background(), decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef"),
			wantErr: ErrContextPortNotFound,
		},
		{
			name: "missing devEui",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				return ctx
			}(),
			wantErr: ErrContextDevEuiNotFound,
		},
		{
			name: "missing fCount",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				return ctx
			}(),
			wantErr: ErrContextFCountNotFound,
		},
		{
			name: "invalid devEui length",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcde") // 15 chars
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 0)
				return ctx
			}(),
			wantErr: ErrContextInvalidDevEui,
		},
		{
			name: "invalid devEui non-hex",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdeg") // 'g' is not hex
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 0)
				return ctx
			}(),
			wantErr: ErrContextInvalidDevEui,
		},
		{
			name: "invalid fCount negative",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(1))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, -1)
				return ctx
			}(),
			wantErr: ErrContextInvalidFCount,
		},
		{
			name: "valid context",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), decoder.PORT_CONTEXT_KEY, uint8(10))
				ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, "0123456789abcdef")
				ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 42)
				return ctx
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContext(tt.ctx)
			if tt.wantErr == nil {
				assert.NoError(t, err, "expected no error but got: %v", err)
				return
			}

			assert.Error(t, err, "expected error but got none")
			assert.ErrorIs(t, err, tt.wantErr, "expected error to match")

		})
	}
}

func TestWithBaseUrl(t *testing.T) {
	// Create a LoracloudClient with a default BaseUrl
	client := LoracloudClient{
		BaseUrl: "https://default.url",
	}

	// Use WithBaseUrl to set a new BaseUrl
	newUrl := "https://custom.url"
	option := WithBaseUrl(newUrl)
	option(&client)

	if client.BaseUrl != newUrl {
		t.Errorf("expected BaseUrl to be %q, got %q", newUrl, client.BaseUrl)
	}
}
