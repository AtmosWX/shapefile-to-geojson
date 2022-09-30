package reader

import (
	"encoding/binary"
	"math"
)

func NewReader(data []byte) *Reader {
	return &Reader{data, 0, 0}
}

type Reader struct {
	Data   []byte
	Offset uint32
	Mode   Mode
}

type Mode int

const (
	BigEndian Mode = iota
	LittleEndian
)

func (r *Reader) Seek(offset uint32) {
	r.Offset = offset
}

func (r *Reader) StepForward(size uint32) {
	r.Offset += size
}

func (r *Reader) StepBackward(size uint32) {
	r.Offset -= size
}

func (r *Reader) SetMode(mode Mode) {
	r.Mode = mode
}

func (r *Reader) ReadBytes(size uint32) []byte {
	buffer := r.Data[r.Offset : r.Offset+size]
	r.StepForward(size)
	return buffer
}

func (r *Reader) ReadString(length uint32) string {
	bytes := r.ReadBytes(length)
	return string(bytes)
}

func (r *Reader) ReadUint() uint32 {
	bytes := r.ReadBytes(4)

	switch r.Mode {
	case LittleEndian:
		return binary.LittleEndian.Uint32(bytes)
	default:
		return binary.BigEndian.Uint32(bytes)
	}
}

func (r *Reader) ReadInt() int32 {
	bytes := r.ReadBytes(4)

	switch r.Mode {
	case LittleEndian:
		return int32(binary.LittleEndian.Uint32(bytes))
	default:
		return int32(binary.BigEndian.Uint32(bytes))
	}
}

func (r *Reader) ReadShortUint() uint16 {
	bytes := r.ReadBytes(2)

	switch r.Mode {
	case LittleEndian:
		return binary.LittleEndian.Uint16(bytes)
	default:
		return binary.BigEndian.Uint16(bytes)
	}
}

func (r *Reader) ReadShortInt() int16 {
	bytes := r.ReadBytes(2)

	switch r.Mode {
	case LittleEndian:
		return int16(binary.LittleEndian.Uint16(bytes))
	default:
		return int16(binary.BigEndian.Uint16(bytes))
	}
}

func (r *Reader) ReadFloat() float32 {
	bytes := r.ReadBytes(4)
	bits := uint32(0)

	switch r.Mode {
	case LittleEndian:
		bits = binary.LittleEndian.Uint32(bytes)
	default:
		bits = binary.BigEndian.Uint32(bytes)
	}

	return math.Float32frombits(bits)
}

func (r *Reader) ReadDouble() float64 {
	bytes := r.ReadBytes(8)
	bits := uint64(0)

	switch r.Mode {
	case LittleEndian:
		bits = binary.LittleEndian.Uint64(bytes)
	default:
		bits = binary.BigEndian.Uint64(bytes)
	}

	return math.Float64frombits(bits)
}

func (r *Reader) StaticReadBytes(size uint32) []byte {
	buffer := r.ReadBytes(size)
	r.StepBackward(size)
	return buffer
}

func (r *Reader) StaticReadString(length uint32) string {
	str := r.ReadString(length)
	r.StepBackward(length)
	return str
}

func (r *Reader) StaticReadUint() uint32 {
	num := r.ReadUint()
	r.StepBackward(4)
	return num
}

func (r *Reader) StaticReadInt() int32 {
	num := r.ReadInt()
	r.StepBackward(4)
	return num
}

func (r *Reader) StaticReadShortUint() uint16 {
	num := r.ReadShortUint()
	r.StepBackward(2)
	return num
}

func (r *Reader) StaticReadShortInt() int16 {
	num := r.ReadShortInt()
	r.StepBackward(2)
	return num
}

func (r *Reader) StaticReadFloat() float32 {
	num := r.ReadFloat()
	r.StepBackward(4)
	return num
}
