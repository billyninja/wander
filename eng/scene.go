package eng

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"math/rand"
	"time"
)

type Window struct {
	Width, Height int32
}

type Camera struct {
	WX int32
	WY int32
}

type Holder interface {
	GetPos() sdl.Rect
}

type Scene struct {
	Window Window

	StartTime    time.Time
	FrameCounter uint64

	PC *Actor

	World         [][]Space
	Width, Height int32
	Cam           Camera
	CullM         []Holder
	Enemies       []*Actor
	GUIBlocks     []*GUIBlock

	TsTxt *sdl.Texture
	SsTxt *sdl.Texture
	Font  *ttf.Font

	ObstCount   uint16
	EnemyCount  uint16
	WidthCells  uint16
	HeightCells uint16
}

func (s *Scene) Init(renderer *sdl.Renderer) {

	sp1 := &SpriteSheet{
		Actions: Actions{
			WD: SpriteAction{
				sdl.Rect{0, 0, tileSize, tileSize},
				3,
			},
			WL: SpriteAction{
				sdl.Rect{0, 32, tileSize, tileSize},
				3,
			},
			WR: SpriteAction{
				sdl.Rect{0, 64, tileSize, tileSize},
				3,
			},
			WU: SpriteAction{
				sdl.Rect{0, 96, tileSize, tileSize},
				3,
			},
		},
		txt: s.SsTxt,
	}

	s.PC = &Actor{
		Object: Object{
			sdl.Rect{500, 500, tileSize, tileSize},
			gfTree,
		},
		Speed:      Vector2d{2, 2},
		Sprite:     sp1,
		CurrAction: sp1.WD,
	}

	s.WidthCells = uint16(s.Width / tileSize)
	s.HeightCells = uint16(s.Height / tileSize)

	/*s.World = make([][]Space, s.WidthCells)

	for i := 0; uint16(i) < s.WidthCells; i++ {
		s.World[i] = make([]Space, s.HeightCells)
		for j, _ := range s.World[i] {
			sp := Space{}
			sp.Gfxs[0] = gfGrass
			s.World[i][j] = sp
		}
	}*/

	s.World = LoadTMX("assets/world.tmx", renderer)

	cl1 := sdl.Color{20, 20, 20, 255}
	cl2 := sdl.Color{50, 50, 50, 255}
	cl3 := sdl.Color{255, 255, 255, 255}

	txt1 := TextEl{
		RelPos:  Vector2d{20, 20},
		Content: "asdqwe111",
		Font:    s.Font,
		Color:   cl3,
	}

	gel1 := GUIEl{
		Pos:     &sdl.Rect{120, 120, 200, 200},
		BGColor: cl1,
		Texts:   []TextEl{txt1},
	}
	gel2 := GUIEl{
		Pos:     &sdl.Rect{120, 420, 200, 200},
		BGColor: cl2,
	}

	gb1 := GUIBlock{
		Elements: []GUIEl{gel1, gel2},
		Pos:      &sdl.Rect{0, 0, 1920, 900},
	}

	s.GUIBlocks = append(s.GUIBlocks, &gb1)

	for _, block := range s.GUIBlocks {
		block.Baked = block.Bake(renderer)
	}

	// for i := 0; uint16(i) < s.ObstCount; i++ {
	// 	ci := rand.Int31n(int32(s.WidthCells))
	// 	cj := rand.Int31n(int32(s.HeightCells))
	// 	s.World[ci][cj].Gfxs[0] = gfWall
	// }

	for i := 0; uint16(i) < s.EnemyCount; i++ {
		ci := rand.Int31n(s.Width)
		cj := rand.Int31n(s.Height)
		s.Enemies = append(s.Enemies, &Actor{
			Object: Object{
				sdl.Rect{ci, cj, tileSize, tileSize},
				gfTree,
			},
			Speed:      Vector2d{1, 1},
			Sprite:     sp1,
			CurrAction: sp1.WD,
		})
	}
}

func (s *Scene) Update() {
	for _, e := range s.Enemies {
		pos := WorldToScreen(e.Pos, s)

		if pos.X > 0 && pos.Y > 0 && pos.X < s.Window.Width && pos.Y < s.Window.Height {
			e.UpdateAI(s)
		}
	}
}

func (s *Scene) Render(renderer *sdl.Renderer) {
	//tgt := renderer.GetRenderTarget()
	//tgt.SetBlendMode(sdl.BLENDMODE_BLEND)

	s.CullM = s.CullM[:0]

	var init int32 = 0
	var Source *sdl.Rect

	var ofX, ofY int32 = 32, 32

	renderer.SetDrawColor(0, 0, 0, 255)

	for sw := init; sw < s.Window.Width; sw += ofX {
		for sh := init; sh < s.Window.Height; sh += ofY {
			ofX = (tileSize - ((s.Cam.WX + sw) % tileSize))
			ofY = (tileSize - ((s.Cam.WY + sh) % tileSize))

			//fmt.Printf(" %d x %d |", ofX, ofY)

			var worldCellX uint16 = uint16((s.Cam.WX + sw) / tileSize)
			var worldCellY uint16 = uint16((s.Cam.WY + sh) / tileSize)

			// TODO TRATAR VALORES NEGATIVOS
			if worldCellX < 0 || worldCellX > s.WidthCells || worldCellY < 0 || worldCellY > s.HeightCells {
				renderer.FillRect(&sdl.Rect{sw, sh, ofX, ofY})
				continue
			}

			rect := Object{
				Pos: sdl.Rect{sw, sh, ofX, ofY},
			}

			for _, gfx := range s.World[worldCellX][worldCellY].Gfxs {

				if gfx != nil {

					if gfx.Txtr == nil {
						continue
					}

					if ofX != int32(tileSize) || ofY != int32(tileSize) {
						Source = &sdl.Rect{gfx.Source.X + (tileSize - ofX), gfx.Source.Y + (tileSize - ofY), ofX, ofY}
					} else {
						Source = gfx.Source
					}

					renderer.Copy(s.TsTxt, Source, &rect.Pos)
				}
			}

			// Updating CullM
			if s.World[worldCellX][worldCellY].Coll {
				s.CullM = append(s.CullM, &rect)
			}
		}
	}

	//drawCullMap(renderer)
	for _, e := range s.Enemies {
		pos := WorldToScreen(e.Pos, s)

		if pos.X > 0 && pos.Y > 0 && pos.X < s.Window.Width && pos.Y < s.Window.Height {
			renderer.Copy(s.SsTxt, e.GetPose(), &pos)
			s.CullM = append(s.CullM, e)
		}
	}

	pos := WorldToScreen(s.PC.Pos, s)
	renderer.Copy(s.SsTxt, s.PC.GetPose(), &pos)

	// First get the surface
	fps := fmt.Sprintf("%v", s.GetFPS())

	surface, _ := s.Font.RenderUTF8_Solid(fps, sdl.Color{255, 255, 255, 255})
	defer surface.Free()

	txtr, _ := renderer.CreateTextureFromSurface(surface)
	defer txtr.Destroy()
	renderer.Copy(txtr, &sdl.Rect{0, 0, surface.W, surface.H}, &sdl.Rect{0, 0, surface.W, surface.H})

	for _, gb := range s.GUIBlocks {
		renderer.Copy(gb.Baked, gb.Pos, gb.Pos)
	}

	renderer.Present()
}

func (s *Scene) GetFPS() float64 {
	x := time.Since(s.StartTime)
	return (float64(s.FrameCounter) / x.Seconds())
}
