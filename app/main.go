package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	_mailHandler "spektr-email-api/mail/delivery/http"
	"spektr-email-api/mail/delivery/rabbit"
	_mailRepo "spektr-email-api/mail/repository/email"
	_mailUsecase "spektr-email-api/mail/usecase"
	"spektr-email-api/pkg/rabbitmq"
	"syscall"
	"time"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	g := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://www.969975-cv27771.tmweb.ru:3000"} // Replace with your allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}                          // Specify allowed HTTP methods
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}                         // Specify allowed headers
	g.Use(cors.New(config))
	conn, err := rabbitmq.NewRabbitMQConn()
	if err != nil {
		return
	}
	defer conn.Close()
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	repo := _mailRepo.NewMailRepository(viper.GetString("email.name"), viper.GetString("email.nameFrom"), os.Getenv("pass"))
	ucase := _mailUsecase.NewMailusecase(repo, timeoutContext)
	emailPublisher := rabbit.NewEmailsPublisher(conn, ucase)
	go func() {
		err := emailPublisher.StartConsumer(
			10,
			"emails-exchange",
			"emails-queue",
			"emails-key",
			"emails-consumer",
		)
		if err != nil {
			return
		}
	}()
	_mailHandler.NewMailHandler(g, ucase)
	server := &http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: g,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for a termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
