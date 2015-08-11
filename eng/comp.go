package eng

import (
	"github.com/veandco/go-sdl2/sdl"
	//"github.com/veandco/go-sdl2/sdl_image"
)

type ObjectType uint8

const (
	UND  ObjectType = iota // UNDEFINED
	PC                     // PLAYER CHARACTER
	FNPC                   // Friendly NPC
	ENPC                   // Enemy NPC
	SOL                    // Some Generic Solid Object
	TRIG                   // GENERIC TRIGGER
)

var MVSpeed = Movement{
	Vector2d{0, -1},
	Vector2d{0, 1},
	Vector2d{-1, 0},
	Vector2d{1, 0},
}

type Vector2d struct {
	X, Y int32
}

type Gfx struct {
	Source *sdl.Rect
	Txtr   *sdl.Texture
}

type Space struct {
	Gfxs     [4]*Gfx // Store up to 4 layers of tiles
	Terrains [4]*Terrain
	Coll     bool
	Warp     string
	Dmg      int16
	Type     ObjectType
}

func (sp *Space) GetType() ObjectType {
	return sp.Type
}

func (sp *Space) GetPos() sdl.Rect {

	// TODO: Dunno what to do about this =(
	return sdl.Rect{0, 0, 0, 0}
}

type Object struct {
	Pos  sdl.Rect
	Gfx  *Gfx
	Type ObjectType
}

func (o *Object) GetType() ObjectType {
	return o.Type
}

func (o *Object) GetPos() sdl.Rect {
	return o.Pos
}

type SpriteAction struct {
	Source sdl.Rect
	NPoses uint8
}

type Actions struct {
	WR SpriteAction
	WL SpriteAction
	WU SpriteAction
	WD SpriteAction
}

type SpriteSheet struct {
	Actions
	txt *sdl.Texture
}

type Movement struct {
	Up    Vector2d
	Down  Vector2d
	Left  Vector2d
	Right Vector2d
}
