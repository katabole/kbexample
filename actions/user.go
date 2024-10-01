package actions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/katabole/kbexample/models"
	usertemplates "github.com/katabole/kbexample/templates/users"
	"github.com/katabole/kbsession"
	"github.com/monoculum/formam"
)

// UsersGET handles GET /users
func (app *App) UsersGET(w http.ResponseWriter, r *http.Request) {
	users, err := app.db.GetUsers()
	if err != nil {
		app.render.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if GetContentType(r) == ContentTypeHTML {
		app.render.HTML(w, r, http.StatusOK, usertemplates.List(users))
	} else {
		app.render.JSON(w, r, http.StatusOK, map[string]interface{}{"users": users})
	}
}

// UserNewGET handles GET /users/new
func (app *App) UserNewGET(w http.ResponseWriter, r *http.Request) {
	app.render.HTML(w, r, http.StatusOK, usertemplates.New(nil))
}

// UserPOST handles POST /users
func (app *App) UserPOST(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if GetContentType(r) == ContentTypeHTML {
		if err := r.ParseForm(); err != nil {
			app.render.Error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := formam.Decode(r.Form, &u); err != nil {
			app.render.Error(w, r, http.StatusBadRequest, err)
			return
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			app.render.Error(w, r, http.StatusBadRequest, err)
			return
		}
	}

	newUser, err := app.db.CreateUser(&u)
	if err != nil {
		app.render.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if GetContentType(r) == ContentTypeHTML {
		app.render.Redirect(w, r, fmt.Sprintf("/users/%d", newUser.ID), http.StatusSeeOther)
	} else {
		app.render.JSON(w, r, http.StatusCreated, newUser)
	}
}

// UserGET handles GET /users/{id}
func (app *App) UserGET(w http.ResponseWriter, r *http.Request) {
	if u := app.getUserHelper(w, r); u != nil {
		if GetContentType(r) == ContentTypeHTML {
			app.render.HTML(w, r, http.StatusOK, usertemplates.Show(u))
		} else {
			app.render.JSON(w, r, http.StatusOK, u)
		}
	}
}

// UserEditGET handles GET /users/{id}/edit
func (app *App) UserEditGET(w http.ResponseWriter, r *http.Request) {
	if u := app.getUserHelper(w, r); u != nil {
		app.render.HTML(w, r, http.StatusOK, usertemplates.New(u))
	}
}

// getUserHelper gets the user, or returns nil in which case it has already sent back an error.
func (app *App) getUserHelper(w http.ResponseWriter, r *http.Request) *models.User {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.render.Error(w, r, http.StatusBadRequest, err)
		return nil
	}

	u, err := app.db.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := fmt.Sprintf("User ID %d not found", id)
			if GetContentType(r) == ContentTypeHTML {
				app.render.HTML(w, r, http.StatusNotFound, usertemplates.NotFound(id))
			} else {
				app.render.JSON(w, r, http.StatusNotFound, map[string]string{"message": msg})
			}
		} else {
			app.render.Error(w, r, http.StatusInternalServerError, err)
		}
		return nil
	}
	return u
}

// UserPUT handles PUT /users/{id}
func (app *App) UserPUT(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.render.Error(w, r, http.StatusBadRequest, err)
		return
	}

	var u models.User
	if GetContentType(r) == ContentTypeHTML {
		if err := r.ParseForm(); err != nil {
			app.render.Error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := formam.Decode(r.Form, &u); err != nil {
			app.render.Error(w, r, http.StatusBadRequest, err)
			return
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			app.render.Error(w, r, http.StatusBadRequest, err)
			return
		}
	}

	// Even if they pass an ID in the body, ignore it and use the one from the URL.
	u.ID = id
	if err := app.db.UpdateUser(&u); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.render.Error(w, r, http.StatusNotFound, errors.New("user not found"))
		} else {
			app.render.Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	if GetContentType(r) == ContentTypeHTML {
		app.render.Redirect(w, r, fmt.Sprintf("/users/%d", u.ID), http.StatusSeeOther)
	} else {
		app.render.JSON(w, r, http.StatusCreated, map[string]string{"message": "User updated"})
	}
}

// UserDELETE handles DELETE /users/{id}
func (app *App) UserDELETE(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.render.Error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := app.db.DeleteUser(id); err != nil {
		app.render.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if GetContentType(r) == ContentTypeHTML {
		kbsession.AddFlash(r, "success", "User deleted")
		app.render.Redirect(w, r, "/users", http.StatusSeeOther)
	} else {
		app.render.JSON(w, r, http.StatusOK, map[string]string{"message": "User deleted"})
	}
}
