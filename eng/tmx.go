package eng

import (
	"encoding/xml"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type ObjGroup struct {
	Name    string    `xml:"name,attr"`
	Objects []TObject `xml:"object"`
}

type TObject struct {
	Name     string     `xml:"name,attr"`
	Type     string     `xml:"type,attr"`
	X        int32      `xml:"x,attr"`
	Y        int32      `xml:"y,attr"`
	W        int32      `xml:"width,attr"`
	H        int32      `xml:"height,attr"`
	PropList []Property `xml:"properties>property"`
}

type TileType struct {
	Id     int    `xml:"id,attr"`
	TerStr string `xml:"terrain,attr"`
	TerArr [4]int
}

type Tile struct {
	Gid int `xml:"gid,attr"`
}

type Layer struct {
	Name   string `xml:"name,attr"`
	Tiles  []Tile `xml:"data>tile"`
	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
}

type TMX struct {
	Layers   []Layer   `xml:"layer"`
	Tilesets []TileSet `xml:"tileset"`

	XMLName     xml.Name   `xml:"map"`
	HeightTiles int        `xml:"height,attr"`
	WidthTiles  int        `xml:"width,attr"`
	TileH       int        `xml:"tileheight,attr"`
	TileW       int        `xml:"tilewidth,attr"`
	ObjGroups   []ObjGroup `xml:"objectgroup"`
}

type TSImg struct {
	Src       string `xml:"source,attr"`
	SrcHeight int    `xml:"height,attr"`
	SrcWidth  int    `xml:"width,attr"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Terrain struct {
	Name     string     `xml:"name,attr"`
	PropList []Property `xml:"properties>property"`
}

type TileSet struct {
	Name  string `xml:"name,attr"`
	Image TSImg  `xml:"image"`
	TileH int    `xml:"tileheight,attr"`
	TileW int    `xml:"tilewidth,attr"`

	TerrTypes    []Terrain  `xml:"terraintypes>terrain"`
	TerrainTiles []TileType `xml:"tile"`
	TTMap        map[int][4]*Terrain

	Txtr *sdl.Texture
	Gids []*Gfx
}

func (ts *TileSet) GetGIDRect(gid int) *sdl.Rect {
	w := ts.Image.SrcWidth / ts.TileW

	var x int32 = int32(((gid - 1) % w) * ts.TileW)
	var y int32 = int32((gid / w) * ts.TileH)

	return &sdl.Rect{x, y, int32(ts.TileW), int32(ts.TileH)}
}

func LoadTMX(mapname string, renderer *sdl.Renderer) ([][]Space, []*Object) {

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

		tmx.Tilesets[i].TTMap = make(map[int][4]*Terrain)

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

		for _, tt := range tmx.Tilesets[i].TerrainTiles {

			var terrList [4]*Terrain

			// Spliting and converting into integer so that it
			// can be used as array idx
			for ti, terr := range strings.Split(tt.TerStr, ",") {
				ttype, _ := strconv.Atoi(terr)
				terrList[ti] = &tmx.Tilesets[i].TerrTypes[ttype]
			}
			tmx.Tilesets[i].TTMap[tt.Id] = terrList
		}

		tmx.Tilesets[i].Txtr = tilesetTxt
	}

	world := make([][]Space, tmx.HeightTiles)
	ts := tmx.Tilesets[0]

	for li, layer := range tmx.Layers {

		for i := 0; i < layer.Height; i++ {

			world[i] = make([]Space, tmx.WidthTiles)

			for j := 0; j < layer.Width; j++ {
				tile := layer.Tiles[(i*layer.Height)+j]

				world[i][j].Terrains = ts.TTMap[tile.Gid]
				world[i][j].Gfxs[li] = &Gfx{
					Txtr:   ts.Txtr,
					Source: ts.GetGIDRect(tile.Gid),
				}

				for _, terr := range world[i][j].Terrains {
					if terr == nil {
						continue
					}

					for _, prop := range terr.PropList {
						switch prop.Name {
						case "COLL":
							world[i][j].Coll = (prop.Value == "1")
							break
						case "WARP":
							world[i][j].Warp = prop.Value
							break
						}
					}
				}
			}
		}
	}

	// Now loading the Tiled Objects
	var objects []*Object

	for _, g := range tmx.ObjGroups {
		for _, o := range g.Objects {
			obj := &Object{
				Pos:  sdl.Rect{o.X, o.Y, o.W, o.H},
				Gfx:  nil,
				Type: TRIG,
			}
			objects = append(objects, obj)
		}
	}

	//tmx = nil
	return world, objects
}

func FreezeLayer(world [][]Space, idx int) *sdl.Texture {
	return &sdl.Texture{}
}
