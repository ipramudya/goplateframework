package db

import (
	"fmt"
	"time"

	"github.com/goplateframework/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func Init(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s options=endpoint=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Dbname,
		cfg.DB.SSLMode,
		cfg.DB.EndpointID,
	)

	db, err := sqlx.Connect(cfg.DB.Driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIddleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.DB.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.DB.ConnMaxIddleTime) * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
