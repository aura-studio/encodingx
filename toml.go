package encodingx

import (
	"github.com/aura-studio/reflectx"
	"github.com/pelletier/go-toml/v2"
)

// TOML implements the Encoding interface for TOML serialization.
// TOML is a configuration file format that is easy to read and write.
type TOML struct{}

func init() {
	register(NewTOML())
}

// NewTOML creates a new TOML encoder instance.
func NewTOML() *TOML {
	return new(TOML)
}

// String returns the type name of the encoder.
func (t TOML) String() string {
	return reflectx.TypeName(t)
}

// Style returns the encoding style type.
// TOML uses EncodingStyleStruct as it serializes structured data.
func (t TOML) Style() EncodingStyleType {
	return EncodingStyleStruct
}

// Marshal serializes the given value to TOML format.
func (TOML) Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}

// Unmarshal deserializes TOML data into the given value.
func (TOML) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

// Reverse returns the encoder itself since TOML is symmetric
// (the same encoder is used for both serialization and deserialization).
func (t TOML) Reverse() Encoding {
	return t
}
