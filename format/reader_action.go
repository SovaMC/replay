package format

func (r *reader) action() Action {
	a := Action{Kind: Kind(r.u8())}
	switch a.Kind {
	case SpawnPlayer:
		a.ID = r.u32()
		a.Name = r.str()
		a.NameTag = r.str()
		a.Skin.Width = int(r.i32())
		a.Skin.Height = int(r.i32())
		a.Skin.Pix = r.bytes()
		a.Skin.Geometry = r.str()
		a.Skin.Model = r.bytes()
		a.Skin.CapeWidth = int(r.i32())
		a.Skin.CapeHeight = int(r.i32())
		a.Skin.Cape = r.bytes()
		a.Pos = r.vec()
		a.Yaw = r.f32()
		a.Pitch = r.f32()
		a.Held = r.stack()
		a.Offhand = r.stack()
		for i := range a.Armour {
			a.Armour[i] = r.stack()
		}
	case SpawnEntity:
		a.ID = r.u32()
		a.Type = r.str()
		a.Pos = r.vec()
		a.Yaw = r.f32()
		a.Pitch = r.f32()
	case Despawn:
		a.ID = r.u32()
	case Move:
		a.ID = r.u32()
		a.Pos = r.vec()
		a.Yaw = r.f32()
		a.Pitch = r.f32()
		a.Teleport = r.bool()
	case Equip:
		a.ID = r.u32()
		a.Held = r.stack()
		a.Offhand = r.stack()
	case Armour:
		a.ID = r.u32()
		for i := range a.Armour {
			a.Armour[i] = r.stack()
		}
	case Flags:
		a.ID = r.u32()
		a.Sneaking = r.bool()
		a.Sprint = r.bool()
		a.NameTag = r.str()
	case ActorEvent, Animate:
		a.ID = r.u32()
		a.Event = r.i32()
		a.Data = r.i32()
	case Block:
		a.Pos = r.vec()
		a.Block = r.u32()
	case LevelEvent:
		a.Event = r.i32()
		a.Pos = r.vec()
		a.Data = r.i32()
	case Sound:
		a.Name = r.str()
		a.Pos = r.vec()
		a.Extra = r.i32()
		a.Type = r.str()
		a.Baby = r.bool()
		a.Global = r.bool()
	default:
		r.err = ErrFormat
	}
	return a
}
