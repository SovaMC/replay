package format

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type writer struct {
	b bytes.Buffer
}

func (w *writer) u8(v byte) { w.b.WriteByte(v) }

func (w *writer) bool(v bool) {
	if v {
		w.u8(1)
		return
	}
	w.u8(0)
}

func (w *writer) u32(v uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	w.b.Write(buf[:])
}

func (w *writer) i32(v int32) { w.u32(uint32(v)) }

func (w *writer) i64(v int64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(v))
	w.b.Write(buf[:])
}

func (w *writer) f32(v float32) { w.u32(math.Float32bits(v)) }

func (w *writer) bytes(v []byte) {
	w.u32(uint32(len(v)))
	w.b.Write(v)
}

func (w *writer) str(v string) { w.bytes([]byte(v)) }

func (w *writer) vec(v mgl32.Vec3) {
	w.f32(v[0])
	w.f32(v[1])
	w.f32(v[2])
}

func (w *writer) stack(s Stack) {
	w.str(s.Name)
	w.i32(int32(s.Meta))
	w.bool(s.Enchanted)
}
