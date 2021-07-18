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

	db "github.com/pyoko/pyogo-rest-template/pkg/models"
	"github.com/pyoko/pyogo-rest-template/pkg/routers"
	"github.com/pyoko/pyogo-rest-template/pkg/settings"
)

func main() {
	// Logger presets
	log.Println("server is starting")

	// Database
	database, err := db.DbConnect()
	if err != nil {
		panic(fmt.Sprintf("server is unable to open the DB connection: %+v", err))
	}

	// Listen for OS interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		OSCall := <-c // waiting for a signal
		log.Println(fmt.Sprintf("server received a system call: %+v", OSCall))
		cancel()
	}()

	if err := serve(ctx, database); err != nil {
		panic(fmt.Sprintf("server failed to serve: %+v", err))
	}
}

func serve(ctx context.Context, db *db.DB) (err error) {
	router := routers.NewRouter(db)
	appRouter := router.Init()

	server := &http.Server{
		Addr: ":"+settings.ReadConfig("SERVER_PORT"),
		Handler: appRouter,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("server critical: %+v", err))
		}
	}()


	log.Println("server started")
	<-ctx.Done()
	log.Println("server  is shutting down")

	// Max time allowed for the in-flight requests to complete
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
		if err == http.ErrServerClosed {
			err = nil
		}
	}()

	// Shutdown HTTP server
	if err = server.Shutdown(ctxShutDown); err != nil {
		log.Println(fmt.Sprintf("server shutdown forcedly after %d seconds of waiting", 10))
	}

	log.Println("server exited gracefully")

	return
}