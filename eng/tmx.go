package eng

import (
	"encoding/xml"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"io/ioutil"
	"os"
)

type Tile struct {
	Gid int `xml:"gid,attr"`
}

type Layer struct {
	Name   string `xml:"name,attr"`
	Tiles  []Tile `xml:"data>tile"`
	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
}

type Map struct {
	Layers      []Layer   `xml:"layer"`
	Tilesets    []TileSet `xml:"tileset"`
	HeightTiles int       `xml:"height,attr"`
	WidthTiles  int       `xml:"width,attr"`
	TileH       int32     `xml:"tileheight,attr"`
	TileW       int32     `xml:"tilewidth,attr"`
}

type TMX struct {
	Layers   []Layer   `xml:"layer"`
	Tilesets []TileSet `xml:"tileset"`

	XMLName     xml.Name `xml:"map"`
	HeightTiles int      `xml:"height,attr"`
	WidthTiles  int      `xml:"width,attr"`
	TileH       int      `xml:"tileheight,attr"`
	TileW       int      `xml:"tilewidth,attr"`
}

type TSImg struct {
	Src       string `xml:"source,attr"`
	SrcHeight int    `xml:"height,attr"`
	SrcWidth  int    `xml:"width,attr"`
}

type TileSet struct {
	Name  string `xml:"name,attr"`
	Image TSImg  `xml:"image"`
	TileH int    `xml:"tileheight,attr"`
	TileW int    `xml:"tilewidth,attr"`

	Txtr *sdl.Texture
	Gids []*Gfx
}

func (ts *TileSet) GetGIDRect(gid int) *sdl.Rect {
	w := ts.Image.SrcWidth / ts.TileW

	var x int32 = int32(((gid - 1) % w) * ts.TileW)
	var y int32 = int32((gid / w) * ts.TileH)

	return &sdl.Rect{x, y, int32(ts.TileW), int32(ts.TileH)}
}

func LoadTMX(mapname string, renderer *sdl.Renderer) [][]Space {

	f, _ := os.Open(mapname)
	output, _ := ioutil.ReadAll(f)
	_ = f.Close()

	tmx := &TMX{}

	err := xml.Unmarshal(output, tmx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the tmx map: %s\n", err)
		os.Exit(11)
	}

	var tilesetTxt *sdl.Texture

	for i := 0; i < len(tmx.Tilesets); i++ {

		tilesetImg, err := img.Load("assets/" + tmx.Tilesets[i].Image.Src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
			os.Exit(3)
		}
		defer tilesetImg.Free()

		tilesetTxt, err = renderer.CreateTextureFromSurface(tilesetImg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			os.Exit(4)
		}
		defer tilesetTxt.Destroy()

		tmx.Tilesets[i].Txtr = tilesetTxt
	}

	world := make([][]Space, 100)
	ts := tmx.Tilesets[0]

	for li, layer := range tmx.Layers {
		for i := 0; i < layer.Height; i++ {

			world[i] = make([]Space, 100)

			for j := 0; j < layer.Width; j++ {
				tile := layer.Tiles[(i*layer.Height)+j]

				world[i][j].Gfxs[li] = &Gfx{
					Txtr:   ts.Txtr,
					Source: ts.GetGIDRect(tile.Gid),
				}
			}
		}
	}

	//tmx = nil
	return world
}

func FreezeLayer(world [][]Space, idx int) *sdl.Texture {
	return &sdl.Texture{}
}
