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

const (
	maxWidth  = 500
	maxHeight = maxWidth
)

type Space struct {
	Free  bool
	Sheet int
	Rectangle
}

func getUsedSpace(spaces []Space, r Rectangle) *Space {
	for i := 0; i < len(spaces); i++ {
		s := &spaces[i]
		if !s.Free {
			s.Free = true
			s.Rectangle = r
			return s
		}
	}
	return nil
}

type Image struct {
	Name  string
	Rect  Rectangle
	Sheet int
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

	sheetCount := 1
	spaces := []Space{{Free: true, Sheet: 0, Rectangle: Rectangle{Width: maxWidth, Height: maxHeight}}}

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
			space = &Space{Free: true, Sheet: sheetCount, Rectangle: Rectangle{Width: maxWidth, Height: maxHeight}}
			sheetCount++
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
			iSpc, jSpc := spaces[i], spaces[j]
			return iSpc.X*iSpc.Y+iSpc.Width*iSpc.Height < jSpc.X*jSpc.Y+jSpc.Width*jSpc.Height
		})
	}

	width := maxWidth
	height := maxHeight
	for i := 0; i < sheetCount; i++ {
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
			draw.Draw(sheet, sheet.Bounds(), img.Image, pos, draw.Src)
		}

		f, err := os.Create(fmt.Sprintf("sheet_%03d.png", i+1))
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
