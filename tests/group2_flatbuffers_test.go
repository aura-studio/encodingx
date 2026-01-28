package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	flatbuffers "github.com/google/flatbuffers/go"
	"pgregory.net/rapid"
)

// ============================================================================
// FlatBuffers 测试类型
// FlatBuffers 需要类型实现 FlatBufferMarshaler 和 FlatBufferUnmarshaler 接口
// ============================================================================

// FlatBufferTestStruct 是实现 FlatBuffers 接口的测试结构体
type FlatBufferTestStruct struct {
	Integer int32
	String  string
	Bool    bool
	Float   float64
}

// MarshalFlatBuffer 实现 FlatBufferMarshaler 接口
func (f *FlatBufferTestStruct) MarshalFlatBuffer(builder *flatbuffers.Builder) error {
	// 创建字符串
	strOffset := builder.CreateString(f.String)

	// 开始构建对象
	builder.StartObject(4)
	builder.PrependInt32Slot(0, f.Integer, 0)
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(strOffset), 0)
	boolVal := byte(0)
	if f.Bool {
		boolVal = 1
	}
	builder.PrependByteSlot(2, boolVal, 0)
	builder.PrependFloat64Slot(3, f.Float, 0)
	obj := builder.EndObject()
	builder.Finish(obj)

	return nil
}

// UnmarshalFlatBuffer 实现 FlatBufferUnmarshaler 接口
func (f *FlatBufferTestStruct) UnmarshalFlatBuffer(data []byte) error {
	// 获取根对象
	n := flatbuffers.GetUOffsetT(data)
	tab := &flatbuffers.Table{}
	tab.Bytes = data
	tab.Pos = n

	// 读取字段
	// Field 0: Integer (int32)
	o := flatbuffers.UOffsetT(tab.Offset(4))
	if o != 0 {
		f.Integer = tab.GetInt32(o + tab.Pos)
	}

	// Field 1: String
	o = flatbuffers.UOffsetT(tab.Offset(6))
	if o != 0 {
		f.String = string(tab.ByteVector(o + tab.Pos))
	}

	// Field 2: Bool
	o = flatbuffers.UOffsetT(tab.Offset(8))
	if o != 0 {
		f.Bool = tab.GetByte(o+tab.Pos) != 0
	}

	// Field 3: Float
	o = flatbuffers.UOffsetT(tab.Offset(10))
	if o != 0 {
		f.Float = tab.GetFloat64(o + tab.Pos)
	}

	return nil
}

// Equal 检查两个 FlatBufferTestStruct 是否相等
func (f FlatBufferTestStruct) Equal(other FlatBufferTestStruct) bool {
	return f.Integer == other.Integer &&
		f.String == other.String &&
		f.Bool == other.Bool &&
		f.Float == other.Float
}

// FlatBufferSimpleStruct 是一个简单的 FlatBuffers 测试结构体
type FlatBufferSimpleStruct struct {
	Value int32
}

// MarshalFlatBuffer 实现 FlatBufferMarshaler 接口
func (f *FlatBufferSimpleStruct) MarshalFlatBuffer(builder *flatbuffers.Builder) error {
	builder.StartObject(1)
	builder.PrependInt32Slot(0, f.Value, 0)
	obj := builder.EndObject()
	builder.Finish(obj)
	return nil
}

// UnmarshalFlatBuffer 实现 FlatBufferUnmarshaler 接口
func (f *FlatBufferSimpleStruct) UnmarshalFlatBuffer(data []byte) error {
	n := flatbuffers.GetUOffsetT(data)
	tab := &flatbuffers.Table{}
	tab.Bytes = data
	tab.Pos = n

	o := flatbuffers.UOffsetT(tab.Offset(4))
	if o != 0 {
		f.Value = tab.GetInt32(o + tab.Pos)
	}
	return nil
}

// Equal 检查两个 FlatBufferSimpleStruct 是否相等
func (f FlatBufferSimpleStruct) Equal(other FlatBufferSimpleStruct) bool {
	return f.Value == other.Value
}

// NonFlatBufferStruct 是一个不实现 FlatBuffers 接口的结构体
// 用于测试错误处理
type NonFlatBufferStruct struct {
	Value int
}

// ============================================================================
// FlatBuffers 编码器单元测试
// Validates: Requirements 17.1, 17.2, 17.3, 17.4
// ============================================================================

// TestFlatBuffersMarshalStruct 测试支持类型的序列化
// Validates: Requirements 17.1
func TestFlatBuffersMarshalStruct(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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

	// 验证返回的是有效的 FlatBuffers 数据（非空）
	if len(data) == 0 {
		t.Error("Marshal should return non-empty data")
	}
}

// TestFlatBuffersUnmarshalStruct 测试支持类型的反序列化
// Validates: Requirements 17.2
func TestFlatBuffersUnmarshalStruct(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
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

// TestFlatBuffersRoundTripStruct 测试结构体序列化后反序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersRoundTripStruct(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(*result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestFlatBuffersSimpleStruct 测试简单结构体序列化/反序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersSimpleStruct(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferSimpleStruct{Value: 12345}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := &FlatBufferSimpleStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(*result) {
		t.Errorf("Simple struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestFlatBuffersString 测试 String() 方法返回类型名称
// Validates: Requirements 17.3
func TestFlatBuffersString(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	name := encoder.String()

	if name != "FlatBuffers" {
		t.Errorf("String() should return 'FlatBuffers', got '%s'", name)
	}
}

// TestFlatBuffersStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 17.3
func TestFlatBuffersStyle(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestFlatBuffersReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 17.3
func TestFlatBuffersReverse(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 FlatBuffers 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestFlatBuffersImplementsEncoding 测试 FlatBuffers 编码器实现 Encoding 接口
// Validates: Requirements 17.3
func TestFlatBuffersImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewFlatBuffers()

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

// TestFlatBuffersRegisteredInEncodingSet 测试 FlatBuffers 编码器注册到 EncodingSet
// 通过 ChainEncoding 间接测试，因为 localEncoding 是内部函数
// Validates: Requirements 17.4
func TestFlatBuffersRegisteredInEncodingSet(t *testing.T) {
	// 创建使用 FlatBuffers 编码器的 ChainEncoding
	chain := encodingx.NewChainEncoding([]string{"FlatBuffers"}, []string{"FlatBuffers"})

	// 测试 Marshal - 如果 localEncoding 找不到 FlatBuffers，会返回错误
	input := &FlatBufferSimpleStruct{Value: 42}
	data, err := chain.Marshal(input)
	if err != nil {
		t.Fatalf("ChainEncoding with FlatBuffers should work, FlatBuffers not registered: %v", err)
	}

	// 测试 Unmarshal
	result := &FlatBufferSimpleStruct{}
	err = chain.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("ChainEncoding Unmarshal with FlatBuffers should work: %v", err)
	}

	// 验证数据正确
	if !input.Equal(*result) {
		t.Errorf("ChainEncoding round trip failed: expected %+v, got %+v", input, result)
	}
}

// TestFlatBuffersMarshalWrongType 测试非 FlatBufferMarshaler 类型返回错误
// Validates: Requirements 17.1
func TestFlatBuffersMarshalWrongType(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	wrongType := &NonFlatBufferStruct{Value: 42}

	// 序列化应该返回错误
	_, err := encoder.Marshal(wrongType)
	if err == nil {
		t.Error("Marshal should fail for non-FlatBufferMarshaler type")
	}
	if err != encodingx.ErrFlatBuffersWrongValueType {
		t.Errorf("Marshal should return ErrFlatBuffersWrongValueType, got: %v", err)
	}
}

// TestFlatBuffersUnmarshalWrongType 测试非 FlatBufferUnmarshaler 类型返回错误
// Validates: Requirements 17.2
func TestFlatBuffersUnmarshalWrongType(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()

	// 先用正确的类型序列化
	original := &FlatBufferSimpleStruct{Value: 42}
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化到错误类型应该返回错误
	wrongType := &NonFlatBufferStruct{}
	err = encoder.Unmarshal(data, wrongType)
	if err == nil {
		t.Error("Unmarshal should fail for non-FlatBufferUnmarshaler type")
	}
	if err != encodingx.ErrFlatBuffersWrongValueType {
		t.Errorf("Unmarshal should return ErrFlatBuffersWrongValueType, got: %v", err)
	}
}

// TestFlatBuffersMarshalEmptyStruct 测试空结构体序列化
// Validates: Requirements 17.1
func TestFlatBuffersMarshalEmptyStruct(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(*result) {
		t.Errorf("Empty struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestFlatBuffersMarshalSpecialCharacters 测试包含特殊字符的结构体序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersMarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(*result) {
		t.Errorf("Special characters round trip failed: original %+v != result %+v", original, result)
	}
}

// TestFlatBuffersMarshalUnicodeString 测试包含 Unicode 字符的结构体序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersMarshalUnicodeString(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(*result) {
		t.Errorf("Unicode string round trip failed: original %+v != result %+v", original, result)
	}
}

// TestFlatBuffersMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersMarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
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

// TestFlatBuffersMarshalNegativeNumbers 测试负数序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersMarshalNegativeNumbers(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
		Integer: -2147483648, // Min int32
		String:  "negative numbers",
		Bool:    false,
		Float:   -1.7976931348623157e+308, // Min float64
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if original.Integer != result.Integer {
		t.Errorf("Negative integer mismatch: expected %d, got %d", original.Integer, result.Integer)
	}
	if original.Float != result.Float {
		t.Errorf("Negative float mismatch: expected %e, got %e", original.Float, result.Float)
	}
}

// TestFlatBuffersMarshalZeroValues 测试零值序列化
// Validates: Requirements 17.1, 17.2
func TestFlatBuffersMarshalZeroValues(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()
	original := &FlatBufferTestStruct{
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
	result := &FlatBufferTestStruct{}
	err = encoder.Unmarshal(data, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(*result) {
		t.Errorf("Zero values round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// FlatBuffers 编码器属性测试
// ============================================================================

// genFlatBufferTestStruct 生成随机 FlatBufferTestStruct 的 rapid 生成器
func genFlatBufferTestStruct() *rapid.Generator[*FlatBufferTestStruct] {
	return rapid.Custom(func(t *rapid.T) *FlatBufferTestStruct {
		return &FlatBufferTestStruct{
			Integer: rapid.Int32Range(-10000, 10000).Draw(t, "integer"),
			String:  rapid.StringN(0, 100, -1).Draw(t, "string"),
			Bool:    rapid.Bool().Draw(t, "bool"),
			Float:   rapid.Float64().Draw(t, "float"),
		}
	})
}

// genFlatBufferSimpleStruct 生成随机 FlatBufferSimpleStruct 的 rapid 生成器
func genFlatBufferSimpleStruct() *rapid.Generator[*FlatBufferSimpleStruct] {
	return rapid.Custom(func(t *rapid.T) *FlatBufferSimpleStruct {
		return &FlatBufferSimpleStruct{
			Value: rapid.Int32().Draw(t, "value"),
		}
	})
}

// TestProperty24_FlatBuffersRoundTripConsistency 属性测试：FlatBuffers Round-Trip 一致性
// **Property 24: FlatBuffers Round-Trip 一致性**
// *For any* 有效的 FlatBuffers 类型，使用 FlatBuffers 编码器序列化后再反序列化，
// 应该产生与原始数据等价的对象。
// **Validates: Requirements 17.5**
func TestProperty24_FlatBuffersRoundTripConsistency(t *testing.T) {
	encoder := encodingx.NewFlatBuffers()

	t.Run("FlatBufferTestStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFlatBufferTestStruct().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			result := &FlatBufferTestStruct{}
			err = encoder.Unmarshal(data, result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(*result) {
				t.Fatalf("Round-trip failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})

	t.Run("FlatBufferSimpleStruct_RoundTrip", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			original := genFlatBufferSimpleStruct().Draw(t, "original")

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			result := &FlatBufferSimpleStruct{}
			err = encoder.Unmarshal(data, result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if !original.Equal(*result) {
				t.Fatalf("Round-trip failed:\n  original: %+v\n  result:   %+v", original, result)
			}
		})
	})
}
