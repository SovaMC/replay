package format

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Stack struct {
	Name      string
	Meta      int16
	Enchanted bool
}

func StackOf(s item.Stack) Stack {
	if s.Empty() {
		return Stack{}
	}
	name, meta := s.Item().EncodeItem()
	return Stack{Name: name, Meta: meta, Enchanted: len(s.Enchantments()) > 0}
}

func (s Stack) Stack() item.Stack {
	if s.Name == "" {
		return item.Stack{}
	}
	it, ok := world.ItemByName(s.Name, s.Meta)
	if !ok {
		return item.Stack{}
	}
	return item.NewStack(it, 1)
}
