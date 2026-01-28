package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
	"github.com/aura-studio/encodingx/tests/testdata"
)

// ============================================================================
// Protobuf 编码器单元测试
// Validates: Requirements 5.1, 5.2, 5.3, 5.4, 14.1, 14.2, 14.3
// ============================================================================

// ============================================================================
// 基本序列化/反序列化测试
// ============================================================================

// TestProtobufMarshalProtoMessage 测试 proto.Message 序列化
// Validates: Requirements 5.1
func TestProtobufMarshalProtoMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建测试消息
	original := testdata.NewTestMessageWithValues(42, "hello world", true)

	// 序列化
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是有效的 Protobuf 字节数组
	if len(data) == 0 {
		t.Error("Marshal should return non-empty byte array")
	}

	// 验证可以反序列化
	result := testdata.NewTestMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证字段值
	if result.GetIntField() != original.GetIntField() {
		t.Errorf("IntField mismatch: expected %d, got %d", original.GetIntField(), result.GetIntField())
	}
	if result.GetStringField() != original.GetStringField() {
		t.Errorf("StringField mismatch: expected %s, got %s", original.GetStringField(), result.GetStringField())
	}
	if result.GetBoolField() != original.GetBoolField() {
		t.Errorf("BoolField mismatch: expected %v, got %v", original.GetBoolField(), result.GetBoolField())
	}
}

// TestProtobufUnmarshalProtoMessage 测试 proto.Message 反序列化
// Validates: Requirements 5.2
func TestProtobufUnmarshalProtoMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建并序列化测试消息
	original := testdata.NewTestMessageWithValues(100, "test string", false)
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := testdata.NewTestMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证字段值
	if result.GetIntField() != 100 {
		t.Errorf("IntField mismatch: expected 100, got %d", result.GetIntField())
	}
	if result.GetStringField() != "test string" {
		t.Errorf("StringField mismatch: expected 'test string', got '%s'", result.GetStringField())
	}
	if result.GetBoolField() != false {
		t.Errorf("BoolField mismatch: expected false, got %v", result.GetBoolField())
	}
}

// TestProtobufRoundTripProtoMessage 测试 proto.Message 序列化后反序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufRoundTripProtoMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	testCases := []struct {
		name        string
		intField    int32
		stringField string
		boolField   bool
	}{
		{"basic values", 42, "hello", true},
		{"zero values", 0, "", false},
		{"negative int", -100, "negative", true},
		{"large int", 2147483647, "max int32", false},
		{"unicode string", 123, "你好世界", true},
		{"special chars", 456, "hello\nworld\ttab", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := testdata.NewTestMessageWithValues(tc.intField, tc.stringField, tc.boolField)

			// 序列化
			data, err := encoder.Marshal(original.Message)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			result := testdata.NewTestMessage()
			err = encoder.Unmarshal(data, result.Message)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证 Round-Trip 一致性
			if result.GetIntField() != original.GetIntField() {
				t.Errorf("IntField mismatch: expected %d, got %d", original.GetIntField(), result.GetIntField())
			}
			if result.GetStringField() != original.GetStringField() {
				t.Errorf("StringField mismatch: expected %s, got %s", original.GetStringField(), result.GetStringField())
			}
			if result.GetBoolField() != original.GetBoolField() {
				t.Errorf("BoolField mismatch: expected %v, got %v", original.GetBoolField(), result.GetBoolField())
			}
		})
	}
}

// ============================================================================
// 嵌套消息测试
// ============================================================================

// TestProtobufMarshalNestedMessage 测试嵌套消息序列化/反序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalNestedMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建嵌套消息
	inner := testdata.NewTestMessageWithValues(100, "inner message", true)
	original := testdata.NewNestedMessageWithValues("outer name", inner)

	// 序列化
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := testdata.NewNestedMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证外层字段
	if result.GetName() != original.GetName() {
		t.Errorf("Name mismatch: expected %s, got %s", original.GetName(), result.GetName())
	}

	// 验证内层消息
	resultInner := result.GetInner()
	if resultInner == nil {
		t.Fatal("Inner message should not be nil")
	}
	if resultInner.GetIntField() != inner.GetIntField() {
		t.Errorf("Inner.IntField mismatch: expected %d, got %d", inner.GetIntField(), resultInner.GetIntField())
	}
	if resultInner.GetStringField() != inner.GetStringField() {
		t.Errorf("Inner.StringField mismatch: expected %s, got %s", inner.GetStringField(), resultInner.GetStringField())
	}
	if resultInner.GetBoolField() != inner.GetBoolField() {
		t.Errorf("Inner.BoolField mismatch: expected %v, got %v", inner.GetBoolField(), resultInner.GetBoolField())
	}
}

// TestProtobufMarshalNestedMessageWithNilInner 测试嵌套消息（内层为空）序列化/反序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalNestedMessageWithNilInner(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建嵌套消息，内层为 nil
	original := testdata.NewNestedMessageWithValues("outer only", nil)

	// 序列化
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := testdata.NewNestedMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证外层字段
	if result.GetName() != original.GetName() {
		t.Errorf("Name mismatch: expected %s, got %s", original.GetName(), result.GetName())
	}
}

// ============================================================================
// 空消息测试
// ============================================================================

// TestProtobufMarshalEmptyMessage 测试空消息序列化/反序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalEmptyMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建空消息
	original := testdata.NewEmptyMessage()

	// 序列化
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 空消息序列化后应该是空字节数组或非常小的字节数组
	// Protobuf 空消息序列化后通常是 0 字节
	if len(data) > 10 {
		t.Errorf("Empty message should serialize to small byte array, got %d bytes", len(data))
	}

	// 反序列化
	result := testdata.NewEmptyMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
}

// TestProtobufMarshalDefaultValues 测试默认值消息序列化/反序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalDefaultValues(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建带默认值的消息
	original := testdata.NewTestMessageWithValues(0, "", false)

	// 序列化
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := testdata.NewTestMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证默认值
	if result.GetIntField() != 0 {
		t.Errorf("IntField should be 0, got %d", result.GetIntField())
	}
	if result.GetStringField() != "" {
		t.Errorf("StringField should be empty, got '%s'", result.GetStringField())
	}
	if result.GetBoolField() != false {
		t.Errorf("BoolField should be false, got %v", result.GetBoolField())
	}
}

// ============================================================================
// 错误处理测试 - 非 proto.Message 类型
// ============================================================================

// TestProtobufMarshalNonProtoMessage 测试序列化非 proto.Message 类型返回错误
// Validates: Requirements 5.3
func TestProtobufMarshalNonProtoMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	testCases := []struct {
		name  string
		value interface{}
	}{
		{"string", "hello"},
		{"int", 42},
		{"float64", 3.14},
		{"bool", true},
		{"byte slice", []byte{1, 2, 3}},
		{"struct", struct{ Name string }{"test"}},
		{"map", map[string]int{"key": 1}},
		{"slice", []int{1, 2, 3}},
		{"nil", nil},
		{"pointer to struct", &struct{ Name string }{"test"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := encoder.Marshal(tc.value)
			if err == nil {
				t.Errorf("Marshal(%v) should return error for non-proto.Message type", tc.value)
			}
			if err != encodingx.ErrProtobufWrongValueType {
				t.Errorf("Marshal(%v) should return ErrProtobufWrongValueType, got %v", tc.value, err)
			}
		})
	}
}

// TestProtobufUnmarshalNonProtoMessage 测试反序列化到非 proto.Message 类型返回错误
// Validates: Requirements 5.4
func TestProtobufUnmarshalNonProtoMessage(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 先创建有效的 Protobuf 数据
	msg := testdata.NewTestMessageWithValues(42, "test", true)
	data, err := encoder.Marshal(msg.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	testCases := []struct {
		name  string
		value interface{}
	}{
		{"string pointer", new(string)},
		{"int pointer", new(int)},
		{"float64 pointer", new(float64)},
		{"bool pointer", new(bool)},
		{"byte slice pointer", new([]byte)},
		{"struct pointer", &struct{ Name string }{}},
		{"map pointer", &map[string]int{}},
		{"slice pointer", &[]int{}},
		{"nil", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := encoder.Unmarshal(data, tc.value)
			if err == nil {
				t.Errorf("Unmarshal to %v should return error for non-proto.Message type", tc.value)
			}
			if err != encodingx.ErrProtobufWrongValueType {
				t.Errorf("Unmarshal to %v should return ErrProtobufWrongValueType, got %v", tc.value, err)
			}
		})
	}
}

// ============================================================================
// String()、Style()、Reverse() 方法测试
// ============================================================================

// TestProtobufString 测试 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestProtobufString(t *testing.T) {
	encoder := encodingx.NewProtobuf()
	name := encoder.String()

	if name != "Protobuf" {
		t.Errorf("String() should return 'Protobuf', got '%s'", name)
	}
}

// TestProtobufStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestProtobufStyle(t *testing.T) {
	encoder := encodingx.NewProtobuf()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestProtobufReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestProtobufReverse(t *testing.T) {
	encoder := encodingx.NewProtobuf()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 Protobuf 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestProtobufImplementsEncoding 测试 Protobuf 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestProtobufImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewProtobuf()

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
// 边界情况和特殊值测试
// ============================================================================

// TestProtobufMarshalLargeString 测试大字符串序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalLargeString(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 创建大字符串
	largeString := make([]byte, 10000)
	for i := range largeString {
		largeString[i] = byte('a' + (i % 26))
	}

	original := testdata.NewTestMessageWithValues(1, string(largeString), true)

	// 序列化
	data, err := encoder.Marshal(original.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	result := testdata.NewTestMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if result.GetStringField() != original.GetStringField() {
		t.Errorf("Large string mismatch: lengths %d vs %d", len(original.GetStringField()), len(result.GetStringField()))
	}
}

// TestProtobufMarshalExtremeIntValues 测试极端整数值序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalExtremeIntValues(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	testCases := []struct {
		name  string
		value int32
	}{
		{"max int32", 2147483647},
		{"min int32", -2147483648},
		{"zero", 0},
		{"positive one", 1},
		{"negative one", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := testdata.NewTestMessageWithValues(tc.value, "test", true)

			// 序列化
			data, err := encoder.Marshal(original.Message)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// 反序列化
			result := testdata.NewTestMessage()
			err = encoder.Unmarshal(data, result.Message)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// 验证
			if result.GetIntField() != tc.value {
				t.Errorf("IntField mismatch: expected %d, got %d", tc.value, result.GetIntField())
			}
		})
	}
}

// TestProtobufUnmarshalEmptyData 测试反序列化空数据
// Validates: Requirements 5.2
func TestProtobufUnmarshalEmptyData(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	// 空数据应该可以反序列化到空消息
	result := testdata.NewTestMessage()
	err := encoder.Unmarshal([]byte{}, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal empty data failed: %v", err)
	}

	// 验证字段为默认值
	if result.GetIntField() != 0 {
		t.Errorf("IntField should be 0 for empty data, got %d", result.GetIntField())
	}
	if result.GetStringField() != "" {
		t.Errorf("StringField should be empty for empty data, got '%s'", result.GetStringField())
	}
	if result.GetBoolField() != false {
		t.Errorf("BoolField should be false for empty data, got %v", result.GetBoolField())
	}
}

// ============================================================================
// 属性测试 - Protobuf Round-Trip 一致性
// ============================================================================

// TestGroup1_Property_7_Protobuf_RoundTrip 属性测试：Protobuf Round-Trip 一致性
// **Property 7: Protobuf Round-Trip 一致性**
// *For any* 有效的 proto.Message，使用 Protobuf 编码器序列化后再反序列化，
// 应该产生与原始消息等价的对象。
// **Validates: Requirements 5.5**
func TestGroup1_Property_7_Protobuf_RoundTrip(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewProtobuf()
	gen := NewTestDataGenerator()

	t.Run("TestMessage_RoundTrip", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			// 生成随机值
			intField := gen.GenerateInt32()
			xmlStruct := gen.GenerateXMLTestStruct()
			stringField := xmlStruct.String
			boolField := gen.GenerateInt32()%2 == 0

			original := testdata.NewTestMessageWithValues(intField, stringField, boolField)

			// 序列化
			data, err := encoder.Marshal(original.Message)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			result := testdata.NewTestMessage()
			err = encoder.Unmarshal(data, result.Message)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if result.GetIntField() != original.GetIntField() {
				t.Errorf("Iteration %d: IntField mismatch: expected %d, got %d",
					i, original.GetIntField(), result.GetIntField())
			}
			if result.GetStringField() != original.GetStringField() {
				t.Errorf("Iteration %d: StringField mismatch: expected %s, got %s",
					i, original.GetStringField(), result.GetStringField())
			}
			if result.GetBoolField() != original.GetBoolField() {
				t.Errorf("Iteration %d: BoolField mismatch: expected %v, got %v",
					i, original.GetBoolField(), result.GetBoolField())
			}
		}
	})

	t.Run("NestedMessage_RoundTrip", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			// 生成随机嵌套消息
			innerIntField := gen.GenerateInt32()
			innerXmlStruct := gen.GenerateXMLTestStruct()
			innerStringField := innerXmlStruct.String
			innerBoolField := gen.GenerateInt32()%2 == 0
			outerXmlStruct := gen.GenerateXMLTestStruct()
			outerName := outerXmlStruct.String

			inner := testdata.NewTestMessageWithValues(innerIntField, innerStringField, innerBoolField)
			original := testdata.NewNestedMessageWithValues(outerName, inner)

			// 序列化
			data, err := encoder.Marshal(original.Message)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			result := testdata.NewNestedMessage()
			err = encoder.Unmarshal(data, result.Message)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if result.GetName() != original.GetName() {
				t.Errorf("Iteration %d: Name mismatch: expected %s, got %s",
					i, original.GetName(), result.GetName())
			}

			resultInner := result.GetInner()
			if resultInner == nil {
				t.Fatalf("Iteration %d: Inner message should not be nil", i)
			}
			if resultInner.GetIntField() != inner.GetIntField() {
				t.Errorf("Iteration %d: Inner.IntField mismatch: expected %d, got %d",
					i, inner.GetIntField(), resultInner.GetIntField())
			}
			if resultInner.GetStringField() != inner.GetStringField() {
				t.Errorf("Iteration %d: Inner.StringField mismatch: expected %s, got %s",
					i, inner.GetStringField(), resultInner.GetStringField())
			}
			if resultInner.GetBoolField() != inner.GetBoolField() {
				t.Errorf("Iteration %d: Inner.BoolField mismatch: expected %v, got %v",
					i, inner.GetBoolField(), resultInner.GetBoolField())
			}
		}
	})

	t.Run("EmptyMessage_RoundTrip", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			original := testdata.NewEmptyMessage()

			// 序列化
			data, err := encoder.Marshal(original.Message)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			result := testdata.NewEmptyMessage()
			err = encoder.Unmarshal(data, result.Message)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 空消息应该成功序列化和反序列化
		}
	})
}

// ============================================================================
// 多消息类型测试
// ============================================================================

// TestProtobufMarshalMultipleMessageTypes 测试多种消息类型的序列化/反序列化
// Validates: Requirements 5.1, 5.2
func TestProtobufMarshalMultipleMessageTypes(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	t.Run("TestMessage", func(t *testing.T) {
		original := testdata.NewTestMessageWithValues(42, "test", true)
		data, err := encoder.Marshal(original.Message)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		result := testdata.NewTestMessage()
		err = encoder.Unmarshal(data, result.Message)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if result.GetIntField() != original.GetIntField() {
			t.Errorf("IntField mismatch")
		}
	})

	t.Run("NestedMessage", func(t *testing.T) {
		inner := testdata.NewTestMessageWithValues(100, "inner", false)
		original := testdata.NewNestedMessageWithValues("outer", inner)
		data, err := encoder.Marshal(original.Message)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		result := testdata.NewNestedMessage()
		err = encoder.Unmarshal(data, result.Message)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if result.GetName() != original.GetName() {
			t.Errorf("Name mismatch")
		}
	})

	t.Run("EmptyMessage", func(t *testing.T) {
		original := testdata.NewEmptyMessage()
		data, err := encoder.Marshal(original.Message)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		result := testdata.NewEmptyMessage()
		err = encoder.Unmarshal(data, result.Message)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
	})
}

// TestProtobufNewProtobufFunction 测试 NewProtobuf() 函数
// Validates: Requirements 14.1, 14.2, 14.3
func TestProtobufNewProtobufFunction(t *testing.T) {
	encoder := encodingx.NewProtobuf()

	if encoder == nil {
		t.Fatal("NewProtobuf() should return non-nil encoder")
	}

	// 验证返回的编码器可以正常工作
	msg := testdata.NewTestMessageWithValues(1, "test", true)
	data, err := encoder.Marshal(msg.Message)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	result := testdata.NewTestMessage()
	err = encoder.Unmarshal(data, result.Message)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.GetIntField() != msg.GetIntField() {
		t.Errorf("IntField mismatch")
	}
}

// TestProtobufErrorType 测试错误类型
// Validates: Requirements 5.3, 5.4
func TestProtobufErrorType(t *testing.T) {
	// 验证 ErrProtobufWrongValueType 是一个有效的错误
	err := encodingx.ErrProtobufWrongValueType
	if err == nil {
		t.Fatal("ErrProtobufWrongValueType should not be nil")
	}

	// 验证错误消息
	errMsg := err.Error()
	if errMsg == "" {
		t.Error("ErrProtobufWrongValueType should have non-empty error message")
	}
}
