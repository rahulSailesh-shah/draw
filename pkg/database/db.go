package database

import (
	"context"
	"log"
	"time"

	cfg "draw/pkg/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Connect() error
	Close() error
	GetDB() *pgxpool.Pool
}

type postgresDB struct {
	ctx    context.Context
	db     *pgxpool.Pool
	config *cfg.DBConfig
}

func NewPostgresDB(ctx context.Context, config *cfg.DBConfig) DB {
	return &postgresDB{
		ctx:    ctx,
		config: config,
	}
}

func (p *postgresDB) Connect() error {
	config, err := pgxpool.ParseConfig("")
	if err != nil {
		return err
	}

	config.ConnConfig.Host = p.config.Host
	config.ConnConfig.Port = uint16(p.config.Port)
	config.ConnConfig.User = p.config.User
	config.ConnConfig.Database = p.config.Name
	config.ConnConfig.Password = p.config.Password
	config.MaxConns = 50

	config.MinConns = 10
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	db, err := pgxpool.NewWithConfig(p.ctx, config)
	if err != nil {
		return err
	}

	p.db = db

	log.Println("Connected to PostgreSQL database successfully")

	return nil
}

func (p *postgresDB) Close() error {
	if p.db == nil {
		log.Println("PostgreSQL database connection is already closed")
		return nil
	}
	p.db.Close()
	log.Println("PostgreSQL database connection closed")
	return nil
}

func (p *postgresDB) GetDB() *pgxpool.Pool {
	return p.db
}
