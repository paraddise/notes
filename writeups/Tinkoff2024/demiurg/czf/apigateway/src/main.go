package main

import (
	"net/http"
	"net/url"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"log"
	"os"
	"crypto/hmac"
	"encoding/hex"
	"strings"
	"errors"
	"io"
	"crypto/sha256"
)

const (
	UserIDHeaderName = "X-UserID"
	FromServiceHeaderName = "X-From-Service"
	ForwardHostHeaderName = "X-Forward-Host"
)


var sessionKey = []byte(os.Getenv("SESSION_SIGN_KEY"))
var responseHeadersToProxy = map[string]bool{
	"Set-Cookie": true,
}

func main() {
	http.HandleFunc("/api/setOffset", proxyForUserInfo)
	http.HandleFunc("/api/setUserColor", proxyForUserInfo)
	http.HandleFunc("/api/setUserMultiplier", proxyForUserInfo)

	http.HandleFunc("/api/buyColor", proxyForShop)
	http.HandleFunc("/api/buyPremium", proxyForShop)
	http.HandleFunc("/api/step", proxyForUserInfo)
	http.HandleFunc("/api/getInfo", proxyForUserInfo)

	http.HandleFunc("/api/register", proxyForUserInfo)
	http.HandleFunc("/api/login", proxyForUserInfo)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	caCert, _ := ioutil.ReadFile("ca.crt")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	serverTLS := &http.Server{
		Addr:      ":9443",
		TLSConfig: tlsConfig,
	}

	go serverTLS.ListenAndServeTLS("server.crt", "server.key")

	server := &http.Server{
		Addr:      ":9080",
		TLSConfig: tlsConfig,
	}

	server.ListenAndServe()
}

func proxyForUserInfo(w http.ResponseWriter, r *http.Request) {
	proxyTo("userinfo:6080", w, r)
}

func proxyForShop(w http.ResponseWriter, r *http.Request) {
	proxyTo("shop:6080", w, r)
}

func proxyTo(host string, w http.ResponseWriter, r *http.Request) {
	userID := ""

	proxyHeader := make(http.Header)
	forwardHost := r.Header.Get(ForwardHostHeaderName)
	if forwardHost != "" {
		proxyHeader.Add(ForwardHostHeaderName, forwardHost)
	}

	fromService := r.Header.Get(FromServiceHeaderName) != ""
	if fromService {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			io.WriteString(w, "not peer certificates")
			return
		}
		userID = r.Header.Get(UserIDHeaderName)
		proxyHeader.Add(UserIDHeaderName, userID)
		proxyHeader.Add(FromServiceHeaderName, "Yes")
	} else {
		cookie, err := r.Cookie("session")
		if err == nil {
			err, ok, payload := parseAndVerifyCookie(cookie)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			if !ok {
				io.WriteString(w, "fail to verify cookie")
				return
			}

			proxyHeader.Add(UserIDHeaderName, payload)
		}
	}

	proxyReq := &http.Request{}
	proxyReq.Method = r.Method
	proxyReq.URL = &url.URL{}
	proxyReq.URL.Scheme = "http"
	proxyReq.URL.Path = r.URL.Path
	proxyReq.URL.Host = host
	proxyReq.URL.RawQuery = r.URL.RawQuery
	proxyReq.Header = proxyHeader

	client := http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Println("error on http.Client.Do:", err.Error())
		return
	}

	for k, v := range resp.Header {
		ok := responseHeadersToProxy[k]
		if ok {
			w.Header().Add(k, v[0])
		}
	}

	w.WriteHeader(resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	w.Write(body)
}

func parseAndVerifyCookie(cookie *http.Cookie) (error, bool, string) {
	value := cookie.Value
	
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return errors.New("more than two parts in session"), false, ""
	}

	payload := parts[0]
	sign := []byte(parts[1])

	mac := hmac.New(sha256.New, sessionKey)
	mac.Write([]byte(payload))

	expectedMAC := mac.Sum(nil)
	return nil, hmac.Equal(sign, []byte(hex.EncodeToString(expectedMAC))), payload
}
