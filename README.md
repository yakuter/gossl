# gossl

GoSSL is an SSL/TLS certificate tool written with Go and built with Love.

## Commands
### help
help: Help command displays default help and existing commands
```bash
go run main.go help
```

### verify
You can verify x509 certificate with provided root CA. Both CA and certificate files need to be PEM encoded.

```bash
go run main.go verify --help
```
```bash
go run main.go verify --cafile ./testdata/ca-cert.pem ./testdata/server-cert.pem
```
```bash
go run main.go verify --hostname 127.0.0.1 --cafile ./testdata/ca-cert.pem ./testdata/server-cert.pem
```

### key
Key command helps you to generate RSA private key with provided bit size.

```bash
go run main.go key --help
```
```bash
go run main.go key 2048
```
```bash
go run main.go key -out private.key 2048
```

### TODO
1. Prepare release with Goreleaser for Windows, MacOS, Linux Deb, Linux RPM environments.
2. Implement this logger: https://github.com/binalyze/httpreq/blob/main/logger.go
3. Add generate command for generating private key, root ca and x509 certificates
