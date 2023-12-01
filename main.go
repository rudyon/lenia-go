package main

import (
	"math"
	"math/rand"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	width         int32 = 800
	height        int32 = 800
	worldWidth    int32 = 100
	worldHeight   int32 = 100
	scalingValue  int32 = height / worldHeight
	randomSquares int   = 1
	squareSize    int   = 40
)

var (
	ra    float64 = 11
	alpha float64 = 0.028
	dt    float64 = 0.05

	b1 float64 = 0.278
	b2 float64 = 0.365
	d1 float64 = 0.267
	d2 float64 = 0.445
)

func main() {
	world := initWorld(worldHeight, worldWidth, randomSquares, squareSize)

	rl.InitWindow(width, height, "lenia-go")
	rl.SetConfigFlags(rl.FlagVsyncHint)

	for !rl.WindowShouldClose() {
		world = updateWorld(world, int(worldHeight), int(worldWidth))

		drawWorld(world, int(scalingValue))
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

func drawWorld(world [][]float64, scalingValue int) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			alpha := world[i][j]
			color := rl.ColorAlpha(rl.White, float32(alpha))
			rl.DrawRectangleV(rl.NewVector2(float32(i*scalingValue), float32(j*scalingValue)), rl.NewVector2(float32(scalingValue), float32(scalingValue)), color)
		}
	}

	rl.EndDrawing()
}

func initWorld(worldHeight, worldWidth int32, numSquares, squareSize int) [][]float64 {
	world := make([][]float64, worldHeight)
	for i := range world {
		world[i] = make([]float64, worldWidth)
	}

	// Initialize random squares
	for s := 0; s < numSquares; s++ {
		startX := rand.Intn(int(worldWidth) - squareSize + 1)
		startY := rand.Intn(int(worldHeight) - squareSize + 1)

		for i := 0; i < squareSize; i++ {
			for j := 0; j < squareSize; j++ {
				world[startY+j][startX+i] = rand.Float64()
			}
		}
	}

	return world
}

func sigma1(x, a float64) float64 {
	return 1.0 / (1.0 + math.Exp(-(x-a)*4/alpha))
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
	nextWorld := make([][]float64, len(world))
	for i := range nextWorld {
		nextWorld[i] = make([]float64, len(world[0]))
	}

	var wg sync.WaitGroup
	for i := 0; i < len(world); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < len(world[0]); j++ {
				outerKernel := calculateOuterKernel(world, i, j, int(ra-1), width, height)
				innerKernel := calculateInnerKernel(world, i, j, int(ra-1)/3, width, height)

				nextWorld[i][j] = 2*s(outerKernel, innerKernel, b1, b2, d1, d2) - 1
			}
		}(i)
	}
	wg.Wait()

	for i := range world {
		for j := range world[i] {
			world[i][j] += dt * nextWorld[i][j]
			world[i][j] = clamp(world[i][j], 1, 0)
		}
	}
	return world
}

func calculateOuterKernel(world [][]float64, x, y, radius int, width, height int) float64 {
	return calculateKernel(world, x, y, radius, width, height, false)
}

func calculateInnerKernel(world [][]float64, x, y, radius int, width, height int) float64 {
	return calculateKernel(world, x, y, radius, width, height, true)
}

func calculateKernel(world [][]float64, x, y, radius int, width, height int, skipCenter bool) float64 {
	sum := 0.0
	totalWeight := 0.0

	for i := euclidMod(x, width) - radius; i <= euclidMod(x, width)+radius; i++ {
		for j := euclidMod(y, height) - radius; j <= euclidMod(y, height)+radius; j++ {
			if !(skipCenter && i == x && j == y) {
				distance := math.Sqrt(math.Pow(float64(i-x), 2) + math.Pow(float64(j-y), 2))
				weight := math.Exp(-0.5 * math.Pow(distance/float64(radius), 2))

				sum += weight * world[euclidMod(i, width)][euclidMod(j, height)]
				totalWeight += weight
			}
		}
	}

	if totalWeight != 0 {
		normalizedSum := sum / totalWeight
		return normalizedSum
	}

	return 0.0
}
