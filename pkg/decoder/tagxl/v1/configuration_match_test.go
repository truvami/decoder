package tagxl

import (
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

		_, err = MatchConfiguration("4c0401230118", "4c0602430118430118")
		assertConfigurationError(t, err, errConfigurationDuplicateTag)
	})

	t.Run("unsupported setter", func(t *testing.T) {
		_, err := MatchConfiguration("4c01014000", "4c010140010f")
		assertConfigurationError(t, err, errConfigurationUnsupportedCommand)

		_, err = MatchConfiguration("4c010180020005", "4c0101430118")
		assertConfigurationError(t, err, errConfigurationUnsupportedCommand)
	})

	t.Run("invalid data rate", func(t *testing.T) {
		_, err := MatchConfiguration("4c0101280108", "4c01014e0108")
		assertConfigurationError(t, err, errConfigurationInvalidDataRate)

		_, err = MatchConfiguration("4c0101280102", "4c01014e0108")
		assertConfigurationError(t, err, errConfigurationInvalidDataRate)
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
	})
}

func TestPort151PayloadConfigCoversComparableReportTags(t *testing.T) {
	tags := make(map[uint8]struct{})
	for _, tag := range port151PayloadConfig().Tags {
		tags[tag.Tag] = struct{}{}
	}
	for setter, spec := range configurationSetterSpecs {
		if _, ok := tags[spec.reportTag]; !ok {
			t.Errorf("report tag 0x%02x (setter 0x%02x) missing from port151PayloadConfig", spec.reportTag, setter)
		}
	}
}

func assertConfigurationError(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}
