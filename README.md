# midarch [![Go Reference](https://pkg.go.dev/badge/github.com/gfads/midarch.svg)](https://pkg.go.dev/github.com/gfads/midarch) [![Go Report Card](https://goreportcard.com/badge/github.com/gfads/midarch)](https://goreportcard.com/report/github.com/gfads/midarch)

## Installing & Configuring

#### 1. Download gMidArch

- Download gMidArch (https://github.com/gfads/midarch.git) into <path>/gmidarch
- Configure GMIDARCH=\<path-to-gmidarch>

If you are going to use SSL then configure SSL environment variables:

- Configure CA_PATH=\<path-to-ca-cert-file>
- Configure CRT_PATH=\<path-to-cert-file>
- Configure KEY_PATH=\<path-to-cert-key-file>

#### 2. Install FDR4

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
