package actions

import (
	"net/http"
	"time"

	"github.com/katabole/kbsession"
	"github.com/markbates/goth/gothic"
)

const CookieLifetime = 2 * time.Hour

func (app *App) AuthCallback(w http.ResponseWriter, r *http.Request) {
	_, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.render.JSONError(w, r, http.StatusUnauthorized, err)
		return
	}

	// TODO(dk): get the membership end of auth working. We use this GetMemberGivenEmail function in a bunch of apps and
	// it probably should be in a membership/tribe client.
	//
	//v, exist, err := GetMemberGivenEmail(c, user.Email)
	//if err != nil {
	//	app.render.JSONError(w, r, http.StatusUnauthorized, err)
	//	return
	//}

	//if !exist {
	//	RenderJSON(w, r, http.StatusUnauthorized, map[string]string{"message": "no member found"})
	//	return
	//}

	s := kbsession.Get(r)
	//s["UserID"] = v.ID
	//s.Values["UserName"] = v.Name
	//s.Values["UserEmail"] = v.Gpmail
	s.Values["LastUsed"] = time.Now().Unix()
	app.render.Redirect(w, r, "/", http.StatusSeeOther)
}

// RequireLogin checks whether or not a user is logged in with an unexpired session cookie
// if the user is not logged in, then the frontend should redirect to /auth/google
func (app *App) RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if s.Values["UserID"] == nil {
			if app.conf.DeployEnv.IsProduction() || app.conf.EnforceAuth {
				app.render.Redirect(w, r, "/auth?provider=google", http.StatusSeeOther)
				return
			}

			// TODO(dk): for now just set the user id to a user that's in the seed data (test_joe)
			// Ideally we would redirect to a react page (for dev only) that lets us choose which seed user we want
			// to be.
			s.Values["UserID"] = 9000
			s.Values["UserName"] = "Joe Schmoe"
			s.Values["UserEmail"] = "joe.schmoe@example.com"
		}

		// Make auth info available for templates
		s.Values["user_id"] = s.Values["UserID"]
		s.Values["user_name"] = s.Values["UserName"]
		s.Values["user_email"] = s.Values["UserEmail"]

		next(w, r)
	}
}

func (app *App) LogoutGET(w http.ResponseWriter, r *http.Request) {
	clear(kbsession.Get(r).Values)
	app.render.Redirect(w, r, "/", http.StatusSeeOther)
}
