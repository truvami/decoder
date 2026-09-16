package tagxl

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/truvami/decoder/pkg/common"
)

func TestMatchConfiguration(t *testing.T) {
	t.Run("captured heartbeat fixture", func(t *testing.T) {
		ok, err := MatchConfiguration("4c0401230118", "4c0301430118")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected match")
		}
	})

	t.Run("outer length and count quirks", func(t *testing.T) {
		ok, err := MatchConfiguration("4cff99230118", "4c0001430118")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected match despite envelope metadata quirks")
		}
	})

	t.Run("all setter mappings", func(t *testing.T) {
		cases := []struct {
			name     string
			sent     string
			observed string
		}{
			{
				name:     "device flags 20 to 40",
				sent:     "4c010120010f",
				observed: "4c010140010f",
			},
			{
				name:     "moving intervals 21 to 41",
				sent:     "4c01012104012c1c20",
				observed: "4c01014104012c1c20",
			},
			{
				name:     "acceleration 22 to 42",
				sent:     "4c01012204006403e8",
				observed: "4c01014204006403e8",
			},
			{
				name:     "heartbeat 23 to 43",
				sent:     "4c0101230118",
				observed: "4c0101430118",
			},
			{
				name:     "advertisement interval 24 to 44",
				sent:     "4c010124013c",
				observed: "4c010144013c",
			},
			{
				name:     "rotation flags 25 to 47",
				sent:     "4c0101250103",
				observed: "4c0101470103",
			},
			{
				name:     "data rate 28 to 4e",
				sent:     "4c0101280102",
				observed: "4c01014e0102",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				ok, err := MatchConfiguration(tc.sent, tc.observed)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !ok {
					t.Fatal("expected match")
				}
			})
		}
	})

	t.Run("mixed and reordered settings", func(t *testing.T) {
		ok, err := MatchConfiguration("4c0a022301182104012c1c20", "4c0a024104012c1c20430118")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected match")
		}
	})

	t.Run("extra read-only report tags", func(t *testing.T) {
		ok, err := MatchConfiguration("4c0401230118", "4c090143011845020fa0460411223344")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected match")
		}
	})

	t.Run("reserved bit masks", func(t *testing.T) {
		t.Run("device flags low four bits", func(t *testing.T) {
			ok, err := MatchConfiguration("4c0101200107", "4c01014001f7")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected match on masked device flags")
			}
		})

		t.Run("rotation flags low two bits", func(t *testing.T) {
			ok, err := MatchConfiguration("4c0101250102", "4c01014701fe")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected match on masked rotation flags")
			}
		})

		t.Run("device flags masked mismatch", func(t *testing.T) {
			ok, err := MatchConfiguration("4c0101200107", "4c0101400108")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok {
				t.Fatal("expected non-match")
			}
		})
	})

	t.Run("differing value", func(t *testing.T) {
		ok, err := MatchConfiguration("4c0401230118", "4c0301430119")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected non-match")
		}
	})

	t.Run("missing reflected field", func(t *testing.T) {
		ok, err := MatchConfiguration("4c0401230118", "4c030145020fa0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("expected non-match")
		}
	})

	t.Run("wrong marker", func(t *testing.T) {
		_, err := MatchConfiguration("4d0401230118", "4c0301430118")
		assertConfigurationError(t, err, errConfigurationWrongMarker)

		_, err = MatchConfiguration("4c0401230118", "4d0301430118")
		assertConfigurationError(t, err, errConfigurationWrongMarker)
	})

	t.Run("truncation", func(t *testing.T) {
		_, err := MatchConfiguration("4c04012301", "4c0301430118")
		assertConfigurationError(t, err, errConfigurationMalformedTLV)

		_, err = MatchConfiguration("4c0401230118", "4c03014301")
		assertConfigurationError(t, err, errConfigurationMalformedTLV)
	})

	t.Run("duplicate comparable tags", func(t *testing.T) {
		_, err := MatchConfiguration("4c0602230118230118", "4c0301430118")
		assertConfigurationError(t, err, errConfigurationDuplicateTag)

		ok, err := MatchConfiguration("4c0401230118", "4c0602430118430118")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected matching duplicate observed setter responses")
		}
	})

	t.Run("unsupported command", func(t *testing.T) {
		_, err := MatchConfiguration("4c01018500", "4c0101430118")
		assertConfigurationError(t, err, errConfigurationUnsupportedCommand)
	})

	t.Run("action-only has no comparable requirement", func(t *testing.T) {
		_, err := MatchConfiguration("4c010180020005", "4c0101430118")
		assertConfigurationError(t, err, errConfigurationNoSetter)
	})

	t.Run("pure getter presence", func(t *testing.T) {
		ok, err := MatchConfiguration("4c01014000", "4c010140010f")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatal("expected getter presence to match")
		}
	})

	t.Run("invalid data rate", func(t *testing.T) {
		_, err := MatchConfiguration("4c0101280108", "4c01014e0108")
		assertConfigurationError(t, err, errConfigurationInvalidDataRate)

		_, err = MatchConfiguration("4c0101280102", "4c01014e0108")
		assertConfigurationError(t, err, errConfigurationInvalidDataRate)

		_, err = CompareConfigurationFor(ConfigurationDialectTagXL, "4c01014e00", "4c01014e0108")
		assertConfigurationError(t, err, errConfigurationInvalidDataRate)

		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c01014e00", "4c01014e0102")
		if err != nil {
			t.Fatalf("valid getter-only data rate: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("wrong observed comparable length", func(t *testing.T) {
		_, err := MatchConfiguration("4c0401230118", "4c050143021819")
		assertConfigurationError(t, err, errConfigurationMalformedTLV)
	})

	t.Run("no setter command", func(t *testing.T) {
		_, err := MatchConfiguration("4c0000", "4c0301430118")
		assertConfigurationError(t, err, errConfigurationNoSetter)
	})

	t.Run("ignored nonsensical read-only tags", func(t *testing.T) {
		t.Run("short battery tag", func(t *testing.T) {
			ok, err := MatchConfiguration("4c0401230118", "4c0601430118450101")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected match despite short read-only battery tag")
			}
		})

		t.Run("unknown extra tag", func(t *testing.T) {
			ok, err := MatchConfiguration("4c0401230118", "4c0601430118ff03010203")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected match despite unknown extra tag")
			}
		})

		t.Run("truncated read-only framing", func(t *testing.T) {
			_, err := MatchConfiguration("4c0401230118", "4c040143011845")
			assertConfigurationError(t, err, errConfigurationMalformedTLV)

			_, err = MatchConfiguration("4c0401230118", "4c0401430118ff02")
			assertConfigurationError(t, err, errConfigurationMalformedTLV)
		})
	})

	t.Run("decoder-owned field validation", func(t *testing.T) {
		t.Run("advertisement interval 0", func(t *testing.T) {
			ok, err := MatchConfiguration("4c0101240100", "4c0101440100")
			if ok {
				t.Fatal("expected validation failure, not verification")
			}
			if !errors.Is(err, common.ErrValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})

		t.Run("heartbeat 169", func(t *testing.T) {
			ok, err := MatchConfiguration("4c01012301a9", "4c01014301a9")
			if ok {
				t.Fatal("expected validation failure, not verification")
			}
			if !errors.Is(err, common.ErrValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})

		t.Run("moving interval 59", func(t *testing.T) {
			ok, err := MatchConfiguration("4c01012104003b1c20", "4c01014104003b1c20")
			if ok {
				t.Fatal("expected validation failure, not verification")
			}
			if !errors.Is(err, common.ErrValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})

		t.Run("acceleration delay 999", func(t *testing.T) {
			ok, err := MatchConfiguration("4c01012204006403e7", "4c01014204006403e7")
			if ok {
				t.Fatal("expected validation failure, not verification")
			}
			if !errors.Is(err, common.ErrValidationFailed) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	})
}

func TestCompareConfiguration(t *testing.T) {
	const sent = "4c0d0320010a210400780708280100"

	t.Run("captured campaign reports are incomplete", func(t *testing.T) {
		reports := []struct {
			name     string
			observed string
		}{
			{name: "data rate only", observed: "4c03014e0100"},
			{name: "flags 1b with telemetry", observed: "4c0d0345020d774b040e45000040011b"},
			{name: "matching flags with telemetry", observed: "4c0d0345020c0b4b048609000040018a"},
		}
		for _, tc := range reports {
			t.Run(tc.name, func(t *testing.T) {
				result, err := CompareConfiguration(sent, tc.observed)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != ConfigurationIncomplete {
					t.Fatalf("got %s, want incomplete", result)
				}
				ok, err := MatchConfiguration(sent, tc.observed)
				if err != nil || ok {
					t.Fatalf("MatchConfiguration: ok=%v err=%v", ok, err)
				}
			})
		}
	})

	t.Run("complete reflection matches", func(t *testing.T) {
		result, err := CompareConfiguration(sent, "4c0d0340010a4104007807084e0100")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("complete contradiction mismatches", func(t *testing.T) {
		result, err := CompareConfiguration(sent, "4c0d0340010a4104007807084e0101")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMismatch {
			t.Fatalf("got %s, want mismatch", result)
		}
	})

	t.Run("present differing field still incomplete when another setter is missing", func(t *testing.T) {
		result, err := CompareConfiguration(sent, "4c0d0340011b")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationIncomplete {
			t.Fatalf("got %s, want incomplete", result)
		}
	})

	t.Run("single setter complete mismatch", func(t *testing.T) {
		result, err := CompareConfiguration("4c0401230118", "4c0301430119")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMismatch {
			t.Fatalf("got %s, want mismatch", result)
		}
	})

	t.Run("missing reflected field is incomplete", func(t *testing.T) {
		result, err := CompareConfiguration("4c0401230118", "4c030145020fa0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationIncomplete {
			t.Fatalf("got %s, want incomplete", result)
		}
	})
}

func TestConfigurationComparisonString(t *testing.T) {
	if ConfigurationMatch.String() != "match" {
		t.Fatalf("match string = %q", ConfigurationMatch.String())
	}
	if ConfigurationMismatch.String() != "mismatch" {
		t.Fatalf("mismatch string = %q", ConfigurationMismatch.String())
	}
	if ConfigurationIncomplete.String() != "incomplete" {
		t.Fatalf("incomplete string = %q", ConfigurationIncomplete.String())
	}
	if ConfigurationComparison(99).String() != "unknown" {
		t.Fatalf("unknown string = %q", ConfigurationComparison(99).String())
	}
}

func TestPort151PayloadConfigCoversComparableTLVTags(t *testing.T) {
	tags := make(map[uint8]struct{})
	for _, tag := range port151PayloadConfig().Tags {
		tags[tag.Tag] = struct{}{}
	}
	for setter, spec := range setterSpecs {
		if _, ok := tags[spec.tlvTag]; !ok {
			t.Errorf("TLV tag 0x%02x (setter 0x%02x) missing from port151PayloadConfig", spec.tlvTag, setter)
		}
	}
}

func TestTagXLDialectOverlapsPort151SetterSpecs(t *testing.T) {
	spec, err := dialectSpecFor(ConfigurationDialectTagXL)
	if err != nil {
		t.Fatalf("dialect: %v", err)
	}
	for setter, port151 := range setterSpecs {
		command, ok := spec.setters[setter]
		if !ok {
			t.Errorf("tagXLDialect missing setter 0x%02x", setter)
			continue
		}
		if command.getterTag != port151.tlvTag {
			t.Errorf("setter 0x%02x getterTag=0x%02x, want 0x%02x", setter, command.getterTag, port151.tlvTag)
		}
		if command.setterLen != port151.valueLen {
			t.Errorf("setter 0x%02x setterLen=%d, want %d", setter, command.setterLen, port151.valueLen)
		}
		if command.responseLen != port151.valueLen {
			t.Errorf("setter 0x%02x responseLen=%d, want %d", setter, command.responseLen, port151.valueLen)
		}
	}
}

func assertConfigurationError(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}

func TestCompareConfigurationForDialects(t *testing.T) {
	t.Run("tag xl firmware getter ignores value", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c01014600", "4c0101460411223344")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("tag xl mixed setter and getter", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c050220010a4600", "4c0a0240010a4604aabbccdd")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("tag xl setter plus matching getter prefers setter equality", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c050220010a4000", "4c010140010b")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMismatch {
			t.Fatalf("got %s, want mismatch", result)
		}
	})

	t.Run("tag xl missing getter is incomplete", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c030246004000", "4c010140010a")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationIncomplete {
			t.Fatalf("got %s, want incomplete", result)
		}
	})

	t.Run("tag xl hold interval setter", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c01012702012c", "4c01014d02012c")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("tag xl rotation enable masks extra bits", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c01012b0181", "4c0101510101")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("smart label flags include ble bit", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c0101200110", "4c0101400110")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}

		result, err = CompareConfigurationFor(ConfigurationDialectTagXL, "4c0101200110", "4c0101400100")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("tag xl should ignore ble bit: got %s", result)
		}
	})

	t.Run("smart label motion compares threshold only", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c01012204006403e8", "4c01014204006407d0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}

		result, err = CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c01012204006403e8", "4c01014204006507d0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMismatch {
			t.Fatalf("got %s, want mismatch", result)
		}
	})

	t.Run("smart label firmware getter", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c01014600", "4c0101460401020304")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("smart label heartbeat 169 is valid", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c01012301a9", "4c01014301a9")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("smart label rejects colliding tag xl rotation setter", func(t *testing.T) {
		_, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c0101250103", "4c0101470103")
		assertConfigurationError(t, err, errConfigurationUnsupportedCommand)
	})

	t.Run("smart label rejects nonresponsive scan timer", func(t *testing.T) {
		_, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c01013203010078", "4c01016203010078")
		assertConfigurationError(t, err, errConfigurationUnsupportedCommand)
	})

	t.Run("smart label rejects secret lte api key getter", func(t *testing.T) {
		_, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, "4c0101b200", "4c0101b20401020304")
		assertConfigurationError(t, err, errConfigurationUnsupportedCommand)
	})

	t.Run("wrong getter length is malformed", func(t *testing.T) {
		_, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c01014600", "4c010146020011")
		assertConfigurationError(t, err, errConfigurationMalformedTLV)
	})

	t.Run("duplicate observed setter values that differ mismatch", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c0101230118", "4c0602430118430119")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMismatch {
			t.Fatalf("got %s, want mismatch", result)
		}
	})

	t.Run("unknown dialect", func(t *testing.T) {
		_, err := CompareConfigurationFor(ConfigurationDialectUnspecified, "4c0101230118", "4c0101430118")
		assertConfigurationError(t, err, errConfigurationUnknownDialect)
	})
}

func TestBuildCurrentConfigurationRequest(t *testing.T) {
	t.Run("converts setters and keeps getters", func(t *testing.T) {
		got, err := BuildCurrentConfigurationRequest(ConfigurationDialectTagXL, "4c0d0320010a2104007807084600")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want, err := hex.DecodeString("4c0703400041004600")
		if err != nil {
			t.Fatalf("decode want: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("got %x, want %x", got, want)
		}
	})

	t.Run("deduplicates setter and getter of the same tag", func(t *testing.T) {
		got, err := BuildCurrentConfigurationRequest(ConfigurationDialectTagXL, "4c050220010a4000")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want, err := hex.DecodeString("4c03014000")
		if err != nil {
			t.Fatalf("decode want: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("got %x, want %x", got, want)
		}
	})

	t.Run("smart label battery thresholds", func(t *testing.T) {
		got, err := BuildCurrentConfigurationRequest(ConfigurationDialectSmartLabelV2, "4c01012902050a")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want, err := hex.DecodeString("4c03014c00")
		if err != nil {
			t.Fatalf("decode want: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("got %x, want %x", got, want)
		}
	})

	t.Run("rejects oversized smart label response", func(t *testing.T) {
		sent := hex.EncodeToString([]byte{
			0x4c, 0x00, 0x00,
			0x2c, 0x10, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
			0x30, 0x07, 0x00, 0x64, 0x00, 0x01, 0x02, 0x03, 0x04,
			0x21, 0x04, 0x01, 0x2c, 0x1c, 0x20,
			0x20, 0x01, 0x0f,
			0x23, 0x01, 0x06,
			0x46, 0x00,
			0x45, 0x00,
			0x49, 0x00,
			0x4a, 0x00,
			0x4b, 0x00,
		})
		_, err := BuildCurrentConfigurationRequest(ConfigurationDialectSmartLabelV2, sent)
		assertConfigurationError(t, err, errConfigurationTooLarge)
	})

	t.Run("round trip generated request", func(t *testing.T) {
		request, err := BuildCurrentConfigurationRequest(ConfigurationDialectTagXL, "4c0101280102")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, hex.EncodeToString(request), "4c01014e0102")
		if err != nil {
			t.Fatalf("compare: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})
}

func TestValidateConfiguration(t *testing.T) {
	if err := ValidateConfiguration(ConfigurationDialectTagXL, "4c01014600"); err != nil {
		t.Fatalf("expected getter-only profile to validate: %v", err)
	}
	if err := ValidateConfiguration(ConfigurationDialectTagXL, "4c01018100"); err != nil {
		t.Fatalf("expected action-only profile to validate: %v", err)
	}
	if err := ValidateConfiguration(ConfigurationDialectSmartLabelV2, "4c0101250101"); err == nil {
		t.Fatal("expected colliding Tag XL command to fail Smart Label validation")
	}

	payload := []byte{0x4c, 0x00, 0x00}
	for range 33 {
		payload = append(payload, tlvTagDeviceFlags, 0x00)
	}
	assertConfigurationError(t, ValidateConfiguration(ConfigurationDialectTagXL, hex.EncodeToString(payload)), errConfigurationTooManyCommands)
}

func TestConfigurationDialectString(t *testing.T) {
	if ConfigurationDialectTagXL.String() != "tagxl" {
		t.Fatalf("tag xl string = %q", ConfigurationDialectTagXL.String())
	}
	if ConfigurationDialectSmartLabelV2.String() != "smartlabel-v2" {
		t.Fatalf("smart label string = %q", ConfigurationDialectSmartLabelV2.String())
	}
	if ConfigurationDialectUnspecified.String() != "unspecified" {
		t.Fatalf("unspecified string = %q", ConfigurationDialectUnspecified.String())
	}
}

func TestTagXLCommandMatrix(t *testing.T) {
	cases := []struct {
		name     string
		sent     string
		observed string
	}{
		{name: "hold 27", sent: "4c01012702012c", observed: "4c01014d02012c"},
		{name: "buffer ack 29", sent: "4c0101290105", observed: "4c01014f0105"},
		{name: "timestamped 2a", sent: "4c01012a0103", observed: "4c0101500103"},
		{name: "rotation enable 2b", sent: "4c01012b0101", observed: "4c0101510101"},
		{name: "time 4c", sent: "4c01014c00", observed: "4c01014c0401020304"},
		{name: "records 52", sent: "4c01015200", observed: "4c01015206000101020304"},
		{name: "battery 45", sent: "4c01014500", observed: "4c010145020fa0"},
		{name: "reset count 49", sent: "4c01014900", observed: "4c01014902000c"},
		{name: "reset cause 4a", sent: "4c01014a00", observed: "4c01014a0400000001"},
		{name: "scan counts 4b", sent: "4c01014b00", observed: "4c01014b0400010002"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CompareConfigurationFor(ConfigurationDialectTagXL, tc.sent, tc.observed)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != ConfigurationMatch {
				t.Fatalf("got %s, want match", result)
			}
		})
	}
}

func TestAnalyzeConfiguration(t *testing.T) {
	cases := []struct {
		name       string
		dialect    ConfigurationDialect
		sent       string
		observable bool
		hasActions bool
		wantErr    error
	}{
		{name: "tag xl setter", dialect: ConfigurationDialectTagXL, sent: "4c0101230118", observable: true},
		{name: "tag xl getter", dialect: ConfigurationDialectTagXL, sent: "4c01014600", observable: true},
		{name: "tag xl alarm", dialect: ConfigurationDialectTagXL, sent: "4c010180020005", hasActions: true},
		{name: "tag xl reset", dialect: ConfigurationDialectTagXL, sent: "4c01018100", hasActions: true},
		{name: "tag xl scan", dialect: ConfigurationDialectTagXL, sent: "4c01018200", hasActions: true},
		{name: "tag xl clear storage", dialect: ConfigurationDialectTagXL, sent: "4c01018300", hasActions: true},
		{name: "tag xl wipe all", dialect: ConfigurationDialectTagXL, sent: "4c01018400", hasActions: true},
		{name: "tag xl mixed setter and reset", dialect: ConfigurationDialectTagXL, sent: "4c0a02210400780e108100", observable: true, hasActions: true},
		{name: "tag xl mixed getter and reset", dialect: ConfigurationDialectTagXL, sent: "4c030140008100", observable: true, hasActions: true},
		{name: "tag xl multiple actions", dialect: ConfigurationDialectTagXL, sent: "4c030281008200", hasActions: true},
		{name: "tag xl empty", dialect: ConfigurationDialectTagXL, sent: "4c0000", wantErr: errConfigurationNoSetter},
		{name: "tag xl unknown", dialect: ConfigurationDialectTagXL, sent: "4c01018500", wantErr: errConfigurationUnsupportedCommand},
		{name: "smart label reset", dialect: ConfigurationDialectSmartLabelV2, sent: "4c01018100", hasActions: true},
		{name: "smart label scan", dialect: ConfigurationDialectSmartLabelV2, sent: "4c01018200", hasActions: true},
		{name: "smart label clear storage", dialect: ConfigurationDialectSmartLabelV2, sent: "4c01018300", hasActions: true},
		{name: "smart label wipe all", dialect: ConfigurationDialectSmartLabelV2, sent: "4c01018400", hasActions: true},
		{name: "smart label alarm", dialect: ConfigurationDialectSmartLabelV2, sent: "4c010180020005", wantErr: errConfigurationUnsupportedCommand},
		{name: "unknown dialect", dialect: ConfigurationDialectUnspecified, sent: "4c01018100", wantErr: errConfigurationUnknownDialect},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			observable, hasActions, err := AnalyzeConfiguration(tc.dialect, tc.sent)
			if tc.wantErr != nil {
				assertConfigurationError(t, err, tc.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if observable != tc.observable || hasActions != tc.hasActions {
				t.Fatalf("got observable=%v actions=%v, want observable=%v actions=%v", observable, hasActions, tc.observable, tc.hasActions)
			}
		})
	}
}

func TestConfigurationActions(t *testing.T) {
	t.Run("wrong action lengths", func(t *testing.T) {
		cases := []struct {
			name    string
			dialect ConfigurationDialect
			sent    string
		}{
			{name: "tag xl alarm empty", dialect: ConfigurationDialectTagXL, sent: "4c01018000"},
			{name: "tag xl reset with value", dialect: ConfigurationDialectTagXL, sent: "4c010181020000"},
			{name: "tag xl scan with value", dialect: ConfigurationDialectTagXL, sent: "4c01018201ff"},
			{name: "tag xl wipe all with value", dialect: ConfigurationDialectTagXL, sent: "4c01018401ff"},
			{name: "smart label reset with value", dialect: ConfigurationDialectSmartLabelV2, sent: "4c01018101ff"},
			{name: "smart label wipe all with value", dialect: ConfigurationDialectSmartLabelV2, sent: "4c01018401ff"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				assertConfigurationError(t, ValidateConfiguration(tc.dialect, tc.sent), errConfigurationMalformedTLV)
			})
		}
	})

	t.Run("duplicate action", func(t *testing.T) {
		assertConfigurationError(t, ValidateConfiguration(ConfigurationDialectTagXL, "4c030281008100"), errConfigurationDuplicateTag)
		assertConfigurationError(t, ValidateConfiguration(ConfigurationDialectSmartLabelV2, "4c030282008200"), errConfigurationDuplicateTag)
	})

	t.Run("compare mixed setter and action", func(t *testing.T) {
		cases := []struct {
			name     string
			dialect  ConfigurationDialect
			sent     string
			observed string
		}{
			{name: "action after setter", dialect: ConfigurationDialectTagXL, sent: "4c0a022104012c1c208100", observed: "4c01014104012c1c20"},
			{name: "action before setter", dialect: ConfigurationDialectTagXL, sent: "4c0a0281002104012c1c20", observed: "4c01014104012c1c20"},
			{name: "action between setter and getter", dialect: ConfigurationDialectTagXL, sent: "4c0d032104012c1c2081004000", observed: "4c05024104012c1c2040010f"},
			{name: "smart label setter and reset", dialect: ConfigurationDialectSmartLabelV2, sent: "4c05022001108100", observed: "4c0101400110"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := CompareConfigurationFor(tc.dialect, tc.sent, tc.observed)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != ConfigurationMatch {
					t.Fatalf("got %s, want match", result)
				}
			})
		}
	})

	t.Run("observed scan acknowledgement is ignored", func(t *testing.T) {
		result, err := CompareConfigurationFor(ConfigurationDialectTagXL, "4c0a022104012c1c208200", "4c05024104012c1c208200")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != ConfigurationMatch {
			t.Fatalf("got %s, want match", result)
		}
	})

	t.Run("compare and build reject action-only", func(t *testing.T) {
		assertConfigurationError(t, mustCompareErr(t, ConfigurationDialectTagXL, "4c01018100", "4c010140010f"), errConfigurationNoSetter)
		assertConfigurationError(t, mustCompareErr(t, ConfigurationDialectSmartLabelV2, "4c01018200", "4c0101400110"), errConfigurationNoSetter)
		_, err := BuildCurrentConfigurationRequest(ConfigurationDialectTagXL, "4c01018100")
		assertConfigurationError(t, err, errConfigurationNoSetter)
	})

	t.Run("readback omits actions", func(t *testing.T) {
		got, err := BuildCurrentConfigurationRequest(ConfigurationDialectTagXL, "4c0d0320010a2104007807088100")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want, err := hex.DecodeString("4c050240004100")
		if err != nil {
			t.Fatalf("decode want: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("got %x, want %x", got, want)
		}
	})

	t.Run("actions count toward command limits", func(t *testing.T) {
		payload := []byte{0x4c, 0x00, 0x00}
		for range 31 {
			payload = append(payload, tlvTagDeviceFlags, 0x00)
		}
		payload = append(payload, actionTagResetDevice, 0x00)
		if err := ValidateConfiguration(ConfigurationDialectTagXL, hex.EncodeToString(payload)); err != nil {
			t.Fatalf("32 commands including one action should validate: %v", err)
		}
		payload = append(payload, actionTagScanNow, 0x00)
		assertConfigurationError(t, ValidateConfiguration(ConfigurationDialectTagXL, hex.EncodeToString(payload)), errConfigurationTooManyCommands)
	})

	t.Run("actions do not count toward response size", func(t *testing.T) {
		sent := hex.EncodeToString([]byte{
			0x4c, 0x00, 0x00,
			0x2c, 0x10, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
			0x30, 0x07, 0x00, 0x64, 0x00, 0x01, 0x02, 0x03, 0x04,
			0x21, 0x04, 0x01, 0x2c, 0x1c, 0x20,
			0x20, 0x01, 0x0f,
			0x23, 0x01, 0x06,
			0x46, 0x00,
			0x81, 0x00,
			0x82, 0x00,
			0x83, 0x00,
		})
		observable, hasActions, err := AnalyzeConfiguration(ConfigurationDialectSmartLabelV2, sent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !observable || !hasActions {
			t.Fatalf("got observable=%v actions=%v", observable, hasActions)
		}
	})
}

func mustCompareErr(t *testing.T, dialect ConfigurationDialect, sent, observed string) error {
	t.Helper()
	_, err := CompareConfigurationFor(dialect, sent, observed)
	return err
}

func TestSmartLabelCommandMatrix(t *testing.T) {
	cases := []struct {
		name     string
		sent     string
		observed string
	}{
		{name: "flags 20", sent: "4c010120011f", observed: "4c010140011f"},
		{name: "intervals 21", sent: "4c01012104012c1c20", observed: "4c01014104012c1c20"},
		{name: "heartbeat 23", sent: "4c0101230106", observed: "4c0101430106"},
		{name: "advertisement 24", sent: "4c010124013c", observed: "4c010144013c"},
		{name: "data rate 28", sent: "4c0101280102", observed: "4c01014e0102"},
		{name: "battery 45", sent: "4c01014500", observed: "4c010145020fa0"},
		{name: "firmware 46", sent: "4c01014600", observed: "4c0101460401020304"},
		{name: "reset count 49", sent: "4c01014900", observed: "4c01014902000c"},
		{name: "reset cause 4a", sent: "4c01014a00", observed: "4c01014a0400000001"},
		{name: "scan counts 4b", sent: "4c01014b00", observed: "4c01014b0400010002"},
		{name: "battery levels 29", sent: "4c01012902050a", observed: "4c01014c02050a"},
		{name: "wifi ap 2a", sent: "4c01012a0103", observed: "4c01014d0103"},
		{name: "batch 2b", sent: "4c01012b0104", observed: "4c0101500104"},
		{name: "ble 2c", sent: "4c01012c10051e0afec10000000000000000000000", observed: "4c01015110051e0afec10000000000000000000000"},
		{name: "motion v2 30", sent: "4c0101300700640102030405", observed: "4c0101600700640102030405"},
		{name: "temp 61", sent: "4c01016100", observed: "4c0101610119"},
		{name: "buffer 63", sent: "4c01016300", observed: "4c01016302000a"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CompareConfigurationFor(ConfigurationDialectSmartLabelV2, tc.sent, tc.observed)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != ConfigurationMatch {
				t.Fatalf("got %s, want match", result)
			}
		})
	}
}
