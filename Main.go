package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
	"math/rand"
)

const PlayerSpeed = 4

type Kuromi struct {
	pict   *ebiten.Image
	xLoc   int
	yLoc   int
	active bool
}
type LaserBeam struct {
	x, y   int
	dx, dy int
	speed  int
	active bool
}

type scrollDemo struct {
	player          *ebiten.Image
	xloc            int
	yloc            int
	dX              int
	background      *ebiten.Image
	backgroundXView int
	isMovingUp      bool
	isMovingDown    bool
	laserBeams      []LaserBeam
	laserBeamImage  *ebiten.Image
	enemy           []Kuromi
	score           int
}

func (demo *scrollDemo) Update() error {
	backgroundWidth := demo.background.Bounds().Dx()
	maxX := backgroundWidth * 2
	demo.backgroundXView -= 4
	demo.backgroundXView %= maxX
	demo.processPlayerInput()
	demo.processShootingInput()
	demo.updateLaserBeams()
	demo.handleEnemyCollisions()
	for i := range demo.enemy {
		if demo.enemy[i].active {
			demo.enemy[i].xLoc -= 2
			if demo.enemy[i].xLoc < 0 {
				demo.enemy[i].active = false
				demo.score--
			}

		}
	}

	return nil
}

func (demo *scrollDemo) Draw(screen *ebiten.Image) {
	drawOps := ebiten.DrawImageOptions{}
	const repeat = 3
	backgroundWidth := demo.background.Bounds().Dx()
	for count := 0; count < repeat; count += 1 {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(backgroundWidth*count),
			float64(-1000))
		drawOps.GeoM.Translate(float64(demo.backgroundXView), 0)
		screen.DrawImage(demo.background, &drawOps)
	}

	// Draw the player image
	drawOps.GeoM.Reset() // Reset transformations
	drawOps.GeoM.Translate(float64(demo.xloc), float64(demo.yloc))
	screen.DrawImage(demo.player, &drawOps)

	for _, beam := range demo.laserBeams {
		if beam.active {
			drawLaserBeam(screen, demo.laserBeamImage, float64(beam.x), float64(beam.y))
		}
	}

	for _, kuromi := range demo.enemy {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(kuromi.xLoc), float64(kuromi.yLoc))
		screen.DrawImage(kuromi.pict, &drawOps)
	}
}

func (s scrollDemo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Hello Kitty World <3")

	backgroundPict, _, err := ebitenutil.NewImageFromFile("pinkgamebackground-2.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}

	playerPict, _, err := ebitenutil.NewImageFromFile("hellokitty-4-2.png")
	if err != nil {
		fmt.Println("Unable to load player image:", err)
	}

	laserBeamImage, _, err := ebitenutil.NewImageFromFile("pinkLaser.png")
	if err != nil {
		log.Fatalf("Unable to load laser beam image: %v", err)
	}

	enemyPict, _, err := ebitenutil.NewImageFromFile("kuromi-2.png")

	allEnemies := make([]Kuromi, 0, 15)
	for i := 0; i < 10; i++ {
		allEnemies = append(allEnemies, NewKuromi(928, 928, enemyPict))
	}

	laserBeams := make([]LaserBeam, 0)

	demo := scrollDemo{
		player:         playerPict,
		background:     backgroundPict,
		xloc:           100,
		yloc:           500,
		dX:             0,
		laserBeamImage: laserBeamImage,
		laserBeams:     laserBeams,
		enemy:          allEnemies,
	}

	err = ebiten.RunGame(&demo)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}

func NewKuromi(MaxWidth int, MaxHeight int,
	image *ebiten.Image) Kuromi {
	return Kuromi{
		pict: image,
		xLoc: rand.Intn(MaxWidth),
		yLoc: rand.Intn(MaxHeight),
	}
}

func (demo *scrollDemo) processPlayerInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		demo.isMovingUp = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		demo.isMovingDown = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		demo.isMovingUp = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		demo.isMovingDown = false
	}
	if demo.isMovingUp {
		demo.yloc -= PlayerSpeed
	}
	if demo.isMovingDown {
		demo.yloc += PlayerSpeed
	}
}

func drawLaserBeam(screen *ebiten.Image, laserBeamImage *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(laserBeamImage, op)
}

func (demo *scrollDemo) processShootingInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		newBeam := LaserBeam{
			x:      demo.xloc, // Starting position
			y:      demo.yloc, // Starting position
			dx:     1,         // Adjust direction
			dy:     0,         // Adjust direction
			speed:  5,         // Adjust speed
			active: true,
		}
		demo.laserBeams = append(demo.laserBeams, newBeam)
	}
}

// I got this function from chatgpt
func (demo *scrollDemo) updateLaserBeams() {
	for i := len(demo.laserBeams) - 1; i >= 0; i-- {
		beam := &demo.laserBeams[i]
		if beam.active {
			beam.x += int(beam.dx * beam.speed)
			beam.y += int(beam.dy * beam.speed)

			if beam.x > 1000 {
				beam.active = false
			}
		} else {

			demo.laserBeams = append(demo.laserBeams[:i], demo.laserBeams[i+1:]...)
		}
	}
}
func (demo *scrollDemo) handleEnemyCollisions() {
	for i := range demo.enemy {
		if demo.enemy[i].active {
			for j := range demo.laserBeams {
				if demo.laserBeams[j].active {
					if collisionOccured(demo.enemy[i], demo.laserBeams[j]) {
						demo.enemy[i].active = false
						demo.laserBeams[j].active = false
						demo.score++ // Increase the score when an enemy is hit
					}
				}
			}
		}
	}
}

func collisionOccured(enemy Kuromi, beam LaserBeam) bool {
	enemyX := enemy.xLoc
	enemyY := enemy.yLoc
	enemyWidth := 928  // Set the width according to your enemy image
	enemyHeight := 928 // Set the height according to your enemy image

	beamX := beam.x
	beamY := beam.y
	beamWidth := 72  // Set the width according to your laser beam image
	beamHeight := 72 // Set the height according to your laser beam image

	// Check for collision by comparing positions and dimensions
	if enemyX < beamX+beamWidth && enemyX+enemyWidth > beamX &&
		enemyY < beamY+beamHeight && enemyY+enemyHeight > beamY {
		return true // Collision occurred
	}

	return false // No collision

}

//make slice in main
//process shooting beam and player input
