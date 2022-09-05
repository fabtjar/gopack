package main

import (
	"github.com/fabtjar/gopack/pkg/pack"
)

func main() {
	images := pack.GetImageData("images")
	pack.Pack(images)
	pack.CreateSheets(images)
	pack.CreateAtlas(images)
}
