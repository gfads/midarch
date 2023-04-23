# Experiments with docker

## Running experiments from images on docker hub

### Download and install docker

- Download and install docker
- Initialize docker swarm
```bash
docker swarm init
```

### Run the experiment

- Chose fibonacci place and sample size by changing the environment varibles inside the docker compose file of the desired stack (e.g., dc-fibomiddleware-ssl.yml)
- From the base folder of the project deploy the desired stack, e.g.:
```bash
docker stack deploy -c evaluation/experiments/docker/dc-fibomiddleware-ssl.yml fibomiddleware-ssl 
```

- Observe the results of the experiments in the service logs, e.g.:
```bash
docker service logs -f fibomiddleware-ssl_client
```

## Running experiments from your own images

### Download and install docker

- Download and install docker
- Initialize docker swarm
```bash
docker swarm init
```

### Generate ssl keys (if experiment with the use of ssl)

- Generate private key for the Certificate Authority
```bash
openssl genrsa -des3 -out myCA.key 2048
```

- Generate root certificate for the Certificate Authority
```bash
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 1825 -out myCA.pem
```

- Generate the server private key for the certificate. Remember to change <your.domain.com> to the name of the server.
  These gMidArch experiments using docker currently uses two kinds of servers, "namingserver" and "server", therefore, "namingserver.key" and "server.key"
```bash
openssl genrsa -out <your.domain.com>.key 2048
```

- Generate a Certificate Signing Request (CSR)
```bash
openssl req -new -key <your.domain.com>.key -out <your.domain.com>.csr
```

- Create a configuration file for the certificate
```bash
cat > <your.domain.com>.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = <your.domain.com>
EOF
```

- Create the certificate using the CSR, the CA private key, the CA certificate and the config file
```bash
openssl x509 -req -in <your.domain.com>.csr -CA myCA.pem -CAkey myCA.key -CAcreateserial \ 
  -out <your.domain.com>.crt -days 825 -sha256 -extfile <your.domain.com>.ext
```

### Build your own images

- Configure the madl file of the desired application
- From the base folder of the project run the desired build scripts, e.g.:
```bash
./evaluation/experiments/docker/fibomiddleware-build.sh
```

### Run the experiment

- From the base folder of the project deploy the desired stack, e.g.:
```bash
docker stack deploy -c evaluation/experiments/docker/dc-fibomiddleware-ssl.yml fibomiddleware-ssl 
```

- Observe the results of the experiments in the service logs, e.g.:
```bash
docker service logs -f fibomiddleware-ssl_client
```