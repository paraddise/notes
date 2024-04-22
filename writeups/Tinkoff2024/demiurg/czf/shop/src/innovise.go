package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"crypto/x509"
	"log"
)

type Client struct {
	Config *tls.Config
}

func NewClient() *Client {
	caCert, _ := ioutil.ReadFile("ca.crt")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, _ := tls.LoadX509KeyPair("client.crt", "client.key")

	return &Client{
		Config: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs: caCertPool,
            Certificates: []tls.Certificate{cert},
		},
	}
}

func (c *Client) Post(urlString string, headers map[string]string) (string, []byte, error) {
	host := "apigateway:9443"

	conn, err := tls.Dial("tcp", host, c.Config)
	if err != nil {
		return "", []byte{}, err
	}
	defer conn.Close()

	request := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\n", urlString, host)
	for key, value := range headers {
		request += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	request += "\r\n"

	log.Println(request)

	log.Println("inno - sent request")

	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", []byte{}, err
	}

	log.Println("inno - read status")

	response := bufio.NewReader(conn)
	statusLine, err := response.ReadString('\n')
	if err != nil {
		return "", []byte{}, err
	}

	log.Println("inno - read headers")

	for {
		line, err := response.ReadString('\n')
		if err != nil {
			return "", []byte{}, err
		}
		if line == "\r\n" {
			break
		}
	}

	log.Println("inno - read body")

	body, err := ioutil.ReadAll(response)
	if err != nil {
		return "", []byte{}, err
	}

	return statusLine, body, nil
}
