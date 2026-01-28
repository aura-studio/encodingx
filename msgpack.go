package encodingx

import (
	"github.com/aura-studio/reflectx"
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPack implements the Encoding interface for MessagePack serialization.
// MessagePack is a binary serialization format that is more compact than JSON.
type MsgPack struct{}

func init() {
	register(NewMsgPack())
}

// NewMsgPack creates a new MsgPack encoder instance.
func NewMsgPack() *MsgPack {
	return new(MsgPack)
}

// String returns the type name of the encoder.
func (m MsgPack) String() string {
	return reflectx.TypeName(m)
}

// Style returns the encoding style type.
// MsgPack uses EncodingStyleStruct as it serializes structured data.
func (m MsgPack) Style() EncodingStyleType {
	return EncodingStyleStruct
}

// Marshal serializes the given value to MsgPack format.
func (MsgPack) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Unmarshal deserializes MsgPack data into the given value.
func (MsgPack) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

// Reverse returns the encoder itself since MsgPack is symmetric
// (the same encoder is used for both serialization and deserialization).
func (m MsgPack) Reverse() Encoding {
	return m
}
