package actions

import (
	"context"
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
	if !s.IsNew {
		if err := s.Save(r, w); err != nil {
			slog.Error("Failed to save session", "err", err)
		}
	}
}
