package data

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//User is a struct that represents a User
type User struct {
	FullName string `json:"fullName"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

//Post is a struct that represents a post
type Post struct {
	gorm.Model
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body"`
	Author User   `json:"author" gorm:"embedded;embeddedPrefix:author_"`
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

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:@tcp(localhost:3306)/posts?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                                                                     // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                                                                    // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                    // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                    // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                   // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		fmt.Print(err)
	}

	var postList = []*Post{}
	db.Find(&postList)

	return postList
}

// AddPost takes a post object add add it to the postList
func AddPost(post *Post) {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:@tcp(localhost:3306)/posts?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                                                                     // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                                                                    // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                    // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                    // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                   // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		fmt.Print(err)
	}

	db.Create(&post)
	fmt.Println(post)

}

// UpdatePost will update the post and sends an error back if post is not found
func UpdatePost(id int, post *Post) error {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:@tcp(localhost:3306)/posts?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                                                                     // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                                                                    // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                    // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                    // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                   // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		fmt.Print(err)
	}

	var tempData = &Post{}
	dataFromDb := db.First(&tempData, id)

	if dataFromDb.RowsAffected <= 0 {
		return ErrorPostNotFound
	}

	tempData.Title = post.Title
	tempData.Body = post.Body
	tempData.Author = post.Author

	db.Save(&tempData)
	return nil
}

// DeletePost delete the post at index id. Returns Error in case post cannot be found
func DeletePost(id int) error {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:@tcp(localhost:3306)/posts?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                                                                     // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                                                                    // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                    // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                    // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                   // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		fmt.Print(err)
	}

	var post = &Post{}
	dataFromDb := db.First(&post, id)

	if dataFromDb.RowsAffected <= 0 {
		return ErrorPostNotFound
	}

	db.Delete(&post)
	return nil
}

// ErrorPostNotFound is an error when post is not found in the list
var ErrorPostNotFound = fmt.Errorf("Product not found")
