#!/bin/bash
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/C=US/ST=State/L=City/O=Organization/CN=client.example.com"
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt
openssl verify -CAfile ca.crt client.crt
cp client.crt client.key ca.crt ../../shop/src/
