package tagxl

import (
	"fmt"
	"testing"

	"github.com/truvami/decoder/pkg/middleware"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		payload  string
		port     int16
		expected interface{}
	}{}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v-%v", test.port, test.payload), func(t *testing.T) {
			decoder := NewTagXLv1Decoder(middleware.NewLoraCloudClient("token"))
			got, err := decoder.Decode(test.payload, test.port, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("got %v", got)

			if got != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, got)
			}
		})
	}
}
