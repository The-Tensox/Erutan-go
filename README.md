
# Grpc_push

```bash
protoc -I protos/realtime --go_out=plugins=grpc:protos/realtime protos/realtime/realtime.proto
```

```bash
cp protos/realtime/realtime.proto ~/Documents/unity/Erutan/Assets/Protos/Realtime
```

```bash

Edit your /etc/ssl/openssl.cnf on the logstash host - add subjectAltName = IP:192.168.2.107 in [v3_ca] section


openssl genrsa -out server1.key 2048
openssl req -new -x509 -sha256 -key server1.key \
              -out server1.crt -days 3650



export GRPC_VERBOSITY=INFO
cp server1.crt ~/Documents/unity/Erutan/Assets
go build -o ./bin/erutan && bin/erutan -s -v -p "" -h "0.0.0.0:50051"
```
