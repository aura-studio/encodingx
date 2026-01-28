package encodingx_test

import (
	"errors"
	"testing"

	"github.com/aura-studio/encodingx"
	"pgregory.net/rapid"
)

// ============================================================================
// Group1 EncodingSet 和辅助函数测试
// Validates: Requirements 12.1, 12.2, 12.3, 13.1, 13.2, 13.3, 13.4, 13.5, 13.6
// ============================================================================

// ============================================================================
// Empty() 函数测试
// Validates: Requirements 12.3
// ============================================================================

// TestGroup1_Empty 测试 Empty() 返回默认 ChainEncoding
// Validates: Requirements 12.3
func TestGroup1_Empty(t *testing.T) {
	t.Run("Empty returns non-nil Encoding", func(t *testing.T) {
		e := encodingx.Empty()
		if e == nil {
			t.Error("Empty() should return non-nil Encoding")
		}
	})

	t.Run("Empty returns ChainEncoding with Mix style", func(t *testing.T) {
		e := encodingx.Empty()
		// ChainEncoding 的 Style() 应该返回 EncodingStyleMix
		if e.Style() != encodingx.EncodingStyleMix {
			t.Errorf("Empty() should return Encoding with Mix style, got %v", e.Style())
		}
	})

	t.Run("Empty returns consistent result", func(t *testing.T) {
		e1 := encodingx.Empty()
		e2 := encodingx.Empty()
		// 两次调用应该返回相同的结果
		if e1.String() != e2.String() {
			t.Errorf("Empty() should return consistent result, got %s and %s", e1.String(), e2.String())
		}
	})

	t.Run("Empty String format", func(t *testing.T) {
		e := encodingx.Empty()
		str := e.String()
		// 验证 String() 返回非空字符串
		if str == "" {
			t.Error("Empty().String() should return non-empty string")
		}
	})

	t.Run("Empty can Marshal and Unmarshal", func(t *testing.T) {
		e := encodingx.Empty()
		// 测试 Empty() 返回的编码器可以正常工作
		input := []byte{1, 2, 3, 4, 5}
		data, err := e.Marshal(input)
		if err != nil {
			t.Errorf("Empty().Marshal() should not return error, got %v", err)
		}
		if !BytesEqual(data, input) {
			t.Errorf("Empty().Marshal() should return same data for []byte, got %v, want %v", data, input)
		}

		// 测试 Unmarshal
		var result encodingx.Bytes
		err = e.Unmarshal(data, &result)
		if err != nil {
			t.Errorf("Empty().Unmarshal() should not return error, got %v", err)
		}
		if !BytesEqual(result.Data, input) {
			t.Errorf("Empty().Unmarshal() should restore data, got %v, want %v", result.Data, input)
		}
	})

	t.Run("Empty Reverse", func(t *testing.T) {
		e := encodingx.Empty()
		rev := e.Reverse()
		if rev == nil {
			t.Error("Empty().Reverse() should return non-nil Encoding")
		}
		// Reverse 的 Style 也应该是 Mix
		if rev.Style() != encodingx.EncodingStyleMix {
			t.Errorf("Empty().Reverse() should return Encoding with Mix style, got %v", rev.Style())
		}
	})
}

// ============================================================================
// localEncoding() 间接测试（通过 ChainEncoding）
// Validates: Requirements 12.1, 12.2
// ============================================================================

// TestGroup1_LocalEncoding_Registered 测试 localEncoding() 查找已注册编码器
// 通过 ChainEncoding 间接测试，因为 localEncoding 是内部函数
// Validates: Requirements 12.1
func TestGroup1_LocalEncoding_Registered(t *testing.T) {
	t.Run("ChainEncoding with registered JSON encoder", func(t *testing.T) {
		// 创建使用 JSON 编码器的 ChainEncoding
		chain := encodingx.NewChainEncoding([]string{"JSON"}, []string{"JSON"})

		// 测试 Marshal - 如果 localEncoding 找不到 JSON，会返回错误
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		data, err := chain.Marshal(input)
		if err != nil {
			t.Errorf("ChainEncoding with registered encoder should work, got error: %v", err)
		}
		if len(data) == 0 {
			t.Error("ChainEncoding.Marshal() should return non-empty data")
		}
	})

	t.Run("ChainEncoding with registered Lazy encoder", func(t *testing.T) {
		// 创建使用 Lazy 编码器的 ChainEncoding
		chain := encodingx.NewChainEncoding([]string{"Lazy"}, []string{"Lazy"})

		// 测试 Marshal
		input := []byte{1, 2, 3, 4, 5}
		data, err := chain.Marshal(input)
		if err != nil {
			t.Errorf("ChainEncoding with Lazy encoder should work, got error: %v", err)
		}
		if !BytesEqual(data, input) {
			t.Errorf("ChainEncoding with Lazy should return same data, got %v, want %v", data, input)
		}
	})

	t.Run("ChainEncoding with registered Base64 encoder", func(t *testing.T) {
		// 创建使用 Base64 编码器的 ChainEncoding
		chain := encodingx.NewChainEncoding([]string{"Base64"}, []string{"Base64"})

		// 测试 Marshal
		input := []byte{1, 2, 3, 4, 5}
		data, err := chain.Marshal(input)
		if err != nil {
			t.Errorf("ChainEncoding with Base64 encoder should work, got error: %v", err)
		}
		if len(data) == 0 {
			t.Error("ChainEncoding.Marshal() with Base64 should return non-empty data")
		}
	})

	t.Run("ChainEncoding with multiple registered encoders", func(t *testing.T) {
		// 创建使用多个编码器的 ChainEncoding: JSON -> Base64
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})

		// 测试 Marshal
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		data, err := chain.Marshal(input)
		if err != nil {
			t.Errorf("ChainEncoding with multiple encoders should work, got error: %v", err)
		}
		if len(data) == 0 {
			t.Error("ChainEncoding.Marshal() should return non-empty data")
		}

		// 测试 Unmarshal
		var result TestStruct
		err = chain.Unmarshal(data, &result)
		if err != nil {
			t.Errorf("ChainEncoding.Unmarshal() should work, got error: %v", err)
		}
		if !input.Equal(result) {
			t.Errorf("ChainEncoding round-trip failed, got %+v, want %+v", result, input)
		}
	})
}

// TestGroup1_LocalEncoding_NotRegistered 测试 localEncoding() 查找未注册编码器返回错误
// 通过 ChainEncoding 间接测试
// Validates: Requirements 12.2
func TestGroup1_LocalEncoding_NotRegistered(t *testing.T) {
	t.Run("ChainEncoding with unregistered encoder returns error on Marshal", func(t *testing.T) {
		// 创建使用不存在编码器的 ChainEncoding
		chain := encodingx.NewChainEncoding([]string{"NonExistentEncoder"}, []string{"NonExistentEncoder"})

		// 测试 Marshal - 应该返回 ErrEncodingMissingEncoding 错误
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		_, err := chain.Marshal(input)
		if err == nil {
			t.Error("ChainEncoding with unregistered encoder should return error")
		}
		if !errors.Is(err, encodingx.ErrEncodingMissingEncoding) {
			t.Errorf("Expected ErrEncodingMissingEncoding, got %v", err)
		}
	})

	t.Run("ChainEncoding with unregistered encoder returns error on Unmarshal", func(t *testing.T) {
		// 创建使用不存在编码器的 ChainEncoding
		chain := encodingx.NewChainEncoding([]string{"Lazy"}, []string{"NonExistentDecoder"})

		// 测试 Unmarshal - 应该返回 ErrEncodingMissingEncoding 错误
		data := []byte{1, 2, 3}
		var result encodingx.Bytes
		err := chain.Unmarshal(data, &result)
		if err == nil {
			t.Error("ChainEncoding with unregistered decoder should return error")
		}
		if !errors.Is(err, encodingx.ErrEncodingMissingEncoding) {
			t.Errorf("Expected ErrEncodingMissingEncoding, got %v", err)
		}
	})

	t.Run("ChainEncoding with mixed registered and unregistered encoders", func(t *testing.T) {
		// 第一个编码器存在，第二个不存在
		chain := encodingx.NewChainEncoding([]string{"JSON", "NonExistent"}, []string{"NonExistent", "JSON"})

		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		_, err := chain.Marshal(input)
		if err == nil {
			t.Error("ChainEncoding with unregistered encoder in chain should return error")
		}
		if !errors.Is(err, encodingx.ErrEncodingMissingEncoding) {
			t.Errorf("Expected ErrEncodingMissingEncoding, got %v", err)
		}
	})
}

// ============================================================================
// Marshal/Unmarshal 辅助函数测试
// Validates: Requirements 13.1, 13.2
// ============================================================================

// TestGroup1_Marshal 测试 Marshal 辅助函数
// Validates: Requirements 13.1
func TestGroup1_Marshal(t *testing.T) {
	t.Run("Marshal with JSON encoder", func(t *testing.T) {
		enc := encodingx.NewJSON()
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		// 使用辅助函数
		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Errorf("Marshal() should not return error, got %v", err)
		}

		// 直接调用编码器
		expected, _ := enc.Marshal(input)

		// 验证结果相同
		if !BytesEqual(data, expected) {
			t.Errorf("Marshal() should return same result as encoder.Marshal(), got %s, want %s",
				string(data), string(expected))
		}
	})

	t.Run("Marshal with Lazy encoder", func(t *testing.T) {
		enc := encodingx.NewLazy()
		input := []byte{1, 2, 3, 4, 5}

		// 使用辅助函数
		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Errorf("Marshal() should not return error, got %v", err)
		}

		// 直接调用编码器
		expected, _ := enc.Marshal(input)

		// 验证结果相同
		if !BytesEqual(data, expected) {
			t.Errorf("Marshal() should return same result as encoder.Marshal()")
		}
	})

	t.Run("Marshal with Base64 encoder", func(t *testing.T) {
		enc := encodingx.NewBase64()
		input := []byte{1, 2, 3, 4, 5}

		// 使用辅助函数
		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Errorf("Marshal() should not return error, got %v", err)
		}

		// 直接调用编码器
		expected, _ := enc.Marshal(input)

		// 验证结果相同
		if !BytesEqual(data, expected) {
			t.Errorf("Marshal() should return same result as encoder.Marshal()")
		}
	})

	t.Run("Marshal returns error when encoder fails", func(t *testing.T) {
		enc := encodingx.NewBase64()
		// Base64 编码器对非字节类型返回错误
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		_, err := encodingx.Marshal(enc, input)
		if err == nil {
			t.Error("Marshal() should return error when encoder fails")
		}
	})

	t.Run("Marshal with ChainEncoding", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		// 使用辅助函数
		data, err := encodingx.Marshal(chain, input)
		if err != nil {
			t.Errorf("Marshal() with ChainEncoding should not return error, got %v", err)
		}

		// 直接调用编码器
		expected, _ := chain.Marshal(input)

		// 验证结果相同
		if !BytesEqual(data, expected) {
			t.Errorf("Marshal() should return same result as ChainEncoding.Marshal()")
		}
	})
}

// TestGroup1_Unmarshal 测试 Unmarshal 辅助函数
// Validates: Requirements 13.2
func TestGroup1_Unmarshal(t *testing.T) {
	t.Run("Unmarshal with JSON encoder", func(t *testing.T) {
		enc := encodingx.NewJSON()
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		data, _ := enc.Marshal(input)

		// 使用辅助函数
		var result TestStruct
		err := encodingx.Unmarshal(enc, data, &result)
		if err != nil {
			t.Errorf("Unmarshal() should not return error, got %v", err)
		}

		// 验证结果
		if !input.Equal(result) {
			t.Errorf("Unmarshal() should restore data, got %+v, want %+v", result, input)
		}
	})

	t.Run("Unmarshal with Lazy encoder", func(t *testing.T) {
		enc := encodingx.NewLazy()
		input := []byte{1, 2, 3, 4, 5}
		data, _ := enc.Marshal(input)

		// 使用辅助函数
		var result encodingx.Bytes
		err := encodingx.Unmarshal(enc, data, &result)
		if err != nil {
			t.Errorf("Unmarshal() should not return error, got %v", err)
		}

		// 验证结果
		if !BytesEqual(result.Data, input) {
			t.Errorf("Unmarshal() should restore data, got %v, want %v", result.Data, input)
		}
	})

	t.Run("Unmarshal returns error when decoder fails", func(t *testing.T) {
		enc := encodingx.NewJSON()
		// 无效的 JSON 数据
		invalidData := []byte("not valid json{{{")

		var result TestStruct
		err := encodingx.Unmarshal(enc, invalidData, &result)
		if err == nil {
			t.Error("Unmarshal() should return error when decoder fails")
		}
	})

	t.Run("Unmarshal with ChainEncoding", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		data, _ := chain.Marshal(input)

		// 使用辅助函数
		var result TestStruct
		err := encodingx.Unmarshal(chain, data, &result)
		if err != nil {
			t.Errorf("Unmarshal() with ChainEncoding should not return error, got %v", err)
		}

		// 验证结果
		if !input.Equal(result) {
			t.Errorf("Unmarshal() should restore data, got %+v, want %+v", result, input)
		}
	})
}

// ============================================================================
// Encode/Decode 辅助函数测试（包括 panic 情况）
// Validates: Requirements 13.3, 13.4, 13.5, 13.6
// ============================================================================

// TestGroup1_Encode 测试 Encode 辅助函数
// Validates: Requirements 13.3, 13.4
func TestGroup1_Encode(t *testing.T) {
	t.Run("Encode success returns byte array", func(t *testing.T) {
		// Validates: Requirements 13.3
		enc := encodingx.NewJSON()
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		// Encode 成功时返回字节数组
		data := encodingx.Encode(enc, input)
		if len(data) == 0 {
			t.Error("Encode() should return non-empty byte array on success")
		}

		// 验证结果与 Marshal 相同
		expected, _ := enc.Marshal(input)
		if !BytesEqual(data, expected) {
			t.Errorf("Encode() should return same result as Marshal(), got %s, want %s",
				string(data), string(expected))
		}
	})

	t.Run("Encode with Lazy encoder", func(t *testing.T) {
		enc := encodingx.NewLazy()
		input := []byte{1, 2, 3, 4, 5}

		data := encodingx.Encode(enc, input)
		if !BytesEqual(data, input) {
			t.Errorf("Encode() with Lazy should return same data, got %v, want %v", data, input)
		}
	})

	t.Run("Encode panics on error", func(t *testing.T) {
		// Validates: Requirements 13.4
		enc := encodingx.NewBase64()
		// Base64 编码器对非字节类型返回错误
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		defer func() {
			if r := recover(); r == nil {
				t.Error("Encode() should panic when encoding fails")
			}
		}()

		// 这应该 panic
		encodingx.Encode(enc, input)
	})

	t.Run("Encode with ChainEncoding success", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		data := encodingx.Encode(chain, input)
		if len(data) == 0 {
			t.Error("Encode() with ChainEncoding should return non-empty data")
		}
	})

	t.Run("Encode with ChainEncoding panics on unregistered encoder", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"NonExistent"}, []string{"NonExistent"})
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		defer func() {
			if r := recover(); r == nil {
				t.Error("Encode() should panic when ChainEncoding has unregistered encoder")
			}
		}()

		encodingx.Encode(chain, input)
	})
}

// TestGroup1_Decode 测试 Decode 辅助函数
// Validates: Requirements 13.5, 13.6
func TestGroup1_Decode(t *testing.T) {
	t.Run("Decode success returns normally", func(t *testing.T) {
		// Validates: Requirements 13.5
		enc := encodingx.NewJSON()
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		data, _ := enc.Marshal(input)

		var result TestStruct
		// Decode 成功时正常返回（不 panic）
		encodingx.Decode(enc, data, &result)

		// 验证结果
		if !input.Equal(result) {
			t.Errorf("Decode() should restore data, got %+v, want %+v", result, input)
		}
	})

	t.Run("Decode with Lazy encoder", func(t *testing.T) {
		enc := encodingx.NewLazy()
		input := []byte{1, 2, 3, 4, 5}
		data, _ := enc.Marshal(input)

		var result encodingx.Bytes
		encodingx.Decode(enc, data, &result)

		if !BytesEqual(result.Data, input) {
			t.Errorf("Decode() with Lazy should restore data, got %v, want %v", result.Data, input)
		}
	})

	t.Run("Decode panics on error", func(t *testing.T) {
		// Validates: Requirements 13.6
		enc := encodingx.NewJSON()
		// 无效的 JSON 数据
		invalidData := []byte("not valid json{{{")

		defer func() {
			if r := recover(); r == nil {
				t.Error("Decode() should panic when decoding fails")
			}
		}()

		var result TestStruct
		// 这应该 panic
		encodingx.Decode(enc, invalidData, &result)
	})

	t.Run("Decode with ChainEncoding success", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}
		data, _ := chain.Marshal(input)

		var result TestStruct
		encodingx.Decode(chain, data, &result)

		if !input.Equal(result) {
			t.Errorf("Decode() with ChainEncoding should restore data, got %+v, want %+v", result, input)
		}
	})

	t.Run("Decode with ChainEncoding panics on unregistered decoder", func(t *testing.T) {
		// 创建一个可以 Marshal 但不能 Unmarshal 的 ChainEncoding
		chain := encodingx.NewChainEncoding([]string{"Lazy"}, []string{"NonExistent"})
		data := []byte{1, 2, 3}

		defer func() {
			if r := recover(); r == nil {
				t.Error("Decode() should panic when ChainEncoding has unregistered decoder")
			}
		}()

		var result encodingx.Bytes
		encodingx.Decode(chain, data, &result)
	})

	t.Run("Decode with invalid Base64 data panics", func(t *testing.T) {
		enc := encodingx.NewBase64()
		// 无效的 Base64 数据
		invalidData := []byte("!!!not-valid-base64!!!")

		defer func() {
			if r := recover(); r == nil {
				t.Error("Decode() should panic when Base64 decoding fails")
			}
		}()

		var result encodingx.Bytes
		encodingx.Decode(enc, invalidData, &result)
	})
}

// ============================================================================
// 集成测试
// ============================================================================

// TestGroup1_Encoding_Integration 集成测试
func TestGroup1_Encoding_Integration(t *testing.T) {
	t.Run("Marshal then Unmarshal round-trip", func(t *testing.T) {
		enc := encodingx.NewJSON()
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		// Marshal
		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Fatalf("Marshal() failed: %v", err)
		}

		// Unmarshal
		var result TestStruct
		err = encodingx.Unmarshal(enc, data, &result)
		if err != nil {
			t.Fatalf("Unmarshal() failed: %v", err)
		}

		// 验证
		if !input.Equal(result) {
			t.Errorf("Round-trip failed, got %+v, want %+v", result, input)
		}
	})

	t.Run("Encode then Decode round-trip", func(t *testing.T) {
		enc := encodingx.NewJSON()
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		// Encode
		data := encodingx.Encode(enc, input)

		// Decode
		var result TestStruct
		encodingx.Decode(enc, data, &result)

		// 验证
		if !input.Equal(result) {
			t.Errorf("Round-trip failed, got %+v, want %+v", result, input)
		}
	})

	t.Run("Empty encoding round-trip with bytes", func(t *testing.T) {
		e := encodingx.Empty()
		input := []byte{1, 2, 3, 4, 5}

		// Marshal
		data, err := encodingx.Marshal(e, input)
		if err != nil {
			t.Fatalf("Marshal() with Empty() failed: %v", err)
		}

		// Unmarshal
		var result encodingx.Bytes
		err = encodingx.Unmarshal(e, data, &result)
		if err != nil {
			t.Fatalf("Unmarshal() with Empty() failed: %v", err)
		}

		// 验证
		if !BytesEqual(result.Data, input) {
			t.Errorf("Round-trip with Empty() failed, got %v, want %v", result.Data, input)
		}
	})

	t.Run("ChainEncoding with JSON and Base64 round-trip", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		// Marshal
		data, err := encodingx.Marshal(chain, input)
		if err != nil {
			t.Fatalf("Marshal() with ChainEncoding failed: %v", err)
		}

		// Unmarshal
		var result TestStruct
		err = encodingx.Unmarshal(chain, data, &result)
		if err != nil {
			t.Fatalf("Unmarshal() with ChainEncoding failed: %v", err)
		}

		// 验证
		if !input.Equal(result) {
			t.Errorf("Round-trip with ChainEncoding failed, got %+v, want %+v", result, input)
		}
	})

	t.Run("Multiple encoders consistency", func(t *testing.T) {
		// 测试不同编码器的 Marshal/Unmarshal 辅助函数一致性
		encoders := []encodingx.Encoding{
			encodingx.NewJSON(),
			encodingx.NewYAML(),
		}

		input := TestStruct{Integer: 42, String: "test", Bool: true, Float: 3.14}

		for _, enc := range encoders {
			// 使用辅助函数
			data1, err1 := encodingx.Marshal(enc, input)
			// 直接调用
			data2, err2 := enc.Marshal(input)

			if (err1 == nil) != (err2 == nil) {
				t.Errorf("Marshal() and encoder.Marshal() should have same error behavior for %s", enc.String())
			}

			if err1 == nil && !BytesEqual(data1, data2) {
				t.Errorf("Marshal() and encoder.Marshal() should return same result for %s", enc.String())
			}
		}
	})
}

// ============================================================================
// 边界情况测试
// ============================================================================

// TestGroup1_Encoding_EdgeCases 边界情况测试
func TestGroup1_Encoding_EdgeCases(t *testing.T) {
	t.Run("Marshal with nil value", func(t *testing.T) {
		enc := encodingx.NewJSON()
		data, err := encodingx.Marshal(enc, nil)
		// JSON 编码 nil 应该返回 "null"
		if err != nil {
			t.Errorf("Marshal(nil) should not return error, got %v", err)
		}
		if string(data) != "null" {
			t.Errorf("Marshal(nil) should return 'null', got %s", string(data))
		}
	})

	t.Run("Marshal with empty struct", func(t *testing.T) {
		enc := encodingx.NewJSON()
		input := struct{}{}
		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Errorf("Marshal(empty struct) should not return error, got %v", err)
		}
		if string(data) != "{}" {
			t.Errorf("Marshal(empty struct) should return '{}', got %s", string(data))
		}
	})

	t.Run("Unmarshal to nil pointer panics", func(t *testing.T) {
		enc := encodingx.NewJSON()
		data := []byte(`{"integer":42}`)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Unmarshal to nil pointer should panic")
			}
		}()

		var result *TestStruct = nil
		encodingx.Decode(enc, data, result)
	})

	t.Run("Empty byte slice Marshal", func(t *testing.T) {
		enc := encodingx.NewLazy()
		input := []byte{}
		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Errorf("Marshal(empty []byte) should not return error, got %v", err)
		}
		if len(data) != 0 {
			t.Errorf("Marshal(empty []byte) should return empty data, got len=%d", len(data))
		}
	})

	t.Run("Large data Marshal and Unmarshal", func(t *testing.T) {
		enc := encodingx.NewLazy()
		input := make([]byte, 10000)
		for i := range input {
			input[i] = byte(i % 256)
		}

		data, err := encodingx.Marshal(enc, input)
		if err != nil {
			t.Fatalf("Marshal(large data) failed: %v", err)
		}

		var result encodingx.Bytes
		err = encodingx.Unmarshal(enc, data, &result)
		if err != nil {
			t.Fatalf("Unmarshal(large data) failed: %v", err)
		}

		if !BytesEqual(result.Data, input) {
			t.Error("Large data round-trip failed")
		}
	})
}

// ============================================================================
// 随机数据测试
// ============================================================================

// TestGroup1_Encoding_RandomData 使用随机数据测试
func TestGroup1_Encoding_RandomData(t *testing.T) {
	gen := NewTestDataGenerator()

	t.Run("Marshal/Unmarshal with random TestStruct", func(t *testing.T) {
		enc := encodingx.NewJSON()
		for i := 0; i < 10; i++ {
			input := gen.GenerateTestStruct()

			data, err := encodingx.Marshal(enc, input)
			if err != nil {
				t.Errorf("Iteration %d: Marshal() failed: %v", i, err)
				continue
			}

			var result TestStruct
			err = encodingx.Unmarshal(enc, data, &result)
			if err != nil {
				t.Errorf("Iteration %d: Unmarshal() failed: %v", i, err)
				continue
			}

			if !input.Equal(result) {
				t.Errorf("Iteration %d: Round-trip failed, got %+v, want %+v", i, result, input)
			}
		}
	})

	t.Run("Encode/Decode with random bytes", func(t *testing.T) {
		enc := encodingx.NewLazy()
		for i := 0; i < 10; i++ {
			input := gen.GenerateBytes(1, 100)

			data := encodingx.Encode(enc, input)

			var result encodingx.Bytes
			encodingx.Decode(enc, data, &result)

			if !BytesEqual(result.Data, input) {
				t.Errorf("Iteration %d: Round-trip failed", i)
			}
		}
	})

	t.Run("ChainEncoding with random data", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		for i := 0; i < 10; i++ {
			input := gen.GenerateTestStruct()

			data, err := encodingx.Marshal(chain, input)
			if err != nil {
				t.Errorf("Iteration %d: Marshal() with ChainEncoding failed: %v", i, err)
				continue
			}

			var result TestStruct
			err = encodingx.Unmarshal(chain, data, &result)
			if err != nil {
				t.Errorf("Iteration %d: Unmarshal() with ChainEncoding failed: %v", i, err)
				continue
			}

			if !input.Equal(result) {
				t.Errorf("Iteration %d: ChainEncoding round-trip failed, got %+v, want %+v", i, result, input)
			}
		}
	})
}

// ============================================================================
// Property 21: Marshal/Unmarshal 辅助函数等价性属性测试
// **Validates: Requirements 13.1, 13.2**
// ============================================================================

// TestGroup1_Property21_MarshalUnmarshalEquivalence 属性测试
// Property 21: For any 编码器和数据，Marshal(e, v) 应该等价于 e.Marshal(v)，
// Unmarshal(e, data, v) 应该等价于 e.Unmarshal(data, v)
// **Validates: Requirements 13.1, 13.2**
func TestGroup1_Property21_MarshalUnmarshalEquivalence(t *testing.T) {
	const iterations = 100
	gen := NewTestDataGenerator()

	// 定义测试用的编码器类型
	type encoderTestCase struct {
		name        string
		encoder     encodingx.Encoding
		genInput    func() interface{}
		genTarget   func() interface{}
		compareFunc func(a, b interface{}) bool
	}

	// 创建测试用例列表
	testCases := []encoderTestCase{
		{
			name:    "JSON with TestStruct",
			encoder: encodingx.NewJSON(),
			genInput: func() interface{} {
				return gen.GenerateTestStruct()
			},
			genTarget: func() interface{} {
				return new(TestStruct)
			},
			compareFunc: func(a, b interface{}) bool {
				aVal, aOk := a.(TestStruct)
				bVal, bOk := b.(*TestStruct)
				if !aOk || !bOk {
					return false
				}
				return aVal.Equal(*bVal)
			},
		},
		{
			name:    "YAML with TestStruct",
			encoder: encodingx.NewYAML(),
			genInput: func() interface{} {
				return gen.GenerateTestStruct()
			},
			genTarget: func() interface{} {
				return new(TestStruct)
			},
			compareFunc: func(a, b interface{}) bool {
				aVal, aOk := a.(TestStruct)
				bVal, bOk := b.(*TestStruct)
				if !aOk || !bOk {
					return false
				}
				return aVal.Equal(*bVal)
			},
		},
		{
			name:    "Lazy with []byte",
			encoder: encodingx.NewLazy(),
			genInput: func() interface{} {
				return gen.GenerateBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(a, b interface{}) bool {
				aVal, aOk := a.([]byte)
				bVal, bOk := b.(*encodingx.Bytes)
				if !aOk || !bOk {
					return false
				}
				return BytesEqual(aVal, bVal.Data)
			},
		},
		{
			name:    "Lazy with Bytes",
			encoder: encodingx.NewLazy(),
			genInput: func() interface{} {
				return gen.GenerateEncodingxBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(a, b interface{}) bool {
				aVal, aOk := a.(encodingx.Bytes)
				bVal, bOk := b.(*encodingx.Bytes)
				if !aOk || !bOk {
					return false
				}
				return BytesEqual(aVal.Data, bVal.Data)
			},
		},
		{
			name:    "Base64 with []byte",
			encoder: encodingx.NewBase64(),
			genInput: func() interface{} {
				return gen.GenerateBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(a, b interface{}) bool {
				aVal, aOk := a.([]byte)
				bVal, bOk := b.(*encodingx.Bytes)
				if !aOk || !bOk {
					return false
				}
				return BytesEqual(aVal, bVal.Data)
			},
		},
		{
			name:    "Base64URL with []byte",
			encoder: encodingx.NewBase64URL(),
			genInput: func() interface{} {
				return gen.GenerateBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(a, b interface{}) bool {
				aVal, aOk := a.([]byte)
				bVal, bOk := b.(*encodingx.Bytes)
				if !aOk || !bOk {
					return false
				}
				return BytesEqual(aVal, bVal.Data)
			},
		},
	}

	// 对每个编码器运行属性测试
	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				input := tc.genInput()

				// 测试 Marshal 等价性
				// Marshal(e, v) 应该等价于 e.Marshal(v)
				helperResult, helperErr := encodingx.Marshal(tc.encoder, input)
				directResult, directErr := tc.encoder.Marshal(input)

				// 验证错误行为一致
				if (helperErr == nil) != (directErr == nil) {
					t.Errorf("Iteration %d: Marshal error behavior mismatch: helper=%v, direct=%v",
						i, helperErr, directErr)
					continue
				}

				// 验证结果一致
				if helperErr == nil && !BytesEqual(helperResult, directResult) {
					t.Errorf("Iteration %d: Marshal result mismatch: helper=%v, direct=%v",
						i, helperResult, directResult)
					continue
				}

				// 如果 Marshal 成功，测试 Unmarshal 等价性
				if helperErr == nil {
					// 使用辅助函数 Unmarshal
					helperTarget := tc.genTarget()
					helperUnmarshalErr := encodingx.Unmarshal(tc.encoder, helperResult, helperTarget)

					// 直接调用编码器 Unmarshal
					directTarget := tc.genTarget()
					directUnmarshalErr := tc.encoder.Unmarshal(directResult, directTarget)

					// 验证错误行为一致
					if (helperUnmarshalErr == nil) != (directUnmarshalErr == nil) {
						t.Errorf("Iteration %d: Unmarshal error behavior mismatch: helper=%v, direct=%v",
							i, helperUnmarshalErr, directUnmarshalErr)
						continue
					}

					// 验证反序列化后的数据与原始数据一致
					if helperUnmarshalErr == nil {
						if !tc.compareFunc(input, helperTarget) {
							t.Errorf("Iteration %d: Helper Unmarshal result doesn't match original input", i)
						}
						if !tc.compareFunc(input, directTarget) {
							t.Errorf("Iteration %d: Direct Unmarshal result doesn't match original input", i)
						}
					}
				}
			}
		})
	}
}

// TestGroup1_Property21_MarshalEquivalence_MultipleEncoders 测试多种编码器的 Marshal 等价性
// **Validates: Requirements 13.1**
func TestGroup1_Property21_MarshalEquivalence_MultipleEncoders(t *testing.T) {
	const iterations = 100
	gen := NewTestDataGenerator()

	t.Run("JSON encoder Marshal equivalence", func(t *testing.T) {
		enc := encodingx.NewJSON()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()

			helperResult, helperErr := encodingx.Marshal(enc, input)
			directResult, directErr := enc.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("YAML encoder Marshal equivalence", func(t *testing.T) {
		enc := encodingx.NewYAML()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()

			helperResult, helperErr := encodingx.Marshal(enc, input)
			directResult, directErr := enc.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Lazy encoder Marshal equivalence", func(t *testing.T) {
		enc := encodingx.NewLazy()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)

			helperResult, helperErr := encodingx.Marshal(enc, input)
			directResult, directErr := enc.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Base64 encoder Marshal equivalence", func(t *testing.T) {
		enc := encodingx.NewBase64()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)

			helperResult, helperErr := encodingx.Marshal(enc, input)
			directResult, directErr := enc.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Base64URL encoder Marshal equivalence", func(t *testing.T) {
		enc := encodingx.NewBase64URL()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)

			helperResult, helperErr := encodingx.Marshal(enc, input)
			directResult, directErr := enc.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})
}

// TestGroup1_Property21_UnmarshalEquivalence_MultipleEncoders 测试多种编码器的 Unmarshal 等价性
// **Validates: Requirements 13.2**
func TestGroup1_Property21_UnmarshalEquivalence_MultipleEncoders(t *testing.T) {
	const iterations = 100
	gen := NewTestDataGenerator()

	t.Run("JSON encoder Unmarshal equivalence", func(t *testing.T) {
		enc := encodingx.NewJSON()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()
			data, _ := enc.Marshal(input)

			var helperResult TestStruct
			helperErr := encodingx.Unmarshal(enc, data, &helperResult)

			var directResult TestStruct
			directErr := enc.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !helperResult.Equal(directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("YAML encoder Unmarshal equivalence", func(t *testing.T) {
		enc := encodingx.NewYAML()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()
			data, _ := enc.Marshal(input)

			var helperResult TestStruct
			helperErr := encodingx.Unmarshal(enc, data, &helperResult)

			var directResult TestStruct
			directErr := enc.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !helperResult.Equal(directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Lazy encoder Unmarshal equivalence", func(t *testing.T) {
		enc := encodingx.NewLazy()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)
			data, _ := enc.Marshal(input)

			var helperResult encodingx.Bytes
			helperErr := encodingx.Unmarshal(enc, data, &helperResult)

			var directResult encodingx.Bytes
			directErr := enc.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult.Data, directResult.Data) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Base64 encoder Unmarshal equivalence", func(t *testing.T) {
		enc := encodingx.NewBase64()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)
			data, _ := enc.Marshal(input)

			var helperResult encodingx.Bytes
			helperErr := encodingx.Unmarshal(enc, data, &helperResult)

			var directResult encodingx.Bytes
			directErr := enc.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult.Data, directResult.Data) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Base64URL encoder Unmarshal equivalence", func(t *testing.T) {
		enc := encodingx.NewBase64URL()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)
			data, _ := enc.Marshal(input)

			var helperResult encodingx.Bytes
			helperErr := encodingx.Unmarshal(enc, data, &helperResult)

			var directResult encodingx.Bytes
			directErr := enc.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult.Data, directResult.Data) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})
}

// TestGroup1_Property21_ChainEncoding_Equivalence 测试 ChainEncoding 的辅助函数等价性
// **Validates: Requirements 13.1, 13.2**
func TestGroup1_Property21_ChainEncoding_Equivalence(t *testing.T) {
	const iterations = 100
	gen := NewTestDataGenerator()

	t.Run("ChainEncoding JSON->Base64 Marshal equivalence", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()

			helperResult, helperErr := encodingx.Marshal(chain, input)
			directResult, directErr := chain.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("ChainEncoding JSON->Base64 Unmarshal equivalence", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()
			data, _ := chain.Marshal(input)

			var helperResult TestStruct
			helperErr := encodingx.Unmarshal(chain, data, &helperResult)

			var directResult TestStruct
			directErr := chain.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !helperResult.Equal(directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Empty ChainEncoding Marshal equivalence", func(t *testing.T) {
		empty := encodingx.Empty()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)

			helperResult, helperErr := encodingx.Marshal(empty, input)
			directResult, directErr := empty.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})

	t.Run("Empty ChainEncoding Unmarshal equivalence", func(t *testing.T) {
		empty := encodingx.Empty()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)
			data, _ := empty.Marshal(input)

			var helperResult encodingx.Bytes
			helperErr := encodingx.Unmarshal(empty, data, &helperResult)

			var directResult encodingx.Bytes
			directErr := empty.Unmarshal(data, &directResult)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult.Data, directResult.Data) {
				t.Errorf("Iteration %d: Result mismatch", i)
			}
		}
	})
}

// TestGroup1_Property21_ErrorBehavior_Equivalence 测试错误情况下的辅助函数等价性
// **Validates: Requirements 13.1, 13.2**
func TestGroup1_Property21_ErrorBehavior_Equivalence(t *testing.T) {
	const iterations = 100
	gen := NewTestDataGenerator()

	t.Run("Base64 with wrong type returns same error", func(t *testing.T) {
		enc := encodingx.NewBase64()
		for i := 0; i < iterations; i++ {
			// Base64 编码器对非字节类型返回错误
			input := gen.GenerateTestStruct()

			_, helperErr := encodingx.Marshal(enc, input)
			_, directErr := enc.Marshal(input)

			// 两者都应该返回错误
			if helperErr == nil || directErr == nil {
				t.Errorf("Iteration %d: Both should return error for wrong type", i)
				continue
			}

			// 错误类型应该一致
			if helperErr.Error() != directErr.Error() {
				t.Errorf("Iteration %d: Error message mismatch: helper=%v, direct=%v",
					i, helperErr, directErr)
			}
		}
	})

	t.Run("Lazy with wrong type returns same error", func(t *testing.T) {
		enc := encodingx.NewLazy()
		for i := 0; i < iterations; i++ {
			// Lazy 编码器对非字节类型返回错误
			input := gen.GenerateTestStruct()

			_, helperErr := encodingx.Marshal(enc, input)
			_, directErr := enc.Marshal(input)

			// 两者都应该返回错误
			if helperErr == nil || directErr == nil {
				t.Errorf("Iteration %d: Both should return error for wrong type", i)
				continue
			}

			// 错误类型应该一致
			if helperErr.Error() != directErr.Error() {
				t.Errorf("Iteration %d: Error message mismatch: helper=%v, direct=%v",
					i, helperErr, directErr)
			}
		}
	})

	t.Run("JSON Unmarshal with invalid data returns same error behavior", func(t *testing.T) {
		enc := encodingx.NewJSON()
		for i := 0; i < iterations; i++ {
			// 生成无效的 JSON 数据
			invalidData := gen.GenerateBytes(10, 50)

			var helperResult TestStruct
			helperErr := encodingx.Unmarshal(enc, invalidData, &helperResult)

			var directResult TestStruct
			directErr := enc.Unmarshal(invalidData, &directResult)

			// 两者的错误行为应该一致（都返回错误或都不返回错误）
			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Error behavior mismatch: helper=%v, direct=%v",
					i, helperErr, directErr)
			}
		}
	})
}

// ============================================================================
// Property 21: Marshal/Unmarshal 辅助函数等价性属性测试（综合测试）
// **Validates: Requirements 13.1, 13.2**
// ============================================================================

// TestGroup1_Property_21_MarshalUnmarshal_Equivalence 综合属性测试
// Property 21: For any 编码器和数据，Marshal(e, v) 应该等价于 e.Marshal(v)，
// Unmarshal(e, data, v) 应该等价于 e.Unmarshal(data, v)
// **Validates: Requirements 13.1, 13.2**
//
// 此测试使用多种编码器（JSON、YAML、Lazy、Base64 等）和随机生成的测试数据，
// 验证辅助函数和直接调用编码器方法的结果完全一致，运行至少 100 次迭代。
func TestGroup1_Property_21_MarshalUnmarshal_Equivalence(t *testing.T) {
	const iterations = 100
	gen := NewTestDataGenerator()

	// 测试用的编码器配置
	type encoderConfig struct {
		name        string
		encoder     encodingx.Encoding
		genInput    func() interface{}
		genTarget   func() interface{}
		compareFunc func(original interface{}, result interface{}) bool
	}

	// 创建多种编码器的测试配置
	configs := []encoderConfig{
		{
			name:    "JSON",
			encoder: encodingx.NewJSON(),
			genInput: func() interface{} {
				return gen.GenerateTestStruct()
			},
			genTarget: func() interface{} {
				return new(TestStruct)
			},
			compareFunc: func(original interface{}, result interface{}) bool {
				orig, ok1 := original.(TestStruct)
				res, ok2 := result.(*TestStruct)
				return ok1 && ok2 && orig.Equal(*res)
			},
		},
		{
			name:    "YAML",
			encoder: encodingx.NewYAML(),
			genInput: func() interface{} {
				return gen.GenerateTestStruct()
			},
			genTarget: func() interface{} {
				return new(TestStruct)
			},
			compareFunc: func(original interface{}, result interface{}) bool {
				orig, ok1 := original.(TestStruct)
				res, ok2 := result.(*TestStruct)
				return ok1 && ok2 && orig.Equal(*res)
			},
		},
		{
			name:    "Lazy_with_bytes",
			encoder: encodingx.NewLazy(),
			genInput: func() interface{} {
				return gen.GenerateBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(original interface{}, result interface{}) bool {
				orig, ok1 := original.([]byte)
				res, ok2 := result.(*encodingx.Bytes)
				return ok1 && ok2 && BytesEqual(orig, res.Data)
			},
		},
		{
			name:    "Lazy_with_Bytes",
			encoder: encodingx.NewLazy(),
			genInput: func() interface{} {
				return gen.GenerateEncodingxBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(original interface{}, result interface{}) bool {
				orig, ok1 := original.(encodingx.Bytes)
				res, ok2 := result.(*encodingx.Bytes)
				return ok1 && ok2 && BytesEqual(orig.Data, res.Data)
			},
		},
		{
			name:    "Base64",
			encoder: encodingx.NewBase64(),
			genInput: func() interface{} {
				return gen.GenerateBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(original interface{}, result interface{}) bool {
				orig, ok1 := original.([]byte)
				res, ok2 := result.(*encodingx.Bytes)
				return ok1 && ok2 && BytesEqual(orig, res.Data)
			},
		},
		{
			name:    "Base64URL",
			encoder: encodingx.NewBase64URL(),
			genInput: func() interface{} {
				return gen.GenerateBytes(1, 100)
			},
			genTarget: func() interface{} {
				return new(encodingx.Bytes)
			},
			compareFunc: func(original interface{}, result interface{}) bool {
				orig, ok1 := original.([]byte)
				res, ok2 := result.(*encodingx.Bytes)
				return ok1 && ok2 && BytesEqual(orig, res.Data)
			},
		},
	}

	// 对每种编码器运行属性测试
	for _, cfg := range configs {
		cfg := cfg // capture range variable
		t.Run(cfg.name, func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				input := cfg.genInput()

				// ========================================
				// 测试 Marshal 等价性
				// Marshal(e, v) 应该等价于 e.Marshal(v)
				// ========================================
				helperMarshalResult, helperMarshalErr := encodingx.Marshal(cfg.encoder, input)
				directMarshalResult, directMarshalErr := cfg.encoder.Marshal(input)

				// 验证错误行为一致
				if (helperMarshalErr == nil) != (directMarshalErr == nil) {
					t.Errorf("Iteration %d: Marshal error behavior mismatch - helper: %v, direct: %v",
						i, helperMarshalErr, directMarshalErr)
					continue
				}

				// 验证 Marshal 结果一致
				if helperMarshalErr == nil && !BytesEqual(helperMarshalResult, directMarshalResult) {
					t.Errorf("Iteration %d: Marshal result mismatch - helper: %v, direct: %v",
						i, helperMarshalResult, directMarshalResult)
					continue
				}

				// ========================================
				// 测试 Unmarshal 等价性
				// Unmarshal(e, data, v) 应该等价于 e.Unmarshal(data, v)
				// ========================================
				if helperMarshalErr == nil {
					// 使用辅助函数 Unmarshal
					helperTarget := cfg.genTarget()
					helperUnmarshalErr := encodingx.Unmarshal(cfg.encoder, helperMarshalResult, helperTarget)

					// 直接调用编码器 Unmarshal
					directTarget := cfg.genTarget()
					directUnmarshalErr := cfg.encoder.Unmarshal(directMarshalResult, directTarget)

					// 验证错误行为一致
					if (helperUnmarshalErr == nil) != (directUnmarshalErr == nil) {
						t.Errorf("Iteration %d: Unmarshal error behavior mismatch - helper: %v, direct: %v",
							i, helperUnmarshalErr, directUnmarshalErr)
						continue
					}

					// 验证 Unmarshal 结果与原始数据一致
					if helperUnmarshalErr == nil {
						if !cfg.compareFunc(input, helperTarget) {
							t.Errorf("Iteration %d: Helper Unmarshal result doesn't match original input", i)
						}
						if !cfg.compareFunc(input, directTarget) {
							t.Errorf("Iteration %d: Direct Unmarshal result doesn't match original input", i)
						}
					}
				}
			}
		})
	}

	// ========================================
	// 额外测试：ChainEncoding 的辅助函数等价性
	// ========================================
	t.Run("ChainEncoding_JSON_Base64", func(t *testing.T) {
		chain := encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"})
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()

			// 测试 Marshal 等价性
			helperResult, helperErr := encodingx.Marshal(chain, input)
			directResult, directErr := chain.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: ChainEncoding Marshal error mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: ChainEncoding Marshal result mismatch", i)
				continue
			}

			// 测试 Unmarshal 等价性
			if helperErr == nil {
				var helperTarget TestStruct
				helperUnmarshalErr := encodingx.Unmarshal(chain, helperResult, &helperTarget)

				var directTarget TestStruct
				directUnmarshalErr := chain.Unmarshal(directResult, &directTarget)

				if (helperUnmarshalErr == nil) != (directUnmarshalErr == nil) {
					t.Errorf("Iteration %d: ChainEncoding Unmarshal error mismatch", i)
					continue
				}
				if helperUnmarshalErr == nil {
					if !input.Equal(helperTarget) {
						t.Errorf("Iteration %d: ChainEncoding helper Unmarshal result mismatch", i)
					}
					if !input.Equal(directTarget) {
						t.Errorf("Iteration %d: ChainEncoding direct Unmarshal result mismatch", i)
					}
				}
			}
		}
	})

	// ========================================
	// 额外测试：Empty ChainEncoding 的辅助函数等价性
	// ========================================
	t.Run("Empty_ChainEncoding", func(t *testing.T) {
		empty := encodingx.Empty()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateBytes(1, 100)

			// 测试 Marshal 等价性
			helperResult, helperErr := encodingx.Marshal(empty, input)
			directResult, directErr := empty.Marshal(input)

			if (helperErr == nil) != (directErr == nil) {
				t.Errorf("Iteration %d: Empty Marshal error mismatch", i)
				continue
			}
			if helperErr == nil && !BytesEqual(helperResult, directResult) {
				t.Errorf("Iteration %d: Empty Marshal result mismatch", i)
				continue
			}

			// 测试 Unmarshal 等价性
			if helperErr == nil {
				var helperTarget encodingx.Bytes
				helperUnmarshalErr := encodingx.Unmarshal(empty, helperResult, &helperTarget)

				var directTarget encodingx.Bytes
				directUnmarshalErr := empty.Unmarshal(directResult, &directTarget)

				if (helperUnmarshalErr == nil) != (directUnmarshalErr == nil) {
					t.Errorf("Iteration %d: Empty Unmarshal error mismatch", i)
					continue
				}
				if helperUnmarshalErr == nil && !BytesEqual(helperTarget.Data, directTarget.Data) {
					t.Errorf("Iteration %d: Empty Unmarshal result mismatch", i)
				}
			}
		}
	})

	// ========================================
	// 额外测试：错误情况下的等价性
	// ========================================
	t.Run("Error_Behavior_Equivalence", func(t *testing.T) {
		// Base64 编码器对非字节类型返回错误
		enc := encodingx.NewBase64()
		for i := 0; i < iterations; i++ {
			input := gen.GenerateTestStruct()

			_, helperErr := encodingx.Marshal(enc, input)
			_, directErr := enc.Marshal(input)

			// 两者都应该返回错误
			if helperErr == nil || directErr == nil {
				t.Errorf("Iteration %d: Both should return error for wrong type", i)
				continue
			}

			// 错误消息应该一致
			if helperErr.Error() != directErr.Error() {
				t.Errorf("Iteration %d: Error message mismatch - helper: %v, direct: %v",
					i, helperErr, directErr)
			}
		}
	})
}

// ============================================================================
// 编码器接口一致性测试
// Validates: Requirements 14.1, 14.2, 14.3
// ============================================================================

// TestGroup1_EncoderInterfaceConsistency 测试所有编码器的接口一致性
// Validates: Requirements 14.1, 14.2, 14.3
func TestGroup1_EncoderInterfaceConsistency(t *testing.T) {
	// 获取所有已知的编码器
	encoders := []encodingx.Encoding{
		encodingx.NewJSON(),
		encodingx.NewYAML(),
		encodingx.NewXML(),
		encodingx.NewCSV(),
		encodingx.NewCSVWithHeaders(),
		encodingx.NewProtobuf(),
		encodingx.NewBase64(),
		encodingx.NewBase64URL(),
		encodingx.NewBinary(),
		encodingx.NewLittleEndian(),
		encodingx.NewBigEndian(),
		encodingx.NewHash(),
		encodingx.NewLazy(),
	}

	for _, enc := range encoders {
		enc := enc // capture range variable
		t.Run(enc.String(), func(t *testing.T) {
			// 测试 String() 返回非空字符串
			// Validates: Requirements 14.1
			name := enc.String()
			if name == "" {
				t.Errorf("%T.String() should return non-empty string", enc)
			}

			// 测试 Style() 返回有效的 EncodingStyleType
			// Validates: Requirements 14.2
			style := enc.Style()
			if style != encodingx.EncodingStyleStruct &&
				style != encodingx.EncodingStyleBytes &&
				style != encodingx.EncodingStyleMix {
				t.Errorf("%s.Style() returned invalid EncodingStyleType: %v", name, style)
			}

			// 测试 Reverse() 返回非 nil 的 Encoding
			// Validates: Requirements 14.3
			reversed := enc.Reverse()
			if reversed == nil {
				t.Errorf("%s.Reverse() should return non-nil Encoding", name)
			}

			// 验证 Reverse() 返回的也是有效的 Encoding
			if reversed != nil {
				revName := reversed.String()
				if revName == "" {
					t.Errorf("%s.Reverse().String() should return non-empty string", name)
				}

				revStyle := reversed.Style()
				if revStyle != encodingx.EncodingStyleStruct &&
					revStyle != encodingx.EncodingStyleBytes &&
					revStyle != encodingx.EncodingStyleMix {
					t.Errorf("%s.Reverse().Style() returned invalid EncodingStyleType: %v", name, revStyle)
				}

				revReversed := reversed.Reverse()
				if revReversed == nil {
					t.Errorf("%s.Reverse().Reverse() should return non-nil Encoding", name)
				}
			}
		})
	}
}

// TestGroup1_EncoderInterfaceConsistency_ChainEncoding 测试 ChainEncoding 的接口一致性
// Validates: Requirements 14.1, 14.2, 14.3
func TestGroup1_EncoderInterfaceConsistency_ChainEncoding(t *testing.T) {
	chains := []encodingx.Encoding{
		encodingx.Empty(),
		encodingx.NewChainEncoding([]string{"JSON"}, []string{"JSON"}),
		encodingx.NewChainEncoding([]string{"Lazy"}, []string{"Lazy"}),
		encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"}),
		encodingx.NewChainEncoding([]string{"JSON", "Base64", "Lazy"}, []string{"Lazy", "Base64", "JSON"}),
	}

	for i, chain := range chains {
		chain := chain
		t.Run(chain.String(), func(t *testing.T) {
			// 测试 String() 返回非空字符串
			name := chain.String()
			if name == "" {
				t.Errorf("ChainEncoding[%d].String() should return non-empty string", i)
			}

			// 测试 Style() 返回 EncodingStyleMix
			style := chain.Style()
			if style != encodingx.EncodingStyleMix {
				t.Errorf("ChainEncoding[%d].Style() should return EncodingStyleMix, got %v", i, style)
			}

			// 测试 Reverse() 返回非 nil 的 Encoding
			reversed := chain.Reverse()
			if reversed == nil {
				t.Errorf("ChainEncoding[%d].Reverse() should return non-nil Encoding", i)
			}

			// 验证 Reverse() 返回的也是有效的 ChainEncoding
			if reversed != nil {
				revStyle := reversed.Style()
				if revStyle != encodingx.EncodingStyleMix {
					t.Errorf("ChainEncoding[%d].Reverse().Style() should return EncodingStyleMix, got %v", i, revStyle)
				}
			}
		})
	}
}

// TestGroup1_EncoderInterfaceConsistency_StringFormat 测试编码器 String() 返回的格式
// Validates: Requirements 14.1
func TestGroup1_EncoderInterfaceConsistency_StringFormat(t *testing.T) {
	expectedNames := map[string]encodingx.Encoding{
		"JSON":           encodingx.NewJSON(),
		"YAML":           encodingx.NewYAML(),
		"XML":            encodingx.NewXML(),
		"CSV":            encodingx.NewCSV(),
		"CSVWithHeaders": encodingx.NewCSVWithHeaders(),
		"Protobuf":       encodingx.NewProtobuf(),
		"Base64":         encodingx.NewBase64(),
		"Base64URL":      encodingx.NewBase64URL(),
		"Binary":         encodingx.NewBinary(),
		"LittleEndian":   encodingx.NewLittleEndian(),
		"BigEndian":      encodingx.NewBigEndian(),
		"Hash":           encodingx.NewHash(),
		"Lazy":           encodingx.NewLazy(),
	}

	for expectedName, enc := range expectedNames {
		actualName := enc.String()
		if actualName != expectedName {
			t.Errorf("Expected encoder name %q, got %q", expectedName, actualName)
		}
	}
}

// TestGroup1_EncoderInterfaceConsistency_StyleValues 测试编码器 Style() 返回的值
// Validates: Requirements 14.2
func TestGroup1_EncoderInterfaceConsistency_StyleValues(t *testing.T) {
	// 结构体风格的编码器
	structStyleEncoders := []encodingx.Encoding{
		encodingx.NewJSON(),
		encodingx.NewYAML(),
		encodingx.NewXML(),
		encodingx.NewCSV(),
		encodingx.NewCSVWithHeaders(),
		encodingx.NewProtobuf(),
		encodingx.NewHash(),
	}

	for _, enc := range structStyleEncoders {
		if enc.Style() != encodingx.EncodingStyleStruct {
			t.Errorf("%s.Style() should return EncodingStyleStruct, got %v", enc.String(), enc.Style())
		}
	}

	// 字节风格的编码器
	bytesStyleEncoders := []encodingx.Encoding{
		encodingx.NewBase64(),
		encodingx.NewBase64URL(),
		encodingx.NewLazy(),
	}

	for _, enc := range bytesStyleEncoders {
		if enc.Style() != encodingx.EncodingStyleBytes {
			t.Errorf("%s.Style() should return EncodingStyleBytes, got %v", enc.String(), enc.Style())
		}
	}

	// Binary 编码器返回 EncodingStyleStruct
	binaryEncoders := []encodingx.Encoding{
		encodingx.NewBinary(),
		encodingx.NewLittleEndian(),
		encodingx.NewBigEndian(),
	}

	for _, enc := range binaryEncoders {
		if enc.Style() != encodingx.EncodingStyleStruct {
			t.Errorf("%s.Style() should return EncodingStyleStruct, got %v", enc.String(), enc.Style())
		}
	}

	// 混合风格的编码器
	mixStyleEncoders := []encodingx.Encoding{
		encodingx.Empty(),
		encodingx.NewChainEncoding([]string{"JSON"}, []string{"JSON"}),
	}

	for _, enc := range mixStyleEncoders {
		if enc.Style() != encodingx.EncodingStyleMix {
			t.Errorf("%s.Style() should return EncodingStyleMix, got %v", enc.String(), enc.Style())
		}
	}
}

// TestGroup1_EncoderInterfaceConsistency_ReverseSymmetry 测试编码器 Reverse() 的对称性
// Validates: Requirements 14.3
func TestGroup1_EncoderInterfaceConsistency_ReverseSymmetry(t *testing.T) {
	// 对于自反编码器（Reverse 返回自身），Reverse().Reverse() 应该等于自身
	selfReverseEncoders := []encodingx.Encoding{
		encodingx.NewJSON(),
		encodingx.NewYAML(),
		encodingx.NewXML(),
		encodingx.NewCSV(),
		encodingx.NewCSVWithHeaders(),
		encodingx.NewProtobuf(),
		encodingx.NewHash(),
		encodingx.NewLazy(),
	}

	for _, enc := range selfReverseEncoders {
		reversed := enc.Reverse()
		doubleReversed := reversed.Reverse()

		// 对于自反编码器，Reverse().Reverse() 应该返回相同类型
		if enc.String() != doubleReversed.String() {
			t.Errorf("%s.Reverse().Reverse().String() should equal %s, got %s",
				enc.String(), enc.String(), doubleReversed.String())
		}
	}

	// 对于 Base64 和 Binary 编码器，Reverse 也应该返回自身
	base64Encoders := []encodingx.Encoding{
		encodingx.NewBase64(),
		encodingx.NewBase64URL(),
		encodingx.NewBinary(),
		encodingx.NewLittleEndian(),
		encodingx.NewBigEndian(),
	}

	for _, enc := range base64Encoders {
		reversed := enc.Reverse()
		if enc.String() != reversed.String() {
			t.Errorf("%s.Reverse().String() should equal %s, got %s",
				enc.String(), enc.String(), reversed.String())
		}
	}
}

// ============================================================================
// Property 20: 编码器接口一致性属性测试
// **Validates: Requirements 14.1, 14.2, 14.3**
// ============================================================================

// TestProperty20_EncoderInterfaceConsistency 属性测试
// **Property 20: 编码器接口一致性**
// *For any* 注册的编码器类型，String() 应该返回类型名称，
// Style() 应该返回有效的 EncodingStyleType，
// Reverse() 应该返回有效的 Encoding 实例。
// **Validates: Requirements 14.1, 14.2, 14.3**
func TestProperty20_EncoderInterfaceConsistency(t *testing.T) {
	// 获取所有已知的编码器
	encoders := []encodingx.Encoding{
		encodingx.NewJSON(),
		encodingx.NewYAML(),
		encodingx.NewXML(),
		encodingx.NewCSV(),
		encodingx.NewCSVWithHeaders(),
		encodingx.NewProtobuf(),
		encodingx.NewBase64(),
		encodingx.NewBase64URL(),
		encodingx.NewBinary(),
		encodingx.NewLittleEndian(),
		encodingx.NewBigEndian(),
		encodingx.NewHash(),
		encodingx.NewLazy(),
		encodingx.Empty(),
		encodingx.NewChainEncoding([]string{"JSON"}, []string{"JSON"}),
		encodingx.NewChainEncoding([]string{"Lazy"}, []string{"Lazy"}),
		encodingx.NewChainEncoding([]string{"JSON", "Base64"}, []string{"Base64", "JSON"}),
	}

	for _, enc := range encoders {
		// Property 20.1: String() 应该返回非空字符串
		name := enc.String()
		if name == "" {
			t.Errorf("Encoder %T: String() should return non-empty string", enc)
		}

		// Property 20.2: Style() 应该返回有效的 EncodingStyleType
		style := enc.Style()
		validStyles := []encodingx.EncodingStyleType{
			encodingx.EncodingStyleStruct,
			encodingx.EncodingStyleBytes,
			encodingx.EncodingStyleMix,
		}
		isValidStyle := false
		for _, vs := range validStyles {
			if style == vs {
				isValidStyle = true
				break
			}
		}
		if !isValidStyle {
			t.Errorf("Encoder %s: Style() returned invalid EncodingStyleType: %v", name, style)
		}

		// Property 20.3: Reverse() 应该返回非 nil 的 Encoding
		reversed := enc.Reverse()
		if reversed == nil {
			t.Errorf("Encoder %s: Reverse() should return non-nil Encoding", name)
			continue
		}

		// Property 20.4: Reverse() 返回的 Encoding 也应该满足接口一致性
		revName := reversed.String()
		if revName == "" {
			t.Errorf("Encoder %s: Reverse().String() should return non-empty string", name)
		}

		revStyle := reversed.Style()
		isValidRevStyle := false
		for _, vs := range validStyles {
			if revStyle == vs {
				isValidRevStyle = true
				break
			}
		}
		if !isValidRevStyle {
			t.Errorf("Encoder %s: Reverse().Style() returned invalid EncodingStyleType: %v", name, revStyle)
		}

		revReversed := reversed.Reverse()
		if revReversed == nil {
			t.Errorf("Encoder %s: Reverse().Reverse() should return non-nil Encoding", name)
		}
	}
}

// TestProperty20_EncoderInterfaceConsistency_Rapid 使用 rapid 进行属性测试
// **Property 20: 编码器接口一致性**
// **Validates: Requirements 14.1, 14.2, 14.3**
func TestProperty20_EncoderInterfaceConsistency_Rapid(t *testing.T) {
	// 定义编码器生成器
	encoderGenerators := []func() encodingx.Encoding{
		func() encodingx.Encoding { return encodingx.NewJSON() },
		func() encodingx.Encoding { return encodingx.NewYAML() },
		func() encodingx.Encoding { return encodingx.NewXML() },
		func() encodingx.Encoding { return encodingx.NewCSV() },
		func() encodingx.Encoding { return encodingx.NewCSVWithHeaders() },
		func() encodingx.Encoding { return encodingx.NewProtobuf() },
		func() encodingx.Encoding { return encodingx.NewBase64() },
		func() encodingx.Encoding { return encodingx.NewBase64URL() },
		func() encodingx.Encoding { return encodingx.NewBinary() },
		func() encodingx.Encoding { return encodingx.NewLittleEndian() },
		func() encodingx.Encoding { return encodingx.NewBigEndian() },
		func() encodingx.Encoding { return encodingx.NewHash() },
		func() encodingx.Encoding { return encodingx.NewLazy() },
		func() encodingx.Encoding { return encodingx.Empty() },
	}

	rapid.Check(t, func(t *rapid.T) {
		// 随机选择一个编码器
		idx := rapid.IntRange(0, len(encoderGenerators)-1).Draw(t, "encoderIndex")
		enc := encoderGenerators[idx]()

		// Property 20.1: String() 应该返回非空字符串
		name := enc.String()
		if name == "" {
			t.Fatalf("Encoder: String() should return non-empty string")
		}

		// Property 20.2: Style() 应该返回有效的 EncodingStyleType
		style := enc.Style()
		if style != encodingx.EncodingStyleStruct &&
			style != encodingx.EncodingStyleBytes &&
			style != encodingx.EncodingStyleMix {
			t.Fatalf("Encoder %s: Style() returned invalid EncodingStyleType: %v", name, style)
		}

		// Property 20.3: Reverse() 应该返回非 nil 的 Encoding
		reversed := enc.Reverse()
		if reversed == nil {
			t.Fatalf("Encoder %s: Reverse() should return non-nil Encoding", name)
		}

		// Property 20.4: Reverse() 返回的 Encoding 也应该满足接口一致性
		revName := reversed.String()
		if revName == "" {
			t.Fatalf("Encoder %s: Reverse().String() should return non-empty string", name)
		}

		revStyle := reversed.Style()
		if revStyle != encodingx.EncodingStyleStruct &&
			revStyle != encodingx.EncodingStyleBytes &&
			revStyle != encodingx.EncodingStyleMix {
			t.Fatalf("Encoder %s: Reverse().Style() returned invalid EncodingStyleType: %v", name, revStyle)
		}

		revReversed := reversed.Reverse()
		if revReversed == nil {
			t.Fatalf("Encoder %s: Reverse().Reverse() should return non-nil Encoding", name)
		}
	})
}
