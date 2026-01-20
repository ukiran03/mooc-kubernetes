package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Image struct {
	mu      sync.Mutex
	name    string
	modTime time.Time
}

func (i *Image) String() string {
	return fmt.Sprintf("Name:   %s\n Last Modified: %s", i.name, i.modTime)
}

func GetImage(img *Image) (bool, *Image) {
	img.mu.Lock()
	defer img.mu.Unlock()

	if !img.modTime.IsZero() && time.Since(img.modTime) < 10*time.Minute {
		fileName, err := downloadImage()
		if err != nil {
			log.Printf("Download failed: %v", err)
			// Treat as Cache hit (fallback to old data)
			return true, img
		}

		info, err := os.Stat(fileName)
		if err != nil {
			log.Printf("Stat failed: %v", err)
			// Treat as Cache hit (fallback to old data)
			return true, img
		}
		img.name = fileName
		img.modTime = info.ModTime()

	}
	return false, img
}

func downloadImage() (string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	f, err := os.Create(pathname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return pathname, nil
}
