package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// MsgPack 编码器单元测试
// Validates: Requirements 15.1, 15.2, 15.3, 15.4
// ============================================================================

// TestMsgPackMarshalStruct 测试普通结构体序列化
// Validates: Requirements 15.1
func TestMsgPackMarshalStruct(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

	// 验证返回的是有效的 MsgPack 数据（非空）
	if len(data) == 0 {
		t.Error("Marshal should return non-empty data")
	}
}

// TestMsgPackUnmarshalStruct 测试普通结构体反序列化
// Validates: Requirements 15.2
func TestMsgPackUnmarshalStruct(t *testing.T) {
	encoder := encodingx.NewMsgPack()
	original := TestStruct{
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
	var result TestStruct
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

// TestMsgPackRoundTripStruct 测试结构体序列化后反序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackRoundTripStruct(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackMarshalNestedStruct 测试嵌套结构体序列化/反序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalNestedStruct(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackString 测试 String() 方法返回类型名称
// Validates: Requirements 15.3
func TestMsgPackString(t *testing.T) {
	encoder := encodingx.NewMsgPack()
	name := encoder.String()

	if name != "MsgPack" {
		t.Errorf("String() should return 'MsgPack', got '%s'", name)
	}
}

// TestMsgPackStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 15.3
func TestMsgPackStyle(t *testing.T) {
	encoder := encodingx.NewMsgPack()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestMsgPackReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 15.3
func TestMsgPackReverse(t *testing.T) {
	encoder := encodingx.NewMsgPack()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 MsgPack 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestMsgPackImplementsEncoding 测试 MsgPack 编码器实现 Encoding 接口
// Validates: Requirements 15.3
func TestMsgPackImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewMsgPack()

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

// TestMsgPackRegisteredInEncodingSet 测试 MsgPack 编码器注册到 EncodingSet
// 通过 ChainEncoding 间接测试，因为 localEncoding 是内部函数
// Validates: Requirements 15.4
func TestMsgPackRegisteredInEncodingSet(t *testing.T) {
	// 创建使用 MsgPack 编码器的 ChainEncoding
	chain := encodingx.NewChainEncoding([]string{"MsgPack"}, []string{"MsgPack"})

	// 测试 Marshal - 如果 localEncoding 找不到 MsgPack，会返回错误
	input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
	data, err := chain.Marshal(input)
	if err != nil {
		t.Fatalf("ChainEncoding with MsgPack should work, MsgPack not registered: %v", err)
	}

	// 测试 Unmarshal
	var result TestStruct
	err = chain.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("ChainEncoding Unmarshal with MsgPack should work: %v", err)
	}

	// 验证数据正确
	if !input.Equal(result) {
		t.Errorf("ChainEncoding round trip failed: expected %+v, got %+v", input, result)
	}
}

// TestMsgPackMarshalEmptyStruct 测试空结构体序列化
// Validates: Requirements 15.1
func TestMsgPackMarshalEmptyStruct(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackMarshalSpecialCharacters 测试包含特殊字符的结构体序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackMarshalUnicodeString 测试包含 Unicode 字符的结构体序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalUnicodeString(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackMarshalSlice 测试切片序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalSlice(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackMarshalMap 测试 map 序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalMap(t *testing.T) {
	encoder := encodingx.NewMsgPack()
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

// TestMsgPackUnmarshalInvalidData 测试无效数据反序列化
// Validates: Requirements 15.2
func TestMsgPackUnmarshalInvalidData(t *testing.T) {
	encoder := encodingx.NewMsgPack()
	// 无效的 MsgPack 数据
	invalidData := []byte{0xFF, 0xFF, 0xFF, 0xFF}

	var result TestStruct
	err := encoder.Unmarshal(invalidData, &result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid MsgPack data")
	}
}

// TestMsgPackMarshalNilValue 测试 nil 值序列化
// Validates: Requirements 15.1
func TestMsgPackMarshalNilValue(t *testing.T) {
	encoder := encodingx.NewMsgPack()

	// 序列化 nil
	data, err := encoder.Marshal(nil)
	if err != nil {
		t.Fatalf("Marshal nil failed: %v", err)
	}

	// 验证返回的数据不为空（MsgPack 会编码 nil 为特定字节）
	if len(data) == 0 {
		t.Error("Marshal nil should return non-empty data")
	}
}

// TestMsgPackMarshalPointer 测试指针序列化
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalPointer(t *testing.T) {
	encoder := encodingx.NewMsgPack()
	original := &TestStruct{
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
	var result TestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Pointer round trip failed: original %+v != result %+v", *original, result)
	}
}

// TestMsgPackMarshalDeeplyNestedStruct 测试深度嵌套结构体
// Validates: Requirements 15.1, 15.2
func TestMsgPackMarshalDeeplyNestedStruct(t *testing.T) {
	encoder := encodingx.NewMsgPack()

	// 创建深度嵌套的结构
	type Level3 struct {
		Value string `msgpack:"value"`
	}
	type Level2 struct {
		Level3 Level3 `msgpack:"level3"`
	}
	type Level1 struct {
		Level2 Level2 `msgpack:"level2"`
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

// ============================================================================
// MsgPack 编码器属性测试
// ============================================================================

// genMsgPackTestStruct 生成随机 TestStruct 的 rapid 生成器
func genMsgPackTestStruct() *rapid.Generator[TestStruct] {
	return rapid.Custom(func(t *rapid.T) TestStruct {
		return TestStruct{
			Integer: rapid.IntRange(-10000, 10000).Draw(t, "integer"),
			String:  rapid.StringN(0, 100, -1).Draw(t, "string"),
			Bool:    rapid.Bool().Draw(t, "bool"),
			Float:   rapid.Float64().Draw(t, "float"),
		}
	})
}

// genMsgPackNestedStruct 生成随机 NestedStruct 的 rapid 生成器
func genMsgPackNestedStruct() *rapid.Generator[NestedStruct] {
	return rapid.Custom(func(t *rapid.T) NestedStruct {
		sliceLen := rapid.IntRange(0, 20).Draw(t, "sliceLen")
		slice := make([]int, sliceLen)
		for i := 0; i < sliceLen; i++ {
			slice[i] = rapid.IntRange(-1000, 1000).Draw(t, "sliceElem")
		}
		return NestedStruct{
			Name:  rapid.StringN(0, 50, -1).Draw(t, "name"),
			Inner: genMsgPackTestStruct().Draw(t, "inner"),
			Slice: slice,
		}
	})
}

// genMsgPackStringIntMap 生成随机 map[string]int 的 rapid 生成器
func genMsgPackStringIntMap() *rapid.Generator[map[string]int] {
	return rapid.Custom(func(t *rapid.T) map[string]int {
		numKeys := rapid.IntRange(0, 20).Draw(t, "numKeys")
		result := make(map[string]int)
		for i := 0; i < numKeys; i++ {
			key := rapid.StringN(1, 20, -1).Draw(t, "key")
			value := rapid.IntRange(-10000, 10000).Draw(t, "value")
			result[key] = value
		}
		return result
	})
}

// TestProperty22_MsgPackRoundTripConsistency 属性测试：MsgPack Round-Trip 一致性
// **Property 22: MsgPack Round-Trip 一致性**
// *For any* 有效的 Go 结构体，使用 MsgPack 编码器序列化后再反序列化，
// 应该产生与原始结构体等价的对象。
// **Validates: Requirements 15.5**
func TestProperty22_MsgPackRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewMsgPack()

	t.Run("TestStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genMsgPackTestStruct().Draw(t, "original")

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

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Fatalf("Round-trip failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("NestedStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genMsgPackNestedStruct().Draw(t, "original")

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

			// 验证 Round-Trip 一致性
			if !original.Equal(result) {
				t.Fatalf("Round-trip failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("Map_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genMsgPackStringIntMap().Draw(t, "original")

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

	t.Run("Slice_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			// 生成随机切片
			sliceLen := rapid.IntRange(0, 20).Draw(t, "sliceLen")
			original := make([]TestStruct, sliceLen)
			for i := 0; i < sliceLen; i++ {
				original[i] = genMsgPackTestStruct().Draw(t, "elem")
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

			// 验证 Round-Trip 一致性
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
			// 测试各种基本类型
			intVal := rapid.Int().Draw(t, "int")
			stringVal := rapid.String().Draw(t, "string")
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
}
