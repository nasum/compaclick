package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createDummyImage is a helper to generate dummy image files for testing
func createDummyImage(t *testing.T, path string, format string) {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create image file: %v", err)
	}
	defer func() { _ = f.Close() }()

	if format == "jpeg" {
		if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
			t.Fatalf("failed to encode jpeg: %v", err)
		}
	} else if format == "png" {
		if err := png.Encode(f, img); err != nil {
			t.Fatalf("failed to encode png: %v", err)
		}
	} else {
		t.Fatalf("unsupported dummy format: %s", format)
	}
}

func TestCompressImage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "compaclick_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	tests := []struct {
		name    string
		format  string
		ext     string
		quality int
		wantErr bool
	}{
		{
			name:    "compress jpeg",
			format:  "jpeg",
			ext:     ".jpg",
			quality: 50,
			wantErr: false,
		},
		{
			name:    "compress png",
			format:  "png",
			ext:     ".png",
			quality: 80,
			wantErr: false,
		},
		{
			name:    "unsupported format",
			format:  "txt",
			ext:     ".txt",
			quality: 50,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputPath := filepath.Join(tempDir, "test_image"+tt.ext)

			if tt.format == "txt" {
				if err := os.WriteFile(inputPath, []byte("not an image"), 0644); err != nil {
					t.Fatalf("failed to write dummy txt: %v", err)
				}
			} else {
				createDummyImage(t, inputPath, tt.format)
			}

			err := compressImage(inputPath, tt.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify that output file is correctly created
				outputPath := filepath.Join(tempDir, "test_image_compressed.jpg")
				info, err := os.Stat(outputPath)
				if err != nil {
					t.Errorf("expected output file does not exist: %v", err)
					return
				}
				if info.Size() == 0 {
					t.Errorf("output file is empty")
				}

				// Verify it is a valid jpeg image
				f, err := os.Open(outputPath)
				if err != nil {
					t.Errorf("failed to open output file: %v", err)
				}
				defer func() { _ = f.Close() }()

				_, format, err := image.Decode(f)
				if err != nil {
					t.Errorf("failed to decode output file: %v", err)
				}
				if format != "jpeg" {
					t.Errorf("output file is not a jpeg, got format: %s", format)
				}
			}
		})
	}
}
