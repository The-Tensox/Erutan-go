
# Erutan-go

Simulating darwinian evolution, fully networked allowing several clients to have a 3D vizualisation.

To be used with [the Unity client](https://github.com/The-Tensox/Erutan-unity)

# Installation

```bash
export UNITY_PROJECT_PATH="/home/louis/Documents/unity/Erutan"
go get github.com/The-Tensox/Erutan-go
protoc -I protos/realtime --go_out=plugins=grpc:protos/realtime protos/realtime/realtime.proto
cp protos/realtime/realtime.proto $UNITY_PROJECT_PATH/Assets/Protos/Realtime
```

# SSL/TLS configuration

```bash
# Edit your /etc/ssl/openssl.cnf on the logstash host - add subjectAltName = IP:192.168.2.107 in [v3_ca] section

# Then
openssl genrsa -out server1.key 2048
openssl req -new -x509 -sha256 -key server1.key \
              -out server1.crt -days 3650

# Copy to Unity project
cp server1.crt $UNITY_PROJECT_PATH/Assets/StreamingAssets
```

# Run

```bash
go build -o ./bin/erutan && bin/erutan

# Or
go ruin main.go
```

# Tests

```bash
go test ~/go/src/github.com/user/erutan/utils/ -v
```

# Debug

```bash
export GRPC_VERBOSITY=INFO
```

# Roadmap

- [ ] 2D -> 3D (map procedurally generated for example)
- [ ] More (useful) characteristics (no point in adding characteristics that doesn't help survival)
- [ ] Environment-based evolution (stay near lakes, need more aquatic food, swim better idk, stay near desert, more resistant to sun ...)
- [ ] Other languages libraries (Python, JS ...) allowing either other front-ends either building bots, client-side heavy computation stuff ...
- [ ] Scalable (kd-tree, centralized state: redis ...)
- [ ] Deployable (GCP, custom, docker, kubernetes)
- [ ] A lot more
