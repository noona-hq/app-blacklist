package db

import (
	"context"
	"strings"
	"time"

	"github.com/noona-hq/blacklist/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	cfg    Config
	logger logger.Logger
	DB     *mongo.Database
}

func New(cfg Config, logger logger.Logger) (*Database, error) {
	d := &Database{
		cfg:    cfg,
		logger: logger,
	}

	if err := d.connect(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Database) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoOps := &options.ClientOptions{}

	mongoOps.SetDirect(d.cfg.DirectConnection)
	mongoOps.SetMaxPoolSize(200)

	trimmedConnectionString := strings.TrimSuffix(d.cfg.Connection, "\n")

	d.logger.Info("Connecting to database")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(trimmedConnectionString), mongoOps)
	if err != nil {
		return err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	d.logger.Info("Database connected")

	d.DB = client.Database(d.cfg.Name)

	return nil
}
