package format

import "github.com/go-gl/mathgl/mgl32"

type Skin struct {
	Width, Height         int
	Pix                   []byte
	Geometry              string
	Model                 []byte
	CapeWidth, CapeHeight int
	Cape                  []byte
}

type Action struct {
	Kind Kind

	ID       uint32
	Type     string
	Name     string
	NameTag  string
	Skin     Skin
	Pos      mgl32.Vec3
	Yaw      float32
	Pitch    float32
	Teleport bool
	Held     Stack
	Offhand  Stack
	Armour   [4]Stack
	Sneaking bool
	Sprint   bool
	Event    int32
	Data     int32
	Extra    int32
	Block    uint32
	Baby     bool
	Global   bool
}
