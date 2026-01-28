package encodingx_test

import (
	"encoding/json"
	"testing"

	"github.com/aura-studio/encodingx"
)

// ============================================================================
// JSON 编码器单元测试
// Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 14.1, 14.2, 14.3
// ============================================================================

// TestJSONMarshalStruct 测试普通结构体序列化
// Validates: Requirements 1.1
func TestJSONMarshalStruct(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := TestStruct{
		Integer: 42,
		String:  "hello world",
		Bool:    true,
		Float:   3.14159,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是有效的 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证字段值
	if int(result["integer"].(float64)) != original.Integer {
		t.Errorf("Integer mismatch: expected %d, got %v", original.Integer, result["integer"])
	}
	if result["string"].(string) != original.String {
		t.Errorf("String mismatch: expected %s, got %v", original.String, result["string"])
	}
	if result["bool"].(bool) != original.Bool {
		t.Errorf("Bool mismatch: expected %v, got %v", original.Bool, result["bool"])
	}
}

// TestJSONUnmarshalStruct 测试普通结构体反序列化
// Validates: Requirements 1.2
func TestJSONUnmarshalStruct(t *testing.T) {
	encoder := encodingx.NewJSON()
	jsonData := []byte(`{"integer":42,"string":"hello world","bool":true,"float":3.14159}`)

	var result TestStruct
	err := encoder.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.Integer != 42 {
		t.Errorf("Integer mismatch: expected 42, got %d", result.Integer)
	}
	if result.String != "hello world" {
		t.Errorf("String mismatch: expected 'hello world', got '%s'", result.String)
	}
	if result.Bool != true {
		t.Errorf("Bool mismatch: expected true, got %v", result.Bool)
	}
	if result.Float != 3.14159 {
		t.Errorf("Float mismatch: expected 3.14159, got %f", result.Float)
	}
}

// TestJSONRoundTripStruct 测试结构体序列化后反序列化
// Validates: Requirements 1.1, 1.2
func TestJSONRoundTripStruct(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := TestStruct{
		Integer: 100,
		String:  "round trip test",
		Bool:    false,
		Float:   2.71828,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestJSONMarshalNestedStruct 测试嵌套结构体序列化/反序列化
// Validates: Requirements 1.1, 1.2
func TestJSONMarshalNestedStruct(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := NestedStruct{
		Name: "nested test",
		Inner: TestStruct{
			Integer: 123,
			String:  "inner struct",
			Bool:    true,
			Float:   1.5,
		},
		Slice: []int{1, 2, 3, 4, 5},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result NestedStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Nested struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestJSONMarshalByteSlice 测试 []byte 类型直接返回
// Validates: Requirements 1.3
func TestJSONMarshalByteSlice(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证直接返回原始字节数组
	if !BytesEqual(data, original) {
		t.Errorf("[]byte should be returned directly: expected %v, got %v", original, data)
	}
}

// TestJSONMarshalBytes 测试 Bytes 类型特殊处理
// Validates: Requirements 1.4
func TestJSONMarshalBytes(t *testing.T) {
	encoder := encodingx.NewJSON()
	originalData := []byte{0x10, 0x20, 0x30, 0x40}
	original := encodingx.MakeBytes(originalData)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回 Bytes.Data 字段
	if !BytesEqual(data, originalData) {
		t.Errorf("Bytes should return Data field: expected %v, got %v", originalData, data)
	}
}

// TestJSONMarshalBytesPointer 测试 *Bytes 类型特殊处理
// Validates: Requirements 1.5
func TestJSONMarshalBytesPointer(t *testing.T) {
	encoder := encodingx.NewJSON()
	originalData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	original := encodingx.NewBytes()
	original.Data = originalData

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回 Bytes.Data 字段
	if !BytesEqual(data, originalData) {
		t.Errorf("*Bytes should return Data field: expected %v, got %v", originalData, data)
	}
}

// TestJSONUnmarshalToBytes 测试反序列化到 *Bytes
// Validates: Requirements 1.6
func TestJSONUnmarshalToBytes(t *testing.T) {
	encoder := encodingx.NewJSON()
	inputData := []byte{0x11, 0x22, 0x33, 0x44, 0x55}

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(inputData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证数据存入 Bytes.Data
	if !BytesEqual(result.Data, inputData) {
		t.Errorf("Unmarshal to *Bytes should store data in Data field: expected %v, got %v", inputData, result.Data)
	}
}

// TestJSONUnmarshalToBytesWithJSONData 测试使用 JSON 数据反序列化到 *Bytes
// Validates: Requirements 1.6
func TestJSONUnmarshalToBytesWithJSONData(t *testing.T) {
	encoder := encodingx.NewJSON()
	jsonData := []byte(`{"key":"value"}`)

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(jsonData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证 JSON 数据直接存入 Bytes.Data
	if !BytesEqual(result.Data, jsonData) {
		t.Errorf("Unmarshal to *Bytes should store raw data: expected %v, got %v", jsonData, result.Data)
	}
}

// TestJSONString 测试 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestJSONString(t *testing.T) {
	encoder := encodingx.NewJSON()
	name := encoder.String()

	if name != "JSON" {
		t.Errorf("String() should return 'JSON', got '%s'", name)
	}
}

// TestJSONStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestJSONStyle(t *testing.T) {
	encoder := encodingx.NewJSON()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestJSONReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestJSONReverse(t *testing.T) {
	encoder := encodingx.NewJSON()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 JSON 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestJSONMarshalEmptyStruct 测试空结构体序列化
// Validates: Requirements 1.1
func TestJSONMarshalEmptyStruct(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := TestStruct{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Empty struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestJSONMarshalEmptyByteSlice 测试空 []byte 序列化
// Validates: Requirements 1.3
func TestJSONMarshalEmptyByteSlice(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := []byte{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证直接返回空字节数组
	if len(data) != 0 {
		t.Errorf("Empty []byte should return empty slice, got %v", data)
	}
}

// TestJSONMarshalEmptyBytes 测试空 Bytes 序列化
// Validates: Requirements 1.4
func TestJSONMarshalEmptyBytes(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := encodingx.MakeBytes([]byte{})

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回空字节数组
	if len(data) != 0 {
		t.Errorf("Empty Bytes should return empty slice, got %v", data)
	}
}

// TestJSONMarshalNilBytesPointer 测试 nil Data 的 *Bytes 序列化
// Validates: Requirements 1.5
func TestJSONMarshalNilBytesPointer(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := encodingx.NewBytes() // Data 为 nil

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回 nil
	if data != nil {
		t.Errorf("*Bytes with nil Data should return nil, got %v", data)
	}
}

// TestJSONUnmarshalEmptyToBytes 测试空数据反序列化到 *Bytes
// Validates: Requirements 1.6
func TestJSONUnmarshalEmptyToBytes(t *testing.T) {
	encoder := encodingx.NewJSON()
	inputData := []byte{}

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(inputData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证空数据存入 Bytes.Data
	if len(result.Data) != 0 {
		t.Errorf("Unmarshal empty data to *Bytes should result in empty Data, got %v", result.Data)
	}
}

// TestJSONMarshalSpecialCharacters 测试包含特殊字符的结构体序列化
// Validates: Requirements 1.1, 1.2
func TestJSONMarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := TestStruct{
		Integer: -999,
		String:  "hello\nworld\t\"quoted\"",
		Bool:    true,
		Float:   -0.001,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Special characters round trip failed: original %+v != result %+v", original, result)
	}
}

// TestJSONMarshalUnicodeString 测试包含 Unicode 字符的结构体序列化
// Validates: Requirements 1.1, 1.2
func TestJSONMarshalUnicodeString(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := TestStruct{
		Integer: 123,
		String:  "你好世界 🌍 مرحبا",
		Bool:    false,
		Float:   42.0,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Unicode string round trip failed: original %+v != result %+v", original, result)
	}
}

// TestJSONMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 1.1, 1.2
func TestJSONMarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := TestStruct{
		Integer: 2147483647, // Max int32
		String:  "large numbers",
		Bool:    true,
		Float:   1.7976931348623157e+308, // Max float64
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if original.Integer != result.Integer {
		t.Errorf("Large integer mismatch: expected %d, got %d", original.Integer, result.Integer)
	}
	if original.Float != result.Float {
		t.Errorf("Large float mismatch: expected %e, got %e", original.Float, result.Float)
	}
}

// TestJSONMarshalBinaryData 测试二进制数据序列化
// Validates: Requirements 1.3
func TestJSONMarshalBinaryData(t *testing.T) {
	encoder := encodingx.NewJSON()
	// 包含所有可能的字节值
	original := make([]byte, 256)
	for i := 0; i < 256; i++ {
		original[i] = byte(i)
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证直接返回原始字节数组
	if !BytesEqual(data, original) {
		t.Errorf("Binary data should be returned directly")
	}
}

// TestJSONImplementsEncoding 测试 JSON 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestJSONImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewJSON()

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

// TestJSONMarshalSlice 测试切片序列化
// Validates: Requirements 1.1, 1.2
func TestJSONMarshalSlice(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := []TestStruct{
		{Integer: 1, String: "first", Bool: true, Float: 1.1},
		{Integer: 2, String: "second", Bool: false, Float: 2.2},
		{Integer: 3, String: "third", Bool: true, Float: 3.3},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(original) != len(result) {
		t.Fatalf("Slice length mismatch: expected %d, got %d", len(original), len(result))
	}
	for i := range original {
		if !original[i].Equal(result[i]) {
			t.Errorf("Slice element %d mismatch: expected %+v, got %+v", i, original[i], result[i])
		}
	}
}

// TestJSONMarshalMap 测试 map 序列化
// Validates: Requirements 1.1, 1.2
func TestJSONMarshalMap(t *testing.T) {
	encoder := encodingx.NewJSON()
	original := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result map[string]int
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(original) != len(result) {
		t.Fatalf("Map length mismatch: expected %d, got %d", len(original), len(result))
	}
	for k, v := range original {
		if result[k] != v {
			t.Errorf("Map value mismatch for key '%s': expected %d, got %d", k, v, result[k])
		}
	}
}

// TestJSONUnmarshalInvalidJSON 测试无效 JSON 反序列化
// Validates: Requirements 1.2
func TestJSONUnmarshalInvalidJSON(t *testing.T) {
	encoder := encodingx.NewJSON()
	invalidJSON := []byte(`{"invalid json`)

	var result TestStruct
	err := encoder.Unmarshal(invalidJSON, &result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid JSON")
	}
}

// TestJSONUnmarshalTypeMismatch 测试类型不匹配的反序列化
// Validates: Requirements 1.2
func TestJSONUnmarshalTypeMismatch(t *testing.T) {
	encoder := encodingx.NewJSON()
	// 尝试将数组反序列化到结构体
	jsonData := []byte(`[1, 2, 3]`)

	var result TestStruct
	err := encoder.Unmarshal(jsonData, &result)
	if err == nil {
		t.Error("Unmarshal should fail for type mismatch")
	}
}

// ============================================================================
// JSON 编码器属性测试
// ============================================================================

// TestGroup1_Property_1_JSON_RoundTrip 属性测试：JSON Round-Trip 一致性
// **Property 1: JSON Round-Trip 一致性**
// *For any* 有效的 Go 结构体，使用 JSON 编码器序列化后再反序列化，
// 应该产生与原始结构体等价的对象。
// **Validates: Requirements 1.7**
func TestGroup1_Property_1_JSON_RoundTrip(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewJSON()

	t.Run("TestStruct_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			original := gen.GenerateTestStruct()

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result TestStruct
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Errorf("Iteration %d: Round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("NestedStruct_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			original := gen.GenerateNestedStruct()

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result NestedStruct
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Errorf("Iteration %d: Round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("Map_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机 map
			original := gen.GenerateStringIntMap(1, 10)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result map[string]int
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if len(original) != len(result) {
				t.Errorf("Iteration %d: Map length mismatch: expected %d, got %d", i, len(original), len(result))
				continue
			}
			for k, v := range original {
				if result[k] != v {
					t.Errorf("Iteration %d: Map value mismatch for key '%s': expected %d, got %d", i, k, v, result[k])
				}
			}
		}
	})

	t.Run("Slice_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机切片
			original := gen.GenerateTestStructSlice(1, 10)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []TestStruct
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if len(original) != len(result) {
				t.Errorf("Iteration %d: Slice length mismatch: expected %d, got %d", i, len(original), len(result))
				continue
			}
			for j := range original {
				if !original[j].Equal(result[j]) {
					t.Errorf("Iteration %d: Slice element %d mismatch:\n  original: %+v\n  result:   %+v", i, j, original[j], result[j])
				}
			}
		}
	})
}

// TestGroup1_Property_2_JSON_BytesSpecialHandling 属性测试：JSON Bytes 类型特殊处理
// **Property 2: JSON Bytes 类型特殊处理**
// *For any* []byte、Bytes 或 *Bytes 类型的值，JSON 编码器的 Marshal 应该直接返回其字节数据，
// 而不是 JSON 编码。
// **Validates: Requirements 1.3, 1.4, 1.5, 1.6**
func TestGroup1_Property_2_JSON_BytesSpecialHandling(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewJSON()

	t.Run("ByteSlice_DirectReturn", func(t *testing.T) {
		// **Validates: Requirements 1.3**
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机字节数组
			original := gen.GenerateBytes(0, 100)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 验证：[]byte 应该直接返回，而不是 JSON 编码
			if !BytesEqual(data, original) {
				t.Errorf("Iteration %d: []byte should be returned directly:\n  original: %v\n  result:   %v", i, original, data)
			}
		}
	})

	t.Run("Bytes_DirectReturn", func(t *testing.T) {
		// **Validates: Requirements 1.4**
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机 Bytes
			originalData := gen.GenerateBytes(0, 100)
			original := encodingx.MakeBytes(originalData)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 验证：Bytes 应该返回 Data 字段，而不是 JSON 编码
			if !BytesEqual(data, originalData) {
				t.Errorf("Iteration %d: Bytes should return Data field:\n  original.Data: %v\n  result:        %v", i, originalData, data)
			}
		}
	})

	t.Run("BytesPointer_DirectReturn", func(t *testing.T) {
		// **Validates: Requirements 1.5**
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机 *Bytes
			originalData := gen.GenerateBytes(0, 100)
			original := encodingx.NewBytes()
			original.Data = originalData

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 验证：*Bytes 应该返回 Data 字段，而不是 JSON 编码
			if !BytesEqual(data, originalData) {
				t.Errorf("Iteration %d: *Bytes should return Data field:\n  original.Data: %v\n  result:        %v", i, originalData, data)
			}
		}
	})

	t.Run("BytesPointer_Unmarshal", func(t *testing.T) {
		// **Validates: Requirements 1.6**
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机字节数据
			inputData := gen.GenerateBytes(0, 100)

			// 反序列化到 *Bytes
			result := encodingx.NewBytes()
			err := encoder.Unmarshal(inputData, result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证：数据应该直接存入 Bytes.Data
			if !BytesEqual(result.Data, inputData) {
				t.Errorf("Iteration %d: Unmarshal to *Bytes should store data in Data field:\n  input:       %v\n  result.Data: %v", i, inputData, result.Data)
			}
		}
	})

	t.Run("Bytes_RoundTrip", func(t *testing.T) {
		// 验证 Bytes 类型的 Round-Trip 一致性
		// **Validates: Requirements 1.3, 1.4, 1.5, 1.6**
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机字节数据
			originalData := gen.GenerateBytes(0, 100)

			// 使用 *Bytes 进行 Round-Trip
			original := encodingx.NewBytes()
			original.Data = originalData

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			result := encodingx.NewBytes()
			err = encoder.Unmarshal(data, result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !BytesEqual(original.Data, result.Data) {
				t.Errorf("Iteration %d: Bytes Round-trip failed:\n  original.Data: %v\n  result.Data:   %v", i, original.Data, result.Data)
			}
		}
	})

	t.Run("EmptyBytes_Handling", func(t *testing.T) {
		// 验证空字节数组的特殊处理
		for i := 0; i < iterations; i++ {
			// 测试空 []byte
			emptySlice := []byte{}
			data, err := encoder.Marshal(emptySlice)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal empty []byte failed: %v", i, err)
			}
			if len(data) != 0 {
				t.Errorf("Iteration %d: Empty []byte should return empty slice, got %v", i, data)
			}

			// 测试空 Bytes
			emptyBytes := encodingx.MakeBytes([]byte{})
			data, err = encoder.Marshal(emptyBytes)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal empty Bytes failed: %v", i, err)
			}
			if len(data) != 0 {
				t.Errorf("Iteration %d: Empty Bytes should return empty slice, got %v", i, data)
			}

			// 测试空 *Bytes
			emptyBytesPtr := encodingx.NewBytes()
			emptyBytesPtr.Data = []byte{}
			data, err = encoder.Marshal(emptyBytesPtr)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal empty *Bytes failed: %v", i, err)
			}
			if len(data) != 0 {
				t.Errorf("Iteration %d: Empty *Bytes should return empty slice, got %v", i, data)
			}
		}
	})

	t.Run("NilBytesPointer_Handling", func(t *testing.T) {
		// 验证 nil Data 的 *Bytes 处理
		for i := 0; i < iterations; i++ {
			nilBytesPtr := encodingx.NewBytes() // Data 为 nil
			data, err := encoder.Marshal(nilBytesPtr)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal nil *Bytes failed: %v", i, err)
			}
			if data != nil {
				t.Errorf("Iteration %d: *Bytes with nil Data should return nil, got %v", i, data)
			}
		}
	})
}
