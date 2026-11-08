package format

func (w *writer) action(a Action) {
	w.u8(byte(a.Kind))
	switch a.Kind {
	case SpawnPlayer:
		w.u32(a.ID)
		w.str(a.Name)
		w.str(a.NameTag)
		w.i32(int32(a.Skin.Width))
		w.i32(int32(a.Skin.Height))
		w.bytes(a.Skin.Pix)
		w.str(a.Skin.Geometry)
		w.bytes(a.Skin.Model)
		w.i32(int32(a.Skin.CapeWidth))
		w.i32(int32(a.Skin.CapeHeight))
		w.bytes(a.Skin.Cape)
		w.vec(a.Pos)
		w.f32(a.Yaw)
		w.f32(a.Pitch)
		w.stack(a.Held)
		w.stack(a.Offhand)
		for _, s := range a.Armour {
			w.stack(s)
		}
	case SpawnEntity:
		w.u32(a.ID)
		w.str(a.Type)
		w.vec(a.Pos)
		w.f32(a.Yaw)
		w.f32(a.Pitch)
	case Despawn:
		w.u32(a.ID)
	case Move:
		w.u32(a.ID)
		w.vec(a.Pos)
		w.f32(a.Yaw)
		w.f32(a.Pitch)
		w.bool(a.Teleport)
	case Equip:
		w.u32(a.ID)
		w.stack(a.Held)
		w.stack(a.Offhand)
	case Armour:
		w.u32(a.ID)
		for _, s := range a.Armour {
			w.stack(s)
		}
	case Flags:
		w.u32(a.ID)
		w.bool(a.Sneaking)
		w.bool(a.Sprint)
		w.str(a.NameTag)
	case ActorEvent, Animate:
		w.u32(a.ID)
		w.i32(a.Event)
		w.i32(a.Data)
	case Block:
		w.vec(a.Pos)
		w.u32(a.Block)
	case LevelEvent:
		w.i32(a.Event)
		w.vec(a.Pos)
		w.i32(a.Data)
	case Sound:
		w.str(a.Name)
		w.vec(a.Pos)
		w.i32(a.Extra)
		w.str(a.Type)
		w.bool(a.Baby)
		w.bool(a.Global)
	}
}
