// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mafia.proto

package mafia_proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ActionType int32

const (
	ActionType_Init      ActionType = 0
	ActionType_Vote      ActionType = 1
	ActionType_DoNothing ActionType = 2
)

var ActionType_name = map[int32]string{
	0: "Init",
	1: "Vote",
	2: "DoNothing",
}

var ActionType_value = map[string]int32{
	"Init":      0,
	"Vote":      1,
	"DoNothing": 2,
}

func (x ActionType) String() string {
	return proto.EnumName(ActionType_name, int32(x))
}

func (ActionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{0}
}

type EventType int32

const (
	EventType_VoteRequest   EventType = 0
	EventType_SystemMessage EventType = 1
	EventType_GameEnd       EventType = 2
)

var EventType_name = map[int32]string{
	0: "VoteRequest",
	1: "SystemMessage",
	2: "GameEnd",
}

var EventType_value = map[string]int32{
	"VoteRequest":   0,
	"SystemMessage": 1,
	"GameEnd":       2,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}

func (EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{1}
}

type Action struct {
	// Types that are valid to be assigned to Data:
	//	*Action_Init_
	//	*Action_Vote_
	//	*Action_DoNothing_
	Data                 isAction_Data `protobuf_oneof:"data"`
	Type                 ActionType    `protobuf:"varint,4,opt,name=type,proto3,enum=mafia.ActionType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()    {}
func (*Action) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{0}
}

func (m *Action) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Action.Unmarshal(m, b)
}
func (m *Action) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Action.Marshal(b, m, deterministic)
}
func (m *Action) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Action.Merge(m, src)
}
func (m *Action) XXX_Size() int {
	return xxx_messageInfo_Action.Size(m)
}
func (m *Action) XXX_DiscardUnknown() {
	xxx_messageInfo_Action.DiscardUnknown(m)
}

var xxx_messageInfo_Action proto.InternalMessageInfo

type isAction_Data interface {
	isAction_Data()
}

type Action_Init_ struct {
	Init *Action_Init `protobuf:"bytes,1,opt,name=init,proto3,oneof"`
}

type Action_Vote_ struct {
	Vote *Action_Vote `protobuf:"bytes,2,opt,name=vote,proto3,oneof"`
}

type Action_DoNothing_ struct {
	DoNothing *Action_DoNothing `protobuf:"bytes,3,opt,name=do_nothing,json=doNothing,proto3,oneof"`
}

func (*Action_Init_) isAction_Data() {}

func (*Action_Vote_) isAction_Data() {}

func (*Action_DoNothing_) isAction_Data() {}

func (m *Action) GetData() isAction_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Action) GetInit() *Action_Init {
	if x, ok := m.GetData().(*Action_Init_); ok {
		return x.Init
	}
	return nil
}

func (m *Action) GetVote() *Action_Vote {
	if x, ok := m.GetData().(*Action_Vote_); ok {
		return x.Vote
	}
	return nil
}

func (m *Action) GetDoNothing() *Action_DoNothing {
	if x, ok := m.GetData().(*Action_DoNothing_); ok {
		return x.DoNothing
	}
	return nil
}

func (m *Action) GetType() ActionType {
	if m != nil {
		return m.Type
	}
	return ActionType_Init
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Action) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Action_Init_)(nil),
		(*Action_Vote_)(nil),
		(*Action_DoNothing_)(nil),
	}
}

type Action_Init struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Action_Init) Reset()         { *m = Action_Init{} }
func (m *Action_Init) String() string { return proto.CompactTextString(m) }
func (*Action_Init) ProtoMessage()    {}
func (*Action_Init) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{0, 0}
}

func (m *Action_Init) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Action_Init.Unmarshal(m, b)
}
func (m *Action_Init) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Action_Init.Marshal(b, m, deterministic)
}
func (m *Action_Init) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Action_Init.Merge(m, src)
}
func (m *Action_Init) XXX_Size() int {
	return xxx_messageInfo_Action_Init.Size(m)
}
func (m *Action_Init) XXX_DiscardUnknown() {
	xxx_messageInfo_Action_Init.DiscardUnknown(m)
}

var xxx_messageInfo_Action_Init proto.InternalMessageInfo

func (m *Action_Init) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Action_Vote struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Action_Vote) Reset()         { *m = Action_Vote{} }
func (m *Action_Vote) String() string { return proto.CompactTextString(m) }
func (*Action_Vote) ProtoMessage()    {}
func (*Action_Vote) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{0, 1}
}

func (m *Action_Vote) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Action_Vote.Unmarshal(m, b)
}
func (m *Action_Vote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Action_Vote.Marshal(b, m, deterministic)
}
func (m *Action_Vote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Action_Vote.Merge(m, src)
}
func (m *Action_Vote) XXX_Size() int {
	return xxx_messageInfo_Action_Vote.Size(m)
}
func (m *Action_Vote) XXX_DiscardUnknown() {
	xxx_messageInfo_Action_Vote.DiscardUnknown(m)
}

var xxx_messageInfo_Action_Vote proto.InternalMessageInfo

func (m *Action_Vote) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Action_DoNothing struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Action_DoNothing) Reset()         { *m = Action_DoNothing{} }
func (m *Action_DoNothing) String() string { return proto.CompactTextString(m) }
func (*Action_DoNothing) ProtoMessage()    {}
func (*Action_DoNothing) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{0, 2}
}

func (m *Action_DoNothing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Action_DoNothing.Unmarshal(m, b)
}
func (m *Action_DoNothing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Action_DoNothing.Marshal(b, m, deterministic)
}
func (m *Action_DoNothing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Action_DoNothing.Merge(m, src)
}
func (m *Action_DoNothing) XXX_Size() int {
	return xxx_messageInfo_Action_DoNothing.Size(m)
}
func (m *Action_DoNothing) XXX_DiscardUnknown() {
	xxx_messageInfo_Action_DoNothing.DiscardUnknown(m)
}

var xxx_messageInfo_Action_DoNothing proto.InternalMessageInfo

type Event struct {
	// Types that are valid to be assigned to Data:
	//	*Event_VoteRequest_
	//	*Event_Message
	//	*Event_GameEnd_
	Data                 isEvent_Data `protobuf_oneof:"data"`
	Type                 EventType    `protobuf:"varint,4,opt,name=type,proto3,enum=mafia.EventType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{1}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

type isEvent_Data interface {
	isEvent_Data()
}

type Event_VoteRequest_ struct {
	VoteRequest *Event_VoteRequest `protobuf:"bytes,1,opt,name=vote_request,json=voteRequest,proto3,oneof"`
}

type Event_Message struct {
	Message *Event_SystemMessage `protobuf:"bytes,2,opt,name=message,proto3,oneof"`
}

type Event_GameEnd_ struct {
	GameEnd *Event_GameEnd `protobuf:"bytes,3,opt,name=game_end,json=gameEnd,proto3,oneof"`
}

func (*Event_VoteRequest_) isEvent_Data() {}

func (*Event_Message) isEvent_Data() {}

func (*Event_GameEnd_) isEvent_Data() {}

func (m *Event) GetData() isEvent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event) GetVoteRequest() *Event_VoteRequest {
	if x, ok := m.GetData().(*Event_VoteRequest_); ok {
		return x.VoteRequest
	}
	return nil
}

func (m *Event) GetMessage() *Event_SystemMessage {
	if x, ok := m.GetData().(*Event_Message); ok {
		return x.Message
	}
	return nil
}

func (m *Event) GetGameEnd() *Event_GameEnd {
	if x, ok := m.GetData().(*Event_GameEnd_); ok {
		return x.GameEnd
	}
	return nil
}

func (m *Event) GetType() EventType {
	if m != nil {
		return m.Type
	}
	return EventType_VoteRequest
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Event) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Event_VoteRequest_)(nil),
		(*Event_Message)(nil),
		(*Event_GameEnd_)(nil),
	}
}

type Event_VoteRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_VoteRequest) Reset()         { *m = Event_VoteRequest{} }
func (m *Event_VoteRequest) String() string { return proto.CompactTextString(m) }
func (*Event_VoteRequest) ProtoMessage()    {}
func (*Event_VoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{1, 0}
}

func (m *Event_VoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event_VoteRequest.Unmarshal(m, b)
}
func (m *Event_VoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event_VoteRequest.Marshal(b, m, deterministic)
}
func (m *Event_VoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_VoteRequest.Merge(m, src)
}
func (m *Event_VoteRequest) XXX_Size() int {
	return xxx_messageInfo_Event_VoteRequest.Size(m)
}
func (m *Event_VoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_VoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Event_VoteRequest proto.InternalMessageInfo

type Event_SystemMessage struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_SystemMessage) Reset()         { *m = Event_SystemMessage{} }
func (m *Event_SystemMessage) String() string { return proto.CompactTextString(m) }
func (*Event_SystemMessage) ProtoMessage()    {}
func (*Event_SystemMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{1, 1}
}

func (m *Event_SystemMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event_SystemMessage.Unmarshal(m, b)
}
func (m *Event_SystemMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event_SystemMessage.Marshal(b, m, deterministic)
}
func (m *Event_SystemMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_SystemMessage.Merge(m, src)
}
func (m *Event_SystemMessage) XXX_Size() int {
	return xxx_messageInfo_Event_SystemMessage.Size(m)
}
func (m *Event_SystemMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_SystemMessage.DiscardUnknown(m)
}

var xxx_messageInfo_Event_SystemMessage proto.InternalMessageInfo

func (m *Event_SystemMessage) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type Event_GameEnd struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Dead                 bool     `protobuf:"varint,2,opt,name=dead,proto3" json:"dead,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_GameEnd) Reset()         { *m = Event_GameEnd{} }
func (m *Event_GameEnd) String() string { return proto.CompactTextString(m) }
func (*Event_GameEnd) ProtoMessage()    {}
func (*Event_GameEnd) Descriptor() ([]byte, []int) {
	return fileDescriptor_7077fdb094f3bc8b, []int{1, 2}
}

func (m *Event_GameEnd) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event_GameEnd.Unmarshal(m, b)
}
func (m *Event_GameEnd) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event_GameEnd.Marshal(b, m, deterministic)
}
func (m *Event_GameEnd) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_GameEnd.Merge(m, src)
}
func (m *Event_GameEnd) XXX_Size() int {
	return xxx_messageInfo_Event_GameEnd.Size(m)
}
func (m *Event_GameEnd) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_GameEnd.DiscardUnknown(m)
}

var xxx_messageInfo_Event_GameEnd proto.InternalMessageInfo

func (m *Event_GameEnd) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Event_GameEnd) GetDead() bool {
	if m != nil {
		return m.Dead
	}
	return false
}

func init() {
	proto.RegisterEnum("mafia.ActionType", ActionType_name, ActionType_value)
	proto.RegisterEnum("mafia.EventType", EventType_name, EventType_value)
	proto.RegisterType((*Action)(nil), "mafia.Action")
	proto.RegisterType((*Action_Init)(nil), "mafia.Action.Init")
	proto.RegisterType((*Action_Vote)(nil), "mafia.Action.Vote")
	proto.RegisterType((*Action_DoNothing)(nil), "mafia.Action.DoNothing")
	proto.RegisterType((*Event)(nil), "mafia.Event")
	proto.RegisterType((*Event_VoteRequest)(nil), "mafia.Event.VoteRequest")
	proto.RegisterType((*Event_SystemMessage)(nil), "mafia.Event.SystemMessage")
	proto.RegisterType((*Event_GameEnd)(nil), "mafia.Event.GameEnd")
}

func init() { proto.RegisterFile("mafia.proto", fileDescriptor_7077fdb094f3bc8b) }

var fileDescriptor_7077fdb094f3bc8b = []byte{
	// 435 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0xc1, 0x6f, 0xd3, 0x30,
	0x14, 0xc6, 0x93, 0x90, 0x75, 0xed, 0x4b, 0x0b, 0xd9, 0x13, 0x12, 0x21, 0x5c, 0xa6, 0x02, 0x52,
	0x35, 0x89, 0x8e, 0x15, 0x09, 0x21, 0xc1, 0x0e, 0x4c, 0x0c, 0xca, 0x61, 0x1c, 0x3c, 0xc4, 0x81,
	0x4b, 0x64, 0xe6, 0x47, 0xf0, 0x21, 0x76, 0x69, 0x4c, 0x44, 0xff, 0x05, 0xfe, 0x30, 0xfe, 0x2e,
	0x64, 0xa7, 0xe9, 0x62, 0xd4, 0x53, 0x3f, 0xfb, 0xfb, 0x3d, 0xf9, 0xf3, 0xe7, 0x06, 0x92, 0x8a,
	0x7f, 0x97, 0x7c, 0xbe, 0x5a, 0x6b, 0xa3, 0xf1, 0xc0, 0x2d, 0xa6, 0x7f, 0x22, 0x18, 0xbc, 0xbd,
	0x31, 0x52, 0x2b, 0x9c, 0x41, 0x2c, 0x95, 0x34, 0x59, 0x78, 0x1c, 0xce, 0x92, 0x05, 0xce, 0x5b,
	0xba, 0x35, 0xe7, 0x1f, 0x95, 0x34, 0xcb, 0x80, 0x39, 0xc2, 0x92, 0x8d, 0x36, 0x94, 0x45, 0xfb,
	0xc8, 0x2f, 0xda, 0x90, 0x25, 0x2d, 0x81, 0xaf, 0x00, 0x84, 0x2e, 0x94, 0x36, 0x3f, 0xa4, 0x2a,
	0xb3, 0x3b, 0x8e, 0x7f, 0xe0, 0xf3, 0xef, 0xf4, 0xa7, 0xd6, 0x5e, 0x06, 0x6c, 0x24, 0xba, 0x05,
	0x3e, 0x85, 0xd8, 0x6c, 0x56, 0x94, 0xc5, 0xc7, 0xe1, 0xec, 0xee, 0xe2, 0xc8, 0x9b, 0xf9, 0xbc,
	0x59, 0x11, 0x73, 0x76, 0x9e, 0x43, 0x6c, 0xa3, 0x21, 0x42, 0xac, 0x78, 0x45, 0x2e, 0xfc, 0x88,
	0x39, 0x6d, 0x3d, 0x1b, 0x66, 0xaf, 0x97, 0xc0, 0x68, 0x77, 0xf0, 0xc5, 0x00, 0x62, 0xc1, 0x0d,
	0x9f, 0xfe, 0x8d, 0xe0, 0xe0, 0xb2, 0x21, 0x65, 0xf0, 0x1c, 0xc6, 0x36, 0x7f, 0xb1, 0xa6, 0x9f,
	0xbf, 0xa8, 0xee, 0x3a, 0xc9, 0xb6, 0x29, 0x1c, 0xe3, 0x2e, 0xca, 0x5a, 0x7f, 0x19, 0xb0, 0xa4,
	0xb9, 0x5d, 0xe2, 0x4b, 0x38, 0xac, 0xa8, 0xae, 0x79, 0xd9, 0x75, 0x94, 0x7b, 0x93, 0xd7, 0x9b,
	0xda, 0x50, 0x75, 0xd5, 0x12, 0xcb, 0x80, 0x75, 0x30, 0x9e, 0xc1, 0xb0, 0xe4, 0x15, 0x15, 0xa4,
	0xc4, 0xb6, 0xac, 0xfb, 0xde, 0xe0, 0x07, 0x5e, 0xd1, 0xa5, 0x12, 0x76, 0xa4, 0x6c, 0x25, 0x3e,
	0xf1, 0x7a, 0x4a, 0xfb, 0x78, 0xaf, 0xa6, 0x09, 0x24, 0xbd, 0xb8, 0xf9, 0x63, 0x98, 0x78, 0x19,
	0x6c, 0x45, 0x86, 0x7e, 0x9b, 0xae, 0x22, 0xab, 0xf3, 0x33, 0x38, 0xdc, 0x9e, 0xb7, 0xcf, 0xb6,
	0x7b, 0x82, 0xb8, 0x70, 0x17, 0x1c, 0x32, 0xa7, 0xbb, 0x22, 0x4f, 0x4e, 0x01, 0x6e, 0x5f, 0x0a,
	0x87, 0xed, 0x1b, 0xa5, 0x81, 0x55, 0x36, 0x46, 0x1a, 0xe2, 0xa4, 0xd7, 0x7f, 0x1a, 0x9d, 0xbc,
	0x81, 0xd1, 0x2e, 0x32, 0xde, 0xf3, 0xc2, 0xa6, 0x01, 0x1e, 0xfd, 0x17, 0x37, 0x0d, 0x31, 0xd9,
	0x85, 0x4b, 0xa3, 0xc5, 0x39, 0x8c, 0xaf, 0xec, 0xb5, 0xaf, 0x69, 0xdd, 0xc8, 0x1b, 0xc2, 0x67,
	0x30, 0x7c, 0x2f, 0x95, 0xb0, 0x00, 0x4e, 0xbc, 0x7f, 0x4e, 0x3e, 0xee, 0x17, 0x34, 0x0d, 0x66,
	0xe1, 0xf3, 0xf0, 0xe2, 0xd1, 0xd7, 0x87, 0xf3, 0x53, 0xb7, 0x5d, 0xb8, 0x6f, 0xe3, 0x75, 0x4f,
	0x7f, 0x1b, 0xb8, 0x9f, 0x17, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x3f, 0x69, 0x79, 0x3d,
	0x03, 0x00, 0x00,
}
