package main

import (
	"math"
	"math/rand"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const width, height int32 = 800, 800
	const world_width, world_height int32 = 100, 100
	const scaling_value int32 = height / world_height

	world := initWorld(world_height, world_width)

	rl.InitWindow(width, height, "lenia-go")
	rl.SetConfigFlags(rl.FlagVsyncHint)

	for !rl.WindowShouldClose() {
		world = updateWorld(world, int(world_height), int(world_width))

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
			rl.DrawRectangleV(rl.NewVector2(float32(i*scaling_value), float32(j*scaling_value)), rl.NewVector2(float32(scaling_value), float32(scaling_value)), color)
		}
	}

	rl.EndDrawing()
}

func initWorld(world_height, world_width int32) [][]float64 {
	var world = make([][]float64, world_height)
	for i := range world {
		world[i] = make([]float64, world_width)
		for j := range world[i] {
			world[i][j] = float64(rand.Float64())
		}
	}
	return world
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

	var wg sync.WaitGroup
	for i := 0; i < len(world); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < len(world[0]); j++ {
				outer_kernel := calculateOuterKernel(world, i, j, 11, width, height)
				inner_kernel := calculateInnerKernel(world, i, j, 11/3, width, height)

				next_world[i][j] = s(outer_kernel, inner_kernel, b1, b2, d1, d2)
			}
		}(i)
	}
	wg.Wait()

	for i := range world {
		for j := range world[i] {
			world[i][j] += 1 * next_world[i][j]
			world[i][j] = clamp(world[i][j], 1, 0)
		}
	}
	return world
}

func calculateKernel(world [][]float64, x, y, radius int, width, height int, wrap func(int, int) (int, int), skipCenter bool) float64 {
	sum := 0.0
	totalWeight := 0.0

	for i := x - radius; i <= x+radius; i++ {
		for j := y - radius; j <= y+radius; j++ {
			// Wrap around the edges
			ii, jj := wrap(i, j)

			if !(skipCenter && ii == x && jj == y) {
				distance := math.Sqrt(math.Pow(float64(ii-x), 2) + math.Pow(float64(jj-y), 2))
				weight := math.Exp(-0.5 * math.Pow(distance/float64(radius), 2))

				sum += weight * world[ii][jj]
				totalWeight += weight
			}
		}
	}

	normalizedSum := sum / totalWeight
	return normalizedSum
}

func calculateOuterKernel(world [][]float64, x, y, radius int, width, height int) float64 {
	wrap := func(i, j int) (int, int) {
		return euclidMod(i, width), euclidMod(j, height)
	}

	return calculateKernel(world, x, y, radius, width, height, wrap, true)
}

func calculateInnerKernel(world [][]float64, x, y, radius int, width, height int) float64 {
	wrap := func(i, j int) (int, int) {
		return euclidMod(i, width), euclidMod(j, height)
	}

	return calculateKernel(world, x, y, radius, width, height, wrap, false)
}
