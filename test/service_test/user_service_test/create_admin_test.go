package article_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"shop/configs/env"
	"shop/configs/pg_conf"
	"shop/internal/api/dto"
	"shop/internal/repository"
	"shop/internal/service"
	"shop/pkg/log"
	"testing"
)

func TestMain(m *testing.M) {
	// Запуск тестов
	os.Exit(m.Run())
}

func setup() (context.Context, service.UserServiceInterface) {
	err := os.Setenv("MONGO_URI", "mongodb://boostbiz1:1000@localhost:27017/?authSource=admin")
	if err != nil {
		fmt.Println("Failed to set env MONGO_URI")
		return nil, nil
	}

	err = os.Setenv("MONGO_DB_NAME", "test")
	if err != nil {
		fmt.Println("Failed to set env MONGO_DB_NAME")
		return nil, nil
	}

	// Загружаем переменные окружения
	env.LoadEnv()

	// Инициализируем логгер и клиент MongoDB
	log.InitLogger()
	pg_conf.InitMongoSingleton()
	logger := log.GetLogger()
	clientDB := pg_conf.GetClient()

	ctx := context.Background()

	repo := repository.NewUserRepository(clientDB, logger)
	serv := service.NewUserService(repo, logger)
	return ctx, serv
}

func TestCreateAdmin(t *testing.T) {
	ctx, serv := setup()
	// Подготовка тестовых данных
	correctDto := dto.CreateUserRequest{
		Email:    "alex@gmail.com",
		Phone:    "375256868812",
		Password: "1234",
	}

	// Запуск теста
	t.Run("Success", func(t *testing.T) {
		createdUser, err := serv.CreateNewAdmin(ctx, &correctDto)
		// Вывод структуры с ключами
		fmt.Printf("CreatedUser: %+v\n", createdUser)

		// Проверяем, что ошибки нет
		assert.NoError(t, err)
		// Проверяем, что user не nil
		assert.NotNil(t, createdUser)

		// Проверяем, что структура содержит нужные ключи
		assert.Equal(t, "alex@gmail.com", createdUser.Email)
		assert.Equal(t, "375256868812", createdUser.Phone)
	})

	// Запуск теста
	//t.Run("Not found", func(t *testing.T) {
	//	article, err := repo.GetArticleById(ctx, invalidID)
	//	// Вывод структуры с ключами
	//	fmt.Printf("Article: %+v\n", article)
	//	fmt.Printf("err: %+v\n", err)
	//	// Проверяем, что вернулась ошибка mongo.ErrNoDocuments
	//	assert.ErrorIs(t, err, m.ErrNoDocuments)
	//})
}

func ptr[T any](v T) *T {
	return &v
}
