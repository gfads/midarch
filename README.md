# midarch [![Godoc](https://godoc.org/github.com/gfads/midarch?status.svg)](https://godoc.org/github.com/gfads/midarch)

## Installing & Configuring

#### 1. Download gMidArch

 - Download gMidArch (https://github.com/gfads/midarch.git) into <path>/gmidarch
 - Configure GOPATH=\<path-to-gopath>
 - Configure GOROOT=\<path-to-goroot>
 - Configure GMIDARCHDIR=\<path-to-gmidarch-gopath>
 
If you are going to use SSL then configure SSL environment variables:

 - Configure CA_PATH=\<path-to-ca-cert-file>
 - Configure CRT_PATH=\<path-to-cert-file>
 - Configure KEY_PATH=\<path-to-cert-key-file>
    
#### 2. Install FDR4

 - Download (https://cocotec.io/fdr/index.html) into <path>/fdr4
 - Install FDR4
 - Configure FDR4=<path>/fdr4

#### 3. Download additional packages used by gMidArch

```
go get gopkg.in/check.v1
go get github.com/kr/pretty
go get github.com/kr/text
go get github.com/vmihailenco/msgpack
go get github.com/vmihailenco/tagparser
```

#### 4. RabbitMQ

 - Download & Install RabbitMQ (https://www.rabbitmq.com/download.html)
 - Download Source code of 'Fibonacci Application': https://github.com/nsrosa70/midarch-go-v11.git
 - Code of Client & Server at: <download-dir>/src/transport/rabbitmq

#### 5. gRPC

 - Download & Install gRPC (https://grpc.io)
 - Download Source code of 'Fibonacci Application': https://github.com/nsrosa70/midarch-go-v11.git
 - Code of Client & Server at: <download-dir>/src/transport/grpc

#### 6. Download & Install AFiRM

 - Access https://github.com/andregpss/afirm.git

#### 7. Source codes of gRPC Client/Server

#### 8. Source code RabbitMQ


## Experiments

#### Scenario 1 (gMidArch - RPC)

1. Move to GMIDARCHDIR/src/apps/artefacts/madls
2. Edit 'midfibonacciserver.madl' (set 'Adaptability' to 'None')
3. Move to GMIDARCHDIR/src/apps/fibomiddleware/naming
4. Compile 'go build namingserver.go'
5. Start Naming Service: './namingserver'
6. Move to GMIDARCHDIR/src/apps/fibomiddleware/server
7. Compile 'go build server.go'
8. Start Fibonacci Server: './server'
9. Move to GMIDARCHDIR/src/apps/fibomiddleware/client
10. Compile 'go build client.go'
11. Start Fibonacci: './client <fibonacci-number> <number-of-requests>'

#### Scenario 2 (Adaptive gMidArch - RPC)

1. Move to GMIDARCHDIR/src/apps/artefacts/madls
2. Edit 'midfibonacciserver.madl' (set 'Adaptability' to 'Evolutive')
3. Move to GMIDARCHDIR/src/apps/fibomiddleware/naming
4. Compile 'go build naming.go'
5. Start Naming Service: './naming'
6. Move to GMIDARCHDIR/src/apps/fibomiddleware/server
7. Compile 'go build server.go'
8. Start Fibonacci Server: './server'
9. Move to GMIDARCHDIR/src/apps/fibomiddleware/client
10. Compile 'go build client.go'
11. Start Fibonacci Client: './client <fibonacci-number> <number-of-requests>'
12. Move GMIDARCHDIR/src/apps/fibomiddleware/injector
13. Compile 'go build injector.go'
14. Start injector: './injector <time-between-injections (in seconds)>'

#### Scenario 3 (gMidArch - MOM)

1. Move to GMIDARCHDIR/src/apps/artefacts/madls
2. Edit 'queueingserver.madl' (set 'Adaptability' to 'None')
3. Move to GMIDARCHDIR/src/apps/pubsub/queueing
4. Compile 'go build queueingserver.go'
5. Start Queueing Service: './queueingserver'
6. Move to GMIDARCHDIR/src/apps/pubsub/server
7. Compile 'go build server.go'
8. Start Fibonacci Server: './server'
9. Move to GMIDARCHDIR/src/apps/pubsub/client
10. Compile 'go build client.go'
11. Start Fibonacci Client: './client <fibonacci-number> <number-of-requests>'