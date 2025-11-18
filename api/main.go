package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	

	var err error
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Database connected successfully")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	router.ServeHTTP(w, r)
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, department, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, department, email string
		if err := rows.Scan(&id, &name, &department, &email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, gin.H{
			"id":         id,
			"name":       name,
			"department": department,
			"email":      email,
		})
	}

	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var user struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Email      string `json:"email"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int
	err := db.QueryRow(
		"INSERT INTO users (name, department, email) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Department, user.Email,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": userID})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Email      string `json:"email"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(
		"UPDATE users SET name=$1, department=$2, email=$3 WHERE id=$4",
		user.Name, user.Department, user.Email, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
