package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project/kocokan/config"
	"github.com/project/kocokan/internal/handler"
	"github.com/project/kocokan/internal/repository"
	"github.com/project/kocokan/internal/service"
	"github.com/project/kocokan/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gormLog := logger.Info
	if cfg.AppEnv == "production" {
		gormLog = logger.Silent
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(gormLog),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	roundRepo := repository.NewRoundRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	groupSvc := service.NewGroupService(groupRepo, participantRepo, roundRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	groupHandler := handler.NewGroupHandler(groupSvc)

	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "./web/static")
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20)
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// Pages
	r.GET("/", func(c *gin.Context) { c.HTML(200, "public.html", nil) })
	r.GET("/app", func(c *gin.Context) { c.HTML(200, "app.html", nil) })

	// Auth API
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/logout", authHandler.Logout)

	// Protected API
	api := r.Group("/api", middleware.Auth(authSvc))
	{
		api.GET("/groups", groupHandler.List)
		api.POST("/groups", groupHandler.Create)
		api.GET("/groups/:id", groupHandler.Get)
		api.PUT("/groups/:id", groupHandler.Update)
		api.DELETE("/groups/:id", groupHandler.Delete)

		api.POST("/groups/:id/participants", groupHandler.AddParticipant)
		api.PUT("/groups/:id/participants/:pid", groupHandler.UpdateParticipant)
		api.DELETE("/groups/:id/participants/:pid", groupHandler.DeleteParticipant)

		api.POST("/groups/:id/draw", groupHandler.Draw)
		api.PUT("/groups/:id/rounds/:rid/winner", groupHandler.UpdateWinner)
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("kocokan running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
