package main

import (
	"github.com/fabtjar/gopack/pkg/pack"
)

func main() {
	images := pack.GetImagesFromDir("images")
	images.Pack()
	images.CreateSheets()
	images.CreateAtlas()
}
