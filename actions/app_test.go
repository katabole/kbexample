package actions

import (
	"context"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/katabole/kbhttp"
	"github.com/katabole/kbsql"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conf Config

func TestMain(m *testing.M) {
	// Use Overload rather than Load since the user likely has dev env vars loaded, but we want to overwrite them.
	if err := godotenv.Overload("../env/test.env"); err != nil {
		log.Fatalf("Error loading test dotenv: %v", err)
	}

	if err := envconfig.Process("", &conf); err != nil {
		log.Fatalf("Error loading app config from environment: %v", err)
	}

	atlasDevDBConf := conf.DBConfig
	atlasDevDBConf.DBName = "atlas_dev"
	if err := kbsql.AtlasSetupDB(conf.DBConfig.URL(), atlasDevDBConf.URL()); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	os.Exit(m.Run())
}

type Fixture struct {
	t      *testing.T
	App    *App
	Client *kbhttp.Client
}

// NewFixture starts a local test server and returns it along with a cleanup function that should be deferred.
func NewFixture(t *testing.T) *Fixture {
	app, err := NewApp(conf)
	require.Nil(t, err)

	app.Start()
	require.NoError(t, kbsql.PostgresCleanDB(app.db.DB))

	baseURL, err := url.Parse("http://" + app.srv.Addr)
	require.Nil(t, err)

	return &Fixture{
		t:      t,
		App:    app,
		Client: kbhttp.NewClient(kbhttp.ClientConfig{BaseURL: baseURL}),
	}
}

func (f *Fixture) Cleanup() {
	assert.Nil(f.t, f.App.Stop(context.Background()))
}
