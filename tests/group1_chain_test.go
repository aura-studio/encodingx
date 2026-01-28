package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// ChainEncoding 单元测试
// Validates: Requirements 10.1, 10.2, 10.3, 10.4, 10.5, 14.1, 14.2, 14.3
// ============================================================================

// ============================================================================
// String() 格式化输出测试
// Validates: Requirements 10.1
// ============================================================================

// TestChainEncodingString 测试 ChainEncoding String() 方法返回格式化的编码链描述
// Validates: Requirements 10.1
func TestChainEncodingString(t *testing.T) {
	// 创建一个简单的 ChainEncoding
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	result := chain.String()

	// 验证格式: [encoder1:encoder2] -> [decoder1:decoder2]
	expected := "[JSON:Base64] -> [Base64:JSON]"
	if result != expected {
		t.Errorf("String() should return '%s', got '%s'", expected, result)
	}
}

// TestChainEncodingStringSingleEncoder 测试单个编码器的 String() 输出
// Validates: Requirements 10.1
func TestChainEncodingStringSingleEncoder(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"JSON"},
	)

	result := chain.String()
	expected := "[JSON] -> [JSON]"
	if result != expected {
		t.Errorf("String() should return '%s', got '%s'", expected, result)
	}
}

// TestChainEncodingStringMultipleEncoders 测试多个编码器的 String() 输出
// Validates: Requirements 10.1
func TestChainEncodingStringMultipleEncoders(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64", "Lazy"},
		[]string{"Lazy", "Base64", "JSON"},
	)

	result := chain.String()
	expected := "[JSON:Base64:Lazy] -> [Lazy:Base64:JSON]"
	if result != expected {
		t.Errorf("String() should return '%s', got '%s'", expected, result)
	}
}

// TestChainEncodingStringEmpty 测试空编码器链的 String() 输出
// Validates: Requirements 10.1
func TestChainEncodingStringEmpty(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{},
		[]string{},
	)

	result := chain.String()
	expected := "[] -> []"
	if result != expected {
		t.Errorf("String() should return '%s', got '%s'", expected, result)
	}
}

// ============================================================================
// Reverse() 编码器/解码器顺序颠倒测试
// Validates: Requirements 10.2
// ============================================================================

// TestChainEncodingReverse 测试 Reverse() 方法返回编码器和解码器顺序颠倒的新 ChainEncoding
// Validates: Requirements 10.2
func TestChainEncodingReverse(t *testing.T) {
	original := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	reversed := original.Reverse()

	// 验证 reversed 的 String() 输出
	// 原始: encoder=[JSON, Base64], decoder=[Base64, JSON]
	// 颠倒后: encoder=[JSON, Base64] (原decoder逆序), decoder=[Base64, JSON] (原encoder逆序)
	expectedStr := "[JSON:Base64] -> [Base64:JSON]"
	if reversed.String() != expectedStr {
		t.Errorf("Reversed chain String() should return '%s', got '%s'", expectedStr, reversed.String())
	}
}

// TestChainEncodingReverseThreeEncoders 测试三个编码器的 Reverse()
// Validates: Requirements 10.2
func TestChainEncodingReverseThreeEncoders(t *testing.T) {
	original := encodingx.NewChainEncoding(
		[]string{"A", "B", "C"},
		[]string{"X", "Y", "Z"},
	)

	reversed := original.Reverse()

	// 原始: encoder=[A, B, C], decoder=[X, Y, Z]
	// 颠倒后: encoder=[Z, Y, X] (原decoder逆序), decoder=[C, B, A] (原encoder逆序)
	expectedStr := "[Z:Y:X] -> [C:B:A]"
	if reversed.String() != expectedStr {
		t.Errorf("Reversed chain String() should return '%s', got '%s'", expectedStr, reversed.String())
	}
}

// TestChainEncodingReverseSymmetry 测试 Reverse() 两次应该产生与原始等价的配置
// Validates: Requirements 10.2
func TestChainEncodingReverseSymmetry(t *testing.T) {
	original := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64", "Lazy"},
		[]string{"Lazy", "Base64", "JSON"},
	)

	// 两次 Reverse 应该回到原始状态
	doubleReversed := original.Reverse().Reverse()

	if original.String() != doubleReversed.String() {
		t.Errorf("Double reverse should return original: expected '%s', got '%s'",
			original.String(), doubleReversed.String())
	}
}

// TestChainEncodingReverseSingleEncoder 测试单个编码器的 Reverse()
// Validates: Requirements 10.2
func TestChainEncodingReverseSingleEncoder(t *testing.T) {
	original := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"JSON"},
	)

	reversed := original.Reverse()

	// 单个编码器的 Reverse 应该保持不变
	expectedStr := "[JSON] -> [JSON]"
	if reversed.String() != expectedStr {
		t.Errorf("Reversed single encoder chain String() should return '%s', got '%s'",
			expectedStr, reversed.String())
	}
}

// ============================================================================
// 链式序列化测试 (Marshal)
// Validates: Requirements 10.3
// ============================================================================

// TestChainEncodingMarshalJSONBase64 测试 JSON -> Base64 链式序列化
// Validates: Requirements 10.3
func TestChainEncodingMarshalJSONBase64(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "hello",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的数据不为空
	if len(data) == 0 {
		t.Error("Marshal should return non-empty data")
	}
}

// TestChainEncodingMarshalJSONBase64Lazy 测试 JSON -> Base64 -> Lazy 链式序列化
// Validates: Requirements 10.3
func TestChainEncodingMarshalJSONBase64Lazy(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64", "Lazy"},
		[]string{"Lazy", "Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 100,
		String:  "test",
		Bool:    false,
		Float:   2.718,
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的数据不为空
	if len(data) == 0 {
		t.Error("Marshal should return non-empty data")
	}
}

// TestChainEncodingMarshalSingleEncoder 测试单个编码器的链式序列化
// Validates: Requirements 10.3
func TestChainEncodingMarshalSingleEncoder(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"JSON"},
	)

	original := TestStruct{
		Integer: 1,
		String:  "single",
		Bool:    true,
		Float:   1.0,
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的数据不为空
	if len(data) == 0 {
		t.Error("Marshal should return non-empty data")
	}
}

// TestChainEncodingMarshalLazy 测试 Lazy 编码器链式序列化
// Validates: Requirements 10.3
func TestChainEncodingMarshalLazy(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"Lazy"},
		[]string{"Lazy"},
	)

	original := []byte{0x01, 0x02, 0x03, 0x04}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的数据与原始数据相等
	if !BytesEqual(data, original) {
		t.Errorf("Lazy chain Marshal should return same data: expected %v, got %v", original, data)
	}
}

// ============================================================================
// 链式反序列化测试 (Unmarshal)
// Validates: Requirements 10.4
// ============================================================================

// TestChainEncodingUnmarshalJSONBase64 测试 Base64 -> JSON 链式反序列化
// Validates: Requirements 10.4
func TestChainEncodingUnmarshalJSONBase64(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "hello",
		Bool:    true,
		Float:   3.14,
	}

	// 先序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 再反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingUnmarshalJSONBase64Lazy 测试 Lazy -> Base64 -> JSON 链式反序列化
// Validates: Requirements 10.4
func TestChainEncodingUnmarshalJSONBase64Lazy(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64", "Lazy"},
		[]string{"Lazy", "Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 100,
		String:  "test",
		Bool:    false,
		Float:   2.718,
	}

	// 先序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 再反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingUnmarshalSingleEncoder 测试单个编码器的链式反序列化
// Validates: Requirements 10.4
func TestChainEncodingUnmarshalSingleEncoder(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"JSON"},
	)

	original := TestStruct{
		Integer: 1,
		String:  "single",
		Bool:    true,
		Float:   1.0,
	}

	// 先序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 再反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingUnmarshalLazy 测试 Lazy 编码器链式反序列化
// Validates: Requirements 10.4
func TestChainEncodingUnmarshalLazy(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"Lazy"},
		[]string{"Lazy"},
	)

	original := []byte{0x01, 0x02, 0x03, 0x04}

	// 先序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 再反序列化
	result := encodingx.NewBytes()
	err = chain.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, original) {
		t.Errorf("Round trip failed: expected %v, got %v", original, result.Data)
	}
}

// ============================================================================
// 中间编码器 Style 错误测试
// Validates: Requirements 10.5
// ============================================================================

// TestChainEncodingMarshalMiddleEncoderStyleError 测试中间编码器 Style 为 EncodingStyleStruct 时返回错误
// Validates: Requirements 10.5
func TestChainEncodingMarshalMiddleEncoderStyleError(t *testing.T) {
	// JSON 的 Style 是 EncodingStyleStruct，作为中间编码器应该返回错误
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "JSON"}, // 第二个 JSON 作为中间编码器
		[]string{"JSON", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "test",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化应该返回 ErrWrongEncodingStyle 错误
	_, err := chain.Marshal(original)
	if err != encodingx.ErrWrongEncodingStyle {
		t.Errorf("Expected ErrWrongEncodingStyle, got %v", err)
	}
}

// TestChainEncodingUnmarshalMiddleDecoderStyleError 测试中间解码器 Style 为 EncodingStyleStruct 时返回错误
// Validates: Requirements 10.5
func TestChainEncodingUnmarshalMiddleDecoderStyleError(t *testing.T) {
	// 创建一个有效的编码链来生成测试数据
	validChain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "test",
		Bool:    true,
		Float:   3.14,
	}

	// 使用有效链序列化
	data, err := validChain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 创建一个中间解码器 Style 为 EncodingStyleStruct 的链
	invalidChain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"JSON", "JSON"}, // 第一个 JSON 作为中间解码器
	)

	// 反序列化应该返回 ErrWrongEncodingStyle 错误
	var result TestStruct
	err = invalidChain.Unmarshal(data, &result)
	if err != encodingx.ErrWrongEncodingStyle {
		t.Errorf("Expected ErrWrongEncodingStyle, got %v", err)
	}
}

// TestChainEncodingMarshalYAMLMiddleError 测试 YAML 作为中间编码器返回错误
// Validates: Requirements 10.5
func TestChainEncodingMarshalYAMLMiddleError(t *testing.T) {
	// YAML 的 Style 是 EncodingStyleStruct，作为中间编码器应该返回错误
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "YAML"}, // YAML 作为中间编码器
		[]string{"YAML", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "test",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化应该返回 ErrWrongEncodingStyle 错误
	_, err := chain.Marshal(original)
	if err != encodingx.ErrWrongEncodingStyle {
		t.Errorf("Expected ErrWrongEncodingStyle, got %v", err)
	}
}

// TestChainEncodingMarshalXMLMiddleError 测试 XML 作为中间编码器返回错误
// Validates: Requirements 10.5
func TestChainEncodingMarshalXMLMiddleError(t *testing.T) {
	// XML 的 Style 是 EncodingStyleStruct，作为中间编码器应该返回错误
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "XML"}, // XML 作为中间编码器
		[]string{"XML", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "test",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化应该返回 ErrWrongEncodingStyle 错误
	_, err := chain.Marshal(original)
	if err != encodingx.ErrWrongEncodingStyle {
		t.Errorf("Expected ErrWrongEncodingStyle, got %v", err)
	}
}

// ============================================================================
// Style() 返回 EncodingStyleMix 测试
// Validates: Requirements 14.2
// ============================================================================

// TestChainEncodingStyle 测试 ChainEncoding Style() 方法返回 EncodingStyleMix
// Validates: Requirements 14.2
func TestChainEncodingStyle(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	style := chain.Style()

	if style != encodingx.EncodingStyleMix {
		t.Errorf("Style() should return EncodingStyleMix, got %v", style)
	}
}

// TestChainEncodingStyleSingleEncoder 测试单个编码器的 ChainEncoding Style()
// Validates: Requirements 14.2
func TestChainEncodingStyleSingleEncoder(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"JSON"},
	)

	style := chain.Style()

	if style != encodingx.EncodingStyleMix {
		t.Errorf("Style() should return EncodingStyleMix, got %v", style)
	}
}

// TestChainEncodingStyleEmpty 测试空 ChainEncoding Style()
// Validates: Requirements 14.2
func TestChainEncodingStyleEmpty(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{},
		[]string{},
	)

	style := chain.Style()

	if style != encodingx.EncodingStyleMix {
		t.Errorf("Style() should return EncodingStyleMix, got %v", style)
	}
}

// ============================================================================
// Encoding 接口实现测试
// Validates: Requirements 14.1, 14.2, 14.3
// ============================================================================

// TestChainEncodingImplementsEncoding 测试 ChainEncoding 实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestChainEncodingImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	// 验证 String() 返回非空字符串
	if encoder.String() == "" {
		t.Error("String() should return non-empty string")
	}

	// 验证 Style() 返回有效的 EncodingStyleType
	style := encoder.Style()
	if style != encodingx.EncodingStyleStruct &&
		style != encodingx.EncodingStyleBytes &&
		style != encodingx.EncodingStyleMix {
		t.Errorf("Style() returned invalid EncodingStyleType: %v", style)
	}

	// 验证 Reverse() 返回非 nil 的 Encoding
	reversed := encoder.Reverse()
	if reversed == nil {
		t.Error("Reverse() should return non-nil Encoding")
	}
}

// ============================================================================
// Empty() 函数测试
// Validates: Requirements 12.3
// ============================================================================

// TestEmpty 测试 Empty() 返回默认的空 ChainEncoding
// Validates: Requirements 12.3
func TestEmpty(t *testing.T) {
	empty := encodingx.Empty()

	// 验证返回的是 ChainEncoding
	if empty == nil {
		t.Fatal("Empty() should return non-nil Encoding")
	}

	// 验证 String() 输出
	expectedStr := "[Lazy] -> [Lazy]"
	if empty.String() != expectedStr {
		t.Errorf("Empty() String() should return '%s', got '%s'", expectedStr, empty.String())
	}

	// 验证 Style() 返回 EncodingStyleMix
	if empty.Style() != encodingx.EncodingStyleMix {
		t.Errorf("Empty() Style() should return EncodingStyleMix, got %v", empty.Style())
	}
}

// TestEmptyRoundTrip 测试 Empty() 返回的 ChainEncoding 的往返一致性
// Validates: Requirements 12.3
func TestEmptyRoundTrip(t *testing.T) {
	empty := encodingx.Empty()

	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// 序列化
	data, err := empty.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = empty.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, original) {
		t.Errorf("Round trip failed: expected %v, got %v", original, result.Data)
	}
}

// ============================================================================
// 边界条件和错误处理测试
// ============================================================================

// TestChainEncodingMarshalUnknownEncoder 测试使用未注册的编码器名称
// Validates: Requirements 10.3
func TestChainEncodingMarshalUnknownEncoder(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"UnknownEncoder"},
		[]string{"UnknownEncoder"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "test",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化应该返回错误
	_, err := chain.Marshal(original)
	if err == nil {
		t.Error("Marshal with unknown encoder should return error")
	}
	if err != encodingx.ErrEncodingMissingEncoding {
		t.Errorf("Expected ErrEncodingMissingEncoding, got %v", err)
	}
}

// TestChainEncodingUnmarshalUnknownDecoder 测试使用未注册的解码器名称
// Validates: Requirements 10.4
func TestChainEncodingUnmarshalUnknownDecoder(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"UnknownDecoder"},
	)

	// 先用 JSON 序列化
	jsonChain := encodingx.NewChainEncoding(
		[]string{"JSON"},
		[]string{"JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "test",
		Bool:    true,
		Float:   3.14,
	}

	data, err := jsonChain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化应该返回错误
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err == nil {
		t.Error("Unmarshal with unknown decoder should return error")
	}
	if err != encodingx.ErrEncodingMissingEncoding {
		t.Errorf("Expected ErrEncodingMissingEncoding, got %v", err)
	}
}

// ============================================================================
// 复杂链式编码测试
// ============================================================================

// TestChainEncodingComplexChain 测试复杂的编码链
// Validates: Requirements 10.3, 10.4
func TestChainEncodingComplexChain(t *testing.T) {
	// JSON -> Base64 -> Lazy 链
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64", "Lazy"},
		[]string{"Lazy", "Base64", "JSON"},
	)

	original := NestedStruct{
		Name: "complex",
		Inner: TestStruct{
			Integer: 123,
			String:  "nested",
			Bool:    true,
			Float:   9.99,
		},
		Slice: []int{1, 2, 3, 4, 5},
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result NestedStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingReverseRoundTrip 测试使用 Reverse() 进行往返
// Validates: Requirements 10.2, 10.3, 10.4
func TestChainEncodingReverseRoundTrip(t *testing.T) {
	// 创建编码链
	encodeChain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	// 获取反向链
	decodeChain := encodeChain.Reverse()

	original := TestStruct{
		Integer: 42,
		String:  "reverse test",
		Bool:    true,
		Float:   3.14159,
	}

	// 使用编码链序列化
	data, err := encodeChain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 使用反向链反序列化
	var result TestStruct
	err = decodeChain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal with reversed chain failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Reverse round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingBase64LazyChain 测试 Base64 -> Lazy 链
// Validates: Requirements 10.3, 10.4
func TestChainEncodingBase64LazyChain(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"Base64", "Lazy"},
		[]string{"Lazy", "Base64"},
	)

	original := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = chain.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !BytesEqual(result.Data, original) {
		t.Errorf("Round trip failed: expected %v, got %v", original, result.Data)
	}
}

// TestChainEncodingMultipleRoundTrips 测试多次往返
// Validates: Requirements 10.3, 10.4
func TestChainEncodingMultipleRoundTrips(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "multiple",
		Bool:    true,
		Float:   3.14,
	}

	// 进行多次往返
	for i := 0; i < 5; i++ {
		// 序列化
		data, err := chain.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed on iteration %d: %v", i, err)
		}

		// 反序列化
		var result TestStruct
		err = chain.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed on iteration %d: %v", i, err)
		}

		// 验证往返一致性
		if !original.Equal(result) {
			t.Errorf("Round trip failed on iteration %d: expected %+v, got %+v", i, original, result)
		}
	}
}

// ============================================================================
// 特殊数据类型测试
// ============================================================================

// TestChainEncodingEmptyStruct 测试空结构体的链式编码
// Validates: Requirements 10.3, 10.4
func TestChainEncodingEmptyStruct(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{} // 空结构体

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingSpecialCharacters 测试包含特殊字符的数据
// Validates: Requirements 10.3, 10.4
func TestChainEncodingSpecialCharacters(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{
		Integer: -999,
		String:  "special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		Bool:    false,
		Float:   -123.456,
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingUnicodeString 测试 Unicode 字符串
// Validates: Requirements 10.3, 10.4
func TestChainEncodingUnicodeString(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	original := TestStruct{
		Integer: 42,
		String:  "Unicode: 你好世界 🌍 مرحبا العالم",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed: expected %+v, got %+v", original, result)
	}
}

// TestChainEncodingLargeData 测试大数据的链式编码
// Validates: Requirements 10.3, 10.4
func TestChainEncodingLargeData(t *testing.T) {
	chain := encodingx.NewChainEncoding(
		[]string{"JSON", "Base64"},
		[]string{"Base64", "JSON"},
	)

	// 生成大字符串
	largeString := ""
	for i := 0; i < 1000; i++ {
		largeString += "abcdefghij"
	}

	original := TestStruct{
		Integer: 999999,
		String:  largeString,
		Bool:    true,
		Float:   999999.999999,
	}

	// 序列化
	data, err := chain.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	if !original.Equal(result) {
		t.Errorf("Round trip failed for large data")
	}
}

// ============================================================================
// ChainEncoding 属性测试
// ============================================================================

// genValidChainConfig 生成有效的编码链配置
// 编码链的规则：
// - 第一个编码器可以是任意类型（EncodingStyleStruct 或 EncodingStyleBytes）
// - 中间编码器必须是 EncodingStyleBytes 类型（如 Base64, Lazy）
// - 最后一个解码器可以是任意类型
func genValidChainConfig() *rapid.Generator[struct {
	encoder []string
	decoder []string
}] {
	return rapid.Custom(func(t *rapid.T) struct {
		encoder []string
		decoder []string
	} {
		// 可用的 EncodingStyleStruct 编码器（只能作为第一个编码器）
		structEncoders := []string{"JSON", "YAML", "XML"}
		// 可用的 EncodingStyleBytes 编码器（可以作为中间编码器）
		bytesEncoders := []string{"Base64", "Base64URL", "Lazy"}

		// 生成编码链长度 (1-3)
		chainLength := rapid.IntRange(1, 3).Draw(t, "chainLength")

		encoder := make([]string, chainLength)
		decoder := make([]string, chainLength)

		// 第一个编码器可以是任意类型
		allEncoders := append(structEncoders, bytesEncoders...)
		firstEncoderIdx := rapid.IntRange(0, len(allEncoders)-1).Draw(t, "firstEncoder")
		encoder[0] = allEncoders[firstEncoderIdx]

		// 中间编码器必须是 EncodingStyleBytes 类型
		for i := 1; i < chainLength; i++ {
			idx := rapid.IntRange(0, len(bytesEncoders)-1).Draw(t, "middleEncoder")
			encoder[i] = bytesEncoders[idx]
		}

		// 解码器顺序与编码器相反
		for i := 0; i < chainLength; i++ {
			decoder[i] = encoder[chainLength-1-i]
		}

		return struct {
			encoder []string
			decoder []string
		}{encoder: encoder, decoder: decoder}
	})
}

// genValidBytesChainConfig 生成只包含 EncodingStyleBytes 编码器的链配置
// 用于测试字节数据的往返
func genValidBytesChainConfig() *rapid.Generator[struct {
	encoder []string
	decoder []string
}] {
	return rapid.Custom(func(t *rapid.T) struct {
		encoder []string
		decoder []string
	} {
		// 只使用 EncodingStyleBytes 编码器
		bytesEncoders := []string{"Base64", "Base64URL", "Lazy"}

		// 生成编码链长度 (1-3)
		chainLength := rapid.IntRange(1, 3).Draw(t, "chainLength")

		encoder := make([]string, chainLength)
		decoder := make([]string, chainLength)

		for i := 0; i < chainLength; i++ {
			idx := rapid.IntRange(0, len(bytesEncoders)-1).Draw(t, "encoder")
			encoder[i] = bytesEncoders[idx]
		}

		// 解码器顺序与编码器相反
		for i := 0; i < chainLength; i++ {
			decoder[i] = encoder[chainLength-1-i]
		}

		return struct {
			encoder []string
			decoder []string
		}{encoder: encoder, decoder: decoder}
	})
}

// genTestStruct 生成随机的 TestStruct
func genTestStruct() *rapid.Generator[TestStruct] {
	return rapid.Custom(func(t *rapid.T) TestStruct {
		return TestStruct{
			Integer: rapid.IntRange(-10000, 10000).Draw(t, "integer"),
			String:  rapid.StringMatching(`[a-zA-Z0-9 ]{1,50}`).Draw(t, "string"),
			Bool:    rapid.Bool().Draw(t, "bool"),
			Float:   rapid.Float64Range(-1000.0, 1000.0).Draw(t, "float"),
		}
	})
}

// genNestedStruct 生成随机的 NestedStruct
func genNestedStruct() *rapid.Generator[NestedStruct] {
	return rapid.Custom(func(t *rapid.T) NestedStruct {
		sliceLen := rapid.IntRange(1, 10).Draw(t, "sliceLen")
		slice := make([]int, sliceLen)
		for i := 0; i < sliceLen; i++ {
			slice[i] = rapid.IntRange(0, 1000).Draw(t, "sliceItem")
		}
		return NestedStruct{
			Name:  rapid.StringMatching(`[a-zA-Z0-9]{1,30}`).Draw(t, "name"),
			Inner: genTestStruct().Draw(t, "inner"),
			Slice: slice,
		}
	})
}

// genByteSliceForChain 生成用于链式编码的字节数组
func genByteSliceForChain() *rapid.Generator[[]byte] {
	return rapid.Custom(func(t *rapid.T) []byte {
		length := rapid.IntRange(1, 256).Draw(t, "length")
		data := make([]byte, length)
		for i := 0; i < length; i++ {
			data[i] = byte(rapid.IntRange(0, 255).Draw(t, "byte"))
		}
		return data
	})
}

// TestProperty16_ChainEncodingReverseRoundTripConsistency 测试 ChainEncoding Reverse Round-Trip 一致性
// **Property 16: ChainEncoding Reverse Round-Trip 一致性**
// *For any* 有效的编码链配置和数据，使用 ChainEncoding 序列化后用其 Reverse() 反序列化，
// 应该产生与原始数据等价的对象。
// **Validates: Requirements 10.6**
func TestProperty16_ChainEncodingReverseRoundTripConsistency(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成有效的编码链配置
		config := genValidChainConfig().Draw(t, "config")

		// 创建 ChainEncoding
		chain := encodingx.NewChainEncoding(config.encoder, config.decoder)

		// 获取反向链
		reversedChain := chain.Reverse()

		// 根据第一个编码器类型选择测试数据
		firstEncoder := config.encoder[0]

		switch firstEncoder {
		case "JSON", "YAML":
			// 使用 TestStruct 进行测试
			original := genTestStruct().Draw(t, "original")

			// 使用原始链序列化
			data, err := chain.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 使用反向链反序列化
			var result TestStruct
			err = reversedChain.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal with reversed chain failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Fatalf("Reverse round-trip failed: original %+v, got %+v", original, result)
			}

		case "XML":
			// XML 需要使用 XMLTestStruct
			original := XMLTestStruct{
				Integer: rapid.IntRange(-10000, 10000).Draw(t, "integer"),
				String:  rapid.StringMatching(`[a-zA-Z0-9]{1,50}`).Draw(t, "string"),
				Bool:    rapid.Bool().Draw(t, "bool"),
				Float:   rapid.Float64Range(-1000.0, 1000.0).Draw(t, "float"),
			}

			// 使用原始链序列化
			data, err := chain.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 使用反向链反序列化
			var result XMLTestStruct
			err = reversedChain.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal with reversed chain failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Fatalf("Reverse round-trip failed: original %+v, got %+v", original, result)
			}

		default:
			// 对于 EncodingStyleBytes 编码器，使用字节数组
			original := genByteSliceForChain().Draw(t, "original")

			// 使用原始链序列化
			data, err := chain.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 使用反向链反序列化
			result := encodingx.NewBytes()
			err = reversedChain.Unmarshal(data, result)
			if err != nil {
				t.Fatalf("Unmarshal with reversed chain failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !BytesEqual(result.Data, original) {
				t.Fatalf("Reverse round-trip failed: original %v, got %v", original, result.Data)
			}
		}
	})
}

// TestProperty16_ChainEncodingReverseRoundTripConsistency_WithNestedStruct 测试嵌套结构体的 Reverse Round-Trip
// **Property 16: ChainEncoding Reverse Round-Trip 一致性**
// **Validates: Requirements 10.6**
func TestProperty16_ChainEncodingReverseRoundTripConsistency_WithNestedStruct(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 使用 JSON 编码器链（支持嵌套结构体）
		bytesEncoders := []string{"Base64", "Lazy"}
		chainLength := rapid.IntRange(1, 2).Draw(t, "chainLength")

		encoder := make([]string, chainLength+1)
		decoder := make([]string, chainLength+1)

		// 第一个编码器是 JSON
		encoder[0] = "JSON"
		for i := 1; i <= chainLength; i++ {
			idx := rapid.IntRange(0, len(bytesEncoders)-1).Draw(t, "encoder")
			encoder[i] = bytesEncoders[idx]
		}

		// 解码器顺序与编码器相反
		for i := 0; i <= chainLength; i++ {
			decoder[i] = encoder[chainLength-i]
		}

		// 创建 ChainEncoding
		chain := encodingx.NewChainEncoding(encoder, decoder)
		reversedChain := chain.Reverse()

		// 生成嵌套结构体
		original := genNestedStruct().Draw(t, "original")

		// 使用原始链序列化
		data, err := chain.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// 使用反向链反序列化
		var result NestedStruct
		err = reversedChain.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal with reversed chain failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !original.Equal(result) {
			t.Fatalf("Reverse round-trip failed for nested struct: original %+v, got %+v", original, result)
		}
	})
}

// TestProperty16_ChainEncodingReverseRoundTripConsistency_BytesOnly 测试纯字节编码链的 Reverse Round-Trip
// **Property 16: ChainEncoding Reverse Round-Trip 一致性**
// **Validates: Requirements 10.6**
func TestProperty16_ChainEncodingReverseRoundTripConsistency_BytesOnly(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成只包含 EncodingStyleBytes 编码器的链配置
		config := genValidBytesChainConfig().Draw(t, "config")

		// 创建 ChainEncoding
		chain := encodingx.NewChainEncoding(config.encoder, config.decoder)
		reversedChain := chain.Reverse()

		// 生成字节数组
		original := genByteSliceForChain().Draw(t, "original")

		// 使用原始链序列化
		data, err := chain.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// 使用反向链反序列化
		result := encodingx.NewBytes()
		err = reversedChain.Unmarshal(data, result)
		if err != nil {
			t.Fatalf("Unmarshal with reversed chain failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		if !BytesEqual(result.Data, original) {
			t.Fatalf("Reverse round-trip failed: original %v, got %v", original, result.Data)
		}
	})
}

// TestProperty17_ChainEncodingReverseSymmetry 测试 ChainEncoding Reverse 对称性
// **Property 17: ChainEncoding Reverse 对称性**
// *For any* ChainEncoding，调用 Reverse() 两次应该产生与原始编码链等价的配置。
// **Validates: Requirements 10.2**
func TestProperty17_ChainEncodingReverseSymmetry(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机的编码器和解码器名称列表
		encoderNames := []string{"JSON", "YAML", "XML", "Base64", "Base64URL", "Lazy"}

		// 生成编码链长度 (1-4)
		chainLength := rapid.IntRange(1, 4).Draw(t, "chainLength")

		encoder := make([]string, chainLength)
		decoder := make([]string, chainLength)

		for i := 0; i < chainLength; i++ {
			encoderIdx := rapid.IntRange(0, len(encoderNames)-1).Draw(t, "encoderIdx")
			decoderIdx := rapid.IntRange(0, len(encoderNames)-1).Draw(t, "decoderIdx")
			encoder[i] = encoderNames[encoderIdx]
			decoder[i] = encoderNames[decoderIdx]
		}

		// 创建原始 ChainEncoding
		original := encodingx.NewChainEncoding(encoder, decoder)

		// 调用 Reverse() 两次
		reversed := original.Reverse()
		doubleReversed := reversed.Reverse()

		// 验证 String() 输出相等（表示配置等价）
		if original.String() != doubleReversed.String() {
			t.Fatalf("Double reverse should return original config: expected '%s', got '%s'",
				original.String(), doubleReversed.String())
		}
	})
}

// TestProperty17_ChainEncodingReverseSymmetry_SingleEncoder 测试单个编码器的 Reverse 对称性
// **Property 17: ChainEncoding Reverse 对称性**
// **Validates: Requirements 10.2**
func TestProperty17_ChainEncodingReverseSymmetry_SingleEncoder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		encoderNames := []string{"JSON", "YAML", "XML", "Base64", "Base64URL", "Lazy"}

		// 随机选择一个编码器
		encoderIdx := rapid.IntRange(0, len(encoderNames)-1).Draw(t, "encoderIdx")
		decoderIdx := rapid.IntRange(0, len(encoderNames)-1).Draw(t, "decoderIdx")

		encoder := []string{encoderNames[encoderIdx]}
		decoder := []string{encoderNames[decoderIdx]}

		// 创建原始 ChainEncoding
		original := encodingx.NewChainEncoding(encoder, decoder)

		// 调用 Reverse() 两次
		reversed := original.Reverse()
		doubleReversed := reversed.Reverse()

		// 验证 String() 输出相等
		if original.String() != doubleReversed.String() {
			t.Fatalf("Double reverse should return original config: expected '%s', got '%s'",
				original.String(), doubleReversed.String())
		}
	})
}

// TestProperty17_ChainEncodingReverseSymmetry_EmptyChain 测试空编码链的 Reverse 对称性
// **Property 17: ChainEncoding Reverse 对称性**
// **Validates: Requirements 10.2**
func TestProperty17_ChainEncodingReverseSymmetry_EmptyChain(t *testing.T) {
	// 空编码链的 Reverse 对称性
	original := encodingx.NewChainEncoding([]string{}, []string{})

	reversed := original.Reverse()
	doubleReversed := reversed.Reverse()

	if original.String() != doubleReversed.String() {
		t.Fatalf("Double reverse of empty chain should return original: expected '%s', got '%s'",
			original.String(), doubleReversed.String())
	}
}

// TestProperty17_ChainEncodingReverseSymmetry_Style 测试 Reverse 后 Style 保持不变
// **Property 17: ChainEncoding Reverse 对称性**
// **Validates: Requirements 10.2**
func TestProperty17_ChainEncodingReverseSymmetry_Style(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		encoderNames := []string{"JSON", "YAML", "XML", "Base64", "Base64URL", "Lazy"}

		// 生成编码链长度 (1-3)
		chainLength := rapid.IntRange(1, 3).Draw(t, "chainLength")

		encoder := make([]string, chainLength)
		decoder := make([]string, chainLength)

		for i := 0; i < chainLength; i++ {
			encoderIdx := rapid.IntRange(0, len(encoderNames)-1).Draw(t, "encoderIdx")
			decoderIdx := rapid.IntRange(0, len(encoderNames)-1).Draw(t, "decoderIdx")
			encoder[i] = encoderNames[encoderIdx]
			decoder[i] = encoderNames[decoderIdx]
		}

		// 创建原始 ChainEncoding
		original := encodingx.NewChainEncoding(encoder, decoder)

		// 调用 Reverse()
		reversed := original.Reverse()
		doubleReversed := reversed.Reverse()

		// 验证 Style() 保持不变
		if original.Style() != reversed.Style() {
			t.Fatalf("Reversed chain should have same Style: expected %v, got %v",
				original.Style(), reversed.Style())
		}

		if original.Style() != doubleReversed.Style() {
			t.Fatalf("Double reversed chain should have same Style: expected %v, got %v",
				original.Style(), doubleReversed.Style())
		}
	})
}
