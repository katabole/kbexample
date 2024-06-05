package actions

import (
	"context"
	"encoding/gob"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionContextKey int

const sessionKey sessionContextKey = 0

type sessionHandler struct {
	app  *App
	next http.Handler
}

func (sh *sessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := sh.app.sessionStore.Get(r, "RootSession")
	if err != nil {
		// This should be extremely rare, it only happens if there's a session that can't be decoded.
		// If there's no existing session yet it's not an error, Get just returns a new one.
		slog.Error("Failed to load session", "err", err)
		http.Error(w, "Failed to load session, check logs for details", http.StatusInternalServerError)
		return
	}
	sh.next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), sessionKey, session)))
}

// SessionHandler decodes the session from the request and adds it to the request context, guaranteeing it'll be
// there and making gorilla/sessions easier to work with.
func (app *App) SessionHandler(next http.Handler) http.Handler {
	return &sessionHandler{app: app, next: next}
}

// Session returns the session from the request context.
func Session(r *http.Request) *sessions.Session {
	return r.Context().Value(sessionKey).(*sessions.Session)
}

// SaveSession saves the final session in the request if it's been accessed (i.e. new or modified).
func SaveSession(w http.ResponseWriter, r *http.Request) {
	s := Session(r)

	// Avoid unnecessarily saving a session if the request didn't come with one and nothing was added to it.
	if s.IsNew && len(s.Values) == 0 {
		return
	}

	if err := s.Save(r, w); err != nil {
		slog.Error("Failed to save session", "err", err)
	}
}

// Flash support
//

func init() {
	// gorilla/sessions uses gob to encode/decode session values, and requires us to register this type we're using to
	// store flash data.
	gob.Register(map[string][]string{})
}

const flashKey = "_flash_"

// AddFlash adds a flash message to the session for display at the next page render.
// The key groups together messages of a similar category and will be interpeted by the template.
// It's common to use keys like "success", "info", "warning", and "error" which map to CSS classes.
// Value is the string to be displayed.
func AddFlash(r *http.Request, key, value string) {
	s := Session(r)
	if flashMap, ok := s.Values[flashKey]; ok && flashMap != nil {
		flashMap := flashMap.(map[string][]string)
		flashMap[key] = append(flashMap[key], value)
	} else {
		s.Values[flashKey] = map[string][]string{key: []string{value}}
	}
}

// Flash grabs the flash messages from the session and removes them so they'll only be rendered once.
func Flash(r *http.Request) map[string][]string {
	s := Session(r)
	if flashMap, ok := s.Values[flashKey]; ok && flashMap != nil {
		delete(s.Values, flashKey)
		return flashMap.(map[string][]string)
	} else {
		return map[string][]string{}
	}
}
