package actions

import (
	"net/http"

	"github.com/dankinder/katabole/kbexample/public/dist"
	"github.com/markbates/goth/gothic"
)

func (app *App) defineRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth", gothic.BeginAuthHandler)
	mux.HandleFunc("GET /auth/provider/callback", app.AuthCallback)

	mux.HandleFunc("GET /", app.HomeGET)
	mux.HandleFunc("GET /logout", app.LogoutGET)

	mux.HandleFunc("PUT /users/{id}", app.UsersPUT)
	mux.HandleFunc("GET /users/{id}", app.UsersGET)

	mux.Handle("GET /assets/{path...}", http.FileServerFS(dist.BuiltAssets))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, dist.BuiltAssets, "assets/images/favicon.ico")
	})
	return mux
}
