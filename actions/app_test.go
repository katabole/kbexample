package actions

import (
	"context"
	"fmt"
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
	// Start our tests by loading config from `env/test.env`.
	// Allow overriding this by setting TEST_ENV, for example for CI
	testEnv := os.Getenv("TEST_ENV")
	if testEnv == "" {
		testEnv = "test"
	}

	if err := godotenv.Overload(fmt.Sprintf("../env/%s.env", testEnv)); err != nil {
		log.Fatalf("Error loading dotenv for %s: %v", testEnv, err)
	}

	if err := envconfig.Process("", &conf); err != nil {
		log.Fatalf("Error loading app config from environment: %v", err)
	}
	log.Printf("Loaded env %s", testEnv)

	// Now create databases if necessary. These may already exist but especially for CI doing it here is convenient.
	if err := kbsql.PostgresCreateDBIfNotExistsByURL(conf.DBConfig.URL()); err != nil {
		log.Fatalf("Error creating database %s: %v", conf.DBConfig.URL(), err)
	}
	atlasDevDBConf := conf.DBConfig
	atlasDevDBConf.DBName = "atlas_dev"
	if err := kbsql.PostgresCreateDBIfNotExistsByURL(atlasDevDBConf.URL()); err != nil {
		log.Fatalf("Error creating database %s: %v", atlasDevDBConf.URL(), err)
	}
	if err := kbsql.AtlasSetupDB(conf.DBConfig.URL(), atlasDevDBConf.URL()); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	os.Exit(m.Run())
}

type Fixture struct {
	t       *testing.T
	App     *App
	Client  *kbhttp.Client
	BaseURL string
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
		t:       t,
		App:     app,
		Client:  kbhttp.NewClient(kbhttp.ClientConfig{BaseURL: baseURL}),
		BaseURL: baseURL.String(),
	}
}

func (f *Fixture) Cleanup() {
	assert.Nil(f.t, f.App.Stop(context.Background()))
}

// URL returns the full URL for the given path.
func (f *Fixture) URL(path string) string {
	return f.BaseURL + path
}
