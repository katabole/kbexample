package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/dankinder/katabole/kbexample/public/dist"
	"github.com/dankinder/katabole/kbexample/templates"
	"github.com/unrolled/render"
)

// Renderer wraps the unrolled/render package in order to provide a few goodies (error rendering, session saving).
type Renderer struct {
	rnd           *render.Render
	isProduction  bool
	assetManifest map[string]string
}

func NewRenderer(isProduction bool) (*Renderer, error) {
	r := &Renderer{
		isProduction: isProduction,
	}

	r.rnd = render.New(render.Options{
		IsDevelopment:   !isProduction,
		Layout:          "layout",
		RequirePartials: true,
		// This is "." and not the default "templates" because we're embedding from inside the template directory
		Directory: ".",
		FileSystem: &render.EmbedFileSystem{
			FS: templates.EmbeddedTemplates,
		},
		Funcs: []template.FuncMap{
			template.FuncMap{
				"asset": r.Asset,
			},
		},
	})

	// In production, preload the asset manifest. In dev/test it'll get loaded every time since it may change.
	if isProduction {
		var err error
		r.assetManifest, err = loadManifest()
		if err != nil {
			return nil, fmt.Errorf("could not load asset manifest: %w", err)
		}
	}

	return r, nil
}

// Data writes out the raw bytes as binary data.
func (r *Renderer) Data(w http.ResponseWriter, req *http.Request, status int, v []byte) error {
	SaveSession(w, req)
	return r.rnd.Data(w, status, v)
}

// HTML builds up the response from the specified template and bindings.
func (r *Renderer) HTML(w http.ResponseWriter, req *http.Request, status int, name string, binding any, htmlOpt ...render.HTMLOptions) error {
	flash := Flash(req)
	SaveSession(w, req)
	if binding == nil {
		binding = map[string]any{}
	}
	return r.rnd.HTML(w, status, name, map[string]any{
		"Flash": flash,
		"Data":  binding,
	}, htmlOpt...)
}

// JSON marshals the given interface object and writes the JSON response.
func (r *Renderer) JSON(w http.ResponseWriter, req *http.Request, status int, v interface{}) error {
	SaveSession(w, req)
	return r.rnd.JSON(w, status, v)
}

// JSONP marshals the given interface object and writes the JSON response.
func (r *Renderer) JSONP(w http.ResponseWriter, req *http.Request, status int, callback string, v interface{}) error {
	SaveSession(w, req)
	return r.rnd.JSONP(w, status, callback, v)
}

// Text writes out a string as plain text.
func (r *Renderer) Text(w http.ResponseWriter, req *http.Request, status int, v string) error {
	SaveSession(w, req)
	return r.rnd.Text(w, status, v)
}

// XML marshals the given interface object and writes the XML response.
func (r *Renderer) XML(w http.ResponseWriter, req *http.Request, status int, v interface{}) error {
	SaveSession(w, req)
	return r.rnd.XML(w, status, v)
}

func (r *Renderer) Redirect(w http.ResponseWriter, req *http.Request, url string, status int) {
	SaveSession(w, req)
	http.Redirect(w, req, url, status)
}

// Error figures out how to render the given error message appropriate to the expected content type, while showing full
// error messages in dev but not in production.
func (r *Renderer) Error(w http.ResponseWriter, req *http.Request, status int, err error) {
	switch GetContentType(req) {
	case ContentTypeHTML:
		r.HTMLError(w, req, status, err)
	default:
		r.JSONError(w, req, status, err)
	}
}

// HTMLError sends the user an HTML error page, hiding error details in production.
func (r *Renderer) HTMLError(w http.ResponseWriter, req *http.Request, status int, err error) {
	if r.isProduction {
		slog.Info("Internal error", "err", err)
		r.HTML(w, req, status, "error", map[string]string{"message": "internal error, see logs for details"})
		return
	}

	r.HTML(w, req, status, "error", map[string]string{"message": err.Error()})
}

// HTMLError sends the user a JSON error payload, hiding error details in production.
func (r *Renderer) JSONError(w http.ResponseWriter, req *http.Request, status int, err error) {
	if r.isProduction {
		slog.Info("Internal error", "err", err)
		r.JSON(w, req, status, map[string]string{"message": "internal error, see logs for details"})
		return
	}

	r.JSON(w, req, status, map[string]string{"message": err.Error()})
}

// Asset loading
//

// Asset returns the path to an asset by looking it up in the asset manifest. Webpack processes each asset (to minify
// CSS, for example) and generates the manifest to tell us where to find it. This function is available in templates as
// the "asset" function.
//
// For example consider public/assets/application.js which gets processed by Webpack into
// public/dist/assets/application.[fingerpring].js. The manifest will have something like:
//
//	{
//	  "assets/application.js": "assets/application.39de836e61570e45cf00.js"
//	}
//
// In a template we would put a script tag like:
//
//	<script src="{{ asset "assets/application.js" }}"></script>
//
// Which would render as:
//
//	<script src="/assets/application.39de836e61570e45cf00.js"></script>
//
// This ensures assets are cacheable and have unique URLs.
func (r *Renderer) Asset(assetPath string) (string, error) {
	p, err := r.relativeAssetPath(assetPath)
	return "/" + p, err
}

func (r *Renderer) relativeAssetPath(assetPath string) (string, error) {
	if r.isProduction {
		if len(r.assetManifest) == 0 {
			// We expect app startup to load the manifest file, so if that's not done yet, something is wrong.
			return "", errors.New("no asset manifest loaded")
		}
		if assetPath, ok := r.assetManifest[assetPath]; ok {
			return assetPath, nil
		}
		return "", fmt.Errorf("asset not found: %s", assetPath)
	}

	manifest, err := loadManifest()
	if err != nil {
		return "", fmt.Errorf("failed to load manifest: %w", err)
	}
	if assetPath, ok := manifest[assetPath]; ok {
		return assetPath, nil
	}
	return "", fmt.Errorf("asset not found: %s", assetPath)
}

func loadManifest() (map[string]string, error) {
	m := map[string]string{}
	if err := json.Unmarshal(dist.Manifest, &m); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}
	return m, nil
}
