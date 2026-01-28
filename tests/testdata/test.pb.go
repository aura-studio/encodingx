// Code generated manually for testing purposes.
// This file provides proto.Message implementations for testing the Protobuf encoder.

package testdata

import (
"sync"

"google.golang.org/protobuf/proto"
"google.golang.org/protobuf/reflect/protodesc"
"google.golang.org/protobuf/reflect/protoreflect"
"google.golang.org/protobuf/reflect/protoregistry"
"google.golang.org/protobuf/types/descriptorpb"
"google.golang.org/protobuf/types/dynamicpb"
)

var (
initOnce sync.Once
fileDesc protoreflect.FileDescriptor
testMessageDesc   protoreflect.MessageDescriptor
nestedMessageDesc protoreflect.MessageDescriptor
emptyMessageDesc  protoreflect.MessageDescriptor
)

func initDescriptors() {
initOnce.Do(func() {
fdp := &descriptorpb.FileDescriptorProto{
Name:    proto.String("test.proto"),
Package: proto.String("testdata"),
Syntax:  proto.String("proto3"),
MessageType: []*descriptorpb.DescriptorProto{
{
Name: proto.String("TestMessage"),
Field: []*descriptorpb.FieldDescriptorProto{
{Name: proto.String("int_field"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), JsonName: proto.String("intField")},
{Name: proto.String("string_field"), Number: proto.Int32(2), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), JsonName: proto.String("stringField")},
{Name: proto.String("bool_field"), Number: proto.Int32(3), Type: descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), JsonName: proto.String("boolField")},
},
},
{
Name: proto.String("NestedMessage"),
Field: []*descriptorpb.FieldDescriptorProto{
{Name: proto.String("name"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), JsonName: proto.String("name")},
{Name: proto.String("inner"), Number: proto.Int32(2), Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".testdata.TestMessage"), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), JsonName: proto.String("inner")},
},
},
{Name: proto.String("EmptyMessage"), Field: []*descriptorpb.FieldDescriptorProto{}},
},
}
var err error
fileDesc, err = protodesc.NewFile(fdp, protoregistry.GlobalFiles)
if err != nil { panic(err) }
testMessageDesc = fileDesc.Messages().ByName("TestMessage")
nestedMessageDesc = fileDesc.Messages().ByName("NestedMessage")
emptyMessageDesc = fileDesc.Messages().ByName("EmptyMessage")
})
}
type TestMessage struct { *dynamicpb.Message }

func NewTestMessage() *TestMessage {
initDescriptors()
return &TestMessage{dynamicpb.NewMessage(testMessageDesc)}
}

func NewTestMessageWithValues(intField int32, stringField string, boolField bool) *TestMessage {
m := NewTestMessage()
m.SetIntField(intField)
m.SetStringField(stringField)
m.SetBoolField(boolField)
return m
}

func (m *TestMessage) GetIntField() int32 {
if m.Message == nil { return 0 }
return int32(m.Message.Get(testMessageDesc.Fields().ByNumber(1)).Int())
}

func (m *TestMessage) SetIntField(v int32) {
m.Message.Set(testMessageDesc.Fields().ByNumber(1), protoreflect.ValueOfInt32(v))
}

func (m *TestMessage) GetStringField() string {
if m.Message == nil { return "" }
return m.Message.Get(testMessageDesc.Fields().ByNumber(2)).String()
}

func (m *TestMessage) SetStringField(v string) {
m.Message.Set(testMessageDesc.Fields().ByNumber(2), protoreflect.ValueOfString(v))
}

func (m *TestMessage) GetBoolField() bool {
if m.Message == nil { return false }
return m.Message.Get(testMessageDesc.Fields().ByNumber(3)).Bool()
}

func (m *TestMessage) SetBoolField(v bool) {
m.Message.Set(testMessageDesc.Fields().ByNumber(3), protoreflect.ValueOfBool(v))
}
type NestedMessage struct { *dynamicpb.Message }

func NewNestedMessage() *NestedMessage {
initDescriptors()
return &NestedMessage{dynamicpb.NewMessage(nestedMessageDesc)}
}

func NewNestedMessageWithValues(name string, inner *TestMessage) *NestedMessage {
m := NewNestedMessage()
m.SetName(name)
if inner != nil { m.SetInner(inner) }
return m
}

func (m *NestedMessage) GetName() string {
if m.Message == nil { return "" }
return m.Message.Get(nestedMessageDesc.Fields().ByNumber(1)).String()
}

func (m *NestedMessage) SetName(v string) {
m.Message.Set(nestedMessageDesc.Fields().ByNumber(1), protoreflect.ValueOfString(v))
}

func (m *NestedMessage) GetInner() *TestMessage {
if m.Message == nil { return nil }
innerMsg := m.Message.Get(nestedMessageDesc.Fields().ByNumber(2)).Message()
if innerMsg == nil || !innerMsg.IsValid() { return nil }
return &TestMessage{innerMsg.Interface().(*dynamicpb.Message)}
}

func (m *NestedMessage) SetInner(v *TestMessage) {
m.Message.Set(nestedMessageDesc.Fields().ByNumber(2), protoreflect.ValueOfMessage(v.Message))
}

type EmptyMessage struct { *dynamicpb.Message }

func NewEmptyMessage() *EmptyMessage {
initDescriptors()
return &EmptyMessage{dynamicpb.NewMessage(emptyMessageDesc)}
}