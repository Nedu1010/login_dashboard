package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/login_flow/auth-service/internal/config"
	"github.com/login_flow/auth-service/internal/handler"
	"github.com/login_flow/auth-service/internal/middleware"
	"github.com/login_flow/auth-service/internal/repository/postgres"
	"github.com/login_flow/auth-service/internal/service"
)

func main() {

	// Database connection

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	db, err := postgres.NewDB((cfg.Database.URL))
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	defer db.Close()

	log.Println("connected to database")

	// Tables setup

	userRepo := postgres.NewUserRepository(db) // Manages "users" table
	tokenRepo := postgres.NewTokenRepository(db)

	authService := service.NewAuthService(userRepo, tokenRepo, cfg) // Login, register, token refresh
	csrfService := service.NewCSRFService(cfg.JWT.Secret)

	authHandler := handler.NewAuthHandler(authService, csrfService, cfg) // /auth/* endpoints
	userHandler := handler.NewUserHandler(authService)

	// Server setup

	gin.SetMode(gin.ReleaseMode) // Production mode (less verbose logging)
	// Create new router (no default middleware)

	r := gin.New()

	// Middleware

	r.Use(gin.Recovery())                             // Catches panics and returns 500 instead of crashing
	r.Use(middleware.Logger())                        // Logs each request (method, path, status, duration)
	r.Use(middleware.CORS(cfg.Server.AllowedOrigins)) // Allows cross-origin requests from frontend

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "running",
		})
	})
	api := r.Group("/api")
	{
		// Authentication routes (PUBLIC - no authentication required)
		auth := api.Group("/auth") // All routes here start with /api/auth
		{
			auth.POST("/register", authHandler.Register)                          // POST /api/auth/register
			auth.POST("/login", authHandler.Login)                                // POST /api/auth/login
			auth.POST("/refresh", authHandler.Refresh)                            // POST /api/auth/refresh
			auth.POST("/logout", middleware.CSRFMiddleware(), authHandler.Logout) // POST /api/auth/logout (CSRF protected)
		}

		// User routes (PROTECTED - require valid access token)
		user := api.Group("/user") // All routes here start with /api/user
		// Use() adds middleware to this group only
		// AuthMiddleware checks for valid access token in cookies
		user.Use(middleware.AuthMiddleware(authService))
		{
			user.GET("/me", userHandler.GetMe) // GET /api/user/me (requires auth)
		}
	}

	r.Run(":8080")
}
