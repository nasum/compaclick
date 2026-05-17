package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	quality := flag.Int("q", 80, "JPEG quality (1-100)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return
	}

	command := args[0]
	if command == "install" {
		if err := install(); err != nil {
			fmt.Println("Install error:", err)
		} else {
			fmt.Println("Installed successfully.")
		}
		return
	} else if command == "uninstall" {
		if err := uninstall(); err != nil {
			fmt.Println("Uninstall error:", err)
		} else {
			fmt.Println("Uninstalled successfully.")
		}
		return
	}

	// Not install/uninstall, so it should be a file path.
	filePath := args[0]
	if err := compressImage(filePath, *quality); err != nil {
		fmt.Println("Error compressing image:", err)
	}
}

func compressImage(inputPath string, quality int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	// We only process jpeg and png
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported format: %s", format)
	}

	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	dir := filepath.Dir(inputPath)

	// Output is always jpg
	outputPath := filepath.Join(dir, fmt.Sprintf("%s_compressed.jpg", baseName))

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()

	options := &jpeg.Options{Quality: quality}
	return jpeg.Encode(outFile, img, options)
}

func install() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath = filepath.Clean(exePath)

	baseKeyPath := `SystemFileAssociations\image\shell\CompaClick`

	// Create base key
	key, _, err := registry.CreateKey(registry.CLASSES_ROOT, baseKeyPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer func() { _ = key.Close() }()

	err = key.SetStringValue("MUIVerb", "Compress image")
	if err != nil {
		return err
	}

	err = key.SetStringValue("SubCommands", "")
	if err != nil {
		return err
	}

	// Add subcommands
	qualities := []struct {
		name string
		q    int
	}{
		{"cmd1", 80},
		{"cmd2", 50},
		{"cmd3", 30},
	}

	for _, q := range qualities {
		cmdKeyPath := fmt.Sprintf(`%s\shell\%s`, baseKeyPath, q.name)
		cmdKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, cmdKeyPath, registry.ALL_ACCESS)
		if err != nil {
			return err
		}

		err = cmdKey.SetStringValue("MUIVerb", fmt.Sprintf("Compress to %d%%", q.q))
		_ = cmdKey.Close()
		if err != nil {
			return err
		}

		execKeyPath := fmt.Sprintf(`%s\command`, cmdKeyPath)
		execKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, execKeyPath, registry.ALL_ACCESS)
		if err != nil {
			return err
		}

		cmdString := fmt.Sprintf(`"%s" -q %d "%%1"`, exePath, q.q)
		err = execKey.SetStringValue("", cmdString)
		_ = execKey.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func uninstall() error {
	baseKeyPath := `SystemFileAssociations\image\shell\CompaClick`

	// registry.DeleteKey doesn't delete recursively, so we must delete children first
	qualities := []string{"cmd1", "cmd2", "cmd3"}

	for _, q := range qualities {
		execKeyPath := fmt.Sprintf(`%s\shell\%s\command`, baseKeyPath, q)
		_ = registry.DeleteKey(registry.CLASSES_ROOT, execKeyPath)

		cmdKeyPath := fmt.Sprintf(`%s\shell\%s`, baseKeyPath, q)
		_ = registry.DeleteKey(registry.CLASSES_ROOT, cmdKeyPath)
	}

	_ = registry.DeleteKey(registry.CLASSES_ROOT, fmt.Sprintf(`%s\shell`, baseKeyPath))
	return registry.DeleteKey(registry.CLASSES_ROOT, baseKeyPath)
}
