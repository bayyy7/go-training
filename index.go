package main

import (
	"example/database"
	"example/handlers"
	"example/middleware"
	"example/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         time.Duration
}

func getDefaultConfig() Config {
	return Config{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		MaxAge: 12 * time.Hour,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	db := database.ConnectDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get DB from GORM:", err)
	}
	defer sqlDB.Close()

	jwtKey := os.Getenv("JWTKEY")
	if jwtKey == "" {
		log.Fatal("JWTKEY environment variable is not set")
	}

	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     getDefaultConfig().AllowedOrigins,
		AllowMethods:     getDefaultConfig().AllowedMethods,
		AllowHeaders:     getDefaultConfig().AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           getDefaultConfig().MaxAge,
	}

	r.Use(cors.New(corsConfig))

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the API",
			"version": "1.0",
		})
	})

	math := r.Group("/math")
	{
		math.GET("/sum", func(ctx *gin.Context) {
			a, err1 := strconv.Atoi(ctx.DefaultQuery("a", "0"))
			b, err2 := strconv.Atoi(ctx.DefaultQuery("b", "0")) // Fixed: changed "a" to "b"

			if err1 != nil || err2 != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid parameters. Both 'a' and 'b' must be integers",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"result": utils.MagicSum(a, b),
			})
		})
		math.POST("/sub", handlers.MathSubHandler)
	}

	authHandler := handlers.NewAuth(db, []byte(jwtKey))
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.AuthLogin)
		authRoutes.POST("/signup", authHandler.AuthSignUp)
		authRoutes.POST("/upsert", authHandler.Upsert)
	}

	accountHandler := handlers.NewAccount(db)
	accountRoutes := r.Group("/account")
	{
		accountRoutes.POST("/create", accountHandler.Create)
		accountRoutes.GET("/read/:id", accountHandler.Read)
		accountRoutes.PATCH("/update/:id", middleware.AuthJWTMiddleware(jwtKey), accountHandler.Update)
		accountRoutes.DELETE("/delete/:id", accountHandler.Delete)
		accountRoutes.GET("/list", accountHandler.List)
		accountRoutes.GET("/my", middleware.AuthJWTMiddleware(jwtKey), accountHandler.My)
		accountRoutes.POST("/topup/:id", accountHandler.TopUp)
		accountRoutes.GET("/balance", middleware.AuthJWTMiddleware(jwtKey), accountHandler.Balance)
		accountRoutes.POST("/transfer", middleware.AuthJWTMiddleware(jwtKey), accountHandler.Transfer)
	}

	transactionHandler := handlers.NewTransaction(db)
	transactionRoutes := r.Group("/transaction")
	{
		transactionRoutes.GET("/last/:id", transactionHandler.LastTransaction)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
