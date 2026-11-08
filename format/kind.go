package format

type Kind byte

const (
	SpawnPlayer Kind = iota
	SpawnEntity
	Despawn
	Move
	Equip
	Armour
	Flags
	ActorEvent
	Animate
	Block
	LevelEvent
	Sound
)

func (k Kind) Entity() bool {
	switch k {
	case SpawnPlayer, SpawnEntity, Despawn, Move, Equip, Armour, Flags:
		return true
	}
	return false
}
