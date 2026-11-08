package playback

import (
	"github.com/SovaMC/replay/format"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type ghost struct {
	player   bool
	uuid     uuid.UUID
	name     string
	typ      string
	skin     format.Skin
	pos      mgl32.Vec3
	yaw      float32
	pitch    float32
	held     format.Stack
	off      format.Stack
	armour   [4]format.Stack
	sneaking bool
	sprint   bool
	nameTag  string
}

func runtime(id uint32) uint64 {
	return runtimeBase + uint64(id)
}

func instance(s format.Stack) protocol.ItemInstance {
	st := s.Stack()
	if st.Empty() {
		return protocol.ItemInstance{}
	}
	rid, meta, _ := world.ItemRuntimeID(st.Item())
	inst := protocol.ItemInstance{Stack: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: rid, MetadataValue: uint32(meta)}, Count: 1}}
	if s.Enchanted {
		inst.Stack.NBTData = map[string]any{"ench": []any{map[string]any{"id": int16(0), "lvl": int16(1)}}}
	}
	return inst
}

func protocolSkin(s format.Skin) protocol.Skin {
	cfg := `{"geometry":{"default":"` + s.Geometry + `"}}`
	return protocol.Skin{
		SkinID:                    uuid.NewString(),
		SkinResourcePatch:         []byte(cfg),
		SkinImageWidth:            uint32(s.Width),
		SkinImageHeight:           uint32(s.Height),
		SkinData:                  s.Pix,
		CapeImageWidth:            uint32(s.CapeWidth),
		CapeImageHeight:           uint32(s.CapeHeight),
		CapeData:                  s.Cape,
		SkinGeometry:              s.Model,
		CapeID:                    uuid.NewString(),
		FullID:                    uuid.NewString(),
		Trusted:                   true,
		OverrideAppearance:        true,
		GeometryDataEngineVersion: []byte(protocol.CurrentVersion),
	}
}

func metadata(g *ghost) protocol.EntityMetadata {
	m := protocol.NewEntityMetadata()
	m[protocol.EntityDataKeyName] = g.nameTag
	if g.player {
		m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagAlwaysShowName)
		m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagShowName)
	}
	if g.sneaking {
		m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagSneaking)
	}
	if g.sprint {
		m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagSprinting)
	}
	m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagHasGravity)
	m[protocol.EntityDataKeyScale] = float32(1)
	return m
}
