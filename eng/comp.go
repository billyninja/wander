package eng

import (
	"github.com/veandco/go-sdl2/sdl"
	//"github.com/veandco/go-sdl2/sdl_image"
)

const (
	tileSize int32 = 32
)

// TEMP
var gfGrass = &Gfx{
	Source: &sdl.Rect{0, 0, tileSize, tileSize},
}

var gfTree = &Gfx{
	Source: &sdl.Rect{0, 32, tileSize, tileSize},
	Coll:   true,
}

var gfWall = &Gfx{
	Source: &sdl.Rect{32, 0, tileSize, tileSize},
	Coll:   true,
}

var gfEnemy = &Gfx{
	Source: &sdl.Rect{576, 512, tileSize, tileSize},
	Coll:   true,
}

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
	Coll   bool
}

type Space struct {
	Gfxs [4]*Gfx
}

type Object struct {
	Pos sdl.Rect
	Gfx *Gfx
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
