package actions

import (
	"net/http"

	"github.com/katabole/kbexample/public/dist"
	"github.com/markbates/goth/gothic"
)

func (app *App) defineRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth", gothic.BeginAuthHandler)
	mux.HandleFunc("GET /auth/provider/callback", app.AuthCallback)

	mux.HandleFunc("GET /{$}", app.HomeGET)
	mux.HandleFunc("GET /logout", app.LogoutGET)

	mux.HandleFunc("GET /users/new", app.UserNewGET)
	mux.HandleFunc("POST /users", app.UserPOST)
	mux.HandleFunc("GET /users/{id}", app.UserGET)
	mux.HandleFunc("GET /users/{id}/edit", app.UserEditGET)
	mux.HandleFunc("PUT /users/{id}", app.UserPUT)
	mux.HandleFunc("POST /users/{id}/update", app.UserPUT)
	mux.HandleFunc("DELETE /users/{id}", app.UserDELETE)
	mux.HandleFunc("POST /users/{id}/delete", app.UserDELETE)
	mux.HandleFunc("GET /users", app.UsersGET)

	// This will match if nothing else does
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := "404 Not Found"
		if GetContentType(r) == ContentTypeHTML {
			app.render.HTML(w, r, http.StatusOK, "error", msg)
		} else {
			app.render.JSON(w, r, http.StatusNotFound, map[string]string{"error": msg})
		}
	})

	mux.Handle("GET /assets/{path...}", http.FileServerFS(dist.BuiltAssets))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, dist.BuiltAssets, "assets/images/favicon.ico")
	})
	return mux
}
