package main

import (
	"math"

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

		// if rl.IsKeyPressed(rl.KeyN) || rl.IsKeyDown(rl.KeySpace) {
		// 	world = updateWorld(world)
		// }

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

func logistic(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func aliveness(world [][]float64, x, y, radius int) float64 {
	weightSum := 0.0
	sum := 0.0
	rows, cols := len(world), len(world[0])
	center_x, center_y := float64(rows/2), float64(cols/2)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i == x && j == y {
				continue
			}
			dx, dy := float64(i)-center_x, float64(j)-center_y
			distance := math.Sqrt(dx*dx + dy*dy)
			weight := logistic((distance - float64(radius)) / 3.0)
			weightSum += weight
			sum += weight * world[i][j]
		}
	}

	return sum / weightSum
}

func neighbors(world [][]float64, x, y, radius int) float64 {
	weightSum := 0.0
	sum := 0.0
	rows, cols := len(world), len(world[0])
	center_x, center_y := float64(rows/2), float64(cols/2)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i == x && j == y {
				continue
			}
			dx, dy := float64(i)-center_x, float64(j)-center_y
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance < float64(3*radius) {
				weight := logistic((distance - float64(radius)) / 3.0)
				weightSum += weight
				sum += weight * world[i][j]
			}
		}
	}

	alivenessValue := (sum - logistic(0)) / (weightSum - logistic(0))
	scaledAliveness := (alivenessValue + 1) / 2 // scale the aliveness value to the range [0, 1]

	return scaledAliveness
}
