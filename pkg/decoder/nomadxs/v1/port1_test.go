package nomadxs

import (
	"testing"
	"time"
)

func TestGetCapturedAt(t *testing.T) {
	tests := []struct {
		name     string
		payload  Port1Payload
		expected *time.Time
	}{
		{
			name: "valid date and time",
			payload: Port1Payload{
				Year:   21,
				Month:  10,
				Day:    5,
				Hour:   14,
				Minute: 30,
				Second: 45,
			},
			expected: func() *time.Time {
				t, _ := time.Parse(time.RFC3339, "2021-10-05T14:30:45Z")
				return &t
			}(),
		},
		{
			name: "invalid date",
			payload: Port1Payload{
				Year:   21,
				Month:  13,
				Day:    32,
				Hour:   25,
				Minute: 61,
				Second: 61,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payload.GetCapturedAt()
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %v, got nil", tt.expected)
				} else if !result.Equal(*tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
