// author: João Maia

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"os"
	"time"
	"wander/eng"
)

var (
	winTitle   string = "Go-SDL2 Wander"
	imageName  string = "assets/textures/ts1.png"
	event      sdl.Event
	currScene  *eng.Scene
	sp1        *eng.SpriteSheet
	camDZUp    sdl.Rect = sdl.Rect{0, 0, winWidth, camDZO}
	camDZDown  sdl.Rect = sdl.Rect{0, winHeight - camDZO, winWidth, camDZO}
	camDZLeft  sdl.Rect = sdl.Rect{0, 0, camDZO, winHeight}
	camDZRight sdl.Rect = sdl.Rect{winWidth - camDZO, 0, camDZO, winHeight}

	font *ttf.Font
)

const (
	winWidth, winHeight int32 = 1920, 1040
	tileSize            int32 = 32
	camOffSetX          int32 = (winWidth / 2) + (tileSize / 2)
	camOffSetY          int32 = (winHeight / 2) + (tileSize / 2)
	camDZO              int32 = 120
)

func handleKeyEvent(key sdl.Keycode) {
	//fmt.Printf("%d pressed \n", key)
	var m eng.Vector2d

	switch key {
	case 1073741906:
		m = eng.MVSpeed.Up
		if eng.CheckColl(eng.WorldToScreen(currScene.PC.Pos, currScene), camDZUp, eng.Vector2d{}) {
			currScene.Cam.WY += (m.Y * currScene.PC.Speed.Y)
		}
	case 1073741905:
		m = eng.MVSpeed.Down
		if eng.CheckColl(eng.WorldToScreen(currScene.PC.Pos, currScene), camDZDown, eng.Vector2d{}) {
			currScene.Cam.WY += (m.Y * currScene.PC.Speed.Y)
		}
	case 1073741904:
		m = eng.MVSpeed.Left
		if eng.CheckColl(eng.WorldToScreen(currScene.PC.Pos, currScene), camDZLeft, eng.Vector2d{}) {
			currScene.Cam.WX += (m.X * currScene.PC.Speed.X)
		}
	case 1073741903:
		m = eng.MVSpeed.Right
		if eng.CheckColl(eng.WorldToScreen(currScene.PC.Pos, currScene), camDZRight, eng.Vector2d{}) {
			currScene.Cam.WX += (m.X * currScene.PC.Speed.X)
		}
	}

	currScene.PC.Move(m, currScene)

}

func catchEvents() bool {
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyDownEvent:
			handleKeyEvent(t.Keysym.Sym)
		}
	}
	return true
}

// Initialize
func main() {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var err error

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(winWidth), int(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}
	defer renderer.Destroy()

	tilesetImg, err := img.Load("assets/textures/ts1.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		os.Exit(3)
	}
	defer tilesetImg.Free()

	tilesetTxt, err := renderer.CreateTextureFromSurface(tilesetImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(4)
	}
	defer tilesetTxt.Destroy()

	spritesheetImg, err := img.Load("assets/textures/actor3.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		os.Exit(3)
	}
	defer spritesheetImg.Free()

	spritesheetTxt, err := renderer.CreateTextureFromSurface(spritesheetImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(4)
	}
	defer spritesheetTxt.Destroy()

	err = ttf.Init()
	font, err = ttf.OpenFont("assets/textures/PressStart2P.ttf", 18)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load font file: %s\n", err)
		os.Exit(6)
	}

	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load music file: %s\n", err)
		os.Exit(7)
	}

	if err := mix.Init(mix.INIT_MP3); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load music file: %s\n", err)
		os.Exit(8)
	}

	/*
		mus, err := mix.LoadMUS("assets/audio/test.mp3")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open music file: %s\n", err)
				os.Exit(9)
			}
			println(&mus)
				err = mus.Play(-1)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Failed to play music: %s\n", err)
						os.Exit(10)
					}*/

	currScene = &eng.Scene{
		Window:     eng.Window{winWidth, winHeight},
		StartTime:  time.Now(),
		Width:      2999,
		Height:     2999,
		Cam:        eng.Camera{0, 0},
		EnemyCount: 220,
		TsTxt:      tilesetTxt,
		SsTxt:      spritesheetTxt,
		Font:       font,
	}

	currScene.Init(renderer)

	var running bool = true

	for running {
		then := time.Now()
		running = catchEvents()
		currScene.Update()
		currScene.Render(renderer)

		dur := time.Since(then)
		sdl.Delay(40 - uint32(dur.Nanoseconds()/1000000))
		currScene.FrameCounter++
	}
}
