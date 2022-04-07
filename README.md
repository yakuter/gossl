# gossl

GoSSL is an SSL/TLS certificate tool written with Go and built with Love.

### Help
Help command displays default help and existing commands
```bash
go run main.go help
```

### Verify
You can verify x509 certificate with provided root CA. Both CA and certificate files need to be PEM encoded.

```bash
go run main.go verify --help
```
```bash
go run main.go verify -cafile ca.pem cert.pem
```
```bash
go run main.go verify -hostname 127.0.0.1 -cafile ca.pem cert.pem
```
