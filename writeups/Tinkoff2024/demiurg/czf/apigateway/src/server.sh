#!/bin/bash
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config openssl-server.conf
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extensions v3_ca -extfile openssl-server.conf
