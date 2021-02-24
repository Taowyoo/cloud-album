package util

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestExampleDecode(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"TestDecodeJPG", "../example_imgs/us-portland-bird.jpg"},
		{"TestDecodeJPG", "../example_imgs/turkey-bodrum.jpg"},
		{"TestDecodeJPG", "../example_imgs/germany-garching-heide.jpg"},
		{"TestDecodeJPG", "../example_imgs/england-portslade-grassfield.jpg"},
		{"TestDecodeJPG", "../example_imgs/england-london-bridge.jpg"},
	}
	for i, tt := range tests {
		t.Run(tt.name+fmt.Sprintf("#%02d", i), func(t *testing.T) {
			fmt.Println("Decoding:", tt.path)
			f, err := os.Open(tt.path)
			if err != nil {
				log.Panic(err)
			}
			res, err := GetExifData(f)
			if err != nil {
				log.Panic(err)
			}
			fmt.Println("Result:\n", res)
		})
	}
}
