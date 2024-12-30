package main

import (
	"math/rand/v2"
)

type StarTwinkle struct {
	current, target Box
}

func NewStarTwinkle() Effect {
	return &StarTwinkle{current: NewBox(), target: NewBox()}
}

func (e *StarTwinkle) String() string {
	return "Star Twinkle"
}

func (e *StarTwinkle) Progress() *Box {
	finished := true

	for i, _ := range e.current.voxels {
		level := e.current.voxels[i].R
		if level > e.target.voxels[i].R {
			e.current.voxels[i] = Voxel{level - 1, level - 1, level - 1}
			finished = false
		} else if level < e.target.voxels[i].R {
			e.current.voxels[i] = Voxel{level + 1, level + 1, level + 1}
			finished = false
		}
	}

	if finished {
		e.target = NewBox()
		for i, _ := range e.target.voxels {
			if rand.Float32() < 0.1 {
				e.target.voxels[i] = Voxel{255, 255, 255}
			}
		}
	}

	return &e.current
}
