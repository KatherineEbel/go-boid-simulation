package main

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth, screenHeight = 640, 360
	boidCount                 = 500
	viewRadius                = 13    // num pixels each boid can see
	adjRate                   = 0.015 // will determine how smooth simulation is
)

var (
	green = color.RGBA{
		R: 10,
		G: 255,
		B: 50,
		A: 255,
	}
	boids   [boidCount]*Boid
	boidMap [screenWidth + 1][screenHeight + 1]int
	rwLock  = sync.RWMutex{}
)

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		for _, boid := range boids {
			screen.Set(int(boid.position.x+1), int(boid.position.y), green)
			screen.Set(int(boid.position.x-1), int(boid.position.y), green)
			screen.Set(int(boid.position.x), int(boid.position.y-1), green)
			screen.Set(int(boid.position.x), int(boid.position.y+1), green)
		}
	}
	return nil
}
func main() {
	// set all boids initial positions on screen to -1
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Boids in a Box"); err != nil {
		log.Fatal(err)
	}
}
