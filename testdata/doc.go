package testdata

/*
Private Keys and Certificates in this folder are generated with the commands below

openssl genrsa 2048 > ca-key.pem

openssl req -new -x509 -nodes -days 365000 \
   -key ca-key.pem \
   -out ca-cert.pem \
   -addext "subjectAltName = IP:127.0.0.1"

openssl req -newkey rsa:2048 -nodes -days 365000 \
   -keyout server-key.pem \
   -out server-req.pem \
   -addext "subjectAltName = IP:127.0.0.1"

openssl x509 -req -days 365000 -set_serial 01 \
   -in server-req.pem \
   -out server-cert.pem \
   -CA ca-cert.pem \
   -CAkey ca-key.pem \
   -extfile v3.ext

v3.ext file content:
subjectAltName = @alt_names

[alt_names]
IP.1 = 127.0.0.1
*/
