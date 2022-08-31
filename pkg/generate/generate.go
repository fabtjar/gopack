package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func clearImagesDir() error {
	err := os.RemoveAll("images")
	if err != nil {
		return fmt.Errorf("faied to clear images directory: %v", err)
	}
	err = os.MkdirAll("images", os.ModePerm)
	if err != nil {
		return fmt.Errorf("faied to create images directory: %v", err)
	}
	return nil
}

func downloadImage(filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("faied to get HTTP response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error with response: %v", resp.Status)
	}

	f, err := os.Create(fmt.Sprintf("images/%s", filename))
	if err != nil {
		return fmt.Errorf("faied to create file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy image to file: %v", err)
	}

	return nil
}

func main() {
	var imageCount int

	if len(os.Args) == 1 {
		log.Fatalf("Missing image_count arg")
	}

	imageCount, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse image_count arg")
	}

	err = clearImagesDir()
	if err != nil {
		log.Fatalf("Failed to clear images dir: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(imageCount)
	for i := 0; i < imageCount; i++ {
		go func(n int) {
			filename := fmt.Sprintf("image_%03d.png", n)
			width := (rand.Intn(4) + 1) * 20
			height := (rand.Intn(4) + 1) * 20
			color := rand.Intn(0xffffff)
			url := fmt.Sprintf("https://dummyimage.com/%dx%d/%06x/000000.png", width, height, color)
			err = downloadImage(filename, url)
			if err != nil {
				log.Fatalf("Failed to download image #%d: %v", n, err)
			}
			wg.Done()
		}(i + 1)
	}
	wg.Wait()
}
