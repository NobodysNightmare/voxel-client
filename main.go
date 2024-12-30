package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type Box struct {
	voxels []Voxel
}

type Voxel struct {
	R, G, B int
}

type Effect interface {
	Progress() *Box
	String() string
}

const width = 20
const height = 20
const depth = 12

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please specify two args: Remote as host:port and effect (twinkle|firework). E.g.")
		fmt.Println("    ./voxel-client voxelbox.coderdojo.red:5005 firework")
		os.Exit(1)
	}

	remote := os.Args[1]

	conn, err := net.Dial("udp", remote)
	handleErr(err)

	var effect Effect
	switch os.Args[2] {
	case "twinkle":
		effect = NewStarTwinkle()
	case "firework":
		effect = NewFirework(6, 15)
	}

	if effect == nil {
		fmt.Println("Unknown effect name", os.Args[2])
		os.Exit(1)
	}

	refreshInterval := time.Second / 30

	fmt.Println("Rendering effect", effect)

	frameCounter := 0
	lastCheckpoint := time.Now()

	for {
		frameStart := time.Now()
		box := effect.Progress()
		_, err = conn.Write(box.ToBytes())
		handleErr(err)

		frameDuration := time.Since(frameStart)
		time.Sleep(refreshInterval - frameDuration)

		frameCounter++
		if frameCounter%100 == 0 {
			currentCheckpoint := time.Now()
			fmt.Printf("\rCurrent FPS: %5.1f", 100/currentCheckpoint.Sub(lastCheckpoint).Seconds())
			lastCheckpoint = currentCheckpoint
		}
	}
}

func handleErr(err error) {
	if err == nil {
		return
	}

	fmt.Println(err)
	os.Exit(1)
}

func NewBox() Box {
	// Note: hardcoding to 20 in all dimensions here as per protocol (even if the physical thing is not as deep)
	v := make([]Voxel, 20*20*20)
	return Box{voxels: v}
}

func (b *Box) ToBytes() []byte {
	result := make([]byte, 0, 20*20*20*3)
	for _, v := range b.voxels {
		result = append(result, v.ToBytes()...)
	}

	return result
}

func (b *Box) At(x, y, z int) *Voxel {
	if x < 0 || y < 0 || z < 0 || x > width-1 || y > depth-1 || z > height-1 {
		return nil
	}

	return &b.voxels[x+width*(height-z-1)+width*height*y]
}

func (v *Voxel) Set(r, g, b int) {
	v.R = r
	v.G = g
	v.B = b
}

func (v *Voxel) Add(r, g, b int) {
	v.R = min(255, max(0, v.R+r))
	v.G = min(255, max(0, v.G+g))
	v.B = min(255, max(0, v.B+b))
}

func (v *Voxel) ToBytes() []byte {
	return []byte{byte(v.R), byte(v.G), byte(v.B)}
}
