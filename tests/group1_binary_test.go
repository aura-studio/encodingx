package encodingx_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// Binary 编码器单元测试
// Validates: Requirements 7.1, 7.2, 7.3, 7.4, 14.1, 14.2, 14.3
// ============================================================================

// ============================================================================
// Binary 编码器测试（小端序）
// ============================================================================

// TestBinaryMarshalInt32 测试 Binary 编码器序列化 int32
// Validates: Requirements 7.1
func TestBinaryMarshalInt32(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := int32(0x12345678)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 4)
	binary.LittleEndian.PutUint32(expected, uint32(original))
	if !BytesEqual(data, expected) {
		t.Errorf("Binary encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBinaryMarshalInt64 测试 Binary 编码器序列化 int64
// Validates: Requirements 7.1
func TestBinaryMarshalInt64(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := int64(0x123456789ABCDEF0)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 8)
	binary.LittleEndian.PutUint64(expected, uint64(original))
	if !BytesEqual(data, expected) {
		t.Errorf("Binary encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBinaryMarshalFloat32 测试 Binary 编码器序列化 float32
// Validates: Requirements 7.1
func TestBinaryMarshalFloat32(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := float32(3.14159)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 4)
	binary.LittleEndian.PutUint32(expected, math.Float32bits(original))
	if !BytesEqual(data, expected) {
		t.Errorf("Binary encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBinaryMarshalFloat64 测试 Binary 编码器序列化 float64
// Validates: Requirements 7.1
func TestBinaryMarshalFloat64(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := float64(3.141592653589793)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 8)
	binary.LittleEndian.PutUint64(expected, math.Float64bits(original))
	if !BytesEqual(data, expected) {
		t.Errorf("Binary encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBinaryUnmarshalInt32 测试 Binary 编码器反序列化 int32
// Validates: Requirements 7.2
func TestBinaryUnmarshalInt32(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := int32(0x12345678)

	// 创建小端序二进制数据
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(original))

	// 反序列化
	var result int32
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("Binary decoding mismatch: expected %d, got %d", original, result)
	}
}

// TestBinaryUnmarshalInt64 测试 Binary 编码器反序列化 int64
// Validates: Requirements 7.2
func TestBinaryUnmarshalInt64(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := int64(0x123456789ABCDEF0)

	// 创建小端序二进制数据
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(original))

	// 反序列化
	var result int64
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("Binary decoding mismatch: expected %d, got %d", original, result)
	}
}

// TestBinaryUnmarshalFloat32 测试 Binary 编码器反序列化 float32
// Validates: Requirements 7.2
func TestBinaryUnmarshalFloat32(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := float32(3.14159)

	// 创建小端序二进制数据
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, math.Float32bits(original))

	// 反序列化
	var result float32
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("Binary decoding mismatch: expected %f, got %f", original, result)
	}
}

// TestBinaryUnmarshalFloat64 测试 Binary 编码器反序列化 float64
// Validates: Requirements 7.2
func TestBinaryUnmarshalFloat64(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := float64(3.141592653589793)

	// 创建小端序二进制数据
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, math.Float64bits(original))

	// 反序列化
	var result float64
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("Binary decoding mismatch: expected %f, got %f", original, result)
	}
}

// TestBinaryRoundTripInt32 测试 Binary 编码器 int32 往返
// Validates: Requirements 7.1, 7.2
func TestBinaryRoundTripInt32(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []int32{0, 1, -1, 127, -128, 32767, -32768, 2147483647, -2147483648}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryRoundTripInt64 测试 Binary 编码器 int64 往返
// Validates: Requirements 7.1, 7.2
func TestBinaryRoundTripInt64(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []int64{0, 1, -1, 9223372036854775807, -9223372036854775808}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryRoundTripFloat32 测试 Binary 编码器 float32 往返
// Validates: Requirements 7.1, 7.2
func TestBinaryRoundTripFloat32(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []float32{0, 1.0, -1.0, 3.14159, -3.14159, math.MaxFloat32, math.SmallestNonzeroFloat32}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %f: %v", original, err)
		}

		// 反序列化
		var result float32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %f: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %f, got %f", original, result)
		}
	}
}

// TestBinaryRoundTripFloat64 测试 Binary 编码器 float64 往返
// Validates: Requirements 7.1, 7.2
func TestBinaryRoundTripFloat64(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []float64{0, 1.0, -1.0, 3.141592653589793, -3.141592653589793, math.MaxFloat64, math.SmallestNonzeroFloat64}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %f: %v", original, err)
		}

		// 反序列化
		var result float64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %f: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %f, got %f", original, result)
		}
	}
}

// TestBinaryString 测试 Binary String() 方法返回类型名称
// Validates: Requirements 14.1
func TestBinaryString(t *testing.T) {
	encoder := encodingx.NewBinary()
	name := encoder.String()

	if name != "Binary" {
		t.Errorf("String() should return 'Binary', got '%s'", name)
	}
}

// TestBinaryStyle 测试 Binary Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestBinaryStyle(t *testing.T) {
	encoder := encodingx.NewBinary()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestBinaryReverse 测试 Binary Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestBinaryReverse(t *testing.T) {
	encoder := encodingx.NewBinary()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 Binary 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestBinaryImplementsEncoding 测试 Binary 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestBinaryImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewBinary()

	// 验证接口方法
	if encoder.String() == "" {
		t.Error("String() should return non-empty string")
	}

	// Style() 应该返回有效的 EncodingStyleType
	style := encoder.Style()
	if style != encodingx.EncodingStyleStruct &&
		style != encodingx.EncodingStyleBytes &&
		style != encodingx.EncodingStyleMix {
		t.Errorf("Style() returned invalid EncodingStyleType: %v", style)
	}

	// Reverse() 应该返回非 nil 的 Encoding
	reversed := encoder.Reverse()
	if reversed == nil {
		t.Error("Reverse() should return non-nil Encoding")
	}
}

// ============================================================================
// LittleEndian 编码器测试
// ============================================================================

// TestLittleEndianMarshalInt32 测试 LittleEndian 编码器序列化 int32
// Validates: Requirements 7.3
func TestLittleEndianMarshalInt32(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := int32(0x12345678)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 4)
	binary.LittleEndian.PutUint32(expected, uint32(original))
	if !BytesEqual(data, expected) {
		t.Errorf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestLittleEndianMarshalInt64 测试 LittleEndian 编码器序列化 int64
// Validates: Requirements 7.3
func TestLittleEndianMarshalInt64(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := int64(0x123456789ABCDEF0)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 8)
	binary.LittleEndian.PutUint64(expected, uint64(original))
	if !BytesEqual(data, expected) {
		t.Errorf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestLittleEndianMarshalFloat32 测试 LittleEndian 编码器序列化 float32
// Validates: Requirements 7.3
func TestLittleEndianMarshalFloat32(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := float32(3.14159)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 4)
	binary.LittleEndian.PutUint32(expected, math.Float32bits(original))
	if !BytesEqual(data, expected) {
		t.Errorf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestLittleEndianMarshalFloat64 测试 LittleEndian 编码器序列化 float64
// Validates: Requirements 7.3
func TestLittleEndianMarshalFloat64(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := float64(3.141592653589793)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是小端序二进制数据
	expected := make([]byte, 8)
	binary.LittleEndian.PutUint64(expected, math.Float64bits(original))
	if !BytesEqual(data, expected) {
		t.Errorf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestLittleEndianUnmarshalInt32 测试 LittleEndian 编码器反序列化 int32
// Validates: Requirements 7.3
func TestLittleEndianUnmarshalInt32(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := int32(0x12345678)

	// 创建小端序二进制数据
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(original))

	// 反序列化
	var result int32
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("LittleEndian decoding mismatch: expected %d, got %d", original, result)
	}
}

// TestLittleEndianUnmarshalInt64 测试 LittleEndian 编码器反序列化 int64
// Validates: Requirements 7.3
func TestLittleEndianUnmarshalInt64(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := int64(0x123456789ABCDEF0)

	// 创建小端序二进制数据
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(original))

	// 反序列化
	var result int64
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("LittleEndian decoding mismatch: expected %d, got %d", original, result)
	}
}

// TestLittleEndianUnmarshalFloat32 测试 LittleEndian 编码器反序列化 float32
// Validates: Requirements 7.3
func TestLittleEndianUnmarshalFloat32(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := float32(3.14159)

	// 创建小端序二进制数据
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, math.Float32bits(original))

	// 反序列化
	var result float32
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("LittleEndian decoding mismatch: expected %f, got %f", original, result)
	}
}

// TestLittleEndianUnmarshalFloat64 测试 LittleEndian 编码器反序列化 float64
// Validates: Requirements 7.3
func TestLittleEndianUnmarshalFloat64(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := float64(3.141592653589793)

	// 创建小端序二进制数据
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, math.Float64bits(original))

	// 反序列化
	var result float64
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("LittleEndian decoding mismatch: expected %f, got %f", original, result)
	}
}

// TestLittleEndianRoundTripInt32 测试 LittleEndian 编码器 int32 往返
// Validates: Requirements 7.3
func TestLittleEndianRoundTripInt32(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	testCases := []int32{0, 1, -1, 127, -128, 32767, -32768, 2147483647, -2147483648}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestLittleEndianRoundTripInt64 测试 LittleEndian 编码器 int64 往返
// Validates: Requirements 7.3
func TestLittleEndianRoundTripInt64(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	testCases := []int64{0, 1, -1, 9223372036854775807, -9223372036854775808}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestLittleEndianRoundTripFloat32 测试 LittleEndian 编码器 float32 往返
// Validates: Requirements 7.3
func TestLittleEndianRoundTripFloat32(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	testCases := []float32{0, 1.0, -1.0, 3.14159, -3.14159, math.MaxFloat32, math.SmallestNonzeroFloat32}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %f: %v", original, err)
		}

		// 反序列化
		var result float32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %f: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %f, got %f", original, result)
		}
	}
}

// TestLittleEndianRoundTripFloat64 测试 LittleEndian 编码器 float64 往返
// Validates: Requirements 7.3
func TestLittleEndianRoundTripFloat64(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	testCases := []float64{0, 1.0, -1.0, 3.141592653589793, -3.141592653589793, math.MaxFloat64, math.SmallestNonzeroFloat64}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %f: %v", original, err)
		}

		// 反序列化
		var result float64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %f: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %f, got %f", original, result)
		}
	}
}

// TestLittleEndianString 测试 LittleEndian String() 方法返回类型名称
// Validates: Requirements 14.1
func TestLittleEndianString(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	name := encoder.String()

	if name != "LittleEndian" {
		t.Errorf("String() should return 'LittleEndian', got '%s'", name)
	}
}

// TestLittleEndianStyle 测试 LittleEndian Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestLittleEndianStyle(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestLittleEndianReverse 测试 LittleEndian Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestLittleEndianReverse(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 LittleEndian 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestLittleEndianImplementsEncoding 测试 LittleEndian 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestLittleEndianImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewLittleEndian()

	// 验证接口方法
	if encoder.String() == "" {
		t.Error("String() should return non-empty string")
	}

	// Style() 应该返回有效的 EncodingStyleType
	style := encoder.Style()
	if style != encodingx.EncodingStyleStruct &&
		style != encodingx.EncodingStyleBytes &&
		style != encodingx.EncodingStyleMix {
		t.Errorf("Style() returned invalid EncodingStyleType: %v", style)
	}

	// Reverse() 应该返回非 nil 的 Encoding
	reversed := encoder.Reverse()
	if reversed == nil {
		t.Error("Reverse() should return non-nil Encoding")
	}
}

// ============================================================================
// BigEndian 编码器测试
// ============================================================================

// TestBigEndianMarshalInt32 测试 BigEndian 编码器序列化 int32
// Validates: Requirements 7.4
func TestBigEndianMarshalInt32(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := int32(0x12345678)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是大端序二进制数据
	expected := make([]byte, 4)
	binary.BigEndian.PutUint32(expected, uint32(original))
	if !BytesEqual(data, expected) {
		t.Errorf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBigEndianMarshalInt64 测试 BigEndian 编码器序列化 int64
// Validates: Requirements 7.4
func TestBigEndianMarshalInt64(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := int64(0x123456789ABCDEF0)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是大端序二进制数据
	expected := make([]byte, 8)
	binary.BigEndian.PutUint64(expected, uint64(original))
	if !BytesEqual(data, expected) {
		t.Errorf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBigEndianMarshalFloat32 测试 BigEndian 编码器序列化 float32
// Validates: Requirements 7.4
func TestBigEndianMarshalFloat32(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := float32(3.14159)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是大端序二进制数据
	expected := make([]byte, 4)
	binary.BigEndian.PutUint32(expected, math.Float32bits(original))
	if !BytesEqual(data, expected) {
		t.Errorf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBigEndianMarshalFloat64 测试 BigEndian 编码器序列化 float64
// Validates: Requirements 7.4
func TestBigEndianMarshalFloat64(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := float64(3.141592653589793)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是大端序二进制数据
	expected := make([]byte, 8)
	binary.BigEndian.PutUint64(expected, math.Float64bits(original))
	if !BytesEqual(data, expected) {
		t.Errorf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
	}
}

// TestBigEndianUnmarshalInt32 测试 BigEndian 编码器反序列化 int32
// Validates: Requirements 7.4
func TestBigEndianUnmarshalInt32(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := int32(0x12345678)

	// 创建大端序二进制数据
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(original))

	// 反序列化
	var result int32
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("BigEndian decoding mismatch: expected %d, got %d", original, result)
	}
}

// TestBigEndianUnmarshalInt64 测试 BigEndian 编码器反序列化 int64
// Validates: Requirements 7.4
func TestBigEndianUnmarshalInt64(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := int64(0x123456789ABCDEF0)

	// 创建大端序二进制数据
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(original))

	// 反序列化
	var result int64
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("BigEndian decoding mismatch: expected %d, got %d", original, result)
	}
}

// TestBigEndianUnmarshalFloat32 测试 BigEndian 编码器反序列化 float32
// Validates: Requirements 7.4
func TestBigEndianUnmarshalFloat32(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := float32(3.14159)

	// 创建大端序二进制数据
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, math.Float32bits(original))

	// 反序列化
	var result float32
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("BigEndian decoding mismatch: expected %f, got %f", original, result)
	}
}

// TestBigEndianUnmarshalFloat64 测试 BigEndian 编码器反序列化 float64
// Validates: Requirements 7.4
func TestBigEndianUnmarshalFloat64(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := float64(3.141592653589793)

	// 创建大端序二进制数据
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, math.Float64bits(original))

	// 反序列化
	var result float64
	err := encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确还原
	if result != original {
		t.Errorf("BigEndian decoding mismatch: expected %f, got %f", original, result)
	}
}

// TestBigEndianRoundTripInt32 测试 BigEndian 编码器 int32 往返
// Validates: Requirements 7.4
func TestBigEndianRoundTripInt32(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	testCases := []int32{0, 1, -1, 127, -128, 32767, -32768, 2147483647, -2147483648}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBigEndianRoundTripInt64 测试 BigEndian 编码器 int64 往返
// Validates: Requirements 7.4
func TestBigEndianRoundTripInt64(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	testCases := []int64{0, 1, -1, 9223372036854775807, -9223372036854775808}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBigEndianRoundTripFloat32 测试 BigEndian 编码器 float32 往返
// Validates: Requirements 7.4
func TestBigEndianRoundTripFloat32(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	testCases := []float32{0, 1.0, -1.0, 3.14159, -3.14159, math.MaxFloat32, math.SmallestNonzeroFloat32}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %f: %v", original, err)
		}

		// 反序列化
		var result float32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %f: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %f, got %f", original, result)
		}
	}
}

// TestBigEndianRoundTripFloat64 测试 BigEndian 编码器 float64 往返
// Validates: Requirements 7.4
func TestBigEndianRoundTripFloat64(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	testCases := []float64{0, 1.0, -1.0, 3.141592653589793, -3.141592653589793, math.MaxFloat64, math.SmallestNonzeroFloat64}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %f: %v", original, err)
		}

		// 反序列化
		var result float64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %f: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %f, got %f", original, result)
		}
	}
}

// TestBigEndianString 测试 BigEndian String() 方法返回类型名称
// Validates: Requirements 14.1
func TestBigEndianString(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	name := encoder.String()

	if name != "BigEndian" {
		t.Errorf("String() should return 'BigEndian', got '%s'", name)
	}
}

// TestBigEndianStyle 测试 BigEndian Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestBigEndianStyle(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestBigEndianReverse 测试 BigEndian Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestBigEndianReverse(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 BigEndian 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestBigEndianImplementsEncoding 测试 BigEndian 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestBigEndianImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewBigEndian()

	// 验证接口方法
	if encoder.String() == "" {
		t.Error("String() should return non-empty string")
	}

	// Style() 应该返回有效的 EncodingStyleType
	style := encoder.Style()
	if style != encodingx.EncodingStyleStruct &&
		style != encodingx.EncodingStyleBytes &&
		style != encodingx.EncodingStyleMix {
		t.Errorf("Style() returned invalid EncodingStyleType: %v", style)
	}

	// Reverse() 应该返回非 nil 的 Encoding
	reversed := encoder.Reverse()
	if reversed == nil {
		t.Error("Reverse() should return non-nil Encoding")
	}
}

// ============================================================================
// Binary vs LittleEndian vs BigEndian 对比测试
// ============================================================================

// TestBinaryVsLittleEndianSameResult 测试 Binary 和 LittleEndian 产生相同结果
// Validates: Requirements 7.1, 7.3
func TestBinaryVsLittleEndianSameResult(t *testing.T) {
	binaryEncoder := encodingx.NewBinary()
	littleEndianEncoder := encodingx.NewLittleEndian()

	testCases := []int32{0, 1, -1, 0x12345678, -0x12345678}

	for _, original := range testCases {
		// Binary 编码
		binaryData, err := binaryEncoder.Marshal(original)
		if err != nil {
			t.Fatalf("Binary Marshal failed for %d: %v", original, err)
		}

		// LittleEndian 编码
		littleEndianData, err := littleEndianEncoder.Marshal(original)
		if err != nil {
			t.Fatalf("LittleEndian Marshal failed for %d: %v", original, err)
		}

		// 验证两者结果相同
		if !BytesEqual(binaryData, littleEndianData) {
			t.Errorf("Binary and LittleEndian should produce same result for %d: Binary=%v, LittleEndian=%v",
				original, binaryData, littleEndianData)
		}
	}
}

// TestLittleEndianVsBigEndianDifferent 测试 LittleEndian 和 BigEndian 产生不同结果
// Validates: Requirements 7.3, 7.4
func TestLittleEndianVsBigEndianDifferent(t *testing.T) {
	littleEndianEncoder := encodingx.NewLittleEndian()
	bigEndianEncoder := encodingx.NewBigEndian()

	// 使用一个非对称的值来确保字节顺序不同
	original := int32(0x12345678)

	// LittleEndian 编码
	littleEndianData, err := littleEndianEncoder.Marshal(original)
	if err != nil {
		t.Fatalf("LittleEndian Marshal failed: %v", err)
	}

	// BigEndian 编码
	bigEndianData, err := bigEndianEncoder.Marshal(original)
	if err != nil {
		t.Fatalf("BigEndian Marshal failed: %v", err)
	}

	// 验证两者结果不同
	if BytesEqual(littleEndianData, bigEndianData) {
		t.Errorf("LittleEndian and BigEndian should produce different results for 0x12345678")
	}

	// 验证字节顺序
	// LittleEndian: 0x78, 0x56, 0x34, 0x12
	// BigEndian: 0x12, 0x34, 0x56, 0x78
	expectedLittleEndian := []byte{0x78, 0x56, 0x34, 0x12}
	expectedBigEndian := []byte{0x12, 0x34, 0x56, 0x78}

	if !BytesEqual(littleEndianData, expectedLittleEndian) {
		t.Errorf("LittleEndian encoding mismatch: expected %v, got %v", expectedLittleEndian, littleEndianData)
	}
	if !BytesEqual(bigEndianData, expectedBigEndian) {
		t.Errorf("BigEndian encoding mismatch: expected %v, got %v", expectedBigEndian, bigEndianData)
	}
}

// TestBigEndianByteOrder 测试 BigEndian 字节顺序正确
// Validates: Requirements 7.4
func TestBigEndianByteOrder(t *testing.T) {
	encoder := encodingx.NewBigEndian()

	// 测试 int32
	int32Val := int32(0x01020304)
	data, err := encoder.Marshal(int32Val)
	if err != nil {
		t.Fatalf("Marshal int32 failed: %v", err)
	}
	expectedInt32 := []byte{0x01, 0x02, 0x03, 0x04}
	if !BytesEqual(data, expectedInt32) {
		t.Errorf("BigEndian int32 byte order wrong: expected %v, got %v", expectedInt32, data)
	}

	// 测试 int64
	int64Val := int64(0x0102030405060708)
	data, err = encoder.Marshal(int64Val)
	if err != nil {
		t.Fatalf("Marshal int64 failed: %v", err)
	}
	expectedInt64 := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	if !BytesEqual(data, expectedInt64) {
		t.Errorf("BigEndian int64 byte order wrong: expected %v, got %v", expectedInt64, data)
	}
}

// TestLittleEndianByteOrder 测试 LittleEndian 字节顺序正确
// Validates: Requirements 7.3
func TestLittleEndianByteOrder(t *testing.T) {
	encoder := encodingx.NewLittleEndian()

	// 测试 int32
	int32Val := int32(0x01020304)
	data, err := encoder.Marshal(int32Val)
	if err != nil {
		t.Fatalf("Marshal int32 failed: %v", err)
	}
	expectedInt32 := []byte{0x04, 0x03, 0x02, 0x01}
	if !BytesEqual(data, expectedInt32) {
		t.Errorf("LittleEndian int32 byte order wrong: expected %v, got %v", expectedInt32, data)
	}

	// 测试 int64
	int64Val := int64(0x0102030405060708)
	data, err = encoder.Marshal(int64Val)
	if err != nil {
		t.Fatalf("Marshal int64 failed: %v", err)
	}
	expectedInt64 := []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}
	if !BytesEqual(data, expectedInt64) {
		t.Errorf("LittleEndian int64 byte order wrong: expected %v, got %v", expectedInt64, data)
	}
}

// ============================================================================
// 固定大小结构体测试
// ============================================================================

// TestBinaryMarshalFixedSizeStruct 测试 Binary 编码器序列化固定大小结构体
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalFixedSizeStruct(t *testing.T) {
	encoder := encodingx.NewBinary()
	original := FixedSizeStruct{
		Int32Val:   0x12345678,
		Int64Val:   0x123456789ABCDEF0,
		Float32Val: 3.14159,
		Float64Val: 2.718281828459045,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result FixedSizeStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestLittleEndianMarshalFixedSizeStruct 测试 LittleEndian 编码器序列化固定大小结构体
// Validates: Requirements 7.3
func TestLittleEndianMarshalFixedSizeStruct(t *testing.T) {
	encoder := encodingx.NewLittleEndian()
	original := FixedSizeStruct{
		Int32Val:   0x12345678,
		Int64Val:   0x123456789ABCDEF0,
		Float32Val: 3.14159,
		Float64Val: 2.718281828459045,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result FixedSizeStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestBigEndianMarshalFixedSizeStruct 测试 BigEndian 编码器序列化固定大小结构体
// Validates: Requirements 7.4
func TestBigEndianMarshalFixedSizeStruct(t *testing.T) {
	encoder := encodingx.NewBigEndian()
	original := FixedSizeStruct{
		Int32Val:   0x12345678,
		Int64Val:   0x123456789ABCDEF0,
		Float32Val: 3.14159,
		Float64Val: 2.718281828459045,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result FixedSizeStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// 边界条件测试
// ============================================================================

// TestBinaryMarshalZeroValues 测试零值序列化
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalZeroValues(t *testing.T) {
	encoder := encodingx.NewBinary()

	// 测试 int32 零值
	var int32Zero int32 = 0
	data, err := encoder.Marshal(int32Zero)
	if err != nil {
		t.Fatalf("Marshal int32 zero failed: %v", err)
	}
	var int32Result int32
	err = encoder.Unmarshal(data, &int32Result)
	if err != nil {
		t.Fatalf("Unmarshal int32 zero failed: %v", err)
	}
	if int32Result != 0 {
		t.Errorf("int32 zero round trip failed: expected 0, got %d", int32Result)
	}

	// 测试 float64 零值
	var float64Zero float64 = 0.0
	data, err = encoder.Marshal(float64Zero)
	if err != nil {
		t.Fatalf("Marshal float64 zero failed: %v", err)
	}
	var float64Result float64
	err = encoder.Unmarshal(data, &float64Result)
	if err != nil {
		t.Fatalf("Unmarshal float64 zero failed: %v", err)
	}
	if float64Result != 0.0 {
		t.Errorf("float64 zero round trip failed: expected 0.0, got %f", float64Result)
	}
}

// TestBinaryMarshalMaxValues 测试最大值序列化
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalMaxValues(t *testing.T) {
	encoder := encodingx.NewBinary()

	// 测试 int32 最大值
	int32Max := int32(2147483647)
	data, err := encoder.Marshal(int32Max)
	if err != nil {
		t.Fatalf("Marshal int32 max failed: %v", err)
	}
	var int32Result int32
	err = encoder.Unmarshal(data, &int32Result)
	if err != nil {
		t.Fatalf("Unmarshal int32 max failed: %v", err)
	}
	if int32Result != int32Max {
		t.Errorf("int32 max round trip failed: expected %d, got %d", int32Max, int32Result)
	}

	// 测试 int64 最大值
	int64Max := int64(9223372036854775807)
	data, err = encoder.Marshal(int64Max)
	if err != nil {
		t.Fatalf("Marshal int64 max failed: %v", err)
	}
	var int64Result int64
	err = encoder.Unmarshal(data, &int64Result)
	if err != nil {
		t.Fatalf("Unmarshal int64 max failed: %v", err)
	}
	if int64Result != int64Max {
		t.Errorf("int64 max round trip failed: expected %d, got %d", int64Max, int64Result)
	}
}

// TestBinaryMarshalMinValues 测试最小值序列化
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalMinValues(t *testing.T) {
	encoder := encodingx.NewBinary()

	// 测试 int32 最小值
	int32Min := int32(-2147483648)
	data, err := encoder.Marshal(int32Min)
	if err != nil {
		t.Fatalf("Marshal int32 min failed: %v", err)
	}
	var int32Result int32
	err = encoder.Unmarshal(data, &int32Result)
	if err != nil {
		t.Fatalf("Unmarshal int32 min failed: %v", err)
	}
	if int32Result != int32Min {
		t.Errorf("int32 min round trip failed: expected %d, got %d", int32Min, int32Result)
	}

	// 测试 int64 最小值
	int64Min := int64(-9223372036854775808)
	data, err = encoder.Marshal(int64Min)
	if err != nil {
		t.Fatalf("Marshal int64 min failed: %v", err)
	}
	var int64Result int64
	err = encoder.Unmarshal(data, &int64Result)
	if err != nil {
		t.Fatalf("Unmarshal int64 min failed: %v", err)
	}
	if int64Result != int64Min {
		t.Errorf("int64 min round trip failed: expected %d, got %d", int64Min, int64Result)
	}
}

// TestBinaryMarshalSpecialFloatValues 测试特殊浮点值序列化
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalSpecialFloatValues(t *testing.T) {
	encoder := encodingx.NewBinary()

	// 测试正无穷
	posInf := math.Inf(1)
	data, err := encoder.Marshal(posInf)
	if err != nil {
		t.Fatalf("Marshal positive infinity failed: %v", err)
	}
	var posInfResult float64
	err = encoder.Unmarshal(data, &posInfResult)
	if err != nil {
		t.Fatalf("Unmarshal positive infinity failed: %v", err)
	}
	if !math.IsInf(posInfResult, 1) {
		t.Errorf("Positive infinity round trip failed: expected +Inf, got %f", posInfResult)
	}

	// 测试负无穷
	negInf := math.Inf(-1)
	data, err = encoder.Marshal(negInf)
	if err != nil {
		t.Fatalf("Marshal negative infinity failed: %v", err)
	}
	var negInfResult float64
	err = encoder.Unmarshal(data, &negInfResult)
	if err != nil {
		t.Fatalf("Unmarshal negative infinity failed: %v", err)
	}
	if !math.IsInf(negInfResult, -1) {
		t.Errorf("Negative infinity round trip failed: expected -Inf, got %f", negInfResult)
	}

	// 测试 NaN
	nan := math.NaN()
	data, err = encoder.Marshal(nan)
	if err != nil {
		t.Fatalf("Marshal NaN failed: %v", err)
	}
	var nanResult float64
	err = encoder.Unmarshal(data, &nanResult)
	if err != nil {
		t.Fatalf("Unmarshal NaN failed: %v", err)
	}
	if !math.IsNaN(nanResult) {
		t.Errorf("NaN round trip failed: expected NaN, got %f", nanResult)
	}
}

// TestBinaryUnmarshalInsufficientData 测试数据不足时反序列化失败
// Validates: Requirements 7.2
func TestBinaryUnmarshalInsufficientData(t *testing.T) {
	encoder := encodingx.NewBinary()

	// 尝试用 2 字节数据反序列化 int32（需要 4 字节）
	insufficientData := []byte{0x01, 0x02}
	var result int32
	err := encoder.Unmarshal(insufficientData, &result)
	if err == nil {
		t.Error("Unmarshal should fail for insufficient data")
	}
}

// TestLittleEndianUnmarshalInsufficientData 测试 LittleEndian 数据不足时反序列化失败
// Validates: Requirements 7.3
func TestLittleEndianUnmarshalInsufficientData(t *testing.T) {
	encoder := encodingx.NewLittleEndian()

	// 尝试用 4 字节数据反序列化 int64（需要 8 字节）
	insufficientData := []byte{0x01, 0x02, 0x03, 0x04}
	var result int64
	err := encoder.Unmarshal(insufficientData, &result)
	if err == nil {
		t.Error("Unmarshal should fail for insufficient data")
	}
}

// TestBigEndianUnmarshalInsufficientData 测试 BigEndian 数据不足时反序列化失败
// Validates: Requirements 7.4
func TestBigEndianUnmarshalInsufficientData(t *testing.T) {
	encoder := encodingx.NewBigEndian()

	// 尝试用 4 字节数据反序列化 float64（需要 8 字节）
	insufficientData := []byte{0x01, 0x02, 0x03, 0x04}
	var result float64
	err := encoder.Unmarshal(insufficientData, &result)
	if err == nil {
		t.Error("Unmarshal should fail for insufficient data")
	}
}

// ============================================================================
// 其他固定大小类型测试
// ============================================================================

// TestBinaryMarshalUint32 测试 Binary 编码器序列化 uint32
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalUint32(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []uint32{0, 1, 255, 65535, 4294967295}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result uint32
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryMarshalUint64 测试 Binary 编码器序列化 uint64
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalUint64(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []uint64{0, 1, 255, 65535, 4294967295, 18446744073709551615}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result uint64
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryMarshalInt16 测试 Binary 编码器序列化 int16
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalInt16(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []int16{0, 1, -1, 127, -128, 32767, -32768}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int16
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryMarshalUint16 测试 Binary 编码器序列化 uint16
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalUint16(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []uint16{0, 1, 255, 65535}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result uint16
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryMarshalInt8 测试 Binary 编码器序列化 int8
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalInt8(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []int8{0, 1, -1, 127, -128}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result int8
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryMarshalUint8 测试 Binary 编码器序列化 uint8
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalUint8(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []uint8{0, 1, 127, 255}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %d: %v", original, err)
		}

		// 反序列化
		var result uint8
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %d: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %d, got %d", original, result)
		}
	}
}

// TestBinaryMarshalBool 测试 Binary 编码器序列化 bool
// Validates: Requirements 7.1, 7.2
func TestBinaryMarshalBool(t *testing.T) {
	encoder := encodingx.NewBinary()
	testCases := []bool{true, false}

	for _, original := range testCases {
		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for %v: %v", original, err)
		}

		// 反序列化
		var result bool
		err = encoder.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed for %v: %v", original, err)
		}

		// 验证往返一致性
		if result != original {
			t.Errorf("Round trip failed: expected %v, got %v", original, result)
		}
	}
}

// ============================================================================
// Binary 编码器属性测试
// ============================================================================

// genInt32 生成任意 int32 值
func genInt32() *rapid.Generator[int32] {
	return rapid.Int32()
}

// genInt64 生成任意 int64 值
func genInt64() *rapid.Generator[int64] {
	return rapid.Int64()
}

// genUint32 生成任意 uint32 值
func genUint32() *rapid.Generator[uint32] {
	return rapid.Uint32()
}

// genUint64 生成任意 uint64 值
func genUint64() *rapid.Generator[uint64] {
	return rapid.Uint64()
}

// genFloat32 生成任意 float32 值（排除 NaN，因为 NaN != NaN）
func genFloat32() *rapid.Generator[float32] {
	return rapid.Custom(func(t *rapid.T) float32 {
		for {
			bits := rapid.Uint32().Draw(t, "float32bits")
			f := math.Float32frombits(bits)
			// 排除 NaN，因为 NaN != NaN 会导致测试失败
			if !math.IsNaN(float64(f)) {
				return f
			}
		}
	})
}

// genFloat64 生成任意 float64 值（排除 NaN，因为 NaN != NaN）
func genFloat64() *rapid.Generator[float64] {
	return rapid.Custom(func(t *rapid.T) float64 {
		for {
			bits := rapid.Uint64().Draw(t, "float64bits")
			f := math.Float64frombits(bits)
			// 排除 NaN，因为 NaN != NaN 会导致测试失败
			if !math.IsNaN(f) {
				return f
			}
		}
	})
}

// TestProperty10_BinaryRoundTripConsistency 测试 Binary Round-Trip 一致性
// **Property 10: Binary Round-Trip 一致性**
// *For any* 有效的固定大小类型，使用 Binary 编码器序列化后再反序列化，
// 应该产生与原始值等价的数据。
// **Validates: Requirements 7.5**
func TestProperty10_BinaryRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewBinary()

	t.Run("int32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genInt32().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			var result int32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("int64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genInt64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var result int64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("uint32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genUint32().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var result uint32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("uint64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genUint64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var result uint64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("float32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFloat32().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var result float32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %f, got %f", original, result)
			}
		})
	})

	t.Run("float64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFloat64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var result float64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %f, got %f", original, result)
			}
		})
	})
}

// TestProperty11_LittleEndianRoundTripConsistency 测试 LittleEndian Round-Trip 一致性
// **Property 11: LittleEndian Round-Trip 一致性**
// *For any* 有效的固定大小类型，使用 LittleEndian 编码器序列化后再反序列化，
// 应该产生与原始值等价的数据。
// **Validates: Requirements 7.6**
func TestProperty11_LittleEndianRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewLittleEndian()

	t.Run("int32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genInt32().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 验证字节顺序是小端序
			expected := make([]byte, 4)
			binary.LittleEndian.PutUint32(expected, uint32(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			// 反序列化
			var result int32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("int64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genInt64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 8)
			binary.LittleEndian.PutUint64(expected, uint64(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result int64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("uint32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genUint32().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 4)
			binary.LittleEndian.PutUint32(expected, original)
			if !BytesEqual(data, expected) {
				t.Fatalf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result uint32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("uint64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genUint64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 8)
			binary.LittleEndian.PutUint64(expected, original)
			if !BytesEqual(data, expected) {
				t.Fatalf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result uint64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("float32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFloat32().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 4)
			binary.LittleEndian.PutUint32(expected, math.Float32bits(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result float32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %f, got %f", original, result)
			}
		})
	})

	t.Run("float64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFloat64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 8)
			binary.LittleEndian.PutUint64(expected, math.Float64bits(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("LittleEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result float64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %f, got %f", original, result)
			}
		})
	})
}

// TestProperty12_BigEndianRoundTripConsistency 测试 BigEndian Round-Trip 一致性
// **Property 12: BigEndian Round-Trip 一致性**
// *For any* 有效的固定大小类型，使用 BigEndian 编码器序列化后再反序列化，
// 应该产生与原始值等价的数据。
// **Validates: Requirements 7.7**
func TestProperty12_BigEndianRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewBigEndian()

	t.Run("int32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genInt32().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 验证字节顺序是大端序
			expected := make([]byte, 4)
			binary.BigEndian.PutUint32(expected, uint32(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			// 反序列化
			var result int32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("int64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genInt64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 8)
			binary.BigEndian.PutUint64(expected, uint64(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result int64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("uint32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genUint32().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 4)
			binary.BigEndian.PutUint32(expected, original)
			if !BytesEqual(data, expected) {
				t.Fatalf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result uint32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("uint64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genUint64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 8)
			binary.BigEndian.PutUint64(expected, original)
			if !BytesEqual(data, expected) {
				t.Fatalf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result uint64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %d, got %d", original, result)
			}
		})
	})

	t.Run("float32", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFloat32().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 4)
			binary.BigEndian.PutUint32(expected, math.Float32bits(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result float32
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %f, got %f", original, result)
			}
		})
	})

	t.Run("float64", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFloat64().Draw(t, "original")

			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			expected := make([]byte, 8)
			binary.BigEndian.PutUint64(expected, math.Float64bits(original))
			if !BytesEqual(data, expected) {
				t.Fatalf("BigEndian encoding mismatch: expected %v, got %v", expected, data)
			}

			var result float64
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != original {
				t.Fatalf("Round-trip failed: original %f, got %f", original, result)
			}
		})
	})
}
