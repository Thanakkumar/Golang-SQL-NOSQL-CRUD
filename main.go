package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io"
	"net/http"
)

// User model
type User struct {
	gorm.Model
	Name  string
	Email string
}

var db *gorm.DB
var err error

func main() {
	// Initialize the SQLite database
	db, err = gorm.Open("mysql", "username:password@tcp(localhost:3306)/user?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// AutoMigrate the User model
	db.AutoMigrate(&User{})

	// Create a new Gin router
	r := gin.Default()

	// Define routes
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUser)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	r.GET("/callapi",CallApi)

	// Run the server
	r.Run(":8080")
}

// GetUsers retrieves all users
func GetUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(200, users)
}

// GetUser retrieves a specific user by ID
func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, user)
	}
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	db.Create(&user)
	c.JSON(200, user)
}

// UpdateUser updates a user by ID
func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&user)
	db.Save(&user)
	c.JSON(200, user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	d := db.Where("id = ?", id).Delete(&user)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func CallApi(c *gin.Context){
	response, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	fmt.Println("get call", response)

	// if there was no error, you should close the body
	defer response.Body.Close()

	// hence this condition is moved into its own block
	if response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	//only way to read response
	body, err := io.ReadAll(response.Body)

	c.JSON(200, string(body))
}
