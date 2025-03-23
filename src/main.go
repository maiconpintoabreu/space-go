package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Vector2 struct {
	X float64
	Y float64
}
type GameState int

const (
	StateStartMenu GameState = iota
	StateInGame
	StateGameOver
)

type Player struct {
	position       Vector2
	topPoint       Vector2
	leftPoint      Vector2
	rightPoint     Vector2
	speed          Vector2
	rotation       float64
	acceleration   float64
	isTurnLeft     bool
	isTurnRight    bool
	isAccelerating bool
	isBreaking     bool
}
type Game struct {
	player                 Player
	keys                   []ebiten.Key
	width                  int32
	height                 int32
	fwidth                 float64
	fheight                float64
	halfWidth              float64
	halfHeight             float64
	frameTimeAccumulator   float64
	isPlayerRotationChange bool
	state                  GameState
}

var (
	screenWidth  int = 640
	screenHeight int = 360
)

const (
	PLAYER_SPEED          float64 = 100.0
	PLAYER_ROTATION_SPEED float64 = 100.0
	SHIP_HALF_HEIGHT      float64 = 5.0 / 0.363970
	ZERO_SPEED            float64 = 0
	PHYSICS_TIME          float64 = 0.02
	DEG2RAD               float64 = 0.01745
	FONT_SIZE             float64 = 10
	TITLE_FONT_SIZE       float64 = FONT_SIZE // Check if needed
)

var (
	textFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	textFaceSource = s
}
func (game *Game) Update() error {

	// Tick

	// Input
	game.keys = inpututil.AppendPressedKeys(game.keys[:0])
	game.player.isTurnLeft = false
	game.player.isTurnRight = false
	game.player.isAccelerating = false
	game.player.isBreaking = false
	for _, key := range game.keys {
		if key == ebiten.KeyLeft {
			game.player.isTurnLeft = true
			game.isPlayerRotationChange = true
		} else if key == ebiten.KeyRight {
			game.player.isTurnRight = true
			game.isPlayerRotationChange = true
		} else if key == ebiten.KeyUp {
			game.player.isAccelerating = true
		} else if key == ebiten.KeyDown {
			game.player.isBreaking = true
		}
	}

	// Physics

	if true {
		game.frameTimeAccumulator = 0 //TODO: Reduce times that physics runs -= PHYSICS_TIME

		var rotation_speed float64 = PLAYER_ROTATION_SPEED * PHYSICS_TIME
		var acceleration float64 = PLAYER_SPEED * PHYSICS_TIME

		if game.player.isTurnLeft {
			game.player.rotation -= rotation_speed
		} else if game.player.isTurnRight {
			game.player.rotation += rotation_speed
		}

		if game.isPlayerRotationChange {
			game.isPlayerRotationChange = false
			if game.player.rotation > 180.0 {
				game.player.rotation -= 360.0
			}
			if game.player.rotation < -180.0 {
				game.player.rotation += 360.0
			}
		}
		if game.player.isAccelerating {
			game.player.isAccelerating = true
			if game.player.acceleration < PLAYER_SPEED {
				game.player.acceleration += acceleration
			}
		} else if game.player.acceleration > ZERO_SPEED {
			game.player.acceleration -= acceleration / 2.0
		} else if game.player.acceleration < ZERO_SPEED {
			game.player.acceleration = ZERO_SPEED
		}
		if game.player.isBreaking {
			if game.player.acceleration > ZERO_SPEED {
				game.player.acceleration -= acceleration
			} else if game.player.acceleration < ZERO_SPEED {
				game.player.acceleration = ZERO_SPEED
			}
		}

		direction := Vector2{
			X: float64(math.Sin(game.player.rotation * DEG2RAD)),
			Y: float64(-math.Cos(game.player.rotation * DEG2RAD)),
		}
		norm_vector := Vector2Normalize(&direction)
		game.player.speed = Vector2Scale(&norm_vector, game.player.acceleration*PHYSICS_TIME)
		game.player.position = Vector2Add(&game.player.position, &game.player.speed)
		// Update Triangle Rotation
		if Vector2Length(&game.player.speed) > 0.0 {
			if game.player.position.X > game.fwidth+SHIP_HALF_HEIGHT {
				game.player.position.X = -SHIP_HALF_HEIGHT
			} else if game.player.position.X < -SHIP_HALF_HEIGHT {
				game.player.position.X = game.fwidth + SHIP_HALF_HEIGHT
			}

			if game.player.position.Y > game.fheight+SHIP_HALF_HEIGHT {
				game.player.position.Y = -SHIP_HALF_HEIGHT
			} else if game.player.position.Y < -SHIP_HALF_HEIGHT {
				game.player.position.Y = game.fheight + SHIP_HALF_HEIGHT
			}
		}
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	switch game.state {
	case StateInGame:
		// Draw In Game UI
		player_speed := fmt.Sprint("Speed: ", game.player.acceleration)
		op := &text.DrawOptions{}
		op.GeoM.Translate(20, 12)
		op.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, player_speed, &text.GoTextFace{
			Source: textFaceSource,
			Size:   TITLE_FONT_SIZE,
		}, op)
		fps := fmt.Sprint("FPS: ", math.Round(ebiten.ActualFPS()))

		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(game.width)-100, 12)
		op.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, fps, &text.GoTextFace{
			Source: textFaceSource,
			Size:   TITLE_FONT_SIZE,
		}, op)

		var cosf float64 = math.Cos(game.player.rotation * DEG2RAD)
		var sinf float64 = math.Sin(game.player.rotation * DEG2RAD)
		game.player.topPoint = Vector2{
			X: game.player.position.X + sinf*SHIP_HALF_HEIGHT,
			Y: game.player.position.Y - cosf*SHIP_HALF_HEIGHT,
		}
		// Temp vector to center the rotation
		v1tmp := Vector2{
			X: game.player.position.X - sinf*SHIP_HALF_HEIGHT,
			Y: game.player.position.Y + cosf*SHIP_HALF_HEIGHT,
		}
		game.player.rightPoint = Vector2{
			X: v1tmp.X - cosf*(SHIP_HALF_HEIGHT-2.0),
			Y: v1tmp.Y - sinf*(SHIP_HALF_HEIGHT-2.0),
		}
		game.player.leftPoint = Vector2{
			X: v1tmp.X + cosf*(SHIP_HALF_HEIGHT-2.0),
			Y: v1tmp.Y + sinf*(SHIP_HALF_HEIGHT-2.0),
		}
		vector.StrokeLine(screen, float32(game.player.topPoint.X), float32(game.player.topPoint.Y), float32(game.player.rightPoint.X), float32(game.player.rightPoint.Y), 1, color.Gray{}, false)
		vector.StrokeLine(screen, float32(game.player.rightPoint.X), float32(game.player.rightPoint.Y), float32(game.player.leftPoint.X), float32(game.player.leftPoint.Y), 1, color.Gray{}, false)
		vector.StrokeLine(screen, float32(game.player.leftPoint.X), float32(game.player.leftPoint.Y), float32(game.player.topPoint.X), float32(game.player.topPoint.Y), 1, color.Gray{}, false)

	case StateStartMenu:

	case StateGameOver:
	}

	// Draw Game
	// screen.DrawTriangle(
	// 	game.player.topPoint,
	// 	game.player.rightPoint,
	// 	game.player.leftPoint,
	// 	Gray,
	// )
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var game Game
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Space Go")

	game.width = int32(screenWidth)
	game.height = int32(screenHeight)
	game.fwidth = float64(game.width)
	game.fheight = float64(game.height)
	game.halfWidth = game.fwidth / 2.0
	game.halfHeight = game.fheight / 2.0

	game.player.position = Vector2{X: game.halfWidth, Y: game.halfHeight - (SHIP_HALF_HEIGHT / 2.0)}
	game.player.acceleration = 0.0
	game.frameTimeAccumulator = 0.0

	game.isPlayerRotationChange = false

	game.state = StateInGame //TODO: Start the game with Start Menu StateStartMenu

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}

}

func Vector2Length(v *Vector2) float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}
func Vector2Normalize(v *Vector2) Vector2 {
	length := Vector2Length(v)

	if length > 0 {
		ilength := 1.0 / length
		return Vector2{v.X * ilength, v.Y * ilength}
	}

	return Vector2{}
}
func Vector2Scale(v *Vector2, scale float64) Vector2 {
	return Vector2{v.X * scale, v.Y * scale}
}
func Vector2Add(v1 *Vector2, v2 *Vector2) Vector2 {
	return Vector2{v1.X + v2.X, v1.Y + v2.Y}
}
