package main

import (
	"fmt"
	"log"
	"time"

	"coworking/space-service/internal/config"
	"coworking/space-service/internal/handlers"
	"coworking/space-service/internal/middleware"
	"coworking/space-service/internal/models"
	"coworking/space-service/internal/repository"
	"coworking/space-service/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "coworking/space-service/docs"
)

// @title Space Service API
// @version 1.0
// @description Coworking Space Service API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// ── Configuration ────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Database ─────────────────────────────────────────────────────────────
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)

	var db *gorm.DB
	var err error

	// Retry loop: wait up to 30 s for PostgreSQL to become ready.
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				if pingErr = sqlDB.Ping(); pingErr == nil {
					log.Println("[db] Connected to PostgreSQL successfully")
					break
				}
			}
			err = pingErr
		}
		log.Printf("[db] Attempt %d/10 – waiting for PostgreSQL: %v", attempt, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("[db] Could not connect to PostgreSQL after 10 attempts: %v", err)
	}

	// Configure connection pool.
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// ── Auto-migration ───────────────────────────────────────────────────────
	if err := db.AutoMigrate(&models.Space{}, &models.Invoice{}); err != nil {
		log.Fatalf("[migrate] AutoMigrate failed: %v", err)
	}
	log.Println("[migrate] Schema is up to date")

	// ── Wire dependencies ────────────────────────────────────────────────────
	spaceRepo := repository.NewSpaceRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)

	spaceSvc := services.NewSpaceService(spaceRepo)
	billingSvc := services.NewBillingService(invoiceRepo)
	reportSvc := services.NewReportService(db)

	spaceHandler := handlers.NewSpaceHandler(spaceSvc)
	billingHandler := handlers.NewBillingHandler(billingSvc)
	reportHandler := handlers.NewReportHandler(reportSvc)

	// ── Router ───────────────────────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health-check (unauthenticated)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true, "message": "space-service is healthy"})
	})

	auth := middleware.AuthMiddleware(cfg.JWTSecret)
	admin := middleware.AdminMiddleware()

	// ── Space routes ─────────────────────────────────────────────────────────
	spaces := r.Group("/spaces")
	spaces.Use(auth)
	{
		// NOTE: /available must be registered BEFORE /:id to avoid route conflict.
		spaces.GET("/available", spaceHandler.GetAvailableSpaces)
		spaces.GET("", spaceHandler.ListSpaces)
		spaces.GET("/:id", spaceHandler.GetSpace)

		// Admin-only mutations
		spaces.POST("", admin, spaceHandler.CreateSpace)
		spaces.PUT("/:id", admin, spaceHandler.UpdateSpace)
		spaces.DELETE("/:id", admin, spaceHandler.DeleteSpace)
	}

	// ── Billing routes ───────────────────────────────────────────────────────
	billing := r.Group("/billing")
	billing.Use(auth)
	{
		billing.GET("/mine", billingHandler.GetMyInvoices)
		billing.GET("/:id", billingHandler.GetInvoice)

		// Admin-only
		billing.GET("", admin, billingHandler.ListInvoices)
		billing.POST("", billingHandler.CreateInvoice) // internal use (auth required)
		billing.PATCH("/:id/pay", admin, billingHandler.MarkAsPaid)
		billing.PATCH("/:id/cancel", admin, billingHandler.CancelInvoice)
	}

	// ── Report routes ────────────────────────────────────────────────────────
	reports := r.Group("/reports")
	reports.Use(auth, admin)
	{
		reports.GET("/occupancy", reportHandler.OccupancyReport)
		reports.GET("/revenue", reportHandler.RevenueReport)
		reports.GET("/spaces", reportHandler.SpacesReport)
	}

	// ── Start server ─────────────────────────────────────────────────────────
	addr := ":" + cfg.Port
	log.Printf("[server] space-service listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[server] Failed to start: %v", err)
	}
}
