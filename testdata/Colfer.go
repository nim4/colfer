package testdata

// This file was generated by colf(1); DO NOT EDIT

import (
	"fmt"
	"io"
	"math"
	"time"
)

// Colfer configuration attributes
var (
	// ColferSizeMax is the upper limit for serial byte sizes.
	ColferSizeMax = 16 * 1024 * 1024
	// ColferListMax is the upper limit for the number of elements in a list.
	ColferListMax = 64 * 1024
)

// ColferMax signals an upper limit breach.
type ColferMax string

// Error honors the error interface.
func (m ColferMax) Error() string { return string(m) }

// ColferError signals a data mismatch as as a byte index.
type ColferError int

// Error honors the error interface.
func (i ColferError) Error() string {
	return fmt.Sprintf("colfer: unknown header at byte %d", i)
}

// ColferTail signals data continuation as a byte index.
type ColferTail int

// Error honors the error interface.
func (i ColferTail) Error() string {
	return fmt.Sprintf("colfer: data continuation at byte %d", i)
}

type O struct {
	B	bool
	U32	uint32
	U64	uint64
	I32	int32
	I64	int64
	F32	float32
	F64	float64
	T	time.Time
	S	string
	A	[]byte
	O	*O
	Os	[]*O
	Ss	[]string
	As	[][]byte
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
// All nil entries in o.Os will be replaced with a new value.
func (o *O) MarshalTo(buf []byte) int {
	var i int

	if o.B {
		buf[i] = 0
		i++
	}

	if x := o.U32; x >= 1<<21 {
		buf[i] = 1 | 0x80
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 5
	} else if x != 0 {
		buf[i] = 1
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if x := o.U64; x >= 1<<49 {
		buf[i] = 2 | 0x80
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 9
	} else if x != 0 {
		buf[i] = 2
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if v := o.I32; v != 0 {
		x := uint32(v)
		if v >= 0 {
			buf[i] = 3
		} else {
			x = ^x + 1
			buf[i] = 3 | 0x80
		}
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if v := o.I64; v != 0 {
		x := uint64(v)
		if v >= 0 {
			buf[i] = 4
		} else {
			x = ^x + 1
			buf[i] = 4 | 0x80
		}
		i++
		for n := 0; x >= 0x80 && n < 8; n++ {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if v := o.F32; v != 0.0 {
		buf[i] = 5
		x := math.Float32bits(v)
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 5
	}

	if v := o.F64; v != 0.0 {
		buf[i] = 6
		x := math.Float64bits(v)
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 9
	}

	if v := o.T; !v.IsZero() {
		s, ns := uint64(v.Unix()), uint(v.Nanosecond())
		if s < 1<<32 {
			buf[i] = 7
			buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
			buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 9
		} else {
			buf[i] = 7 | 0x80
			buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(s>>56), byte(s>>48), byte(s>>40), byte(s>>32)
			buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
			buf[i+9], buf[i+10], buf[i+11], buf[i+12] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 13
		}
	}

	if l := len(o.S); l != 0 {
		buf[i] = 8
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.S)
	}

	if l := len(o.A); l != 0 {
		buf[i] = 9
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		i += copy(buf[i:], o.A)
	}

	if v := o.O; v != nil {
		buf[i] = 10
		i++
		i += v.MarshalTo(buf[i:])
	}

	if l := len(o.Os); l != 0 {
		buf[i] = 11
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		for vi, v := range o.Os {
			if v == nil {
				v = new(O)
				o.Os[vi] = v
			}
			i += v.MarshalTo(buf[i:])
		}
	}

	if l := len(o.Ss); l != 0 {
		buf[i] = 12
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		for _, a := range o.Ss {
			x = uint(len(a))
			for x >= 0x80 {
				buf[i] = byte(x | 0x80)
				x >>= 7
				i++
			}
			buf[i] = byte(x)
			i++
			i += copy(buf[i:], a)
		}
	}

	if l := len(o.As); l != 0 {
		buf[i] = 13
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		for _, a := range o.As {
			x = uint(len(a))
			for x >= 0x80 {
				buf[i] = byte(x | 0x80)
				x >>= 7
				i++
			}
			buf[i] = byte(x)
			i++
			i += copy(buf[i:], a)
		}
	}

	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is testdata.ColferMax.
func (o *O) MarshalLen() (int, error) {
	l := 1

	if o.B {
		l++
	}

	if x := o.U32; x >= 1<<21 {
		l += 5
	} else if x != 0 {
		l += 2
		for x >= 0x80 {
			x >>= 7
			l++
		}
	}

	if x := o.U64; x >= 1<<49 {
		l += 9
	} else if x != 0 {
		l += 2
		for x >= 0x80 {
			x >>= 7
			l++
		}
	}

	if v := o.I32; v != 0 {
		l += 2
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
	}

	if v := o.I64; v != 0 {
		l += 2
		x := uint64(v)
		if v < 0 {
			x = ^x + 1
		}
		for n := 0; x >= 0x80 && n < 8; n++ {
			x >>= 7
			l++
		}
	}

	if o.F32 != 0.0 {
		l += 5
	}

	if o.F64 != 0.0 {
		l += 9
	}

	if v := o.T; !v.IsZero() {
		if s := uint(v.Unix()); s < 1<<32 {
			l += 9
		} else {
			l += 13
		}
	}

	if x := len(o.S); x != 0 {
		l += x
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
	}

	if x := len(o.A); x != 0 {
		l += x
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
	}

	if v := o.O; v != nil {
		vl, err := v.MarshalLen()
		if err != nil {
			return -1, err
		}
		l += vl + 1
	}

	if x := len(o.Os); x != 0 {
		if x > ColferListMax {
			return -1, ColferMax(fmt.Sprintf("colfer: field testdata.o.os exceeds %d elements", ColferListMax))
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
		for _, v := range o.Os {
			if v == nil {
				l++
				continue
			}
			vl, err := v.MarshalLen()
			if err != nil {
				return -1, err
			}
			l += vl
		}
	}

	if x := len(o.Ss); x != 0 {
		if x > ColferListMax {
			return -1, ColferMax(fmt.Sprintf("colfer: field testdata.o.ss exceeds %d elements", ColferListMax))
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
		for _, a := range o.Ss {
			x = len(a)
			l += x
			for x >= 0x80 {
				x >>= 7
				l++
			}
			l++
		}
	}

	if x := len(o.As); x != 0 {
		if x > ColferListMax {
			return -1, ColferMax(fmt.Sprintf("colfer: field testdata.o.as exceeds %d elements", ColferListMax))
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
		for _, a := range o.As {
			x = len(a)
			l += x
			for x >= 0x80 {
				x >>= 7
				l++
			}
			l++
		}
	}

	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct testdata.o exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// All nil entries in o.Os will be replaced with a new value.
// The error return option is testdata.ColferMax.
func (o *O) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, testdata.ColferError and testdata.ColferMax.
func (o *O) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1

	if header == 0 {
		o.B = true
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}

	if header == 1 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.U32 = x

		header = data[i]
		i++
	} else if header == 1|0x80 {
		if i+4 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		o.U32 = uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		header = data[i+4]
		i += 5
	}

	if header == 2 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.U64 = x

		header = data[i]
		i++
	} else if header == 2|0x80 {
		if i+8 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		o.U64 = uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32 | uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		header = data[i+8]
		i += 9
	}

	if header == 3 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.I32 = int32(x)

		header = data[i]
		i++
	} else if header == 3|0x80 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint32(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint32(data[i])
				i++
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.I32 = int32(^x + 1)

		header = data[i]
		i++
	}

	if header == 4 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.I64 = int64(x)

		header = data[i]
		i++
	} else if header == 4|0x80 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				b := uint64(data[i])
				i++
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}

				if b < 0x80 || shift == 56 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		o.I64 = int64(^x + 1)

		header = data[i]
		i++
	}

	if header == 5 {
		if i+4 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		o.F32 = math.Float32frombits(x)

		header = data[i+4]
		i += 5
	}

	if header == 6 {
		if i+8 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		o.F64 = math.Float64frombits(x)

		header = data[i+8]
		i += 9
	}

	if header == 7 {
		if i+8 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		v := int64(x)
		o.T = time.Unix(v>>32, v&(1<<32-1)).In(time.UTC)

		header = data[i+8]
		i += 9
	} else if header == 7|0x80 {
		if i+12 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		nsec := uint64(data[i+8])<<24 | uint64(data[i+9])<<16 | uint64(data[i+10])<<8 | uint64(data[i+11])
		o.T = time.Unix(int64(sec), int64(nsec)).In(time.UTC)

		header = data[i+12]
		i += 13
	}

	if header == 8 {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.s size %d exceeds %d bytes", x, ColferSizeMax))
		}
		to := i + int(x)
		if to >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		o.S = string(data[i:to])

		header = data[to]
		i = to + 1
	}

	if header == 9 {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.a size %d exceeds %d bytes", x, ColferSizeMax))
		}
		l := int(x)
		if i + l >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		v := make([]byte, l)
		i += copy(v, data[i:])
		o.A = v

		header = data[i]
		i++
	}

	if header == 10 {
		o.O = new(O)
		n, err := o.O.Unmarshal(data[i:])
		if err != nil {
			return 0, err
		}
		i += n

		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}

	if header == 11 {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					return 0, io.EOF
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.os length %d exceeds %d elements", x, ColferListMax))
		}

		l := int(x)
		a := make([]*O, l)
		malloc := make([]O, l)
		for ai, _ := range a {
			v := &malloc[ai]
			a[ai] = v

			n, err := v.Unmarshal(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
		o.Os = a

		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}

	if header == 12 {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.ss length %d exceeds %d elements", x, ColferListMax))
		}
		a := make([]string, int(x))
		o.Ss = a
		for ai := range a {
			if i >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
				}
				return 0, io.EOF
			}
			x := uint(data[i])
			i++

			if x >= 0x80 {
				x &= 0x7f
				for shift := uint(7); ; shift += 7 {
					if i >= len(data) {
						if i >= ColferSizeMax {
							return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
						}
						return 0, io.EOF
					}
					b := uint(data[i])
					i++

					if b < 0x80 {
						x |= b << shift
						break
					}
					x |= (b & 0x7f) << shift
				}
			}
			if x > uint(ColferSizeMax) {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.ss element %d size %d exceeds %d bytes", ai, x, ColferSizeMax))
			}
			to := i + int(x)
			if to >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
				}
				return 0, io.EOF
			}
			a[ai] = string(data[i:to])
			i = to
		}

		if i >= len(data) {
			return 0, io.EOF
		}
		header = data[i]
		i++
	}

	if header == 13 {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint(data[i])
		i++

		if x >= 0x80 {
			x &= 0x7f
			for shift := uint(7); ; shift += 7 {
				if i >= len(data) {
					if i >= ColferSizeMax {
						return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
					}
					return 0, io.EOF
				}
				b := uint(data[i])
				i++

				if b < 0x80 {
					x |= b << shift
					break
				}
				x |= (b & 0x7f) << shift
			}
		}
		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.as length %d exceeds %d elements", x, ColferListMax))
		}
		a := make([][]byte, int(x))
		o.As = a
		for ai := range a {
			if i >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
				}
				return 0, io.EOF
			}
			x := uint(data[i])
			i++

			if x >= 0x80 {
				x &= 0x7f
				for shift := uint(7); ; shift += 7 {
					if i >= len(data) {
						if i >= ColferSizeMax {
							return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
						}
						return 0, io.EOF
					}
					b := uint(data[i])
					i++

					if b < 0x80 {
						x |= b << shift
						break
					}
					x |= (b & 0x7f) << shift
				}
			}
			if x > uint(ColferSizeMax) {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o.as element %d size %d exceeds %d bytes", ai, x, ColferSizeMax))
			}
			v := make([]byte, int(x))

			if i + len(v) >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
				}
				return 0, io.EOF
			}

			i += copy(v, data[i:])
			a[ai] = v
		}

		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}

	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: testdata.o size %d exceeds %d bytes", i, ColferSizeMax))
	}
	return i, nil
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, testdata.ColferError, testdata.ColferTail and testdata.ColferMax.
func (o *O) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}
