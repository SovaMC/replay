package format

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type reader struct {
	r   *bytes.Reader
	err error
}

func (r *reader) read(n int) []byte {
	if r.err != nil {
		return make([]byte, max(n, 0))
	}
	if n < 0 || n > r.r.Len() {
		r.err = io.ErrUnexpectedEOF
		return make([]byte, max(n, 0))
	}
	buf := make([]byte, n)
	_, r.err = io.ReadFull(r.r, buf)
	return buf
}

func (r *reader) u8() byte { return r.read(1)[0] }

func (r *reader) bool() bool { return r.u8() != 0 }

func (r *reader) u32() uint32 { return binary.LittleEndian.Uint32(r.read(4)) }

func (r *reader) i32() int32 { return int32(r.u32()) }

func (r *reader) i64() int64 { return int64(binary.LittleEndian.Uint64(r.read(8))) }

func (r *reader) f32() float32 { return math.Float32frombits(r.u32()) }

func (r *reader) bytes() []byte { return r.read(int(r.u32())) }

func (r *reader) str() string { return string(r.bytes()) }

func (r *reader) vec() mgl32.Vec3 {
	return mgl32.Vec3{r.f32(), r.f32(), r.f32()}
}

func (r *reader) stack() Stack {
	return Stack{Name: r.str(), Meta: int16(r.i32()), Enchanted: r.bool()}
}
