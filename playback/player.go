package playback

import (
	"github.com/SovaMC/replay/format"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	runtimeBase = 1 << 42
	eyeHeight   = 1.62
	jumpTicks   = 200
	tailTicks   = 20
)

type Viewer interface {
	Packet(pk packet.Packet)
	Online() bool
	World() *world.World
}

type Config struct {
	Data   *format.Data
	World  *world.World
	Viewer Viewer
	OnStop func()
}

type Player struct {
	conf    Config
	tick    int
	last    int
	paused  bool
	speed   float64
	acc     float64
	ghosts  map[uint32]*ghost
	stopped bool
}

func (conf Config) New() *Player {
	return &Player{conf: conf, last: conf.Data.Last(), speed: 1, ghosts: map[uint32]*ghost{}}
}

func (p *Player) Stop() {
	if p.stopped {
		return
	}
	p.stopped = true
	p.despawnAll()
	if p.conf.OnStop != nil {
		p.conf.OnStop()
	}
}

func (p *Player) Stopped() bool { return p.stopped }

func (p *Player) Paused() bool { return p.paused }

func (p *Player) Pause()  { p.paused = true }
func (p *Player) Resume() { p.paused = false }

func (p *Player) SetSpeed(v float64) { p.speed = v }

func (p *Player) Tick() int { return p.tick }

func (p *Player) Length() int { return p.last }

func (p *Player) Step() {
	if p.stopped {
		return
	}
	if !p.conf.Viewer.Online() || p.conf.Viewer.World() != p.conf.World {
		p.Stop()
		return
	}
	if p.paused {
		return
	}
	p.acc += p.speed
	for p.acc >= 1 {
		p.acc--
		for _, a := range p.conf.Data.Ticks[p.tick] {
			p.apply(a, true)
		}
		p.tick++
	}
	if p.tick > p.last+tailTicks {
		p.Stop()
	}
}

func (p *Player) Jump(target int) {
	target = max(0, min(target, p.last))
	p.despawnAll()
	p.ghosts = map[uint32]*ghost{}
	p.tick, p.acc = 0, 0
	for p.tick < target {
		for _, a := range p.conf.Data.Ticks[p.tick] {
			if a.Kind.Entity() {
				p.apply(a, false)
			}
		}
		p.tick++
	}
	for id, g := range p.ghosts {
		p.spawn(id, g)
	}
}

func (p *Player) Forward()  { p.Jump(p.tick + jumpTicks) }
func (p *Player) Backward() { p.Jump(p.tick - jumpTicks) }

func (p *Player) send(pks ...packet.Packet) {
	for _, pk := range pks {
		p.conf.Viewer.Packet(pk)
	}
}
