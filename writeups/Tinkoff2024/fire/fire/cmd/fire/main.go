package main

import (
	"fire/internal/controllers"
	"fire/internal/middlewares"
	"fmt"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func init() {
	viper.SetEnvPrefix("fire")
	viper.AutomaticEnv()
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func GetDBConfig() *DBConfig {
	return &DBConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
	}
}

func main() {
	dbConfig := GetDBConfig()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	store := cookie.NewStore([]byte(viper.GetString("SECRET_KEY")))
	r.Use(sessions.Sessions("session", store))

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	unauthorized := r.Group("/")
	unauthorized.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", nil)
	})
	unauthorized.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	auth := r.Group("/")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/signin", controllers.Signin)
	}

	authorized := r.Group("/")
	authorized.Use(middlewares.AuthRequired)
	{
		authorized.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
		authorized.GET("/profile", func(c *gin.Context) {
			c.HTML(http.StatusOK, "profile.html", nil)
		})
		authorized.GET("/logout", controllers.Logout)
		authorized.POST("/message/send", controllers.SendMessage)
		authorized.POST("/profile/age", controllers.UpdateAge)
		authorized.POST("/profile/password", controllers.UpdatePassword)
		authorized.GET("/message/get", controllers.GetMessage)
		authorized.GET("/profile/get", controllers.GetProfile)
	}

	r.Run()
}
