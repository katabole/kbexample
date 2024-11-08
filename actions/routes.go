package actions

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katabole/kbexample/public/dist"
	"github.com/markbates/goth/gothic"
)

func (app *App) defineRoutes(r *chi.Mux) {
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

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.render.Error(w, r, http.StatusNotFound, errors.New("404 Not Found"))
	})

	r.Handle("/assets/*", http.FileServer(http.FS(dist.BuiltAssets)))
	for _, f := range []string{"favicon.ico", "robots.txt"} {
		r.Get("/"+f, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, dist.BuiltAssets, "assets/"+f)
		})
	}
}
