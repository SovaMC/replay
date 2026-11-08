package format

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Decode(blob []byte) (*Data, error) {
	z, err := zlib.NewReader(bytes.NewReader(blob))
	if err != nil {
		return nil, ErrFormat
	}
	raw, err := io.ReadAll(z)
	if err != nil {
		return nil, ErrFormat
	}
	if !bytes.HasPrefix(raw, []byte(magic)) || len(raw) <= len(magic) || raw[len(magic)] != version {
		return nil, ErrFormat
	}
	r := &reader{r: bytes.NewReader(raw[len(magic)+1:])}
	d := &Data{Ticks: map[int][]Action{}}
	d.Map = r.str()
	d.Time = int(r.i32())
	d.Spawn = r.vec()
	d.Created = r.i64()
	d.Duration = int(r.u32())
	count := int(r.u32())
	for range count {
		tick := int(r.u32())
		n := int(r.u32())
		if r.err != nil || n < 0 {
			return nil, ErrFormat
		}
		for range n {
			d.Ticks[tick] = append(d.Ticks[tick], r.action())
		}
	}
	if r.err != nil {
		return nil, ErrFormat
	}
	return d, nil
}
