package main

import (
    "net/http"
    "fmt"
    "github.com/gin-gonic/gin"
    "encoding/json"
    "log"
)

const (
	UserIDHeaderName = "X-UserID"
	FromServiceHeaderName = "X-From-Service"
)

func main() {
    router := gin.Default()

    router.POST("/api/buyColor", changeColor)
    router.POST("/api/buyPremium", buyPremium)

    router.Run("0.0.0.0:6080")
}

type errorResp struct {
    Success bool
    Error string
}

func changeColor(c *gin.Context) {
    userID := c.GetHeader(UserIDHeaderName)
    color := c.Query("color")

    headers := map[string]string{
        UserIDHeaderName: userID,
        "Connection": "close",
        FromServiceHeaderName: "Yes",
    }

    price := 5
    client := NewClient()
    _, body, err := client.Post("https://apigateway:9443/api/getInfo", headers)
    if err != nil {
        resp := errorResp{
            Success: false,
            Error: err.Error(),
        }
        c.IndentedJSON(http.StatusOK, resp)
        return
    }

    type userInfoT struct {
		Offset     int
		Multiplier int
		Color      int
	}
    log.Println(string(body))
    var userInfo userInfoT
    err = json.Unmarshal(body, &userInfo)
    if err != nil {
        resp := errorResp{
            Success: false,
            Error: err.Error(),
        }
        c.IndentedJSON(http.StatusOK, resp)
        return
    }
    
    if userInfo.Offset >= price {
        newOffset := userInfo.Offset - price
        params := fmt.Sprintf("prevOffset=%d&offset=%d", userInfo.Offset, newOffset)

        client.Post("https://apigateway:9443/api/setOffset?" + params, headers)

        params = fmt.Sprintf("color=%s", color)
        client.Post("https://apigateway:9443/api/setUserColor?" + params, headers)

        resp := errorResp{
            Success: true,
            Error: "",
        }
        c.IndentedJSON(http.StatusOK, resp)
    } else {
        resp := errorResp{
            Success: false,
            Error: "Недостаточно средств, нужно 5 кликов",
        }
        c.IndentedJSON(http.StatusOK, resp)
    }
}

func buyPremium(c *gin.Context) {
    userID := c.GetHeader(UserIDHeaderName)
    headers := map[string]string{
        UserIDHeaderName: userID,
        "Connection": "close",
        FromServiceHeaderName: "Yes",
    }
    price := 213370
    client := NewClient()
    _, body, err := client.Post("https://apigateway:9443/api/getInfo", headers)
    if err != nil {
        resp := errorResp{
            Success: false,
            Error: err.Error(),
        }
        c.IndentedJSON(http.StatusOK, resp)
        return
    }

    type userInfoT struct {
		Offset     int
		Multiplier int
		Color      int
    }
    log.Println(string(body))
    var userInfo userInfoT
    err = json.Unmarshal(body, &userInfo)
    if err != nil {
        resp := errorResp{
            Success: false,
            Error: err.Error(),
        }
        c.IndentedJSON(http.StatusOK, resp)
        return
    }
    
    if userInfo.Offset >= price {
        newOffset := userInfo.Offset - price
        params := fmt.Sprintf("prevOffset=%d&offset=%d", userInfo.Offset, newOffset)

        client.Post("https://apigateway:9443/api/setOffset?" + params, headers)

        params = fmt.Sprintf("multiplier=%d", 10)
        client.Post("https://apigateway:9443/api/setUserMultiplier?" + params, headers)

        resp := errorResp{
            Success: true,
            Error: "",
        }
        c.IndentedJSON(http.StatusOK, resp)
    } else {
        resp := errorResp{
            Success: false,
            Error: "Недостаточно средств, нужно 213370 кликов",
        }
        c.IndentedJSON(http.StatusOK, resp)
    }
}