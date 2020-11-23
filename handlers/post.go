package handlers

import (
	"log"
	"net/http"
	"regexp"
	"samplerest/data"
	"strconv"
)

// Post handler is a struct that handles all apis regarding the Posts
type Post struct {
	l *log.Logger
}

// NewPost is used to inject logger into the Post handler
func NewPost(l *log.Logger) *Post {
	return &Post{l}
}

func (post *Post) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	post.l.Println("ServeHTTP")
	if r.Method == http.MethodGet {
		post.getPosts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		post.addPost(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		regex := regexp.MustCompile(`/([0-9]+)`)
		idData := regex.FindAllStringSubmatch(r.URL.Path, -1)

		if len(idData) != 1 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}
		if len(idData[0]) != 2 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		idString := idData[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}
		post.l.Println("got id: ", id)

		post.updatePost(id, rw, r)

	}

	// catch all others
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (post *Post) getPosts(rw http.ResponseWriter, r *http.Request) {
	post.l.Println("Handle GET Posts")
	listOfPosts := data.GetPosts()
	err := listOfPosts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal Json", http.StatusInternalServerError)
	}
}

func (post *Post) addPost(rw http.ResponseWriter, r *http.Request) {
	post.l.Println("Handle POST Post", r.Body)

	postData := &data.Post{}

	err := postData.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall json", http.StatusBadRequest)
	}

	post.l.Printf("Post: %#v", postData)

	data.AddPost(postData)
}

func (post *Post) updatePost(id int, rw http.ResponseWriter, r *http.Request) {
	post.l.Println("Handle PUT Post", r.Body)

	postData := &data.Post{}

	err := postData.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall json", http.StatusBadRequest)
		return
	}
	err = data.UpdatePost(id, postData)
	if err != nil {
		http.Error(rw, "Post not found", http.StatusNotFound)
		return
	}
}
