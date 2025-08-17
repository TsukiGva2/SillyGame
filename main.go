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
}

type PlayerType struct {
	Position    rl.Vector2
	MovingSpeed float32
	Speed       rl.Vector2
	FOV         float32
	A           float32
	Hand        rl.Texture2D
}

type MouseType struct {
	Yaw float64
}

type Shader struct {
	Loaded   rl.Shader
	Location int32
	Dest     rl.RenderTexture2D
}

type RendererType struct {
	Texture       rl.Texture2D
	TexturePixels int32
	Shader        Shader
}

type GameType struct {
	World         WorldType
	Player        PlayerType
	Renderer      RendererType
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
		1,0,0,1,1,2,2,2,2,0,0,0,0,0,0,1,
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
}

func (world *WorldType) At(x, y int32) uint8 {
	return world.Map[int32(world.Size.X) * y + x]
}

func (renderer *RendererType) DrawTextureSlice(
    index uint8,
    offset int32,
    size, pos rl.Vector2) {

    sourceRec := rl.NewRectangle(
        float32(index) * float32(renderer.TexturePixels) + float32(offset),
        0,
        1,
        float32(renderer.TexturePixels),
    )

    destRec := rl.NewRectangle(
        pos.X,
        pos.Y,
        size.X,
        size.Y,
    )

    rl.DrawTexturePro(renderer.Texture, sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.RayWhite)
}


func (renderer *RendererType) Raycaster(
	player PlayerType,
	world WorldType,
) {
	x := float64(player.Position.X)
	y := float64(player.Position.Y)

	for i := 0; i < SCREEN_WIDTH; i++ {
		angle := float64(
			(player.A - player.FOV/2) +
				player.FOV * float32(i)/float32(SCREEN_WIDTH));

		mapX := int32(x)
		mapY := int32(y)

		rayDirX := math.Cos(angle)
		rayDirY := math.Sin(angle)

		deltaDistX := 1e30
		if rayDirX != 0 {
			deltaDistX = math.Abs(1 / rayDirX)
		}
		deltaDistY := 1e30
		if rayDirY != 0 {
			deltaDistY = math.Abs(1 / rayDirY)
		}

		var sideDistX, sideDistY float64

		var stepX, stepY int32

		if rayDirX < 0 {
			stepX = -1
			sideDistX = (x - float64(mapX)) * deltaDistX
		} else {
			stepX = 1
			sideDistX = (float64(mapX) + 1.0 - x) * deltaDistX
		}

		if rayDirY < 0 {
			stepY = -1
			sideDistY = (y - float64(mapY)) * deltaDistY
		} else {
			stepY = 1
			sideDistY = (float64(mapY) + 1.0 - y) * deltaDistY
		}

		hit := 0
		var verticalHit bool

		for hit == 0 {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				verticalHit = true
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				verticalHit = false
			}

			collision := world.At(mapX, mapY)
			if collision > 0 {
				var perpWallDist float64
				var wallX float64

				if verticalHit {
					perpWallDist = (float64(mapX) - x + (1.0 - float64(stepX)) / 2.0) / rayDirX
					wallX = y + perpWallDist * rayDirY
				} else {
					perpWallDist = (float64(mapY) - y + (1.0 - float64(stepY)) / 2.0) / rayDirY
					wallX = x + perpWallDist * rayDirX
				}

				wallX -= math.Floor(wallX)

				colh := float64(SCREEN_HEIGHT) / perpWallDist

				renderer.DrawTextureSlice(
					collision, int32(wallX*float64(renderer.TexturePixels)),
					rl.NewVector2(1, float32(colh)),
					rl.NewVector2(float32(i),
						float32(SCREEN_HEIGHT/2-colh/2)))

				hit = 1
			}
		}
	}

	return
}

func (player *PlayerType) Update(world WorldType, fc int32, mouse MouseType) {
	deltaTime := rl.GetFrameTime()

	player.Speed = rl.NewVector2(0, 0)
	player.MovingSpeed = 3

	if rl.IsKeyDown(rl.KeyLeftShift) {
		player.MovingSpeed *= 2
	}

	if rl.IsKeyDown(rl.KeyW) {
		player.Speed.X += float32(math.Cos(float64(player.A)))
		player.Speed.Y += float32(math.Sin(float64(player.A)))
	}

	if rl.IsKeyDown(rl.KeyS) {
		player.Speed.X += float32(-math.Cos(float64(player.A)))
		player.Speed.Y += float32(-math.Sin(float64(player.A)))
	}

	if rl.IsKeyDown(rl.KeyD) {
		player.Speed.X += float32(-math.Sin(float64(player.A)))
		player.Speed.Y += float32(math.Cos(float64(player.A)))
	}

	if rl.IsKeyDown(rl.KeyA) {
		player.Speed.X += float32(math.Sin(float64(player.A)))
		player.Speed.Y += float32(-math.Cos(float64(player.A)))
	}

	if rl.Vector2Length(player.Speed) > 1.0 {
        player.Speed = rl.Vector2Normalize(player.Speed)
    }

	player.Position.X += player.Speed.X * player.MovingSpeed * deltaTime
	player.Position.Y += player.Speed.Y * player.MovingSpeed * deltaTime

	player.A = float32(mouse.Yaw)

	//mouse := rl.GetMousePosition()
	//player.A = (mouse.X / SCREEN_WIDTH) * (2 * math.Pi)
}

func (renderer *RendererType) Render(game *GameType) {

	time := float32(rl.GetTime())
	rl.SetShaderValue(
		renderer.Shader.Loaded, renderer.Shader.Location,
		[]float32{time}, rl.ShaderUniformFloat)

	rl.BeginTextureMode(renderer.Shader.Dest)

	rl.ClearBackground(rl.RayWhite)

	renderer.Raycaster(
		game.Player,
		game.World,
	)

	game.Player.DrawHand()

	rl.EndTextureMode()
}

func (player *PlayerType) DrawHand() {
	weaponX := SCREEN_WIDTH - player.Hand.Width
	weaponY := SCREEN_HEIGHT - player.Hand.Height
	rl.DrawTexture(player.Hand, weaponX, weaponY, rl.White)
}

func (game *GameType) DrawExtra() {
}

func (game *GameType) Draw() {

	game.Renderer.Render(game)

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.BeginShaderMode(game.Renderer.Shader.Loaded)

	rl.DrawTextureRec(
		game.Renderer.Shader.Dest.Texture,
		rl.NewRectangle(
			0, 0, float32(game.Renderer.Shader.Dest.Texture.Width),
			-float32(game.Renderer.Shader.Dest.Texture.Height)),
		rl.NewVector2(0, 0),
		rl.White,
	)

	rl.EndShaderMode()

	// extra draw step
	game.DrawExtra()

	rl.EndDrawing()
}

func (game *GameType) Kill() {
	rl.UnloadTexture(game.Player.Hand)
	rl.UnloadTexture(game.Renderer.Texture)
	rl.CloseWindow()
}

func (game *GameType) Setup() {
	game.World.Setup()
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "sillygame")

	// TODO: implement menus/settings

	game.Renderer.Shader.Loaded = rl.LoadShader("", "vhs.fs")
	game.Renderer.Shader.Location = rl.GetShaderLocation(
		game.Renderer.Shader.Loaded, "time")
	game.Renderer.Shader.Dest = rl.LoadRenderTexture(SCREEN_WIDTH, SCREEN_HEIGHT)

	rl.SetTargetFPS(60)

	rl.DisableCursor()

	rl.SetMousePosition(SCREEN_WIDTH/2, SCREEN_HEIGHT/2);

	game.Player.Position.X = 1
	game.Player.Position.Y = 1
	game.Player.FOV = math.Pi / 4
	game.Player.A = 0
	game.Player.MovingSpeed = 2

	game.Player.Hand = rl.LoadTexture("hand_overlay.png")

	game.Renderer.Texture = rl.LoadTexture("walltext.png")
	game.Renderer.TexturePixels = 64
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

