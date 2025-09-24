package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterInstanceRequest_EventType int32

const (
	RegisterInstanceRequest_UNKNOWN   RegisterInstanceRequest_EventType = 0
	RegisterInstanceRequest_READY     RegisterInstanceRequest_EventType = 1
	RegisterInstanceRequest_NOT_READY RegisterInstanceRequest_EventType = 2
	RegisterInstanceRequest_STOPPED   RegisterInstanceRequest_EventType = 3
)

var (
	RegisterInstanceRequest_EventType_name = map[int32]string{
		0: "UNKNOWN",
		1: "READY",
		2: "NOT_READY",
		3: "STOPPED",
	}
	RegisterInstanceRequest_EventType_value = map[string]int32{
		"UNKNOWN":   0,
		"READY":     1,
		"NOT_READY": 2,
		"STOPPED":   3,
	}
)

func (x RegisterInstanceRequest_EventType) Enum() *RegisterInstanceRequest_EventType {
	p := new(RegisterInstanceRequest_EventType)
	*p = x
	return p
}

func (x RegisterInstanceRequest_EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RegisterInstanceRequest_EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_balancer_balancer_proto_enumTypes[0].Descriptor()
}

func (RegisterInstanceRequest_EventType) Type() protoreflect.EnumType {
	return &file_proto_balancer_balancer_proto_enumTypes[0]
}

func (x RegisterInstanceRequest_EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

func (RegisterInstanceRequest_EventType) EnumDescriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{0, 0}
}

type RegisterInstanceResponse_Status int32

const (
	RegisterInstanceResponse_SUCCESS RegisterInstanceResponse_Status = 0
	RegisterInstanceResponse_ERROR   RegisterInstanceResponse_Status = 1
)

var (
	RegisterInstanceResponse_Status_name = map[int32]string{
		0: "SUCCESS",
		1: "ERROR",
	}
	RegisterInstanceResponse_Status_value = map[string]int32{
		"SUCCESS": 0,
		"ERROR":   1,
	}
)

func (x RegisterInstanceResponse_Status) Enum() *RegisterInstanceResponse_Status {
	p := new(RegisterInstanceResponse_Status)
	*p = x
	return p
}

func (x RegisterInstanceResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RegisterInstanceResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_balancer_balancer_proto_enumTypes[1].Descriptor()
}

func (RegisterInstanceResponse_Status) Type() protoreflect.EnumType {
	return &file_proto_balancer_balancer_proto_enumTypes[1]
}

func (x RegisterInstanceResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

func (RegisterInstanceResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{1, 0}
}

type BlockUserResponse_Status int32

const (
	BlockUserResponse_SUCCESS BlockUserResponse_Status = 0
	BlockUserResponse_ERROR   BlockUserResponse_Status = 1
)

var (
	BlockUserResponse_Status_name = map[int32]string{
		0: "SUCCESS",
		1: "ERROR",
	}
	BlockUserResponse_Status_value = map[string]int32{
		"SUCCESS": 0,
		"ERROR":   1,
	}
)

func (x BlockUserResponse_Status) Enum() *BlockUserResponse_Status {
	p := new(BlockUserResponse_Status)
	*p = x
	return p
}

func (x BlockUserResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BlockUserResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_balancer_balancer_proto_enumTypes[2].Descriptor()
}

func (BlockUserResponse_Status) Type() protoreflect.EnumType {
	return &file_proto_balancer_balancer_proto_enumTypes[2]
}

func (x BlockUserResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

func (BlockUserResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{5, 0}
}

type RegisterInstanceRequest struct {
	state         protoimpl.MessageState            `protogen:"open.v1"`
	EventType     RegisterInstanceRequest_EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=balancer.v1.RegisterInstanceRequest_EventType" json:"event_type,omitempty"`
	InstanceId    string                            `protobuf:"bytes,2,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	ChallengeType string                            `protobuf:"bytes,3,opt,name=challenge_type,json=challengeType,proto3" json:"challenge_type,omitempty"`
	Host          string                            `protobuf:"bytes,4,opt,name=host,proto3" json:"host,omitempty"`
	PortNumber    int32                             `protobuf:"varint,5,opt,name=port_number,json=portNumber,proto3" json:"port_number,omitempty"`
	Timestamp     int64                             `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterInstanceRequest) Reset() {
	*x = RegisterInstanceRequest{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterInstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterInstanceRequest) ProtoMessage() {}

func (x *RegisterInstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*RegisterInstanceRequest) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterInstanceRequest) GetEventType() RegisterInstanceRequest_EventType {
	if x != nil {
		return x.EventType
	}
	return RegisterInstanceRequest_UNKNOWN
}

func (x *RegisterInstanceRequest) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

func (x *RegisterInstanceRequest) GetChallengeType() string {
	if x != nil {
		return x.ChallengeType
	}
	return ""
}

func (x *RegisterInstanceRequest) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *RegisterInstanceRequest) GetPortNumber() int32 {
	if x != nil {
		return x.PortNumber
	}
	return 0
}

func (x *RegisterInstanceRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type RegisterInstanceResponse struct {
	state         protoimpl.MessageState          `protogen:"open.v1"`
	Status        RegisterInstanceResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=balancer.v1.RegisterInstanceResponse_Status" json:"status,omitempty"`
	Message       string                          `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterInstanceResponse) Reset() {
	*x = RegisterInstanceResponse{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterInstanceResponse) ProtoMessage() {}

func (x *RegisterInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*RegisterInstanceResponse) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterInstanceResponse) GetStatus() RegisterInstanceResponse_Status {
	if x != nil {
		return x.Status
	}
	return RegisterInstanceResponse_SUCCESS
}

func (x *RegisterInstanceResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type CheckUserBlockedRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckUserBlockedRequest) Reset() {
	*x = CheckUserBlockedRequest{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckUserBlockedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckUserBlockedRequest) ProtoMessage() {}

func (x *CheckUserBlockedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CheckUserBlockedRequest) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{2}
}

func (x *CheckUserBlockedRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type CheckUserBlockedResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	IsBlocked     bool                   `protobuf:"varint,1,opt,name=is_blocked,json=isBlocked,proto3" json:"is_blocked,omitempty"`
	Reason        string                 `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	BlockedUntil  int64                  `protobuf:"varint,3,opt,name=blocked_until,json=blockedUntil,proto3" json:"blocked_until,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckUserBlockedResponse) Reset() {
	*x = CheckUserBlockedResponse{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckUserBlockedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckUserBlockedResponse) ProtoMessage() {}

func (x *CheckUserBlockedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CheckUserBlockedResponse) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{3}
}

func (x *CheckUserBlockedResponse) GetIsBlocked() bool {
	if x != nil {
		return x.IsBlocked
	}
	return false
}

func (x *CheckUserBlockedResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *CheckUserBlockedResponse) GetBlockedUntil() int64 {
	if x != nil {
		return x.BlockedUntil
	}
	return 0
}

type BlockUserRequest struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	UserId          string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	DurationMinutes int32                  `protobuf:"varint,2,opt,name=duration_minutes,json=durationMinutes,proto3" json:"duration_minutes,omitempty"`
	Reason          string                 `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *BlockUserRequest) Reset() {
	*x = BlockUserRequest{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockUserRequest) ProtoMessage() {}

func (x *BlockUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*BlockUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{4}
}

func (x *BlockUserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *BlockUserRequest) GetDurationMinutes() int32 {
	if x != nil {
		return x.DurationMinutes
	}
	return 0
}

func (x *BlockUserRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type BlockUserResponse struct {
	state         protoimpl.MessageState   `protogen:"open.v1"`
	Status        BlockUserResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=balancer.v1.BlockUserResponse_Status" json:"status,omitempty"`
	Message       string                   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BlockUserResponse) Reset() {
	*x = BlockUserResponse{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockUserResponse) ProtoMessage() {}

func (x *BlockUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*BlockUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{5}
}

func (x *BlockUserResponse) GetStatus() BlockUserResponse_Status {
	if x != nil {
		return x.Status
	}
	return BlockUserResponse_SUCCESS
}

func (x *BlockUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetInstancesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetInstancesRequest) Reset() {
	*x = GetInstancesRequest{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInstancesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInstancesRequest) ProtoMessage() {}

func (x *GetInstancesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetInstancesRequest) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{6}
}

type InstanceInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	InstanceId    string                 `protobuf:"bytes,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	ChallengeType string                 `protobuf:"bytes,2,opt,name=challenge_type,json=challengeType,proto3" json:"challenge_type,omitempty"`
	Host          string                 `protobuf:"bytes,3,opt,name=host,proto3" json:"host,omitempty"`
	PortNumber    int32                  `protobuf:"varint,4,opt,name=port_number,json=portNumber,proto3" json:"port_number,omitempty"`
	Status        string                 `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	LastSeen      int64                  `protobuf:"varint,6,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InstanceInfo) Reset() {
	*x = InstanceInfo{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InstanceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceInfo) ProtoMessage() {}

func (x *InstanceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*InstanceInfo) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{7}
}

func (x *InstanceInfo) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

func (x *InstanceInfo) GetChallengeType() string {
	if x != nil {
		return x.ChallengeType
	}
	return ""
}

func (x *InstanceInfo) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *InstanceInfo) GetPortNumber() int32 {
	if x != nil {
		return x.PortNumber
	}
	return 0
}

func (x *InstanceInfo) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *InstanceInfo) GetLastSeen() int64 {
	if x != nil {
		return x.LastSeen
	}
	return 0
}

type GetInstancesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Instances     []*InstanceInfo        `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
	Count         int32                  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetInstancesResponse) Reset() {
	*x = GetInstancesResponse{}
	mi := &file_proto_balancer_balancer_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInstancesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInstancesResponse) ProtoMessage() {}

func (x *GetInstancesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_balancer_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetInstancesResponse) Descriptor() ([]byte, []int) {
	return file_proto_balancer_balancer_proto_rawDescGZIP(), []int{8}
}

func (x *GetInstancesResponse) GetInstances() []*InstanceInfo {
	if x != nil {
		return x.Instances
	}
	return nil
}

func (x *GetInstancesResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_proto_balancer_balancer_proto protoreflect.FileDescriptor

const file_proto_balancer_balancer_proto_rawDesc = "" +
	"\n" +
	"\x1dproto/balancer/balancer.proto\x12\vbalancer.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\xc4\x02\n" +
	"\x17RegisterInstanceRequest\x12M\n" +
	"\n" +
	"event_type\x18\x01 \x01(\x0e2..balancer.v1.RegisterInstanceRequest.EventTypeR\teventType\x12\x1f\n" +
	"\vinstance_id\x18\x02 \x01(\tR\n" +
	"instanceId\x12%\n" +
	"\x0echallenge_type\x18\x03 \x01(\tR\rchallengeType\x12\x12\n" +
	"\x04host\x18\x04 \x01(\tR\x04host\x12\x1f\n" +
	"\vport_number\x18\x05 \x01(\x05R\n" +
	"portNumber\x12\x1c\n" +
	"\ttimestamp\x18\x06 \x01(\x03R\ttimestamp\"?\n" +
	"\tEventType\x12\v\n" +
	"\aUNKNOWN\x10\x00\x12\t\n" +
	"\x05READY\x10\x01\x12\r\n" +
	"\tNOT_READY\x10\x02\x12\v\n" +
	"\aSTOPPED\x10\x03\"\x9c\x01\n" +
	"\x18RegisterInstanceResponse\x12D\n" +
	"\x06status\x18\x01 \x01(\x0e2,.balancer.v1.RegisterInstanceResponse.StatusR\x06status\x12\x18\n" +
	"\amessage\x18\x03 \x01(\tR\amessage\" \n" +
	"\x06Status\x12\v\n" +
	"\aSUCCESS\x10\x00\x12\t\n" +
	"\x05ERROR\x10\x01\"2\n" +
	"\x17CheckUserBlockedRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\"v\n" +
	"\x18CheckUserBlockedResponse\x12\x1d\n" +
	"\n" +
	"is_blocked\x18\x01 \x01(\bR\tisBlocked\x12\x16\n" +
	"\x06reason\x18\x02 \x01(\tR\x06reason\x12#\n" +
	"\rblocked_until\x18\x03 \x01(\x03R\fblockedUntil\"n\n" +
	"\x10BlockUserRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12)\n" +
	"\x10duration_minutes\x18\x02 \x01(\x05R\x0fdurationMinutes\x12\x16\n" +
	"\x06reason\x18\x03 \x01(\tR\x06reason\"\x8e\x01\n" +
	"\x11BlockUserResponse\x12=\n" +
	"\x06status\x18\x01 \x01(\x0e2%.balancer.v1.BlockUserResponse.StatusR\x06status\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\" \n" +
	"\x06Status\x12\v\n" +
	"\aSUCCESS\x10\x00\x12\t\n" +
	"\x05ERROR\x10\x01\"\x15\n" +
	"\x13GetInstancesRequest\"\xc0\x01\n" +
	"\fInstanceInfo\x12\x1f\n" +
	"\vinstance_id\x18\x01 \x01(\tR\n" +
	"instanceId\x12%\n" +
	"\x0echallenge_type\x18\x02 \x01(\tR\rchallengeType\x12\x12\n" +
	"\x04host\x18\x03 \x01(\tR\x04host\x12\x1f\n" +
	"\vport_number\x18\x04 \x01(\x05R\n" +
	"portNumber\x12\x16\n" +
	"\x06status\x18\x05 \x01(\tR\x06status\x12\x1b\n" +
	"\tlast_seen\x18\x06 \x01(\x03R\blastSeen\"e\n" +
	"\x14GetInstancesResponse\x127\n" +
	"\tinstances\x18\x01 \x03(\v2\x19.balancer.v1.InstanceInfoR\tinstances\x12\x14\n" +
	"\x05count\x18\x02 \x01(\x05R\x05count2\x80\x03\n" +
	"\x0fBalancerService\x12e\n" +
	"\x10RegisterInstance\x12$.balancer.v1.RegisterInstanceRequest\x1a%.balancer.v1.RegisterInstanceResponse\"\x00(\x010\x01\x12a\n" +
	"\x10CheckUserBlocked\x12$.balancer.v1.CheckUserBlockedRequest\x1a%.balancer.v1.CheckUserBlockedResponse\"\x00\x12L\n" +
	"\tBlockUser\x12\x1d.balancer.v1.BlockUserRequest\x1a\x1e.balancer.v1.BlockUserResponse\"\x00\x12U\n" +
	"\fGetInstances\x12 .balancer.v1.GetInstancesRequest\x1a!.balancer.v1.GetInstancesResponse\"\x00B'Z%captcha-service/gen/proto/balancer/v1b\x06proto3"

var (
	file_proto_balancer_balancer_proto_rawDescOnce sync.Once
	file_proto_balancer_balancer_proto_rawDescData []byte
)

func file_proto_balancer_balancer_proto_rawDescGZIP() []byte {
	file_proto_balancer_balancer_proto_rawDescOnce.Do(func() {
		file_proto_balancer_balancer_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_balancer_balancer_proto_rawDesc), len(file_proto_balancer_balancer_proto_rawDesc)))
	})
	return file_proto_balancer_balancer_proto_rawDescData
}

var file_proto_balancer_balancer_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_proto_balancer_balancer_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_balancer_balancer_proto_goTypes = []any{
	(RegisterInstanceRequest_EventType)(0), // 0: balancer.v1.RegisterInstanceRequest.EventType
	(RegisterInstanceResponse_Status)(0),   // 1: balancer.v1.RegisterInstanceResponse.Status
	(BlockUserResponse_Status)(0),          // 2: balancer.v1.BlockUserResponse.Status
	(*RegisterInstanceRequest)(nil),        // 3: balancer.v1.RegisterInstanceRequest
	(*RegisterInstanceResponse)(nil),       // 4: balancer.v1.RegisterInstanceResponse
	(*CheckUserBlockedRequest)(nil),        // 5: balancer.v1.CheckUserBlockedRequest
	(*CheckUserBlockedResponse)(nil),       // 6: balancer.v1.CheckUserBlockedResponse
	(*BlockUserRequest)(nil),               // 7: balancer.v1.BlockUserRequest
	(*BlockUserResponse)(nil),              // 8: balancer.v1.BlockUserResponse
	(*GetInstancesRequest)(nil),            // 9: balancer.v1.GetInstancesRequest
	(*InstanceInfo)(nil),                   // 10: balancer.v1.InstanceInfo
	(*GetInstancesResponse)(nil),           // 11: balancer.v1.GetInstancesResponse
}
var file_proto_balancer_balancer_proto_depIdxs = []int32{
	0,  // 0: balancer.v1.RegisterInstanceRequest.event_type:type_name -> balancer.v1.RegisterInstanceRequest.EventType
	1,  // 1: balancer.v1.RegisterInstanceResponse.status:type_name -> balancer.v1.RegisterInstanceResponse.Status
	2,  // 2: balancer.v1.BlockUserResponse.status:type_name -> balancer.v1.BlockUserResponse.Status
	10, // 3: balancer.v1.GetInstancesResponse.instances:type_name -> balancer.v1.InstanceInfo
	3,  // 4: balancer.v1.BalancerService.RegisterInstance:input_type -> balancer.v1.RegisterInstanceRequest
	5,  // 5: balancer.v1.BalancerService.CheckUserBlocked:input_type -> balancer.v1.CheckUserBlockedRequest
	7,  // 6: balancer.v1.BalancerService.BlockUser:input_type -> balancer.v1.BlockUserRequest
	9,  // 7: balancer.v1.BalancerService.GetInstances:input_type -> balancer.v1.GetInstancesRequest
	4,  // 8: balancer.v1.BalancerService.RegisterInstance:output_type -> balancer.v1.RegisterInstanceResponse
	6,  // 9: balancer.v1.BalancerService.CheckUserBlocked:output_type -> balancer.v1.CheckUserBlockedResponse
	8,  // 10: balancer.v1.BalancerService.BlockUser:output_type -> balancer.v1.BlockUserResponse
	11, // 11: balancer.v1.BalancerService.GetInstances:output_type -> balancer.v1.GetInstancesResponse
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_balancer_balancer_proto_init() }
func file_proto_balancer_balancer_proto_init() {
	if File_proto_balancer_balancer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_balancer_balancer_proto_rawDesc), len(file_proto_balancer_balancer_proto_rawDesc)),
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_balancer_balancer_proto_goTypes,
		DependencyIndexes: file_proto_balancer_balancer_proto_depIdxs,
		EnumInfos:         file_proto_balancer_balancer_proto_enumTypes,
		MessageInfos:      file_proto_balancer_balancer_proto_msgTypes,
	}.Build()
	File_proto_balancer_balancer_proto = out.File
	file_proto_balancer_balancer_proto_goTypes = nil
	file_proto_balancer_balancer_proto_depIdxs = nil
}
