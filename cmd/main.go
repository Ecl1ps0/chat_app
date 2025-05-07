package main

import (
	"ChatApp/internal/audio"
	repository6 "ChatApp/internal/audio/repository"
	usecase6 "ChatApp/internal/audio/usecase"
	"ChatApp/internal/auth"
	"ChatApp/internal/auth/repository"
	"ChatApp/internal/auth/usecase"
	"ChatApp/internal/chat"
	repository2 "ChatApp/internal/chat/repository"
	usecase2 "ChatApp/internal/chat/usecase"
	"ChatApp/internal/image"
	repository5 "ChatApp/internal/image/repository"
	usecase5 "ChatApp/internal/image/usecase"
	"ChatApp/internal/message"
	repository3 "ChatApp/internal/message/repository"
	usecase3 "ChatApp/internal/message/usecase"
	"ChatApp/internal/user"
	repository4 "ChatApp/internal/user/repository"
	usecase4 "ChatApp/internal/user/usecase"
	"ChatApp/pkg/mongo"
	server2 "ChatApp/pkg/server"
	"context"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := mongo.InitDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatalf("fail to connect database: %s", err.Error())
	}

	database := db.Database("ChatApp")
	router := http.NewServeMux()

	authRepo := repository.NewAuthRepository(database, "users")
	chatRepo := repository2.NewChatRepository(database, "chats")
	messageRepo := repository3.NewMessageRepository(database, "messages")
	userRepo := repository4.NewUserRepository(database, "users")
	imageRepo := repository5.NewImageRepository(database, "images")
	audioRepo := repository6.NewAudioRepository(database, "audios")

	authUc := usecase.NewAuthUsecase(authRepo)
	chatUc := usecase2.NewChatUsecase(chatRepo)
	messageUc := usecase3.NewMessageUsecase(messageRepo)
	imageUc := usecase5.NewImageUsecase(imageRepo)
	userUc := usecase4.NewUserUsecase(userRepo)
	audioUc := usecase6.NewAudioUsecase(audioRepo)

	authHandler := auth.NewAuthHandler(authUc)
	authHandler.AuthRouterInit(router)

	chatHandler := chat.NewChatHandler(chatUc, messageUc)
	chatHandler.ChatRouterInit(router, authHandler.AuthMiddleware)

	messageHandler := message.NewMessageHandler(messageUc)
	messageHandler.MessageRouterInit(router, authHandler.AuthMiddleware)

	userHandler := user.NewUserHandler(userUc, imageUc)
	userHandler.UserRouterInit(router, authHandler.AuthMiddleware)

	imageHandler := image.NewImageHandler(usecase5.NewImageUsecase(imageRepo))
	imageHandler.ImageRouterInit(router, authHandler.AuthMiddleware)

	audioHandler := audio.NewAudioHandler(audioUc)
	audioHandler.AudioRouterInit(router, authHandler.AuthMiddleware)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://16.171.65.117", "http://13.49.21.134", "https://chat-app-front-sigma.vercel.app"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(router)

	server := new(server2.Server)
	go func() {
		if err := server.Run("8080", handler); err != nil {
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
