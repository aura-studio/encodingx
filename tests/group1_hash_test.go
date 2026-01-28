package encodingx_test

import (
	"encoding/json"
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// Hash 编码器单元测试
// Validates: Requirements 8.1, 8.2, 8.3, 14.1, 14.2, 14.3
// ============================================================================

// ============================================================================
// Hash 编码器序列化测试
// ============================================================================

// TestHashMarshalHashMarshaller 测试 Hash 编码器序列化实现 HashMarshaller 接口的类型
// Validates: Requirements 8.1
func TestHashMarshalHashMarshaller(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"key1", "value1"},
		{"key2", "value2"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 JSON 格式的键值对数组
	var result [][]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// 验证键值对数量
	if len(result) != len(pairs) {
		t.Errorf("Expected %d pairs, got %d", len(pairs), len(result))
	}

	// 验证键值对内容
	for i, pair := range result {
		if len(pair) != 2 {
			t.Errorf("Pair %d should have 2 elements, got %d", i, len(pair))
			continue
		}
		expectedKey := pairs[i][0].(string)
		expectedValue := pairs[i][1].(string)
		if pair[0] != expectedKey {
			t.Errorf("Pair %d key mismatch: expected %s, got %v", i, expectedKey, pair[0])
		}
		if pair[1] != expectedValue {
			t.Errorf("Pair %d value mismatch: expected %s, got %v", i, expectedValue, pair[1])
		}
	}
}

// TestHashMarshalEmptyPairs 测试 Hash 编码器序列化空键值对
// Validates: Requirements 8.1
func TestHashMarshalEmptyPairs(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是空 JSON 数组
	var result [][]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty array, got %d elements", len(result))
	}
}

// TestHashMarshalSinglePair 测试 Hash 编码器序列化单个键值对
// Validates: Requirements 8.1
func TestHashMarshalSinglePair(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"singleKey", "singleValue"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 JSON 格式的键值对数组
	var result [][]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 pair, got %d", len(result))
	}

	if result[0][0] != "singleKey" || result[0][1] != "singleValue" {
		t.Errorf("Pair content mismatch: expected [singleKey, singleValue], got %v", result[0])
	}
}

// TestHashMarshalMultiplePairs 测试 Hash 编码器序列化多个键值对
// Validates: Requirements 8.1
func TestHashMarshalMultiplePairs(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"name", "Alice"},
		{"age", "30"},
		{"city", "Beijing"},
		{"country", "China"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 JSON 格式的键值对数组
	var result [][]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 pairs, got %d", len(result))
	}

	// 验证每个键值对
	for i, pair := range result {
		expectedKey := pairs[i][0].(string)
		expectedValue := pairs[i][1].(string)
		if pair[0] != expectedKey {
			t.Errorf("Pair %d key mismatch: expected %s, got %v", i, expectedKey, pair[0])
		}
		if pair[1] != expectedValue {
			t.Errorf("Pair %d value mismatch: expected %s, got %v", i, expectedValue, pair[1])
		}
	}
}

// ============================================================================
// Hash 编码器反序列化测试
// ============================================================================

// TestHashUnmarshalHashUnmarshaller 测试 Hash 编码器反序列化到实现 HashUnmarshaller 接口的类型
// Validates: Requirements 8.2
func TestHashUnmarshalHashUnmarshaller(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"key1", "value1"},
		{"key2", "value2"},
	}

	// 创建 JSON 格式的键值对数组
	data, err := json.Marshal(pairs)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证 UnmarshalHash 方法被正确调用
	resultPairs := result.GetPairs()
	if len(resultPairs) != len(pairs) {
		t.Errorf("Expected %d pairs, got %d", len(pairs), len(resultPairs))
	}

	// 验证键值对内容
	for i, pair := range resultPairs {
		if len(pair) != 2 {
			t.Errorf("Pair %d should have 2 elements, got %d", i, len(pair))
			continue
		}
		expectedKey := pairs[i][0].(string)
		expectedValue := pairs[i][1].(string)
		if pair[0] != expectedKey {
			t.Errorf("Pair %d key mismatch: expected %s, got %v", i, expectedKey, pair[0])
		}
		if pair[1] != expectedValue {
			t.Errorf("Pair %d value mismatch: expected %s, got %v", i, expectedValue, pair[1])
		}
	}
}

// TestHashUnmarshalEmptyPairs 测试 Hash 编码器反序列化空键值对
// Validates: Requirements 8.2
func TestHashUnmarshalEmptyPairs(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{}

	// 创建 JSON 格式的空数组
	data, err := json.Marshal(pairs)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证结果为空
	resultPairs := result.GetPairs()
	if len(resultPairs) != 0 {
		t.Errorf("Expected empty pairs, got %d", len(resultPairs))
	}
}

// TestHashUnmarshalSinglePair 测试 Hash 编码器反序列化单个键值对
// Validates: Requirements 8.2
func TestHashUnmarshalSinglePair(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"onlyKey", "onlyValue"},
	}

	// 创建 JSON 格式的键值对数组
	data, err := json.Marshal(pairs)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证结果
	resultPairs := result.GetPairs()
	if len(resultPairs) != 1 {
		t.Errorf("Expected 1 pair, got %d", len(resultPairs))
	}

	if resultPairs[0][0] != "onlyKey" || resultPairs[0][1] != "onlyValue" {
		t.Errorf("Pair content mismatch: expected [onlyKey, onlyValue], got %v", resultPairs[0])
	}
}

// ============================================================================
// Hash 编码器 Round-Trip 测试
// ============================================================================

// TestHashRoundTrip 测试 Hash 编码器序列化/反序列化往返
// Validates: Requirements 8.1, 8.2
func TestHashRoundTrip(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"username", "john_doe"},
		{"email", "john@example.com"},
		{"role", "admin"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	resultPairs := result.GetPairs()
	if !HashPairsEqual(pairs, resultPairs) {
		t.Errorf("Round trip failed: expected %v, got %v", pairs, resultPairs)
	}
}

// TestHashRoundTripEmptyPairs 测试空键值对的往返
// Validates: Requirements 8.1, 8.2
func TestHashRoundTripEmptyPairs(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	resultPairs := result.GetPairs()
	if len(resultPairs) != 0 {
		t.Errorf("Round trip failed for empty pairs: expected empty, got %v", resultPairs)
	}
}

// TestHashRoundTripManyPairs 测试多个键值对的往返
// Validates: Requirements 8.1, 8.2
func TestHashRoundTripManyPairs(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"field1", "value1"},
		{"field2", "value2"},
		{"field3", "value3"},
		{"field4", "value4"},
		{"field5", "value5"},
		{"field6", "value6"},
		{"field7", "value7"},
		{"field8", "value8"},
		{"field9", "value9"},
		{"field10", "value10"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	resultPairs := result.GetPairs()
	if !HashPairsEqual(pairs, resultPairs) {
		t.Errorf("Round trip failed for many pairs: expected %v, got %v", pairs, resultPairs)
	}
}

// ============================================================================
// Hash 编码器未实现接口类型测试
// ============================================================================

// TestHashMarshalNonHashMarshaller 测试 Hash 编码器序列化未实现 HashMarshaller 的类型返回 nil
// Validates: Requirements 8.3
func TestHashMarshalNonHashMarshaller(t *testing.T) {
	encoder := encodingx.NewHash()

	// 测试普通结构体
	ts := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
	data, err := encoder.Marshal(ts)
	if err != nil {
		t.Fatalf("Marshal should not return error for non-HashMarshaller, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for non-HashMarshaller, got: %v", data)
	}
}

// TestHashMarshalString 测试 Hash 编码器序列化字符串类型返回 nil
// Validates: Requirements 8.3
func TestHashMarshalString(t *testing.T) {
	encoder := encodingx.NewHash()

	data, err := encoder.Marshal("hello world")
	if err != nil {
		t.Fatalf("Marshal should not return error for string, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for string, got: %v", data)
	}
}

// TestHashMarshalInt 测试 Hash 编码器序列化整数类型返回 nil
// Validates: Requirements 8.3
func TestHashMarshalInt(t *testing.T) {
	encoder := encodingx.NewHash()

	data, err := encoder.Marshal(42)
	if err != nil {
		t.Fatalf("Marshal should not return error for int, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for int, got: %v", data)
	}
}

// TestHashMarshalSlice 测试 Hash 编码器序列化切片类型返回 nil
// Validates: Requirements 8.3
func TestHashMarshalSlice(t *testing.T) {
	encoder := encodingx.NewHash()

	data, err := encoder.Marshal([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("Marshal should not return error for slice, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for slice, got: %v", data)
	}
}

// TestHashMarshalMap 测试 Hash 编码器序列化 map 类型返回 nil
// Validates: Requirements 8.3
func TestHashMarshalMap(t *testing.T) {
	encoder := encodingx.NewHash()

	data, err := encoder.Marshal(map[string]int{"key": 1})
	if err != nil {
		t.Fatalf("Marshal should not return error for map, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for map, got: %v", data)
	}
}

// TestHashMarshalByteSlice 测试 Hash 编码器序列化 []byte 类型返回 nil
// Validates: Requirements 8.3
func TestHashMarshalByteSlice(t *testing.T) {
	encoder := encodingx.NewHash()

	data, err := encoder.Marshal([]byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("Marshal should not return error for []byte, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for []byte, got: %v", data)
	}
}

// TestHashMarshalNil 测试 Hash 编码器序列化 nil 返回 nil
// Validates: Requirements 8.3
func TestHashMarshalNil(t *testing.T) {
	encoder := encodingx.NewHash()

	data, err := encoder.Marshal(nil)
	if err != nil {
		t.Fatalf("Marshal should not return error for nil, got: %v", err)
	}
	if data != nil {
		t.Errorf("Marshal should return nil for nil, got: %v", data)
	}
}

// TestHashUnmarshalNonHashUnmarshaller 测试 Hash 编码器反序列化到未实现 HashUnmarshaller 的类型
// Validates: Requirements 8.3
func TestHashUnmarshalNonHashUnmarshaller(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"key", "value"},
	}
	data, _ := json.Marshal(pairs)

	// 测试反序列化到普通结构体
	var ts TestStruct
	err := encoder.Unmarshal(data, &ts)
	if err != nil {
		t.Fatalf("Unmarshal should not return error for non-HashUnmarshaller, got: %v", err)
	}
	// 结构体应该保持零值
	if ts.Integer != 0 || ts.String != "" || ts.Bool != false || ts.Float != 0 {
		t.Errorf("Unmarshal should not modify non-HashUnmarshaller struct")
	}
}

// ============================================================================
// Hash 编码器接口方法测试
// ============================================================================

// TestHashString 测试 Hash String() 方法返回类型名称
// Validates: Requirements 14.1
func TestHashString(t *testing.T) {
	encoder := encodingx.NewHash()
	name := encoder.String()

	if name != "Hash" {
		t.Errorf("String() should return 'Hash', got '%s'", name)
	}
}

// TestHashStyle 测试 Hash Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestHashStyle(t *testing.T) {
	encoder := encodingx.NewHash()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestHashReverse 测试 Hash Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestHashReverse(t *testing.T) {
	encoder := encodingx.NewHash()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 Hash 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestHashImplementsEncoding 测试 Hash 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestHashImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewHash()

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
// Hash 编码器边界条件测试
// ============================================================================

// TestHashMarshalSpecialCharacters 测试包含特殊字符的键值对
// Validates: Requirements 8.1
func TestHashMarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"key with spaces", "value with spaces"},
		{"key\"with\"quotes", "value\"with\"quotes"},
		{"key\nwith\nnewlines", "value\nwith\nnewlines"},
		{"key\twith\ttabs", "value\twith\ttabs"},
		{"中文键", "中文值"},
		{"emoji🎉", "emoji🚀"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	resultPairs := result.GetPairs()
	if !HashPairsEqual(pairs, resultPairs) {
		t.Errorf("Round trip failed for special characters: expected %v, got %v", pairs, resultPairs)
	}
}

// TestHashMarshalEmptyStrings 测试空字符串键值对
// Validates: Requirements 8.1
func TestHashMarshalEmptyStrings(t *testing.T) {
	encoder := encodingx.NewHash()
	pairs := [][]interface{}{
		{"", ""},
		{"emptyValue", ""},
		{"", "emptyKey"},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	resultPairs := result.GetPairs()
	if !HashPairsEqual(pairs, resultPairs) {
		t.Errorf("Round trip failed for empty strings: expected %v, got %v", pairs, resultPairs)
	}
}

// TestHashMarshalLongStrings 测试长字符串键值对
// Validates: Requirements 8.1
func TestHashMarshalLongStrings(t *testing.T) {
	encoder := encodingx.NewHash()

	// 生成长字符串
	longKey := make([]byte, 1000)
	longValue := make([]byte, 1000)
	for i := range longKey {
		longKey[i] = 'a' + byte(i%26)
		longValue[i] = 'A' + byte(i%26)
	}

	pairs := [][]interface{}{
		{string(longKey), string(longValue)},
	}
	original := NewHashableStruct(pairs)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := NewHashableStruct(nil)
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证往返一致性
	resultPairs := result.GetPairs()
	if !HashPairsEqual(pairs, resultPairs) {
		t.Errorf("Round trip failed for long strings")
	}
}

// TestHashUnmarshalInvalidJSON 测试反序列化无效 JSON
// Validates: Requirements 8.2
func TestHashUnmarshalInvalidJSON(t *testing.T) {
	encoder := encodingx.NewHash()
	invalidData := []byte("not valid json")

	result := NewHashableStruct(nil)
	err := encoder.Unmarshal(invalidData, result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid JSON")
	}
}

// TestHashUnmarshalMalformedJSON 测试反序列化格式错误的 JSON
// Validates: Requirements 8.2
func TestHashUnmarshalMalformedJSON(t *testing.T) {
	encoder := encodingx.NewHash()
	malformedData := []byte(`{"key": "value"}`) // 不是数组格式

	result := NewHashableStruct(nil)
	err := encoder.Unmarshal(malformedData, result)
	if err == nil {
		t.Error("Unmarshal should fail for malformed JSON (not an array)")
	}
}

// ============================================================================
// Hash 编码器使用随机数据测试
// ============================================================================

// TestHashRoundTripRandomData 测试使用随机数据的往返
// Validates: Requirements 8.1, 8.2
func TestHashRoundTripRandomData(t *testing.T) {
	encoder := encodingx.NewHash()
	gen := NewTestDataGenerator()

	// 运行多次随机测试
	for i := 0; i < 10; i++ {
		pairs := gen.GenerateHashPairs(1, 10)
		original := NewHashableStruct(pairs)

		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed for iteration %d: %v", i, err)
		}

		// 反序列化
		result := NewHashableStruct(nil)
		err = encoder.Unmarshal(data, result)
		if err != nil {
			t.Fatalf("Unmarshal failed for iteration %d: %v", i, err)
		}

		// 验证往返一致性
		resultPairs := result.GetPairs()
		if !HashPairsEqual(pairs, resultPairs) {
			t.Errorf("Round trip failed for iteration %d: expected %v, got %v", i, pairs, resultPairs)
		}
	}
}

// ============================================================================
// Hash 编码器 NewHash 构造函数测试
// ============================================================================

// TestNewHash 测试 NewHash 构造函数
// Validates: Requirements 14.1, 14.2, 14.3
func TestNewHash(t *testing.T) {
	encoder := encodingx.NewHash()

	if encoder == nil {
		t.Fatal("NewHash() should return non-nil encoder")
	}

	// 验证返回的是 Hash 编码器
	if encoder.String() != "Hash" {
		t.Errorf("NewHash() should return Hash encoder, got %s", encoder.String())
	}
}

// TestHashMultipleInstances 测试多个 Hash 实例的独立性
// Validates: Requirements 14.1, 14.2, 14.3
func TestHashMultipleInstances(t *testing.T) {
	encoder1 := encodingx.NewHash()
	encoder2 := encodingx.NewHash()

	// 两个实例应该有相同的行为
	if encoder1.String() != encoder2.String() {
		t.Errorf("Multiple Hash instances should have same String(): %s vs %s",
			encoder1.String(), encoder2.String())
	}

	if encoder1.Style() != encoder2.Style() {
		t.Errorf("Multiple Hash instances should have same Style(): %v vs %v",
			encoder1.Style(), encoder2.Style())
	}

	// 测试两个实例的序列化结果相同
	pairs := [][]interface{}{{"key", "value"}}
	original := NewHashableStruct(pairs)

	data1, err1 := encoder1.Marshal(original)
	data2, err2 := encoder2.Marshal(original)

	if err1 != nil || err2 != nil {
		t.Fatalf("Marshal failed: err1=%v, err2=%v", err1, err2)
	}

	if !BytesEqual(data1, data2) {
		t.Errorf("Multiple Hash instances should produce same output: %s vs %s",
			string(data1), string(data2))
	}
}

// ============================================================================
// Hash 编码器属性测试
// ============================================================================

// genHashPairs 生成随机的 Hash 键值对切片
// 生成 0-10 个键值对，每个键值对包含字母数字字符串
func genHashPairs() *rapid.Generator[[][]interface{}] {
	return rapid.Custom(func(t *rapid.T) [][]interface{} {
		numPairs := rapid.IntRange(0, 10).Draw(t, "numPairs")
		pairs := make([][]interface{}, numPairs)
		for i := 0; i < numPairs; i++ {
			key := rapid.StringMatching(`[a-zA-Z0-9]{1,20}`).Draw(t, "key")
			value := rapid.StringMatching(`[a-zA-Z0-9]{1,30}`).Draw(t, "value")
			pairs[i] = []interface{}{key, value}
		}
		return pairs
	})
}

// TestProperty13_HashRoundTripConsistency 测试 Hash Round-Trip 一致性
// **Property 13: Hash Round-Trip 一致性**
// *For any* 实现 HashMarshaller 和 HashUnmarshaller 接口的类型，
// 使用 Hash 编码器序列化后再反序列化，应该产生与原始数据等价的对象。
// **Validates: Requirements 8.4**
func TestProperty13_HashRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewHash()

	rapid.Check(t, func(t *rapid.T) {
		// 生成随机的键值对
		pairs := genHashPairs().Draw(t, "pairs")
		original := NewHashableStruct(pairs)

		// 序列化
		data, err := encoder.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// 反序列化
		result := NewHashableStruct(nil)
		err = encoder.Unmarshal(data, result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// 验证 Round-Trip 一致性
		resultPairs := result.GetPairs()
		if !HashPairsEqual(pairs, resultPairs) {
			t.Fatalf("Round-trip failed: original %v, got %v", pairs, resultPairs)
		}
	})
}
