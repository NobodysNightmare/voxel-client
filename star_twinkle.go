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
	speed := 5

	for i, _ := range e.current.voxels {
		level := e.current.voxels[i].R
		target := e.target.voxels[i].R
		nextLevel := level
		if level > target {
			nextLevel = max(level-speed, target)
			finished = false
		} else if level < target {
			nextLevel = min(level+speed, target)
			finished = false
		}

		e.current.voxels[i] = Voxel{nextLevel, nextLevel, nextLevel}
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
