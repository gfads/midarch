# midarch [![Go Reference](https://pkg.go.dev/badge/github.com/gfads/midarch.svg)](https://pkg.go.dev/github.com/gfads/midarch) [![Go Report Card](https://goreportcard.com/badge/github.com/gfads/midarch)](https://goreportcard.com/report/github.com/gfads/midarch)

## Installing & Configuring

### 1. Download gMidArch

- Download gMidArch (https://github.com/gfads/midarch.git) into any place you like (\<path-to-gmidarch>)

### 2. Configure gMidArch

#### 2.1 Environment variable

- Configure environment variable GMIDARCH=\<path-to-gmidarch>

#### 2.2 Certificates

If you are going to use TLS then configure TLS environment variables:

- Configure CA_PATH=\<path-to-ca-cert-file>
- Configure CRT_PATH=\<path-to-cert-file>
- Configure KEY_PATH=\<path-to-cert-key-file>

If you don't have certificades and want to generate certificates just for developing and/or testing:

```bash
mkdir -p ./examples/certs/
cd ./examples/certs/

# Generate CA key
openssl genrsa -out myCA.key 1024
# Generate the CA Certificate
## Ask for common name, country, ...
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 3650 -out myCA.pem
## Or provide in script common name, country, ...
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 3650 -out myCA.pem -subj '/CN=MidArchCA/C=BR/ST=Pernambuco/L=Recife/O=MidArch'

# Generate Server Certificate
openssl req -new -nodes -out server.csr -newkey -sha256 -keyout server.key -subj '/CN=localhost/C=BR/ST=Pernambuco/L=Recife/O=MidArch'
# Sign the Server Certificate
openssl x509 -req -in server.csr -CA myCA.pem -CAkey myCA.key -CAcreateserial -sha256 -days 3650 -out server.pem
```

Where:

- CN = <Common Name (e.g. server FQDN or YOUR name)>
- C = <Country Name (2 letter code)>
- ST = <State or Province Name (eg, city)>
- L = <Locality Name (eg, city)>
- O = <Organization Name (eg, company)>
- OU = <Organizational Unit Name (eg, section)>

### 3. Install FDR4

- Download (https://cocotec.io/fdr/index.html) into <path>/fdr4
- Install FDR4
- Configure FDR4=<path>/fdr4

## Experiments

#### Scenario 1 (gMidArch - RPC)

1. Go to GMIDARCH/pkg/apps/artefacts/madls
2. Edit 'midfibonacciserver.madl' (set 'Adaptability' to 'None')
3. Go to GMIDARCH/examples/fibonaccidistributed/naming
4. Compile 'go build naming.go'
5. Start Naming Service: './naming'
6. Go to GMIDARCH/examples/fibonaccidistributed/server
7. Compile 'go build server.go'
8. Start Fibonacci Server: './server'
9. Go to GMIDARCH/examples/fibonaccidistributed/client
10. Compile 'go build client.go'
11. Start Fibonacci: './client <fibonacci-number> <number-of-requests>'

#### Scenario 2 (Adaptive gMidArch - RPC)

1. Go to GMIDARCH/pkg/apps/artefacts/madls
2. Edit 'midfibonacciserver.madl' (set 'Adaptability' to 'Evolutive')
3. Go to GMIDARCH/examples/fibonaccidistributed/naming
4. Compile 'go build naming.go'
5. Start Naming Service: './naming'
6. Go to GMIDARCH/examples/fibonaccidistributed/server
7. Compile 'go build server.go'
8. Start Fibonacci Server: './server'
9. Go to GMIDARCH/examples/fibonaccidistributed/client
10. Compile 'go build client.go'
11. Start Fibonacci Client: './client <fibonacci-number> <number-of-requests>'
