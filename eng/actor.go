package eng

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
)

const (
	tileSize int32 = 32
)

type Actor struct {
	Object
	Chasing bool
	Speed   Vector2d
	MvStack Vector2d

	//SpriteSheet and Animation
	Sprite     *SpriteSheet
	CurrAction SpriteAction
	CurrPose   uint8
}

func (a *Actor) Animate(m Vector2d) {
	var nAction SpriteAction

	// FIX-ME
	if m.X < 0 {
		nAction = a.Sprite.WL
	}
	if m.X > 0 {
		nAction = a.Sprite.WR
	}
	if m.Y < 0 {
		nAction = a.Sprite.WU
	}
	if m.Y > 0 {
		nAction = a.Sprite.WD
	}

	if nAction != a.CurrAction {
		a.CurrPose = 0
		a.CurrAction = nAction
	} else {
		if a.CurrPose/24 < (a.CurrAction.NPoses - 1) {
			a.CurrPose += 4 // use increment to dictate the Speed for frame change
		} else {
			a.CurrPose = 0
		}
	}

}

func (a *Actor) GetPos() sdl.Rect {
	return a.Object.Pos
}

func (a *Actor) GetPose() *sdl.Rect {

	return &sdl.Rect{
		a.CurrAction.Source.X + (int32(a.CurrPose/24) * tileSize),
		a.CurrAction.Source.Y,
		tileSize,
		tileSize,
	}
}

func (a *Actor) UpdateAI(currScene *Scene) {
	var m Vector2d

	if Distance(a.Pos, currScene.PC.Pos) < 160 {
		a.Chasing = true
		a.MvStack = Vector2d{
			(currScene.PC.Pos.X - a.Pos.X) / a.Speed.X,
			(currScene.PC.Pos.Y - a.Pos.Y) / a.Speed.Y,
		}
	} else {
		a.Chasing = false
	}

	if a.MvStack.X > 0 {
		m = MVSpeed.Right
		a.MvStack.X -= 1
	}

	if a.MvStack.X < 0 {
		m = MVSpeed.Left
		a.MvStack.X += 1
	}

	if a.MvStack.Y > 0 {
		m = MVSpeed.Down
		a.MvStack.Y -= 1
	}

	if a.MvStack.Y < 0 {
		m = MVSpeed.Up
		a.MvStack.Y += 1
	}

	a.Move(m, currScene)

	// logic for npc wandering
	if !a.Chasing && a.MvStack.X == 0 && a.MvStack.Y == 0 {
		a.MvStack = Vector2d{
			rand.Int31n(3*tileSize) * int32(math.Copysign(1.0, float64(rand.Int31n(1))-1.0)),
			rand.Int31n(3*tileSize) * int32(math.Copysign(1.0, float64(rand.Int31n(1))-1.0)),
		}
	}
}

func (a *Actor) Move(m Vector2d, currScene *Scene) {

	nPos := a.Pos

	if a.Sprite != nil {
		a.Animate(m)
	}

	fSpd := Vector2d{(m.X * a.Speed.X), (m.Y * a.Speed.Y)}
	nPos.X += fSpd.X
	nPos.Y += fSpd.Y

	if !CheckAllColl(WorldToScreen(nPos, currScene), currScene, fSpd) {
		a.Pos = nPos
	}
}
