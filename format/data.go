package format

import (
	"errors"
	"slices"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	magic   = "SOVAREPLAY"
	version = 1
)

var ErrFormat = errors.New("replay: unsupported format")

type Data struct {
	Map      string
	Time     int
	Spawn    mgl32.Vec3
	Created  int64
	Duration int
	Ticks    map[int][]Action
}

func (d *Data) Last() int {
	ticks := make([]int, 0, len(d.Ticks))
	for t := range d.Ticks {
		ticks = append(ticks, t)
	}
	if len(ticks) == 0 {
		return 0
	}
	return slices.Max(ticks)
}
