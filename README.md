# gossl

GoSSL is an SSL/TLS certificate tool written with Go and built with ❤️

## Features
- Verify a server certificate with a root CA
- Generate RSA key
- Generate x509 RSA root CA
- Generate x509 RSA server certificate
- Generate SSH key pair

## Commands
### help
`help` command displays default help and existing commands
```bash
./gossl help
```

### verify
`verify` command verifies x509 certificate with provided root CA in PEM format.

```bash
./gossl verify --help
./gossl verify --cafile ./testdata/ca-cert.pem --certfile ./testdata/server-cert.pem
./gossl verify --dns 127.0.0.1 --cafile ./testdata/ca-cert.pem --certfile ./testdata/server-cert.pem
```

### key
`key` command generates RSA private key with provided bit size.

```bash
./gossl key --help
./gossl key --bits 2048
./gossl key --bits 2048 -out private.key
```

### cert
`cert` command generates x509 certificate with provided private key.

```bash
./gossl cert --help
```
// Generate Root CA
```bash
./gossl cert \
    --key private.key \
    --out ca.pem \
    --days 365 \
    --serial 12345 \
    --isCA 
```
// Generate Certificate
```bash
./gossl cert \
    --key private.key \
    --out cert.pem \
    --days 365 \
    --serial 12345
```

### ssh
`ssh` command generates SSH Key Pair with provided bit size just like `ssh-keygen` tool.

```bash
./gossl key --help
./gossl key --bits 2048
./gossl key --bits 2048 -out ./id_rsa
// output will be written to ./id_rsa and ./id_rsa_pub files
```

### ssh-upload
`ssh-upload` connects remote SSH server and appends provided public key (eg, id_rsa.pub) to `authorized_keys` file just like `ssh-copy-id` tool.

```bash
./gossl ssh-upload --help
./gossl ssh-upload --pubkey /home/user/.ssh/id_rsa.pub remoteUser@remoteIP
// This command will ask for password to connect SSH server
```

ssh-upload --pubkey /Users/yakuter/.ssh/id_rsa.pub root@178.18.195.42

### TODO
1. Add generate command for generating private key, root ca and x509 certificates in one command
2. Add test for cert
3. Add cert template format read from yaml file
4. Add verification of a CA and http endpoint
5. Add test for utils package
6. Add test for help package
7. Add test for ssh-upload
