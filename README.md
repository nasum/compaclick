# CompaClick

CompaClick is a lightweight Windows context-menu utility that allows you to easily compress JPEG and PNG images directly from the Explorer right-click menu.

## Features

- **Quick compression** via Windows right-click context menu.
- **Three compression quality presets**: 80%, 50%, and 30%.
- Automatically outputs the compressed image in the same directory (appends `_compressed.jpg` to the filename).
- Supports both `.jpg` and `.png` input files.

## Installation

1. Download the installer (`CompaClick_Setup.exe`) from the Releases page.
2. Run the installer. It will automatically add CompaClick to your right-click context menu for image files.

> **Note for Windows 11 Users:**
> Due to Windows 11's context menu design, CompaClick will appear under **"Show more options"** (その他のオプションを表示). Alternatively, you can access it directly by holding down the `Shift` key while right-clicking a file.

## Uninstallation

You can uninstall CompaClick easily through the Windows Settings > "Apps & features" (アプリと機能). This will also cleanly remove the registry keys for the context menu.

## Building from Source

### Prerequisites
- Go 1.26.0 or higher
- Inno Setup 6 (for building the installer)

### Steps
1. Clone this repository.
2. Run the build script:
   ```cmd
   build.bat
   ```
3. The Go binary will be compiled, and `CompaClick_Setup.exe` will be generated in the `output/` folder.

## Linting and Testing

This project uses `golangci-lint` and standard Go tests.
To run the tests and linter:

```cmd
go test -v ./...
golangci-lint run
```
