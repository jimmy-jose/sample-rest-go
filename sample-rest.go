package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"samplerest/handlers"
	"strconv"
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

	getRouter := serveMux.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", handlerPosts.GetPosts)

	putRouter := serveMux.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", handlerPosts.UpdatePost)
	putRouter.Use(handlerPosts.MiddlewarePostValidation)

	postRouter := serveMux.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", handlerPosts.AddPost)
	postRouter.Use(handlerPosts.MiddlewarePostValidation)

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

func getPosts(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// handling the error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to an integer!"))
		return
	}

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found!"))
		return
	}

	post := posts[id]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

func addPost(w http.ResponseWriter, r *http.Request) {
	//get item from request body
	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)
	posts = append(posts, newPost)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// handling the error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to an integer!"))
		return
	}

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found!"))
		return
	}

	var updatedPost Post
	json.NewDecoder(r.Body).Decode(&updatedPost)

	posts[id] = updatedPost

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

func patchPost(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// handling the error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to an integer!"))
		return
	}

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found!"))
		return
	}

	post := &posts[id]
	json.NewDecoder(r.Body).Decode(post)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

func deletePost(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// handling the error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to an integer!"))
		return
	}

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found!"))
		return
	}

	posts = append(posts[:id], posts[id+1:]...)

	w.WriteHeader(200)

}
