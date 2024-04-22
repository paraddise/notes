package controllers

import (
	"fmt"
	"net/http"
	"context"
	"time"
	"database/sql"
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetMessage(c *gin.Context) {
	type Message struct {
		Title    string
		Body     string
		Username string
	}

	db := c.MustGet("db").(*sql.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

	var messages []Message

	query := fmt.Sprintf("SELECT title, body, username FROM messages \n")

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occurred"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Title, &m.Body, &m.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occurred"})
			return
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func SendMessage(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	session := sessions.Default(c)
	username := session.Get("user_id")
	title := c.PostForm("title")
	body := c.PostForm("body")

	title = strings.ReplaceAll(title, "🔥", "💨")
    body = strings.ReplaceAll(body, "🔥", "💨")


	go func(db *sql.DB, title, body, username interface{}) {

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        query := fmt.Sprintf("INSERT INTO messages (title, body, username) VALUES ($1, $2, $3) \n")
        _, err := db.ExecContext(ctx, query, title, body, username)

		if err != nil {
            log.Printf("Failed to insert message: %v", err)
        }
    }(db, title, body, username)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Message has been sent")})
}