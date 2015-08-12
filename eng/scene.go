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
	GetType() ObjectType
}

type CullingMap struct {
	PCs   []Holder // PLAYER CHARACTER
	FNPCs []Holder // Friendly NPC
	ENPCs []Holder // Enemy NPC
	SOLs  []Holder // Some Generic Solid Object
	TRIGs []Holder // GENERIC TRIGGER
}

func (cm *CullingMap) Add(obj Holder, s *Scene) bool {

	obj_pos := WorldToScreen(obj.GetPos(), s.Cam)

	// If not within viewport boundries...
	if !(obj_pos.X > 0 &&
		obj_pos.Y > 0 &&
		obj_pos.X < s.Window.Width &&
		obj_pos.Y < s.Window.Height) {
		return false
	}

	switch obj.GetType() {
	case ENPC:
		cm.ENPCs = append(cm.ENPCs, obj)
		return true
	case FNPC:
		cm.FNPCs = append(cm.FNPCs, obj)
		return true
	case SOL:
		cm.SOLs = append(cm.SOLs, obj)
		return true
	case TRIG:
		cm.TRIGs = append(cm.TRIGs, obj)
		return true
	default:
		return true
	}
}

func (cm *CullingMap) Zero() {
	cm.PCs = []Holder{}
	cm.FNPCs = []Holder{}
	cm.ENPCs = []Holder{}
	cm.SOLs = []Holder{}
	cm.TRIGs = []Holder{}
}

type Scene struct {
	Window Window

	StartTime    time.Time
	FrameCounter uint64

	PC *Actor

	World         [][]Space
	Width, Height int32
	Cam           Camera
	CullM         CullingMap
	Objects       []*Object
	Enemies       []*Actor
	GUIBlocks     []*GUIBlock

	TsTxt *sdl.Texture
	SsTxt *sdl.Texture
	Font  *ttf.Font

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
			sdl.Rect{400, 120, tileSize, tileSize},
			&Gfx{},
			PC,
		},
		Speed:      Vector2d{2, 2},
		Sprite:     sp1,
		CurrAction: sp1.WD,
	}

	s.World, s.Objects = LoadTMX("assets/world.tmx", renderer)

	s.WidthCells = uint16(len(s.World))
	s.HeightCells = uint16(len(s.World[0]))

	// Defining some sample GUI elements
	// Very early work
	cl1 := sdl.Color{20, 20, 20, 255}
	cl2 := sdl.Color{50, 50, 50, 255}
	cl3 := sdl.Color{255, 255, 255, 255}

	txt1 := TextEl{
		RelPos:  Vector2d{20, 20},
		Content: "testing some gui",
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

	//println(&gb1)
	s.GUIBlocks = append(s.GUIBlocks, &gb1)

	for _, block := range s.GUIBlocks {
		block.Baked = block.Bake(renderer)
	}

	// TODO: load from tiled objects
	for i := 0; uint16(i) < s.EnemyCount; i++ {

		ci := rand.Int31n(int32(s.WidthCells) * tileSize)
		cj := rand.Int31n(int32(s.HeightCells) * tileSize)

		s.Enemies = append(s.Enemies, &Actor{
			Object: Object{
				sdl.Rect{ci, cj, tileSize, tileSize},
				&Gfx{},
				ENPC,
			},
			Speed:      Vector2d{1, 1},
			Sprite:     sp1,
			CurrAction: sp1.WD,
		})
	}
}

func (s *Scene) Update() {

	// Update AI for Enemies that are within Screen
	for _, e := range s.Enemies {
		pos := WorldToScreen(e.Pos, s.Cam)

		if pos.X > 0 && pos.Y > 0 && pos.X < s.Window.Width && pos.Y < s.Window.Height {
			e.UpdateAI(s)
		}
	}
}

func (s *Scene) Render(renderer *sdl.Renderer) {

	// Empty CullM
	s.CullM.Zero()

	var init int32 = 0
	var Source *sdl.Rect

	var ofX, ofY int32 = tileSize, tileSize

	renderer.SetDrawColor(0, 0, 0, 255)

	// Rendering the map
	for sh := init; sh < s.Window.Height; sh += ofY {

		for sw := init; sw < s.Window.Width; sw += ofX {

			ofX = (tileSize - ((s.Cam.WX + sw) % tileSize))
			ofY = (tileSize - ((s.Cam.WY + sh) % tileSize))

			var worldCellX uint16 = uint16((s.Cam.WX + sw) / tileSize)
			var worldCellY uint16 = uint16((s.Cam.WY + sh) / tileSize)

			// Draw black box for out of bounds areas
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

			// Updating CullM with SOLID/COLLIDABLE terrain types
			/*if s.World[worldCellX][worldCellY].Coll {
				s.CullM.Add(s.World[worldCellX][worldCellY], s)
			}*/
		}
	}

	// Rendering the enemies
	for _, e := range s.Enemies {
		pos := WorldToScreen(e.Pos, s.Cam)

		if in := s.CullM.Add(e, s); in {
			renderer.Copy(s.SsTxt, e.GetPose(), &pos)
		}
	}

	// Rendering the player character
	pos := WorldToScreen(s.PC.Pos, s.Cam)
	renderer.Copy(s.SsTxt, s.PC.GetPose(), &pos)

	// Rendering FRAME RATE COUNTER
	fps := fmt.Sprintf("%v", s.GetFPS())
	surface, _ := s.Font.RenderUTF8_Solid(fps, sdl.Color{255, 255, 255, 255})
	defer surface.Free()

	txtr, _ := renderer.CreateTextureFromSurface(surface)
	defer txtr.Destroy()
	renderer.Copy(txtr, &sdl.Rect{0, 0, surface.W, surface.H}, &sdl.Rect{0, 0, surface.W, surface.H})

	// Rendering Game Objects
	for _, obj := range s.Objects {

		if in := s.CullM.Add(obj, s); in {
			obj_pos := WorldToScreen(obj.Pos, s.Cam)
			renderer.SetDrawColor(255, 0, 0, 125)
			renderer.FillRect(&obj_pos)
		}
	}

	// Rendering GUI Blocks
	for _, gb := range s.GUIBlocks {
		renderer.Copy(gb.Baked, gb.Pos, gb.Pos)
	}

	renderer.Present()
}

func (s *Scene) GetFPS() float64 {
	x := time.Since(s.StartTime)
	return (float64(s.FrameCounter) / x.Seconds())
}
