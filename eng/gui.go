package eng

/*

Scene
	GUIBlocks []GUIBlock

*/

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type GUIBlock struct {
	Elements []GUIEl
	Baked    *sdl.Texture
	Pos      *sdl.Rect
}

type GUIEl struct {
	Layer   int
	Pos     *sdl.Rect
	BGColor sdl.Color
	Texts   []TextEl
}

type TextEl struct {
	RelPos  Vector2d
	Font    *ttf.Font
	Content string
	Color   sdl.Color
}

func (gb *GUIBlock) Update() {
	for _, el := range gb.Elements {
		println("Updating: ")
		println(el.Layer, el.Pos.X, el.Pos.Y)
	}
}

func (gb *GUIBlock) Bake(renderer *sdl.Renderer) *sdl.Texture {

	finalTxtr, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, 800, 600)
	originalTarget := renderer.GetRenderTarget()
	renderer.SetRenderTarget(finalTxtr)
	defer renderer.SetRenderTarget(originalTarget)

	renderer.SetDrawColor(1, 1, 1, 0)
	renderer.FillRect(gb.Pos)

	for _, el := range gb.Elements {
		println("Baking: ")
		println(el.Layer, el.Pos.X, el.Pos.Y)

		renderer.SetDrawColor(el.BGColor.R, el.BGColor.G, el.BGColor.B, el.BGColor.A)
		renderer.FillRect(el.Pos)

		for _, txt := range el.Texts {
			texture, W, H := txt.Bake(renderer)
			renderer.Copy(
				texture,
				&sdl.Rect{0, 0, W, H},
				&sdl.Rect{el.Pos.X + txt.RelPos.X, el.Pos.Y + txt.RelPos.Y, W, H})
		}
	}

	finalTxtr.SetBlendMode(sdl.BLENDMODE_BLEND)
	finalTxtr.SetAlphaMod(216)

	return finalTxtr
}

func (t *TextEl) Bake(renderer *sdl.Renderer) (*sdl.Texture, int32, int32) {
	surface, _ := t.Font.RenderUTF8_Solid(t.Content, t.Color)
	defer surface.Free()

	txtr, _ := renderer.CreateTextureFromSurface(surface)
	//defer txtr.Destroy()

	return txtr, surface.W, surface.H
}
