package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const width, height int32 = 800, 800
	const world_width, world_height int32 = 50, 50
	const scaling_value int32 = height / world_height

	world := initWorld(world_height, world_width, int64(time.Second))

	rl.InitWindow(width, height, "lenia-go")
	rl.SetTargetFPS(30)

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 1
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 0
		}

		if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyDown(rl.KeySpace) {
			fmt.Println("updating")
			world = updateWorld(world, int(world_height), int(world_width))
			fmt.Println("done")
		}

		drawWorld(world, int(scaling_value))
	}
}

func euclidMod(a, b int) int {
	return (a%b + b) % b
}

func clamp(x, max, min float64) float64 {
	if x > max {
		return max
	} else if x < min {
		return min
	}
	return x
}

func drawWorld(world [][]float64, scaling_value int) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			alpha := world[i][j]
			color := rl.ColorAlpha(rl.White, float32(alpha))
			rl.DrawRectangle(int32(i*scaling_value), int32(j*scaling_value), int32(scaling_value), int32(scaling_value), color)
		}
	}

	rl.EndDrawing()
}

func initWorld(world_height, world_width int32, seed int64) [][]float64 {
	rand.Seed(seed)
	var world = make([][]float64, world_height)
	for i := range world {
		world[i] = make([]float64, world_width)
		for j := range world[i] {
			world[i][j] = float64(rand.Float64())
		}
	}
	return world
}

func kernel(world [][]float64, x, y, radius int, width, height int) float64 {
	sum := 0.0
	totalCells := 0

	for i := x - radius; i <= x+radius; i++ {
		for j := y - radius; j <= y+radius; j++ {
			distance := math.Sqrt(math.Pow(float64(euclidMod(i, width)-x), 2) + math.Pow(float64(euclidMod(j, height)-y), 2))
			if distance <= float64(radius) {
				sum += world[euclidMod(i, width)][euclidMod(j, height)]
				totalCells++
			}
		}
	}

	normalizedSum := sum / float64(totalCells)
	return normalizedSum
}

func sigma1(x, a float64) float64 {
	return 1.0 / +math.Exp(-(x-a)*4/0.028)
}

func sigma2(x, a, b float64) float64 {
	return sigma1(x, a) * (1 - sigma1(x, b))
}

func sigmam(x, y, m float64) float64 {
	return x*(1-sigma1(m, 0.5)) + y*sigma1(m, 0.5)
}

func s(n, m, b1, b2, d1, d2 float64) float64 {
	return sigma2(n, sigmam(b1, d1, m), sigmam(b2, d2, m))
}

func updateWorld(world [][]float64, width, height int) [][]float64 {
	const b1, b2, d1, d2 float64 = 0.278, 0.365, 0.267, 0.445
	next_world := make([][]float64, len(world))
	for i := range next_world {
		next_world[i] = make([]float64, len(world[0]))
	}

	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			outer_kernel := kernel(world, i, j, 21, width, height)
			inner_kernel := kernel(world, i, j, 21/3, width, height)

			next_world[i][j] = s(outer_kernel, inner_kernel, b1, b2, d1, d2)
		}
	}

	for i := range world {
		for j := range world[i] {
			world[i][j] += 0.1 * next_world[i][j]
			world[i][j] = clamp(world[i][j], 1, 0)
		}
	}
	return world
}
