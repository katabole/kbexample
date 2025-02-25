package actions

import "net/http"

// HomeGET is a default handler to serve up a home page.
func (app *App) HomeGET(w http.ResponseWriter, r *http.Request) {
	app.render.HTML(w, r, HTMLParams{Template: "home"})
}
