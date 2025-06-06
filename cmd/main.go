package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IgorBrizack/rate-limiter-system-design/internal/controller"
	"github.com/IgorBrizack/rate-limiter-system-design/internal/database"
	"github.com/IgorBrizack/rate-limiter-system-design/internal/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env.")
	}

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}

	cacheDB := database.NewRedisClient()
	database := database.NewDatabase()
	DB := database.DB()

	userController := controller.NewController(cacheDB, DB)

	router := gin.Default()

	rateLimiter := limiter.NewMiddleware(cacheDB)

	router.GET("/users", rateLimiter.LeakingBucketHandler(1, 1.0), userController.GetUsers)
	router.POST("/users", rateLimiter.LeakingBucketHandler(1, 1.0), userController.CreateUser)

	fmt.Printf("Running on %s\n", port)
	router.Run(":" + port)
}
