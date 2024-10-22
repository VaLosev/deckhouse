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
// source: bootstrap.proto

package dhctl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BootstrapRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*BootstrapRequest_Start
	//	*BootstrapRequest_Continue
	Message isBootstrapRequest_Message `protobuf_oneof:"message"`
}

func (x *BootstrapRequest) Reset() {
	*x = BootstrapRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapRequest) ProtoMessage() {}

func (x *BootstrapRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapRequest.ProtoReflect.Descriptor instead.
func (*BootstrapRequest) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{0}
}

func (m *BootstrapRequest) GetMessage() isBootstrapRequest_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *BootstrapRequest) GetStart() *BootstrapStart {
	if x, ok := x.GetMessage().(*BootstrapRequest_Start); ok {
		return x.Start
	}
	return nil
}

func (x *BootstrapRequest) GetContinue() *BootstrapContinue {
	if x, ok := x.GetMessage().(*BootstrapRequest_Continue); ok {
		return x.Continue
	}
	return nil
}

type isBootstrapRequest_Message interface {
	isBootstrapRequest_Message()
}

type BootstrapRequest_Start struct {
	Start *BootstrapStart `protobuf:"bytes,1,opt,name=start,proto3,oneof"`
}

type BootstrapRequest_Continue struct {
	Continue *BootstrapContinue `protobuf:"bytes,2,opt,name=continue,proto3,oneof"`
}

func (*BootstrapRequest_Start) isBootstrapRequest_Message() {}

func (*BootstrapRequest_Continue) isBootstrapRequest_Message() {}

type BootstrapResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*BootstrapResponse_Result
	//	*BootstrapResponse_PhaseEnd
	//	*BootstrapResponse_Logs
	Message isBootstrapResponse_Message `protobuf_oneof:"message"`
}

func (x *BootstrapResponse) Reset() {
	*x = BootstrapResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapResponse) ProtoMessage() {}

func (x *BootstrapResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapResponse.ProtoReflect.Descriptor instead.
func (*BootstrapResponse) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{1}
}

func (m *BootstrapResponse) GetMessage() isBootstrapResponse_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *BootstrapResponse) GetResult() *BootstrapResult {
	if x, ok := x.GetMessage().(*BootstrapResponse_Result); ok {
		return x.Result
	}
	return nil
}

func (x *BootstrapResponse) GetPhaseEnd() *BootstrapPhaseEnd {
	if x, ok := x.GetMessage().(*BootstrapResponse_PhaseEnd); ok {
		return x.PhaseEnd
	}
	return nil
}

func (x *BootstrapResponse) GetLogs() *Logs {
	if x, ok := x.GetMessage().(*BootstrapResponse_Logs); ok {
		return x.Logs
	}
	return nil
}

type isBootstrapResponse_Message interface {
	isBootstrapResponse_Message()
}

type BootstrapResponse_Result struct {
	Result *BootstrapResult `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}

type BootstrapResponse_PhaseEnd struct {
	PhaseEnd *BootstrapPhaseEnd `protobuf:"bytes,2,opt,name=phase_end,json=phaseEnd,proto3,oneof"`
}

type BootstrapResponse_Logs struct {
	Logs *Logs `protobuf:"bytes,3,opt,name=logs,proto3,oneof"`
}

func (*BootstrapResponse_Result) isBootstrapResponse_Message() {}

func (*BootstrapResponse_PhaseEnd) isBootstrapResponse_Message() {}

func (*BootstrapResponse_Logs) isBootstrapResponse_Message() {}

type BootstrapStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnectionConfig              string                 `protobuf:"bytes,1,opt,name=connection_config,json=connectionConfig,proto3" json:"connection_config,omitempty"`
	InitConfig                    string                 `protobuf:"bytes,2,opt,name=init_config,json=initConfig,proto3" json:"init_config,omitempty"`
	ClusterConfig                 string                 `protobuf:"bytes,3,opt,name=cluster_config,json=clusterConfig,proto3" json:"cluster_config,omitempty"`
	ProviderSpecificClusterConfig string                 `protobuf:"bytes,4,opt,name=provider_specific_cluster_config,json=providerSpecificClusterConfig,proto3" json:"provider_specific_cluster_config,omitempty"`
	InitResources                 string                 `protobuf:"bytes,5,opt,name=init_resources,json=initResources,proto3" json:"init_resources,omitempty"`
	Resources                     string                 `protobuf:"bytes,6,opt,name=resources,proto3" json:"resources,omitempty"`
	State                         string                 `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	PostBootstrapScript           string                 `protobuf:"bytes,8,opt,name=post_bootstrap_script,json=postBootstrapScript,proto3" json:"post_bootstrap_script,omitempty"`
	Options                       *BootstrapStartOptions `protobuf:"bytes,9,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *BootstrapStart) Reset() {
	*x = BootstrapStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapStart) ProtoMessage() {}

func (x *BootstrapStart) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapStart.ProtoReflect.Descriptor instead.
func (*BootstrapStart) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{2}
}

func (x *BootstrapStart) GetConnectionConfig() string {
	if x != nil {
		return x.ConnectionConfig
	}
	return ""
}

func (x *BootstrapStart) GetInitConfig() string {
	if x != nil {
		return x.InitConfig
	}
	return ""
}

func (x *BootstrapStart) GetClusterConfig() string {
	if x != nil {
		return x.ClusterConfig
	}
	return ""
}

func (x *BootstrapStart) GetProviderSpecificClusterConfig() string {
	if x != nil {
		return x.ProviderSpecificClusterConfig
	}
	return ""
}

func (x *BootstrapStart) GetInitResources() string {
	if x != nil {
		return x.InitResources
	}
	return ""
}

func (x *BootstrapStart) GetResources() string {
	if x != nil {
		return x.Resources
	}
	return ""
}

func (x *BootstrapStart) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *BootstrapStart) GetPostBootstrapScript() string {
	if x != nil {
		return x.PostBootstrapScript
	}
	return ""
}

func (x *BootstrapStart) GetOptions() *BootstrapStartOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

type BootstrapPhaseEnd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CompletedPhase      string            `protobuf:"bytes,1,opt,name=completed_phase,json=completedPhase,proto3" json:"completed_phase,omitempty"`
	CompletedPhaseState map[string][]byte `protobuf:"bytes,2,rep,name=completed_phase_state,json=completedPhaseState,proto3" json:"completed_phase_state,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NextPhase           string            `protobuf:"bytes,3,opt,name=next_phase,json=nextPhase,proto3" json:"next_phase,omitempty"`
	NextPhaseCritical   bool              `protobuf:"varint,4,opt,name=next_phase_critical,json=nextPhaseCritical,proto3" json:"next_phase_critical,omitempty"`
}

func (x *BootstrapPhaseEnd) Reset() {
	*x = BootstrapPhaseEnd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapPhaseEnd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapPhaseEnd) ProtoMessage() {}

func (x *BootstrapPhaseEnd) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapPhaseEnd.ProtoReflect.Descriptor instead.
func (*BootstrapPhaseEnd) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{3}
}

func (x *BootstrapPhaseEnd) GetCompletedPhase() string {
	if x != nil {
		return x.CompletedPhase
	}
	return ""
}

func (x *BootstrapPhaseEnd) GetCompletedPhaseState() map[string][]byte {
	if x != nil {
		return x.CompletedPhaseState
	}
	return nil
}

func (x *BootstrapPhaseEnd) GetNextPhase() string {
	if x != nil {
		return x.NextPhase
	}
	return ""
}

func (x *BootstrapPhaseEnd) GetNextPhaseCritical() bool {
	if x != nil {
		return x.NextPhaseCritical
	}
	return false
}

type BootstrapContinue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Continue Continue `protobuf:"varint,1,opt,name=continue,proto3,enum=dhctl.Continue" json:"continue,omitempty"`
	Err      string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *BootstrapContinue) Reset() {
	*x = BootstrapContinue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapContinue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapContinue) ProtoMessage() {}

func (x *BootstrapContinue) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapContinue.ProtoReflect.Descriptor instead.
func (*BootstrapContinue) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{4}
}

func (x *BootstrapContinue) GetContinue() Continue {
	if x != nil {
		return x.Continue
	}
	return Continue_CONTINUE_UNSPECIFIED
}

func (x *BootstrapContinue) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type BootstrapStartOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommanderMode    bool                 `protobuf:"varint,1,opt,name=commander_mode,json=commanderMode,proto3" json:"commander_mode,omitempty"`
	CommanderUuid    string               `protobuf:"bytes,2,opt,name=commander_uuid,json=commanderUuid,proto3" json:"commander_uuid,omitempty"`
	LogWidth         int32                `protobuf:"varint,3,opt,name=log_width,json=logWidth,proto3" json:"log_width,omitempty"`
	ResourcesTimeout *durationpb.Duration `protobuf:"bytes,4,opt,name=resources_timeout,json=resourcesTimeout,proto3" json:"resources_timeout,omitempty"`
	DeckhouseTimeout *durationpb.Duration `protobuf:"bytes,5,opt,name=deckhouse_timeout,json=deckhouseTimeout,proto3" json:"deckhouse_timeout,omitempty"`
	CommonOptions    *OperationOptions    `protobuf:"bytes,10,opt,name=common_options,json=commonOptions,proto3" json:"common_options,omitempty"`
}

func (x *BootstrapStartOptions) Reset() {
	*x = BootstrapStartOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapStartOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapStartOptions) ProtoMessage() {}

func (x *BootstrapStartOptions) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapStartOptions.ProtoReflect.Descriptor instead.
func (*BootstrapStartOptions) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{5}
}

func (x *BootstrapStartOptions) GetCommanderMode() bool {
	if x != nil {
		return x.CommanderMode
	}
	return false
}

func (x *BootstrapStartOptions) GetCommanderUuid() string {
	if x != nil {
		return x.CommanderUuid
	}
	return ""
}

func (x *BootstrapStartOptions) GetLogWidth() int32 {
	if x != nil {
		return x.LogWidth
	}
	return 0
}

func (x *BootstrapStartOptions) GetResourcesTimeout() *durationpb.Duration {
	if x != nil {
		return x.ResourcesTimeout
	}
	return nil
}

func (x *BootstrapStartOptions) GetDeckhouseTimeout() *durationpb.Duration {
	if x != nil {
		return x.DeckhouseTimeout
	}
	return nil
}

func (x *BootstrapStartOptions) GetCommonOptions() *OperationOptions {
	if x != nil {
		return x.CommonOptions
	}
	return nil
}

type BootstrapResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	Err   string `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *BootstrapResult) Reset() {
	*x = BootstrapResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bootstrap_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BootstrapResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BootstrapResult) ProtoMessage() {}

func (x *BootstrapResult) ProtoReflect() protoreflect.Message {
	mi := &file_bootstrap_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BootstrapResult.ProtoReflect.Descriptor instead.
func (*BootstrapResult) Descriptor() ([]byte, []int) {
	return file_bootstrap_proto_rawDescGZIP(), []int{6}
}

func (x *BootstrapResult) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *BootstrapResult) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_bootstrap_proto protoreflect.FileDescriptor

var file_bootstrap_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x10, 0x42, 0x6f, 0x6f, 0x74, 0x73,
	0x74, 0x72, 0x61, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x64, 0x68, 0x63,
	0x74, 0x6c, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x63, 0x6f,
	0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64,
	0x68, 0x63, 0x74, 0x6c, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x43, 0x6f,
	0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x48, 0x00, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e,
	0x75, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xac, 0x01,
	0x0a, 0x11, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x42, 0x6f, 0x6f, 0x74,
	0x73, 0x74, 0x72, 0x61, 0x70, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x48, 0x00, 0x52, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x65,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c,
	0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45,
	0x6e, 0x64, 0x48, 0x00, 0x52, 0x08, 0x70, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x12, 0x21,
	0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x64,
	0x68, 0x63, 0x74, 0x6c, 0x2e, 0x4c, 0x6f, 0x67, 0x73, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x6f, 0x67,
	0x73, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x95, 0x03, 0x0a,
	0x0e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12,
	0x2b, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x0a, 0x0b,
	0x69, 0x6e, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x69, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x47, 0x0a, 0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1d,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x25, 0x0a,
	0x0e, 0x69, 0x6e, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x70, 0x6f, 0x73, 0x74,
	0x5f, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x70, 0x6f, 0x73, 0x74, 0x42, 0x6f, 0x6f,
	0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x36, 0x0a, 0x07,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x22, 0xba, 0x02, 0x0a, 0x11, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72,
	0x61, 0x70, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68,
	0x61, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x31, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x73,
	0x74, 0x72, 0x61, 0x70, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x2e, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x13, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x50, 0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6e, 0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x6e, 0x65, 0x78,
	0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73,
	0x65, 0x43, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x1a, 0x46, 0x0a, 0x18, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x52, 0x0a, 0x11, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x43, 0x6f,
	0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c,
	0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69,
	0x6e, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0xd2, 0x02, 0x0a, 0x15, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74,
	0x72, 0x61, 0x70, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x6d, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x6c, 0x6f, 0x67, 0x57, 0x69, 0x64, 0x74, 0x68, 0x12, 0x46, 0x0a, 0x11, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x46, 0x0a, 0x11, 0x64, 0x65, 0x63, 0x6b, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x64, 0x65, 0x63, 0x6b, 0x68, 0x6f,
	0x75, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x3e, 0x0a, 0x0e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x0d, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x39, 0x0a, 0x0f, 0x42, 0x6f,
	0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x65, 0x72, 0x72, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x62, 0x2f, 0x64, 0x68, 0x63, 0x74,
	0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bootstrap_proto_rawDescOnce sync.Once
	file_bootstrap_proto_rawDescData = file_bootstrap_proto_rawDesc
)

func file_bootstrap_proto_rawDescGZIP() []byte {
	file_bootstrap_proto_rawDescOnce.Do(func() {
		file_bootstrap_proto_rawDescData = protoimpl.X.CompressGZIP(file_bootstrap_proto_rawDescData)
	})
	return file_bootstrap_proto_rawDescData
}

var file_bootstrap_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_bootstrap_proto_goTypes = []interface{}{
	(*BootstrapRequest)(nil),      // 0: dhctl.BootstrapRequest
	(*BootstrapResponse)(nil),     // 1: dhctl.BootstrapResponse
	(*BootstrapStart)(nil),        // 2: dhctl.BootstrapStart
	(*BootstrapPhaseEnd)(nil),     // 3: dhctl.BootstrapPhaseEnd
	(*BootstrapContinue)(nil),     // 4: dhctl.BootstrapContinue
	(*BootstrapStartOptions)(nil), // 5: dhctl.BootstrapStartOptions
	(*BootstrapResult)(nil),       // 6: dhctl.BootstrapResult
	nil,                           // 7: dhctl.BootstrapPhaseEnd.CompletedPhaseStateEntry
	(*Logs)(nil),                  // 8: dhctl.Logs
	(Continue)(0),                 // 9: dhctl.Continue
	(*durationpb.Duration)(nil),   // 10: google.protobuf.Duration
	(*OperationOptions)(nil),      // 11: dhctl.OperationOptions
}
var file_bootstrap_proto_depIdxs = []int32{
	2,  // 0: dhctl.BootstrapRequest.start:type_name -> dhctl.BootstrapStart
	4,  // 1: dhctl.BootstrapRequest.continue:type_name -> dhctl.BootstrapContinue
	6,  // 2: dhctl.BootstrapResponse.result:type_name -> dhctl.BootstrapResult
	3,  // 3: dhctl.BootstrapResponse.phase_end:type_name -> dhctl.BootstrapPhaseEnd
	8,  // 4: dhctl.BootstrapResponse.logs:type_name -> dhctl.Logs
	5,  // 5: dhctl.BootstrapStart.options:type_name -> dhctl.BootstrapStartOptions
	7,  // 6: dhctl.BootstrapPhaseEnd.completed_phase_state:type_name -> dhctl.BootstrapPhaseEnd.CompletedPhaseStateEntry
	9,  // 7: dhctl.BootstrapContinue.continue:type_name -> dhctl.Continue
	10, // 8: dhctl.BootstrapStartOptions.resources_timeout:type_name -> google.protobuf.Duration
	10, // 9: dhctl.BootstrapStartOptions.deckhouse_timeout:type_name -> google.protobuf.Duration
	11, // 10: dhctl.BootstrapStartOptions.common_options:type_name -> dhctl.OperationOptions
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_bootstrap_proto_init() }
func file_bootstrap_proto_init() {
	if File_bootstrap_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_bootstrap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapRequest); i {
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
		file_bootstrap_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapResponse); i {
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
		file_bootstrap_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapStart); i {
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
		file_bootstrap_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapPhaseEnd); i {
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
		file_bootstrap_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapContinue); i {
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
		file_bootstrap_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapStartOptions); i {
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
		file_bootstrap_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BootstrapResult); i {
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
	file_bootstrap_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*BootstrapRequest_Start)(nil),
		(*BootstrapRequest_Continue)(nil),
	}
	file_bootstrap_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*BootstrapResponse_Result)(nil),
		(*BootstrapResponse_PhaseEnd)(nil),
		(*BootstrapResponse_Logs)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_bootstrap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bootstrap_proto_goTypes,
		DependencyIndexes: file_bootstrap_proto_depIdxs,
		MessageInfos:      file_bootstrap_proto_msgTypes,
	}.Build()
	File_bootstrap_proto = out.File
	file_bootstrap_proto_rawDesc = nil
	file_bootstrap_proto_goTypes = nil
	file_bootstrap_proto_depIdxs = nil
}
