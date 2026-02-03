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
	ErrHexTierWrongValueType = errors.New("encoding HexTier converts on wrong type value")
	ErrHexTierInvalidData    = errors.New("encoding HexTier invalid data")
)

// HexTier is a gradient hex encoding with zlib compression.
// Output lengths: 16 -> 32 -> 64 -> 128 -> 256 -> ... characters (doubles each tier)
// No upper limit - automatically expands to fit data.
// Format: [4 bytes length prefix (big-endian)] + [compressed data] + [padding]
type HexTier struct{}

func init() {
	register(NewHexTier())
}

func NewHexTier() *HexTier {
	return new(HexTier)
}

func (h HexTier) String() string {
	return reflectx.TypeName(h)
}

func (HexTier) Style() EncodingStyleType {
	return EncodingStyleBytes
}

// findTierSize finds the smallest tier that can hold the payload
// Tiers: 8, 16, 32, 64, 128, 256, ... bytes (no upper limit)
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

func (HexTier) Marshal(v interface{}) ([]byte, error) {
	var data []byte
	switch v := v.(type) {
	case []byte:
		data = v
	case Bytes:
		data = v.Data
	case *Bytes:
		data = v.Data
	default:
		return nil, ErrHexTierWrongValueType
	}

	// Compress data
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data)
	w.Close()
	compressed := buf.Bytes()

	// Find suitable tier (need 4 bytes for length prefix)
	payloadLen := len(compressed) + 4
	tierSize := findTierSize(payloadLen)

	// Build output: [4 bytes length (big-endian)] + [compressed] + [zero padding]
	output := make([]byte, tierSize)
	binary.BigEndian.PutUint32(output[:4], uint32(len(compressed)))
	copy(output[4:], compressed)

	return []byte(hex.EncodeToString(output)), nil
}

func (HexTier) Unmarshal(data []byte, v interface{}) error {
	// Decode hex
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	// Validate length is a valid tier
	if !isTierSize(len(decoded)) {
		return ErrHexTierInvalidData
	}

	if len(decoded) < 4 {
		return ErrHexTierInvalidData
	}

	// Extract compressed data length
	compressedLen := int(binary.BigEndian.Uint32(decoded[:4]))
	if compressedLen > len(decoded)-4 {
		return ErrHexTierInvalidData
	}
	compressed := decoded[4 : 4+compressedLen]

	// Decompress
	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return err
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	switch v := v.(type) {
	case *Bytes:
		v.Data = decompressed
		return nil
	default:
		return ErrHexTierWrongValueType
	}
}

func (h HexTier) Reverse() Encoding {
	return h
}
