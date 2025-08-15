package main

import (
	//os
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDTH = 800
	SCREEN_HEIGHT = 600
)

type WorldType struct {
	Map      []uint8
	Size     rl.Vector2
	CellSize rl.Vector2
}

type PlayerType struct {
	Position rl.Vector2
	Speed    rl.Vector2
	FOV      float32
	A        float32
}

type MouseType struct {
	Yaw float64
}

type GameType struct {
	World         WorldType
	Player        PlayerType
	FramesCounter int32
	Mouse         MouseType
}

func (world *WorldType) Setup() {
	world.Map = nil
	world.Map = []uint8{
		1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
		1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,
		1,0,0,1,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,1,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,1,1,1,1,0,0,0,0,0,0,0,0,1,
		1,0,1,1,0,0,0,0,0,0,0,0,0,0,0,1,
		1,0,0,1,0,0,0,0,0,0,0,0,0,0,0,1,
		1,0,0,1,1,1,1,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,
		1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
	}
	world.Size = rl.NewVector2(16, 16)
	world.CellSize = rl.NewVector2(
		float32(SCREEN_WIDTH) / world.Size.X,
		float32(SCREEN_HEIGHT) / world.Size.Y,
	)
}

func (world *WorldType) At(x, y int32) uint8 {
	return world.Map[int32(world.Size.X) * y + x]
}

func (world *WorldType) Draw() {
}

func (player *PlayerType) RayCast(world WorldType) {
	x := float64(player.Position.X)
	y := float64(player.Position.Y)

	var c float64

	for i := 0; i < SCREEN_WIDTH; i++ {
		angle := float64(
			(player.A - player.FOV/2) +
				player.FOV * float32(i)/float32(SCREEN_WIDTH));

		for c = 0; c < 20; c += 0.05 {
			if world.At(
				int32(x + c * math.Cos(angle)),
				int32(y + c * math.Sin(angle)),
			) > 0 {
				colh := SCREEN_HEIGHT / c
				rl.DrawRectangle(
					 int32(i),
					 int32(SCREEN_HEIGHT/2 - colh/2),
					 1, int32(colh), rl.Red,
				)
				break
			}
		}
	}

	return
}

func (player *PlayerType) Draw(world WorldType) {
	player.RayCast(world)
}

func (player *PlayerType) Update(world WorldType, fc int32, mouse MouseType) {
	if rl.IsKeyPressed(rl.KeyW) {
		player.Speed = rl.NewVector2(
			float32(math.Cos(float64(player.A))),
			float32(math.Sin(float64(player.A))),
		)
	}
	if rl.IsKeyPressed(rl.KeyS) {
		player.Speed = rl.NewVector2(
			float32(-math.Cos(float64(player.A))),
			float32(-math.Sin(float64(player.A))),
		)
	}

	if fc % 5 == 0 {
		player.Position.X += player.Speed.X
		player.Position.Y += player.Speed.Y
		player.Speed = rl.NewVector2(0, 0)
	}

	player.A = float32(mouse.Yaw)

	//mouse := rl.GetMousePosition()
	//player.A = (mouse.X / SCREEN_WIDTH) * (2 * math.Pi)
}

func (game *GameType) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	game.World.Draw()
	game.Player.Draw(game.World)

	rl.EndDrawing()
}

func (game *GameType) Kill() {
	rl.CloseWindow()
}

func (game *GameType) Setup() {
	// TODO: implement menus/settings
	game.World.Setup()

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "sillygame")
	rl.SetTargetFPS(60)

	rl.DisableCursor()

	rl.SetMousePosition(SCREEN_WIDTH/2, SCREEN_HEIGHT/2);

	game.Player.Position.X = 1
	game.Player.Position.Y = 1
	game.Player.FOV = math.Pi / 4
	game.Player.A = 0
}

func (mouse *MouseType) MouseUpdate() {
	mousePos := rl.GetMousePosition()
	deltaX := float64(mousePos.X - (SCREEN_WIDTH/2))

	sensitivity := float64(0.003)

	mouse.Yaw += deltaX * sensitivity

	// Wrap yaw
	if mouse.Yaw > math.Pi * 2 {
		mouse.Yaw -= math.Pi * 2
	}

	if mouse.Yaw < 0 {
		mouse.Yaw += math.Pi * 2
	}

	rl.SetMousePosition(SCREEN_WIDTH/2, SCREEN_HEIGHT/2);
}

func (game *GameType) HandleMouse() {
	game.Mouse.MouseUpdate()
}

func (game *GameType) Update() {
	game.HandleMouse()

	game.Player.Update(
		game.World, game.FramesCounter,
		game.Mouse,
	)

	game.FramesCounter++
}

func (game *GameType) Loop() {
	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}

func main() {
	var game GameType
	defer game.Kill()
	game.Setup()
	game.Loop()
}

