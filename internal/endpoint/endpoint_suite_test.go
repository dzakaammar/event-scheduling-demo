//go:build integration

package endpoint_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var (
	pool             *dockertest.Pool
	postgresResource *dockertest.Resource

	dbUrl string
	db    *sqlx.DB
)

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	postgresResource = runDockerPostgres(pool)

	migrateDB()

	code := m.Run()

	if err := pool.Purge(postgresResource); err != nil {
		log.Fatalf("Could not purge postgres resource: %s", err)
	}

	os.Exit(code)
}

func TestEndpoints(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Endpoint Suite")
}

func migrateDB() {
	m, err := migrate.New(
		"file://../../sql/migrations",
		dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func runDockerPostgres(pool *dockertest.Pool) *dockertest.Resource {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "event_scheduling_demo_it",
		Repository: "postgres",
		Tag:        "12-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=admin",
			"POSTGRES_DB=event_scheduling_test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start postgres resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	dbUrl = fmt.Sprintf("postgres://admin:secret@%s/event_scheduling_test?sslmode=disable", hostAndPort)

	if err = pool.Retry(func() error {
		db, err = sqlx.Open("pgx", dbUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to postgres docker: %s", err)
	}

	return resource
}
