package testdata

import (
"testing"

"google.golang.org/protobuf/proto"
)

func TestProtoMessageInterface(t *testing.T) {
msg := NewTestMessageWithValues(42, "hello", true)
var _ proto.Message = msg.Message

data, err := proto.Marshal(msg.Message)
if err != nil {
t.Fatalf("Failed to marshal TestMessage: %v", err)
}

msg2 := NewTestMessage()
err = proto.Unmarshal(data, msg2.Message)
if err != nil {
t.Fatalf("Failed to unmarshal TestMessage: %v", err)
}

if msg2.GetIntField() != msg.GetIntField() {
t.Errorf("IntField mismatch: got %d, want %d", msg2.GetIntField(), msg.GetIntField())
}
if msg2.GetStringField() != msg.GetStringField() {
t.Errorf("StringField mismatch: got %s, want %s", msg2.GetStringField(), msg.GetStringField())
}
if msg2.GetBoolField() != msg.GetBoolField() {
t.Errorf("BoolField mismatch: got %v, want %v", msg2.GetBoolField(), msg.GetBoolField())
}
}

func TestNestedMessage(t *testing.T) {
inner := NewTestMessageWithValues(100, "inner", false)
msg := NewNestedMessageWithValues("nested", inner)

data, err := proto.Marshal(msg.Message)
if err != nil {
t.Fatalf("Failed to marshal NestedMessage: %v", err)
}

msg2 := NewNestedMessage()
err = proto.Unmarshal(data, msg2.Message)
if err != nil {
t.Fatalf("Failed to unmarshal NestedMessage: %v", err)
}

if msg2.GetName() != msg.GetName() {
t.Errorf("Name mismatch: got %s, want %s", msg2.GetName(), msg.GetName())
}

inner2 := msg2.GetInner()
if inner2 == nil {
t.Fatal("Inner is nil")
}
if inner2.GetIntField() != inner.GetIntField() {
t.Errorf("Inner.IntField mismatch: got %d, want %d", inner2.GetIntField(), inner.GetIntField())
}
}

func TestEmptyMessage(t *testing.T) {
msg := NewEmptyMessage()

data, err := proto.Marshal(msg.Message)
if err != nil {
t.Fatalf("Failed to marshal EmptyMessage: %v", err)
}

msg2 := NewEmptyMessage()
err = proto.Unmarshal(data, msg2.Message)
if err != nil {
t.Fatalf("Failed to unmarshal EmptyMessage: %v", err)
}
}