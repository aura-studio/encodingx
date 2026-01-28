package encodingx_test

import (
	"encoding/xml"
	"testing"

	"github.com/aura-studio/encodingx"
)

// ============================================================================
// XML 编码器单元测试
// Validates: Requirements 3.1, 3.2, 14.1, 14.2, 14.3
// ============================================================================

// XMLTestStructWithName 是带有 XMLName 字段的测试结构体
type XMLTestStructWithName struct {
	XMLName xml.Name `xml:"root"`
	Integer int      `xml:"integer"`
	String  string   `xml:"string"`
	Bool    bool     `xml:"bool"`
	Float   float64  `xml:"float"`
}

// Equal 检查两个 XMLTestStructWithName 是否相等
func (x XMLTestStructWithName) Equal(other XMLTestStructWithName) bool {
	return x.Integer == other.Integer &&
		x.String == other.String &&
		x.Bool == other.Bool &&
		x.Float == other.Float
}

// XMLNestedStruct 是带有 XML 标签的嵌套结构体
type XMLNestedStruct struct {
	XMLName xml.Name      `xml:"nested"`
	Name    string        `xml:"name"`
	Inner   XMLTestStruct `xml:"inner"`
	Items   []int         `xml:"items>item"`
}

// Equal 检查两个 XMLNestedStruct 是否相等
func (x XMLNestedStruct) Equal(other XMLNestedStruct) bool {
	if x.Name != other.Name {
		return false
	}
	if !x.Inner.Equal(other.Inner) {
		return false
	}
	if len(x.Items) != len(other.Items) {
		return false
	}
	for i := range x.Items {
		if x.Items[i] != other.Items[i] {
			return false
		}
	}
	return true
}

// XMLStructWithAttr 是带有 XML 属性的结构体
type XMLStructWithAttr struct {
	XMLName xml.Name `xml:"element"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

// Equal 检查两个 XMLStructWithAttr 是否相等
func (x XMLStructWithAttr) Equal(other XMLStructWithAttr) bool {
	return x.ID == other.ID &&
		x.Name == other.Name &&
		x.Value == other.Value
}

// XMLStructWithCDATA 是带有 CDATA 的结构体
type XMLStructWithCDATA struct {
	XMLName xml.Name `xml:"data"`
	Content string   `xml:",cdata"`
}

// Equal 检查两个 XMLStructWithCDATA 是否相等
func (x XMLStructWithCDATA) Equal(other XMLStructWithCDATA) bool {
	return x.Content == other.Content
}

// ============================================================================
// 基本序列化/反序列化测试
// ============================================================================

// TestXMLMarshalStruct 测试普通结构体序列化
// Validates: Requirements 3.1
func TestXMLMarshalStruct(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{
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

	// 验证返回的是有效的 XML
	var result XMLTestStruct
	if err := xml.Unmarshal(data, &result); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// 验证字段值
	if result.Integer != original.Integer {
		t.Errorf("Integer mismatch: expected %d, got %d", original.Integer, result.Integer)
	}
	if result.String != original.String {
		t.Errorf("String mismatch: expected %s, got %s", original.String, result.String)
	}
	if result.Bool != original.Bool {
		t.Errorf("Bool mismatch: expected %v, got %v", original.Bool, result.Bool)
	}
}

// TestXMLUnmarshalStruct 测试普通结构体反序列化
// Validates: Requirements 3.2
func TestXMLUnmarshalStruct(t *testing.T) {
	encoder := encodingx.NewXML()
	xmlData := []byte(`<XMLTestStruct><integer>42</integer><string>hello world</string><bool>true</bool><float>3.14159</float></XMLTestStruct>`)

	var result XMLTestStruct
	err := encoder.Unmarshal(xmlData, &result)
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

// TestXMLRoundTripStruct 测试结构体序列化后反序列化
// Validates: Requirements 3.1, 3.2
func TestXMLRoundTripStruct(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{
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
	var result XMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// 带 XML 标签的结构体测试
// ============================================================================

// TestXMLMarshalStructWithXMLName 测试带 XMLName 字段的结构体序列化
// Validates: Requirements 3.1
func TestXMLMarshalStructWithXMLName(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStructWithName{
		Integer: 42,
		String:  "with xml name",
		Bool:    true,
		Float:   1.5,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证 XML 根元素名称
	if string(data[:6]) != "<root>" {
		t.Errorf("XML should start with <root>, got: %s", string(data[:20]))
	}

	// 反序列化
	var result XMLTestStructWithName
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestXMLMarshalNestedStruct 测试嵌套结构体序列化/反序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalNestedStruct(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLNestedStruct{
		Name: "nested test",
		Inner: XMLTestStruct{
			Integer: 123,
			String:  "inner struct",
			Bool:    true,
			Float:   1.5,
		},
		Items: []int{1, 2, 3, 4, 5},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLNestedStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Nested struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestXMLMarshalStructWithAttr 测试带 XML 属性的结构体序列化/反序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalStructWithAttr(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLStructWithAttr{
		ID:    42,
		Name:  "test-element",
		Value: "element content",
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证 XML 包含属性
	xmlStr := string(data)
	if !contains(xmlStr, "id=") || !contains(xmlStr, "name=") {
		t.Errorf("XML should contain attributes, got: %s", xmlStr)
	}

	// 反序列化
	var result XMLStructWithAttr
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Struct with attr round trip failed: original %+v != result %+v", original, result)
	}
}

// TestXMLMarshalStructWithCDATA 测试带 CDATA 的结构体序列化/反序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalStructWithCDATA(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLStructWithCDATA{
		Content: "This is <special> content with & characters",
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLStructWithCDATA
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Struct with CDATA round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// String()、Style()、Reverse() 方法测试
// ============================================================================

// TestXMLString 测试 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestXMLString(t *testing.T) {
	encoder := encodingx.NewXML()
	name := encoder.String()

	if name != "XML" {
		t.Errorf("String() should return 'XML', got '%s'", name)
	}
}

// TestXMLStyle 测试 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestXMLStyle(t *testing.T) {
	encoder := encodingx.NewXML()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestXMLReverse 测试 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestXMLReverse(t *testing.T) {
	encoder := encodingx.NewXML()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 XML 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestXMLImplementsEncoding 测试 XML 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestXMLImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewXML()

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
// 边界情况和特殊字符测试
// ============================================================================

// TestXMLMarshalEmptyStruct 测试空结构体序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalEmptyStruct(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Empty struct round trip failed: original %+v != result %+v", original, result)
	}
}

// TestXMLMarshalSpecialCharacters 测试包含 XML 特殊字符的结构体序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalSpecialCharacters(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{
		Integer: -999,
		String:  "hello &amp; world <test>",
		Bool:    true,
		Float:   -0.001,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Special characters round trip failed: original %+v != result %+v", original, result)
	}
}

// TestXMLMarshalUnicodeString 测试包含 Unicode 字符的结构体序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalUnicodeString(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{
		Integer: 123,
		String:  "你好世界",
		Bool:    false,
		Float:   42.0,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLTestStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !original.Equal(result) {
		t.Errorf("Unicode string round trip failed: original %+v != result %+v", original, result)
	}
}

// TestXMLMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{
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
	var result XMLTestStruct
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

// TestXMLMarshalNegativeNumbers 测试负数序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalNegativeNumbers(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLTestStruct{
		Integer: -2147483648, // Min int32
		String:  "negative numbers",
		Bool:    false,
		Float:   -1.5e+100,
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLTestStruct
	err = encoder.Unmarshal(data, &result)
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

// TestXMLUnmarshalInvalidXML 测试无效 XML 反序列化
// Validates: Requirements 3.2
func TestXMLUnmarshalInvalidXML(t *testing.T) {
	encoder := encodingx.NewXML()
	invalidXML := []byte(`<invalid xml`)

	var result XMLTestStruct
	err := encoder.Unmarshal(invalidXML, &result)
	if err == nil {
		t.Error("Unmarshal should fail for invalid XML")
	}
}

// TestXMLUnmarshalMalformedXML 测试格式错误的 XML 反序列化
// Validates: Requirements 3.2
func TestXMLUnmarshalMalformedXML(t *testing.T) {
	encoder := encodingx.NewXML()
	malformedXML := []byte(`<root><unclosed>`)

	var result XMLTestStruct
	err := encoder.Unmarshal(malformedXML, &result)
	if err == nil {
		t.Error("Unmarshal should fail for malformed XML")
	}
}

// TestXMLMarshalSlice 测试切片序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalSlice(t *testing.T) {
	encoder := encodingx.NewXML()

	// XML 切片需要包装在一个根元素中
	type XMLSliceWrapper struct {
		Items []XMLTestStruct `xml:"item"`
	}

	original := XMLSliceWrapper{
		Items: []XMLTestStruct{
			{Integer: 1, String: "first", Bool: true, Float: 1.1},
			{Integer: 2, String: "second", Bool: false, Float: 2.2},
			{Integer: 3, String: "third", Bool: true, Float: 3.3},
		},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLSliceWrapper
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(original.Items) != len(result.Items) {
		t.Fatalf("Slice length mismatch: expected %d, got %d", len(original.Items), len(result.Items))
	}
	for i := range original.Items {
		if !original.Items[i].Equal(result.Items[i]) {
			t.Errorf("Slice element %d mismatch: expected %+v, got %+v", i, original.Items[i], result.Items[i])
		}
	}
}

// TestXMLMarshalEmptySlice 测试空切片序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalEmptySlice(t *testing.T) {
	encoder := encodingx.NewXML()

	type XMLSliceWrapper struct {
		Items []XMLTestStruct `xml:"item"`
	}

	original := XMLSliceWrapper{
		Items: []XMLTestStruct{},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLSliceWrapper
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result.Items) != 0 {
		t.Errorf("Empty slice should remain empty, got %d items", len(result.Items))
	}
}

// TestXMLMarshalEmptyNestedItems 测试空嵌套项序列化
// Validates: Requirements 3.1, 3.2
func TestXMLMarshalEmptyNestedItems(t *testing.T) {
	encoder := encodingx.NewXML()
	original := XMLNestedStruct{
		Name: "empty items",
		Inner: XMLTestStruct{
			Integer: 0,
			String:  "",
			Bool:    false,
			Float:   0,
		},
		Items: []int{},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result XMLNestedStruct
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if result.Name != original.Name {
		t.Errorf("Name mismatch: expected %s, got %s", original.Name, result.Name)
	}
}

// ============================================================================
// 辅助函数
// ============================================================================

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// XML 编码器属性测试
// ============================================================================

// generateXMLNestedStruct 生成一个随机的 XMLNestedStruct
func generateXMLNestedStruct(gen *TestDataGenerator) XMLNestedStruct {
	sliceLen := 1 + gen.GenerateInt32()%10
	if sliceLen < 1 {
		sliceLen = 1
	}
	items := make([]int, sliceLen)
	for j := range items {
		items[j] = int(gen.GenerateInt32() % 1000)
	}
	return XMLNestedStruct{
		Name:  gen.GenerateXMLTestStruct().String,
		Inner: gen.GenerateXMLTestStruct(),
		Items: items,
	}
}

// generateXMLStructWithAttr 生成一个随机的 XMLStructWithAttr
func generateXMLStructWithAttr(gen *TestDataGenerator) XMLStructWithAttr {
	xmlStruct := gen.GenerateXMLTestStruct()
	return XMLStructWithAttr{
		ID:    int(gen.GenerateInt32() % 10000),
		Name:  xmlStruct.String,
		Value: gen.GenerateXMLTestStruct().String,
	}
}

// generateXMLTestStructWithName 生成一个随机的 XMLTestStructWithName
func generateXMLTestStructWithName(gen *TestDataGenerator) XMLTestStructWithName {
	xmlStruct := gen.GenerateXMLTestStruct()
	return XMLTestStructWithName{
		Integer: xmlStruct.Integer,
		String:  xmlStruct.String,
		Bool:    xmlStruct.Bool,
		Float:   xmlStruct.Float,
	}
}

// generateXMLTestStructSlice 生成一个随机的 XMLTestStruct 切片
func generateXMLTestStructSlice(gen *TestDataGenerator, minLen, maxLen int) []XMLTestStruct {
	length := minLen + int(gen.GenerateInt32()%int32(maxLen-minLen+1))
	if length < minLen {
		length = minLen
	}
	result := make([]XMLTestStruct, length)
	for i := range result {
		result[i] = gen.GenerateXMLTestStruct()
	}
	return result
}

// TestGroup1_Property_4_XML_RoundTrip 属性测试：XML Round-Trip 一致性
// **Property 4: XML Round-Trip 一致性**
// *For any* 有效的 Go 结构体（带 XML 标签），使用 XML 编码器序列化后再反序列化，
// 应该产生与原始结构体等价的对象。
// **Validates: Requirements 3.3**
func TestGroup1_Property_4_XML_RoundTrip(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewXML()

	t.Run("XMLTestStruct_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			original := gen.GenerateXMLTestStruct()

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLTestStruct
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

	t.Run("XMLNestedStruct_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机嵌套结构体
			original := generateXMLNestedStruct(gen)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLNestedStruct
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

	t.Run("XMLStructWithAttr_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			original := generateXMLStructWithAttr(gen)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLStructWithAttr
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

	t.Run("XMLSlice_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机切片
			items := generateXMLTestStructSlice(gen, 1, 10)

			// XML 切片需要包装在一个根元素中
			type XMLSliceWrapper struct {
				Items []XMLTestStruct `xml:"item"`
			}
			original := XMLSliceWrapper{Items: items}

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLSliceWrapper
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if len(original.Items) != len(result.Items) {
				t.Errorf("Iteration %d: Slice length mismatch: expected %d, got %d", i, len(original.Items), len(result.Items))
				continue
			}
			for j := range original.Items {
				if !original.Items[j].Equal(result.Items[j]) {
					t.Errorf("Iteration %d: Slice element %d mismatch:\n  original: %+v\n  result:   %+v", i, j, original.Items[j], result.Items[j])
				}
			}
		}
	})

	t.Run("EmptyStruct_RoundTrip", func(t *testing.T) {
		// 测试空结构体的 Round-Trip
		for i := 0; i < iterations; i++ {
			original := XMLTestStruct{}

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLTestStruct
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
		type XMLSliceWrapper struct {
			Items []XMLTestStruct `xml:"item"`
		}

		for i := 0; i < iterations; i++ {
			original := XMLSliceWrapper{Items: []XMLTestStruct{}}

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLSliceWrapper
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性（空切片）
			if len(result.Items) != 0 {
				t.Errorf("Iteration %d: Empty slice round-trip failed: expected empty slice, got %v", i, result.Items)
			}
		}
	})

	t.Run("XMLTestStructWithName_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			original := generateXMLTestStructWithName(gen)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result XMLTestStructWithName
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
}
