package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	"gopkg.in/yaml.v3"
)

// ============================================================================
// YAML 编码器单元测试
// Validates: Requirements 2.1, 2.2, 14.1, 14.2, 14.3
// ============================================================================

// TestYAMLMarshalStruct 测试普通结构体序列化
// Validates: Requirements 2.1
func TestYAMLMarshalStruct(t *testing.T) {
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

	// 验证返回的是有效的 YAML
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("Result is not valid YAML: %v", err)
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

// TestYAMLUnmarshalStruct 测试普通结构体反序列化
// Validates: Requirements 2.2
func TestYAMLUnmarshalStruct(t *testing.T) {
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

// TestYAMLRoundTripStruct 测试结构体序列化后反序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLRoundTripStruct(t *testing.T) {
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

// TestYAMLMarshalNestedStruct 测试嵌套结构体序列化/反序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalNestedStruct(t *testing.T) {
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

// TestYAMLString 测试 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestYAMLString(t *testing.T) {
	encoder := encodingx.NewYAML()
	name := encoder.String()

	if name != "YAML" {
		t.Errorf("String() should return 'YAML', got '%s'", name)
	}
}

// TestYAMLStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestYAMLStyle(t *testing.T) {
	encoder := encodingx.NewYAML()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestYAMLReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestYAMLReverse(t *testing.T) {
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

// TestYAMLMarshalEmptyStruct 测试空结构体序列化
// Validates: Requirements 2.1
func TestYAMLMarshalEmptyStruct(t *testing.T) {
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

// TestYAMLMarshalSpecialCharacters 测试包含特殊字符的结构体序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalSpecialCharacters(t *testing.T) {
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

// TestYAMLMarshalUnicodeString 测试包含 Unicode 字符的结构体序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalUnicodeString(t *testing.T) {
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

// TestYAMLMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalLargeNumbers(t *testing.T) {
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

// TestYAMLImplementsEncoding 测试 YAML 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestYAMLImplementsEncoding(t *testing.T) {
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

// TestYAMLMarshalSlice 测试切片序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalSlice(t *testing.T) {
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

// TestYAMLMarshalMap 测试 map 序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalMap(t *testing.T) {
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

// TestYAMLUnmarshalInvalidYAML 测试无效 YAML 反序列化
// Validates: Requirements 2.2
func TestYAMLUnmarshalInvalidYAML(t *testing.T) {
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

// TestYAMLMarshalNegativeNumbers 测试负数序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalNegativeNumbers(t *testing.T) {
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

// TestYAMLMarshalZeroValues 测试零值序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalZeroValues(t *testing.T) {
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

// TestYAMLMarshalDeeplyNestedStruct 测试深度嵌套结构体
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalDeeplyNestedStruct(t *testing.T) {
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

// TestYAMLMarshalSliceOfMaps 测试 map 切片序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalSliceOfMaps(t *testing.T) {
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

// TestYAMLMarshalPointerFields 测试包含指针字段的结构体
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalPointerFields(t *testing.T) {
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

// TestYAMLMarshalNilPointerFields 测试包含 nil 指针字段的结构体
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalNilPointerFields(t *testing.T) {
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

// TestYAMLMarshalEmptySlice 测试空切片序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalEmptySlice(t *testing.T) {
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

// TestYAMLMarshalEmptyMap 测试空 map 序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalEmptyMap(t *testing.T) {
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

// TestYAMLMarshalMultilineString 测试多行字符串序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalMultilineString(t *testing.T) {
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

// TestYAMLMarshalYAMLSpecialChars 测试 YAML 特殊字符序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalYAMLSpecialChars(t *testing.T) {
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

// TestYAMLMarshalQuotedString 测试带引号的字符串序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalQuotedString(t *testing.T) {
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

// TestYAMLMarshalFloatSpecialValues 测试浮点数特殊值序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalFloatSpecialValues(t *testing.T) {
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

// TestYAMLMarshalIntegerBoundaries 测试整数边界值序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalIntegerBoundaries(t *testing.T) {
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

// TestYAMLMarshalBoolValues 测试布尔值序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalBoolValues(t *testing.T) {
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

// TestYAMLNewYAMLConstructor 测试 NewYAML 构造函数
// Validates: Requirements 14.1, 14.2, 14.3
func TestYAMLNewYAMLConstructor(t *testing.T) {
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

// TestYAMLMarshalInterfaceSlice 测试 interface{} 切片序列化
// Validates: Requirements 2.1, 2.2
func TestYAMLMarshalInterfaceSlice(t *testing.T) {
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

// ============================================================================
// YAML 编码器属性测试
// ============================================================================

// TestGroup1_Property_3_YAML_RoundTrip 属性测试：YAML Round-Trip 一致性
// **Property 3: YAML Round-Trip 一致性**
// *For any* 有效的 Go 结构体，使用 YAML 编码器序列化后再反序列化，
// 应该产生与原始结构体等价的对象。
// **Validates: Requirements 2.3**
func TestGroup1_Property_3_YAML_RoundTrip(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewYAML()

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
				t.Fatalf("Unmarshal failed: %v", err)
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

	t.Run("EmptyStruct_RoundTrip", func(t *testing.T) {
		// 测试空结构体的 Round-Trip
		for i := 0; i < iterations; i++ {
			original := TestStruct{}

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
				t.Errorf("Iteration %d: Empty struct round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("EmptySlice_RoundTrip", func(t *testing.T) {
		// 测试空切片的 Round-Trip
		for i := 0; i < iterations; i++ {
			original := []TestStruct{}

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

			// 验证 Round-Trip 一致性（空切片）
			if len(result) != 0 {
				t.Errorf("Iteration %d: Empty slice round-trip failed: expected empty slice, got %v", i, result)
			}
		}
	})

	t.Run("EmptyMap_RoundTrip", func(t *testing.T) {
		// 测试空 map 的 Round-Trip
		for i := 0; i < iterations; i++ {
			original := map[string]int{}

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

			// 验证 Round-Trip 一致性（空 map）
			if len(result) != 0 {
				t.Errorf("Iteration %d: Empty map round-trip failed: expected empty map, got %v", i, result)
			}
		}
	})
}
