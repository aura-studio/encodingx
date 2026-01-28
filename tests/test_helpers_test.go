package encodingx_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/aura-studio/encodingx"
)

// ============================================================================
// 测试辅助结构体
// ============================================================================

// TestStruct 是基本测试结构体，包含常见的基本类型字段
// 用于 JSON、YAML、XML、CSV 等编码器的测试
// Validates: Requirements 1.1, 2.1, 3.1, 4.1
type TestStruct struct {
	Integer int     `json:"integer" yaml:"integer" xml:"integer" csv:"integer"`
	String  string  `json:"string" yaml:"string" xml:"string" csv:"string"`
	Bool    bool    `json:"bool" yaml:"bool" xml:"bool" csv:"bool"`
	Float   float64 `json:"float" yaml:"float" xml:"float" csv:"float"`
}

// Equal 检查两个 TestStruct 是否相等
func (t TestStruct) Equal(other TestStruct) bool {
	return t.Integer == other.Integer &&
		t.String == other.String &&
		t.Bool == other.Bool &&
		t.Float == other.Float
}

// NestedStruct 是嵌套结构体测试类型
// 用于测试编码器对嵌套结构的处理能力
type NestedStruct struct {
	Name  string     `json:"name" yaml:"name" xml:"name"`
	Inner TestStruct `json:"inner" yaml:"inner" xml:"inner"`
	Slice []int      `json:"slice" yaml:"slice" xml:"slice"`
}

// Equal 检查两个 NestedStruct 是否相等
func (n NestedStruct) Equal(other NestedStruct) bool {
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

// CSVRecord 是 CSV 专用测试结构体
// 用于 CSV 和 CSVWithHeaders 编码器的测试
// Validates: Requirements 4.1
type CSVRecord struct {
	ID    int     `csv:"id"`
	Name  string  `csv:"name"`
	Value float64 `csv:"value"`
}

// Equal 检查两个 CSVRecord 是否相等
func (c CSVRecord) Equal(other CSVRecord) bool {
	return c.ID == other.ID &&
		c.Name == other.Name &&
		c.Value == other.Value
}

// XMLTestStruct 是 XML 专用测试结构体，带有 XMLName 字段
type XMLTestStruct struct {
	Integer int     `xml:"integer"`
	String  string  `xml:"string"`
	Bool    bool    `xml:"bool"`
	Float   float64 `xml:"float"`
}

// Equal 检查两个 XMLTestStruct 是否相等
func (x XMLTestStruct) Equal(other XMLTestStruct) bool {
	return x.Integer == other.Integer &&
		x.String == other.String &&
		x.Bool == other.Bool &&
		x.Float == other.Float
}

// HashableStruct 实现 HashMarshaller 和 HashUnmarshaller 接口
// 用于 Hash 编码器的测试
type HashableStruct struct {
	pairs [][]interface{}
}

// MarshalHash 实现 HashMarshaller 接口
func (h *HashableStruct) MarshalHash() [][]interface{} {
	return h.pairs
}

// UnmarshalHash 实现 HashUnmarshaller 接口
func (h *HashableStruct) UnmarshalHash(pairs [][]interface{}) {
	h.pairs = pairs
}

// GetPairs 返回内部的 pairs 数据
func (h *HashableStruct) GetPairs() [][]interface{} {
	return h.pairs
}

// SetPairs 设置内部的 pairs 数据
func (h *HashableStruct) SetPairs(pairs [][]interface{}) {
	h.pairs = pairs
}

// NewHashableStruct 创建一个新的 HashableStruct
func NewHashableStruct(pairs [][]interface{}) *HashableStruct {
	return &HashableStruct{pairs: pairs}
}

// ============================================================================
// 随机数据生成器
// ============================================================================

// TestDataGenerator 是测试数据生成器
// 用于属性测试中生成随机测试数据
type TestDataGenerator struct {
	rng *rand.Rand
}

// NewTestDataGenerator 创建一个新的测试数据生成器
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewTestDataGeneratorWithSeed 创建一个带有指定种子的测试数据生成器
// 用于可重现的测试
func NewTestDataGeneratorWithSeed(seed int64) *TestDataGenerator {
	return &TestDataGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// GenerateTestStruct 生成一个随机的 TestStruct
func (g *TestDataGenerator) GenerateTestStruct() TestStruct {
	return TestStruct{
		Integer: g.rng.Intn(10000) - 5000, // -5000 到 4999
		String:  g.generateString(1, 50),
		Bool:    g.rng.Intn(2) == 1,
		Float:   g.rng.Float64() * 1000,
	}
}

// GenerateNestedStruct 生成一个随机的 NestedStruct
func (g *TestDataGenerator) GenerateNestedStruct() NestedStruct {
	sliceLen := g.rng.Intn(10) + 1
	slice := make([]int, sliceLen)
	for i := range slice {
		slice[i] = g.rng.Intn(1000)
	}
	return NestedStruct{
		Name:  g.generateString(1, 30),
		Inner: g.GenerateTestStruct(),
		Slice: slice,
	}
}

// GenerateCSVRecord 生成一个随机的 CSVRecord
func (g *TestDataGenerator) GenerateCSVRecord() CSVRecord {
	return CSVRecord{
		ID:    g.rng.Intn(10000),
		Name:  g.generateAlphanumericString(1, 20), // CSV 使用字母数字字符串避免特殊字符问题
		Value: g.rng.Float64() * 1000,
	}
}

// GenerateCSVRecords 生成一个随机的 CSVRecord 切片
func (g *TestDataGenerator) GenerateCSVRecords(minLen, maxLen int) []*CSVRecord {
	length := minLen + g.rng.Intn(maxLen-minLen+1)
	records := make([]*CSVRecord, length)
	for i := range records {
		record := g.GenerateCSVRecord()
		records[i] = &record
	}
	return records
}

// GenerateXMLTestStruct 生成一个随机的 XMLTestStruct
func (g *TestDataGenerator) GenerateXMLTestStruct() XMLTestStruct {
	return XMLTestStruct{
		Integer: g.rng.Intn(10000) - 5000,
		String:  g.generateAlphanumericString(1, 50), // XML 使用字母数字字符串避免特殊字符问题
		Bool:    g.rng.Intn(2) == 1,
		Float:   g.rng.Float64() * 1000,
	}
}

// GenerateBytes 生成指定长度范围的随机字节数组
func (g *TestDataGenerator) GenerateBytes(minLen, maxLen int) []byte {
	length := minLen + g.rng.Intn(maxLen-minLen+1)
	data := make([]byte, length)
	g.rng.Read(data)
	return data
}

// GenerateEncodingxBytes 生成一个随机的 encodingx.Bytes
func (g *TestDataGenerator) GenerateEncodingxBytes(minLen, maxLen int) encodingx.Bytes {
	return encodingx.MakeBytes(g.GenerateBytes(minLen, maxLen))
}

// GenerateHashPairs 生成随机的 Hash 键值对
func (g *TestDataGenerator) GenerateHashPairs(minPairs, maxPairs int) [][]interface{} {
	numPairs := minPairs + g.rng.Intn(maxPairs-minPairs+1)
	pairs := make([][]interface{}, numPairs)
	for i := range pairs {
		key := g.generateAlphanumericString(1, 10)
		value := g.generateAlphanumericString(1, 20)
		pairs[i] = []interface{}{key, value}
	}
	return pairs
}

// GenerateInt32 生成一个随机的 int32
func (g *TestDataGenerator) GenerateInt32() int32 {
	return int32(g.rng.Int31())
}

// GenerateInt64 生成一个随机的 int64
func (g *TestDataGenerator) GenerateInt64() int64 {
	return g.rng.Int63()
}

// GenerateFloat32 生成一个随机的 float32
func (g *TestDataGenerator) GenerateFloat32() float32 {
	return g.rng.Float32()
}

// GenerateFloat64 生成一个随机的 float64
func (g *TestDataGenerator) GenerateFloat64() float64 {
	return g.rng.Float64()
}

// GenerateUint32 生成一个随机的 uint32
func (g *TestDataGenerator) GenerateUint32() uint32 {
	return g.rng.Uint32()
}

// GenerateUint64 生成一个随机的 uint64
func (g *TestDataGenerator) GenerateUint64() uint64 {
	return g.rng.Uint64()
}

// GenerateStringIntMap 生成一个随机的 map[string]int
func (g *TestDataGenerator) GenerateStringIntMap(minKeys, maxKeys int) map[string]int {
	numKeys := minKeys + g.rng.Intn(maxKeys-minKeys+1)
	result := make(map[string]int)
	for i := 0; i < numKeys; i++ {
		key := g.generateAlphanumericString(1, 20)
		value := g.rng.Intn(10000)
		result[key] = value
	}
	return result
}

// GenerateTestStructSlice 生成一个随机的 TestStruct 切片
func (g *TestDataGenerator) GenerateTestStructSlice(minLen, maxLen int) []TestStruct {
	length := minLen + g.rng.Intn(maxLen-minLen+1)
	result := make([]TestStruct, length)
	for i := range result {
		result[i] = g.GenerateTestStruct()
	}
	return result
}

// ============================================================================
// 内部辅助方法
// ============================================================================

// generateString 生成指定长度范围的随机字符串
func (g *TestDataGenerator) generateString(minLen, maxLen int) string {
	length := minLen + g.rng.Intn(maxLen-minLen+1)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[g.rng.Intn(len(charset))]
	}
	return string(result)
}

// generateAlphanumericString 生成只包含字母和数字的随机字符串
// 用于 CSV 和 XML 测试，避免特殊字符导致的解析问题
func (g *TestDataGenerator) generateAlphanumericString(minLen, maxLen int) string {
	length := minLen + g.rng.Intn(maxLen-minLen+1)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[g.rng.Intn(len(charset))]
	}
	return string(result)
}

// ============================================================================
// 测试辅助函数
// ============================================================================

// BytesEqual 比较两个字节数组是否相等
func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// CSVRecordsEqual 比较两个 CSVRecord 切片是否相等
func CSVRecordsEqual(a, b []*CSVRecord) bool {
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

// HashPairsEqual 比较两个 Hash 键值对切片是否相等
func HashPairsEqual(a, b [][]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			// 使用字符串比较，因为 JSON 反序列化后类型可能不同
			aStr, aOk := a[i][j].(string)
			bStr, bOk := b[i][j].(string)
			if aOk && bOk {
				if aStr != bStr {
					return false
				}
			} else {
				// 对于非字符串类型，使用简单的相等比较
				if a[i][j] != b[i][j] {
					return false
				}
			}
		}
	}
	return true
}

// ============================================================================
// 固定大小类型测试结构体（用于 Binary 编码器测试）
// ============================================================================

// FixedSizeStruct 是用于 Binary 编码器测试的固定大小结构体
type FixedSizeStruct struct {
	Int32Val   int32
	Int64Val   int64
	Float32Val float32
	Float64Val float64
}

// Equal 检查两个 FixedSizeStruct 是否相等
func (f FixedSizeStruct) Equal(other FixedSizeStruct) bool {
	return f.Int32Val == other.Int32Val &&
		f.Int64Val == other.Int64Val &&
		f.Float32Val == other.Float32Val &&
		f.Float64Val == other.Float64Val
}

// GenerateFixedSizeStruct 生成一个随机的 FixedSizeStruct
func (g *TestDataGenerator) GenerateFixedSizeStruct() FixedSizeStruct {
	return FixedSizeStruct{
		Int32Val:   g.GenerateInt32(),
		Int64Val:   g.GenerateInt64(),
		Float32Val: g.GenerateFloat32(),
		Float64Val: g.GenerateFloat64(),
	}
}

// ============================================================================
// 测试辅助结构体和生成器的自测试
// ============================================================================

// TestTestDataGenerator 验证测试数据生成器的基本功能
func TestTestDataGenerator(t *testing.T) {
	gen := NewTestDataGenerator()

	t.Run("GenerateTestStruct", func(t *testing.T) {
		ts := gen.GenerateTestStruct()
		// 验证生成的结构体字段有值
		if ts.String == "" {
			t.Error("GenerateTestStruct should generate non-empty string")
		}
	})

	t.Run("GenerateNestedStruct", func(t *testing.T) {
		ns := gen.GenerateNestedStruct()
		if ns.Name == "" {
			t.Error("GenerateNestedStruct should generate non-empty name")
		}
		if len(ns.Slice) == 0 {
			t.Error("GenerateNestedStruct should generate non-empty slice")
		}
	})

	t.Run("GenerateCSVRecord", func(t *testing.T) {
		record := gen.GenerateCSVRecord()
		if record.Name == "" {
			t.Error("GenerateCSVRecord should generate non-empty name")
		}
	})

	t.Run("GenerateCSVRecords", func(t *testing.T) {
		records := gen.GenerateCSVRecords(3, 5)
		if len(records) < 3 || len(records) > 5 {
			t.Errorf("GenerateCSVRecords should generate 3-5 records, got %d", len(records))
		}
	})

	t.Run("GenerateBytes", func(t *testing.T) {
		data := gen.GenerateBytes(10, 20)
		if len(data) < 10 || len(data) > 20 {
			t.Errorf("GenerateBytes should generate 10-20 bytes, got %d", len(data))
		}
	})

	t.Run("GenerateEncodingxBytes", func(t *testing.T) {
		bytes := gen.GenerateEncodingxBytes(5, 15)
		if len(bytes.Data) < 5 || len(bytes.Data) > 15 {
			t.Errorf("GenerateEncodingxBytes should generate 5-15 bytes, got %d", len(bytes.Data))
		}
	})

	t.Run("GenerateHashPairs", func(t *testing.T) {
		pairs := gen.GenerateHashPairs(2, 4)
		if len(pairs) < 2 || len(pairs) > 4 {
			t.Errorf("GenerateHashPairs should generate 2-4 pairs, got %d", len(pairs))
		}
		for _, pair := range pairs {
			if len(pair) != 2 {
				t.Error("Each hash pair should have exactly 2 elements")
			}
		}
	})

	t.Run("GenerateFixedSizeStruct", func(t *testing.T) {
		fs := gen.GenerateFixedSizeStruct()
		// 验证生成的结构体（这里只是确保不会 panic）
		_ = fs.Int32Val
		_ = fs.Int64Val
		_ = fs.Float32Val
		_ = fs.Float64Val
	})
}

// TestTestDataGeneratorWithSeed 验证带种子的生成器可重现性
func TestTestDataGeneratorWithSeed(t *testing.T) {
	seed := int64(12345)
	gen1 := NewTestDataGeneratorWithSeed(seed)
	gen2 := NewTestDataGeneratorWithSeed(seed)

	ts1 := gen1.GenerateTestStruct()
	ts2 := gen2.GenerateTestStruct()

	if !ts1.Equal(ts2) {
		t.Error("Same seed should produce same TestStruct")
	}
}

// TestStructEquality 验证结构体的 Equal 方法
func TestStructEquality(t *testing.T) {
	t.Run("TestStruct.Equal", func(t *testing.T) {
		ts1 := TestStruct{Integer: 1, String: "test", Bool: true, Float: 1.5}
		ts2 := TestStruct{Integer: 1, String: "test", Bool: true, Float: 1.5}
		ts3 := TestStruct{Integer: 2, String: "test", Bool: true, Float: 1.5}

		if !ts1.Equal(ts2) {
			t.Error("Equal TestStructs should be equal")
		}
		if ts1.Equal(ts3) {
			t.Error("Different TestStructs should not be equal")
		}
	})

	t.Run("NestedStruct.Equal", func(t *testing.T) {
		ns1 := NestedStruct{
			Name:  "test",
			Inner: TestStruct{Integer: 1, String: "inner", Bool: true, Float: 1.5},
			Slice: []int{1, 2, 3},
		}
		ns2 := NestedStruct{
			Name:  "test",
			Inner: TestStruct{Integer: 1, String: "inner", Bool: true, Float: 1.5},
			Slice: []int{1, 2, 3},
		}
		ns3 := NestedStruct{
			Name:  "different",
			Inner: TestStruct{Integer: 1, String: "inner", Bool: true, Float: 1.5},
			Slice: []int{1, 2, 3},
		}

		if !ns1.Equal(ns2) {
			t.Error("Equal NestedStructs should be equal")
		}
		if ns1.Equal(ns3) {
			t.Error("Different NestedStructs should not be equal")
		}
	})

	t.Run("CSVRecord.Equal", func(t *testing.T) {
		r1 := CSVRecord{ID: 1, Name: "test", Value: 1.5}
		r2 := CSVRecord{ID: 1, Name: "test", Value: 1.5}
		r3 := CSVRecord{ID: 2, Name: "test", Value: 1.5}

		if !r1.Equal(r2) {
			t.Error("Equal CSVRecords should be equal")
		}
		if r1.Equal(r3) {
			t.Error("Different CSVRecords should not be equal")
		}
	})
}

// TestHashableStruct 验证 HashableStruct 的功能
func TestHashableStruct(t *testing.T) {
	pairs := [][]interface{}{
		{"key1", "value1"},
		{"key2", "value2"},
	}

	h := NewHashableStruct(pairs)

	// 测试 MarshalHash
	marshaled := h.MarshalHash()
	if len(marshaled) != 2 {
		t.Errorf("MarshalHash should return 2 pairs, got %d", len(marshaled))
	}

	// 测试 UnmarshalHash
	newPairs := [][]interface{}{
		{"newKey", "newValue"},
	}
	h.UnmarshalHash(newPairs)
	if len(h.GetPairs()) != 1 {
		t.Errorf("UnmarshalHash should set 1 pair, got %d", len(h.GetPairs()))
	}

	// 测试 SetPairs
	h.SetPairs(pairs)
	if len(h.GetPairs()) != 2 {
		t.Errorf("SetPairs should set 2 pairs, got %d", len(h.GetPairs()))
	}
}

// TestBytesEqual 验证 BytesEqual 辅助函数
func TestBytesEqual(t *testing.T) {
	a := []byte{1, 2, 3}
	b := []byte{1, 2, 3}
	c := []byte{1, 2, 4}
	d := []byte{1, 2}

	if !BytesEqual(a, b) {
		t.Error("Equal byte slices should be equal")
	}
	if BytesEqual(a, c) {
		t.Error("Different byte slices should not be equal")
	}
	if BytesEqual(a, d) {
		t.Error("Different length byte slices should not be equal")
	}
}

// TestCSVRecordsEqual 验证 CSVRecordsEqual 辅助函数
func TestCSVRecordsEqual(t *testing.T) {
	r1 := &CSVRecord{ID: 1, Name: "test", Value: 1.5}
	r2 := &CSVRecord{ID: 1, Name: "test", Value: 1.5}
	r3 := &CSVRecord{ID: 2, Name: "test", Value: 1.5}

	a := []*CSVRecord{r1}
	b := []*CSVRecord{r2}
	c := []*CSVRecord{r3}

	if !CSVRecordsEqual(a, b) {
		t.Error("Equal CSVRecord slices should be equal")
	}
	if CSVRecordsEqual(a, c) {
		t.Error("Different CSVRecord slices should not be equal")
	}
}

// TestHashPairsEqual 验证 HashPairsEqual 辅助函数
func TestHashPairsEqual(t *testing.T) {
	a := [][]interface{}{{"key1", "value1"}, {"key2", "value2"}}
	b := [][]interface{}{{"key1", "value1"}, {"key2", "value2"}}
	c := [][]interface{}{{"key1", "value1"}, {"key2", "different"}}

	if !HashPairsEqual(a, b) {
		t.Error("Equal hash pairs should be equal")
	}
	if HashPairsEqual(a, c) {
		t.Error("Different hash pairs should not be equal")
	}
}
