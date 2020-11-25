package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"samplerest/handlers"
	"time"

	"github.com/gorilla/mux"
)

//User is a struct that represents a User
type User struct {
	FullName string `json:"fullName"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

//Post is a struct that represents a post
type Post struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author User   `json:"author"`
}

var posts []Post = []Post{}

func main() {

	l := log.New(os.Stdout, "Posts-api", log.LstdFlags)
	handlerPosts := handlers.NewPost(l)

	serveMux := mux.NewRouter()
	serveMux.Use(commonMiddleware)

	getRouter := serveMux.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", handlerPosts.GetPosts)

	putRouter := serveMux.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", handlerPosts.UpdatePost)
	putRouter.Use(handlerPosts.MiddlewarePostValidation)

	postRouter := serveMux.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", handlerPosts.AddPost)
	postRouter.Use(handlerPosts.MiddlewarePostValidation)

	delRouter := serveMux.Methods("DELETE").Subrouter()
	delRouter.HandleFunc("/{id:[0-9]+}", handlerPosts.DeletePost)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	l.Println("Received terminate, graceful shutdown", sig)

	timeOutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeOutContext)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
