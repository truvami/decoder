package smartlabel

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/truvami/decoder/pkg/loracloud"
)


func TestDecode(t *testing.T) {

	middleware := loracloud.NewLoracloudMiddleware("access_token")

	tests := []struct {
		payload     string
		port        int16
		devEui      string
		expected    interface{}
		expectedErr string
	}{
		// TODO: test 192, 197
		{
			payload:     "00",
			port:        0,
			devEui:      "",
			expected:    nil,
			expectedErr: "port 0 not supported",
		},
		{
			payload: "0e021e0384012c01003c012c03e8",
			port:    11,
			expected: Port11ConfigurationPayload{
				Flags:               30,
				GNSSEnabled:         1,
				WiFiEnabled:         1,
				AccEnabled:          1,
				StaticSF:            "SF9",
				SteadyIntervalS:     900,
				MovingIntervalS:     300,
				HeartbeatIntervalH:  1,
				LEDBlinkIntervalS:   60,
				AccThresholdMS:      300,
				AccDelayMS:          1000,
			},
		},
		{
			payload: "11021e0384012c01003c012c03e8e43420ea",
			port:    11,
			expected: Port11ConfigurationPayload{
				Flags:               30,
				GNSSEnabled:         1,
				WiFiEnabled:         1,
				AccEnabled:          1,
				StaticSF:            "SF9",
				SteadyIntervalS:     900,
				MovingIntervalS:     300,
				HeartbeatIntervalH:  1,
				LEDBlinkIntervalS:   60,
				AccThresholdMS:      300,
				AccDelayMS:          1000,
				GitHash:             "e43420ea",
			},
		},
		{
			payload: "0a010f05095f4100000000",
			port:    11,
			expected: Port11HeartbeatPayload{
				Battery:           3.845,
				Temperature:       23.99,
				RH:                32.5,
				GNSSScansCount:    0,
				WiFiScansCount:    0,
			},
		},
		{
			payload:     "",
			port:        11,
			devEui:      "",
			expected:    nil,
			expectedErr: "data length is less than 2",
		},
		{
			payload:     "0eXX",
			port:        11,
			devEui:      "",
			expected:    nil,
			expectedErr: "encoding/hex",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPort%vWith%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewSmartLabelv1Decoder(middleware)
			got, _, err := decoder.Decode(test.payload, test.port, test.devEui)
			if err != nil && len(test.expectedErr) == 0 {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if !reflect.DeepEqual(got, test.expected) && len(test.expectedErr) == 0 {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}

			if len(test.expectedErr) > 0 && err != nil && !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("expected error: %v, got: %v", test.expectedErr, err)
			}
		})
	}
}

func TestInvalidPort(t *testing.T) {
	decoder := NewSmartLabelv1Decoder(loracloud.NewLoracloudMiddleware("appEui"))
	_, _, err := decoder.Decode("00", 0, "")
	if err == nil {
		t.Fatal("expected port not supported")
	}
}

func TestWithAutoPadding(t *testing.T) {
	middleware := loracloud.NewLoracloudMiddleware("access_token")
	
	decoder := NewSmartLabelv1Decoder(
		middleware,
		WithAutoPadding(true),
	)
	
	// Type assert to access the internal field
	if d, ok := decoder.(*SmartLabelv1Decoder); ok {
		if !d.autoPadding {
			t.Error("expected autoPadding to be true")
		}
	} else {
		t.Error("failed to type assert decoder")
	}
}

func TestGetPort11PayloadType(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		want        string
		expectedErr string
	}{
		{
			name:        "Empty data",
			data:        "",
			want:        "",
			expectedErr: "data length is less than 2",
		},
		{
			name:        "Single byte",
			data:        "0",
			want:        "",
			expectedErr: "data length is less than 2",
		},
		{
			name:        "Configuration payload 0E",
			data:        "0E",
			want:        "configuration",
			expectedErr: "",
		},
		{
			name:        "Configuration payload 11",
			data:        "11",
			want:        "configuration",
			expectedErr: "",
		},
		{
			name:        "Heartbeat payload",
			data:        "0A",
			want:        "heartbeat",
			expectedErr: "",
		},
		{
			name:        "Invalid payload",
			data:        "FF",
			want:        "",
			expectedErr: "invalid payload for port 11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPort11PayloadType(tt.data)
			if err != nil && tt.expectedErr == "" {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tt.expectedErr != "" {
				t.Errorf("expected error containing %q, got nil", tt.expectedErr)
			}
			if err != nil && tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectedErr, err.Error())
			}
			if got != tt.want {
				t.Errorf("getPort11PayloadType() = %v, want %v", got, tt.want)
			}
		})
	}
}
