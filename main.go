package main

import (
	"godb/handler"
	"godb/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"github.com/gin-contrib/cors"
	"gorm.io/gorm"
)

func main() {
	// Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Database
	db := NewDatabase()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get DB from GORM:", err)
	}
	defer sqlDB.Close()

	// secret-key
	signingKey := os.Getenv("SIGNING_KEY")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:8082", "http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	// grouping route with /auth
	authHandler := handler.NewAuth(db, []byte(signingKey))
	authRoute := r.Group("/auth")
	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/upsert", authHandler.Upsert)
	authRoute.POST("/change-password", middleware.AuthMiddleware(signingKey), authHandler.ChangePassword)

	accountHandler := handler.NewAccount(db)
	accountRoutes := r.Group("/account")
	accountRoutes.POST("/create", accountHandler.Create)
	accountRoutes.GET("/read/:id", accountHandler.Read)
	accountRoutes.PATCH("/update/:id", accountHandler.Update)
	accountRoutes.DELETE("/delete/:id", accountHandler.Delete)
	accountRoutes.GET("/list", accountHandler.List)
	accountRoutes.POST("/topup", accountHandler.TopUp)
	accountRoutes.GET("/balance", middleware.AuthMiddleware(signingKey), accountHandler.Balance)
	accountRoutes.GET("/my", middleware.AuthMiddleware(signingKey), accountHandler.My)
	accountRoutes.POST("/transfer", middleware.AuthMiddleware(signingKey), accountHandler.Transfer)
	accountRoutes.GET("/mutation", middleware.AuthMiddleware(signingKey), accountHandler.MutationList)

	transcatHandler := handler.NewTransactionCategories(db)
	transcatRoutes := r.Group("/transcat")
	transcatRoutes.POST("/create", transcatHandler.Create)
	transcatRoutes.GET("/read/:id", transcatHandler.Read)
	transcatRoutes.PATCH("/update/:id", transcatHandler.Update)
	transcatRoutes.DELETE("/delete/:id", transcatHandler.Delete)
	transcatRoutes.GET("/list", transcatHandler.List)

	transactionHandler := handler.NewTransactionHandler(db)
	transactionRoutes := r.Group("/transaction")
	transactionRoutes.POST("/new", transactionHandler.NewTransaction)
	transactionRoutes.GET("/list/:account_id", transactionHandler.TransactionList)

	r.Run(":8081") 
}

func NewDatabase() *gorm.DB {
	// dsn := "host=localhost port=5432 user=postgres dbname=digi sslmode=disable TimeZone=Asia/Jakarta"
	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get DB object: %v", err)
	}

	var currentDB string
	err = sqlDB.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		log.Fatalf("failed to query current database: %v", err)
	}

	log.Printf("Current Database: %s\n", currentDB)

	return db
}
