package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

//User is a struct that represents a User
type User struct {
	FullName  string `json:"fullName"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

//Post is a struct that represents a post
type Post struct {
	Title     string `json:"title" validate:"required"`
	Body      string `json:"body"`
	Author    User   `json:"author"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// Posts is a custom type which has an inbuilt toJson Method
type Posts []*Post

// Validate is used to validate the struct
func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ToJSON encodes the posts slice (p) into json and writes it to the io writer
func (p *Posts) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

// FromJSON decodes the JSON from the io.Reader
func (p *Post) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

// GetPosts returns all posts
func GetPosts() Posts {
	return postList
}

// AddPost takes a post object add add it to the postList
func AddPost(post *Post) {
	postList = append(postList, post)
}

// UpdatePost will update the post and sends an error back if post is not found
func UpdatePost(id int, post *Post) error {
	if isPresent(id) {
		postList[id] = post
		return nil
	}
	return ErrorPostNotFound
}

// DeletePost delete the post at index id. Returns Error in case post cannot be found
func DeletePost(id int) error {
	if isPresent(id) {
		postList = append(postList[:id], postList[id+1:]...)
		return nil
	}
	return ErrorPostNotFound
}

// ErrorPostNotFound is an error when post is not found in the list
var ErrorPostNotFound = fmt.Errorf("Product not found")

func isPresent(id int) bool {
	if id >= len(postList) {
		return false
	}
	return true
}

var postList = []*Post{
	{
		Title: "ONE",
		Body:  "One Body",
		Author: User{
			FullName:  "One Author",
			UserName:  "oneuser",
			Email:     "one@gmail.com",
			CreatedOn: time.Now().UTC().String(),
			UpdatedOn: time.Now().UTC().String(),
			DeletedOn: time.Now().UTC().String(),
		},
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
		DeletedOn: time.Now().UTC().String(),
	},
	{
		Title: "Two",
		Body:  "Two Body",
		Author: User{
			FullName:  "Two Author",
			UserName:  "twouser",
			Email:     "two@gmail.com",
			CreatedOn: time.Now().UTC().String(),
			UpdatedOn: time.Now().UTC().String(),
			DeletedOn: time.Now().UTC().String(),
		},
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
		DeletedOn: time.Now().UTC().String(),
	},
}
