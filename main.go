package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wizzyszn/go_bank/config"
	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/handlers"
	"github.com/wizzyszn/go_bank/middleware"
	"github.com/wizzyszn/go_bank/repository"
	"github.com/wizzyszn/go_bank/service"
)

func main() {
	log.Println("Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	log.Println("Connecting to database...")

	dbConfig := db.NewConfig(cfg.GetDNS())
	database, err := db.New(dbConfig)

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer database.Close()
	log.Println("Database connected successfully")

	//Initializing Repositories
	log.Println("Initializing repositories...")

	accountRepo := repository.NewAccountRepository(database)
	transactionRepo := repository.NewTransactionRepository(database)
	sessionRepo := repository.NewSessionRepository(database)

	// Initializing Services
	authService := service.NewAuthService(database, accountRepo, sessionRepo, cfg.Security.SessionDuration)
	transactionService := service.NewTransactionService(database, accountRepo, transactionRepo)

	// Initializing Handlers
	log.Println("Initializing Handlers...")

	authHandler := handlers.NewAuthHandler(authService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	accountHandler := handlers.NewAccountHandler(authService, transactionService)
	healthHandler := handlers.NewHealthHandler(database)

	// Initializing middlewares
	log.Println("Initializing middlewares...")

	authMiddleware := middleware.NewAuthMiddleware(authService)
	rateLimtiter := middleware.NewRateLimiter(30, 100)

	var corsConfig middleware.CORSConfig

	if cfg.IsDevelopment() {
		corsConfig = middleware.DefaultCorsConfig()
	} else {
		corsConfig = middleware.ProductionCORSConfig([]string{
			"https://willdecide.com",
		})
	}

	mux := http.NewServeMux()

	//HEALTH ENDPOINTS
	mux.HandleFunc("/health", middleware.Chain(healthHandler.Health, middleware.Logger))

	mux.HandleFunc("/ready", middleware.Chain(healthHandler.Ready, middleware.Logger))

	mux.HandleFunc("/live", middleware.Chain(healthHandler.Live, middleware.Logger))

	//PUBLIC AUTHENTICATION ENDPOINTS
	mux.HandleFunc("/api/register", middleware.Chain(authHandler.Register, middleware.Logger, middleware.CORS(corsConfig), rateLimtiter.RateLimit))

	mux.HandleFunc("/api/login", middleware.Chain(authHandler.Login, middleware.Logger, middleware.CORS(corsConfig), rateLimtiter.RateLimit))

	//PROTECTED AUTHENTICATION ENDPOINTS
	mux.HandleFunc("/api/logout", middleware.Chain(authHandler.Logout, middleware.Logger, middleware.CORS(corsConfig), authMiddleware.Authenticate))

	mux.HandleFunc("/api/me", middleware.Chain(authHandler.GetMe, middleware.Logger, middleware.CORS(corsConfig), authMiddleware.Authenticate))

	//PROTECTED ACCOUNT ENDPOINTS
	mux.HandleFunc("/api/account", func(w http.ResponseWriter, r *http.Request) {
		handler := middleware.Chain(accountHandler.GetAccount, middleware.Logger, middleware.CORS(corsConfig), authMiddleware.Authenticate)

		switch r.Method {
		case http.MethodGet:
			handler(w, r)
		case http.MethodPatch:
			middleware.Chain(accountHandler.UpdateAccount, middleware.Logger, middleware.CORS(corsConfig), authMiddleware.Authenticate)(w, r)
		case http.MethodOptions:
			middleware.CORS(corsConfig)(handler)(w, r)
		default:
			http.Error(w, r.Method+" Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/account/balance", middleware.Chain(
		accountHandler.GetBalance,
		middleware.Logger,
		middleware.CORS(corsConfig),
		authMiddleware.Authenticate,
	))
	mux.HandleFunc("/api/account/balance", middleware.Chain(
		accountHandler.GetBalance,
		middleware.Logger,
		middleware.CORS(corsConfig),
		authMiddleware.Authenticate,
	))

	// PROTECTED TRANSACTION ENDPOIINTS
	mux.HandleFunc("/api/deposit", middleware.Chain(
		transactionHandler.Deposit,
		middleware.Logger,
		middleware.CORS(corsConfig),
		authMiddleware.Authenticate,
		rateLimtiter.RateLimit,
	))
	mux.HandleFunc("/api/withdraw", middleware.Chain(
		transactionHandler.Withdraw,
		middleware.Logger,
		middleware.CORS(corsConfig),
		authMiddleware.Authenticate,
		rateLimtiter.RateLimit,
	))
	mux.HandleFunc("/api/transfer", middleware.Chain(
		transactionHandler.Transfer,
		middleware.Logger,
		middleware.CORS(corsConfig),
		authMiddleware.Authenticate,
		rateLimtiter.RateLimit,
	))
	mux.HandleFunc("/api/transactions", middleware.Chain(
		transactionHandler.GetTransations,
		middleware.Logger,
		middleware.CORS(corsConfig),
		authMiddleware.Authenticate,
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				count, err := authService.CleanupExpiredSessions()
				if err != nil {
					log.Printf("Error cleaning up sessions: %v", err)
				} else {
					log.Printf("Cleaned up %d expired sessions", count)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	server := http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		log.Printf("Environment: %s", cfg.Server.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Giving outstanding requests 30 seconds to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")

}
