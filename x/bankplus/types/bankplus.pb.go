// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/bankplus/v1/bankplus.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// InactiveAddr models the blocked address for the bankplus module
type InactiveAddr struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *InactiveAddr) Reset()         { *m = InactiveAddr{} }
func (m *InactiveAddr) String() string { return proto.CompactTextString(m) }
func (*InactiveAddr) ProtoMessage()    {}
func (*InactiveAddr) Descriptor() ([]byte, []int) {
	return fileDescriptor_79e8c66834b4419a, []int{0}
}
func (m *InactiveAddr) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InactiveAddr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InactiveAddr.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InactiveAddr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InactiveAddr.Merge(m, src)
}
func (m *InactiveAddr) XXX_Size() int {
	return m.Size()
}
func (m *InactiveAddr) XXX_DiscardUnknown() {
	xxx_messageInfo_InactiveAddr.DiscardUnknown(m)
}

var xxx_messageInfo_InactiveAddr proto.InternalMessageInfo

func (m *InactiveAddr) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*InactiveAddr)(nil), "lbm.bankplus.v1.InactiveAddr")
}

func init() { proto.RegisterFile("lbm/bankplus/v1/bankplus.proto", fileDescriptor_79e8c66834b4419a) }

var fileDescriptor_79e8c66834b4419a = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcb, 0x49, 0xca, 0xd5,
	0x4f, 0x4a, 0xcc, 0xcb, 0x2e, 0xc8, 0x29, 0x2d, 0xd6, 0x2f, 0x33, 0x84, 0xb3, 0xf5, 0x0a, 0x8a,
	0xf2, 0x4b, 0xf2, 0x85, 0xf8, 0x73, 0x92, 0x72, 0xf5, 0xe0, 0x62, 0x65, 0x86, 0x52, 0x22, 0xe9,
	0xf9, 0xe9, 0xf9, 0x60, 0x39, 0x7d, 0x10, 0x0b, 0xa2, 0x4c, 0x49, 0x8f, 0x8b, 0xc7, 0x33, 0x2f,
	0x31, 0xb9, 0x24, 0xb3, 0x2c, 0xd5, 0x31, 0x25, 0xa5, 0x48, 0x48, 0x82, 0x8b, 0x3d, 0x31, 0x25,
	0xa5, 0x28, 0xb5, 0xb8, 0x58, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc6, 0xb5, 0x62, 0x79,
	0xb1, 0x40, 0x9e, 0xd1, 0xc9, 0xe9, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c,
	0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2,
	0x34, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x73, 0x32, 0xf3, 0x52,
	0xf5, 0x73, 0x92, 0x72, 0x75, 0x8b, 0x53, 0xb2, 0xf5, 0x2b, 0x10, 0xce, 0x2c, 0xa9, 0x2c, 0x48,
	0x2d, 0x4e, 0x62, 0x03, 0x5b, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x57, 0x32, 0x20, 0x7a,
	0xc3, 0x00, 0x00, 0x00,
}

func (this *InactiveAddr) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*InactiveAddr)
	if !ok {
		that2, ok := that.(InactiveAddr)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Address != that1.Address {
		return false
	}
	return true
}
func (m *InactiveAddr) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InactiveAddr) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InactiveAddr) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintBankplus(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBankplus(dAtA []byte, offset int, v uint64) int {
	offset -= sovBankplus(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *InactiveAddr) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovBankplus(uint64(l))
	}
	return n
}

func sovBankplus(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBankplus(x uint64) (n int) {
	return sovBankplus(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InactiveAddr) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBankplus
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: InactiveAddr: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InactiveAddr: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBankplus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBankplus
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBankplus
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBankplus(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBankplus
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBankplus(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBankplus
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBankplus
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBankplus
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthBankplus
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBankplus
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBankplus
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBankplus        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBankplus          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBankplus = fmt.Errorf("proto: unexpected end of group")
)
