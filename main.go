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

	var world = make([][]uint8, world_height)
	for i := range world {
	    world[i] = make([]uint8, world_width)
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
				if world[i][j] == 1 {
					rl.DrawRectangle(int32(i)*scaling_value, int32(j)*scaling_value, scaling_value, scaling_value, rl.White)
				}
			}
		}

		rl.EndDrawing()
	}
}

func neighbors(world [][]uint8, x int, y int) int {
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

func updateWorld(world [][]uint8) [][]uint8 {
	next_world := make([][]uint8, len(world))
	for i := range next_world {
		next_world[i] = make([]uint8, len(world[0]))
	}

	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			if i > 0 && i < len(world)-1 && j > 0 && j < len(world[0])-1 {
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

