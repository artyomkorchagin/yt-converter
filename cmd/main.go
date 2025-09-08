package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/artyomkorchagin/yt-converter/config"
	"github.com/artyomkorchagin/yt-converter/internal/logger"
	"github.com/artyomkorchagin/yt-converter/internal/router"
)

func main() {
	var err error

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	var zapLogger *zap.Logger

	if cfg.LogLevel == "DEV" {
		zapLogger, err = logger.NewDevelopmentLogger()
	} else {
		zapLogger, err = logger.NewLogger()
	}

	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting application")

	handler := router.NewHandler(zapLogger)
	r := handler.InitRouter()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}

	port := cfg.Port
	srv := &http.Server{
		Addr:      cfg.Host + ":" + port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	go func() {
		zapLogger.Info("Server starting", zap.String("port", port))
		if err := srv.ListenAndServeTLS("./cert.pem", "./privkey.pem"); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server shutdown failed", zap.Error(err))
	}

	zapLogger.Info("Server exited")
	zapLogger.Info("Shutdown completed")
}
