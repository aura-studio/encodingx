package encodingx

import (
	"errors"

	"github.com/aura-studio/reflectx"
	flatbuffers "github.com/google/flatbuffers/go"
)

var (
	// ErrFlatBuffersWrongValueType is returned when Marshal or Unmarshal is called
	// with a type that doesn't implement the required FlatBuffers interfaces.
	ErrFlatBuffersWrongValueType = errors.New("encoding flatbuffers converts on wrong type value")
)

// FlatBufferMarshaler is the interface that types must implement to be marshaled
// using the FlatBuffers encoder. Types implementing this interface should use
// a FlatBuffers Builder to serialize their data.
type FlatBufferMarshaler interface {
	MarshalFlatBuffer(builder *flatbuffers.Builder) error
}

// FlatBufferUnmarshaler is the interface that types must implement to be unmarshaled
// using the FlatBuffers encoder. Types implementing this interface should initialize
// themselves from the provided byte slice.
type FlatBufferUnmarshaler interface {
	UnmarshalFlatBuffer(data []byte) error
}

// FlatBuffers implements the Encoding interface for FlatBuffers serialization.
// FlatBuffers is a high-performance serialization library that provides
// zero-copy access to serialized data.
//
// Note: FlatBuffers requires types to implement FlatBufferMarshaler and
// FlatBufferUnmarshaler interfaces for serialization/deserialization.
// Unlike JSON or MsgPack, FlatBuffers doesn't support arbitrary struct
// serialization out of the box.
type FlatBuffers struct{}

func init() {
	register(NewFlatBuffers())
}

// NewFlatBuffers creates a new FlatBuffers encoder instance.
func NewFlatBuffers() *FlatBuffers {
	return new(FlatBuffers)
}

// String returns the type name of the encoder.
func (f FlatBuffers) String() string {
	return reflectx.TypeName(f)
}

// Style returns the encoding style type.
// FlatBuffers uses EncodingStyleStruct as it serializes structured data.
func (f FlatBuffers) Style() EncodingStyleType {
	return EncodingStyleStruct
}

// Marshal serializes the given value to FlatBuffers format.
// The value must implement the FlatBufferMarshaler interface.
// Returns ErrFlatBuffersWrongValueType if the value doesn't implement the interface.
func (FlatBuffers) Marshal(v interface{}) ([]byte, error) {
	marshaler, ok := v.(FlatBufferMarshaler)
	if !ok {
		return nil, ErrFlatBuffersWrongValueType
	}

	builder := flatbuffers.NewBuilder(0)
	if err := marshaler.MarshalFlatBuffer(builder); err != nil {
		return nil, err
	}

	return builder.FinishedBytes(), nil
}

// Unmarshal deserializes FlatBuffers data into the given value.
// The value must implement the FlatBufferUnmarshaler interface.
// Returns ErrFlatBuffersWrongValueType if the value doesn't implement the interface.
func (FlatBuffers) Unmarshal(data []byte, v interface{}) error {
	unmarshaler, ok := v.(FlatBufferUnmarshaler)
	if !ok {
		return ErrFlatBuffersWrongValueType
	}

	return unmarshaler.UnmarshalFlatBuffer(data)
}

// Reverse returns the encoder itself since FlatBuffers is symmetric
// (the same encoder is used for both serialization and deserialization).
func (f FlatBuffers) Reverse() Encoding {
	return f
}
