package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func neighbors(world [100][100]uint8, x int, y int) int {
	sum := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			sum += int(world[x+i][y+j])
		}
	}

	return sum - int(world[x][y])
}

func main() {
	rl.InitWindow(800, 800, "lenia-go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(30)

	const world_width, world_height = 100, 100

	var world [world_width][world_height]uint8
	var next_world [world_width][world_height]uint8

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			world[rl.GetMouseX()/8][rl.GetMouseY()/8] = 1
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			world[rl.GetMouseX()/8][rl.GetMouseY()/8] = 0
		}

		if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyDown(rl.KeySpace) {
			for i := 0; i < world_width; i++ {
				for j := 0; j < world_height; j++ {
					if i > 0 && i < len(world)-1 && j > 0 && j < len(world[0])-1 {
						N := neighbors(world, i, j) // count the neighbors

						// rules
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
			// update world state
			world = next_world

		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		for i := 0; i < world_width; i++ {
			for j := 0; j < world_height; j++ {
				if world[i][j] == 1 {
					rl.DrawRectangle(int32(i)*8, int32(j)*8, 8, 8, rl.White)
				}
			}
		}

		rl.EndDrawing()
	}
}
