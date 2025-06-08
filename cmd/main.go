package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop/configs/env"
	"shop/configs/pg_conf"
	"shop/internal/api/middlewares"
	"shop/internal/api/routes"
	"shop/pkg/log"
	"syscall"
	"time"
)

func main() {
	initApp()
	app := fiber.New()

	// Middleware: CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins, consider limiting this for production
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Middleware: Global error handling
	//app.Use(middlewares.ErrorHandlerMiddleware(container.Logger))
	////Middleware: Request logging
	app.Use(middlewares.RequestLoggerMiddleware())
	app.Use(middlewares.LimitQueryParamsMiddleware)

	groupApi := app.Group("/api")

	// Register routes
	routes.RegisterSizeRoutes(groupApi)
	routes.RegisterCharacteristicRoutes(groupApi)
	routes.RegisterCharDefaultValueRoutes(groupApi)
	routes.RegisterNodeTypeRoutes(groupApi)
	routes.RegisterNodeRoutes(groupApi)
	routes.RegisterCardRoutes(groupApi)
	routes.RegisterOrderRoutes(groupApi)

	// Start the server
	port := env.GetEnv("SERV_PORT", "3000")
	log.Info("Starting server", zap.String("port", port))

	// Graceful shutdown handling in a goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Call graceful shutdown handler
	handleGracefulShutdown(app)
}

// handleGracefulShutdown handles signal-based graceful shutdown
func handleGracefulShutdown(app *fiber.App) {
	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-signalChan

	// Gracefully shutdown the server
	shutdownTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info("Shutting down server...")
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Failed to gracefully shutdown server", zap.Error(err))
	} else {
		log.Info("Server shut down gracefully")
	}
}

func initApp() {
	// Load environment variables
	env.LoadEnv()
	log.InitLogger()
	pg_conf.InitPostgresSingleton()
}
