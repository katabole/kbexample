package dist

import "embed"

//go:embed assets/*
var BuiltAssets embed.FS

//go:embed manifest.json
var Manifest []byte
