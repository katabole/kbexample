package actions

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/hashicorp/go-multierror"
	"github.com/katabole/kbexample/models"
	"github.com/katabole/kbsession"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/unrolled/secure"
)

type Config struct {
	ServerAddr    string        `envconfig:"SERVER_ADDR"`
	SessionSecret string        `envconfig:"SESSION_SECRET"`
	DeployEnv     Environment   `envconfig:"DEPLOY_ENV"`
	EnforceAuth   bool          `envconfig:"ENFORCE_AUTH"`
	SiteURL       string        `envconfig:"SITE_URL"`
	DBConfig      models.Config `envconfig:"DB"`

	GoogleOAuthKey    string `envconfig:"GOOGLE_OAUTH_KEY"`
	GoogleOAuthSecret string `envconfig:"GOOGLE_OAUTH_SECRET"`
}

type App struct {
	conf   Config
	srv    *http.Server
	render *Renderer
	db     *models.DB
}

func NewApp(conf Config) (*App, error) {
	app := &App{conf: conf}

	// Set up the database
	var err error
	app.db, err = models.NewDB(conf.DBConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create database: %w", err)
	}

	// Configure our session store. For test/dev it can be a dummy but for production it must be secure.
	var sessionStore sessions.Store
	if conf.DeployEnv.IsProduction() {
		if conf.SessionSecret == "" {
			return nil, fmt.Errorf("SESSION_SECRET must be set")
		}
		s := sessions.NewCookieStore([]byte(conf.SessionSecret))
		s.Options.Secure = true
		s.Options.HttpOnly = true
		sessionStore = s
	} else {
		if conf.SessionSecret == "" {
			conf.SessionSecret = "not-so-super-secret"
		}
		sessionStore = sessions.NewCookieStore([]byte(conf.SessionSecret))
	}

	// Set up oauth, which is configured globally here and applied in routes.go
	gothic.Store = sessionStore
	goth.UseProviders(
		google.New(conf.GoogleOAuthKey, conf.GoogleOAuthSecret, fmt.Sprintf(conf.SiteURL+"/auth/google/callback")),
	)

	// Define our router middleware (logging, etc.), then define routes
	router := chi.NewRouter()
	// TODO(dk): form POST example with csrf protection https://github.com/gorilla/csrf
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(secure.New(secure.Options{
		IsDevelopment:   !conf.DeployEnv.IsProduction(),
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	}).Handler)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(kbsession.NewMiddleware(sessionStore))

	app.defineRoutes(router)

	app.srv = &http.Server{
		Addr:    conf.ServerAddr,
		Handler: router,
	}
	app.render, err = NewRenderer(conf.DeployEnv.IsProduction())
	if err != nil {
		return nil, fmt.Errorf("could not create renderer: %w", err)
	}

	return app, nil
}

// Start begins listening for connections and serving clients in the background.
func (app *App) Start() {
	go func() {
		if err := app.srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Warn("Server encountered an unexpected error while stopping", "err", err)
		}
	}()
	slog.Info("Server listening", "addr", app.conf.ServerAddr)
}

// Stop gracefully shuts down the server and closes the database. Set a timeout on the provided context to force
// shutdown after a certain amount of time.
func (app *App) Stop(ctx context.Context) error {
	var result error
	if err := app.srv.Shutdown(ctx); err != nil {
		result = multierror.Append(result, fmt.Errorf("could not shutdown server: %w", err))
	}
	if err := app.db.Close(); err != nil {
		result = multierror.Append(result, fmt.Errorf("could not close database: %w", err))
	}
	return result
}
