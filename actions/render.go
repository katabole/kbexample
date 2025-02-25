package actions

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/katabole/kbexample/build"
	"github.com/katabole/kbexample/templates"
	"github.com/katabole/kbsession"
	"github.com/olivere/vite"
	"github.com/unrolled/render"
)

// Renderer wraps the unrolled/render package in order to provide a few goodies (error rendering, session saving).
type Renderer struct {
	rnd          *render.Render
	isProduction bool
	viteFragment *vite.Fragment
}

func NewRenderer(isProduction bool) (*Renderer, error) {
	var err error

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
	})

	// This object is used by layout.html to render the Vite fragment.
	// That fragment is js/css imports, which use hot module reloading in dev and are static in production.
	r.viteFragment, err = vite.HTMLFragment(vite.Config{
		FS:           build.DistDir(),
		IsDev:        !isProduction,
		ViteEntry:    "js/main.js",
		ViteTemplate: vite.Vanilla,
		ViteManifest: "manifest.json",
	})
	if err != nil {
		return nil, fmt.Errorf("could not create Vite fragment: %w", err)
	}

	return r, nil
}

// Data writes out the raw bytes as binary data.
func (r *Renderer) Data(w http.ResponseWriter, req *http.Request, status int, v []byte) error {
	kbsession.Save(w, req)
	return r.rnd.Data(w, status, v)
}

// HTMLParams provides all the HTML function needs to render an HTML template.
type HTMLParams struct {
	// HTTP Status, defaults to http.StatusOK (200).
	Status int
	// Name of the template to render, e.g. "users/new" for templates/users/new.html.
	Template string
	// Data to pass to the template, often a map of key/values.
	Data any
	// Options for the renderer.
	HTMLOptions []render.HTMLOptions
	// A page title used by the layout.
	Title string
}

// HTML builds up the response from the specified parameters.
func (r *Renderer) HTML(w http.ResponseWriter, req *http.Request, params HTMLParams) error {
	flash := kbsession.Flash(req)
	session := kbsession.Get(req)
	kbsession.Save(w, req)

	if params.Status == 0 {
		params.Status = http.StatusOK
	}
	if params.Data == nil {
		params.Data = map[string]any{}
	}
	return r.rnd.HTML(w, params.Status, params.Template, map[string]any{
		"Flash":   flash,
		"Session": session,
		"Title":   params.Title,
		"Vite":    r.viteFragment,
		"Data":    params.Data,
	}, params.HTMLOptions...)
}

// JSON marshals the given interface object and writes the JSON response.
func (r *Renderer) JSON(w http.ResponseWriter, req *http.Request, status int, v interface{}) error {
	kbsession.Save(w, req)
	return r.rnd.JSON(w, status, v)
}

// JSONP marshals the given interface object and writes the JSON response.
func (r *Renderer) JSONP(w http.ResponseWriter, req *http.Request, status int, callback string, v interface{}) error {
	kbsession.Save(w, req)
	return r.rnd.JSONP(w, status, callback, v)
}

// Text writes out a string as plain text.
func (r *Renderer) Text(w http.ResponseWriter, req *http.Request, status int, v string) error {
	kbsession.Save(w, req)
	return r.rnd.Text(w, status, v)
}

// XML marshals the given interface object and writes the XML response.
func (r *Renderer) XML(w http.ResponseWriter, req *http.Request, status int, v interface{}) error {
	kbsession.Save(w, req)
	return r.rnd.XML(w, status, v)
}

func (r *Renderer) Redirect(w http.ResponseWriter, req *http.Request, url string, status int) {
	kbsession.Save(w, req)
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
		r.HTML(w, req, HTMLParams{
			Status:   status,
			Template: "error",
			Data:     map[string]string{"Message": "internal error, see logs for details"},
		})
		return
	}

	r.HTML(w, req, HTMLParams{Status: status, Template: "error", Data: map[string]string{"Message": err.Error()}})
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
