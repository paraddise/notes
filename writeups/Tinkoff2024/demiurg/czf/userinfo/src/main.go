package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"encoding/hex"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"strconv"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v4/pgxpool"
)

const (
	UserIDHeaderName = "X-UserID"
	FromServiceHeaderName = "X-From-Service"
)

var sessionKey = []byte(os.Getenv("SESSION_SIGN_KEY"))

type App struct {
	DB *pgxpool.Pool
}

func main() {
	dsn := "postgres://postgres:password@db:5432/postgres?sslmode=disable"
	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	app := &App{
		DB: dbpool,
	}

	router := gin.Default()

	router.POST("/api/setOffset", app.setOffset)
	router.POST("/api/setUserColor", app.setUserColor)
	router.POST("/api/setUserMultiplier", app.setUserMultiplier)

	router.POST("/api/step", app.step)
	router.POST("/api/getInfo", app.getInfo)

	router.POST("/api/register", app.register)
	router.POST("/api/login", app.login)

	router.Run("0.0.0.0:6080")
}

type errorResp struct {
	Success bool
	Error   string
}

func (app *App) setOffset(c *gin.Context) {
	if forbidRequestFromUser(c) {
		return
	}

	userID := c.GetHeader(UserIDHeaderName)
	prevOffsetStr := c.Query("prevOffset")

	prevOffset, err := strconv.Atoi(prevOffsetStr)
	if err != nil {
		log.Println("strconv.Atoi:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on validation",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Println("strconv.Atoi:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on validation",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}


	query := "UPDATE users SET \"offset\" = $1 WHERE user_id = $2 and \"offset\" = $3"
	_, err = app.DB.Exec(context.Background(), query, offset, userID, prevOffset)

	if err != nil {
		log.Println("setOffset:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := errorResp{
		Success: true,
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (app *App) setUserColor(c *gin.Context) {
	if forbidRequestFromUser(c) {
		return
	}

	userID := c.GetHeader(UserIDHeaderName)
	colorStr := c.Query("color")

	colorToEnum := map[string]int{
		"black": 1,
		"white": 0,
		"aqua": 2,
	}

	color, ok := colorToEnum[colorStr]
	if !ok {
		resp := errorResp{
			Success: false,
			Error:   "not found color enum",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	query := "UPDATE users SET color = $1 WHERE user_id = $2"
	_, err := app.DB.Exec(context.Background(), query, color, userID)

	if err != nil {
		log.Println("setUserColor:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := errorResp{
		Success: true,
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (app *App) setUserMultiplier(c *gin.Context) {
	if forbidRequestFromUser(c) {
		return
	}

	userID := c.GetHeader(UserIDHeaderName)
	multiplierStr := c.Query("multiplier")

	muttiplier, err := strconv.Atoi(multiplierStr)
	if err != nil {
		log.Println("strconv.Atoi:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on validation",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	query := "UPDATE users SET multiplier = $1 WHERE user_id = $2"
	_, err = app.DB.Exec(context.Background(), query, muttiplier, userID)

	if err != nil {
		log.Println("setUserMultiplier:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := errorResp{
		Success: true,
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (app *App) step(c *gin.Context) {
	userID := c.GetHeader(UserIDHeaderName)

    var multiplier int
    var offset int
	query := "SELECT multiplier, \"offset\" FROM users WHERE user_id = $1"
	rows := app.DB.QueryRow(context.Background(), query, userID)

	err := rows.Scan(&multiplier, &offset)
	if err != nil {
		log.Println("step:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

    now := time.Now().Unix()
    
    offset += 1 * multiplier

	query = "UPDATE users SET \"offset\" = $1, update_at = $2 WHERE user_id = $3 and update_at < $4"
	_, err = app.DB.Exec(context.Background(), query, offset, now, userID, now)

	if err != nil {
		log.Println("step:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := errorResp{
		Success: true,
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (app *App) getInfo(c *gin.Context) {
	userID := c.GetHeader(UserIDHeaderName)
	query := "SELECT \"offset\", multiplier, color FROM users WHERE user_id = $1"

	type currInfoRespT struct {
		Offset     int
		Multiplier int
		Color      int
		Success bool
		Prize string
	}

	var currInfoResp currInfoRespT
	rows := app.DB.QueryRow(context.Background(), query, userID)
	err := rows.Scan(&currInfoResp.Offset, &currInfoResp.Multiplier, &currInfoResp.Color)
	if err != nil {
		log.Println("currentInfo:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	if currInfoResp.Offset >= 413370 {
		currInfoResp.Prize = GetFlag(c)
	}

	currInfoResp.Success = true
	c.IndentedJSON(http.StatusOK, currInfoResp)
}

func (app *App) register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	query := "SELECT user_id FROM users WHERE username = $1"
	rows := app.DB.QueryRow(context.Background(), query, username)

	type userInfoRowT struct {
		UserID string
	}

	var userInfoRow userInfoRowT
	err := rows.Scan(&userInfoRow.UserID)
	if err == nil {
		resp := errorResp{
			Success: false,
			Error:   "user already exists",
		}

		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	query = "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err = app.DB.Exec(context.Background(), query, username, password)
	if err != nil {
		log.Println("register:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	resp := errorResp{
		Success: true,
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func (app *App) login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	query := "SELECT password, user_id FROM users WHERE username = $1"
	rows := app.DB.QueryRow(context.Background(), query, username)

	type userInfoRowT struct {
		Password string
		UserID string
	}

	var userInfoRow userInfoRowT
	err := rows.Scan(&userInfoRow.Password, &userInfoRow.UserID)

	if err != nil {
		log.Println("login:", err.Error())

		resp := errorResp{
			Success: false,
			Error:   "error on query",
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	if password != userInfoRow.Password {
		resp := errorResp{
			Success: true,
			Error:   "password not matched" + "," + password + "," + userInfoRow.Password,
		}
		c.IndentedJSON(http.StatusOK, resp)
		return
	}

	sessionValue := createAndSignSession(string(userInfoRow.UserID))
	c.SetCookie("session", sessionValue, 3600000, "", "", false, false)

	resp := errorResp{
		Success: true,
		Error:   "password matched",
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func createAndSignSession(userID string) string {
	payload := []byte(userID)

	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte(payload))

	expectedMAC := mac.Sum(nil)

	return string(payload) + "." + string(hex.EncodeToString(expectedMAC))
}

func forbidRequestFromUser(c *gin.Context) (bool) {
	fromService := c.GetHeader(FromServiceHeaderName)
	if fromService != "Yes" {
		c.Status(http.StatusForbidden)
		return true
	}
	return false
}

func GetFlag(c *gin.Context) (string) {
	return os.Getenv("FLAG")
} 