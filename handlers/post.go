package handlers

import (
	"context"
	"log"
	"net/http"
	"samplerest/data"
	"strconv"

	"github.com/gorilla/mux"
)

// Post handler is a struct that handles all apis regarding the Posts
type Post struct {
	l *log.Logger
}

// NewPost is used to inject logger into the Post handler
func NewPost(l *log.Logger) *Post {
	return &Post{l}
}

// GetPosts returns all the posts
func (post *Post) GetPosts(rw http.ResponseWriter, r *http.Request) {

	post.l.Println("Handle GET Posts")
	listOfPosts := data.GetPosts()
	err := listOfPosts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal Json", http.StatusInternalServerError)
	}
}

// AddPost adds a post object to the list
func (post *Post) AddPost(rw http.ResponseWriter, r *http.Request) {
	post.l.Println("Handle POST Post")

	postData := r.Context().Value(KeyPost{}).(data.Post)
	data.AddPost(&postData)
}

// UpdatePost is used to update a post
func (post Post) UpdatePost(rw http.ResponseWriter, r *http.Request) {
	post.l.Println("Handle PUT Post", r.Body)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to parse id", http.StatusBadRequest)
		return
	}

	postData := r.Context().Value(KeyPost{}).(data.Post)

	err = data.UpdatePost(id, &postData)
	if err != nil {
		http.Error(rw, "Post not found", http.StatusNotFound)
		return
	}
}

type KeyPost struct{}

// MiddlewarePostValidation handles the post validations
func (post Post) MiddlewarePostValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		postData := data.Post{}
		err := postData.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshall json", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyPost{}, postData)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
