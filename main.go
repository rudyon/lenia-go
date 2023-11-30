package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const window_width, window_height = 800, 800
	const world_width, world_height = 50, 50
	const scaling_value = window_height / world_height

	rl.InitWindow(window_width, window_height, "lenia-go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(30)

	var world = make([][]float64, world_height)
	for i := range world {
		world[i] = make([]float64, world_width)
	}

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 1
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 0
		}

		if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyDown(rl.KeySpace) {
			world = updateWorld(world)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		for i := 0; i < len(world); i++ {
			for j := 0; j < len(world[0]); j++ {
				rl.DrawRectangle(int32(i)*scaling_value, int32(j)*scaling_value, scaling_value, scaling_value, rl.ColorAlpha(rl.White, float32(world[i][j])))
			}
		}

		rl.EndDrawing()
	}
}

func neighbors(world [][]float64, x, y, radius int) int {
	sum := 0
	rows, cols := len(world), len(world[0])

	for i := x - radius; i <= x+radius; i++ {
		for j := y - radius; j <= y+radius; j++ {
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
	nextWorld := make([][]float64, len(world))
	for i := range nextWorld {
		nextWorld[i] = make([]float64, len(world[0]))
	}

	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			N := neighbors(world, i, j, 4) // You can adjust the radius as needed
			B := neighbors(world, i, j, 2) // You can adjust the radius as needed
			S := world[i][j]
			nextWorld[i][j] = 0.99*S + 0.01*float64(N+B-N*B)
		}
	}

	return nextWorld
}
