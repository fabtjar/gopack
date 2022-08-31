package pack

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Rectangle struct {
	X, Y, Width, Height int
}

func (r Rectangle) isLarger(other Rectangle) bool {
	return r.Width > other.Width || r.Height > other.Height
}

func (r Rectangle) top() int {
	return r.Y
}

func (r Rectangle) bottom() int {
	return r.Y + r.Height
}

func (r Rectangle) left() int {
	return r.X
}

func (r Rectangle) right() int {
	return r.X + r.Width
}

func (r Rectangle) overlaps(other Rectangle) bool {
	return r.left() < other.right() && r.right() > other.left() && r.top() < other.bottom() && r.bottom() > other.top()
}

type Space struct {
	Free bool
	Rectangle
}

func (s Space) cut(cut Rectangle) []Space {
	var spaces []Space

	top := Rectangle{s.left(), s.top(), s.Width, cut.top() - s.top()}
	bottom := Rectangle{s.left(), cut.bottom(), s.Width, s.bottom() - cut.bottom()}
	left := Rectangle{s.left(), s.top(), cut.left() - s.left(), s.Height}
	right := Rectangle{cut.right(), s.top(), s.right() - cut.right(), s.Height}

	for _, r := range []Rectangle{top, bottom, left, right} {
		if r.Width > 0 && r.Height > 0 {
			spaces = append(spaces, Space{Free: true, Rectangle: r})
		}
	}

	return spaces
}

const (
	maxWidth  = 500
	maxHeight = maxWidth
)

type Image struct {
	Name string
	Rect Rectangle
	image.Image
}

func getImages(dir string) []Image {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Failed to read images dir.")
	}

	var images []Image
	for _, file := range files {
		img := getImage(fmt.Sprintf("%s/%s", dir, file.Name()))
		size := img.Bounds().Max
		images = append(images, Image{Name: file.Name(), Rect: Rectangle{0, 0, size.X, size.Y}, Image: img})
	}
	return images
}

func getImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to read file.")
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal("Failed to decode image.")
	}
	return img
}

func Pack(dir string) {
	images := getImages(dir)

	sort.Slice(images, func(i, j int) bool {
		iImg, jImg := images[i], images[j]
		return iImg.Rect.Width*iImg.Rect.Height > jImg.Rect.Width*jImg.Rect.Height
	})

	spaces := []Space{{Free: true, Rectangle: Rectangle{Width: maxWidth, Height: maxHeight}}}

	for i := range images {
		image := &images[i].Rect
		spaceFound := false
		for j := range spaces {
			space := &spaces[j]
			if space.Free && !image.isLarger(space.Rectangle) {
				space.Free = false
				spaceFound = true
				image.X = space.X
				image.Y = space.Y
				spaces = append(spaces, space.cut(*image)...)

				for k := range spaces {
					overlap := &spaces[k]
					if overlap.Free && image.overlaps(overlap.Rectangle) {
						overlap.Free = false
						spaces = append(spaces, overlap.cut(*image)...)
					}
				}

				break
			}
		}
		if !spaceFound {
			log.Fatalf("Failed to find space for image %v (#%d)\n", image, i)
		}
		sort.Slice(spaces, func(i, j int) bool {
			iSpc, jSpc := spaces[i], spaces[j]
			return iSpc.X*iSpc.Y+iSpc.Width*iSpc.Height < jSpc.X*jSpc.Y+jSpc.Width*jSpc.Height
		})
	}

	width := maxWidth
	height := maxHeight
	sheet := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			sheet.Set(x, y, color.White)
		}
	}

	for _, img := range images {
		pos := image.Point{-img.Rect.X, -img.Rect.Y}
		draw.Draw(sheet, sheet.Bounds(), img.Image, pos, draw.Src)
	}

	f, err := os.Create("sheet_001.png")
	if err != nil {
		log.Fatalf("Failed to create sheet file: %v", err)
	}
	defer f.Close()
	err = png.Encode(f, sheet)
	if err != nil {
		log.Fatalf("Failed to write to sheet file: %v", err)
	}
}
