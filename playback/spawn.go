package playback

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (p *Player) spawn(id uint32, g *ghost) {
	rid := runtime(id)
	if !g.player {
		p.send(&packet.AddActor{EntityUniqueID: int64(rid), EntityRuntimeID: rid, EntityType: g.typ, Position: g.pos, Pitch: g.pitch, Yaw: g.yaw, HeadYaw: g.yaw, BodyYaw: g.yaw, EntityMetadata: metadata(g)})
		return
	}
	entry := protocol.PlayerListEntry{ActionType: protocol.PlayerListActionAdd, UUID: g.uuid, EntityUniqueID: int64(rid), Username: g.name, Skin: protocolSkin(g.skin)}
	p.send(
		&packet.PlayerList{Entries: []protocol.PlayerListEntry{entry}},
		&packet.AddPlayer{
			UUID: g.uuid, Username: g.name, EntityRuntimeID: rid, Position: g.pos, Pitch: g.pitch, Yaw: g.yaw, HeadYaw: g.yaw,
			HeldItem: instance(g.held), GameType: packet.GameTypeSurvival, EntityMetadata: metadata(g),
			AbilityData: protocol.AbilityData{EntityUniqueID: int64(rid), Layers: []protocol.AbilityLayer{{Type: protocol.AbilityLayerTypeBase, Abilities: protocol.AbilityCount - 1}}},
		},
		&packet.PlayerList{Entries: []protocol.PlayerListEntry{{ActionType: protocol.PlayerListActionRemove, UUID: g.uuid}}},
	)
	p.equip(rid, g)
	p.armour(rid, g)
}

func (p *Player) equip(rid uint64, g *ghost) {
	p.send(
		&packet.MobEquipment{EntityRuntimeID: rid, NewItem: instance(g.held), WindowID: protocol.WindowIDInventory},
		&packet.MobEquipment{EntityRuntimeID: rid, NewItem: instance(g.off), WindowID: protocol.WindowIDOffHand},
	)
}

func (p *Player) armour(rid uint64, g *ghost) {
	p.send(&packet.MobArmourEquipment{EntityRuntimeID: rid, Helmet: instance(g.armour[0]), Chestplate: instance(g.armour[1]), Leggings: instance(g.armour[2]), Boots: instance(g.armour[3])})
}

func (p *Player) despawn(id uint32) {
	p.send(&packet.RemoveActor{EntityUniqueID: int64(runtime(id))})
}

func (p *Player) despawnAll() {
	if !p.conf.Viewer.Online() {
		return
	}
	for id := range p.ghosts {
		p.despawn(id)
	}
}
