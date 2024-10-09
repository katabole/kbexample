package actions

import (
	"net/http"
	"time"

	"github.com/katabole/kbsession"
	"github.com/markbates/goth/gothic"
)

const CookieLifetime = 2 * time.Hour

func (app *App) AuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.render.JSONError(w, r, http.StatusUnauthorized, err)
		return
	}

	s := kbsession.Get(r)
	s.Values["UserName"] = user.Name
	s.Values["UserEmail"] = user.Email
	s.Values["LastUsed"] = time.Now().Unix()
	app.render.Redirect(w, r, "/", http.StatusSeeOther)
}

// RequireLogin checks whether or not a user is logged in with an unexpired session cookie
// if the user is not logged in, then the frontend should redirect to /auth/google
func (app *App) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := kbsession.Get(r)
		lastUsed := s.Values["LastUsed"]

		if lastUsed != nil {
			lastUsedEpoch := time.Unix(lastUsed.(int64), 0)
			if time.Since(lastUsedEpoch) > CookieLifetime {
				clear(s.Values)
			} else {
				s.Values["LastUsed"] = time.Now().Unix()
			}
		}

		if s.Values["UserEmail"] == nil {
			if app.conf.DeployEnv.IsProduction() || app.conf.EnforceAuth {
				app.render.Redirect(w, r, "/auth?provider=google", http.StatusSeeOther)
				return
			}

			// In dev/test, just use a test user.
			s.Values["UserName"] = "Joe Schmoe"
			s.Values["UserEmail"] = "joe.schmoe@example.com"
		}
		next.ServeHTTP(w, r)
	})
}

func (app *App) LogoutGET(w http.ResponseWriter, r *http.Request) {
	clear(kbsession.Get(r).Values)
	app.render.Redirect(w, r, "/", http.StatusSeeOther)
}
