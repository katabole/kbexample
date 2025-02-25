package build

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed dist/manifest.json
var Manifest []byte

//go:embed dist
var distDir embed.FS

// TODO update
// DistDir returns an FS representing the `dist` directory, which has the final assets (js, css, images, etc.) compiled
// and put there by the build pipeline.
func DistDir() fs.FS {
	f, err := fs.Sub(distDir, "dist")
	if err != nil {
		panic(fmt.Sprintf("error opening the assets directory in the build directory (hint: try running vite): %w", err))
	}
	return f
}

func AssetsDir() fs.FS {
	f, err := fs.Sub(DistDir(), "assets")
	if err != nil {
		panic(fmt.Sprintf("error opening the assets directory in the build directory (hint: try running vite): %w", err))
	}
	return f
}
