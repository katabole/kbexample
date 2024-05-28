package actions

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/dankinder/gobase/gbexample/models"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conf Config

func TestMain(m *testing.M) {
	if err := godotenv.Load("../env/test.env"); err != nil {
		log.Fatalf("Error loading test dotenv: %v", err)
	}

	if err := envconfig.Process("", &conf); err != nil {
		log.Fatalf("Error loading app config from environment: %v", err)
	}

	atlasDevDBConf := conf.DBConfig
	atlasDevDBConf.DBName = "atlas_dev"
	if err := SetupDB(conf.DBConfig, atlasDevDBConf); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	os.Exit(m.Run())
}

// Setup starts a local test server and returns it along with a cleanup function that should be deferred.
func Setup(t *testing.T) (*App, func()) {
	app, err := NewApp(conf)
	require.Nil(t, err)

	app.Start()
	CleanDB(t, app.db.DB)
	return app, func() { assert.Nil(t, app.Stop(context.Background())) }
}

// SetupDB runs atlas to ensure the database is up to date.
func SetupDB(dbConf models.Config, atlasDevDBConf models.Config) error {
	cmd := exec.Command("atlas", "schema", "apply",
		"--to", "file://schema.sql",
		"--url", dbConf.URL(),
		"--dev-url", atlasDevDBConf.URL(),
		"--auto-approve")
	cmd.Dir = ".."
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running atlas schema apply: %v\n\nOutput: %s\n", err, output)
	}
	return nil
}

// CleanDB resets sequences and wipes all tables in this database.
func CleanDB(t *testing.T, db *sqlx.DB) {
	var sequences []string
	err := db.Select(&sequences,
		`SELECT sequence_name
		FROM information_schema.sequences
		WHERE sequence_schema NOT IN ('information_schema, pg_catalog')`)
	require.Nil(t, err)

	for _, s := range sequences {
		_, err = db.Exec("ALTER SEQUENCE " + s + " RESTART WITH 1")
		require.Nil(t, err)
	}

	var tables []struct {
		Name   string `db:"table_name"`
		Schema string `db:"table_schema"`
	}
	err = db.Select(&tables,
		`SELECT table_name, table_schema
		FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema', 'pg_catalog') AND table_type = 'BASE TABLE'`)
	require.Nil(t, err)

	for _, table := range tables {
		_, err = db.Exec("DELETE FROM " + table.Schema + "." + table.Name + " CASCADE")
		require.Nil(t, err)
	}
}
