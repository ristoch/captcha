package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ClientEvent_EventType int32

const (
	ClientEvent_FRONTEND_EVENT    ClientEvent_EventType = 0
	ClientEvent_CONNECTION_CLOSED ClientEvent_EventType = 1
	ClientEvent_BALANCER_EVENT    ClientEvent_EventType = 2
)

var (
	ClientEvent_EventType_name = map[int32]string{
		0: "FRONTEND_EVENT",
		1: "CONNECTION_CLOSED",
		2: "BALANCER_EVENT",
	}
	ClientEvent_EventType_value = map[string]int32{
		"FRONTEND_EVENT":    0,
		"CONNECTION_CLOSED": 1,
		"BALANCER_EVENT":    2,
	}
)

func (x ClientEvent_EventType) Enum() *ClientEvent_EventType {
	p := new(ClientEvent_EventType)
	*p = x
	return p
}

func (x ClientEvent_EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientEvent_EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_captcha_captcha_proto_enumTypes[0].Descriptor()
}

func (ClientEvent_EventType) Type() protoreflect.EnumType {
	return &file_proto_captcha_captcha_proto_enumTypes[0]
}

func (x ClientEvent_EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

func (ClientEvent_EventType) EnumDescriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{4, 0}
}

type ChallengeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Complexity    int32                  `protobuf:"varint,1,opt,name=complexity,proto3" json:"complexity,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChallengeRequest) Reset() {
	*x = ChallengeRequest{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChallengeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChallengeRequest) ProtoMessage() {}

func (x *ChallengeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ChallengeRequest) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{0}
}

func (x *ChallengeRequest) GetComplexity() int32 {
	if x != nil {
		return x.Complexity
	}
	return 0
}

func (x *ChallengeRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ChallengeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChallengeId   string                 `protobuf:"bytes,1,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	Html          string                 `protobuf:"bytes,2,opt,name=html,proto3" json:"html,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChallengeResponse) Reset() {
	*x = ChallengeResponse{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChallengeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChallengeResponse) ProtoMessage() {}

func (x *ChallengeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ChallengeResponse) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{1}
}

func (x *ChallengeResponse) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *ChallengeResponse) GetHtml() string {
	if x != nil {
		return x.Html
	}
	return ""
}

type ValidateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChallengeId   string                 `protobuf:"bytes,1,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	Answer        string                 `protobuf:"bytes,2,opt,name=answer,proto3" json:"answer,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidateRequest) Reset() {
	*x = ValidateRequest{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateRequest) ProtoMessage() {}

func (x *ValidateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ValidateRequest) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{2}
}

func (x *ValidateRequest) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *ValidateRequest) GetAnswer() string {
	if x != nil {
		return x.Answer
	}
	return ""
}

type ValidateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Valid         bool                   `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	Confidence    int32                  `protobuf:"varint,2,opt,name=confidence,proto3" json:"confidence,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidateResponse) Reset() {
	*x = ValidateResponse{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateResponse) ProtoMessage() {}

func (x *ValidateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ValidateResponse) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{3}
}

func (x *ValidateResponse) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

func (x *ValidateResponse) GetConfidence() int32 {
	if x != nil {
		return x.Confidence
	}
	return 0
}

type ClientEvent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	EventType     ClientEvent_EventType  `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=captcha.v1.ClientEvent_EventType" json:"event_type,omitempty"`
	ChallengeId   string                 `protobuf:"bytes,2,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	Data          []byte                 `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	UserId        string                 `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClientEvent) Reset() {
	*x = ClientEvent{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientEvent) ProtoMessage() {}

func (x *ClientEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ClientEvent) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{4}
}

func (x *ClientEvent) GetEventType() ClientEvent_EventType {
	if x != nil {
		return x.EventType
	}
	return ClientEvent_FRONTEND_EVENT
}

func (x *ClientEvent) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *ClientEvent) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ClientEvent) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ServerEvent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Event         isServerEvent_Event    `protobuf_oneof:"event"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ServerEvent) Reset() {
	*x = ServerEvent{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEvent) ProtoMessage() {}

func (x *ServerEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ServerEvent) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{5}
}

func (x *ServerEvent) GetEvent() isServerEvent_Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *ServerEvent) GetResult() *ServerEvent_ChallengeResult {
	if x != nil {
		if x, ok := x.Event.(*ServerEvent_Result); ok {
			return x.Result
		}
	}
	return nil
}

func (x *ServerEvent) GetClientJs() *ServerEvent_RunClientJS {
	if x != nil {
		if x, ok := x.Event.(*ServerEvent_ClientJs); ok {
			return x.ClientJs
		}
	}
	return nil
}

func (x *ServerEvent) GetClientData() *ServerEvent_SendClientData {
	if x != nil {
		if x, ok := x.Event.(*ServerEvent_ClientData); ok {
			return x.ClientData
		}
	}
	return nil
}

type isServerEvent_Event interface {
	isServerEvent_Event()
}

type ServerEvent_Result struct {
	Result *ServerEvent_ChallengeResult `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}

type ServerEvent_ClientJs struct {
	ClientJs *ServerEvent_RunClientJS `protobuf:"bytes,2,opt,name=client_js,json=clientJs,proto3,oneof"`
}

type ServerEvent_ClientData struct {
	ClientData *ServerEvent_SendClientData `protobuf:"bytes,3,opt,name=client_data,json=clientData,proto3,oneof"`
}

func (*ServerEvent_Result) isServerEvent_Event() {}

func (*ServerEvent_ClientJs) isServerEvent_Event() {}

func (*ServerEvent_ClientData) isServerEvent_Event() {}

type ServerEvent_ChallengeResult struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	ChallengeId       string                 `protobuf:"bytes,1,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	ConfidencePercent int32                  `protobuf:"varint,2,opt,name=confidence_percent,json=confidencePercent,proto3" json:"confidence_percent,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *ServerEvent_ChallengeResult) Reset() {
	*x = ServerEvent_ChallengeResult{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerEvent_ChallengeResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEvent_ChallengeResult) ProtoMessage() {}

func (x *ServerEvent_ChallengeResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ServerEvent_ChallengeResult) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{5, 0}
}

func (x *ServerEvent_ChallengeResult) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *ServerEvent_ChallengeResult) GetConfidencePercent() int32 {
	if x != nil {
		return x.ConfidencePercent
	}
	return 0
}

type ServerEvent_RunClientJS struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChallengeId   string                 `protobuf:"bytes,1,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	JsCode        string                 `protobuf:"bytes,2,opt,name=js_code,json=jsCode,proto3" json:"js_code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ServerEvent_RunClientJS) Reset() {
	*x = ServerEvent_RunClientJS{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerEvent_RunClientJS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEvent_RunClientJS) ProtoMessage() {}

func (x *ServerEvent_RunClientJS) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ServerEvent_RunClientJS) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{5, 1}
}

func (x *ServerEvent_RunClientJS) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *ServerEvent_RunClientJS) GetJsCode() string {
	if x != nil {
		return x.JsCode
	}
	return ""
}

type ServerEvent_SendClientData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChallengeId   string                 `protobuf:"bytes,1,opt,name=challenge_id,json=challengeId,proto3" json:"challenge_id,omitempty"`
	Data          []byte                 `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ServerEvent_SendClientData) Reset() {
	*x = ServerEvent_SendClientData{}
	mi := &file_proto_captcha_captcha_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerEvent_SendClientData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEvent_SendClientData) ProtoMessage() {}

func (x *ServerEvent_SendClientData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_captcha_captcha_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ServerEvent_SendClientData) Descriptor() ([]byte, []int) {
	return file_proto_captcha_captcha_proto_rawDescGZIP(), []int{5, 2}
}

func (x *ServerEvent_SendClientData) GetChallengeId() string {
	if x != nil {
		return x.ChallengeId
	}
	return ""
}

func (x *ServerEvent_SendClientData) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_captcha_captcha_proto protoreflect.FileDescriptor

const file_proto_captcha_captcha_proto_rawDesc = "" +
	"\n" +
	"\x1bproto/captcha/captcha.proto\x12\n" +
	"captcha.v1\"K\n" +
	"\x10ChallengeRequest\x12\x1e\n" +
	"\n" +
	"complexity\x18\x01 \x01(\x05R\n" +
	"complexity\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\"J\n" +
	"\x11ChallengeResponse\x12!\n" +
	"\fchallenge_id\x18\x01 \x01(\tR\vchallengeId\x12\x12\n" +
	"\x04html\x18\x02 \x01(\tR\x04html\"L\n" +
	"\x0fValidateRequest\x12!\n" +
	"\fchallenge_id\x18\x01 \x01(\tR\vchallengeId\x12\x16\n" +
	"\x06answer\x18\x02 \x01(\tR\x06answer\"H\n" +
	"\x10ValidateResponse\x12\x14\n" +
	"\x05valid\x18\x01 \x01(\bR\x05valid\x12\x1e\n" +
	"\n" +
	"confidence\x18\x02 \x01(\x05R\n" +
	"confidence\"\xeb\x01\n" +
	"\vClientEvent\x12@\n" +
	"\n" +
	"event_type\x18\x01 \x01(\x0e2!.captcha.v1.ClientEvent.EventTypeR\teventType\x12!\n" +
	"\fchallenge_id\x18\x02 \x01(\tR\vchallengeId\x12\x12\n" +
	"\x04data\x18\x03 \x01(\fR\x04data\x12\x17\n" +
	"\auser_id\x18\x04 \x01(\tR\x06userId\"J\n" +
	"\tEventType\x12\x12\n" +
	"\x0eFRONTEND_EVENT\x10\x00\x12\x15\n" +
	"\x11CONNECTION_CLOSED\x10\x01\x12\x12\n" +
	"\x0eBALANCER_EVENT\x10\x02\"\xe1\x03\n" +
	"\vServerEvent\x12A\n" +
	"\x06result\x18\x01 \x01(\v2'.captcha.v1.ServerEvent.ChallengeResultH\x00R\x06result\x12B\n" +
	"\tclient_js\x18\x02 \x01(\v2#.captcha.v1.ServerEvent.RunClientJSH\x00R\bclientJs\x12I\n" +
	"\vclient_data\x18\x03 \x01(\v2&.captcha.v1.ServerEvent.SendClientDataH\x00R\n" +
	"clientData\x1ac\n" +
	"\x0fChallengeResult\x12!\n" +
	"\fchallenge_id\x18\x01 \x01(\tR\vchallengeId\x12-\n" +
	"\x12confidence_percent\x18\x02 \x01(\x05R\x11confidencePercent\x1aI\n" +
	"\vRunClientJS\x12!\n" +
	"\fchallenge_id\x18\x01 \x01(\tR\vchallengeId\x12\x17\n" +
	"\ajs_code\x18\x02 \x01(\tR\x06jsCode\x1aG\n" +
	"\x0eSendClientData\x12!\n" +
	"\fchallenge_id\x18\x01 \x01(\tR\vchallengeId\x12\x12\n" +
	"\x04data\x18\x02 \x01(\fR\x04dataB\a\n" +
	"\x05event2\xfc\x01\n" +
	"\x0eCaptchaService\x12M\n" +
	"\fNewChallenge\x12\x1c.captcha.v1.ChallengeRequest\x1a\x1d.captcha.v1.ChallengeResponse\"\x00\x12P\n" +
	"\x11ValidateChallenge\x12\x1b.captcha.v1.ValidateRequest\x1a\x1c.captcha.v1.ValidateResponse\"\x00\x12I\n" +
	"\x0fMakeEventStream\x12\x17.captcha.v1.ClientEvent\x1a\x17.captcha.v1.ServerEvent\"\x00(\x010\x01B&Z$captcha-service/gen/proto/captcha/v1b\x06proto3"

var (
	file_proto_captcha_captcha_proto_rawDescOnce sync.Once
	file_proto_captcha_captcha_proto_rawDescData []byte
)

func file_proto_captcha_captcha_proto_rawDescGZIP() []byte {
	file_proto_captcha_captcha_proto_rawDescOnce.Do(func() {
		file_proto_captcha_captcha_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_captcha_captcha_proto_rawDesc), len(file_proto_captcha_captcha_proto_rawDesc)))
	})
	return file_proto_captcha_captcha_proto_rawDescData
}

var file_proto_captcha_captcha_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_captcha_captcha_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_captcha_captcha_proto_goTypes = []any{
	(ClientEvent_EventType)(0),          // 0: captcha.v1.ClientEvent.EventType
	(*ChallengeRequest)(nil),            // 1: captcha.v1.ChallengeRequest
	(*ChallengeResponse)(nil),           // 2: captcha.v1.ChallengeResponse
	(*ValidateRequest)(nil),             // 3: captcha.v1.ValidateRequest
	(*ValidateResponse)(nil),            // 4: captcha.v1.ValidateResponse
	(*ClientEvent)(nil),                 // 5: captcha.v1.ClientEvent
	(*ServerEvent)(nil),                 // 6: captcha.v1.ServerEvent
	(*ServerEvent_ChallengeResult)(nil), // 7: captcha.v1.ServerEvent.ChallengeResult
	(*ServerEvent_RunClientJS)(nil),     // 8: captcha.v1.ServerEvent.RunClientJS
	(*ServerEvent_SendClientData)(nil),  // 9: captcha.v1.ServerEvent.SendClientData
}
var file_proto_captcha_captcha_proto_depIdxs = []int32{
	0, // 0: captcha.v1.ClientEvent.event_type:type_name -> captcha.v1.ClientEvent.EventType
	7, // 1: captcha.v1.ServerEvent.result:type_name -> captcha.v1.ServerEvent.ChallengeResult
	8, // 2: captcha.v1.ServerEvent.client_js:type_name -> captcha.v1.ServerEvent.RunClientJS
	9, // 3: captcha.v1.ServerEvent.client_data:type_name -> captcha.v1.ServerEvent.SendClientData
	1, // 4: captcha.v1.CaptchaService.NewChallenge:input_type -> captcha.v1.ChallengeRequest
	3, // 5: captcha.v1.CaptchaService.ValidateChallenge:input_type -> captcha.v1.ValidateRequest
	5, // 6: captcha.v1.CaptchaService.MakeEventStream:input_type -> captcha.v1.ClientEvent
	2, // 7: captcha.v1.CaptchaService.NewChallenge:output_type -> captcha.v1.ChallengeResponse
	4, // 8: captcha.v1.CaptchaService.ValidateChallenge:output_type -> captcha.v1.ValidateResponse
	6, // 9: captcha.v1.CaptchaService.MakeEventStream:output_type -> captcha.v1.ServerEvent
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_captcha_captcha_proto_init() }
func file_proto_captcha_captcha_proto_init() {
	if File_proto_captcha_captcha_proto != nil {
		return
	}
	file_proto_captcha_captcha_proto_msgTypes[5].OneofWrappers = []any{
		(*ServerEvent_Result)(nil),
		(*ServerEvent_ClientJs)(nil),
		(*ServerEvent_ClientData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_captcha_captcha_proto_rawDesc), len(file_proto_captcha_captcha_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_captcha_captcha_proto_goTypes,
		DependencyIndexes: file_proto_captcha_captcha_proto_depIdxs,
		EnumInfos:         file_proto_captcha_captcha_proto_enumTypes,
		MessageInfos:      file_proto_captcha_captcha_proto_msgTypes,
	}.Build()
	File_proto_captcha_captcha_proto = out.File
	file_proto_captcha_captcha_proto_goTypes = nil
	file_proto_captcha_captcha_proto_depIdxs = nil
}
