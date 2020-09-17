package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{
			x: rand.Float64() * screenWidth,
			y: rand.Float64() * screenHeight,
		},
		velocity: Vector2D{
			x: (rand.Float64() * 2) - 1.0,
			y: (rand.Float64() * 2) - 1.0,
		},
		id: bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}

func (b *Boid) moveOne() {
	accel := b.calculateAcceleration() // make sure to cal acceleration first since it tries to acquire lock to thread also
	lock.Lock()
	b.velocity = b.velocity.Add(accel).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.Add(b.velocity)
	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{
			x: -b.velocity.x,
			y: b.velocity.y,
		}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{
			x: b.velocity.x,
			y: -b.velocity.y,
		}
	}
	lock.Unlock()
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) calculateAcceleration() Vector2D {
	// upper right and lower left of box
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPos, avgVel := Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0

	lock.Lock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVel = avgVel.Add(boids[otherBoidId].velocity)
					avgPos = avgPos.Add(boids[otherBoidId].position)
				}
			}
		}
	}
	lock.Unlock()
	accel := Vector2D{0, 0}
	if count > 0 {
		avgPos, avgVel = avgPos.DivisionV(count), avgVel.DivisionV(count)
		accelAlignment := avgVel.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPos.Subtract(b.position).MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion)
	}
	return accel
}
