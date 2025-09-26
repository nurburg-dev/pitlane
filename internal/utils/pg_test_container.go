package utils

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PGTestContainer struct {
	container  *postgres.PostgresContainer
	connString string
	host       string
	port       string
	pool       *pgxpool.Pool
}

const (
	TestDBName     = "testdb"
	TestDBUser     = "testuser"
	TestDBPassword = "testpass"
)

func GetPGTestContainer(ctx context.Context) (*PGTestContainer, error) {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(TestDBName),
		postgres.WithUsername(TestDBUser),
		postgres.WithPassword(TestDBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(time.Minute),
		),
	)
	if err != nil {
		return nil, err
	}
	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	port := mappedPort.Port()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if pingErr := pool.Ping(ctx); pingErr != nil {
		return nil, pingErr
	}
	return &PGTestContainer{
		container:  postgresContainer,
		connString: connString,
		pool:       pool,
		host:       host,
		port:       port,
	}, nil
}

func (c *PGTestContainer) GetPool() *pgxpool.Pool {
	return c.pool
}

func (c *PGTestContainer) GetHost() string {
	return c.host
}

func (c *PGTestContainer) GetPort() string {
	return c.port
}

func (c *PGTestContainer) GetUsername() string {
	return TestDBUser
}

func (c *PGTestContainer) GetDatabase() string {
	return TestDBName
}

func (c *PGTestContainer) GetPassword() string {
	return TestDBPassword
}

func (c *PGTestContainer) Close(ctx context.Context) error {
	c.pool.Close()
	return c.container.Terminate(ctx)
}
