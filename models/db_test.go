package models

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/katabole/kbsql"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conf Config

func TestMain(m *testing.M) {
	if err := godotenv.Load("../env/test.env"); err != nil {
		log.Fatalf("Error loading test dotenv: %v", err)
	}

	if err := envconfig.Process("DB", &conf); err != nil {
		log.Fatalf("Error loading app config from environment: %v", err)
	}

	atlasDevDBConf := conf
	atlasDevDBConf.DBName = "atlas_dev"
	if err := kbsql.AtlasSetupDB(conf.URL(), atlasDevDBConf.URL()); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	os.Exit(m.Run())
}

type Fixture struct {
	t  *testing.T
	db *DB
}

// Setup starts a local test server and returns it along with a cleanup function that should be deferred.
func NewFixture(t *testing.T) *Fixture {
	db, err := NewDB(conf)
	require.NoError(t, err)
	require.NoError(t, kbsql.PostgresCleanDB(db.DB))

	return &Fixture{
		t:  t,
		db: db,
	}
}

func (f *Fixture) Cleanup() {
	assert.NoError(f.t, f.db.Close())
}
