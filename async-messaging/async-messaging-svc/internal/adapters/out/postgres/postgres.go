package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/logger"

	"github.com/lib/pq"
)

type PostgresDB struct {
	db     *sql.DB
	config *config.DB
}

var instance *PostgresDB
var once sync.Once

func GetInstance(c *config.DB) (*PostgresDB, error) {
	var err error
	once.Do(func() {
		instance = &PostgresDB{
			config: c,
		}
		err = instance.connect()
	})

	if err != nil {
		logger.LoggerError().Error(err)
		return nil, err
	}
	return instance, nil
}

func (p *PostgresDB) connect() error {
	var err error

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", p.config.User, p.config.Pass, p.config.Host, p.config.Port, p.config.Name)

	p.db, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.LoggerError().Errorf("Error al conectarse a la base postgres: %v", err)
		return err
	}

	err = p.db.Ping()
	if err != nil {
		logger.LoggerError().Errorf("Error al conectarse a la base postgres: %v", err)
		return err
	}

	go p.reconnectOnFailure()
	return nil
}

func (p *PostgresDB) reconnectOnFailure() {
	for {
		time.Sleep(10 * time.Second)
		err := p.db.Ping()
		if err != nil {
			logger.LoggerInfo().Info("Se perdió conexión con la base de datos Postgres, reconectando...")
			p.db.Close()
			p.connect()
		}
	}
}

func (p *PostgresDB) GetDB() *sql.DB {
	return p.db
}

func (p *PostgresDB) MapPostgresError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateKey
		case "23503":
			return ErrForeignKeyViolation
		case "23502":
			return ErrNotNullViolation
		case "22P02":
			return ErrInvalidInput
		case "08003":
			return ErrConnectionDoesNotExist
		case "57P01":
			return ErrServerShutdown
		case "42601":
			return ErrSyntaxError
		case "40001":
			return ErrDeadlockDetected
		default:
			return ErrInternalServer
		}
	}
	return err
}

func (p *PostgresDB) Close() error {
	fmt.Println("🧹 Cerrando conexión a Postgres...")
	if p.db != nil {
		return p.db.Close()
	}
	fmt.Println("✅ Postgres cerrado")
	return nil
}
