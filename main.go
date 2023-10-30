package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const window_width, window_height = 800, 800
	const world_width, world_height = 400, 400
	const scaling_value = window_height / world_height

	rl.InitWindow(window_width, window_height, "lenia-go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(30)

	var world = make([][]uint8, world_height)
	for i := range world {
		world[i] = make([]uint8, world_width)
		for j := range world[i] {
			world[i][j] = 0
		}
	}

	var next_world = make([][]uint8, len(world))
	copy(next_world, world)

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 1
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			world[rl.GetMouseX()/scaling_value][rl.GetMouseY()/scaling_value] = 0
		}

		if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyDown(rl.KeySpace) {
			world = updateWorld(world, next_world)

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
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			sum += int(world[x+i][y+j])
		}
	}

	return sum - int(world[x][y])
}

func updateWorld(world [][]uint8, next_world [][]uint8) [][]uint8 {
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			if i > 0 && i < len(world)-1 && j > 0 && j < len(world[0])-1 {
				N := neighbors(world, i, j)

				if world[i][j] == 1 {
					if N == 2 {
						next_world[i][j] = world[i][j]
					}
					if N == 3 {
						next_world[i][j] = world[i][j]
					}
					if N > 3 {
						next_world[i][j] = 0
					}
					if N < 2 {
						next_world[i][j] = 0
					}
				} else {
					if N == 3 {
						next_world[i][j] = 1
					}
				}

			}
		}
	}

	world = next_world
	return world
}
