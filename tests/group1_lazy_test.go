package encodingx_test

import (
"testing"

"github.com/aura-studio/encodingx"
"pgregory.net/rapid"
)

// ============================================================================
// Lazy 编码器单元测试
// Validates: Requirements 9.1, 9.2, 9.3, 9.4, 9.5, 14.1, 14.2, 14.3
// ============================================================================

// ============================================================================
// Lazy 编码器序列化测试 - []byte 类型
// ============================================================================

// TestLazyMarshalByteSlice 测试 []byte 序列化返回副本
// Validates: Requirements 9.1
func TestLazyMarshalByteSlice(t *testing.T) {
	encoder := encodingx.NewLazy()
	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的数据与原始数据相等
	if !BytesEqual(data, original) {
		t.Errorf("Marshal result mismatch: expected %v, got %v", original, data)
	}
}

// TestLazyMarshalByteSliceReturnsCopy 测试 []byte 序列化返回的是副本而非原始引用
// Validates: Requirements 9.1
func TestLazyMarshalByteSliceReturnsCopy(t *testing.T) {
	encoder := encodingx.NewLazy()
	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 修改返回的数据
	data[0] = 0xFF

	// 验证原始数据未被修改
	if !BytesEqual(original, originalCopy) {
		t.Errorf("Original data was modified: expected %v, got %v", originalCopy, original)
	}
}

// TestLazyMarshalByteSliceModifyOriginal 测试修改原始数据不影响序列化结果
// Validates: Requirements 9.1
func TestLazyMarshalByteSliceModifyOriginal(t *testing.T) {
	encoder := encodingx.NewLazy()
	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 保存序列化结果的副本
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	// 修改原始数据
	original[0] = 0xFF

	// 验证序列化结果未被修改
	if !BytesEqual(data, dataCopy) {
		t.Errorf("Serialized data was modified when original changed: expected %v, got %v", dataCopy, data)
	}
}

// TestLazyMarshalEmptyByteSlice 测试空 []byte 序列化
// Validates: Requirements 9.1
func TestLazyMarshalEmptyByteSlice(t *testing.T) {
	encoder := encodingx.NewLazy()
	original := []byte{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal empty []byte failed: %v", err)
	}

	// 验证返回空字节数组
	if len(data) != 0 {
		t.Errorf("Marshal empty []byte should return empty slice, got %v", data)
	}
}

// ============================================================================
// Lazy 编码器序列化测试 - Bytes 类型
// ============================================================================

// TestLazyMarshalBytes 测试 Bytes 类型序列化
// Validates: Requirements 9.2
func TestLazyMarshalBytes(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0x10, 0x20, 0x30, 0x40}
	original := encodingx.MakeBytes(originalData)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 Bytes.Data 的副本
	if !BytesEqual(data, originalData) {
		t.Errorf("Marshal result mismatch: expected %v, got %v", originalData, data)
	}
}

// TestLazyMarshalBytesReturnsCopy 测试 Bytes 类型序列化返回的是副本
// Validates: Requirements 9.2
func TestLazyMarshalBytesReturnsCopy(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0x10, 0x20, 0x30, 0x40}
	original := encodingx.MakeBytes(originalData)
	originalDataCopy := make([]byte, len(originalData))
	copy(originalDataCopy, originalData)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 修改返回的数据
	data[0] = 0xFF

	// 验证原始 Bytes.Data 未被修改
	if !BytesEqual(original.Data, originalDataCopy) {
		t.Errorf("Original Bytes.Data was modified: expected %v, got %v", originalDataCopy, original.Data)
	}
}

// TestLazyMarshalEmptyBytes 测试空 Bytes 类型序列化
// Validates: Requirements 9.2
func TestLazyMarshalEmptyBytes(t *testing.T) {
	encoder := encodingx.NewLazy()
	original := encodingx.MakeBytes([]byte{})

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal empty Bytes failed: %v", err)
	}

	// 验证返回空字节数组
	if len(data) != 0 {
		t.Errorf("Marshal empty Bytes should return empty slice, got %v", data)
	}
}

// ============================================================================
// Lazy 编码器序列化测试 - *Bytes 类型
// ============================================================================

// TestLazyMarshalBytesPointer 测试 *Bytes 类型序列化
// Validates: Requirements 9.3
func TestLazyMarshalBytesPointer(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	original := encodingx.NewBytes()
	original.Data = originalData

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是 Bytes.Data 的副本
	if !BytesEqual(data, originalData) {
		t.Errorf("Marshal result mismatch: expected %v, got %v", originalData, data)
	}
}

// TestLazyMarshalBytesPointerReturnsCopy 测试 *Bytes 类型序列化返回的是副本
// Validates: Requirements 9.3
func TestLazyMarshalBytesPointerReturnsCopy(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	original := encodingx.NewBytes()
	original.Data = originalData
	originalDataCopy := make([]byte, len(originalData))
	copy(originalDataCopy, originalData)

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 修改返回的数据
	data[0] = 0xFF

	// 验证原始 Bytes.Data 未被修改
	if !BytesEqual(original.Data, originalDataCopy) {
		t.Errorf("Original *Bytes.Data was modified: expected %v, got %v", originalDataCopy, original.Data)
	}
}

// TestLazyMarshalNilBytesPointer 测试 nil Data 的 *Bytes 序列化
// Validates: Requirements 9.3
func TestLazyMarshalNilBytesPointer(t *testing.T) {
	encoder := encodingx.NewLazy()
	original := encodingx.NewBytes() // Data 为 nil

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal nil *Bytes failed: %v", err)
	}

	// 验证返回空字节数组
	if len(data) != 0 {
		t.Errorf("Marshal nil *Bytes.Data should return empty slice, got %v", data)
	}
}

// ============================================================================
// Lazy 编码器反序列化测试 - *Bytes 类型
// ============================================================================

// TestLazyUnmarshalToBytes 测试反序列化到 *Bytes
// Validates: Requirements 9.4
func TestLazyUnmarshalToBytes(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0x11, 0x22, 0x33, 0x44, 0x55}

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(originalData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证数据被正确存入 Bytes.Data
	if !BytesEqual(result.Data, originalData) {
		t.Errorf("Unmarshal result mismatch: expected %v, got %v", originalData, result.Data)
	}
}

// TestLazyUnmarshalToBytesStoresCopy 测试反序列化存储的是数据副本
// Validates: Requirements 9.4
func TestLazyUnmarshalToBytesStoresCopy(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0x11, 0x22, 0x33, 0x44, 0x55}
	originalDataCopy := make([]byte, len(originalData))
	copy(originalDataCopy, originalData)

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(originalData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 修改原始数据
	originalData[0] = 0xFF

	// 验证 result.Data 未被修改
	if !BytesEqual(result.Data, originalDataCopy) {
		t.Errorf("Unmarshal result was modified when original changed: expected %v, got %v", originalDataCopy, result.Data)
	}
}

// TestLazyUnmarshalToBytesModifyResult 测试修改反序列化结果不影响原始数据
// Validates: Requirements 9.4
func TestLazyUnmarshalToBytesModifyResult(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0x11, 0x22, 0x33, 0x44, 0x55}
	originalDataCopy := make([]byte, len(originalData))
	copy(originalDataCopy, originalData)

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(originalData, result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 修改反序列化结果
	result.Data[0] = 0xFF

	// 验证原始数据未被修改
	if !BytesEqual(originalData, originalDataCopy) {
		t.Errorf("Original data was modified when result changed: expected %v, got %v", originalDataCopy, originalData)
	}
}

// TestLazyUnmarshalEmptyData 测试反序列化空数据
// Validates: Requirements 9.4
func TestLazyUnmarshalEmptyData(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{}

	result := encodingx.NewBytes()
	err := encoder.Unmarshal(originalData, result)
	if err != nil {
		t.Fatalf("Unmarshal empty data failed: %v", err)
	}

	// 验证结果为空
	if len(result.Data) != 0 {
		t.Errorf("Unmarshal empty data should result in empty Data, got %v", result.Data)
	}
}

// ============================================================================
// Lazy 编码器错误处理测试 - 非字节类型
// ============================================================================

// TestLazyMarshalWrongType 测试非字节类型返回错误
// Validates: Requirements 9.5
func TestLazyMarshalWrongType(t *testing.T) {
	encoder := encodingx.NewLazy()

	// 测试结构体类型
	_, err := encoder.Marshal(TestStruct{Integer: 42, String: "test"})
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for struct, got %v", err)
	}

	// 测试字符串类型
	_, err = encoder.Marshal("hello")
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for string, got %v", err)
	}

	// 测试整数类型
	_, err = encoder.Marshal(42)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for int, got %v", err)
	}

	// 测试浮点数类型
	_, err = encoder.Marshal(3.14)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for float, got %v", err)
	}

	// 测试布尔类型
	_, err = encoder.Marshal(true)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for bool, got %v", err)
	}

	// 测试 map 类型
	_, err = encoder.Marshal(map[string]int{"key": 1})
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for map, got %v", err)
	}

	// 测试切片类型（非 []byte）
	_, err = encoder.Marshal([]int{1, 2, 3})
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for []int, got %v", err)
	}

	// 测试 nil
	_, err = encoder.Marshal(nil)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for nil, got %v", err)
	}
}

// TestLazyUnmarshalWrongType 测试反序列化到非 *Bytes 类型返回错误
// Validates: Requirements 9.5
func TestLazyUnmarshalWrongType(t *testing.T) {
	encoder := encodingx.NewLazy()
	data := []byte{0x01, 0x02, 0x03}

	// 测试反序列化到结构体
	var ts TestStruct
	err := encoder.Unmarshal(data, &ts)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for struct pointer, got %v", err)
	}

	// 测试反序列化到字符串指针
	var s string
	err = encoder.Unmarshal(data, &s)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for string pointer, got %v", err)
	}

	// 测试反序列化到 []byte 指针（不是 *Bytes）
	var bs []byte
	err = encoder.Unmarshal(data, &bs)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for []byte pointer, got %v", err)
	}

	// 测试反序列化到整数指针
	var i int
	err = encoder.Unmarshal(data, &i)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for int pointer, got %v", err)
	}

	// 测试反序列化到 Bytes（非指针）
	var b encodingx.Bytes
	err = encoder.Unmarshal(data, b)
	if err != encodingx.ErrLazyWrongValueType {
		t.Errorf("Expected ErrLazyWrongValueType for Bytes (non-pointer), got %v", err)
	}
}

// ============================================================================
// Lazy 编码器 Round-Trip 测试
// ============================================================================

// TestLazyRoundTrip 测试 Lazy 编码器序列化/反序列化往返
// Validates: Requirements 9.1, 9.4
func TestLazyRoundTrip(t *testing.T) {
	encoder := encodingx.NewLazy()
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

// TestLazyRoundTripWithBytes 测试使用 Bytes 类型的往返
// Validates: Requirements 9.2, 9.4
func TestLazyRoundTripWithBytes(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0x12, 0x34, 0x56, 0x78}
	original := encodingx.MakeBytes(originalData)

	// 序列化
	encoded, err := encoder.Marshal(original)
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
		t.Errorf("Round trip with Bytes failed: expected %v, got %v", originalData, result.Data)
	}
}

// TestLazyRoundTripWithBytesPointer 测试使用 *Bytes 类型的往返
// Validates: Requirements 9.3, 9.4
func TestLazyRoundTripWithBytesPointer(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{0xAB, 0xCD, 0xEF}
	original := encodingx.NewBytes()
	original.Data = originalData

	// 序列化
	encoded, err := encoder.Marshal(original)
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
		t.Errorf("Round trip with *Bytes failed: expected %v, got %v", originalData, result.Data)
	}
}

// TestLazyRoundTripEmpty 测试空数据的往返
// Validates: Requirements 9.1, 9.4
func TestLazyRoundTripEmpty(t *testing.T) {
	encoder := encodingx.NewLazy()
	originalData := []byte{}

	// 序列化
	encoded, err := encoder.Marshal(originalData)
	if err != nil {
		t.Fatalf("Marshal empty data failed: %v", err)
	}

	// 反序列化
	result := encodingx.NewBytes()
	err = encoder.Unmarshal(encoded, result)
	if err != nil {
		t.Fatalf("Unmarshal empty data failed: %v", err)
	}

	// 验证往返一致性
	if len(result.Data) != 0 {
		t.Errorf("Round trip empty data failed: expected empty, got %v", result.Data)
	}
}

// ============================================================================
// Lazy 编码器接口方法测试
// ============================================================================

// TestLazyString 测试 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestLazyString(t *testing.T) {
	encoder := encodingx.NewLazy()
	name := encoder.String()

	if name != "Lazy" {
		t.Errorf("String() should return 'Lazy', got '%s'", name)
	}
}

// TestLazyStyle 测试 Style() 方法返回 EncodingStyleBytes
// Validates: Requirements 14.2
func TestLazyStyle(t *testing.T) {
	encoder := encodingx.NewLazy()
	style := encoder.Style()

	if style != encodingx.EncodingStyleBytes {
		t.Errorf("Style() should return EncodingStyleBytes, got %v", style)
	}
}

// TestLazyReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestLazyReverse(t *testing.T) {
	encoder := encodingx.NewLazy()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 Lazy 编码器
	if reversed.Style() != encodingx.EncodingStyleBytes {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestLazyImplementsEncoding 测试 Lazy 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestLazyImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewLazy()

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
// Lazy 编码器构造函数测试
// ============================================================================

// TestNewLazy 测试 NewLazy 构造函数
// Validates: Requirements 14.1, 14.2, 14.3
func TestNewLazy(t *testing.T) {
	encoder := encodingx.NewLazy()

	if encoder == nil {
		t.Fatal("NewLazy() should return non-nil encoder")
	}

	// 验证返回的是 Lazy 编码器
	if encoder.String() != "Lazy" {
		t.Errorf("NewLazy() should return Lazy encoder, got %s", encoder.String())
	}
}

// TestLazyMultipleInstances 测试多个 Lazy 实例的独立性
// Validates: Requirements 14.1, 14.2, 14.3
func TestLazyMultipleInstances(t *testing.T) {
	encoder1 := encodingx.NewLazy()
	encoder2 := encodingx.NewLazy()

	// 两个实例应该有相同的行为
	if encoder1.String() != encoder2.String() {
		t.Errorf("Multiple Lazy instances should have same String(): %s vs %s",
			encoder1.String(), encoder2.String())
	}

	if encoder1.Style() != encoder2.Style() {
		t.Errorf("Multiple Lazy instances should have same Style(): %v vs %v",
			encoder1.Style(), encoder2.Style())
	}

	// 测试两个实例的序列化结果相同
	testData := []byte{0x01, 0x02, 0x03}

	data1, err1 := encoder1.Marshal(testData)
	data2, err2 := encoder2.Marshal(testData)

	if err1 != nil || err2 != nil {
		t.Fatalf("Marshal failed: err1=%v, err2=%v", err1, err2)
	}

	if !BytesEqual(data1, data2) {
		t.Errorf("Multiple Lazy instances should produce same output: %v vs %v", data1, data2)
	}
}

// ============================================================================
// Lazy 编码器边界条件测试
// ============================================================================

// TestLazyLargeData 测试大数据编码/解码
// Validates: Requirements 9.1, 9.4
func TestLazyLargeData(t *testing.T) {
	encoder := encodingx.NewLazy()
	// 生成 1KB 的数据
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

// TestLazyAllByteValues 测试所有可能的字节值
// Validates: Requirements 9.1, 9.4
func TestLazyAllByteValues(t *testing.T) {
	encoder := encodingx.NewLazy()
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

// TestLazySingleByte 测试单字节编码/解码
// Validates: Requirements 9.1, 9.4
func TestLazySingleByte(t *testing.T) {
	encoder := encodingx.NewLazy()

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

// ============================================================================
// Lazy 编码器使用随机数据测试
// ============================================================================

// TestLazyRoundTripRandomData 测试使用随机数据的往返
// Validates: Requirements 9.1, 9.4
func TestLazyRoundTripRandomData(t *testing.T) {
	encoder := encodingx.NewLazy()
	gen := NewTestDataGenerator()

	// 运行多次随机测试
	for i := 0; i < 20; i++ {
		originalData := gen.GenerateBytes(0, 500)

		// 序列化
		encoded, err := encoder.Marshal(originalData)
		if err != nil {
			t.Fatalf("Marshal failed for iteration %d: %v", i, err)
		}

		// 反序列化
		result := encodingx.NewBytes()
		err = encoder.Unmarshal(encoded, result)
		if err != nil {
			t.Fatalf("Unmarshal failed for iteration %d: %v", i, err)
		}

		// 验证往返一致性
		if !BytesEqual(result.Data, originalData) {
			t.Errorf("Round trip failed for iteration %d: expected %v, got %v", i, originalData, result.Data)
		}
	}
}

// TestLazyDeepCopyVerification 测试深拷贝验证
// Validates: Requirements 9.1, 9.2, 9.3, 9.4
func TestLazyDeepCopyVerification(t *testing.T) {
	encoder := encodingx.NewLazy()

	t.Run("ByteSlice", func(t *testing.T) {
		original := []byte{0x01, 0x02, 0x03}
		encoded, _ := encoder.Marshal(original)

		// 修改原始数据
		original[0] = 0xFF
		// 验证编码结果未受影响
		if encoded[0] == 0xFF {
			t.Error("Marshal should return a copy, not a reference")
		}
	})

	t.Run("Bytes", func(t *testing.T) {
		original := encodingx.MakeBytes([]byte{0x01, 0x02, 0x03})
		encoded, _ := encoder.Marshal(original)

		// 修改原始数据
		original.Data[0] = 0xFF
		// 验证编码结果未受影响
		if encoded[0] == 0xFF {
			t.Error("Marshal Bytes should return a copy, not a reference")
		}
	})

	t.Run("BytesPointer", func(t *testing.T) {
		original := encodingx.NewBytes()
		original.Data = []byte{0x01, 0x02, 0x03}
		encoded, _ := encoder.Marshal(original)

		// 修改原始数据
		original.Data[0] = 0xFF
		// 验证编码结果未受影响
		if encoded[0] == 0xFF {
			t.Error("Marshal *Bytes should return a copy, not a reference")
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03}
		result := encodingx.NewBytes()
		encoder.Unmarshal(data, result)

		// 修改原始数据
		data[0] = 0xFF
		// 验证反序列化结果未受影响
		if result.Data[0] == 0xFF {
			t.Error("Unmarshal should store a copy, not a reference")
		}
	})
}

// ============================================================================
// Lazy 编码器属性测试
// ============================================================================

// TestProperty14_LazyRoundTripConsistency 测试 Lazy Round-Trip 一致性
// **Property 14: Lazy Round-Trip 一致性**
// *For any* 有效的字节数组，使用 Lazy 编码器序列化后再反序列化到 *Bytes，
// 应该产生与原始数据等价的字节。
// **Validates: Requirements 9.6**
func TestProperty14_LazyRoundTripConsistency(t *testing.T) {
encoder := encodingx.NewLazy()

rapid.Check(t, func(t *rapid.T) {
// 生成随机字节数据
original := genByteSlice().Draw(t, "original")

// 序列化
data, err := encoder.Marshal(original)
if err != nil {
t.Fatalf("Marshal failed: %v", err)
}

// 反序列化
result := encodingx.NewBytes()
err = encoder.Unmarshal(data, result)
if err != nil {
t.Fatalf("Unmarshal failed: %v", err)
}

// 验证 Round-Trip 一致性
if !BytesEqual(result.Data, original) {
t.Fatalf("Round-trip failed: original %v, got %v", original, result.Data)
}
})
}

// TestProperty15_LazyDeepCopy 测试 Lazy 深拷贝属性
// **Property 15: Lazy 深拷贝属性**
// *For any* 字节数组，Lazy 编码器的 Marshal 返回的数据应该是原始数据的独立副本，
// 修改返回值不应影响原始数据。
// **Validates: Requirements 9.1, 9.2, 9.3, 9.4**
func TestProperty15_LazyDeepCopy(t *testing.T) {
encoder := encodingx.NewLazy()

rapid.Check(t, func(t *rapid.T) {
// 生成随机字节数据（至少1字节以便修改）
length := rapid.IntRange(1, 1000).Draw(t, "length")
original := make([]byte, length)
for i := 0; i < length; i++ {
original[i] = byte(rapid.IntRange(0, 255).Draw(t, "byte"))
}

// 保存原始数据的副本用于比较
originalCopy := make([]byte, len(original))
copy(originalCopy, original)

// 序列化
data, err := encoder.Marshal(original)
if err != nil {
t.Fatalf("Marshal failed: %v", err)
}

// 修改返回的数据
if len(data) > 0 {
data[0] = ^data[0] // 翻转第一个字节
}

// 验证原始数据未被修改
if !BytesEqual(original, originalCopy) {
t.Fatalf("Deep copy failed: modifying Marshal result affected original data")
}

// 测试 Bytes 类型的深拷贝
bytesValue := encodingx.MakeBytes(original)
data2, err := encoder.Marshal(bytesValue)
if err != nil {
t.Fatalf("Marshal Bytes failed: %v", err)
}

// 修改返回的数据
if len(data2) > 0 {
data2[0] = ^data2[0]
}

// 验证 Bytes.Data 未被修改
if !BytesEqual(bytesValue.Data, originalCopy) {
t.Fatalf("Deep copy failed for Bytes: modifying Marshal result affected Bytes.Data")
}

// 测试 *Bytes 类型的深拷贝
bytesPtr := encodingx.NewBytes()
bytesPtr.Data = make([]byte, len(original))
copy(bytesPtr.Data, original)

data3, err := encoder.Marshal(bytesPtr)
if err != nil {
t.Fatalf("Marshal *Bytes failed: %v", err)
}

// 修改返回的数据
if len(data3) > 0 {
data3[0] = ^data3[0]
}

// 验证 *Bytes.Data 未被修改
if !BytesEqual(bytesPtr.Data, originalCopy) {
t.Fatalf("Deep copy failed for *Bytes: modifying Marshal result affected *Bytes.Data")
}

// 测试 Unmarshal 的深拷贝
result := encodingx.NewBytes()
err = encoder.Unmarshal(original, result)
if err != nil {
t.Fatalf("Unmarshal failed: %v", err)
}

// 修改 result.Data
if len(result.Data) > 0 {
result.Data[0] = ^result.Data[0]
}

// 验证原始数据未被修改
if !BytesEqual(original, originalCopy) {
t.Fatalf("Deep copy failed for Unmarshal: modifying result affected original data")
}
})
}

