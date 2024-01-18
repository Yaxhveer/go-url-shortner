package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
    
	"github.com/Yaxhveer/go-url-shortner/handler"
	"github.com/Yaxhveer/go-url-shortner/services"
	"github.com/Yaxhveer/go-url-shortner/storage"
	"github.com/joho/godotenv"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
}

func main() {

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
    defer stop()

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    PORT := os.Getenv("PORT")
    log.Println(PORT)

    db, err := storage.NewPostgresStore(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Postgres Database created")

    err = db.CreateSchema()
    if err != nil {
        log.Fatal(err)
    }

    service := services.NewService(db)
    log.Println("Service created")

    hdlr := handler.NewHandler(service)
    
    server := handler.NewServer(ctx, hdlr, PORT)

    err = server.Run(ctx)
    if err != nil {
        log.Fatal(err)
    }
}