// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: proto/metrics.proto

package proto

import (
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

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    string  `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Name    string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Gauge   float64 `protobuf:"fixed64,3,opt,name=gauge,proto3" json:"gauge,omitempty"`
	Counter int64   `protobuf:"varint,4,opt,name=counter,proto3" json:"counter,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Metric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metric) GetGauge() float64 {
	if x != nil {
		return x.Gauge
	}
	return 0
}

func (x *Metric) GetCounter() int64 {
	if x != nil {
		return x.Counter
	}
	return 0
}

type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric  *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Success bool    `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UpdateMetricResponse) Reset() {
	*x = UpdateMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricResponse) ProtoMessage() {}

func (x *UpdateMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateMetricResponse) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *UpdateMetricResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type UpdateMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *UpdateMetricsRequest) Reset() {
	*x = UpdateMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsRequest) ProtoMessage() {}

func (x *UpdateMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricsRequest) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateMetricsRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdateMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UpdateMetricsResponse) Reset() {
	*x = UpdateMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsResponse) ProtoMessage() {}

func (x *UpdateMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricsResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateMetricsResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_proto_metrics_proto protoreflect.FileDescriptor

var file_proto_metrics_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x22, 0x60, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x22, 0x41, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x06, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x5c, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x22, 0x44, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x07,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x31, 0x0a, 0x15, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0xb9, 0x01,
	0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x51, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x12, 0x1f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x20, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x20, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x32, 0x70, 0x32, 0x66, 0x2f, 0x70, 0x72,
	0x61, 0x63, 0x74, 0x69, 0x63, 0x75, 0x6d, 0x2d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrics_proto_rawDescOnce sync.Once
	file_proto_metrics_proto_rawDescData = file_proto_metrics_proto_rawDesc
)

func file_proto_metrics_proto_rawDescGZIP() []byte {
	file_proto_metrics_proto_rawDescOnce.Do(func() {
		file_proto_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrics_proto_rawDescData)
	})
	return file_proto_metrics_proto_rawDescData
}

var file_proto_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_metrics_proto_goTypes = []interface{}{
	(*Metric)(nil),                // 0: grpcmetric.Metric
	(*UpdateMetricRequest)(nil),   // 1: grpcmetric.UpdateMetricRequest
	(*UpdateMetricResponse)(nil),  // 2: grpcmetric.UpdateMetricResponse
	(*UpdateMetricsRequest)(nil),  // 3: grpcmetric.UpdateMetricsRequest
	(*UpdateMetricsResponse)(nil), // 4: grpcmetric.UpdateMetricsResponse
}
var file_proto_metrics_proto_depIdxs = []int32{
	0, // 0: grpcmetric.UpdateMetricRequest.metric:type_name -> grpcmetric.Metric
	0, // 1: grpcmetric.UpdateMetricResponse.metric:type_name -> grpcmetric.Metric
	0, // 2: grpcmetric.UpdateMetricsRequest.metrics:type_name -> grpcmetric.Metric
	1, // 3: grpcmetric.MetricsService.UpdateMetric:input_type -> grpcmetric.UpdateMetricRequest
	3, // 4: grpcmetric.MetricsService.UpdateMetrics:input_type -> grpcmetric.UpdateMetricsRequest
	2, // 5: grpcmetric.MetricsService.UpdateMetric:output_type -> grpcmetric.UpdateMetricResponse
	4, // 6: grpcmetric.MetricsService.UpdateMetrics:output_type -> grpcmetric.UpdateMetricsResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_metrics_proto_init() }
func file_proto_metrics_proto_init() {
	if File_proto_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metrics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_proto_metrics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricRequest); i {
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
		file_proto_metrics_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricResponse); i {
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
		file_proto_metrics_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricsRequest); i {
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
		file_proto_metrics_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricsResponse); i {
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
			RawDescriptor: file_proto_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metrics_proto_goTypes,
		DependencyIndexes: file_proto_metrics_proto_depIdxs,
		MessageInfos:      file_proto_metrics_proto_msgTypes,
	}.Build()
	File_proto_metrics_proto = out.File
	file_proto_metrics_proto_rawDesc = nil
	file_proto_metrics_proto_goTypes = nil
	file_proto_metrics_proto_depIdxs = nil
}