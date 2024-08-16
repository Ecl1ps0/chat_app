package main

import (
	"ChatApp/internal/auth"
	"ChatApp/internal/auth/repository"
	"ChatApp/internal/auth/usercases"
	"ChatApp/pkg/mongo"
	server2 "ChatApp/pkg/server"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env: %s", err.Error())
	}

	db, err := mongo.InitDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatalf("fail to connect database: %s", err.Error())
	}

	router := http.NewServeMux()

	authRepo := repository.NewAuthRepository(db.Database("ChatApp"), "users")
	authHandler := auth.NewAuthHandler(usercases.NewAuthUsecase(*authRepo))
	authHandler.AuthRouterInit(router)

	server := new(server2.Server)
	go func() {
		if err := server.Run("8080", router); err != nil {
			log.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	log.Println("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Server shutting down")

	if err = server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err = db.Disconnect(context.TODO()); err != nil {
		log.Fatalf("fail to disconnect with database: %s", err.Error())
	}
}
