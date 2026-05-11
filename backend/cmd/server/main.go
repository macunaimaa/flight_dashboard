package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"

	"github.com/macunaimaa/dashboard/backend/internal/config"
	"github.com/macunaimaa/dashboard/backend/internal/handler"
	"github.com/macunaimaa/dashboard/backend/internal/ingestion"
	"github.com/macunaimaa/dashboard/backend/internal/middleware"
	"github.com/macunaimaa/dashboard/backend/internal/pkg/logger"
	pkgmongo "github.com/macunaimaa/dashboard/backend/internal/pkg/mongo"
	"github.com/macunaimaa/dashboard/backend/internal/repository"
	"github.com/macunaimaa/dashboard/backend/internal/service"
	"github.com/macunaimaa/dashboard/backend/internal/ws"
)

func main() {
	// Load .env from project root (../.. from cmd/server)
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	cfg := config.Load()
	l := logger.New(cfg.LogLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// MongoDB
	mongoClient, err := pkgmongo.Connect(ctx, cfg.MongoURI, l)
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB)
	if err := pkgmongo.EnsureIndexes(ctx, db, l); err != nil {
		log.Fatalf("indexes: %v", err)
	}

	// Redis
	redisClient := goredis.NewClient(&goredis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis: %v", err)
	}
	l.Info("connected to Redis", "addr", cfg.RedisAddr)
	defer redisClient.Close()

	// Repositories
	aircraftRepo := repository.NewMongoAircraftRepo(db)
	tenantRepo := repository.NewMongoTenantRepo(db)
	auditRepo := repository.NewMongoAuditRepo(db)
	cache := repository.NewRedisCache(redisClient)

	// Services
	authSvc := service.NewAuthService(tenantRepo, cfg.JWTSecret)
	aircraftSvc := service.NewAircraftService(aircraftRepo, cache)
	auditSvc := service.NewAuditService(auditRepo)
	newsSvc := service.NewNewsService(cfg.NewsAPIKey)

	// WebSocket Hub
	hub := ws.NewHub(l)
	go hub.Run()

	// Ingestion Pipeline — select data source
	var dataClient ingestion.StatesFetcher
	source := cfg.DataSource
	if cfg.UseMockData {
		source = "mock" // backwards compat
	}
	switch source {
	case "mock":
		l.Info("using MOCK data source")
		dataClient = ingestion.NewMockClient()
	case "opensky":
		l.Info("using OpenSky data source")
		dataClient = ingestion.NewOpenSkyClient(cfg.OpenSkyUsername, cfg.OpenSkyPassword)
	default:
		l.Info("using adsb.lol data source")
		dataClient = ingestion.NewADSBLolClient()
	}
	pipeline := ingestion.NewPipeline(dataClient, aircraftRepo, tenantRepo, cache, hub, l)
	go pipeline.Start(cfg.OpenSkyPollInterval)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	aircraftHandler := handler.NewAircraftHandler(aircraftSvc, auditSvc)
	wsHandler := handler.NewWSHandler(hub, authSvc, l)
	newsHandler := handler.NewNewsHandler(newsSvc)

	// Router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(l))
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// Public routes
	r.Get("/api/v1/health", handler.Health)
	r.Post("/api/v1/auth/login", authHandler.Login)

	// WebSocket — auth handled inside the handler via query param token
	r.HandleFunc("/api/v1/ws", wsHandler.ServeHTTP)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authSvc))

		r.Get("/api/v1/aircraft", aircraftHandler.List)
		r.Get("/api/v1/aircraft/bbox", aircraftHandler.GetInBBox)
		r.Get("/api/v1/aircraft/{icao24}", aircraftHandler.GetByICAO)
		r.Get("/api/v1/aircraft/{icao24}/track", aircraftHandler.GetTrack)
		r.Get("/api/v1/news/search", newsHandler.Search)
	})

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		l.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("shutting down...")
	pipeline.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}
