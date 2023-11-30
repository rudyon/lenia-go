package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const window_width, window_height = 800, 800
	const world_width, world_height = 40, 40
	const scaling_value = window_height / world_height

	world := initWorld(world_height, world_width)

	rl.InitWindow(window_width, window_height, "lenia-go")
	rl.SetTargetFPS(30)

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 1
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 0
		}

		if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyDown(rl.KeySpace) {
			world = updateWorld(world)
		}

		drawWorld(world, scaling_value)
	}
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

func initWorld(world_height int, world_width int) [][]float64 {
	var world = make([][]float64, world_height)
	for i := range world {
		world[i] = make([]float64, world_width)
		for j := range world[i] {
			world[i][j] = float64(rand.Float64())
		}
	}
	return world
}

func neighbors(world [][]float64, x int, y int) int {
	sum := 0
	rows, cols := len(world), len(world[0])

	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i == x && j == y || i < 0 || j < 0 || i >= rows || j >= cols {
				continue
			}
			if world[i][j] == 1 {
				sum++
			}
		}
	}

	return sum
}

func updateWorld(world [][]float64) [][]float64 {
	next_world := make([][]float64, len(world))
	for i := range next_world {
		next_world[i] = make([]float64, len(world[0]))
	}

	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			if i >= 0 && i < len(world) && j >= 0 && j < len(world[0]) {
				N := neighbors(world, i, j)

				switch {
				case world[i][j] == 1 && (N == 2 || N == 3):
					next_world[i][j] = 1
				case world[i][j] == 0 && N == 3:
					next_world[i][j] = 1
				default:
					next_world[i][j] = 0
				}
			}
		}
	}

	return next_world
}
