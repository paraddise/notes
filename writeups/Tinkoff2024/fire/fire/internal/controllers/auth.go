package controllers

import (
	"net/http"
	"context"
	"fmt"
	"time"
	"database/sql"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	username := c.PostForm("username")
	password := c.PostForm("password")

	username = strings.ReplaceAll(username, "🔥", "💨")
	password = strings.ReplaceAll(password, "🔥", "💨")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password cannot be empty"})
		return
	}

	if len(username) < 9 || len(password) < 9 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password must be at least 9 characters"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

	var dbUsername, dbPassword string

	query := fmt.Sprintf("SELECT username, password FROM Users WHERE username = $1 \n")
	
	err := db.QueryRowContext(ctx, query, username).Scan(&dbUsername, &dbPassword)

	if err != nil {
        if err != sql.ErrNoRows {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Username is already taken"})
			return
        }
    }

	query = fmt.Sprintf("INSERT INTO Users (username, password, age) VALUES ($1, $2, $3) \n")

    _, err = db.ExecContext(ctx, query, username, password, 18)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully signuped"})
}

func Signin(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password cannot be empty"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

	var dbUsername, dbPassword string
	query := fmt.Sprintf("SELECT username, password FROM Users WHERE username = $1 \n")
	
	err := db.QueryRowContext(ctx, query, username).Scan(&dbUsername, &dbPassword)

	if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login creadentials"})
			return
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occurred"})
			return
        }
        return
    }

	if password != dbPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
	 	return
	}


	session := sessions.Default(c)
	session.Set("user_id", dbUsername)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed in"})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
	}

	c.Redirect(http.StatusMovedPermanently, "/signin")
}