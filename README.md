# gossl

GoSSL is an SSL/TLS certificate tool written with Go and built with ❤️

## Features
- Verify a server certificate with a root CA certificate
- Generate RSA key
- Generate x509 RSA CA certificate
- Generate x509 RSA server certificate
- Generate SSH key pair

## Commands
### help
help: Help command displays default help and existing commands
```bash
./gossl help
```

### verify
You can verify x509 certificate with provided root CA. Both CA and certificate files need to be PEM encoded.

```bash
./gossl verify --help
./gossl verify --cafile ./testdata/ca-cert.pem --certfile ./testdata/server-cert.pem
./gossl verify --dns 127.0.0.1 --cafile ./testdata/ca-cert.pem --certfile ./testdata/server-cert.pem
```

### key
Key command helps you to generate RSA private key with provided bit size.

```bash
./gossl key --help
./gossl key --bits 2048
./gossl key --bits 2048 -out private.key
```

### cert
Cert command generates x509 certificate with provided private key.

```bash
./gossl cert --help
```
```bash
// CA Certificate
./gossl cert \
    --key private.key \
    --out ca.pem \
    --days 365 \
    --serial 12345 \
    --isCA 
```
// Server Certificate
```bash
./gossl cert \
    --key private.key \
    --out cert.pem \
    --days 365 \
    --serial 12345
```

### ssh
SSH command helps you to generate SSH Key Pair with provided bit size.

```bash
./gossl key --help
./gossl key --bits 2048
./gossl key --bits 2048 -out ./id_rsa
// output will be written to ./id_rsa and ./id_rsa_pub files
```

### TODO
1. Add generate command for generating private key, root ca and x509 certificates in one command
2. Add test for cert
3. Add cert template format read from yaml file
4. Add verification of a CA and http endpoint
5. Add test for utils package
6. Add test for help package
7. Add ssh-copy-id feature to upload ssh key to remote server easily