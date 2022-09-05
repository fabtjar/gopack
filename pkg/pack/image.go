package pack

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

type Image struct {
	Name     string    `json:"name"`
	Location string    `json:"-"`
	Rect     Rectangle `json:"rect"`
	Sheet    int       `json:"sheet"`
}

func GetImageData(dir string) []Image {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Failed to read images dir.")
	}

	var images []Image
	for _, file := range files {
		name := file.Name()
		location := fmt.Sprintf("%s/%s", dir, name)
		img := getImage(location)
		size := img.Bounds().Max
		images = append(images, Image{Name: name, Location: location, Rect: Rectangle{0, 0, size.X, size.Y}})
	}
	return images
}

func getImage(name string) image.Image {
	f, err := os.Open(name)
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
