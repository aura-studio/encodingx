package encodingx

import (
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/aura-studio/reflectx"
)

var (
	ErrHexWrongValueType = errors.New("encoding Hex converts on wrong type value")
	ErrHexInvalidData    = errors.New("encoding Hex invalid data")
)

// Hex is a gradient hex encoding without compression.
// Output lengths: 16 -> 32 -> 64 -> 128 -> 256 -> ... characters (doubles each tier)
// No upper limit - automatically expands to fit data.
// Format: [4 bytes length prefix (big-endian)] + [raw data] + [padding]
type Hex struct{}

func init() {
	register(NewHex())
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

// findHexTierSize finds the smallest tier that can hold the payload
// Tiers: 8, 16, 32, 64, 128, 256, ... bytes (no upper limit)
func findHexTierSize(payloadLen int) int {
	tierSize := 8
	for tierSize < payloadLen {
		tierSize *= 2
	}
	return tierSize
}

// isHexTierSize checks if size is a valid tier (power of 2, >= 8)
func isHexTierSize(size int) bool {
	if size < 8 {
		return false
	}
	return size&(size-1) == 0
}

func (Hex) Marshal(v interface{}) ([]byte, error) {
	var data []byte
	switch v := v.(type) {
	case []byte:
		data = v
	case Bytes:
		data = v.Data
	case *Bytes:
		data = v.Data
	default:
		return nil, ErrHexWrongValueType
	}

	// Find suitable tier (need 4 bytes for length prefix)
	payloadLen := len(data) + 4
	tierSize := findHexTierSize(payloadLen)

	// Build output: [4 bytes length (big-endian)] + [data] + [zero padding]
	output := make([]byte, tierSize)
	binary.BigEndian.PutUint32(output[:4], uint32(len(data)))
	copy(output[4:], data)

	return []byte(hex.EncodeToString(output)), nil
}

func (Hex) Unmarshal(data []byte, v interface{}) error {
	// Decode hex
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return err
	}

	// Validate length is a valid tier
	if !isHexTierSize(len(decoded)) {
		return ErrHexInvalidData
	}

	if len(decoded) < 4 {
		return ErrHexInvalidData
	}

	// Extract data length
	dataLen := int(binary.BigEndian.Uint32(decoded[:4]))
	if dataLen > len(decoded)-4 {
		return ErrHexInvalidData
	}
	rawData := decoded[4 : 4+dataLen]

	switch v := v.(type) {
	case *Bytes:
		v.Data = rawData
		return nil
	default:
		return ErrHexWrongValueType
	}
}

func (h Hex) Reverse() Encoding {
	return h
}
