package platform

import (
	"database/sql"
	"os"
	"strconv"
	"time"
)

const (
	defaultDBPoolMinConns           = 2
	defaultDBPoolMaxConns           = 3
	defaultDBPoolMaxConnIdleSeconds = 300
	defaultDBPoolMaxConnLifeSeconds = 1800
)

type DBPoolConfig struct {
	MinConns        int
	MaxConns        int
	MaxConnIdleTime time.Duration
	MaxConnLifetime time.Duration
}

func LoadDBPoolConfigFromEnv() DBPoolConfig {
	minConns := getDBPoolInt("DB_POOL_MIN_CONNS", defaultDBPoolMinConns)
	maxConns := getDBPoolInt("DB_POOL_MAX_CONNS", defaultDBPoolMaxConns)
	if maxConns < 1 {
		maxConns = defaultDBPoolMaxConns
	}
	if minConns < 0 {
		minConns = 0
	}
	if minConns > maxConns {
		minConns = maxConns
	}

	return DBPoolConfig{
		MinConns:        minConns,
		MaxConns:        maxConns,
		MaxConnIdleTime: time.Duration(getDBPoolInt("DB_POOL_MAX_CONN_IDLE_SECONDS", defaultDBPoolMaxConnIdleSeconds)) * time.Second,
		MaxConnLifetime: time.Duration(getDBPoolInt("DB_POOL_MAX_CONN_LIFETIME_SECONDS", defaultDBPoolMaxConnLifeSeconds)) * time.Second,
	}
}

func ApplyDBPoolConfig(db *sql.DB, cfg DBPoolConfig) {
	db.SetMaxIdleConns(cfg.MinConns)
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)
}

func getDBPoolInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
