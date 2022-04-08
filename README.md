# gossl

GoSSL is an SSL/TLS certificate tool written with Go and built with ❤️

## Commands
### help
help: Help command displays default help and existing commands
```bash
gossl help
```

### verify
You can verify x509 certificate with provided root CA. Both CA and certificate files need to be PEM encoded.

```bash
gossl verify --help
```
```bash
gossl verify --cafile ./testdata/ca-cert.pem ./testdata/server-cert.pem
```
```bash
gossl verify --dns 127.0.0.1 --cafile ./testdata/ca-cert.pem ./testdata/server-cert.pem
```

### key
Key command helps you to generate RSA private key with provided bit size.

```bash
gossl key --help
```
```bash
// Default output is os.Stdout
gossl key 2048
```
```bash
gossl key -out private.key 2048
```

### TODO
1. Implement this logger: https://github.com/binalyze/httpreq/blob/main/logger.go
2. Add generate command for generating private key, root ca and x509 certificates
