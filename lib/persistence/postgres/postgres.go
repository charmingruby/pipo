package postgres

import (
	"context"
	"fmt"

	"github.com/charmingruby/pipo/lib/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	Conn   *sqlx.DB
	logger *logger.Logger
}

type ConnectionInput struct {
	User         string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SSL          string
}

func New(logger *logger.Logger, in ConnectionInput) (*Client, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		in.User,
		in.Password,
		in.Host,
		in.Port,
		in.DatabaseName,
		in.SSL,
	)

	dbDriver := "postgres"

	db, err := sqlx.Connect(dbDriver, connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Client{Conn: db, logger: logger}, nil
}

func (c *Client) Close(ctx context.Context) error {
	c.Conn.Close()

	select {
	case <-ctx.Done():
		c.logger.Error("failed to close postgres connection", "error", ctx.Err())
		return ctx.Err()
	default:
		c.logger.Debug("postgres connection closed")
		return nil
	}
}
