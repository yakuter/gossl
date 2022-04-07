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

### TODO
1. Prepare release with Goreleaser for Windows, MacOS, Linux Deb, Linux RPM environments.
2. Implement this logger: https://github.com/binalyze/httpreq/blob/main/logger.go
3. Add generate command for generating private key, root ca and x509 certificates
