// Code generated by protoc-gen-gogo.
// source: checker.proto
// DO NOT EDIT!

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	checker.proto
	spanx.proto
	vape.proto

It has these top-level messages:
	CheckResourceResponse
	ResourceResponse
	CheckResourceRequest
	ResultsResource
	TestCheckRequest
	TestCheckResponse
	PutRoleRequest
	PutRoleResponse
	GetCredentialsRequest
	GetCredentialsResponse
	GetUserRequest
	GetUserResponse
	ListUsersRequest
	ListUsersResponse
*/
package service

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/opsee/protobuf/opseeproto/types"
import opsee_types1 "github.com/opsee/protobuf/opseeproto/types"
import _ "github.com/opsee/protobuf/opseeproto"
import opsee "github.com/opsee/basic/schema"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type CheckResourceResponse struct {
	Id    string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Check *opsee.Check `protobuf:"bytes,2,opt,name=check" json:"check,omitempty"`
	Error string       `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *CheckResourceResponse) Reset()         { *m = CheckResourceResponse{} }
func (m *CheckResourceResponse) String() string { return proto.CompactTextString(m) }
func (*CheckResourceResponse) ProtoMessage()    {}

func (m *CheckResourceResponse) GetCheck() *opsee.Check {
	if m != nil {
		return m.Check
	}
	return nil
}

type ResourceResponse struct {
	Responses []*CheckResourceResponse `protobuf:"bytes,1,rep,name=responses" json:"responses,omitempty"`
}

func (m *ResourceResponse) Reset()         { *m = ResourceResponse{} }
func (m *ResourceResponse) String() string { return proto.CompactTextString(m) }
func (*ResourceResponse) ProtoMessage()    {}

func (m *ResourceResponse) GetResponses() []*CheckResourceResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

type CheckResourceRequest struct {
	Checks []*opsee.Check `protobuf:"bytes,1,rep,name=checks" json:"checks,omitempty"`
}

func (m *CheckResourceRequest) Reset()         { *m = CheckResourceRequest{} }
func (m *CheckResourceRequest) String() string { return proto.CompactTextString(m) }
func (*CheckResourceRequest) ProtoMessage()    {}

func (m *CheckResourceRequest) GetChecks() []*opsee.Check {
	if m != nil {
		return m.Checks
	}
	return nil
}

type ResultsResource struct {
	Results []*opsee.CheckResult `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}

func (m *ResultsResource) Reset()         { *m = ResultsResource{} }
func (m *ResultsResource) String() string { return proto.CompactTextString(m) }
func (*ResultsResource) ProtoMessage()    {}

func (m *ResultsResource) GetResults() []*opsee.CheckResult {
	if m != nil {
		return m.Results
	}
	return nil
}

type TestCheckRequest struct {
	MaxHosts int32                   `protobuf:"varint,1,opt,name=max_hosts,proto3" json:"max_hosts,omitempty"`
	Deadline *opsee_types1.Timestamp `protobuf:"bytes,2,opt,name=deadline" json:"deadline,omitempty"`
	Check    *opsee.Check            `protobuf:"bytes,3,opt,name=check" json:"check,omitempty"`
}

func (m *TestCheckRequest) Reset()         { *m = TestCheckRequest{} }
func (m *TestCheckRequest) String() string { return proto.CompactTextString(m) }
func (*TestCheckRequest) ProtoMessage()    {}

func (m *TestCheckRequest) GetDeadline() *opsee_types1.Timestamp {
	if m != nil {
		return m.Deadline
	}
	return nil
}

func (m *TestCheckRequest) GetCheck() *opsee.Check {
	if m != nil {
		return m.Check
	}
	return nil
}

type TestCheckResponse struct {
	Responses []*opsee.CheckResponse `protobuf:"bytes,1,rep,name=responses" json:"responses,omitempty"`
	Error     string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *TestCheckResponse) Reset()         { *m = TestCheckResponse{} }
func (m *TestCheckResponse) String() string { return proto.CompactTextString(m) }
func (*TestCheckResponse) ProtoMessage()    {}

func (m *TestCheckResponse) GetResponses() []*opsee.CheckResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

func init() {
	proto.RegisterType((*CheckResourceResponse)(nil), "opsee.CheckResourceResponse")
	proto.RegisterType((*ResourceResponse)(nil), "opsee.ResourceResponse")
	proto.RegisterType((*CheckResourceRequest)(nil), "opsee.CheckResourceRequest")
	proto.RegisterType((*ResultsResource)(nil), "opsee.ResultsResource")
	proto.RegisterType((*TestCheckRequest)(nil), "opsee.TestCheckRequest")
	proto.RegisterType((*TestCheckResponse)(nil), "opsee.TestCheckResponse")
}
func (this *CheckResourceResponse) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*CheckResourceResponse)
	if !ok {
		that2, ok := that.(CheckResourceResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	if !this.Check.Equal(that1.Check) {
		return false
	}
	if this.Error != that1.Error {
		return false
	}
	return true
}
func (this *ResourceResponse) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ResourceResponse)
	if !ok {
		that2, ok := that.(ResourceResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.Responses) != len(that1.Responses) {
		return false
	}
	for i := range this.Responses {
		if !this.Responses[i].Equal(that1.Responses[i]) {
			return false
		}
	}
	return true
}
func (this *CheckResourceRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*CheckResourceRequest)
	if !ok {
		that2, ok := that.(CheckResourceRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.Checks) != len(that1.Checks) {
		return false
	}
	for i := range this.Checks {
		if !this.Checks[i].Equal(that1.Checks[i]) {
			return false
		}
	}
	return true
}
func (this *ResultsResource) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*ResultsResource)
	if !ok {
		that2, ok := that.(ResultsResource)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.Results) != len(that1.Results) {
		return false
	}
	for i := range this.Results {
		if !this.Results[i].Equal(that1.Results[i]) {
			return false
		}
	}
	return true
}
func (this *TestCheckRequest) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*TestCheckRequest)
	if !ok {
		that2, ok := that.(TestCheckRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.MaxHosts != that1.MaxHosts {
		return false
	}
	if !this.Deadline.Equal(that1.Deadline) {
		return false
	}
	if !this.Check.Equal(that1.Check) {
		return false
	}
	return true
}
func (this *TestCheckResponse) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*TestCheckResponse)
	if !ok {
		that2, ok := that.(TestCheckResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if len(this.Responses) != len(that1.Responses) {
		return false
	}
	for i := range this.Responses {
		if !this.Responses[i].Equal(that1.Responses[i]) {
			return false
		}
	}
	if this.Error != that1.Error {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Checker service

type CheckerClient interface {
	TestCheck(ctx context.Context, in *TestCheckRequest, opts ...grpc.CallOption) (*TestCheckResponse, error)
	CreateCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error)
	RetrieveCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error)
	UpdateCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error)
	DeleteCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error)
}

type checkerClient struct {
	cc *grpc.ClientConn
}

func NewCheckerClient(cc *grpc.ClientConn) CheckerClient {
	return &checkerClient{cc}
}

func (c *checkerClient) TestCheck(ctx context.Context, in *TestCheckRequest, opts ...grpc.CallOption) (*TestCheckResponse, error) {
	out := new(TestCheckResponse)
	err := grpc.Invoke(ctx, "/opsee.Checker/TestCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerClient) CreateCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error) {
	out := new(ResourceResponse)
	err := grpc.Invoke(ctx, "/opsee.Checker/CreateCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerClient) RetrieveCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error) {
	out := new(ResourceResponse)
	err := grpc.Invoke(ctx, "/opsee.Checker/RetrieveCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerClient) UpdateCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error) {
	out := new(ResourceResponse)
	err := grpc.Invoke(ctx, "/opsee.Checker/UpdateCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerClient) DeleteCheck(ctx context.Context, in *CheckResourceRequest, opts ...grpc.CallOption) (*ResourceResponse, error) {
	out := new(ResourceResponse)
	err := grpc.Invoke(ctx, "/opsee.Checker/DeleteCheck", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Checker service

type CheckerServer interface {
	TestCheck(context.Context, *TestCheckRequest) (*TestCheckResponse, error)
	CreateCheck(context.Context, *CheckResourceRequest) (*ResourceResponse, error)
	RetrieveCheck(context.Context, *CheckResourceRequest) (*ResourceResponse, error)
	UpdateCheck(context.Context, *CheckResourceRequest) (*ResourceResponse, error)
	DeleteCheck(context.Context, *CheckResourceRequest) (*ResourceResponse, error)
}

func RegisterCheckerServer(s *grpc.Server, srv CheckerServer) {
	s.RegisterService(&_Checker_serviceDesc, srv)
}

func _Checker_TestCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(TestCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CheckerServer).TestCheck(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Checker_CreateCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CheckResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CheckerServer).CreateCheck(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Checker_RetrieveCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CheckResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CheckerServer).RetrieveCheck(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Checker_UpdateCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CheckResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CheckerServer).UpdateCheck(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Checker_DeleteCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CheckResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CheckerServer).DeleteCheck(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Checker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "opsee.Checker",
	HandlerType: (*CheckerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TestCheck",
			Handler:    _Checker_TestCheck_Handler,
		},
		{
			MethodName: "CreateCheck",
			Handler:    _Checker_CreateCheck_Handler,
		},
		{
			MethodName: "RetrieveCheck",
			Handler:    _Checker_RetrieveCheck_Handler,
		},
		{
			MethodName: "UpdateCheck",
			Handler:    _Checker_UpdateCheck_Handler,
		},
		{
			MethodName: "DeleteCheck",
			Handler:    _Checker_DeleteCheck_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func NewPopulatedCheckResourceResponse(r randyChecker, easy bool) *CheckResourceResponse {
	this := &CheckResourceResponse{}
	this.Id = randStringChecker(r)
	if r.Intn(10) != 0 {
		this.Check = opsee.NewPopulatedCheck(r, easy)
	}
	this.Error = randStringChecker(r)
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedResourceResponse(r randyChecker, easy bool) *ResourceResponse {
	this := &ResourceResponse{}
	if r.Intn(10) != 0 {
		v1 := r.Intn(5)
		this.Responses = make([]*CheckResourceResponse, v1)
		for i := 0; i < v1; i++ {
			this.Responses[i] = NewPopulatedCheckResourceResponse(r, easy)
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedCheckResourceRequest(r randyChecker, easy bool) *CheckResourceRequest {
	this := &CheckResourceRequest{}
	if r.Intn(10) != 0 {
		v2 := r.Intn(5)
		this.Checks = make([]*opsee.Check, v2)
		for i := 0; i < v2; i++ {
			this.Checks[i] = opsee.NewPopulatedCheck(r, easy)
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedResultsResource(r randyChecker, easy bool) *ResultsResource {
	this := &ResultsResource{}
	if r.Intn(10) != 0 {
		v3 := r.Intn(5)
		this.Results = make([]*opsee.CheckResult, v3)
		for i := 0; i < v3; i++ {
			this.Results[i] = opsee.NewPopulatedCheckResult(r, easy)
		}
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedTestCheckRequest(r randyChecker, easy bool) *TestCheckRequest {
	this := &TestCheckRequest{}
	this.MaxHosts = int32(r.Int31())
	if r.Intn(2) == 0 {
		this.MaxHosts *= -1
	}
	if r.Intn(10) != 0 {
		this.Deadline = opsee_types1.NewPopulatedTimestamp(r, easy)
	}
	if r.Intn(10) != 0 {
		this.Check = opsee.NewPopulatedCheck(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedTestCheckResponse(r randyChecker, easy bool) *TestCheckResponse {
	this := &TestCheckResponse{}
	if r.Intn(10) != 0 {
		v4 := r.Intn(5)
		this.Responses = make([]*opsee.CheckResponse, v4)
		for i := 0; i < v4; i++ {
			this.Responses[i] = opsee.NewPopulatedCheckResponse(r, easy)
		}
	}
	this.Error = randStringChecker(r)
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyChecker interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneChecker(r randyChecker) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringChecker(r randyChecker) string {
	v5 := r.Intn(100)
	tmps := make([]rune, v5)
	for i := 0; i < v5; i++ {
		tmps[i] = randUTF8RuneChecker(r)
	}
	return string(tmps)
}
func randUnrecognizedChecker(r randyChecker, maxFieldNumber int) (data []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		data = randFieldChecker(data, r, fieldNumber, wire)
	}
	return data
}
func randFieldChecker(data []byte, r randyChecker, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		data = encodeVarintPopulateChecker(data, uint64(key))
		v6 := r.Int63()
		if r.Intn(2) == 0 {
			v6 *= -1
		}
		data = encodeVarintPopulateChecker(data, uint64(v6))
	case 1:
		data = encodeVarintPopulateChecker(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		data = encodeVarintPopulateChecker(data, uint64(key))
		ll := r.Intn(100)
		data = encodeVarintPopulateChecker(data, uint64(ll))
		for j := 0; j < ll; j++ {
			data = append(data, byte(r.Intn(256)))
		}
	default:
		data = encodeVarintPopulateChecker(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return data
}
func encodeVarintPopulateChecker(data []byte, v uint64) []byte {
	for v >= 1<<7 {
		data = append(data, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	data = append(data, uint8(v))
	return data
}
