// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: commander_attach.proto

package dhctl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommanderAttachRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*CommanderAttachRequest_Start
	//	*CommanderAttachRequest_Continue
	Message isCommanderAttachRequest_Message `protobuf_oneof:"message"`
}

func (x *CommanderAttachRequest) Reset() {
	*x = CommanderAttachRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachRequest) ProtoMessage() {}

func (x *CommanderAttachRequest) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachRequest.ProtoReflect.Descriptor instead.
func (*CommanderAttachRequest) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{0}
}

func (m *CommanderAttachRequest) GetMessage() isCommanderAttachRequest_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *CommanderAttachRequest) GetStart() *CommanderAttachStart {
	if x, ok := x.GetMessage().(*CommanderAttachRequest_Start); ok {
		return x.Start
	}
	return nil
}

func (x *CommanderAttachRequest) GetContinue() *CommanderAttachContinue {
	if x, ok := x.GetMessage().(*CommanderAttachRequest_Continue); ok {
		return x.Continue
	}
	return nil
}

type isCommanderAttachRequest_Message interface {
	isCommanderAttachRequest_Message()
}

type CommanderAttachRequest_Start struct {
	Start *CommanderAttachStart `protobuf:"bytes,1,opt,name=start,proto3,oneof"`
}

type CommanderAttachRequest_Continue struct {
	Continue *CommanderAttachContinue `protobuf:"bytes,2,opt,name=continue,proto3,oneof"`
}

func (*CommanderAttachRequest_Start) isCommanderAttachRequest_Message() {}

func (*CommanderAttachRequest_Continue) isCommanderAttachRequest_Message() {}

type CommanderAttachResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*CommanderAttachResponse_Result
	//	*CommanderAttachResponse_PhaseEnd
	//	*CommanderAttachResponse_Logs
	Message isCommanderAttachResponse_Message `protobuf_oneof:"message"`
}

func (x *CommanderAttachResponse) Reset() {
	*x = CommanderAttachResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachResponse) ProtoMessage() {}

func (x *CommanderAttachResponse) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachResponse.ProtoReflect.Descriptor instead.
func (*CommanderAttachResponse) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{1}
}

func (m *CommanderAttachResponse) GetMessage() isCommanderAttachResponse_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *CommanderAttachResponse) GetResult() *CommanderAttachResult {
	if x, ok := x.GetMessage().(*CommanderAttachResponse_Result); ok {
		return x.Result
	}
	return nil
}

func (x *CommanderAttachResponse) GetPhaseEnd() *CommanderAttachPhaseEnd {
	if x, ok := x.GetMessage().(*CommanderAttachResponse_PhaseEnd); ok {
		return x.PhaseEnd
	}
	return nil
}

func (x *CommanderAttachResponse) GetLogs() *Logs {
	if x, ok := x.GetMessage().(*CommanderAttachResponse_Logs); ok {
		return x.Logs
	}
	return nil
}

type isCommanderAttachResponse_Message interface {
	isCommanderAttachResponse_Message()
}

type CommanderAttachResponse_Result struct {
	Result *CommanderAttachResult `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}

type CommanderAttachResponse_PhaseEnd struct {
	PhaseEnd *CommanderAttachPhaseEnd `protobuf:"bytes,2,opt,name=phase_end,json=phaseEnd,proto3,oneof"`
}

type CommanderAttachResponse_Logs struct {
	Logs *Logs `protobuf:"bytes,3,opt,name=logs,proto3,oneof"`
}

func (*CommanderAttachResponse_Result) isCommanderAttachResponse_Message() {}

func (*CommanderAttachResponse_PhaseEnd) isCommanderAttachResponse_Message() {}

func (*CommanderAttachResponse_Logs) isCommanderAttachResponse_Message() {}

type CommanderAttachStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnectionConfig  string                       `protobuf:"bytes,1,opt,name=connection_config,json=connectionConfig,proto3" json:"connection_config,omitempty"`
	ScanOnly          *bool                        `protobuf:"varint,2,opt,name=scan_only,json=scanOnly,proto3,oneof" json:"scan_only,omitempty"`
	ResourcesTemplate string                       `protobuf:"bytes,3,opt,name=resources_template,json=resourcesTemplate,proto3" json:"resources_template,omitempty"`
	ResourcesValues   *structpb.Struct             `protobuf:"bytes,4,opt,name=resources_values,json=resourcesValues,proto3" json:"resources_values,omitempty"`
	Options           *CommanderAttachStartOptions `protobuf:"bytes,5,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *CommanderAttachStart) Reset() {
	*x = CommanderAttachStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachStart) ProtoMessage() {}

func (x *CommanderAttachStart) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachStart.ProtoReflect.Descriptor instead.
func (*CommanderAttachStart) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{2}
}

func (x *CommanderAttachStart) GetConnectionConfig() string {
	if x != nil {
		return x.ConnectionConfig
	}
	return ""
}

func (x *CommanderAttachStart) GetScanOnly() bool {
	if x != nil && x.ScanOnly != nil {
		return *x.ScanOnly
	}
	return false
}

func (x *CommanderAttachStart) GetResourcesTemplate() string {
	if x != nil {
		return x.ResourcesTemplate
	}
	return ""
}

func (x *CommanderAttachStart) GetResourcesValues() *structpb.Struct {
	if x != nil {
		return x.ResourcesValues
	}
	return nil
}

func (x *CommanderAttachStart) GetOptions() *CommanderAttachStartOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

type CommanderAttachPhaseEnd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CompletedPhase      string            `protobuf:"bytes,1,opt,name=completed_phase,json=completedPhase,proto3" json:"completed_phase,omitempty"`
	CompletedPhaseState map[string][]byte `protobuf:"bytes,2,rep,name=completed_phase_state,json=completedPhaseState,proto3" json:"completed_phase_state,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CompletedPhaseData  string            `protobuf:"bytes,3,opt,name=completed_phase_data,json=completedPhaseData,proto3" json:"completed_phase_data,omitempty"`
	NextPhase           string            `protobuf:"bytes,4,opt,name=next_phase,json=nextPhase,proto3" json:"next_phase,omitempty"`
	NextPhaseCritical   bool              `protobuf:"varint,5,opt,name=next_phase_critical,json=nextPhaseCritical,proto3" json:"next_phase_critical,omitempty"`
}

func (x *CommanderAttachPhaseEnd) Reset() {
	*x = CommanderAttachPhaseEnd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachPhaseEnd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachPhaseEnd) ProtoMessage() {}

func (x *CommanderAttachPhaseEnd) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachPhaseEnd.ProtoReflect.Descriptor instead.
func (*CommanderAttachPhaseEnd) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{3}
}

func (x *CommanderAttachPhaseEnd) GetCompletedPhase() string {
	if x != nil {
		return x.CompletedPhase
	}
	return ""
}

func (x *CommanderAttachPhaseEnd) GetCompletedPhaseState() map[string][]byte {
	if x != nil {
		return x.CompletedPhaseState
	}
	return nil
}

func (x *CommanderAttachPhaseEnd) GetCompletedPhaseData() string {
	if x != nil {
		return x.CompletedPhaseData
	}
	return ""
}

func (x *CommanderAttachPhaseEnd) GetNextPhase() string {
	if x != nil {
		return x.NextPhase
	}
	return ""
}

func (x *CommanderAttachPhaseEnd) GetNextPhaseCritical() bool {
	if x != nil {
		return x.NextPhaseCritical
	}
	return false
}

type CommanderAttachContinue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Continue Continue `protobuf:"varint,1,opt,name=continue,proto3,enum=dhctl.Continue" json:"continue,omitempty"`
	Err      string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *CommanderAttachContinue) Reset() {
	*x = CommanderAttachContinue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachContinue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachContinue) ProtoMessage() {}

func (x *CommanderAttachContinue) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachContinue.ProtoReflect.Descriptor instead.
func (*CommanderAttachContinue) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{4}
}

func (x *CommanderAttachContinue) GetContinue() Continue {
	if x != nil {
		return x.Continue
	}
	return Continue_CONTINUE_UNSPECIFIED
}

func (x *CommanderAttachContinue) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type CommanderAttachStartOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommanderMode    bool                 `protobuf:"varint,1,opt,name=commander_mode,json=commanderMode,proto3" json:"commander_mode,omitempty"`
	CommanderUuid    string               `protobuf:"bytes,2,opt,name=commander_uuid,json=commanderUuid,proto3" json:"commander_uuid,omitempty"`
	LogWidth         int32                `protobuf:"varint,3,opt,name=log_width,json=logWidth,proto3" json:"log_width,omitempty"`
	ResourcesTimeout *durationpb.Duration `protobuf:"bytes,4,opt,name=resources_timeout,json=resourcesTimeout,proto3" json:"resources_timeout,omitempty"`
	DeckhouseTimeout *durationpb.Duration `protobuf:"bytes,5,opt,name=deckhouse_timeout,json=deckhouseTimeout,proto3" json:"deckhouse_timeout,omitempty"`
}

func (x *CommanderAttachStartOptions) Reset() {
	*x = CommanderAttachStartOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachStartOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachStartOptions) ProtoMessage() {}

func (x *CommanderAttachStartOptions) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachStartOptions.ProtoReflect.Descriptor instead.
func (*CommanderAttachStartOptions) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{5}
}

func (x *CommanderAttachStartOptions) GetCommanderMode() bool {
	if x != nil {
		return x.CommanderMode
	}
	return false
}

func (x *CommanderAttachStartOptions) GetCommanderUuid() string {
	if x != nil {
		return x.CommanderUuid
	}
	return ""
}

func (x *CommanderAttachStartOptions) GetLogWidth() int32 {
	if x != nil {
		return x.LogWidth
	}
	return 0
}

func (x *CommanderAttachStartOptions) GetResourcesTimeout() *durationpb.Duration {
	if x != nil {
		return x.ResourcesTimeout
	}
	return nil
}

func (x *CommanderAttachStartOptions) GetDeckhouseTimeout() *durationpb.Duration {
	if x != nil {
		return x.DeckhouseTimeout
	}
	return nil
}

type CommanderAttachResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State  string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	Err    string `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *CommanderAttachResult) Reset() {
	*x = CommanderAttachResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commander_attach_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommanderAttachResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommanderAttachResult) ProtoMessage() {}

func (x *CommanderAttachResult) ProtoReflect() protoreflect.Message {
	mi := &file_commander_attach_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommanderAttachResult.ProtoReflect.Descriptor instead.
func (*CommanderAttachResult) Descriptor() ([]byte, []int) {
	return file_commander_attach_proto_rawDescGZIP(), []int{6}
}

func (x *CommanderAttachResult) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *CommanderAttachResult) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *CommanderAttachResult) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_commander_attach_proto protoreflect.FileDescriptor

var file_commander_attach_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x61, 0x74, 0x74, 0x61,
	0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x96, 0x01, 0x0a, 0x16,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x3c, 0x0a, 0x08, 0x63,
	0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41,
	0x74, 0x74, 0x61, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x48, 0x00, 0x52,
	0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xbe, 0x01, 0x0a, 0x17, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x36, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x48, 0x00,
	0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x3d, 0x0a, 0x09, 0x70, 0x68, 0x61, 0x73,
	0x65, 0x5f, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x64, 0x68,
	0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x74, 0x74,
	0x61, 0x63, 0x68, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x08, 0x70,
	0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x12, 0x21, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x4c, 0x6f,
	0x67, 0x73, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xa4, 0x02, 0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2b,
	0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x20, 0x0a, 0x09, 0x73,
	0x63, 0x61, 0x6e, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00,
	0x52, 0x08, 0x73, 0x63, 0x61, 0x6e, 0x4f, 0x6e, 0x6c, 0x79, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a,
	0x12, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x42, 0x0a, 0x10,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52,
	0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x12, 0x3c, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x22, 0xf8, 0x02, 0x0a,
	0x17, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68,
	0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73,
	0x65, 0x12, 0x6b, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x70,
	0x68, 0x61, 0x73, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x37, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64,
	0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x13, 0x63, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x30,
	0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x68, 0x61, 0x73,
	0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x63, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12,
	0x2e, 0x0a, 0x13, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x72,
	0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x6e, 0x65,
	0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x43, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x1a,
	0x46, 0x0a, 0x18, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x58, 0x0a, 0x17, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e,
	0x75, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x69, 0x6e, 0x75, 0x65, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72,
	0x72, 0x22, 0x98, 0x02, 0x0a, 0x1b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41,
	0x74, 0x74, 0x61, 0x63, 0x68, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x6d,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x57, 0x69, 0x64, 0x74, 0x68, 0x12, 0x46, 0x0a, 0x11,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x54, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x12, 0x46, 0x0a, 0x11, 0x64, 0x65, 0x63, 0x6b, 0x68, 0x6f, 0x75, 0x73,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x64, 0x65, 0x63, 0x6b,
	0x68, 0x6f, 0x75, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x57, 0x0a, 0x15,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x65, 0x72, 0x72, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x62, 0x2f, 0x64, 0x68, 0x63, 0x74,
	0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_commander_attach_proto_rawDescOnce sync.Once
	file_commander_attach_proto_rawDescData = file_commander_attach_proto_rawDesc
)

func file_commander_attach_proto_rawDescGZIP() []byte {
	file_commander_attach_proto_rawDescOnce.Do(func() {
		file_commander_attach_proto_rawDescData = protoimpl.X.CompressGZIP(file_commander_attach_proto_rawDescData)
	})
	return file_commander_attach_proto_rawDescData
}

var file_commander_attach_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_commander_attach_proto_goTypes = []interface{}{
	(*CommanderAttachRequest)(nil),      // 0: dhctl.CommanderAttachRequest
	(*CommanderAttachResponse)(nil),     // 1: dhctl.CommanderAttachResponse
	(*CommanderAttachStart)(nil),        // 2: dhctl.CommanderAttachStart
	(*CommanderAttachPhaseEnd)(nil),     // 3: dhctl.CommanderAttachPhaseEnd
	(*CommanderAttachContinue)(nil),     // 4: dhctl.CommanderAttachContinue
	(*CommanderAttachStartOptions)(nil), // 5: dhctl.CommanderAttachStartOptions
	(*CommanderAttachResult)(nil),       // 6: dhctl.CommanderAttachResult
	nil,                                 // 7: dhctl.CommanderAttachPhaseEnd.CompletedPhaseStateEntry
	(*Logs)(nil),                        // 8: dhctl.Logs
	(*structpb.Struct)(nil),             // 9: google.protobuf.Struct
	(Continue)(0),                       // 10: dhctl.Continue
	(*durationpb.Duration)(nil),         // 11: google.protobuf.Duration
}
var file_commander_attach_proto_depIdxs = []int32{
	2,  // 0: dhctl.CommanderAttachRequest.start:type_name -> dhctl.CommanderAttachStart
	4,  // 1: dhctl.CommanderAttachRequest.continue:type_name -> dhctl.CommanderAttachContinue
	6,  // 2: dhctl.CommanderAttachResponse.result:type_name -> dhctl.CommanderAttachResult
	3,  // 3: dhctl.CommanderAttachResponse.phase_end:type_name -> dhctl.CommanderAttachPhaseEnd
	8,  // 4: dhctl.CommanderAttachResponse.logs:type_name -> dhctl.Logs
	9,  // 5: dhctl.CommanderAttachStart.resources_values:type_name -> google.protobuf.Struct
	5,  // 6: dhctl.CommanderAttachStart.options:type_name -> dhctl.CommanderAttachStartOptions
	7,  // 7: dhctl.CommanderAttachPhaseEnd.completed_phase_state:type_name -> dhctl.CommanderAttachPhaseEnd.CompletedPhaseStateEntry
	10, // 8: dhctl.CommanderAttachContinue.continue:type_name -> dhctl.Continue
	11, // 9: dhctl.CommanderAttachStartOptions.resources_timeout:type_name -> google.protobuf.Duration
	11, // 10: dhctl.CommanderAttachStartOptions.deckhouse_timeout:type_name -> google.protobuf.Duration
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_commander_attach_proto_init() }
func file_commander_attach_proto_init() {
	if File_commander_attach_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_commander_attach_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commander_attach_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commander_attach_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachStart); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commander_attach_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachPhaseEnd); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commander_attach_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachContinue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commander_attach_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachStartOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commander_attach_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommanderAttachResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_commander_attach_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*CommanderAttachRequest_Start)(nil),
		(*CommanderAttachRequest_Continue)(nil),
	}
	file_commander_attach_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*CommanderAttachResponse_Result)(nil),
		(*CommanderAttachResponse_PhaseEnd)(nil),
		(*CommanderAttachResponse_Logs)(nil),
	}
	file_commander_attach_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_commander_attach_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commander_attach_proto_goTypes,
		DependencyIndexes: file_commander_attach_proto_depIdxs,
		MessageInfos:      file_commander_attach_proto_msgTypes,
	}.Build()
	File_commander_attach_proto = out.File
	file_commander_attach_proto_rawDesc = nil
	file_commander_attach_proto_goTypes = nil
	file_commander_attach_proto_depIdxs = nil
}
