package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	birdSize     = 40
	pipeWidth    = 80
	pipeGap      = 150
	pipeSpeed    = 200
)

type Game struct {
	birdX      float64
	birdY      float64
	birdVel    float64
	birdAngle  float64
	pipes      []Pipe
	score      int
	highScore  int
	gameState  string
	lastSpawn  time.Time
}

type Pipe struct {
	x          float64
	height     float64
	scored     bool
}

func (g *Game) Update() error {
	switch g.gameState {
	case "start":
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.gameState = "play"
		}
	case "play":
		// Bird physics
		g.birdVel += 600 * 1/60.0
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.birdVel = -300
		}
		g.birdY += g.birdVel * 1/60.0

		// Bird rotation
		targetAngle := g.birdVel * 0.15
		if targetAngle < -90 { targetAngle = -90 }
		if targetAngle > 30 { targetAngle = 30 }
		g.birdAngle = g.birdAngle*0.9 + targetAngle*0.1

		// Pipe spawning
		if time.Since(g.lastSpawn).Seconds() > 3 {
			g.pipes = append(g.pipes, Pipe{
				x:      screenWidth,
				height: 100 + rand.Float64()*(screenHeight-200-pipeGap),
			})
			g.lastSpawn = time.Now()
		}

		// Pipe movement and scoring
		for i := range g.pipes {
			g.pipes[i].x -= pipeSpeed * 1/60.0

			// Scoring
			if !g.pipes[i].scored && g.birdX > g.pipes[i].x {
				g.score++
				g.pipes[i].scored = true
				if g.score > g.highScore {
					g.highScore = g.score
				}
			}
		}

		// Remove off-screen pipes
		if len(g.pipes) > 0 && g.pipes[0].x < -pipeWidth {
			g.pipes = g.pipes[1:]
		}

		// Collision detection
		for _, pipe := range g.pipes {
			// Pipe collision check
			if g.birdX+birdSize/2 > pipe.x &&          // Bird right edge > pipe left
				g.birdX-birdSize/2 < pipe.x+pipeWidth && // Bird left edge < pipe right
				(g.birdY-birdSize/2 < pipe.height ||     // Bird above top pipe
				 g.birdY+birdSize/2 > pipe.height+pipeGap) { // Bird below bottom pipe
				g.gameState = "gameover"
			}
		}

		// Ground/ceiling collision
		if g.birdY > screenHeight-birdSize/2 || g.birdY < birdSize/2 {
			g.gameState = "gameover"
		}

	case "gameover":
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.reset()
			g.gameState = "play"
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{135, 206, 235, 255})

	// Bird
	bird := ebiten.NewImage(birdSize, birdSize)
	bird.Fill(color.RGBA{255, 255, 0, 255})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-birdSize/2, -birdSize/2)
	opts.GeoM.Rotate(g.birdAngle * math.Pi / 180)
	opts.GeoM.Translate(g.birdX, g.birdY)
	screen.DrawImage(bird, opts)

	// Pipes
	for _, pipe := range g.pipes {
		// Top pipe
		top := ebiten.NewImage(pipeWidth, int(pipe.height))
		top.Fill(color.RGBA{0, 255, 0, 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(pipe.x, 0)
		screen.DrawImage(top, opts)

		// Bottom pipe
		bottom := ebiten.NewImage(pipeWidth, int(screenHeight-pipe.height-pipeGap))
		bottom.Fill(color.RGBA{0, 255, 0, 255})
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(pipe.x, pipe.height+pipeGap)
		screen.DrawImage(bottom, opts)
	}

	// UI
	switch g.gameState {
	case "start":
		ebitenutil.DebugPrintAt(screen, "PRESS SPACE TO START", screenWidth/2-70, screenHeight/2)
	case "play":
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 10, 10)
	case "gameover":
		ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-40, screenHeight/2-20)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d  High Score: %d", g.score, g.highScore), screenWidth/2-80, screenHeight/2+10)
		ebitenutil.DebugPrintAt(screen, "PRESS SPACE TO RESTART", screenWidth/2-90, screenHeight/2+40)
	}
}

func (g *Game) Layout(int, int) (int, int) { return screenWidth, screenHeight }

func (g *Game) reset() {
	g.birdX = screenWidth / 2
	g.birdY = screenHeight / 2
	g.birdVel = 0
	g.pipes = nil
	g.score = 0
	g.lastSpawn = time.Now()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy Bird")
	game := &Game{gameState: "start"}
	game.reset()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
