package endpoint_test

import (
	"fmt"
	"log"
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

	db *sqlx.DB
)

func TestEndpoints(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Endpoint Suite", Label("integration"))
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	pool, err = dockertest.NewPool("")
	Expect(err).ShouldNot(HaveOccurred())

	dbURL := runDockerPostgres(pool)

	return []byte(dbURL)
}, func(dbURL []byte) {
	dbConn, err := sqlx.Open("pgx", string(dbURL))
	Expect(err).ShouldNot(HaveOccurred())

	db = dbConn
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	Expect(pool.Purge(postgresResource)).ShouldNot(HaveOccurred())
})

func runDockerPostgres(pool *dockertest.Pool) string {
	var err error
	postgresResource, err = pool.RunWithOptions(&dockertest.RunOptions{
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

	hostAndPort := postgresResource.GetHostPort("5432/tcp")
	dbUrl := fmt.Sprintf("postgres://admin:secret@%s/event_scheduling_test?sslmode=disable", hostAndPort)

	if err := pool.Retry(func() error {
		retryConn, err := sqlx.Open("pgx", dbUrl)
		if err != nil {
			return err
		}

		err = retryConn.Ping()
		if err != nil {
			return err
		}

		return migrateDB(dbUrl)
	}); err != nil {
		log.Fatalf("Could not connect to postgres docker: %s", err)
	}

	return dbUrl
}

func migrateDB(dbUrl string) error {
	m, err := migrate.New(
		"file://../../sql/migrations",
		dbUrl)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
