package colfer

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generate writes the code into file "Colfer.go".
func Generate(basedir string, packages []*Package) error {
	t := template.New("go-code").Delims("<:", ":>")
	template.Must(t.Parse(goCode))
	template.Must(t.New("marshal-field").Parse(goMarshalField))
	template.Must(t.New("marshal-field-len").Parse(goMarshalFieldLen))
	template.Must(t.New("unmarshal-field").Parse(goUnmarshalField))

	for _, p := range packages {
		p.NameNative = p.Name[strings.LastIndexByte(p.Name, '/')+1:]
	}

	for _, p := range packages {
		for _, s := range p.Structs {
			for _, f := range s.Fields {
				switch f.Type {
				default:
					if f.TypeRef == nil {
						f.TypeNative = f.Type
					} else {
						f.TypeNative = f.TypeRef.NameTitle()
						if f.TypeRef.Pkg != p {
							f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
						}
					}
				case "timestamp":
					f.TypeNative = "time.Time"
				case "text":
					f.TypeNative = "string"
				case "binary":
					f.TypeNative = "[]byte"
				}
			}
		}

		pkgdir, err := makePkgDir(p, basedir)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(pkgdir, "Colfer.go"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err = t.Execute(f, p); err != nil {
			return err
		}
	}
	return nil
}

func makePkgDir(p *Package, basedir string) (path string, err error) {
	pkgdir := strings.Replace(p.Name, "/", string(filepath.Separator), -1)
	path = filepath.Join(basedir, pkgdir)
	err = os.MkdirAll(path, os.ModeDir|os.ModePerm)
	return
}

const goCode = `package <:.NameNative:>

// This file was generated by colf(1); DO NOT EDIT

import (
	"fmt"
	"io"
<:- if .HasFloat:>
	"math"
<:- end:>
<:- if .HasTimestamp:>
	"time"
<:- end:>
)

// Colfer configuration attributes
var (
	// ColferSizeMax is the upper limit for serial byte sizes.
	ColferSizeMax = 16 * 1024 * 1024

<:- if .HasList:>
	// ColferListMax is the upper limit for the number of elements in a list.
	ColferListMax = 64 * 1024
<:- end:>
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
<:range .Structs:>
type <:.NameTitle:> struct {
<:range .Fields:>	<:.NameTitle:>	<:if .TypeList:>[]<:end:><:if .TypeRef:>*<:end:><:.TypeNative:>
<:end:>}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MarshalTo will panic.
<:- range .Fields:><:if and .TypeList .TypeRef:>
// All nil entries in o.<:.NameTitle:> will be replaced with a new value.
<:- end:><:end:>
func (o *<:.NameTitle:>) MarshalTo(buf []byte) int {
	var i int
<:range .Fields:><:template "marshal-field" .:><:end:>
	buf[i] = 0x7f
	i++
	return i
}

// MarshalLen returns the Colfer serial byte size.
// The error return option is <:.Pkg.NameNative:>.ColferMax.
func (o *<:.NameTitle:>) MarshalLen() (int, error) {
	l := 1
<:range .Fields:><:template "marshal-field-len" .:><:end:>
	if l > ColferSizeMax {
		return l, ColferMax(fmt.Sprintf("colfer: struct <:.String:> exceeds %d bytes", ColferSizeMax))
	}
	return l, nil
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
<:- range .Fields:><:if and .TypeList .TypeRef:>
// All nil entries in o.<:.NameTitle:> will be replaced with a new value.
<:- end:><:end:>
// The error return option is <:.Pkg.NameNative:>.ColferMax.
func (o *<:.NameTitle:>) MarshalBinary() (data []byte, err error) {
	l, err := o.MarshalLen()
	if err != nil {
		return nil, err
	}
	data = make([]byte, l)
	o.MarshalTo(data)
	return data, nil
}

// Unmarshal decodes data as Colfer and returns the number of bytes read.
// The error return options are io.EOF, <:.Pkg.NameNative:>.ColferError and <:.Pkg.NameNative:>.ColferMax.
func (o *<:.NameTitle:>) Unmarshal(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, io.EOF
	}
	header := data[0]
	i := 1
<:range .Fields:><:template "unmarshal-field" .:><:end:>
	if header != 0x7f {
		return 0, ColferError(i - 1)
	}
	if i >= ColferSizeMax {
		return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> size %d exceeds %d bytes", i, ColferSizeMax))
	}
	return i, nil
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, <:.Pkg.NameNative:>.ColferError, <:.Pkg.NameNative:>.ColferTail and <:.Pkg.NameNative:>.ColferMax.
func (o *<:.NameTitle:>) UnmarshalBinary(data []byte) error {
	i, err := o.Unmarshal(data)
	if i < len(data) && err == nil {
		return ColferTail(i)
	}
	return err
}
<:end:>`

const goMarshalField = `<:if eq .Type "bool":>
	if o.<:.NameTitle:> {
		buf[i] = <:.Index:>
		i++
	}
<:else if eq .Type "uint32":>
	if x := o.<:.NameTitle:>; x >= 1<<21 {
		buf[i] = <:.Index:> | 0x80
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 5
	} else if x != 0 {
		buf[i] = <:.Index:>
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}
<:else if eq .Type "uint64":>
	if x := o.<:.NameTitle:>; x >= 1<<49 {
		buf[i] = <:.Index:> | 0x80
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 9
	} else if x != 0 {
		buf[i] = <:.Index:>
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}
<:else if eq .Type "int32":>
	if v := o.<:.NameTitle:>; v != 0 {
		x := uint32(v)
		if v >= 0 {
			buf[i] = <:.Index:>
		} else {
			x = ^x + 1
			buf[i] = <:.Index:> | 0x80
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
<:else if eq .Type "int64":>
	if v := o.<:.NameTitle:>; v != 0 {
		x := uint64(v)
		if v >= 0 {
			buf[i] = <:.Index:>
		} else {
			x = ^x + 1
			buf[i] = <:.Index:> | 0x80
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
<:else if eq .Type "float32":>
	if v := o.<:.NameTitle:>; v != 0.0 {
		buf[i] = <:.Index:>
		x := math.Float32bits(v)
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 5
	}
<:else if eq .Type "float64":>
	if v := o.<:.NameTitle:>; v != 0.0 {
		buf[i] = <:.Index:>
		x := math.Float64bits(v)
		buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 9
	}
<:else if eq .Type "timestamp":>
	if v := o.<:.NameTitle:>; !v.IsZero() {
		s, ns := uint64(v.Unix()), uint(v.Nanosecond())
		if s < 1<<32 {
			buf[i] = <:.Index:>
			buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
			buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 9
		} else {
			buf[i] = <:.Index:> | 0x80
			buf[i+1], buf[i+2], buf[i+3], buf[i+4] = byte(s>>56), byte(s>>48), byte(s>>40), byte(s>>32)
			buf[i+5], buf[i+6], buf[i+7], buf[i+8] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
			buf[i+9], buf[i+10], buf[i+11], buf[i+12] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 13
		}
	}
<:else if eq .Type "text" "binary":>
	if l := len(o.<:.NameTitle:>); l != 0 {
		buf[i] = <:.Index:>
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
 <:- if .TypeList:>
		for _, a := range o.<:.NameTitle:> {
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
 <:- else:>
		i += copy(buf[i:], o.<:.NameTitle:>)
 <:- end:>
	}
<:else if .TypeList:>
	if l := len(o.<:.NameTitle:>); l != 0 {
		buf[i] = <:.Index:>
		i++
		x := uint(l)
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		for vi, v := range o.<:.NameTitle:> {
			if v == nil {
				v = new(<:.TypeNative:>)
				o.<:.NameTitle:>[vi] = v
			}
			i += v.MarshalTo(buf[i:])
		}
	}
<:else:>
	if v := o.<:.NameTitle:>; v != nil {
		buf[i] = <:.Index:>
		i++
		i += v.MarshalTo(buf[i:])
	}
<:end:>`

const goMarshalFieldLen = `<:if eq .Type "bool":>
	if o.<:.NameTitle:> {
		l++
	}
<:else if eq .Type "uint32":>
	if x := o.<:.NameTitle:>; x >= 1<<21 {
		l += 5
	} else if x != 0 {
		l += 2
		for x >= 0x80 {
			x >>= 7
			l++
		}
	}
<:else if eq .Type "uint64":>
	if x := o.<:.NameTitle:>; x >= 1<<49 {
		l += 9
	} else if x != 0 {
		l += 2
		for x >= 0x80 {
			x >>= 7
			l++
		}
	}
<:else if eq .Type "int32":>
	if v := o.<:.NameTitle:>; v != 0 {
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
<:else if eq .Type "int64":>
	if v := o.<:.NameTitle:>; v != 0 {
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
<:else if eq .Type "float32":>
	if o.<:.NameTitle:> != 0.0 {
		l += 5
	}
<:else if eq .Type "float64":>
	if o.<:.NameTitle:> != 0.0 {
		l += 9
	}
<:else if eq .Type "timestamp":>
	if v := o.<:.NameTitle:>; !v.IsZero() {
		if s := uint(v.Unix()); s < 1<<32 {
			l += 9
		} else {
			l += 13
		}
	}
<:else if eq .Type "text" "binary":>
	if x := len(o.<:.NameTitle:>); x != 0 {
 <:- if .TypeList:>
		if x > ColferListMax {
			return -1, ColferMax(fmt.Sprintf("colfer: field <:.String:> exceeds %d elements", ColferListMax))
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
		for _, a := range o.<:.NameTitle:> {
			x = len(a)
			l += x
			for x >= 0x80 {
				x >>= 7
				l++
			}
			l++
		}
 <:- else:>
		l += x
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
 <:- end:>
	}
<:else if .TypeList:>
	if x := len(o.<:.NameTitle:>); x != 0 {
		if x > ColferListMax {
			return -1, ColferMax(fmt.Sprintf("colfer: field <:.String:> exceeds %d elements", ColferListMax))
		}
		for x >= 0x80 {
			x >>= 7
			l++
		}
		l += 2
		for _, v := range o.<:.NameTitle:> {
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
<:else:>
	if v := o.<:.NameTitle:>; v != nil {
		vl, err := v.MarshalLen()
		if err != nil {
			return -1, err
		}
		l += vl + 1
	}
<:end:>`

const goUnmarshalField = `<:if eq .Type "bool":>
	if header == <:.Index:> {
		o.<:.NameTitle:> = true
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}
<:else if eq .Type "uint32":>
	if header == <:.Index:> {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
		o.<:.NameTitle:> = x

		header = data[i]
		i++
	} else if header == <:.Index:>|0x80 {
		if i+4 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		o.<:.NameTitle:> = uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		header = data[i+4]
		i += 5
	}
<:else if eq .Type "uint64":>
	if header == <:.Index:> {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
		o.<:.NameTitle:> = x

		header = data[i]
		i++
	} else if header == <:.Index:>|0x80 {
		if i+8 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		o.<:.NameTitle:> = uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32 | uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		header = data[i+8]
		i += 9
	}
<:else if eq .Type "int32":>
	if header == <:.Index:> {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
		o.<:.NameTitle:> = int32(x)

		header = data[i]
		i++
	} else if header == <:.Index:>|0x80 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
		o.<:.NameTitle:> = int32(^x + 1)

		header = data[i]
		i++
	}
<:else if eq .Type "int64":>
	if header == <:.Index:> {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
		o.<:.NameTitle:> = int64(x)

		header = data[i]
		i++
	} else if header == <:.Index:>|0x80 {
		if i+1 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
		o.<:.NameTitle:> = int64(^x + 1)

		header = data[i]
		i++
	}
<:else if eq .Type "float32":>
	if header == <:.Index:> {
		if i+4 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		o.<:.NameTitle:> = math.Float32frombits(x)

		header = data[i+4]
		i += 5
	}
<:else if eq .Type "float64":>
	if header == <:.Index:> {
		if i+8 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		o.<:.NameTitle:> = math.Float64frombits(x)

		header = data[i+8]
		i += 9
	}
<:else if eq .Type "timestamp":>
	if header == <:.Index:> {
		if i+8 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		v := int64(x)
		o.<:.NameTitle:> = time.Unix(v>>32, v&(1<<32-1)).In(time.UTC)

		header = data[i+8]
		i += 9
	} else if header == <:.Index:>|0x80 {
		if i+12 >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		nsec := uint64(data[i+8])<<24 | uint64(data[i+9])<<16 | uint64(data[i+10])<<8 | uint64(data[i+11])
		o.<:.NameTitle:> = time.Unix(int64(sec), int64(nsec)).In(time.UTC)

		header = data[i+12]
		i += 13
	}
<:else if eq .Type "text":>
	if header == <:.Index:> {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
 <:- if .TypeList:>
		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> length %d exceeds %d elements", x, ColferListMax))
		}
		a := make([]string, int(x))
		o.<:.NameTitle:> = a
		for ai := range a {
			if i >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
							return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
				return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> element %d size %d exceeds %d bytes", ai, x, ColferSizeMax))
			}
			to := i + int(x)
			if to >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
 <:- else:>
		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> size %d exceeds %d bytes", x, ColferSizeMax))
		}
		to := i + int(x)
		if to >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		o.<:.NameTitle:> = string(data[i:to])

		header = data[to]
		i = to + 1
	}
 <:- end:>
<:else if eq .Type "binary":>
	if header == <:.Index:> {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
						return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
 <:- if not .TypeList:>
		if x > uint(ColferSizeMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> size %d exceeds %d bytes", x, ColferSizeMax))
		}
		l := int(x)
		if i + l >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		v := make([]byte, l)
		i += copy(v, data[i:])
		o.<:.NameTitle:> = v

		header = data[i]
		i++
 <:- else:>
		if x > uint(ColferListMax) {
			return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> length %d exceeds %d elements", x, ColferListMax))
		}
		a := make([][]byte, int(x))
		o.<:.NameTitle:> = a
		for ai := range a {
			if i >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
							return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
				return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> element %d size %d exceeds %d bytes", ai, x, ColferSizeMax))
			}
			v := make([]byte, int(x))

			if i + len(v) >= len(data) {
				if i >= ColferSizeMax {
					return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
				}
				return 0, io.EOF
			}

			i += copy(v, data[i:])
			a[ai] = v
		}

		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
 <:- end:>
	}
<:else if .TypeList:>
	if header == <:.Index:> {
		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
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
			return 0, ColferMax(fmt.Sprintf("colfer: <:.String:> length %d exceeds %d elements", x, ColferListMax))
		}

		l := int(x)
		a := make([]*<:.TypeNative:>, l)
		malloc := make([]<:.TypeNative:>, l)
		for ai, _ := range a {
			v := &malloc[ai]
			a[ai] = v

			n, err := v.Unmarshal(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
		o.<:.NameTitle:> = a

		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}
<:else:>
	if header == <:.Index:> {
		o.<:.NameTitle:> = new(<:.TypeNative:>)
		n, err := o.<:.NameTitle:>.Unmarshal(data[i:])
		if err != nil {
			return 0, err
		}
		i += n

		if i >= len(data) {
			if i >= ColferSizeMax {
				return 0, ColferMax(fmt.Sprintf("colfer: <:.Struct.String:> size %d exceeds %d bytes", i, ColferSizeMax))
			}
			return 0, io.EOF
		}
		header = data[i]
		i++
	}
<:end:>`
