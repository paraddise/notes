package controllers

import (
	"net/http"
	"context"
	"fmt"
	"time"
	"database/sql"
	"strconv"
	"strings"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	session := sessions.Default(c)
	username := session.Get("user_id")

	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

	query := fmt.Sprintf("SELECT age FROM users WHERE username = $1 \n")

	var age string
	err := db.QueryRowContext(ctx, query, username).Scan(&age)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"age": age, "username": username})
}

func UpdatePassword(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	session := sessions.Default(c)
	username := session.Get("user_id")
	// username := c.PostForm("username")
	password := c.PostForm("password")

	password = strings.ReplaceAll(password, "🔥", "💨")

	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	go func(db *sql.DB, password, username interface{}) {

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
	
		query := fmt.Sprintf("UPDATE users SET password = '%s' \n", password)// hashPassword(password))
		query = query + "WHERE username = $1"
	
		_, err := db.ExecContext(ctx, query, username)

		if err != nil {
            log.Printf("Failed to insert message: %v", err)
        }

    }(db, password, username)

	Logout(c)
}

func UpdateAge(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	session := sessions.Default(c)
	username := session.Get("user_id")
	age := c.PostForm("age")

	if age == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Age cannot be empty"})
		return
	}

	var ageInt int
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age format"})
		return
	}

	go func(db *sql.DB, ageInt, username interface{}) {

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		query := fmt.Sprintf("UPDATE users SET age = %d \n", ageInt)
		query = query + "WHERE username = $1"

		_, err = db.ExecContext(ctx, query, username)

		if err != nil {
            log.Printf("Failed to insert message: %v", err)
        }

    }(db, ageInt, username)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Age was succsessfuly updated")})
}