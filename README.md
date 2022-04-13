<p align="center"><img src="https://www.yakuter.com/wp-content/yuklemeler/Goossl.png" width="300"></p>

<p align="center">
<img src="https://img.shields.io/github/workflow/status/yakuter/gossl/Test/main" alt="Build Status">
<img src="https://img.shields.io/github/downloads/yakuter/gossl/total" alt="Total Downloads">
<img src="https://img.shields.io/codecov/c/github/yakuter/gossl/main" alt="Codecov branch">
<img src="https://img.shields.io/github/go-mod/go-version/yakuter/gossl" alt="Go Version">
<a href="https://pkg.go.dev/github.com/yakuter/gossl"><img src="https://pkg.go.dev/badge/github.com/yakuter/gossl.svg" alt="Go Version"></a><br>
<em>Gopher design by <a href="https://twitter.com/tgybalci">Tugay BALCI</a></em>
</p>

# GoSSL
GoSSL is a cross platform, easy to use SSL/TLS toolset written with Go and built with ❤️

## Features
- Generate RSA private and public key
- Generate x509 RSA Certificate Request (CSR)
- Generate x509 RSA Root CA
- Generate x509 RSA Certificate
- Verify a Certificate with a Root CA
- Generate SSH key pair
- Copy SSH public key to remote SSH server

## Install
Executable binaries can be downloaded at [Releases](https://github.com/yakuter/gossl/releases) page according to user's operating system and architecture. After download, extract compressed files and start using GoSSL via terminal.

### MacOS Homebrew Install
MacOS users can install GoSSl via Homebrew with the commands below.
```bash
brew tap yakuter/homebrew-tap
brew install gossl
```

## Commands
### help
`help` command displays default help and existing commands
```bash
./gossl help
```

### key
`key` command generates RSA private key with provided bit size.

```bash
./gossl key --help
./gossl key --bits 2048
./gossl key --bits 2048 --out private.key
./gossl key --bits 2048 --out private.key --withpub
```

### cert
`cert` command generates x509 SSL/TLS Certificate Request (CSR), Root CA and Certificate with provided private key.

Help
```bash
./gossl cert --help
```
Generate Certificate Request (CSR)
```bash
./gossl cert \
    --key private.key \
    --out cert.csr \
    --days 365 \
    --serial 12345 \
    --isCSR
```
Generate Root CA
```bash
./gossl cert \
    --key private.key \
    --out ca.pem \
    --days 365 \
    --serial 12345 \
    --isCA 
```
Generate Certificate
```bash
./gossl cert \
    --key private.key \
    --out cert.pem \
    --days 365 \
    --serial 12345
```

### verify
`verify` command verifies x509 certificate with provided root CA in PEM format.

```bash
./gossl verify --help
./gossl verify --cafile ./testdata/ca-cert.pem --certfile ./testdata/server-cert.pem
./gossl verify --dns 127.0.0.1 --cafile ./testdata/ca-cert.pem --certfile ./testdata/server-cert.pem
```

### ssh
`ssh` command generates SSH key pair with provided bit size just like `ssh-keygen` tool. These key pairs are used for automating logins, single sign-on, and for authenticating hosts.

```bash
./gossl key --help
./gossl key --bits 2048
./gossl key --bits 2048 -out ./id_rsa
// output will be written to ./id_rsa and ./id_rsa_pub files
```

### ssh-copy
`ssh-copy` connects remote SSH server, creates `/home/user/.ssh` directory and `authorized_keys` file in it and appends provided public key (eg, id_rsa.pub) to `authorized_keys` file just like `ssh-copy-id` tool.

```bash
./gossl ssh-copy --help
./gossl ssh-copy --pubkey /home/user/.ssh/id_rsa.pub remoteUser@remoteIP
// This command will ask for password to connect SSH server
```

### TODO
1. Add generate command for generating private key, root ca and x509 certificates in one command
2. Add cert template format read from yaml file
3. Add verification of a CA and http endpoint
4. Add certificate converter command like DER to PEM etc.
