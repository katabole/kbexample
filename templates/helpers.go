package templates

import (
	"context"
)

type contextKey int8

const (
	flashKey contextKey = 0
	assetKey contextKey = 1
)

// SetGlobals is a helper function to set global variables on each page render.
// Under the hood this sets the given arguments on the templ context object that can be used while
// rendering templates.
//
// If the list of global variables needed gets too long, this function could always be refactored
// into smaller set functions
func SetGlobals(ctx context.Context, manifest map[string]string) context.Context {
	ctx = context.WithValue(ctx, assetKey, manifest)
	return ctx
}

func getAsset(ctx context.Context, asset string) string {
	assets, ok := ctx.Value(assetKey).(map[string]string)

	if !ok {
		return ""
	}
	return "/" + assets[asset]
}
