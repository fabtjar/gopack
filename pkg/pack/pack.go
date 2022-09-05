package pack

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sort"

	calc "github.com/fabtjar/gopack/pkg/math"
)

const (
	maxWidth  = 500
	maxHeight = maxWidth
)

type Space struct {
	Free  bool
	Sheet int
	Rectangle
}

type Images []Image

func (images Images) CreateAtlas() {
	mar, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}

	f, err := os.Create("atlas.json")
	if err != nil {
		log.Fatalf("Failed to create atlas: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(string(mar))
	if err != nil {
		log.Fatalf("Failed to write to atlas: %v", err)
	}
}

func (images Images) getSheetCount() int {
	var n int
	for _, img := range images {
		if img.Sheet > n {
			n = img.Sheet
		}
	}
	return n
}

func (images Images) CreateSheets() {
	sheetCount := images.getSheetCount()

	width := maxWidth
	height := maxHeight
	for i := 1; i <= sheetCount; i++ {
		sheet := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				sheet.Set(x, y, color.White)
			}
		}

		for _, img := range images {
			if img.Sheet != i {
				continue
			}
			pos := image.Point{-img.Rect.X, -img.Rect.Y}
			file := getImage(img.Location)
			draw.Draw(sheet, sheet.Bounds(), file, pos, draw.Src)
		}

		f, err := os.Create(fmt.Sprintf("sheet_%03d.png", i))
		if err != nil {
			log.Fatalf("Failed to create sheet file: %v", err)
		}
		defer f.Close()
		err = png.Encode(f, sheet)
		if err != nil {
			log.Fatalf("Failed to write to sheet file: %v", err)
		}
	}
}

func (images Images) Pack() {
	sort.Slice(images, func(i, j int) bool {
		a, b := images[i].Rect, images[j].Rect
		return calc.Max(a.Width, a.Height) > calc.Max(b.Width, b.Height)
	})

	sheetNumber := 1
	spaces := []Space{{Free: true, Sheet: sheetNumber, Rectangle: Rectangle{Width: maxWidth, Height: maxHeight}}}

	for i := range images {
		image := &images[i]
		var space *Space
		for _, s := range spaces {
			if s.Free && s.Rectangle.canFit(image.Rect) {
				space = &s
				break
			}
		}
		if space == nil {
			space = &Space{Free: true, Sheet: sheetNumber, Rectangle: Rectangle{Width: maxWidth, Height: maxHeight}}
			sheetNumber++
			spaces = append(spaces, *space)
		}

		image.Rect.X = space.X
		image.Rect.Y = space.Y
		image.Sheet = space.Sheet

		for j := range spaces {
			overlap := &spaces[j]
			if overlap.Free && image.Rect.overlaps(overlap.Rectangle) {
				overlap.Free = false
				for _, r := range overlap.cut(image.Rect) {
					spaces = append(spaces, Space{Free: true, Sheet: overlap.Sheet, Rectangle: r})
				}
			}
		}

		sort.Slice(spaces, func(i, j int) bool {
			a, b := spaces[i], spaces[j]
			if a.Sheet != b.Sheet {
				return a.Sheet > b.Sheet
			}
			return a.X*a.X+a.Y*a.Y < b.X*b.X+b.Y*b.Y
		})
	}
}
