package main

import (
	"math"
	"math/rand/v2"
	"slices"
)

type Firework struct {
	ParticleCount int
	Darkening     int

	box        Box
	particles  []*Particle
	frameCount int
}

type Particle struct {
	X, Y    int
	R, G, B int

	z, TargetZ           int
	radius, TargetRadius int

	Finished bool
}

func NewFirework(particleCount, darkening int) Effect {
	return &Firework{ParticleCount: particleCount, Darkening: darkening, box: NewBox()}
}

func (e *Firework) String() string {
	return "Firework"
}

func (e *Firework) Progress() *Box {
	if len(e.particles) < e.ParticleCount {
		// TODO: prevent black particles
		e.particles = append(e.particles, &Particle{
			X:            rand.IntN(width),
			Y:            rand.IntN(depth),
			TargetZ:      10 + rand.IntN(height-13),
			TargetRadius: 1 + rand.IntN(3),
			R:            allOrNone(),
			G:            allOrNone(),
			B:            allOrNone(),
		})
	}

	for i, _ := range e.box.voxels {
		e.box.voxels[i].Add(-e.Darkening, -e.Darkening, -e.Darkening)
	}

	for _, p := range e.particles {
		p.Render(&e.box)
		if !p.Finished {
			p.Progress(e.frameCount)
		}
	}

	e.particles = slices.DeleteFunc(e.particles, func(p *Particle) bool { return p.Finished })

	e.frameCount++
	return &e.box
}

func (p *Particle) Render(box *Box) {
	for x := p.X - p.radius; x <= p.X+p.radius; x++ {
		for y := p.Y - p.radius; y <= p.Y+p.radius; y++ {
			for z := p.z - p.radius; z <= p.z+p.radius; z++ {
				radius := math.Sqrt(float64((x-p.X)*(x-p.X) + (y-p.Y)*(y-p.Y) + (z-p.z)*(z-p.z)))
				if float64(p.radius) < radius {
					continue
				}

				voxel := box.At(x, y, z)
				if voxel != nil {
					voxel.Add(p.R, p.G, p.B)
				}
			}
		}
	}
}

func (p *Particle) Progress(frameCount int) {
	if p.z < p.TargetZ {
		if frameCount%3 == 0 {
			p.z++
		}
	} else {
		if p.radius < p.TargetRadius {
			if frameCount%3 == 0 {
				p.radius++
			}
		} else {
			p.Finished = true
		}
	}
}

func allOrNone() int {
	if rand.Float32() < 0.5 {
		return 0
	}

	return 255
}

func abs(i int) int {
	if i > 0 {
		return i
	}

	return -i
}
