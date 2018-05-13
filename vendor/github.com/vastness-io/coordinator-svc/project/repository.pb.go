// Code generated by protoc-gen-go. DO NOT EDIT.
// source: repository.proto

package project

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Repository struct {
	Name     string    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Owner    string    `protobuf:"bytes,2,opt,name=owner" json:"owner,omitempty"`
	Type     string    `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	Branches []*Branch `protobuf:"bytes,4,rep,name=branches" json:"branches,omitempty"`
}

func (m *Repository) Reset()                    { *m = Repository{} }
func (m *Repository) String() string            { return proto.CompactTextString(m) }
func (*Repository) ProtoMessage()               {}
func (*Repository) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *Repository) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Repository) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Repository) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Repository) GetBranches() []*Branch {
	if m != nil {
		return m.Branches
	}
	return nil
}

func init() {
	proto.RegisterType((*Repository)(nil), "project.Repository")
}

func init() { proto.RegisterFile("repository.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 142 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x4a, 0x2d, 0xc8,
	0x2f, 0xce, 0x2c, 0xc9, 0x2f, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2f, 0x28,
	0xca, 0xcf, 0x4a, 0x4d, 0x2e, 0x91, 0xe2, 0x49, 0x2a, 0x4a, 0xcc, 0x4b, 0xce, 0x80, 0x08, 0x2b,
	0x95, 0x73, 0x71, 0x05, 0xc1, 0x95, 0x0a, 0x09, 0x71, 0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30,
	0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42, 0x22, 0x5c, 0xac, 0xf9, 0xe5, 0x79, 0xa9, 0x45,
	0x12, 0x4c, 0x60, 0x41, 0x08, 0x07, 0xa4, 0xb2, 0xa4, 0xb2, 0x20, 0x55, 0x82, 0x19, 0xa2, 0x12,
	0xc4, 0x16, 0xd2, 0xe6, 0xe2, 0x80, 0x98, 0x9d, 0x5a, 0x2c, 0xc1, 0xa2, 0xc0, 0xac, 0xc1, 0x6d,
	0xc4, 0xaf, 0x07, 0xb5, 0x55, 0xcf, 0x09, 0x2c, 0x11, 0x04, 0x57, 0x90, 0xc4, 0x06, 0xb6, 0xdf,
	0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xf2, 0x9b, 0xc1, 0x5f, 0xaa, 0x00, 0x00, 0x00,
}