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

func CheckAllColl(sbj Holder, s *Scene, spd Vector2d) bool {

	a_pos := sbj.GetPos()

	// NOTE: ORDER IS IMPORTANT HERE
	// Since the func return (exits) on the first collision

	sbj_type := sbj.GetType()

	// If the subject is a Player Character...
	if sbj_type == PC {
		println("pc")

		// Check against Abstract objects e.g: door triggers
		for _, b := range s.CullM.TRIGs {

			b_pos := WorldToScreen(b.GetPos(), s.Cam)

			if CheckColl(a_pos, b_pos, spd) {
				println("collided!!")
				// b.Collided(sbj)
				return true
			}
		}
	}

	for _, b := range s.CullM.SOLs {

		b_pos := WorldToScreen(b.GetPos(), s.Cam)

		if CheckColl(a_pos, b_pos, spd) {
			return true
		}
	}

	for _, b := range s.CullM.ENPCs {

		b_pos := WorldToScreen(b.GetPos(), s.Cam)

		if CheckColl(a_pos, b_pos, spd) {
			return true
		}
	}

	for _, b := range s.CullM.FNPCs {

		b_pos := WorldToScreen(b.GetPos(), s.Cam)

		if CheckColl(a_pos, b_pos, spd) {
			return true
		}
	}

	return false
}

func v2Rect(p Vector2d) sdl.Rect {
	return sdl.Rect{p.X, p.Y, tileSize, tileSize}
}

func WorldToScreen(pos sdl.Rect, cam Camera) sdl.Rect {
	return sdl.Rect{
		X: pos.X - cam.WX,
		Y: pos.Y - cam.WY,
		W: pos.W,
		H: pos.H,
	}
}

func Distance(p1, p2 sdl.Rect) uint32 {
	return uint32(math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2)))
}
