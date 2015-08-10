package eng

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

func CheckColl(a, b sdl.Rect, s Vector2d) bool {

	if a.X-s.X == b.X && a.Y-s.Y == b.Y {
		return false
	}

	//The sides of the rectangles
	var leftA, leftB, rightA, rightB, topA, topB, bottomA, bottomB int32

	//Calculate the sides of rect A
	leftA = a.X
	rightA = a.X + a.W
	topA = a.Y
	bottomA = a.Y + a.H

	//Calculate the sides of rect B
	leftB = b.X
	rightB = b.X + b.W
	topB = b.Y
	bottomB = b.Y + b.H

	//If any of the sides from A are outside of B
	if bottomA <= topB {
		return false
	}

	if topA >= bottomB {
		return false
	}

	if rightA <= leftB {
		return false
	}

	if leftA >= rightB {
		return false
	}

	return true
}

func CheckAllColl(sbj sdl.Rect, s *Scene, spd Vector2d) bool {
	for _, b := range s.CullM {

		pos := b.GetPos()

		// TODO - pensar em um refactory que elimine esse teste
		switch b.(type) {
		case *Actor:
			pos = WorldToScreen(pos, s)
			break
		}

		if CheckColl(sbj, pos, spd) {
			return true
		}
	}
	return false
}

func v2Rect(p Vector2d) sdl.Rect {
	return sdl.Rect{p.X, p.Y, tileSize, tileSize}
}

func WorldToScreen(pos sdl.Rect, currScene *Scene) sdl.Rect {
	return sdl.Rect{
		X: pos.X - currScene.Cam.WX,
		Y: pos.Y - currScene.Cam.WY,
		W: pos.W,
		H: pos.H,
	}
}

func Distance(p1, p2 sdl.Rect) uint32 {
	return uint32(math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2)))
}
