package actions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dankinder/gobase/gbexample/models"
)

func (app *App) UsersGET(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.render.JSONError(w, r, http.StatusBadRequest, err)
		return
	}
	u, err := app.db.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.render.JSON(w, r, http.StatusNotFound, map[string]string{"message": "user not found"})
		} else {
			app.render.JSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	switch AcceptContentType(r) {
	case AcceptHTML:
		app.render.HTML(w, r, http.StatusOK, "users/show", u)
	default:
		app.render.JSON(w, r, http.StatusOK, u)
	}
}

func (app *App) UsersPUT(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.render.JSONError(w, r, http.StatusBadRequest, err)
		return
	}

	// Decode the body into a User struct
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.render.JSONError(w, r, http.StatusBadRequest, err)
		return
	}
	u.ID = id
	if err := app.db.SaveUser(&u); err != nil {
		app.render.JSONError(w, r, http.StatusInternalServerError, err)
		return
	}
}
