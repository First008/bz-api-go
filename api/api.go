package api

import (
	"bulutzincir/cmd"
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

	u, err := url.Parse("http://127.0.0.1:5984/")
	if err != nil {
		panic(err)
	}

	client, err := couchdb.NewAuthClient("admin", "qwer", u)
	if err != nil {
		panic(err)
	}

	port := ":" + os.Getenv("PORT")
	fmt.Println(port)

	var db = cmd.NewAuth(client)
	var tk = cmd.NewToken()
	var service = cmd.NewProfile(db, tk)

	var router = gin.Default()

	router.GET("/login", service.Login)
	router.POST("/todo", cmd.TokenAuthMiddleware(), service.CreateTodo)
	router.POST("/register", service.CreateAccount)
	router.POST("/logout", cmd.TokenAuthMiddleware(), service.Logout)
	router.GET("/whoami", cmd.TokenAuthMiddleware(), service.ReturnIdentity)
	// router.POST("/refresh", service.Refresh)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
