// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: proto/oteleport.proto

package proto

import (
	v12 "go.opentelemetry.io/proto/otlp/logs/v1"
	v11 "go.opentelemetry.io/proto/otlp/metrics/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TracesData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceSpans []*v1.ResourceSpans `protobuf:"bytes,1,rep,name=resource_spans,json=resourceSpans,proto3" json:"resource_spans,omitempty"`
	SignalType    string              `protobuf:"bytes,2,opt,name=signal_type,json=signalType,proto3" json:"signal_type,omitempty"`
}

func (x *TracesData) Reset() {
	*x = TracesData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TracesData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TracesData) ProtoMessage() {}

func (x *TracesData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TracesData.ProtoReflect.Descriptor instead.
func (*TracesData) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{0}
}

func (x *TracesData) GetResourceSpans() []*v1.ResourceSpans {
	if x != nil {
		return x.ResourceSpans
	}
	return nil
}

func (x *TracesData) GetSignalType() string {
	if x != nil {
		return x.SignalType
	}
	return ""
}

type MetricsData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceMetrics []*v11.ResourceMetrics `protobuf:"bytes,1,rep,name=resource_metrics,json=resourceMetrics,proto3" json:"resource_metrics,omitempty"`
	SignalType      string                 `protobuf:"bytes,2,opt,name=signal_type,json=signalType,proto3" json:"signal_type,omitempty"`
}

func (x *MetricsData) Reset() {
	*x = MetricsData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsData) ProtoMessage() {}

func (x *MetricsData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsData.ProtoReflect.Descriptor instead.
func (*MetricsData) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{1}
}

func (x *MetricsData) GetResourceMetrics() []*v11.ResourceMetrics {
	if x != nil {
		return x.ResourceMetrics
	}
	return nil
}

func (x *MetricsData) GetSignalType() string {
	if x != nil {
		return x.SignalType
	}
	return ""
}

type LogsData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceLogs []*v12.ResourceLogs `protobuf:"bytes,1,rep,name=resource_logs,json=resourceLogs,proto3" json:"resource_logs,omitempty"`
	SignalType   string              `protobuf:"bytes,2,opt,name=signal_type,json=signalType,proto3" json:"signal_type,omitempty"`
}

func (x *LogsData) Reset() {
	*x = LogsData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsData) ProtoMessage() {}

func (x *LogsData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsData.ProtoReflect.Descriptor instead.
func (*LogsData) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{2}
}

func (x *LogsData) GetResourceLogs() []*v12.ResourceLogs {
	if x != nil {
		return x.ResourceLogs
	}
	return nil
}

func (x *LogsData) GetSignalType() string {
	if x != nil {
		return x.SignalType
	}
	return ""
}

type FetchTracesDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTimeUnixNano uint64 `protobuf:"fixed64,1,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`
	EndTimeUnixNano   uint64 `protobuf:"fixed64,2,opt,name=end_time_unix_nano,json=endTimeUnixNano,proto3" json:"end_time_unix_nano,omitempty"`
	Cursor            string `protobuf:"bytes,3,opt,name=cursor,proto3" json:"cursor,omitempty"`
	Limit             int64  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *FetchTracesDataRequest) Reset() {
	*x = FetchTracesDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchTracesDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchTracesDataRequest) ProtoMessage() {}

func (x *FetchTracesDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchTracesDataRequest.ProtoReflect.Descriptor instead.
func (*FetchTracesDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{3}
}

func (x *FetchTracesDataRequest) GetStartTimeUnixNano() uint64 {
	if x != nil {
		return x.StartTimeUnixNano
	}
	return 0
}

func (x *FetchTracesDataRequest) GetEndTimeUnixNano() uint64 {
	if x != nil {
		return x.EndTimeUnixNano
	}
	return 0
}

func (x *FetchTracesDataRequest) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

func (x *FetchTracesDataRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type FetchTracesDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceSpans []*v1.ResourceSpans `protobuf:"bytes,1,rep,name=resource_spans,json=resourceSpans,proto3" json:"resource_spans,omitempty"`
	NextCursor    string              `protobuf:"bytes,2,opt,name=next_cursor,json=nextCursor,proto3" json:"next_cursor,omitempty"`
	HasMore       bool                `protobuf:"varint,3,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
}

func (x *FetchTracesDataResponse) Reset() {
	*x = FetchTracesDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchTracesDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchTracesDataResponse) ProtoMessage() {}

func (x *FetchTracesDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchTracesDataResponse.ProtoReflect.Descriptor instead.
func (*FetchTracesDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{4}
}

func (x *FetchTracesDataResponse) GetResourceSpans() []*v1.ResourceSpans {
	if x != nil {
		return x.ResourceSpans
	}
	return nil
}

func (x *FetchTracesDataResponse) GetNextCursor() string {
	if x != nil {
		return x.NextCursor
	}
	return ""
}

func (x *FetchTracesDataResponse) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

type FetchMetricsDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTimeUnixNano uint64 `protobuf:"fixed64,1,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`
	EndTimeUnixNano   uint64 `protobuf:"fixed64,2,opt,name=end_time_unix_nano,json=endTimeUnixNano,proto3" json:"end_time_unix_nano,omitempty"`
	Cursor            string `protobuf:"bytes,3,opt,name=cursor,proto3" json:"cursor,omitempty"`
	Limit             int64  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *FetchMetricsDataRequest) Reset() {
	*x = FetchMetricsDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchMetricsDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchMetricsDataRequest) ProtoMessage() {}

func (x *FetchMetricsDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchMetricsDataRequest.ProtoReflect.Descriptor instead.
func (*FetchMetricsDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{5}
}

func (x *FetchMetricsDataRequest) GetStartTimeUnixNano() uint64 {
	if x != nil {
		return x.StartTimeUnixNano
	}
	return 0
}

func (x *FetchMetricsDataRequest) GetEndTimeUnixNano() uint64 {
	if x != nil {
		return x.EndTimeUnixNano
	}
	return 0
}

func (x *FetchMetricsDataRequest) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

func (x *FetchMetricsDataRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type FetchMetricsDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceMetrics []*v11.ResourceMetrics `protobuf:"bytes,1,rep,name=resource_metrics,json=resourceMetrics,proto3" json:"resource_metrics,omitempty"`
	NextCursor      string                 `protobuf:"bytes,2,opt,name=next_cursor,json=nextCursor,proto3" json:"next_cursor,omitempty"`
	HasMore         bool                   `protobuf:"varint,3,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
}

func (x *FetchMetricsDataResponse) Reset() {
	*x = FetchMetricsDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchMetricsDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchMetricsDataResponse) ProtoMessage() {}

func (x *FetchMetricsDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchMetricsDataResponse.ProtoReflect.Descriptor instead.
func (*FetchMetricsDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{6}
}

func (x *FetchMetricsDataResponse) GetResourceMetrics() []*v11.ResourceMetrics {
	if x != nil {
		return x.ResourceMetrics
	}
	return nil
}

func (x *FetchMetricsDataResponse) GetNextCursor() string {
	if x != nil {
		return x.NextCursor
	}
	return ""
}

func (x *FetchMetricsDataResponse) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

type FetchLogsDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTimeUnixNano uint64 `protobuf:"fixed64,1,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`
	EndTimeUnixNano   uint64 `protobuf:"fixed64,2,opt,name=end_time_unix_nano,json=endTimeUnixNano,proto3" json:"end_time_unix_nano,omitempty"`
	Cursor            string `protobuf:"bytes,3,opt,name=cursor,proto3" json:"cursor,omitempty"`
	Limit             int64  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *FetchLogsDataRequest) Reset() {
	*x = FetchLogsDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchLogsDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchLogsDataRequest) ProtoMessage() {}

func (x *FetchLogsDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchLogsDataRequest.ProtoReflect.Descriptor instead.
func (*FetchLogsDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{7}
}

func (x *FetchLogsDataRequest) GetStartTimeUnixNano() uint64 {
	if x != nil {
		return x.StartTimeUnixNano
	}
	return 0
}

func (x *FetchLogsDataRequest) GetEndTimeUnixNano() uint64 {
	if x != nil {
		return x.EndTimeUnixNano
	}
	return 0
}

func (x *FetchLogsDataRequest) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

func (x *FetchLogsDataRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type FetchLogsDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceLogs []*v12.ResourceLogs `protobuf:"bytes,1,rep,name=resource_logs,json=resourceLogs,proto3" json:"resource_logs,omitempty"`
	NextCursor   string              `protobuf:"bytes,2,opt,name=next_cursor,json=nextCursor,proto3" json:"next_cursor,omitempty"`
	HasMore      bool                `protobuf:"varint,3,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
}

func (x *FetchLogsDataResponse) Reset() {
	*x = FetchLogsDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oteleport_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchLogsDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchLogsDataResponse) ProtoMessage() {}

func (x *FetchLogsDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oteleport_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchLogsDataResponse.ProtoReflect.Descriptor instead.
func (*FetchLogsDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_oteleport_proto_rawDescGZIP(), []int{8}
}

func (x *FetchLogsDataResponse) GetResourceLogs() []*v12.ResourceLogs {
	if x != nil {
		return x.ResourceLogs
	}
	return nil
}

func (x *FetchLogsDataResponse) GetNextCursor() string {
	if x != nil {
		return x.NextCursor
	}
	return ""
}

func (x *FetchLogsDataResponse) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

var File_proto_oteleport_proto protoreflect.FileDescriptor

var file_proto_oteleport_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x6f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x1a, 0x28, 0x6f, 0x70, 0x65,
	0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2c, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x26, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74,
	0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x73, 0x2f, 0x76, 0x31,
	0x2f, 0x6c, 0x6f, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x0a,
	0x54, 0x72, 0x61, 0x63, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x52, 0x0a, 0x0e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74,
	0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x70, 0x61, 0x6e, 0x73, 0x52,
	0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x70, 0x61, 0x6e, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x22,
	0x8a, 0x01, 0x0a, 0x0b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x5a, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x6f, 0x70, 0x65, 0x6e,
	0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x0f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x6c, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x22, 0x7b, 0x0a, 0x08,
	0x4c, 0x6f, 0x67, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x4e, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x29, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x6c, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x22, 0xa4, 0x01, 0x0a, 0x16, 0x46, 0x65,
	0x74, 0x63, 0x68, 0x54, 0x72, 0x61, 0x63, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x06, 0x52, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e, 0x69,
	0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x12, 0x2b, 0x0a, 0x12, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x06, 0x52, 0x0f, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e, 0x69, 0x78, 0x4e, 0x61,
	0x6e, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x22, 0xa9, 0x01, 0x0a, 0x17, 0x46, 0x65, 0x74, 0x63, 0x68, 0x54, 0x72, 0x61, 0x63, 0x65, 0x73,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x52, 0x0a, 0x0e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x70, 0x61, 0x6e,
	0x73, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x70, 0x61, 0x6e, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x43, 0x75, 0x72, 0x73, 0x6f,
	0x72, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x6f, 0x72, 0x65, 0x22, 0xa5, 0x01, 0x0a,
	0x17, 0x46, 0x65, 0x74, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x06, 0x52, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x55, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x12, 0x2b, 0x0a, 0x12, 0x65, 0x6e, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x06, 0x52, 0x0f, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e,
	0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x14,
	0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x22, 0xb2, 0x01, 0x0a, 0x18, 0x46, 0x65, 0x74, 0x63, 0x68, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x5a, 0x0a, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x6f, 0x70,
	0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x0f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1f, 0x0a,
	0x0b, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x19,
	0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x6f, 0x72, 0x65, 0x22, 0xa2, 0x01, 0x0a, 0x14, 0x46, 0x65,
	0x74, 0x63, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x06,
	0x52, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e, 0x69, 0x78, 0x4e,
	0x61, 0x6e, 0x6f, 0x12, 0x2b, 0x0a, 0x12, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x06, 0x52,
	0x0f, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f,
	0x12, 0x16, 0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0xa3,
	0x01, 0x0a, 0x15, 0x46, 0x65, 0x74, 0x63, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x29, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x78, 0x74,
	0x5f, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e,
	0x65, 0x78, 0x74, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x61, 0x73,
	0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73,
	0x4d, 0x6f, 0x72, 0x65, 0x32, 0xd9, 0x02, 0x0a, 0x10, 0x4f, 0x74, 0x65, 0x72, 0x6c, 0x70, 0x6f,
	0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6c, 0x0a, 0x0f, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x54, 0x72, 0x61, 0x63, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2a, 0x2e, 0x6f,
	0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x31, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x54, 0x72, 0x61, 0x63, 0x65, 0x73, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65,
	0x74, 0x63, 0x68, 0x54, 0x72, 0x61, 0x63, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6f, 0x0a, 0x10, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2b, 0x2e, 0x6f, 0x74,
	0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65,
	0x74, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x66, 0x0a, 0x0d, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x4c, 0x6f, 0x67, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x28, 0x2e, 0x6f, 0x74, 0x65, 0x6c,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x65, 0x74, 0x63, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x4c, 0x6f,
	0x67, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d,
	0x61, 0x73, 0x68, 0x69, 0x69, 0x6b, 0x65, 0x2f, 0x6f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_oteleport_proto_rawDescOnce sync.Once
	file_proto_oteleport_proto_rawDescData = file_proto_oteleport_proto_rawDesc
)

func file_proto_oteleport_proto_rawDescGZIP() []byte {
	file_proto_oteleport_proto_rawDescOnce.Do(func() {
		file_proto_oteleport_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_oteleport_proto_rawDescData)
	})
	return file_proto_oteleport_proto_rawDescData
}

var file_proto_oteleport_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_oteleport_proto_goTypes = []any{
	(*TracesData)(nil),               // 0: oteleport.proto.v1.TracesData
	(*MetricsData)(nil),              // 1: oteleport.proto.v1.MetricsData
	(*LogsData)(nil),                 // 2: oteleport.proto.v1.LogsData
	(*FetchTracesDataRequest)(nil),   // 3: oteleport.proto.v1.FetchTracesDataRequest
	(*FetchTracesDataResponse)(nil),  // 4: oteleport.proto.v1.FetchTracesDataResponse
	(*FetchMetricsDataRequest)(nil),  // 5: oteleport.proto.v1.FetchMetricsDataRequest
	(*FetchMetricsDataResponse)(nil), // 6: oteleport.proto.v1.FetchMetricsDataResponse
	(*FetchLogsDataRequest)(nil),     // 7: oteleport.proto.v1.FetchLogsDataRequest
	(*FetchLogsDataResponse)(nil),    // 8: oteleport.proto.v1.FetchLogsDataResponse
	(*v1.ResourceSpans)(nil),         // 9: opentelemetry.proto.trace.v1.ResourceSpans
	(*v11.ResourceMetrics)(nil),      // 10: opentelemetry.proto.metrics.v1.ResourceMetrics
	(*v12.ResourceLogs)(nil),         // 11: opentelemetry.proto.logs.v1.ResourceLogs
}
var file_proto_oteleport_proto_depIdxs = []int32{
	9,  // 0: oteleport.proto.v1.TracesData.resource_spans:type_name -> opentelemetry.proto.trace.v1.ResourceSpans
	10, // 1: oteleport.proto.v1.MetricsData.resource_metrics:type_name -> opentelemetry.proto.metrics.v1.ResourceMetrics
	11, // 2: oteleport.proto.v1.LogsData.resource_logs:type_name -> opentelemetry.proto.logs.v1.ResourceLogs
	9,  // 3: oteleport.proto.v1.FetchTracesDataResponse.resource_spans:type_name -> opentelemetry.proto.trace.v1.ResourceSpans
	10, // 4: oteleport.proto.v1.FetchMetricsDataResponse.resource_metrics:type_name -> opentelemetry.proto.metrics.v1.ResourceMetrics
	11, // 5: oteleport.proto.v1.FetchLogsDataResponse.resource_logs:type_name -> opentelemetry.proto.logs.v1.ResourceLogs
	3,  // 6: oteleport.proto.v1.OterlportService.FetchTracesData:input_type -> oteleport.proto.v1.FetchTracesDataRequest
	5,  // 7: oteleport.proto.v1.OterlportService.FetchMetricsData:input_type -> oteleport.proto.v1.FetchMetricsDataRequest
	7,  // 8: oteleport.proto.v1.OterlportService.FetchLogsData:input_type -> oteleport.proto.v1.FetchLogsDataRequest
	4,  // 9: oteleport.proto.v1.OterlportService.FetchTracesData:output_type -> oteleport.proto.v1.FetchTracesDataResponse
	6,  // 10: oteleport.proto.v1.OterlportService.FetchMetricsData:output_type -> oteleport.proto.v1.FetchMetricsDataResponse
	8,  // 11: oteleport.proto.v1.OterlportService.FetchLogsData:output_type -> oteleport.proto.v1.FetchLogsDataResponse
	9,  // [9:12] is the sub-list for method output_type
	6,  // [6:9] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_proto_oteleport_proto_init() }
func file_proto_oteleport_proto_init() {
	if File_proto_oteleport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_oteleport_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*TracesData); i {
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
		file_proto_oteleport_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*MetricsData); i {
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
		file_proto_oteleport_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*LogsData); i {
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
		file_proto_oteleport_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*FetchTracesDataRequest); i {
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
		file_proto_oteleport_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*FetchTracesDataResponse); i {
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
		file_proto_oteleport_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*FetchMetricsDataRequest); i {
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
		file_proto_oteleport_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*FetchMetricsDataResponse); i {
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
		file_proto_oteleport_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*FetchLogsDataRequest); i {
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
		file_proto_oteleport_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*FetchLogsDataResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_oteleport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_oteleport_proto_goTypes,
		DependencyIndexes: file_proto_oteleport_proto_depIdxs,
		MessageInfos:      file_proto_oteleport_proto_msgTypes,
	}.Build()
	File_proto_oteleport_proto = out.File
	file_proto_oteleport_proto_rawDesc = nil
	file_proto_oteleport_proto_goTypes = nil
	file_proto_oteleport_proto_depIdxs = nil
}
