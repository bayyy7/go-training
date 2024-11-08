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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db := database.ConnectDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get DB from GORM:", err)
	}
	defer sqlDB.Close()

	jwtKey := os.Getenv("JWTKEY")

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hahahihi",
		})
	})

	math := r.Group("/math")
	{
		math.GET("/sum", func(ctx *gin.Context) {
			a, _ := strconv.Atoi(ctx.DefaultQuery("a", "0"))
			b, _ := strconv.Atoi(ctx.DefaultQuery("a", "0"))

			ctx.JSON(http.StatusOK, gin.H{
				"message": utils.MagicSum(a, b),
			})
		})

		math.POST("/sub", handlers.MathSubHandler)
	}

	authHandler := handlers.NewAuth(db, []byte(jwtKey))
	authRoutes := r.Group("/auth")
	authRoutes.POST("/login", authHandler.AuthLogin)
	authRoutes.POST("/signup", authHandler.AuthSignUp)

	accountHandler := handlers.NewAccount(db)
	accountRoutes := r.Group("/account")
	accountRoutes.POST("/create", accountHandler.Create)
	accountRoutes.GET("/read/:id", accountHandler.Read)
	accountRoutes.PATCH("/update/:id", middleware.AuthJWTMiddleware(jwtKey), accountHandler.Update)
	accountRoutes.DELETE("/delete/:id", accountHandler.Delete)
	accountRoutes.GET("/list", accountHandler.List)
	accountRoutes.GET("/my", middleware.AuthJWTMiddleware(jwtKey), accountHandler.My)
	accountRoutes.POST("/topup/:id", accountHandler.TopUp)
	accountRoutes.GET("/balance", middleware.AuthJWTMiddleware(jwtKey), accountHandler.Balance)
	accountRoutes.POST("/transfer", middleware.AuthJWTMiddleware(jwtKey), accountHandler.Transfer)

	transactionHandler := handlers.NewTransaction(db)
	transactionRoutes := r.Group("/transaction")
	transactionRoutes.GET("/last/:id", transactionHandler.LastTransaction)
	r.Run()
}
