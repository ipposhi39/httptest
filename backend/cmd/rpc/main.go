package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/httptest/backend/pkg/config"
	"github.com/httptest/backend/pkg/logger"
	"github.com/httptest/backend/pkg/util"
	"github.com/httptest/backend/rpc/injector"
	"github.com/inconshreveable/log15"
)

var (
	applicationName = "httptest-rpc"
)

func init() {
	util.InitLocale()
}

func main() {
	c := config.Prepare()
	logger.InitLogger(c.Logger)

	log15.Info("listening....", "method", "main.init", "port", c.HTTP.Port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.HTTP.Port),
		Handler: mux(c),
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalln("Server closed with error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log15.Error("Failed to gracefully shutdown:", err)
	}
	log15.Info("Server shutdown")
}

func mux(c config.AppConfig) *http.ServeMux {
	firebaseHandler := injector.InitializeFirebaseHandler(c.HTTP, c.Postgres, c.Firebase, "wdc-rpc-firebase")

	mux := http.NewServeMux()
	mux.Handle("/", firebaseHandler)
	return mux
}
