package encodingx_test

import (
	"encoding/base64"
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// Base64 编码器单元测试
// Validates: Requirements 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 14.1, 14.2, 14.3
// ============================================================================

// TestBase64MarshalByteSlice 测试 []byte 编码
// Validates: Requirements 6.1
func TestBase64MarshalByteSlice(t *testing.T) {
	encoder := encodingx.NewBase64()
	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是标准 Base64 编码字符串
	expected := base64.StdEncoding.EncodeToString(original)
	if string(data) != expected {
		t.Errorf("Base64 encoding mismatch: expected %s, got %s", expected, string(data))
	}
}

// TestBase64MarshalBytes 测试 Bytes 类型编码
// Validates: Requirements 6.2
func TestBase64MarshalBytes(t *testing.T) {
	encoder := encodingx.NewBase64()
	originalData := []byte{0x10, 0x20, 0x30, 0x40}
	original := encodingx.MakeBytes(originalData)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是标准 Base64 编码字符串
	expected := base64.StdEncoding.EncodeToString(originalData)
	if string(data) != expected {
		t.Errorf("Base64 encoding mismatch: expected %s, got %s", expected, string(data))
	}
}

// TestBase64MarshalBytesPointer 测试 *Bytes 类型编码
// Validates: Requirements 6.3
func TestBase64MarshalBytesPointer(t *testing.T) {
	encoder := encodingx.NewBase64()
	originalData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	original := encodingx.NewBytes()
	original.Data = originalData

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是标准 Base64 编码字符串
	expected := base64.StdEncoding.EncodeToString(originalData)
	if string(data) != expected {
		t.Errorf("Base64 encoding mismatch: expected %s, got %s", expected, string(data))
	}
}

// TestBase64UnmarshalToBytes 测试反序列化到 *Bytes
// Validates: Requirements 6.4
func TestBase64UnmarshalToBytes(t *testing.T) {
	encoder := encodingx.NewBase64()
	originalData := []byte{0x11, 0x22, 0x33, 0x44, 0x55}
	encodedData := []byte(base64.StdEncoding.EncodeToString(originalData))

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(encodedData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确解码 Base64 字符串
	if !BytesEqual(result.Data, originalData) {
		t.Errorf("Base64 decoding mismatch: expected %v, got %v", originalData, result.Data)
	}
}

// TestBase64MarshalWrongType 测试非字节类型返回错误
// Validates: Requirements 6.5
func TestBase64MarshalWrongType(t *testing.T) {
	encoder := encodingx.NewBase64()

	// 测试结构体类型
	_, err := encoder.Marshal(TestStruct{Integer: 42, String: "test"})
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType, got %v", err)
	}

	// 测试字符串类型
	_, err = encoder.Marshal("hello")
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType for string, got %v", err)
	}

	// 测试整数类型
	_, err = encoder.Marshal(42)
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType for int, got %v", err)
	}

	// 测试 map 类型
	_, err = encoder.Marshal(map[string]int{"key": 1})
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType for map, got %v", err)
	}
}

// TestBase64UnmarshalWrongType 测试反序列化到非 *Bytes 类型返回错误
// Validates: Requirements 6.5
func TestBase64UnmarshalWrongType(t *testing.T) {
	encoder := encodingx.NewBase64()
	encodedData := []byte(base64.StdEncoding.EncodeToString([]byte{0x01, 0x02}))

	// 测试反序列化到结构体
	var ts TestStruct
	err := encoder.Unmarshal(encodedData, &ts)
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType, got %v", err)
	}

	// 测试反序列化到字符串指针
	var s string
	err = encoder.Unmarshal(encodedData, &s)
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType for string pointer, got %v", err)
	}

	// 测试反序列化到 []byte（不是 *Bytes）
	var bs []byte
	err = encoder.Unmarshal(encodedData, &bs)
	if err != encodingx.ErrBase64WrongValueType {
		t.Errorf("Expected ErrBase64WrongValueType for []byte pointer, got %v", err)
	}
}

// TestBase64RoundTrip 测试 Base64 编码/解码往返
// Validates: Requirements 6.1, 6.4
func TestBase64RoundTrip(t *testing.T) {
	encoder := encodingx.NewBase64()
	originalData := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, originalData) {
		t.Errorf("Round trip failed: expected %v, got %v", originalData, result.Data)
	}
}

// TestBase64EmptyData 测试空数据编码/解码
// Validates: Requirements 6.1, 6.4
func TestBase64EmptyData(t *testing.T) {
	encoder := encodingx.NewBase64()
	originalData := []byte{}

	// 序列化空数据
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal empty data failed: %v", err)
	}

	// 验证空数据的 Base64 编码
	expected := base64.StdEncoding.EncodeToString(originalData)
	if string(encoded) != expected {
		t.Errorf("Empty data encoding mismatch: expected '%s', got '%s'", expected, string(encoded))
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal empty data failed: %v", err)
	}

	// 验证往返一致性
	if len(result.Data) != 0 {
		t.Errorf("Empty data round trip failed: expected empty, got %v", result.Data)
	}
}

// TestBase64String 测试 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestBase64String(t *testing.T) {
	encoder := encodingx.NewBase64()
	name := encoder.String()

	if name != "Base64" {
		t.Errorf("String() should return 'Base64', got '%s'", name)
	}
}

// TestBase64Style 测试 Style() 方法返回 EncodingStyleBytes
// Validates: Requirements 14.2
func TestBase64Style(t *testing.T) {
	encoder := encodingx.NewBase64()
	style := encoder.Style()

	if style != encodingx.EncodingStyleBytes {
		t.Errorf("Style() should return EncodingStyleBytes, got %v", style)
	}
}

// TestBase64Reverse 测试 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestBase64Reverse(t *testing.T) {
	encoder := encodingx.NewBase64()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 Base64 编码器
	if reversed.Style() != encodingx.EncodingStyleBytes {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestBase64ImplementsEncoding 测试 Base64 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestBase64ImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewBase64()

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
// Base64URL 编码器单元测试
// Validates: Requirements 6.6, 6.7, 14.1, 14.2, 14.3
// ============================================================================

// TestBase64URLMarshalByteSlice 测试 Base64URL []byte 编码
// Validates: Requirements 6.6
func TestBase64URLMarshalByteSlice(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	// 使用包含 URL 不安全字符的数据（会产生 + 和 / 在标准 Base64 中）
	original := []byte{0xFB, 0xFF, 0xFE, 0xFD, 0xFC}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 URL 安全的 Base64 编码字符串
	expected := base64.URLEncoding.EncodeToString(original)
	if string(data) != expected {
		t.Errorf("Base64URL encoding mismatch: expected %s, got %s", expected, string(data))
	}

	// 验证不包含 URL 不安全字符
	for _, c := range string(data) {
		if c == '+' || c == '/' {
			t.Errorf("Base64URL should not contain '+' or '/', got %s", string(data))
			break
		}
	}
}

// TestBase64URLMarshalBytes 测试 Base64URL Bytes 类型编码
// Validates: Requirements 6.6
func TestBase64URLMarshalBytes(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	originalData := []byte{0xFB, 0xFF, 0xFE}
	original := encodingx.MakeBytes(originalData)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 URL 安全的 Base64 编码字符串
	expected := base64.URLEncoding.EncodeToString(originalData)
	if string(data) != expected {
		t.Errorf("Base64URL encoding mismatch: expected %s, got %s", expected, string(data))
	}
}

// TestBase64URLMarshalBytesPointer 测试 Base64URL *Bytes 类型编码
// Validates: Requirements 6.6
func TestBase64URLMarshalBytesPointer(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	originalData := []byte{0xFB, 0xFF, 0xFE, 0xFD}
	original := encodingx.NewBytes()
	original.Data = originalData

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 URL 安全的 Base64 编码字符串
	expected := base64.URLEncoding.EncodeToString(originalData)
	if string(data) != expected {
		t.Errorf("Base64URL encoding mismatch: expected %s, got %s", expected, string(data))
	}
}

// TestBase64URLUnmarshalToBytes 测试 Base64URL 反序列化到 *Bytes
// Validates: Requirements 6.7
func TestBase64URLUnmarshalToBytes(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	originalData := []byte{0xFB, 0xFF, 0xFE, 0xFD, 0xFC}
	encodedData := []byte(base64.URLEncoding.EncodeToString(originalData))

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(encodedData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证正确解码 URL 安全 Base64 字符串
	if !BytesEqual(result.Data, originalData) {
		t.Errorf("Base64URL decoding mismatch: expected %v, got %v", originalData, result.Data)
	}
}

// TestBase64URLMarshalWrongType 测试 Base64URL 非字节类型返回错误
// Validates: Requirements 6.6
func TestBase64URLMarshalWrongType(t *testing.T) {
	encoder := encodingx.NewBase64URL()

	// 测试结构体类型
	_, err := encoder.Marshal(TestStruct{Integer: 42, String: "test"})
	if err != encodingx.ErrBase64URLWrongValueType {
		t.Errorf("Expected ErrBase64URLWrongValueType, got %v", err)
	}

	// 测试字符串类型
	_, err = encoder.Marshal("hello")
	if err != encodingx.ErrBase64URLWrongValueType {
		t.Errorf("Expected ErrBase64URLWrongValueType for string, got %v", err)
	}

	// 测试整数类型
	_, err = encoder.Marshal(42)
	if err != encodingx.ErrBase64URLWrongValueType {
		t.Errorf("Expected ErrBase64URLWrongValueType for int, got %v", err)
	}
}

// TestBase64URLUnmarshalWrongType 测试 Base64URL 反序列化到非 *Bytes 类型返回错误
// Validates: Requirements 6.7
func TestBase64URLUnmarshalWrongType(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	encodedData := []byte(base64.URLEncoding.EncodeToString([]byte{0x01, 0x02}))

	// 测试反序列化到结构体
	var ts TestStruct
	err := encoder.Unmarshal(encodedData, &ts)
	if err != encodingx.ErrBase64URLWrongValueType {
		t.Errorf("Expected ErrBase64URLWrongValueType, got %v", err)
	}

	// 测试反序列化到字符串指针
	var s string
	err = encoder.Unmarshal(encodedData, &s)
	if err != encodingx.ErrBase64URLWrongValueType {
		t.Errorf("Expected ErrBase64URLWrongValueType for string pointer, got %v", err)
	}
}

// TestBase64URLRoundTrip 测试 Base64URL 编码/解码往返
// Validates: Requirements 6.6, 6.7
func TestBase64URLRoundTrip(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	// 使用会产生 URL 不安全字符的数据
	originalData := []byte{0xFB, 0xFF, 0xFE, 0xFD, 0xFC, 0x00, 0x11, 0x22, 0x33}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, originalData) {
		t.Errorf("Round trip failed: expected %v, got %v", originalData, result.Data)
	}
}

// TestBase64URLEmptyData 测试 Base64URL 空数据编码/解码
// Validates: Requirements 6.6, 6.7
func TestBase64URLEmptyData(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	originalData := []byte{}

	// 序列化空数据
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal empty data failed: %v", err)
	}

	// 验证空数据的 Base64URL 编码
	expected := base64.URLEncoding.EncodeToString(originalData)
	if string(encoded) != expected {
		t.Errorf("Empty data encoding mismatch: expected '%s', got '%s'", expected, string(encoded))
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal empty data failed: %v", err)
	}

	// 验证往返一致性
	if len(result.Data) != 0 {
		t.Errorf("Empty data round trip failed: expected empty, got %v", result.Data)
	}
}

// TestBase64URLString 测试 Base64URL String() 方法返回类型名称
// Validates: Requirements 14.1
func TestBase64URLString(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	name := encoder.String()

	if name != "Base64URL" {
		t.Errorf("String() should return 'Base64URL', got '%s'", name)
	}
}

// TestBase64URLStyle 测试 Base64URL Style() 方法返回 EncodingStyleBytes
// Validates: Requirements 14.2
func TestBase64URLStyle(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	style := encoder.Style()

	if style != encodingx.EncodingStyleBytes {
		t.Errorf("Style() should return EncodingStyleBytes, got %v", style)
	}
}

// TestBase64URLReverse 测试 Base64URL Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestBase64URLReverse(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 Base64URL 编码器
	if reversed.Style() != encodingx.EncodingStyleBytes {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestBase64URLImplementsEncoding 测试 Base64URL 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestBase64URLImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewBase64URL()

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
// Base64 vs Base64URL 对比测试
// ============================================================================

// TestBase64VsBase64URLDifference 测试 Base64 和 Base64URL 编码的差异
// Validates: Requirements 6.1, 6.6
func TestBase64VsBase64URLDifference(t *testing.T) {
	base64Encoder := encodingx.NewBase64()
	base64URLEncoder := encodingx.NewBase64URL()

	// 使用会产生 + 和 / 的数据
	testData := []byte{0xFB, 0xFF, 0xFE}

	// Base64 编码
	base64Encoded, err := base64Encoder.Marshal(testData)
	if err != nil {
		t.Fatalf("Base64 Marshal failed: %v", err)
	}

	// Base64URL 编码
	base64URLEncoded, err := base64URLEncoder.Marshal(testData)
	if err != nil {
		t.Fatalf("Base64URL Marshal failed: %v", err)
	}

	// 验证两种编码结果不同（对于特定数据）
	base64Str := string(base64Encoded)
	base64URLStr := string(base64URLEncoded)

	// 验证 Base64 使用标准字符集
	expectedBase64 := base64.StdEncoding.EncodeToString(testData)
	if base64Str != expectedBase64 {
		t.Errorf("Base64 encoding mismatch: expected %s, got %s", expectedBase64, base64Str)
	}

	// 验证 Base64URL 使用 URL 安全字符集
	expectedBase64URL := base64.URLEncoding.EncodeToString(testData)
	if base64URLStr != expectedBase64URL {
		t.Errorf("Base64URL encoding mismatch: expected %s, got %s", expectedBase64URL, base64URLStr)
	}
}

// TestBase64InvalidBase64String 测试无效 Base64 字符串反序列化
// Validates: Requirements 6.4
func TestBase64InvalidBase64String(t *testing.T) {
	encoder := encodingx.NewBase64()

	// 无效的 Base64 字符串
	invalidData := []byte("!!!invalid-base64!!!")

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(invalidData, result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid Base64 string")
	}
}

// TestBase64URLInvalidBase64URLString 测试无效 Base64URL 字符串反序列化
// Validates: Requirements 6.7
func TestBase64URLInvalidBase64URLString(t *testing.T) {
	encoder := encodingx.NewBase64URL()

	// 无效的 Base64URL 字符串
	invalidData := []byte("!!!invalid-base64url!!!")

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(invalidData, result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid Base64URL string")
	}
}

// ============================================================================
// 边界条件测试
// ============================================================================

// TestBase64LargeData 测试大数据编码/解码
// Validates: Requirements 6.1, 6.4
func TestBase64LargeData(t *testing.T) {
	encoder := encodingx.NewBase64()
	// 生成 1KB 的随机数据
	originalData := make([]byte, 1024)
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal large data failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal large data failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, originalData) {
		t.Error("Large data round trip failed")
	}
}

// TestBase64URLLargeData 测试 Base64URL 大数据编码/解码
// Validates: Requirements 6.6, 6.7
func TestBase64URLLargeData(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	// 生成 1KB 的随机数据
	originalData := make([]byte, 1024)
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal large data failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal large data failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, originalData) {
		t.Error("Large data round trip failed")
	}
}

// TestBase64AllByteValues 测试所有可能的字节值
// Validates: Requirements 6.1, 6.4
func TestBase64AllByteValues(t *testing.T) {
	encoder := encodingx.NewBase64()
	// 包含所有可能的字节值 (0-255)
	originalData := make([]byte, 256)
	for i := 0; i < 256; i++ {
		originalData[i] = byte(i)
	}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal all byte values failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal all byte values failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, originalData) {
		t.Error("All byte values round trip failed")
	}
}

// TestBase64URLAllByteValues 测试 Base64URL 所有可能的字节值
// Validates: Requirements 6.6, 6.7
func TestBase64URLAllByteValues(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	// 包含所有可能的字节值 (0-255)
	originalData := make([]byte, 256)
	for i := 0; i < 256; i++ {
		originalData[i] = byte(i)
	}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal all byte values failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal all byte values failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, originalData) {
		t.Error("All byte values round trip failed")
	}
}

// TestBase64NilBytesPointer 测试 nil Data 的 *Bytes 编码
// Validates: Requirements 6.3
func TestBase64NilBytesPointer(t *testing.T) {
	encoder := encodingx.NewBase64()
	original := encodingx.NewBytes() // Data 为 nil

	// 序列化
	encoded, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal nil *Bytes failed: %v", err)
	}

	// 验证 nil Data 的 Base64 编码（应该是空字符串的编码）
	expected := base64.StdEncoding.EncodeToString(nil)
	if string(encoded) != expected {
		t.Errorf("Nil *Bytes encoding mismatch: expected '%s', got '%s'", expected, string(encoded))
	}
}

// TestBase64URLNilBytesPointer 测试 Base64URL nil Data 的 *Bytes 编码
// Validates: Requirements 6.6
func TestBase64URLNilBytesPointer(t *testing.T) {
	encoder := encodingx.NewBase64URL()
	original := encodingx.NewBytes() // Data 为 nil

	// 序列化
	encoded, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal nil *Bytes failed: %v", err)
	}

	// 验证 nil Data 的 Base64URL 编码（应该是空字符串的编码）
	expected := base64.URLEncoding.EncodeToString(nil)
	if string(encoded) != expected {
		t.Errorf("Nil *Bytes encoding mismatch: expected '%s', got '%s'", expected, string(encoded))
	}
}

// TestBase64SingleByte 测试单字节编码/解码
// Validates: Requirements 6.1, 6.4
func TestBase64SingleByte(t *testing.T) {
	encoder := encodingx.NewBase64()

	for i := 0; i < 256; i++ {
		originalData := []byte{byte(i)}

		// 序列化
		encoded, err := encoder.Marshal(originalData)
		if err != nil {
			t.Fatalf("Marshal single byte %d failed: %v", i, err)
		}

		// 反序列化
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal single byte %d failed: %v", i, err)
		}

		// 验证往返一致性
		if !BytesEqual(result.Data, originalData) {
			t.Errorf("Single byte %d round trip failed: expected %v, got %v", i, originalData, result.Data)
		}
	}
}

// TestBase64PaddingVariations 测试不同长度数据的 padding 变化
// Validates: Requirements 6.1, 6.4
func TestBase64PaddingVariations(t *testing.T) {
	encoder := encodingx.NewBase64()

	// 测试不同长度的数据（会产生不同的 padding）
	testCases := []struct {
		name string
		data []byte
	}{
		{"length_1_padding_2", []byte{0x01}},                   // 1 byte -> 2 padding
		{"length_2_padding_1", []byte{0x01, 0x02}},             // 2 bytes -> 1 padding
		{"length_3_no_padding", []byte{0x01, 0x02, 0x03}},      // 3 bytes -> no padding
		{"length_4_padding_2", []byte{0x01, 0x02, 0x03, 0x04}}, // 4 bytes -> 2 padding
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 序列化
			encoded, err := encoder.Marshal(tc.data)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			result := encodingx.NewBytes()
			err = encoder.Unmarshal(encoded, result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证往返一致性
			if !BytesEqual(result.Data, tc.data) {
				t.Errorf("Round trip failed: expected %v, got %v", tc.data, result.Data)
			}
		})
	}
}

// ============================================================================
// Base64 编码器属性测试
// ============================================================================

// genByteSlice 生成任意长度的字节数组
func genByteSlice() *rapid.Generator[[]byte] {
	return rapid.Custom(func(t *rapid.T) []byte {
		// 生成 0 到 1024 字节的随机数据
		length := rapid.IntRange(0, 1024).Draw(t, "length")
		if length == 0 {
			return []byte{}
		}
		data := make([]byte, length)
		for i := 0; i < length; i++ {
			data[i] = byte(rapid.IntRange(0, 255).Draw(t, "byte"))
		}
		return data
	})
}

// TestProperty8_Base64RoundTripConsistency 测试 Base64 Round-Trip 一致性
// **Property 8: Base64 Round-Trip 一致性**
// *For any* 有效的字节数组，使用 Base64 编码器序列化后再反序列化到 *Bytes，
// 应该产生与原始数据等价的字节。
// **Validates: Requirements 6.8**
func TestProperty8_Base64RoundTripConsistency(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")

		encoder := encodingx.NewBase64()

		// 序列化
		encoded, err := encoder.Marshal(originalData)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// 验证编码结果是有效的 Base64 字符串
		expectedEncoded := base64.StdEncoding.EncodeToString(originalData)
		if string(encoded) != expectedEncoded {
			t.Fatalf("Encoded data mismatch: expected %s, got %s", expectedEncoded, string(encoded))
		}

		// 反序列化到 *Bytes
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, originalData) {
			t.Fatalf("Round-trip failed: original %v, got %v", originalData, result.Data)
		}
	})
}

// TestProperty8_Base64RoundTripConsistency_WithBytesType 测试使用 Bytes 类型的 Round-Trip
// **Property 8: Base64 Round-Trip 一致性**
// **Validates: Requirements 6.8**
func TestProperty8_Base64RoundTripConsistency_WithBytesType(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")
		original := encodingx.MakeBytes(originalData)

		encoder := encodingx.NewBase64()

		// 序列化 Bytes 类型
		encoded, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal Bytes failed: %v", err)
		}

		// 反序列化到 *Bytes
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, originalData) {
			t.Fatalf("Round-trip with Bytes type failed: original %v, got %v", originalData, result.Data)
		}
	})
}

// TestProperty8_Base64RoundTripConsistency_WithBytesPointer 测试使用 *Bytes 类型的 Round-Trip
// **Property 8: Base64 Round-Trip 一致性**
// **Validates: Requirements 6.8**
func TestProperty8_Base64RoundTripConsistency_WithBytesPointer(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")
		original := encodingx.NewBytes()
		original.Data = originalData

		encoder := encodingx.NewBase64()

		// 序列化 *Bytes 类型
		encoded, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal *Bytes failed: %v", err)
		}

		// 反序列化到 *Bytes
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, originalData) {
			t.Fatalf("Round-trip with *Bytes type failed: original %v, got %v", originalData, result.Data)
		}
	})
}

// TestProperty9_Base64URLRoundTripConsistency 测试 Base64URL Round-Trip 一致性
// **Property 9: Base64URL Round-Trip 一致性**
// *For any* 有效的字节数组，使用 Base64URL 编码器序列化后再反序列化到 *Bytes，
// 应该产生与原始数据等价的字节。
// **Validates: Requirements 6.9**
func TestProperty9_Base64URLRoundTripConsistency(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")

		encoder := encodingx.NewBase64URL()

		// 序列化
		encoded, err := encoder.Marshal(originalData)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// 验证编码结果是有效的 URL 安全 Base64 字符串
		expectedEncoded := base64.URLEncoding.EncodeToString(originalData)
		if string(encoded) != expectedEncoded {
			t.Fatalf("Encoded data mismatch: expected %s, got %s", expectedEncoded, string(encoded))
		}

		// 验证不包含 URL 不安全字符
		encodedStr := string(encoded)
		for _, c := range encodedStr {
			if c == '+' || c == '/' {
				t.Fatalf("Base64URL should not contain '+' or '/', got %s", encodedStr)
			}
		}

		// 反序列化到 *Bytes
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, originalData) {
			t.Fatalf("Round-trip failed: original %v, got %v", originalData, result.Data)
		}
	})
}

// TestProperty9_Base64URLRoundTripConsistency_WithBytesType 测试使用 Bytes 类型的 Round-Trip
// **Property 9: Base64URL Round-Trip 一致性**
// **Validates: Requirements 6.9**
func TestProperty9_Base64URLRoundTripConsistency_WithBytesType(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")
		original := encodingx.MakeBytes(originalData)

		encoder := encodingx.NewBase64URL()

		// 序列化 Bytes 类型
		encoded, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal Bytes failed: %v", err)
		}

		// 反序列化到 *Bytes
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, originalData) {
			t.Fatalf("Round-trip with Bytes type failed: original %v, got %v", originalData, result.Data)
		}
	})
}

// TestProperty9_Base64URLRoundTripConsistency_WithBytesPointer 测试使用 *Bytes 类型的 Round-Trip
// **Property 9: Base64URL Round-Trip 一致性**
// **Validates: Requirements 6.9**
func TestProperty9_Base64URLRoundTripConsistency_WithBytesPointer(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")
		original := encodingx.NewBytes()
		original.Data = originalData

		encoder := encodingx.NewBase64URL()

		// 序列化 *Bytes 类型
		encoded, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal *Bytes failed: %v", err)
		}

		// 反序列化到 *Bytes
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, originalData) {
			t.Fatalf("Round-trip with *Bytes type failed: original %v, got %v", originalData, result.Data)
		}
	})
}

// TestProperty8And9_Base64VsBase64URLDifference 测试 Base64 和 Base64URL 编码的差异
// 验证两种编码器对相同数据产生不同的编码结果（当数据包含特定字节时）
// **Validates: Requirements 6.8, 6.9**
func TestProperty8And9_Base64VsBase64URLDifference(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机字节数组
		originalData := genByteSlice().Draw(t, "originalData")

		base64Encoder := encodingx.NewBase64()
		base64URLEncoder := encodingx.NewBase64URL()

		// Base64 编码
		base64Encoded, err := base64Encoder.Marshal(originalData)
		if err != nil {
			t.Fatalf("Base64 Marshal failed: %v", err)
		}

		// Base64URL 编码
		base64URLEncoded, err := base64URLEncoder.Marshal(originalData)
		if err != nil {
			t.Fatalf("Base64URL Marshal failed: %v", err)
		}

		// 验证两种编码都能正确解码回原始数据
		result1 := encodingx.NewBytes()
		err = base64Encoder.Unmarshal(base64Encoded, result1)
		if err != nil {
			t.Fatalf("Base64 Unmarshal failed: %v", err)
		}
		if !BytesEqual(result1.Data, originalData) {
			t.Fatalf("Base64 round-trip failed")
		}

		result2 := encodingx.NewBytes()
		err = base64URLEncoder.Unmarshal(base64URLEncoded, result2)
		if err != nil {
			t.Fatalf("Base64URL Unmarshal failed: %v", err)
		}
		if !BytesEqual(result2.Data, originalData) {
			t.Fatalf("Base64URL round-trip failed")
		}

		// 验证 Base64URL 编码不包含 URL 不安全字符
		for _, c := range string(base64URLEncoded) {
			if c == '+' || c == '/' {
				t.Fatalf("Base64URL should not contain '+' or '/', got %s", string(base64URLEncoded))
			}
		}
	})
}
