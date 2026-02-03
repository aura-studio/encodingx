package encodingx

import (
	"bytes"
	"compress/zlib"
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
