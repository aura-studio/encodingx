package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// TOML 编码器单元测试
// Validates: Requirements 16.1, 16.2, 16.3, 16.4
// ============================================================================

// TOMLTestStruct is a TOML-compatible test struct
// TOML requires exported fields with valid TOML types
type TOMLTestStruct struct {
	Integer int     `toml:"integer"`
	String  string  `toml:"string"`
	Bool    bool    `toml:"bool"`
	Float   float64 `toml:"float"`
}

// Equal checks if two TOMLTestStruct are equal
func (t TOMLTestStruct) Equal(other TOMLTestStruct) bool {
	return t.Integer == other.Integer &&
		t.String == other.String &&
		t.Bool == other.Bool &&
		t.Float == other.Float
}

// TOMLNestedStruct is a TOML-compatible nested struct
type TOMLNestedStruct struct {
	Name  string         `toml:"name"`
	Inner TOMLTestStruct `toml:"inner"`
	Slice []int          `toml:"slice"`
}

// Equal checks if two TOMLNestedStruct are equal
func (n TOMLNestedStruct) Equal(other TOMLNestedStruct) bool {
	if n.Name != other.Name {
		return false
	}
	if !n.Inner.Equal(other.Inner) {
		return false
	}
	if len(n.Slice) != len(other.Slice) {
		return false
	}
	for i := range n.Slice {
		if n.Slice[i] != other.Slice[i] {
			return false
		}
	}
	return true
}

// TestTOMLMarshalStruct 测试普通结构体序列化
// Validates: Requirements 16.1
func TestTOMLMarshalStruct(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{
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

	// 验证返回的是有效的 TOML 数据（非空）
	if len(data) == 0 {
		t.Error("Marshal should return non-empty data")
	}
}

// TestTOMLUnmarshalStruct 测试普通结构体反序列化
// Validates: Requirements 16.2
func TestTOMLUnmarshalStruct(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{
		Integer: 42,
		String:  "hello world",
		Bool:    true,
		Float:   3.14159,
	}

	// 先序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TOMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if result.Integer != original.Integer {
		t.Errorf("Integer mismatch: expected %d, got %d", original.Integer, result.Integer)
	}
	if result.String != original.String {
		t.Errorf("String mismatch: expected '%s', got '%s'", original.String, result.String)
	}
	if result.Bool != original.Bool {
		t.Errorf("Bool mismatch: expected %v, got %v", original.Bool, result.Bool)
	}
	if result.Float != original.Float {
		t.Errorf("Float mismatch: expected %f, got %f", original.Float, result.Float)
	}
}

// TestTOMLRoundTripStruct 测试结构体序列化后反序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLRoundTripStruct(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{
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
	var result TOMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestTOMLMarshalNestedStruct 测试嵌套结构体序列化/反序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalNestedStruct(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLNestedStruct{
		Name: "nested test",
		Inner: TOMLTestStruct{
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
	var result TOMLNestedStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Nested struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestTOMLString 测试 String() 方法返回类型名称
// Validates: Requirements 16.3
func TestTOMLString(t *testing.T) {
	encoder := encodingx.NewTOML()
	name := encoder.String()

	if name != "TOML" {
		t.Errorf("String() should return 'TOML', got '%s'", name)
	}
}

// TestTOMLStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 16.3
func TestTOMLStyle(t *testing.T) {
	encoder := encodingx.NewTOML()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestTOMLReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 16.3
func TestTOMLReverse(t *testing.T) {
	encoder := encodingx.NewTOML()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 TOML 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestTOMLImplementsEncoding 测试 TOML 编码器实现 Encoding 接口
// Validates: Requirements 16.3
func TestTOMLImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewTOML()

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

// TestTOMLRegisteredInEncodingSet 测试 TOML 编码器注册到 EncodingSet
// 通过 ChainEncoding 间接测试，因为 localEncoding 是内部函数
// Validates: Requirements 16.4
func TestTOMLRegisteredInEncodingSet(t *testing.T) {
	// 创建使用 TOML 编码器的 ChainEncoding
	chain := encodingx.NewChainEncoding([]string{"TOML"}, []string{"TOML"})

	// 测试 Marshal - 如果 localEncoding 找不到 TOML，会返回错误
	input := TOMLTestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
	data, err := chain.Marshal(input)
	if err != nil {
		t.Fatalf("ChainEncoding with TOML should work, TOML not registered: %v", err)
	}

	// 测试 Unmarshal
	var result TOMLTestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("ChainEncoding Unmarshal with TOML should work: %v", err)
	}

	// 验证数据正确
	if !input.Equal(result) {
		t.Errorf("ChainEncoding round trip failed: expected %+v, got %+v", input, result)
	}
}

// TestTOMLMarshalEmptyStruct 测试空结构体序列化
// Validates: Requirements 16.1
func TestTOMLMarshalEmptyStruct(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TOMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Empty struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestTOMLMarshalSpecialCharacters 测试包含特殊字符的结构体序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{
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
	var result TOMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Special characters round trip failed: original %+v != result %+v", original, result)
	}
}

// TestTOMLMarshalUnicodeString 测试包含 Unicode 字符的结构体序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalUnicodeString(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{
		Integer: 123,
		String:  "你好世界 مرحبا",
		Bool:    false,
		Float:   42.0,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TOMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Unicode string round trip failed: original %+v != result %+v", original, result)
	}
}

// TestTOMLMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := TOMLTestStruct{
		Integer: 2147483647, // Max int32
		String:  "large numbers",
		Bool:    true,
		Float:   1.7976931348623157e+100, // Large but not max float64 (TOML has limits)
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TOMLTestStruct
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

// TestTOMLMarshalMap 测试 map 序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalMap(t *testing.T) {
	encoder := encodingx.NewTOML()
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

// TestTOMLUnmarshalInvalidData 测试无效数据反序列化
// Validates: Requirements 16.2
func TestTOMLUnmarshalInvalidData(t *testing.T) {
	encoder := encodingx.NewTOML()
	// 无效的 TOML 数据
	invalidData := []byte("this is not valid toml [[[")

	var result TOMLTestStruct
	err := encoder.Unmarshal(invalidData, &result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid TOML data")
	}
}

// TestTOMLMarshalPointer 测试指针序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalPointer(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := &TOMLTestStruct{
		Integer: 42,
		String:  "pointer test",
		Bool:    true,
		Float:   3.14,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result TOMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Pointer round trip failed: original %+v != result %+v", *original, result)
	}
}

// TestTOMLMarshalDeeplyNestedStruct 测试深度嵌套结构体
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalDeeplyNestedStruct(t *testing.T) {
	encoder := encodingx.NewTOML()

	// 创建深度嵌套的结构
	type Level3 struct {
		Value string `toml:"value"`
	}
	type Level2 struct {
		Level3 Level3 `toml:"level3"`
	}
	type Level1 struct {
		Level2 Level2 `toml:"level2"`
	}

	original := Level1{
		Level2: Level2{
			Level3: Level3{
				Value: "deeply nested",
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

// TestTOMLMarshalStringMap 测试 map[string]string 序列化
// Validates: Requirements 16.1, 16.2
func TestTOMLMarshalStringMap(t *testing.T) {
	encoder := encodingx.NewTOML()
	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result map[string]string
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
			t.Errorf("Map value mismatch for key '%s': expected '%s', got '%s'", k, v, result[k])
		}
	}
}

// ============================================================================
// TOML 编码器属性测试
// ============================================================================

// genTOMLTestStruct 生成随机 TOMLTestStruct 的 rapid 生成器
// TOML has some limitations, so we use constrained generators
func genTOMLTestStruct() *rapid.Generator[TOMLTestStruct] {
	return rapid.Custom(func(t *rapid.T) TOMLTestStruct {
		return TOMLTestStruct{
			Integer: rapid.IntRange(-10000, 10000).Draw(t, "integer"),
			String:  rapid.StringN(0, 100, -1).Draw(t, "string"),
			Bool:    rapid.Bool().Draw(t, "bool"),
			// Use a constrained float range to avoid TOML serialization issues with special values
			Float: rapid.Float64Range(-1e100, 1e100).Draw(t, "float"),
		}
	})
}

// genTOMLNestedStruct 生成随机 TOMLNestedStruct 的 rapid 生成器
func genTOMLNestedStruct() *rapid.Generator[TOMLNestedStruct] {
	return rapid.Custom(func(t *rapid.T) TOMLNestedStruct {
		sliceLen := rapid.IntRange(0, 20).Draw(t, "sliceLen")
		slice := make([]int, sliceLen)
		for i := 0; i < sliceLen; i++ {
			slice[i] = rapid.IntRange(-1000, 1000).Draw(t, "sliceElem")
		}
		return TOMLNestedStruct{
			Name:  rapid.StringN(0, 50, -1).Draw(t, "name"),
			Inner: genTOMLTestStruct().Draw(t, "inner"),
			Slice: slice,
		}
	})
}

// genTOMLStringIntMap 生成随机 map[string]int 的 rapid 生成器
// TOML requires valid string keys
func genTOMLStringIntMap() *rapid.Generator[map[string]int] {
	return rapid.Custom(func(t *rapid.T) map[string]int {
		numKeys := rapid.IntRange(0, 20).Draw(t, "numKeys")
		result := make(map[string]int)
		for i := 0; i < numKeys; i++ {
			// Use alphanumeric keys to ensure TOML compatibility
			key := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9]{0,19}`).Draw(t, "key")
			value := rapid.IntRange(-10000, 10000).Draw(t, "value")
			result[key] = value
		}
		return result
	})
}

// TestProperty23_TOMLRoundTripConsistency 属性测试：TOML Round-Trip 一致性
// **Property 23: TOML Round-Trip 一致性**
// *For any* 有效的 Go 结构体（符合 TOML 规范），使用 TOML 编码器序列化后再反序列化，
// 应该产生与原始结构体等价的对象。
// **Validates: Requirements 16.5**
func TestProperty23_TOMLRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewTOML()

	t.Run("TOMLTestStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genTOMLTestStruct().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			var result TOMLTestStruct
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Fatalf("Round-trip failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("TOMLNestedStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genTOMLNestedStruct().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			var result TOMLNestedStruct
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Fatalf("Round-trip failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("Map_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genTOMLStringIntMap().Draw(t, "original")

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

			// 验证 Round-Trip 一致性
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

	t.Run("PrimitiveTypes_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// TOML requires top-level tables, so we wrap primitives in a struct
			type IntWrapper struct {
				Value int `toml:"value"`
			}
			type StringWrapper struct {
				Value string `toml:"value"`
			}
			type BoolWrapper struct {
				Value bool `toml:"value"`
			}
			type Float64Wrapper struct {
				Value float64 `toml:"value"`
			}

			// 测试 int
			intVal := IntWrapper{Value: rapid.Int().Draw(t, "int")}
			intData, err := encoder.Marshal(intVal)
			if err != nil {
				t.Fatalf("Marshal int failed: %v", err)
			}
			var intResult IntWrapper
			err = encoder.Unmarshal(intData, &intResult)
			if err != nil {
				t.Fatalf("Unmarshal int failed: %v", err)
			}
			if intVal.Value != intResult.Value {
				t.Fatalf("Int round-trip failed: expected %d, got %d", intVal.Value, intResult.Value)
			}

			// 测试 string
			stringVal := StringWrapper{Value: rapid.String().Draw(t, "string")}
			stringData, err := encoder.Marshal(stringVal)
			if err != nil {
				t.Fatalf("Marshal string failed: %v", err)
			}
			var stringResult StringWrapper
			err = encoder.Unmarshal(stringData, &stringResult)
			if err != nil {
				t.Fatalf("Unmarshal string failed: %v", err)
			}
			if stringVal.Value != stringResult.Value {
				t.Fatalf("String round-trip failed: expected '%s', got '%s'", stringVal.Value, stringResult.Value)
			}

			// 测试 bool
			boolVal := BoolWrapper{Value: rapid.Bool().Draw(t, "bool")}
			boolData, err := encoder.Marshal(boolVal)
			if err != nil {
				t.Fatalf("Marshal bool failed: %v", err)
			}
			var boolResult BoolWrapper
			err = encoder.Unmarshal(boolData, &boolResult)
			if err != nil {
				t.Fatalf("Unmarshal bool failed: %v", err)
			}
			if boolVal.Value != boolResult.Value {
				t.Fatalf("Bool round-trip failed: expected %v, got %v", boolVal.Value, boolResult.Value)
			}

			// 测试 float64 (constrained to avoid special values)
			float64Val := Float64Wrapper{Value: rapid.Float64Range(-1e100, 1e100).Draw(t, "float64")}
			float64Data, err := encoder.Marshal(float64Val)
			if err != nil {
				t.Fatalf("Marshal float64 failed: %v", err)
			}
			var float64Result Float64Wrapper
			err = encoder.Unmarshal(float64Data, &float64Result)
			if err != nil {
				t.Fatalf("Unmarshal float64 failed: %v", err)
			}
			if float64Val.Value != float64Result.Value {
				t.Fatalf("Float64 round-trip failed: expected %f, got %f", float64Val.Value, float64Result.Value)
			}
		})
	})
}
