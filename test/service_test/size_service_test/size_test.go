package article_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"shop/configs/env"
	"shop/configs/pg_conf"
	"shop/internal/repository"
	"shop/internal/service"
	"shop/pkg/log"
	"testing"
)

func TestMain(m *testing.M) {
	// Запуск тестов
	os.Exit(m.Run())
}

func setup() (context.Context, service.SizeServiceInterface) {
	err := os.Setenv("POSTGRES_URI", "user=alex password=1000 host=127.0.0.1 port=5432 dbname=shop")
	if err != nil {
		fmt.Println("Failed to set env POSTGRES_URI")
		return nil, nil
	}

	// Загружаем переменные окружения
	env.LoadEnv()

	log.InitLogger()
	pg_conf.InitPostgresSingleton()
	clientDB := pg_conf.GetDB()

	ctx := context.Background()

	repo := repository.NewSizeRepository(clientDB)
	serv := service.NewSizeService(repo)
	return ctx, serv
}

func TestGetAllSize(t *testing.T) {
	_, serv := setup()
	// Запуск теста
	t.Run("Success", func(t *testing.T) {
		sizeWithPagination, err := serv.GetAllSizes(1, 10)
		// Вывод структуры с ключами
		fmt.Printf("sizeWithPagination: %+v\n", sizeWithPagination)

		// Проверяем, что ошибки нет
		assert.NoError(t, err)
		// Проверяем, что user не nil
		assert.NotNil(t, sizeWithPagination)

	})

}

func ptr[T any](v T) *T {
	return &v
}
