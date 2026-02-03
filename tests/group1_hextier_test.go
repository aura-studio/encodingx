package encodingx_test

import (
	"strings"
	"testing"

	"github.com/aura-studio/encodingx"
)

// ============================================================================
// HexTier 编码器单元测试
// ============================================================================

// TestHexTierMarshalUnmarshal 测试 HexTier 编码/解码往返
func TestHexTierMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"tiny", "hi"},
		{"small", "hello world"},
		{"medium", strings.Repeat("x", 20)},
		{"large", strings.Repeat("x", 50)},
		{"xlarge", strings.Repeat("x", 100)},
		{"xxlarge", strings.Repeat("x", 500)},
		{"json-like", `{"code":0,"msg":"ok","url":"https://example.com/path"}`},
	}

	enc := encodingx.NewHexTier()

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
			if !isPowerOfTwo(byteLen) || byteLen < 8 {
				t.Errorf("output byte length %d is not a valid tier", byteLen)
			}

			t.Logf("input len=%d, hex output len=%d", len(tt.input), hexLen)

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

// TestHexTierWithBytes 测试 Bytes 类型编码
func TestHexTierWithBytes(t *testing.T) {
	enc := encodingx.NewHexTier()
	input := encodingx.Bytes{Data: []byte("test data")}

	encoded, err := enc.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal Bytes failed: %v", err)
	}

	var result encodingx.Bytes
	err = enc.Unmarshal(encoded, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if string(result.Data) != string(input.Data) {
		t.Errorf("got %q, want %q", string(result.Data), string(input.Data))
	}
}

// TestHexTierWithBytesPointer 测试 *Bytes 类型编码
func TestHexTierWithBytesPointer(t *testing.T) {
	enc := encodingx.NewHexTier()
	input := encodingx.NewBytes()
	input.Data = []byte("test data pointer")

	encoded, err := enc.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal *Bytes failed: %v", err)
	}

	result := encodingx.NewBytes()
	err = enc.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if string(result.Data) != string(input.Data) {
		t.Errorf("got %q, want %q", string(result.Data), string(input.Data))
	}
}

// TestHexTierMarshalWrongType 测试非字节类型返回错误
func TestHexTierMarshalWrongType(t *testing.T) {
	enc := encodingx.NewHexTier()

	_, err := enc.Marshal(123)
	if err != encodingx.ErrHexWrongValueType {
		t.Errorf("expected ErrHexWrongValueType, got %v", err)
	}

	_, err = enc.Marshal("string")
	if err != encodingx.ErrHexWrongValueType {
		t.Errorf("expected ErrHexWrongValueType, got %v", err)
	}

	_, err = enc.Marshal(map[string]int{"key": 1})
	if err != encodingx.ErrHexWrongValueType {
		t.Errorf("expected ErrHexWrongValueType, got %v", err)
	}
}

// TestHexTierUnmarshalInvalidHex 测试无效十六进制字符串
func TestHexTierUnmarshalInvalidHex(t *testing.T) {
	enc := encodingx.NewHexTier()
	var result encodingx.Bytes

	err := enc.Unmarshal([]byte("not-hex!"), &result)
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

// TestHexTierUnmarshalInvalidTierSize 测试无效 tier 大小
func TestHexTierUnmarshalInvalidTierSize(t *testing.T) {
	enc := encodingx.NewHexTier()
	var result encodingx.Bytes

	// 6 bytes (12 hex chars) is not a valid tier
	err := enc.Unmarshal([]byte("000000000000"), &result)
	if err != encodingx.ErrHexInvalidData {
		t.Errorf("expected ErrHexInvalidData, got %v", err)
	}
}

// TestHexTierGradient 测试梯度特性
func TestHexTierGradient(t *testing.T) {
	enc := encodingx.NewHexTier()

	// Test that larger inputs produce larger or equal tiers
	prevLen := 0
	for _, size := range []int{1, 10, 50, 200, 1000} {
		input := strings.Repeat("a", size)
		encoded, err := enc.Marshal([]byte(input))
		if err != nil {
			t.Fatalf("Marshal failed for size %d: %v", size, err)
		}
		if len(encoded) < prevLen {
			t.Errorf("tier should not shrink: size %d produced %d, prev was %d", size, len(encoded), prevLen)
		}
		prevLen = len(encoded)
	}
}

// TestHexTierEmptyData 测试空数据编码/解码
func TestHexTierEmptyData(t *testing.T) {
	enc := encodingx.NewHexTier()

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

// TestHexTierString 测试 String() 方法
func TestHexTierString(t *testing.T) {
	enc := encodingx.NewHexTier()
	name := enc.String()

	if name != "HexTier" {
		t.Errorf("String() should return 'HexTier', got '%s'", name)
	}
}

// TestHexTierStyle 测试 Style() 方法
func TestHexTierStyle(t *testing.T) {
	enc := encodingx.NewHexTier()
	style := enc.Style()

	if style != encodingx.EncodingStyleBytes {
		t.Errorf("Style() should return EncodingStyleBytes, got %v", style)
	}
}

// TestHexTierReverse 测试 Reverse() 方法
func TestHexTierReverse(t *testing.T) {
	enc := encodingx.NewHexTier()
	reversed := enc.Reverse()

	if reversed.String() != enc.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}
}

// TestHexTierImplementsEncoding 测试实现 Encoding 接口
func TestHexTierImplementsEncoding(t *testing.T) {
	var enc encodingx.Encoding = encodingx.NewHexTier()

	if enc.String() == "" {
		t.Error("String() should return non-empty string")
	}

	style := enc.Style()
	if style != encodingx.EncodingStyleStruct &&
		style != encodingx.EncodingStyleBytes &&
		style != encodingx.EncodingStyleMix {
		t.Errorf("Style() returned invalid EncodingStyleType: %v", style)
	}

	reversed := enc.Reverse()
	if reversed == nil {
		t.Error("Reverse() should return non-nil Encoding")
	}
}

// TestHexTierLargeData 测试大数据编码/解码
func TestHexTierLargeData(t *testing.T) {
	enc := encodingx.NewHexTier()

	// 10KB data
	originalData := make([]byte, 10*1024)
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	encoded, err := enc.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal large data failed: %v", err)
	}

	var result encodingx.Bytes
	err = enc.Unmarshal(encoded, &result)
	if err != nil {
		t.Fatalf("Unmarshal large data failed: %v", err)
	}

	if !BytesEqual(result.Data, originalData) {
		t.Error("Large data round trip failed")
	}
}

// isPowerOfTwo checks if n is a power of 2
func isPowerOfTwo(n int) bool {
	return n > 0 && n&(n-1) == 0
}
