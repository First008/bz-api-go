package api

import (
	"bulutzincir/auth"

	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zemirco/couchdb"
)

func Gin() {
	dbUrl := os.Getenv("DB_URL")
	couchDbUrl, err := url.Parse(dbUrl)

	if err != nil {
		panic(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbUserPassword := os.Getenv("DB_USER_PW")

	client, err := couchdb.NewAuthClient(dbUser, dbUserPassword, couchDbUrl)
	if err != nil {
		panic(err)
	}

	port := ":" + os.Getenv("PORT")
	fmt.Println(port)

	var db = auth.NewAuth(client)
	var tk = auth.NewToken()
	var service = auth.NewProfile(db, tk)

	var router = gin.Default()

	router.GET("/login", service.Login)
	router.POST("/todo", auth.TokenAuthMiddleware(), service.CreateTodo)
	router.POST("/register", service.CreateAccount)
	router.POST("/logout", auth.TokenAuthMiddleware(), service.Logout)
	router.GET("/whoami", auth.TokenAuthMiddleware(), service.ReturnIdentity)

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
