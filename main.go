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

type GameType struct {
	World         WorldType
	Player        PlayerType
	FramesCounter int32
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

func (world *WorldType) DrawCell(pos rl.Vector2, color rl.Color) {
	rl.DrawRectangle(
		int32(pos.X + (world.CellSize.X * pos.X)),
		int32(pos.Y + (world.CellSize.Y * pos.Y)),
		int32(world.CellSize.X),
		int32(world.CellSize.Y),
		color,
	)
}

func (world *WorldType) Draw() {
	var i int32
	var j int32
	xl := int32(world.Size.X)
	yl := int32(world.Size.Y)

	for i = 0; i < xl; i++ {
		for j = 0; j < yl; j++ {
			if world.At(j, i) > 0 {
				world.DrawCell(rl.NewVector2(
					float32(j), float32(i)), rl.Red)
			}
		}
	}
}

func (player *PlayerType) RayCast(world WorldType) (d float32) {
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
				int32(y + c * math.Sin(angle))) > 0 {
					break
				}
		}

		plx := player.Position.X
		ply := player.Position.Y
		lx := float32(x + c * math.Cos(angle))
		ly := float32(y + c * math.Sin(angle))

		rl.DrawLine(
			int32(plx + (world.CellSize.X * plx)),
			int32(ply + (world.CellSize.Y * ply)),
			int32(lx + (world.CellSize.X * lx)),
			int32(ly + (world.CellSize.Y * ly)),
			rl.Red,
		)

		d = float32(c)
	}

	return
}

func (player *PlayerType) Draw(world WorldType) {
	world.DrawCell(player.Position, rl.DarkBlue)
	player.RayCast(world)
}

func (player *PlayerType) Update(world WorldType, fc int32) {
	if rl.IsKeyPressed(rl.KeyRight) && player.Speed.X == 0 {
		player.Speed = rl.NewVector2(1, 0)
	}
	if rl.IsKeyPressed(rl.KeyLeft) && player.Speed.X == 0 {
		player.Speed = rl.NewVector2(-1, 0)
	}
	if rl.IsKeyPressed(rl.KeyUp) && player.Speed.Y == 0 {
		player.Speed = rl.NewVector2(0, -1)
	}
	if rl.IsKeyPressed(rl.KeyDown) && player.Speed.Y == 0 {
		player.Speed = rl.NewVector2(0, 1)
	}

	if fc % 5 == 0 {
		player.Position.X += player.Speed.X
		player.Position.Y += player.Speed.Y
		player.Speed = rl.NewVector2(0, 0)
	}

	mouse := rl.GetMousePosition()

	player.A = (mouse.X / SCREEN_WIDTH) * (2*math.Pi)
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

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "RayMod")
	rl.SetTargetFPS(60)

	game.Player.Position.X = 0
	game.Player.Position.Y = 0
	game.Player.FOV = math.Pi / 3
}

func (game *GameType) Update() {
	game.Player.Update(game.World, game.FramesCounter)
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

