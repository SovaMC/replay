package format

import (
	"bytes"
	"compress/zlib"
)

func (d *Data) Encode() ([]byte, error) {
	w := &writer{}
	w.b.WriteString(magic)
	w.u8(version)
	w.str(d.Map)
	w.i32(int32(d.Time))
	w.vec(d.Spawn)
	w.i64(d.Created)
	w.u32(uint32(d.Duration))
	w.u32(uint32(len(d.Ticks)))
	for tick := 0; tick <= d.Duration; tick++ {
		actions, ok := d.Ticks[tick]
		if !ok {
			continue
		}
		w.u32(uint32(tick))
		w.u32(uint32(len(actions)))
		for _, a := range actions {
			w.action(a)
		}
	}
	var out bytes.Buffer
	z, err := zlib.NewWriterLevel(&out, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := z.Write(w.b.Bytes()); err != nil {
		return nil, err
	}
	if err := z.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
