package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sachin141097/students-api/internal/config"
	"github.com/sachin141097/students-api/internal/http/handlers/student"
	"github.com/sachin141097/students-api/internal/storage/sqllite"
)

func main() {
	cfg := config.MustLoad()
	storage, err := sqllite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized")
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PUT /api/students/{id}", student.Update(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.Delete(storage))
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.Addr))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()
	<-done
	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
