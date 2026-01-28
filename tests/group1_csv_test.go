package encodingx_test

import (
	"testing"

	"github.com/aura-studio/encodingx"
)

// ============================================================================
// CSV 编码器单元测试
// Validates: Requirements 4.1, 4.2, 4.3, 4.4, 14.1, 14.2, 14.3
// ============================================================================

// CSVTestRecord 是 CSV 测试专用结构体
type CSVTestRecord struct {
	ID    int     `csv:"id"`
	Name  string  `csv:"name"`
	Value float64 `csv:"value"`
	Flag  bool    `csv:"flag"`
}

// Equal 检查两个 CSVTestRecord 是否相等
func (c CSVTestRecord) Equal(other CSVTestRecord) bool {
	return c.ID == other.ID &&
		c.Name == other.Name &&
		c.Value == other.Value &&
		c.Flag == other.Flag
}

// CSVTestRecordSliceEqual 比较两个 CSVTestRecord 切片是否相等
func CSVTestRecordSliceEqual(a, b []*CSVTestRecord) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(*b[i]) {
			return false
		}
	}
	return true
}

// ============================================================================
// CSV 编码器（无表头）基本序列化/反序列化测试
// ============================================================================

// TestCSVMarshalSlice 测试 CSV 无表头序列化结构体切片
// Validates: Requirements 4.1
func TestCSVMarshalSlice(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{
		{ID: 1, Name: "Alice", Value: 100.5, Flag: true},
		{ID: 2, Name: "Bob", Value: 200.75, Flag: false},
		{ID: 3, Name: "Charlie", Value: 300.25, Flag: true},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是有效的 CSV（不带表头）
	csvStr := string(data)
	// CSV 无表头应该直接是数据行
	if len(csvStr) == 0 {
		t.Error("Marshal should return non-empty CSV data")
	}

	// 验证不包含表头（表头应该是 id,name,value,flag）
	// 无表头的 CSV 第一行应该是数据，不是字段名
	if csvStr[:2] == "id" {
		t.Error("CSV without headers should not start with header row")
	}
}

// TestCSVUnmarshalSlice 测试 CSV 无表头反序列化
// Validates: Requirements 4.2
func TestCSVUnmarshalSlice(t *testing.T) {
	encoder := encodingx.NewCSV()
	// CSV 数据（无表头）
	csvData := []byte("1,Alice,100.5,true\n2,Bob,200.75,false\n")

	var result []*CSVTestRecord
	err := encoder.Unmarshal(csvData, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(result))
	}

	if result[0].ID != 1 || result[0].Name != "Alice" {
		t.Errorf("First record mismatch: got %+v", result[0])
	}
	if result[1].ID != 2 || result[1].Name != "Bob" {
		t.Errorf("Second record mismatch: got %+v", result[1])
	}
}

// TestCSVRoundTripSlice 测试 CSV 无表头序列化后反序列化
// Validates: Requirements 4.1, 4.2
func TestCSVRoundTripSlice(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{
		{ID: 10, Name: "Test1", Value: 1.1, Flag: true},
		{ID: 20, Name: "Test2", Value: 2.2, Flag: false},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !CSVTestRecordSliceEqual(original, result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// CSVWithHeaders 编码器（带表头）基本序列化/反序列化测试
// ============================================================================

// TestCSVWithHeadersMarshalSlice 测试 CSVWithHeaders 带表头序列化结构体切片
// Validates: Requirements 4.3
func TestCSVWithHeadersMarshalSlice(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	original := []*CSVTestRecord{
		{ID: 1, Name: "Alice", Value: 100.5, Flag: true},
		{ID: 2, Name: "Bob", Value: 200.75, Flag: false},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证返回的是有效的 CSV（带表头）
	csvStr := string(data)
	if len(csvStr) == 0 {
		t.Error("Marshal should return non-empty CSV data")
	}

	// 验证包含表头
	if csvStr[:2] != "id" {
		t.Errorf("CSV with headers should start with header row, got: %s", csvStr[:20])
	}
}

// TestCSVWithHeadersUnmarshalSlice 测试 CSVWithHeaders 带表头反序列化
// Validates: Requirements 4.4
func TestCSVWithHeadersUnmarshalSlice(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	// CSV 数据（带表头）
	csvData := []byte("id,name,value,flag\n1,Alice,100.5,true\n2,Bob,200.75,false\n")

	var result []*CSVTestRecord
	err := encoder.Unmarshal(csvData, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(result))
	}

	if result[0].ID != 1 || result[0].Name != "Alice" {
		t.Errorf("First record mismatch: got %+v", result[0])
	}
	if result[1].ID != 2 || result[1].Name != "Bob" {
		t.Errorf("Second record mismatch: got %+v", result[1])
	}
}

// TestCSVWithHeadersRoundTripSlice 测试 CSVWithHeaders 带表头序列化后反序列化
// Validates: Requirements 4.3, 4.4
func TestCSVWithHeadersRoundTripSlice(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	original := []*CSVTestRecord{
		{ID: 100, Name: "HeaderTest1", Value: 11.11, Flag: true},
		{ID: 200, Name: "HeaderTest2", Value: 22.22, Flag: false},
		{ID: 300, Name: "HeaderTest3", Value: 33.33, Flag: true},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !CSVTestRecordSliceEqual(original, result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// CSV String()、Style()、Reverse() 方法测试
// ============================================================================

// TestCSVString 测试 CSV 编码器 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestCSVString(t *testing.T) {
	encoder := encodingx.NewCSV()
	name := encoder.String()

	if name != "CSV" {
		t.Errorf("String() should return 'CSV', got '%s'", name)
	}
}

// TestCSVStyle 测试 CSV 编码器 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestCSVStyle(t *testing.T) {
	encoder := encodingx.NewCSV()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestCSVReverse 测试 CSV 编码器 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestCSVReverse(t *testing.T) {
	encoder := encodingx.NewCSV()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 CSV 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestCSVImplementsEncoding 测试 CSV 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestCSVImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewCSV()

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
// CSVWithHeaders String()、Style()、Reverse() 方法测试
// ============================================================================

// TestCSVWithHeadersString 测试 CSVWithHeaders 编码器 String() 方法返回类型名称
// Validates: Requirements 14.1
func TestCSVWithHeadersString(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	name := encoder.String()

	if name != "CSVWithHeaders" {
		t.Errorf("String() should return 'CSVWithHeaders', got '%s'", name)
	}
}

// TestCSVWithHeadersStyle 测试 CSVWithHeaders 编码器 Style() 方法返回 EncodingStyleStruct
// Validates: Requirements 14.2
func TestCSVWithHeadersStyle(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	style := encoder.Style()

	if style != encodingx.EncodingStyleStruct {
		t.Errorf("Style() should return EncodingStyleStruct, got %v", style)
	}
}

// TestCSVWithHeadersReverse 测试 CSVWithHeaders 编码器 Reverse() 方法返回自身
// Validates: Requirements 14.3
func TestCSVWithHeadersReverse(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	reversed := encoder.Reverse()

	// Reverse() 应该返回自身
	if reversed.String() != encoder.String() {
		t.Errorf("Reverse() should return self, got different encoder: %s", reversed.String())
	}

	// 验证 reversed 也是 CSVWithHeaders 编码器
	if reversed.Style() != encodingx.EncodingStyleStruct {
		t.Errorf("Reversed encoder should have same style")
	}
}

// TestCSVWithHeadersImplementsEncoding 测试 CSVWithHeaders 编码器实现 Encoding 接口
// Validates: Requirements 14.1, 14.2, 14.3
func TestCSVWithHeadersImplementsEncoding(t *testing.T) {
	var encoder encodingx.Encoding = encodingx.NewCSVWithHeaders()

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

// TestCSVMarshalEmptySlice 测试空切片序列化
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalEmptySlice(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 空切片应该返回空数据
	if len(data) != 0 {
		t.Errorf("Empty slice should marshal to empty data, got: %s", string(data))
	}
}

// TestCSVWithHeadersMarshalEmptySlice 测试 CSVWithHeaders 空切片序列化
// Validates: Requirements 4.3, 4.4
func TestCSVWithHeadersMarshalEmptySlice(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	original := []*CSVTestRecord{}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 带表头的空切片应该只有表头行
	csvStr := string(data)
	if len(csvStr) == 0 {
		// 空切片可能返回空数据或只有表头
		return
	}
	// 如果有数据，应该是表头
	if csvStr[:2] != "id" {
		t.Errorf("CSV with headers empty slice should have header row, got: %s", csvStr)
	}
}

// TestCSVMarshalSingleRecord 测试单条记录序列化
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalSingleRecord(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{
		{ID: 42, Name: "SingleRecord", Value: 99.99, Flag: true},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}
	if !original[0].Equal(*result[0]) {
		t.Errorf("Single record round trip failed: original %+v != result %+v", original[0], result[0])
	}
}

// TestCSVWithHeadersMarshalSingleRecord 测试 CSVWithHeaders 单条记录序列化
// Validates: Requirements 4.3, 4.4
func TestCSVWithHeadersMarshalSingleRecord(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	original := []*CSVTestRecord{
		{ID: 42, Name: "SingleRecord", Value: 99.99, Flag: true},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}
	if !original[0].Equal(*result[0]) {
		t.Errorf("Single record round trip failed: original %+v != result %+v", original[0], result[0])
	}
}

// TestCSVMarshalLargeNumbers 测试大数值序列化
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalLargeNumbers(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{
		{ID: 2147483647, Name: "MaxInt", Value: 1.7976931348623157e+100, Flag: true},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}
	if result[0].ID != original[0].ID {
		t.Errorf("Large integer mismatch: expected %d, got %d", original[0].ID, result[0].ID)
	}
}

// TestCSVMarshalNegativeNumbers 测试负数序列化
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalNegativeNumbers(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{
		{ID: -999, Name: "Negative", Value: -123.456, Flag: false},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}
	if result[0].ID != original[0].ID {
		t.Errorf("Negative integer mismatch: expected %d, got %d", original[0].ID, result[0].ID)
	}
	if result[0].Value != original[0].Value {
		t.Errorf("Negative float mismatch: expected %f, got %f", original[0].Value, result[0].Value)
	}
}

// TestCSVMarshalZeroValues 测试零值序列化
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalZeroValues(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVTestRecord{
		{ID: 0, Name: "", Value: 0.0, Flag: false},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}
	if !original[0].Equal(*result[0]) {
		t.Errorf("Zero values round trip failed: original %+v != result %+v", original[0], result[0])
	}
}

// TestCSVMarshalManyRecords 测试多条记录序列化
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalManyRecords(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := make([]*CSVTestRecord, 100)
	for i := 0; i < 100; i++ {
		original[i] = &CSVTestRecord{
			ID:    i,
			Name:  "Record",
			Value: float64(i) * 1.5,
			Flag:  i%2 == 0,
		}
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 100 {
		t.Fatalf("Expected 100 records, got %d", len(result))
	}
	for i := 0; i < 100; i++ {
		if !original[i].Equal(*result[i]) {
			t.Errorf("Record %d mismatch: original %+v != result %+v", i, original[i], result[i])
		}
	}
}

// TestCSVWithHeadersMarshalManyRecords 测试 CSVWithHeaders 多条记录序列化
// Validates: Requirements 4.3, 4.4
func TestCSVWithHeadersMarshalManyRecords(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	original := make([]*CSVTestRecord, 50)
	for i := 0; i < 50; i++ {
		original[i] = &CSVTestRecord{
			ID:    i * 10,
			Name:  "HeaderRecord",
			Value: float64(i) * 2.5,
			Flag:  i%3 == 0,
		}
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVTestRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if len(result) != 50 {
		t.Fatalf("Expected 50 records, got %d", len(result))
	}
	for i := 0; i < 50; i++ {
		if !original[i].Equal(*result[i]) {
			t.Errorf("Record %d mismatch: original %+v != result %+v", i, original[i], result[i])
		}
	}
}

// ============================================================================
// CSV 和 CSVWithHeaders 输出格式对比测试
// ============================================================================

// TestCSVVsCSVWithHeadersOutput 测试 CSV 和 CSVWithHeaders 输出格式差异
// Validates: Requirements 4.1, 4.3
func TestCSVVsCSVWithHeadersOutput(t *testing.T) {
	csvEncoder := encodingx.NewCSV()
	csvWithHeadersEncoder := encodingx.NewCSVWithHeaders()

	records := []*CSVTestRecord{
		{ID: 1, Name: "Test", Value: 1.5, Flag: true},
	}

	// CSV 无表头序列化
	csvData, err := csvEncoder.Marshal(records)
	if err != nil {
		t.Fatalf("CSV Marshal failed: %v", err)
	}

	// CSVWithHeaders 带表头序列化
	csvWithHeadersData, err := csvWithHeadersEncoder.Marshal(records)
	if err != nil {
		t.Fatalf("CSVWithHeaders Marshal failed: %v", err)
	}

	// CSVWithHeaders 输出应该比 CSV 长（因为有表头行）
	if len(csvWithHeadersData) <= len(csvData) {
		t.Errorf("CSVWithHeaders output should be longer than CSV output")
	}

	// 验证 CSVWithHeaders 包含表头
	csvWithHeadersStr := string(csvWithHeadersData)
	if csvWithHeadersStr[:2] != "id" {
		t.Errorf("CSVWithHeaders should start with header, got: %s", csvWithHeadersStr[:20])
	}
}

// ============================================================================
// 使用 CSVRecord 测试（来自 test_helpers_test.go）
// ============================================================================

// TestCSVMarshalCSVRecord 测试使用 CSVRecord 结构体
// Validates: Requirements 4.1, 4.2
func TestCSVMarshalCSVRecord(t *testing.T) {
	encoder := encodingx.NewCSV()
	original := []*CSVRecord{
		{ID: 1, Name: "Alice", Value: 100.5},
		{ID: 2, Name: "Bob", Value: 200.75},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !CSVRecordsEqual(original, result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// TestCSVWithHeadersMarshalCSVRecord 测试使用 CSVRecord 结构体（带表头）
// Validates: Requirements 4.3, 4.4
func TestCSVWithHeadersMarshalCSVRecord(t *testing.T) {
	encoder := encodingx.NewCSVWithHeaders()
	original := []*CSVRecord{
		{ID: 1, Name: "Alice", Value: 100.5},
		{ID: 2, Name: "Bob", Value: 200.75},
	}

	// 序列化
	data, err := encoder.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var result []*CSVRecord
	err = encoder.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证
	if !CSVRecordsEqual(original, result) {
		t.Errorf("Round trip failed: original %+v != result %+v", original, result)
	}
}

// ============================================================================
// CSV 编码器属性测试
// ============================================================================

// generateCSVTestRecords 生成一个随机的 CSVTestRecord 切片
func generateCSVTestRecords(gen *TestDataGenerator, minLen, maxLen int) []*CSVTestRecord {
	length := 1 + int(gen.GenerateInt32()%int32(maxLen-minLen+1))
	if length < minLen {
		length = minLen
	}
	if length > maxLen {
		length = maxLen
	}
	records := make([]*CSVTestRecord, length)
	for i := range records {
		record := gen.GenerateCSVRecord()
		records[i] = &CSVTestRecord{
			ID:    record.ID,
			Name:  record.Name,
			Value: record.Value,
			Flag:  gen.GenerateInt32()%2 == 0,
		}
	}
	return records
}

// TestGroup1_Property_5_CSV_RoundTrip 属性测试：CSV Round-Trip 一致性
// **Property 5: CSV Round-Trip 一致性**
// *For any* 有效的结构体切片，使用 CSV 编码器序列化后再反序列化，
// 应该产生与原始切片等价的对象。
// **Validates: Requirements 4.5**
func TestGroup1_Property_5_CSV_RoundTrip(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewCSV()

	t.Run("CSVTestRecord_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机的 CSVTestRecord 切片（至少1条记录）
			original := generateCSVTestRecords(gen, 1, 10)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("CSVRecord_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 使用 test_helpers 中的 GenerateCSVRecords
			original := gen.GenerateCSVRecords(1, 10)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVRecordsEqual(original, result) {
				t.Errorf("Iteration %d: Round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("SingleRecord_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 测试单条记录的 Round-Trip
			original := generateCSVTestRecords(gen, 1, 1)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Single record round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("ManyRecords_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations/10; i++ { // 减少迭代次数因为每次生成更多记录
			// 测试多条记录的 Round-Trip
			original := generateCSVTestRecords(gen, 10, 50)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Many records round-trip failed: length mismatch or content mismatch", i)
			}
		}
	})
}

// TestGroup1_Property_6_CSVWithHeaders_RoundTrip 属性测试：CSVWithHeaders Round-Trip 一致性
// **Property 6: CSVWithHeaders Round-Trip 一致性**
// *For any* 有效的结构体切片，使用 CSVWithHeaders 编码器序列化后再反序列化，
// 应该产生与原始切片等价的对象。
// **Validates: Requirements 4.6**
func TestGroup1_Property_6_CSVWithHeaders_RoundTrip(t *testing.T) {
	const iterations = 100
	encoder := encodingx.NewCSVWithHeaders()

	t.Run("CSVTestRecord_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机的 CSVTestRecord 切片（至少1条记录）
			original := generateCSVTestRecords(gen, 1, 10)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("CSVRecord_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 使用 test_helpers 中的 GenerateCSVRecords
			original := gen.GenerateCSVRecords(1, 10)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVRecordsEqual(original, result) {
				t.Errorf("Iteration %d: Round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("SingleRecord_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 测试单条记录的 Round-Trip
			original := generateCSVTestRecords(gen, 1, 1)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Single record round-trip failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})

	t.Run("ManyRecords_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations/10; i++ { // 减少迭代次数因为每次生成更多记录
			// 测试多条记录的 Round-Trip
			original := generateCSVTestRecords(gen, 10, 50)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Many records round-trip failed: length mismatch or content mismatch", i)
			}
		}
	})

	t.Run("HeadersPresent_RoundTrip", func(t *testing.T) {
		gen := NewTestDataGenerator()
		for i := 0; i < iterations; i++ {
			// 生成随机记录
			original := generateCSVTestRecords(gen, 1, 5)

			// 序列化
			data, err := encoder.Marshal(original)
			if err != nil {
				t.Fatalf("Iteration %d: Marshal failed: %v", i, err)
			}

			// 验证输出包含表头
			csvStr := string(data)
			if len(csvStr) > 0 && csvStr[:2] != "id" {
				t.Errorf("Iteration %d: CSVWithHeaders should start with header row, got: %s", i, csvStr[:min(20, len(csvStr))])
			}

			// 反序列化
			var result []*CSVTestRecord
			err = encoder.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("Iteration %d: Unmarshal failed: %v", i, err)
			}

			// 验证 Round-Trip 一致性
			if !CSVTestRecordSliceEqual(original, result) {
				t.Errorf("Iteration %d: Round-trip with headers failed:\n  original: %+v\n  result:   %+v", i, original, result)
			}
		}
	})
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
