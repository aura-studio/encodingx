package encodingx

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"

	"github.com/aura-studio/reflectx"
)

var (
	ErrHexWrongValueType = errors.New("encoding Hex converts on wrong type value")
	ErrHexInvalidData    = errors.New("encoding Hex invalid data")
)

// ============================================================================
// Hex - Pure hex encoding (no length prefix)
// Format: [raw data] -> hex string
// Output length = data_length * 2 hex characters
// ============================================================================

type Hex struct{}

func init() {
	register(NewHex())
	register(NewHexTier())
	register(NewHexTierRand())
	register(NewHexZlib())
}

func NewHex() *Hex {
	return new(Hex)
}

func (h Hex) String() string {
	return reflectx.TypeName(h)
}

func (Hex) Style() EncodingStyleType {
	return EncodingStyleBytes
}

func (Hex) Marshal(v any) ([]byte, error) {
	data, err := toBytes(v)
	if err != nil {
		return nil, ErrHexWrongValueType
	}

	return []byte(hex.EncodeToString(data)), nil
}

func (Hex) Unmarshal(data []byte, v any) error {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	return fromBytes(decoded, v)
}

func (h Hex) Reverse() Encoding {
	return h
}

// ============================================================================
// HexTier - Hex encoding with tier padding (no compression)
// Format: [4 bytes length (big-endian)] + [raw data] + [zero padding]
// Tiers: 8, 16, 32, 64, 128, 256, ... bytes (power of 2)
// ============================================================================

type HexTier struct{}

func NewHexTier() *HexTier {
	return new(HexTier)
}

func (h HexTier) String() string {
	return reflectx.TypeName(h)
}

func (HexTier) Style() EncodingStyleType {
	return EncodingStyleBytes
}

func (HexTier) Marshal(v any) ([]byte, error) {
	data, err := toBytes(v)
	if err != nil {
		return nil, ErrHexWrongValueType
	}

	tierSize := findTierSize(len(data) + 4)
	output := make([]byte, tierSize)
	binary.BigEndian.PutUint32(output[:4], uint32(len(data)))
	copy(output[4:], data)

	return []byte(hex.EncodeToString(output)), nil
}

func (HexTier) Unmarshal(data []byte, v any) error {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	if !isTierSize(len(decoded)) {
		return ErrHexInvalidData
	}

	if len(decoded) < 4 {
		return ErrHexInvalidData
	}

	dataLen := int(binary.BigEndian.Uint32(decoded[:4]))
	if dataLen > len(decoded)-4 {
		return ErrHexInvalidData
	}

	return fromBytes(decoded[4:4+dataLen], v)
}

func (h HexTier) Reverse() Encoding {
	return h
}

// ============================================================================
// HexTierRand - Hex encoding with random padding and rolling XOR obfuscation
// Format: [4 bytes random key] + [4 bytes XORed length] + [XORed data] + [random padding]
// Tiers: 16, 32, 64, 128, 256, ... bytes (power of 2, min 16 for key+length)
// Each byte uses different XOR key (rolling), making output appear fully random
// ============================================================================

type HexTierRand struct{}

func NewHexTierRand() *HexTierRand {
	return new(HexTierRand)
}

func (h HexTierRand) String() string {
	return reflectx.TypeName(h)
}

func (HexTierRand) Style() EncodingStyleType {
	return EncodingStyleBytes
}

func (HexTierRand) Marshal(v any) ([]byte, error) {
	data, err := toBytes(v)
	if err != nil {
		return nil, ErrHexWrongValueType
	}

	// Generate 4-byte random key
	key := make([]byte, 4)
	rand.Read(key)

	tierSize := findTierSizeMin(len(data)+8, 16) // min 16 for key(4)+len(4)+data
	output := make([]byte, tierSize)

	// Fill entire buffer with random bytes first (random padding)
	rand.Read(output)

	// First 4 bytes: random key (unobfuscated)
	copy(output[:4], key)

	// Next 4 bytes: length XORed with rolling key
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(data)))
	for i := 0; i < 4; i++ {
		output[4+i] = lenBytes[i] ^ key[i]
	}

	// XOR data with rolling key
	for i, b := range data {
		output[8+i] = b ^ key[i%4]
	}

	return []byte(hex.EncodeToString(output)), nil
}

func (HexTierRand) Unmarshal(data []byte, v any) error {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	if !isTierSize(len(decoded)) {
		return ErrHexInvalidData
	}

	if len(decoded) < 8 {
		return ErrHexInvalidData
	}

	// Extract 4-byte key
	key := decoded[:4]

	// Extract length (XOR decode with rolling key)
	lenBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		lenBytes[i] = decoded[4+i] ^ key[i]
	}
	dataLen := int(binary.BigEndian.Uint32(lenBytes))

	if dataLen > len(decoded)-8 {
		return ErrHexInvalidData
	}

	// XOR decode the data with rolling key
	result := make([]byte, dataLen)
	for i := 0; i < dataLen; i++ {
		result[i] = decoded[8+i] ^ key[i%4]
	}

	return fromBytes(result, v)
}

func (h HexTierRand) Reverse() Encoding {
	return h
}

// findTierSizeMin finds the smallest tier (power of 2) >= minTier that can hold the payload
func findTierSizeMin(payloadLen, minTier int) int {
	tierSize := minTier
	for tierSize < payloadLen {
		tierSize *= 2
	}
	return tierSize
}

// ============================================================================
// HexZlib - Hex encoding with zlib compression and tier padding
// Format: [4 bytes length (big-endian)] + [zlib compressed data] + [zero padding]
// Tiers: 8, 16, 32, 64, 128, 256, ... bytes (power of 2)
// ============================================================================

type HexZlib struct{}

func NewHexZlib() *HexZlib {
	return new(HexZlib)
}

func (h HexZlib) String() string {
	return reflectx.TypeName(h)
}

func (HexZlib) Style() EncodingStyleType {
	return EncodingStyleBytes
}

func (HexZlib) Marshal(v any) ([]byte, error) {
	data, err := toBytes(v)
	if err != nil {
		return nil, ErrHexWrongValueType
	}

	// Compress
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data)
	w.Close()
	compressed := buf.Bytes()

	tierSize := findTierSize(len(compressed) + 4)
	output := make([]byte, tierSize)
	binary.BigEndian.PutUint32(output[:4], uint32(len(compressed)))
	copy(output[4:], compressed)

	return []byte(hex.EncodeToString(output)), nil
}

func (HexZlib) Unmarshal(data []byte, v any) error {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	if !isTierSize(len(decoded)) {
		return ErrHexInvalidData
	}

	if len(decoded) < 4 {
		return ErrHexInvalidData
	}

	compressedLen := int(binary.BigEndian.Uint32(decoded[:4]))
	if compressedLen > len(decoded)-4 {
		return ErrHexInvalidData
	}

	// Decompress
	r, err := zlib.NewReader(bytes.NewReader(decoded[4 : 4+compressedLen]))
	if err != nil {
		return err
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return fromBytes(decompressed, v)
}

func (h HexZlib) Reverse() Encoding {
	return h
}

// ============================================================================
// Helper functions
// ============================================================================

// findTierSize finds the smallest tier (power of 2, >= 8) that can hold the payload
func findTierSize(payloadLen int) int {
	tierSize := 8
	for tierSize < payloadLen {
		tierSize *= 2
	}
	return tierSize
}

// isTierSize checks if size is a valid tier (power of 2, >= 8)
func isTierSize(size int) bool {
	if size < 8 {
		return false
	}
	return size&(size-1) == 0
}

// toBytes converts value to []byte
func toBytes(v any) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil
	case Bytes:
		return v.Data, nil
	case *Bytes:
		return v.Data, nil
	default:
		return nil, ErrHexWrongValueType
	}
}

// fromBytes sets []byte to target value
func fromBytes(data []byte, v any) error {
	switch v := v.(type) {
	case *Bytes:
		v.Data = data
		return nil
	default:
		return ErrHexWrongValueType
	}
}
