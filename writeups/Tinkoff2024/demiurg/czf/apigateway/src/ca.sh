#!/bin/bash
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 1024 -out ca.crt -subj "/C=US/ST=New York/L=New York/O=Example Company/OU=IT Department/CN=example.com"
