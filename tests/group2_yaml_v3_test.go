package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	"gopkg.in/yaml.v3"
	"pgregory.net/rapid"
)

// ============================================================================
// YAML v3 兼容性测试
// 验证 YAML v2 升级到 v3 后的向后兼容性
// Validates: Requirements 18.1, 18.2, 18.3, 18.4
// ============================================================================

// TestYAMLv3MarshalStruct 测试 YAML v3 普通结构体序列化
// 复用组1测试用例，验证 v3 兼容性
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalStruct(t *testing.T) {
	encoder := encodingx.NewYAML()
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

	// 验证返回的是有效的 YAML v3 格式
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("Result is not valid YAML v3: %v", err)
	}

	// 验证字段值
	if int(result["integer"].(int)) != original.Integer {
		t.Errorf("Integer mismatch: expected %d, got %v", original.Integer, result["integer"])
	}
	if result["string"].(string) != original.String {
		t.Errorf("String mismatch: expected %s, got %v", original.String, result["string"])
	}
	if result["bool"].(bool) != original.Bool {
		t.Errorf("Bool mismatch: expected %v, got %v", original.Bool, result["bool"])
	}
}

// TestYAMLv3UnmarshalStruct 测试 YAML v3 普通结构体反序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3UnmarshalStruct(t *testing.T) {
	encoder := encodingx.NewYAML()
	yamlData := []byte(`integer: 42
string: hello world
bool: true
float: 3.14159`)

	var result TestStruct
	err := encoder.Unmarshal(yamlData, &result)
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

// TestYAMLv3RoundTripStruct 测试 YAML v3 结构体序列化后反序列化
// Validates: Requirements 18.2, 18.3, 18.4
func TestYAMLv3RoundTripStruct(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3MarshalNestedStruct 测试 YAML v3 嵌套结构体序列化/反序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalNestedStruct(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3String 测试 YAML v3 String() 方法返回类型名称
// Validates: Requirements 18.3
func TestYAMLv3String(t *testing.T) {
	encoder := encodingx.NewYAML()
	name := encoder.String()

	if name != "YAML" {
		t.Errorf("String() should return 'YAML', got '%s'", name)
	}
}

// TestYAMLv3Style 测试 YAML v3 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 18.3
func TestYAMLv3Style(t *testing.T) {
	encoder := encodingx.NewYAML()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestYAMLv3Reverse 测试 YAML v3 Reverse() 方法返回自身
// Validates: Requirements 18.3
func TestYAMLv3Reverse(t *testing.T) {
	encoder := encodingx.NewYAML()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 YAML 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestYAMLv3MarshalEmptyStruct 测试 YAML v3 空结构体序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalEmptyStruct(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3MarshalSpecialCharacters 测试 YAML v3 包含特殊字符的结构体序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3MarshalUnicodeString 测试 YAML v3 包含 Unicode 字符的结构体序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalUnicodeString(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3MarshalLargeNumbers 测试 YAML v3 大数值序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3ImplementsEncoding 测试 YAML v3 编码器实现 Encoding 接口
// Validates: Requirements 18.3
func TestYAMLv3ImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewYAML()

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

// TestYAMLv3MarshalSlice 测试 YAML v3 切片序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalSlice(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3MarshalMap 测试 YAML v3 map 序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalMap(t *testing.T) {
	encoder := encodingx.NewYAML()
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

// TestYAMLv3UnmarshalInvalidYAML 测试 YAML v3 无效 YAML 反序列化
// Validates: Requirements 18.2
func TestYAMLv3UnmarshalInvalidYAML(t *testing.T) {
	encoder := encodingx.NewYAML()
	// 无效的 YAML：缩进错误
	invalidYAML := []byte(`
key: value
  invalid: indentation
`)

	var result map[string]interface{}
	err := encoder.Unmarshal(invalidYAML, &result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid YAML")
	}
}

// TestYAMLv3MarshalNegativeNumbers 测试 YAML v3 负数序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalNegativeNumbers(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := TestStruct{
		Integer: -12345,
		String:  "negative test",
		Bool:    false,
		Float:   -99.99,
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
		t.Errorf("Negative numbers round trip failed: original %+v != result %+v", original, result)
	}
}

// TestYAMLv3MarshalZeroValues 测试 YAML v3 零值序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalZeroValues(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := TestStruct{
		Integer: 0,
		String:  "",
		Bool:    false,
		Float:   0.0,
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
		t.Errorf("Zero values round trip failed: original %+v != result %+v", original, result)
	}
}

// TestYAMLv3MarshalDeeplyNestedStruct 测试 YAML v3 深度嵌套结构体
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalDeeplyNestedStruct(t *testing.T) {
	encoder := encodingx.NewYAML()

	// 创建深度嵌套的结构
	type Level3 struct {
		Value string `yaml:"value"`
	}
	type Level2 struct {
		Level3 Level3 `yaml:"level3"`
	}
	type Level1 struct {
		Level2 Level2 `yaml:"level2"`
	}

	original := Level1{
		Level2: Level2{
			Level3: Level3{
				Value: "deep value",
			},
		},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result Level1
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if original.Level2.Level3.Value != result.Level2.Level3.Value {
		t.Errorf("Deeply nested struct round trip failed: expected '%s', got '%s'",
			original.Level2.Level3.Value, result.Level2.Level3.Value)
	}
}

// TestYAMLv3MarshalSliceOfMaps 测试 YAML v3 map 切片序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalSliceOfMaps(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := []map[string]string{
		{"key1": "value1", "key2": "value2"},
		{"key3": "value3", "key4": "value4"},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []map[string]string
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(original) != len(result) {
		t.Fatalf("Slice length mismatch: expected %d, got %d", len(original), len(result))
	}
	for i := range original {
		for k, v := range original[i] {
			if result[i][k] != v {
				t.Errorf("Map value mismatch at index %d, key '%s': expected '%s', got '%s'",
					i, k, v, result[i][k])
			}
		}
	}
}

// TestYAMLv3MarshalPointerFields 测试 YAML v3 包含指针字段的结构体
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalPointerFields(t *testing.T) {
	encoder := encodingx.NewYAML()

	type StructWithPointers struct {
		IntPtr    *int    `yaml:"int_ptr"`
		StringPtr *string `yaml:"string_ptr"`
	}

	intVal := 42
	strVal := "pointer value"
	original := StructWithPointers{
		IntPtr:    &intVal,
		StringPtr: &strVal,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result StructWithPointers
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if result.IntPtr == nil || *result.IntPtr != intVal {
		t.Errorf("IntPtr mismatch: expected %d, got %v", intVal, result.IntPtr)
	}
	if result.StringPtr == nil || *result.StringPtr != strVal {
		t.Errorf("StringPtr mismatch: expected '%s', got %v", strVal, result.StringPtr)
	}
}

// TestYAMLv3MarshalNilPointerFields 测试 YAML v3 包含 nil 指针字段的结构体
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalNilPointerFields(t *testing.T) {
	encoder := encodingx.NewYAML()

	type StructWithPointers struct {
		IntPtr    *int    `yaml:"int_ptr,omitempty"`
		StringPtr *string `yaml:"string_ptr,omitempty"`
	}

	original := StructWithPointers{
		IntPtr:    nil,
		StringPtr: nil,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result StructWithPointers
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if result.IntPtr != nil {
		t.Errorf("IntPtr should be nil, got %v", result.IntPtr)
	}
	if result.StringPtr != nil {
		t.Errorf("StringPtr should be nil, got %v", result.StringPtr)
	}
}

// TestYAMLv3MarshalEmptySlice 测试 YAML v3 空切片序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalEmptySlice(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := []TestStruct{}

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
	if len(result) != 0 {
		t.Errorf("Empty slice round trip failed: expected empty slice, got %v", result)
	}
}

// TestYAMLv3MarshalEmptyMap 测试 YAML v3 空 map 序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalEmptyMap(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := map[string]int{}

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
	if len(result) != 0 {
		t.Errorf("Empty map round trip failed: expected empty map, got %v", result)
	}
}

// TestYAMLv3MarshalMultilineString 测试 YAML v3 多行字符串序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalMultilineString(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := TestStruct{
		Integer: 1,
		String:  "line1\nline2\nline3",
		Bool:    true,
		Float:   1.0,
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
		t.Errorf("Multiline string round trip failed: original %+v != result %+v", original, result)
	}
}

// TestYAMLv3MarshalYAMLSpecialChars 测试 YAML v3 特殊字符序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalYAMLSpecialChars(t *testing.T) {
	encoder := encodingx.NewYAML()
	// YAML 特殊字符：冒号、井号、破折号等
	original := TestStruct{
		Integer: 1,
		String:  "key: value # comment - item",
		Bool:    true,
		Float:   1.0,
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
		t.Errorf("YAML special chars round trip failed: original %+v != result %+v", original, result)
	}
}

// TestYAMLv3MarshalQuotedString 测试 YAML v3 带引号的字符串序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalQuotedString(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := TestStruct{
		Integer: 1,
		String:  `"double quoted" and 'single quoted'`,
		Bool:    true,
		Float:   1.0,
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
		t.Errorf("Quoted string round trip failed: original %+v != result %+v", original, result)
	}
}

// TestYAMLv3MarshalFloatSpecialValues 测试 YAML v3 浮点数特殊值序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalFloatSpecialValues(t *testing.T) {
	encoder := encodingx.NewYAML()

	testCases := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"negative_zero", -0.0},
		{"small_positive", 0.000001},
		{"small_negative", -0.000001},
		{"large_positive", 1e100},
		{"large_negative", -1e100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := TestStruct{
				Integer: 1,
				String:  "float test",
				Bool:    true,
				Float:   tc.value,
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
			if original.Float != result.Float {
				t.Errorf("Float value mismatch: expected %v, got %v", original.Float, result.Float)
			}
		})
	}
}

// TestYAMLv3MarshalIntegerBoundaries 测试 YAML v3 整数边界值序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalIntegerBoundaries(t *testing.T) {
	encoder := encodingx.NewYAML()

	testCases := []struct {
		name  string
		value int
	}{
		{"zero", 0},
		{"one", 1},
		{"negative_one", -1},
		{"max_int32", 2147483647},
		{"min_int32", -2147483648},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := TestStruct{
				Integer: tc.value,
				String:  "integer test",
				Bool:    true,
				Float:   1.0,
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
				t.Errorf("Integer value mismatch: expected %d, got %d", original.Integer, result.Integer)
			}
		})
	}
}

// TestYAMLv3MarshalBoolValues 测试 YAML v3 布尔值序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalBoolValues(t *testing.T) {
	encoder := encodingx.NewYAML()

	testCases := []struct {
		name  string
		value bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := TestStruct{
				Integer: 1,
				String:  "bool test",
				Bool:    tc.value,
				Float:   1.0,
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
			if original.Bool != result.Bool {
				t.Errorf("Bool value mismatch: expected %v, got %v", original.Bool, result.Bool)
			}
		})
	}
}

// TestYAMLv3NewYAMLConstructor 测试 YAML v3 NewYAML 构造函数
// Validates: Requirements 18.3
func TestYAMLv3NewYAMLConstructor(t *testing.T) {
	encoder := encodingx.NewYAML()

	if encoder == nil {
		t.Fatal("NewYAML() should return non-nil encoder")
	}

	// 验证返回的编码器可以正常工作
	original := TestStruct{Integer: 1, String: "test", Bool: true, Float: 1.0}
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !original.Equal(result) {
		t.Errorf("NewYAML encoder round trip failed")
	}
}

// TestYAMLv3MarshalInterfaceSlice 测试 YAML v3 interface{} 切片序列化
// Validates: Requirements 18.2, 18.3
func TestYAMLv3MarshalInterfaceSlice(t *testing.T) {
	encoder := encodingx.NewYAML()
	original := []interface{}{
		"string",
		42,
		true,
		3.14,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []interface{}
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证长度
	if len(original) != len(result) {
		t.Fatalf("Slice length mismatch: expected %d, got %d", len(original), len(result))
	}

	// 验证第一个元素（字符串）
	if result[0].(string) != original[0].(string) {
		t.Errorf("String element mismatch: expected '%v', got '%v'", original[0], result[0])
	}
}

// TestYAMLv3RegisteredInEncodingSet 测试 YAML v3 编码器注册到 EncodingSet
// 通过 ChainEncoding 间接测试
// Validates: Requirements 18.3
func TestYAMLv3RegisteredInEncodingSet(t *testing.T) {
	// 创建使用 YAML 编码器的 ChainEncoding
	chain := encodingx.NewChainEncoding([]string{"YAML"}, []string{"YAML"})

	// 测试 Marshal - 如果 localEncoding 找不到 YAML，会返回错误
	input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
	data, err := chain.Marshal(input)
	if err != nil {
		t.Fatalf("ChainEncoding with YAML should work, YAML not registered: %v", err)
	}

	// 测试 Unmarshal
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("ChainEncoding Unmarshal with YAML should work: %v", err)
	}

	// 验证数据正确
	if !input.Equal(result) {
		t.Errorf("ChainEncoding round trip failed: expected %+v, got %+v", input, result)
	}
}

// ============================================================================
// YAML v3 兼容性属性测试
// ============================================================================

// genYAMLSafeString 生成 YAML 安全的字符串（避免特殊字符导致解析问题）
// YAML 对某些特殊字符（如换行符在字符串开头）有特殊处理
func genYAMLSafeString(minLen, maxLen int) *rapid.Generator[string] {
	// 使用字母数字和常见安全字符
	const safeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 _-"
	return rapid.Custom(func(t *rapid.T) string {
		length := rapid.IntRange(minLen, maxLen).Draw(t, "length")
		if length == 0 {
			return ""
		}
		result := make([]byte, length)
		for i := 0; i < length; i++ {
			idx := rapid.IntRange(0, len(safeChars)-1).Draw(t, "charIdx")
			result[i] = safeChars[idx]
		}
		return string(result)
	})
}

// genYAMLv3TestStruct 生成随机 TestStruct 的 rapid 生成器
// 使用 YAML 安全的字符串以确保向后兼容性测试的准确性
func genYAMLv3TestStruct() *rapid.Generator[TestStruct] {
	return rapid.Custom(func(t *rapid.T) TestStruct {
		return TestStruct{
			Integer: rapid.IntRange(-10000, 10000).Draw(t, "integer"),
			String:  genYAMLSafeString(0, 100).Draw(t, "string"),
			Bool:    rapid.Bool().Draw(t, "bool"),
			Float:   rapid.Float64().Draw(t, "float"),
		}
	})
}

// genYAMLv3NestedStruct 生成随机 NestedStruct 的 rapid 生成器
func genYAMLv3NestedStruct() *rapid.Generator[NestedStruct] {
	return rapid.Custom(func(t *rapid.T) NestedStruct {
		sliceLen := rapid.IntRange(0, 20).Draw(t, "sliceLen")
		slice := make([]int, sliceLen)
		for i := 0; i < sliceLen; i++ {
			slice[i] = rapid.IntRange(-1000, 1000).Draw(t, "sliceElem")
		}
		return NestedStruct{
			Name:  genYAMLSafeString(0, 50).Draw(t, "name"),
			Inner: genYAMLv3TestStruct().Draw(t, "inner"),
			Slice: slice,
		}
	})
}

// genYAMLv3StringIntMap 生成随机 map[string]int 的 rapid 生成器
// 使用 YAML 安全的字符串作为 map 键
func genYAMLv3StringIntMap() *rapid.Generator[map[string]int] {
	return rapid.Custom(func(t *rapid.T) map[string]int {
		numKeys := rapid.IntRange(0, 20).Draw(t, "numKeys")
		result := make(map[string]int)
		for i := 0; i < numKeys; i++ {
			// 使用 YAML 安全的字符串作为键
			key := genYAMLSafeString(1, 20).Draw(t, "key")
			value := rapid.IntRange(-10000, 10000).Draw(t, "value")
			result[key] = value
		}
		return result
	})
}

// TestProperty25_YAMLv3BackwardCompatibility 属性测试：YAML v3 向后兼容性
// **Property 25: YAML v3 向后兼容性**
// *For any* 在 YAML v2 下有效的结构体，升级到 YAML v3 后序列化和反序列化行为应该保持一致。
// **Validates: Requirements 18.1, 18.2, 18.3, 18.4**
func TestProperty25_YAMLv3BackwardCompatibility(t *testing.T) {
	encoder := encodingx.NewYAML()

	t.Run("TestStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genYAMLv3TestStruct().Draw(t, "original")

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

			// 验证 Round-Trip 一致性（向后兼容性）
			if !original.Equal(result) {
				t.Fatalf("YAML v3 backward compatibility failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("NestedStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genYAMLv3NestedStruct().Draw(t, "original")

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

			// 验证 Round-Trip 一致性（向后兼容性）
			if !original.Equal(result) {
				t.Fatalf("YAML v3 backward compatibility failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("Map_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genYAMLv3StringIntMap().Draw(t, "original")

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

			// 验证 Round-Trip 一致性（向后兼容性）
			if len(original) != len(result) {
				t.Fatalf("Map length mismatch: expected %d, got %d", len(original), len(result))
			}
			for k, v := range original {
				if result[k] != v {
					t.Fatalf("Map value mismatch for key '%s': expected %d, got %d", k, v, result[k])
				}
			}
		})
	})

	t.Run("Slice_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// 生成随机切片
			sliceLen := rapid.IntRange(0, 20).Draw(t, "sliceLen")
			original := make([]TestStruct, sliceLen)
			for i := 0; i < sliceLen; i++ {
				original[i] = genYAMLv3TestStruct().Draw(t, "elem")
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

			// 验证 Round-Trip 一致性（向后兼容性）
			if len(original) != len(result) {
				t.Fatalf("Slice length mismatch: expected %d, got %d", len(original), len(result))
			}
			for i := range original {
				if !original[i].Equal(result[i]) {
					t.Fatalf("Slice element %d mismatch:\n  original: %+v\n  result:   %+v", i, original[i], result[i])
				}
			}
		})
	})

	t.Run("PrimitiveTypes_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// 测试各种基本类型的向后兼容性
			intVal := rapid.Int().Draw(t, "int")
			stringVal := genYAMLSafeString(0, 100).Draw(t, "string")
			boolVal := rapid.Bool().Draw(t, "bool")
			float64Val := rapid.Float64().Draw(t, "float64")

			// 测试 int
			intData, err := encoder.Marshal(intVal)
			if err != nil {
				t.Fatalf("Marshal int failed: %v", err)
			}
			var intResult int
			err = encoder.Unmarshal(intData, &intResult)
			if err != nil {
				t.Fatalf("Unmarshal int failed: %v", err)
			}
			if intVal != intResult {
				t.Fatalf("Int round-trip failed: expected %d, got %d", intVal, intResult)
			}

			// 测试 string
			stringData, err := encoder.Marshal(stringVal)
			if err != nil {
				t.Fatalf("Marshal string failed: %v", err)
			}
			var stringResult string
			err = encoder.Unmarshal(stringData, &stringResult)
			if err != nil {
				t.Fatalf("Unmarshal string failed: %v", err)
			}
			if stringVal != stringResult {
				t.Fatalf("String round-trip failed: expected '%s', got '%s'", stringVal, stringResult)
			}

			// 测试 bool
			boolData, err := encoder.Marshal(boolVal)
			if err != nil {
				t.Fatalf("Marshal bool failed: %v", err)
			}
			var boolResult bool
			err = encoder.Unmarshal(boolData, &boolResult)
			if err != nil {
				t.Fatalf("Unmarshal bool failed: %v", err)
			}
			if boolVal != boolResult {
				t.Fatalf("Bool round-trip failed: expected %v, got %v", boolVal, boolResult)
			}

			// 测试 float64
			float64Data, err := encoder.Marshal(float64Val)
			if err != nil {
				t.Fatalf("Marshal float64 failed: %v", err)
			}
			var float64Result float64
			err = encoder.Unmarshal(float64Data, &float64Result)
			if err != nil {
				t.Fatalf("Unmarshal float64 failed: %v", err)
			}
			if float64Val != float64Result {
				t.Fatalf("Float64 round-trip failed: expected %f, got %f", float64Val, float64Result)
			}
		})
	})

	t.Run("EmptyStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// 测试空结构体的向后兼容性
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

			// 验证 Round-Trip 一致性（向后兼容性）
			if !original.Equal(result) {
				t.Fatalf("Empty struct backward compatibility failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("EmptySlice_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// 测试空切片的向后兼容性
			original := []TestStruct{}

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

			// 验证 Round-Trip 一致性（向后兼容性）
			if len(result) != 0 {
				t.Fatalf("Empty slice backward compatibility failed: expected empty slice, got %v", result)
			}
		})
	})

	t.Run("EmptyMap_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// 测试空 map 的向后兼容性
			original := map[string]int{}

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

			// 验证 Round-Trip 一致性（向后兼容性）
			if len(result) != 0 {
				t.Fatalf("Empty map backward compatibility failed: expected empty map, got %v", result)
			}
		})
	})
}
