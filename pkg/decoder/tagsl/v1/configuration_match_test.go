package tagsl

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/truvami/decoder/pkg/common"
)

const (
	// Encoder port-128 fixture: moving 3600, steady 7200, config 86400, gnss 120,
	// accel 300/1500, battery 21600, batch 10, buffer 4096.
	tagSLSentHex = "01010100000e1000001c20000151800078012c05dc00005460000a1000"

	// Matching 28-byte port-4 report with different device state / firmware / hardware.
	tagSLReportCoreHex = "00000e1000001c20000151800078012c05dc01090909000000005460"

	// Matching 32-byte report including batch/buffer.
	tagSLReportFullHex = "00000e1000001c20000151800078012c05dc02020100010200005460000a1000"

	// Decoder port-4 fixture: moving 60, steady 300, config 86400, gnss 120,
	// accel 300/1500, battery 21600.
	tagSLDecoderReportCoreHex = "0000003c0000012c000151800078012c05dc02020100010200005460"
	tagSLDecoderSentHex       = "0000000000003c0000012c000151800078012c05dc00005460000a1000"
	tagSLDecoderReportFullHex = "0000003c0000012c000151800078012c05dc02020100010200005460000a1000"
)

func TestMatchConfiguration(t *testing.T) {
	t.Run("encoder fixture matches 28-byte report", func(t *testing.T) {
		assertConfigurationMatch(t, tagSLSentHex, tagSLReportCoreHex)
	})

	t.Run("encoder fixture matches 32-byte report", func(t *testing.T) {
		assertConfigurationMatch(t, tagSLSentHex, tagSLReportFullHex)
	})

	t.Run("encoder fixture matches 30-byte report", func(t *testing.T) {
		assertConfigurationMatch(t, tagSLSentHex, tagSLReportCoreHex+"000a")
	})

	t.Run("decoder fixture matches 28-byte report", func(t *testing.T) {
		assertConfigurationMatch(t, tagSLDecoderSentHex, tagSLDecoderReportCoreHex)
	})

	t.Run("decoder fixture matches 32-byte report", func(t *testing.T) {
		assertConfigurationMatch(t, tagSLDecoderSentHex, tagSLDecoderReportFullHex)
	})

	t.Run("absent optional tails are skipped", func(t *testing.T) {
		sent := mustMutateHex(t, tagSLSentHex, map[int][]byte{
			25: {0x00, 0x14},
			27: {0x08, 0x00},
		})
		assertConfigurationMatch(t, sent, tagSLReportCoreHex)
	})

	t.Run("batch compared when present", func(t *testing.T) {
		observed := mustMutateHex(t, tagSLReportFullHex, map[int][]byte{28: {0x00, 0x0b}})
		assertConfigurationMismatch(t, tagSLSentHex, observed)
	})

	t.Run("buffer compared when present", func(t *testing.T) {
		observed := mustMutateHex(t, tagSLReportFullHex, map[int][]byte{30: {0x10, 0x01}})
		assertConfigurationMismatch(t, tagSLSentHex, observed)
	})

	t.Run("ignored flags firmware and hardware", func(t *testing.T) {
		sent := mustMutateHex(t, tagSLSentHex, map[int][]byte{
			0: {0x00},
			1: {0x00},
			2: {0x00},
		})
		assertConfigurationMatch(t, sent, tagSLReportCoreHex)
	})

	t.Run("field mismatches", func(t *testing.T) {
		cases := []struct {
			name   string
			offset int
			value  []byte
		}{
			{name: "moving interval", offset: 0, value: []byte{0x00, 0x00, 0x0e, 0x11}},
			{name: "steady interval", offset: 4, value: []byte{0x00, 0x00, 0x1c, 0x21}},
			{name: "config interval", offset: 8, value: []byte{0x00, 0x01, 0x51, 0x81}},
			{name: "gnss timeout", offset: 12, value: []byte{0x00, 0x79}},
			{name: "accelerometer threshold", offset: 14, value: []byte{0x01, 0x2d}},
			{name: "accelerometer delay", offset: 16, value: []byte{0x05, 0xdd}},
			{name: "battery interval", offset: 24, value: []byte{0x00, 0x00, 0x54, 0x61}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				observed := mustMutateHex(t, tagSLReportCoreHex, map[int][]byte{tc.offset: tc.value})
				assertConfigurationMismatch(t, tagSLSentHex, observed)
			})
		}
	})

	t.Run("invalid hex", func(t *testing.T) {
		_, err := MatchConfiguration("zz", tagSLReportCoreHex)
		if err == nil {
			t.Fatal("expected hex error")
		}

		_, err = MatchConfiguration(tagSLSentHex, "abc")
		if err == nil {
			t.Fatal("expected hex error")
		}
	})

	t.Run("wrong sent length", func(t *testing.T) {
		_, err := MatchConfiguration(tagSLSentHex[:len(tagSLSentHex)-2], tagSLReportCoreHex)
		assertConfigurationError(t, err, common.ErrInvalidPayloadLength)
	})

	t.Run("incomplete observed tails", func(t *testing.T) {
		_, err := MatchConfiguration(tagSLSentHex, tagSLReportCoreHex+"00")
		assertConfigurationError(t, err, common.ErrInvalidPayloadLength)

		_, err = MatchConfiguration(tagSLSentHex, tagSLReportFullHex[:len(tagSLReportFullHex)-2])
		assertConfigurationError(t, err, common.ErrInvalidPayloadLength)
	})

	t.Run("decoder-owned field validation", func(t *testing.T) {
		t.Run("sent moving interval 4", func(t *testing.T) {
			sent := mustMutateHex(t, tagSLSentHex, map[int][]byte{3: {0x00, 0x00, 0x00, 0x04}})
			ok, err := MatchConfiguration(sent, tagSLReportCoreHex)
			if ok {
				t.Fatal("expected validation failure, not verification")
			}
			if !errors.Is(err, common.ErrValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})

		t.Run("observed moving interval 4", func(t *testing.T) {
			observed := mustMutateHex(t, tagSLReportCoreHex, map[int][]byte{0: {0x00, 0x00, 0x00, 0x04}})
			ok, err := MatchConfiguration(tagSLSentHex, observed)
			if ok {
				t.Fatal("expected validation failure, not verification")
			}
			if !errors.Is(err, common.ErrValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	})
}

func assertConfigurationMatch(t *testing.T, sent, observed string) {
	t.Helper()
	ok, err := MatchConfiguration(sent, observed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected match")
	}
}

func assertConfigurationMismatch(t *testing.T, sent, observed string) {
	t.Helper()
	ok, err := MatchConfiguration(sent, observed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected non-match")
	}
}

func assertConfigurationError(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}

func mustMutateHex(t *testing.T, payloadHex string, replacements map[int][]byte) string {
	t.Helper()
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	for offset, value := range replacements {
		copy(payload[offset:], value)
	}
	return hex.EncodeToString(payload)
}
