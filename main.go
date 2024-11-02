package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Struct untuk Post
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Struct untuk User
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Data default
var Posts = []Post{
	{ID: 1, Title: "Judul Postingan Pertama", Content: "Ini adalah postingan pertama di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Title: "Judul Postingan Kedua", Content: "Ini adalah postingan kedua di blog ini.", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

var users = []User{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
}

// Middleware untuk autentikasi
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Memisahkan tipe dan nilai dari header
		payload := auth[len("Basic "):]
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Memeriksa user dan password
		credentials := string(decoded)
		valid := false
		for _, user := range users {
			if credentials == user.Username+":"+user.Password {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetupRouter mengatur semua route
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(authMiddleware())

	// Endpoint untuk GET /posts
	r.GET("/posts", func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam != "" {
			id, err := strconv.Atoi(idParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
				return
			}
			for _, post := range Posts {
				if post.ID == id {
					c.JSON(http.StatusOK, gin.H{"post": post})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Postingan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"posts": Posts})
	})

	// Endpoint untuk POST /posts
	r.POST("/posts", func(c *gin.Context) {
		var newPost Post
		if err := c.ShouldBindJSON(&newPost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		newPost.ID = len(Posts) + 1
		newPost.CreatedAt = time.Now()
		newPost.UpdatedAt = time.Now()
		Posts = append(Posts, newPost)
		c.JSON(http.StatusCreated, gin.H{"message": "Postingan berhasil ditambahkan", "post": newPost})
	})

	return r
}

// Fungsi utama
func main() {
	r := SetupRouter()
	r.Run(":8080") // Menjalankan server pada port 8080
}
