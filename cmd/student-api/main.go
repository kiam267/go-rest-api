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

	"github.com/kiam267/student-api/internal/config"
	"github.com/kiam267/student-api/internal/http/handlers/student"
)

func main() {
	//  Load Config
	cfg := config.MustLoad()

	// database setup 
	// setup router
	router:= http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	// setup server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started", slog.String("address",cfg.Addr,))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)


  go func ()  {
		err:= server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done
  slog.Info("Shutting down the server")
	
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
 
 if err:= server.Shutdown(ctx);err != nil {
	slog.Error("Failed to Shutdown server", slog.String("error", err.Error()))
 }
 slog.Info("server Shutdown successfully");

}