package actions

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/katabole/kbexample/build"
	"github.com/markbates/goth/gothic"
	"github.com/olivere/vite"
)

// defineRoutes is the part of app setup where routes/endpoints are defined.
// For how to define these routes on the chi Mux, see https://go-chi.io/#/pages/routing
func (app *App) defineRoutes(r *chi.Mux) error {
	r.Get("/auth", gothic.BeginAuthHandler)
	r.Get("/auth/google/callback", app.AuthCallback)

	r.Get("/", app.HomeGET)
	r.Get("/logout", app.LogoutGET)

	r.Group(func(r chi.Router) {
		r.Use(app.RequireLogin)
		r.Get("/users/new", app.UserNewGET)
		r.Post("/users", app.UserPOST)
		r.Get("/users/{id}", app.UserGET)
		r.Get("/users/{id}/edit", app.UserEditGET)
		r.Put("/users/{id}", app.UserPUT)
		r.Post("/users/{id}/update", app.UserPUT)
		r.Delete("/users/{id}", app.UserDELETE)
		r.Post("/users/{id}/delete", app.UserDELETE)
		r.Get("/users", app.UsersGET)
	})

	// For any special file that needs to be served not under /assets/, add the route here.
	for _, f := range []string{"robots.txt"} {
		r.Get("/"+f, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, build.DistDir(), f)
		})
	}

	assetHandler, err := app.getAssetHandler()
	if err != nil {
		return fmt.Errorf("could not create asset handler: %w", err)
	}
	r.Handle("/assets/*", assetHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.render.Error(w, r, http.StatusNotFound, errors.New("404 Not Found"))
	})
	return nil
}

// getAssetHandler returns the handler for serving static assets. In dev/test this directly serves the local "public"
// directory. In production build, Vite copies those assets into the build/dist directory where they get embedded into
// the binary, and the vite handler serves them.
func (app *App) getAssetHandler() (http.Handler, error) {
	if app.conf.DeployEnv.IsProduction() {
		viteHandler, err := vite.NewHandler(vite.Config{
			FS:           build.DistDir(),
			IsDev:        app.conf.DeployEnv.IsProduction(),
			ViteEntry:    "js/main.js",
			ViteTemplate: vite.Vanilla,
			ViteManifest: "dist/manifest.json",
		})
		if err != nil {
			return nil, fmt.Errorf("could not create Vite handler: %w", err)
		}
		return viteHandler, nil
	} else {
		return http.FileServerFS(os.DirFS("public")), nil
	}
}
