package encodingx_test

import (
	"strings"
	"testing"

	"github.com/aura-studio/encodingx"
)

// ============================================================================
// HexTierRand 编码器单元测试
// ============================================================================

// TestHexTierRandMarshalUnmarshal 测试 HexTierRand 编码/解码往返
func TestHexTierRandMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"tiny", "hi"},
		{"portal", "portal"},
		{"small", "hello world"},
		{"medium", strings.Repeat("x", 20)},
		{"large", strings.Repeat("x", 50)},
		{"xlarge", strings.Repeat("x", 100)},
		{"json-like", `{"code":0,"msg":"ok","url":"https://example.com/path"}`},
	}

	enc := encodingx.NewHexTierRand()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			encoded, err := enc.Marshal([]byte(tt.input))
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// Verify hex output length is power of 2 * 2 (hex doubles bytes)
			hexLen := len(encoded)
			byteLen := hexLen / 2
			if !isPowerOfTwo(byteLen) || byteLen < 16 {
				t.Errorf("output byte length %d is not a valid tier (min 16)", byteLen)
			}

			t.Logf("input=%q len=%d, hex output=%s len=%d", tt.input, len(tt.input), string(encoded), hexLen)

			// Unmarshal
			var result encodingx.Bytes
			err = enc.Unmarshal(encoded, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Verify roundtrip
			if string(result.Data) != tt.input {
				t.Errorf("roundtrip failed: got %q, want %q", string(result.Data), tt.input)
			}
		})
	}
}

// TestHexTierRandRandomness 测试每次编码结果不同
func TestHexTierRandRandomness(t *testing.T) {
	enc := encodingx.NewHexTierRand()
	input := []byte("portal")

	results := make(map[string]bool)
	for i := 0; i < 10; i++ {
		encoded, err := enc.Marshal(input)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		results[string(encoded)] = true
	}

	// Should have multiple different outputs
	if len(results) < 5 {
		t.Errorf("expected more randomness, got only %d unique outputs", len(results))
	}
	t.Logf("got %d unique outputs from 10 encodings", len(results))
}

// TestHexTierRandEmptyData 测试空数据编码/解码
func TestHexTierRandEmptyData(t *testing.T) {
	enc := encodingx.NewHexTierRand()

	encoded, err := enc.Marshal([]byte{})
	if err != nil {
		t.Fatalf("Marshal empty data failed: %v", err)
	}

	var result encodingx.Bytes
	err = enc.Unmarshal(encoded, &result)
	if err != nil {
		t.Fatalf("Unmarshal empty data failed: %v", err)
	}

	if len(result.Data) != 0 {
		t.Errorf("expected empty data, got %v", result.Data)
	}
}

// TestHexTierRandString 测试 String() 方法
func TestHexTierRandString(t *testing.T) {
	enc := encodingx.NewHexTierRand()
	name := enc.String()

	if name != "HexTierRand" {
		t.Errorf("String() should return 'HexTierRand', got '%s'", name)
	}
}
