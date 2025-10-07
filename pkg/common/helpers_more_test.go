package common

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestValidateLength_WithTags(t *testing.T) {
	payload := "0000" // 2 bytes (too short; min is 3 when Tags present)
	cfg := &PayloadConfig{
		Tags: []TagConfig{
			{Name: "X", Tag: 0x01, Optional: true},
		},
	}

	if err := ValidateLength(&payload, cfg); err == nil || !errors.Is(err, ErrInvalidPayloadLength) {
		t.Fatalf("expected ErrInvalidPayloadLength on too short with Tags, got %v", err)
	}

	ok := "000000" // 3 bytes, ok
	if err := ValidateLength(&ok, cfg); err != nil {
		t.Fatalf("unexpected error for valid length with Tags: %v", err)
	}
}

func TestValidateLength_WithFields(t *testing.T) {
	cfg := &PayloadConfig{
		Fields: []FieldConfig{
			// one required field (min length = 2)
			{Name: "A", Start: 0, Length: 2, Optional: false},
			// one optional field (max length ends up at 4)
			{Name: "B", Start: 2, Length: 2, Optional: true},
		},
	}

	short := "00" // 1 byte
	if err := ValidateLength(&short, cfg); err == nil || !errors.Is(err, ErrInvalidPayloadLength) {
		t.Fatalf("expected ErrInvalidPayloadLength on too short with Fields, got %v", err)
	}

	ok := "00000000" // 4 bytes
	if err := ValidateLength(&ok, cfg); err != nil {
		t.Fatalf("unexpected error for valid length with Fields: %v", err)
	}

	long := "0000000000" // 5 bytes
	if err := ValidateLength(&long, cfg); err == nil || !errors.Is(err, ErrInvalidPayloadLength) {
		t.Fatalf("expected ErrInvalidPayloadLength on too long with Fields, got %v", err)
	}
}

func TestBoolToBytesAndBytesToBool(t *testing.T) {
	if got := BoolToBytes(true, 0); len(got) != 1 || got[0] != 0x01 {
		t.Fatalf("BoolToBytes(true,0) = %v want [1]", got)
	}
	if got := BoolToBytes(false, 0); len(got) != 1 || got[0] != 0x00 {
		t.Fatalf("BoolToBytes(false,0) = %v want [0]", got)
	}
	// bit must be 0..7
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on invalid bit > 7")
		}
	}()
	_ = BoolToBytes(true, 8)
}

func TestUintAndIntToBytes(t *testing.T) {
	if got := UintToBytes(0x1234, 2); len(got) != 2 || got[0] != 0x12 || got[1] != 0x34 {
		t.Fatalf("UintToBytes(0x1234,2)=%x want 12 34", got)
	}
	if got := IntToBytes(0x0102, 2); len(got) != 2 || got[0] != 0x01 || got[1] != 0x02 {
		t.Fatalf("IntToBytes(0x0102,2)=%x want 01 02", got)
	}
}

func TestFloatRoundTrip(t *testing.T) {
	f32 := float32(123.456)
	b32 := Float32ToBytes(f32)
	if got := BytesToFloat32(b32); math.Abs(float64(got-f32)) > 1e-5 {
		t.Fatalf("float32 roundtrip got %f want %f", got, f32)
	}

	f64 := float64(789.12345)
	b64 := Float64ToBytes(f64)
	if got := BytesToFloat64(b64); math.Abs(got-f64) > 1e-9 {
		t.Fatalf("float64 roundtrip got %f want %f", got, f64)
	}
}

func TestTimePointerAndCompare(t *testing.T) {
	ts := 1234.5
	tp := TimePointer(ts)
	if tp == nil {
		t.Fatalf("TimePointer returned nil")
	}
	// Compare same pointer
	if !TimePointerCompare(tp, tp) {
		t.Fatalf("expected same pointers to be equal")
	}
	// Compare nil cases
	if !TimePointerCompare(nil, nil) {
		t.Fatalf("expected nil,nil to be equal")
	}
	if TimePointerCompare(tp, nil) {
		t.Fatalf("expected tp vs nil to be not equal")
	}
}

func TestEncode_SimpleStruct(t *testing.T) {
	type Data struct {
		A uint16
		B string
		C bool
	}
	cfg := PayloadConfig{
		Fields: []FieldConfig{
			{Name: "A", Start: 0, Length: 2},
			{Name: "B", Start: 2, Length: 2},
			{Name: "C", Start: 4, Length: 1},
		},
		TargetType: reflect.TypeOf(Data{}),
	}

	d := Data{
		A: 0x0102,
		B: "0A0B",
		C: true,
	}
	got, err := Encode(d, cfg)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	if got != "01020a0b01" {
		t.Fatalf("Encode got %q want 01020a0b01", got)
	}
}

func TestPointerHelpers(t *testing.T) {
	if v := Uint8Ptr(7); v == nil || *v != 7 {
		t.Fatalf("Uint8Ptr failed")
	}
	if v := Uint16Ptr(9); v == nil || *v != 9 {
		t.Fatalf("Uint16Ptr failed")
	}
	if v := Uint32Ptr(11); v == nil || *v != 11 {
		t.Fatalf("Uint32Ptr failed")
	}
	if v := Int8Ptr(-3); v == nil || *v != -3 {
		t.Fatalf("Int8Ptr failed")
	}
	if v := BoolPtr(true); v == nil || *v != true {
		t.Fatalf("BoolPtr failed")
	}
	if v := StringPtr("x"); v == nil || *v != "x" {
		t.Fatalf("StringPtr failed")
	}
	if v := Float32Ptr(1.5); v == nil || *v != 1.5 {
		t.Fatalf("Float32Ptr failed")
	}
	if v := Float64Ptr(2.5); v == nil || *v != 2.5 {
		t.Fatalf("Float64Ptr failed")
	}
	if v := DurationPtr(3 * time.Second); v == nil || *v != 3*time.Second {
		t.Fatalf("DurationPtr failed")
	}
}

func TestDerefValue(t *testing.T) {
	val := reflect.ValueOf(123)
	// For non-pointer reflect.Value input, DerefValue returns the same reflect.Value.
	if got := DerefValue(val); got != val {
		t.Fatalf("DerefValue non-ptr should return same reflect.Value")
	}

	x := 42
	ptrVal := reflect.ValueOf(&x)
	if got := DerefValue(ptrVal); got.(int) != 42 {
		t.Fatalf("DerefValue ptr failed: %v", got)
	}

	var nilPtr *int
	nilVal := reflect.ValueOf(nilPtr)
	if got := DerefValue(nilVal); got != nilVal {
		t.Fatalf("DerefValue nil ptr should return same reflect.Value")
	}
}
