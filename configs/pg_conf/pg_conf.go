package pg_conf

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"os"
	"shop/pkg/log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	postgresDB    *sql.DB
	once          sync.Once
	postgresURI   string
	checkInterval = 15 * time.Second
)

// InitPostgresSingleton инициализирует подключение к PostgreSQL в виде синглтона и запускает мониторинг соединения.
func InitPostgresSingleton() {
	once.Do(func() {
		postgresURI = os.Getenv("POSTGRES_URI")

		if postgresURI == "" {
			log.Fatal("POSTGRES_URI is not set in environment variables")
		}
		db, err := connectPostgres(postgresURI)
		if err != nil {
			log.Fatal("Failed to connect to Postgres", zap.Error(err))
		}

		log.Info("Connected to Postgres successfully!")
		postgresDB = db

		//// Запускаем горутину для периодической проверки соединения и реконнекта
		go monitorConnection()
	})
}

// connectPostgres подключается к PostgreSQL
func connectPostgres(uri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(7)
	db.SetConnMaxLifetime(1 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// monitorConnection каждые 15 секунд проверяет соединение и при необходимости пытается переподключиться.
func monitorConnection() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err := postgresDB.PingContext(ctx)
		cancel()

		if err != nil {
			log.Error("Postgres connection lost, trying to reconnect...", zap.Error(err))
			reconnect()
		}
	}
}

// reconnect пытается переподключиться к БД до успешного результата.
func reconnect() {
	for {
		db, err := connectPostgres(postgresURI)
		if err == nil {
			oldDB := postgresDB
			postgresDB = db
			if oldDB != nil {
				oldDB.Close()
			}
			log.Info("Reconnected to Postgres successfully!")
			return
		} else {
			log.Error("Failed to reconnect to Postgres, will retry in 15 seconds", zap.Error(err))
			time.Sleep(checkInterval)
		}
	}
}

// GetDB возвращает инициализированный экземпляр *sql.DB.
func GetDB() *sql.DB {
	if postgresDB == nil {
		log.Fatal("Postgres DB is not initialized. Call InitPostgresSingleton() first.")
	}
	return postgresDB
}

// ClosePostgresDB закрывает соединение с БД.
func ClosePostgresDB() {
	if postgresDB != nil {
		if err := postgresDB.Close(); err != nil {
			log.Error("Failed to close Postgres connection", zap.Error(err))
		} else {
			log.Info("Postgres connection closed successfully.")
		}
	}
}
