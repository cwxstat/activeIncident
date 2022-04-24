```bash

docker run -it --rm -v ${PWD}:/work -w /work debian bash


apt-get update && apt-get install -y curl &&
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64 -o /usr/local/bin/cfssl && \
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64 -o /usr/local/bin/cfssljson && \
chmod +x /usr/local/bin/cfssl && \
chmod +x /usr/local/bin/cfssljson

#generate ca in /tmp
cfssl gencert -initca ./tls/ca-csr.json | cfssljson -bare /tmp/ca

#generate certificate in /tmp
cfssl gencert \
  -ca=/tmp/ca.pem \
  -ca-key=/tmp/ca-key.pem \
  -config=./tls/ca-config.json \
  -hostname="mongo,mongo.mongodb.svc.cluster.local,mongo.default.svc,localhost,127.0.0.1,mongo.pigbot.svc.cluster.local,34.111.92.27" \
  -profile=default \
  ./tls/ca-csr.json | cfssljson -bare /tmp/mongo-certs


```
