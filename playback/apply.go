package playback

import (
	"math"

	"github.com/SovaMC/replay/format"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (p *Player) apply(a format.Action, render bool) {
	switch a.Kind {
	case format.SpawnPlayer:
		g := &ghost{player: true, uuid: uuid.New(), name: a.Name, skin: a.Skin, pos: a.Pos, yaw: a.Yaw, pitch: a.Pitch, held: a.Held, off: a.Offhand, armour: a.Armour, nameTag: a.NameTag}
		p.ghosts[a.ID] = g
		if render {
			p.spawn(a.ID, g)
		}
	case format.SpawnEntity:
		g := &ghost{typ: a.Type, pos: a.Pos, yaw: a.Yaw, pitch: a.Pitch}
		p.ghosts[a.ID] = g
		if render {
			p.spawn(a.ID, g)
		}
	case format.Despawn:
		if _, ok := p.ghosts[a.ID]; ok {
			delete(p.ghosts, a.ID)
			if render {
				p.despawn(a.ID)
			}
		}
	case format.Move:
		p.move(a, render)
	case format.Equip:
		if g, ok := p.ghosts[a.ID]; ok && g.player {
			g.held, g.off = a.Held, a.Offhand
			if render {
				p.equip(runtime(a.ID), g)
			}
		}
	case format.Armour:
		if g, ok := p.ghosts[a.ID]; ok && g.player {
			g.armour = a.Armour
			if render {
				p.armour(runtime(a.ID), g)
			}
		}
	case format.Flags:
		if g, ok := p.ghosts[a.ID]; ok {
			g.sneaking, g.sprint, g.nameTag = a.Sneaking, a.Sprint, a.NameTag
			if render {
				p.send(&packet.SetActorData{EntityRuntimeID: runtime(a.ID), EntityMetadata: metadata(g)})
			}
		}
	case format.ActorEvent:
		if _, ok := p.ghosts[a.ID]; ok && render {
			p.send(&packet.ActorEvent{EntityRuntimeID: runtime(a.ID), EventType: uint8(a.Event), EventData: a.Data})
		}
	case format.Animate:
		if _, ok := p.ghosts[a.ID]; ok && render {
			p.send(&packet.Animate{EntityRuntimeID: runtime(a.ID), ActionType: uint8(a.Event)})
		}
	case format.LevelEvent:
		p.send(&packet.LevelEvent{EventType: a.Event, Position: a.Pos, EventData: a.Data})
	case format.Sound:
		p.send(&packet.LevelSoundEvent{SoundType: a.Name, Position: a.Pos, ExtraData: a.Extra, EntityType: a.Type, BabyMob: a.Baby, DisableRelativeVolume: a.Global})
	case format.Block:
		p.block(a)
	}
}

func (p *Player) move(a format.Action, render bool) {
	g, ok := p.ghosts[a.ID]
	if !ok {
		return
	}
	g.pos, g.yaw, g.pitch = a.Pos, a.Yaw, a.Pitch
	if !render {
		return
	}
	rid := runtime(a.ID)
	if g.player {
		mode := byte(packet.MoveModeNormal)
		if a.Teleport {
			mode = packet.MoveModeTeleport
		}
		p.send(&packet.MovePlayer{EntityRuntimeID: rid, Position: a.Pos.Add(mgl32.Vec3{0, eyeHeight}), Pitch: a.Pitch, Yaw: a.Yaw, HeadYaw: a.Yaw, Mode: mode})
		return
	}
	var flags byte
	if a.Teleport {
		flags = packet.MoveFlagTeleport
	}
	p.send(&packet.MoveActorAbsolute{EntityRuntimeID: rid, Position: a.Pos, Rotation: mgl32.Vec3{a.Pitch, a.Yaw, a.Yaw}, Flags: flags})
}

func (p *Player) block(a format.Action) {
	b, ok := world.BlockByRuntimeID(a.Block)
	if !ok {
		return
	}
	pos := cube.Pos{int(math.Floor(float64(a.Pos[0]))), int(math.Floor(float64(a.Pos[1]))), int(math.Floor(float64(a.Pos[2])))}
	p.conf.World.Do(func(tx *world.Tx) { tx.SetBlock(pos, b, nil) })
}
