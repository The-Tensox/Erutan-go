
# Erutan-go

Trying to simulate evolution, synchronized over gRPC to clients that render a 3D visualisation.

To be used with [the Unity client](https://github.com/The-Tensox/Erutan-unity)

# Installation

```bash
# Unity project path
export UNITY_PROJECT_PATH="/home/louis/Documents/unity/Erutan"
go get github.com/The-Tensox/Erutan-go
cd $GOPATH/src/github.com/The-Tensox/Erutan-go
protoc --go_out=plugins=grpc:. protobuf/protometry/*.proto --go_opt=paths=source_relative
protoc --go_out=plugins=grpc:. protobuf/*.proto --go_opt=paths=source_relative

# If you updated the .proto, copy to unity project
cp protobuf/*.proto $UNITY_PROJECT_PATH/Assets/protobuf
cp protobuf/protometry/*.proto $UNITY_PROJECT_PATH/Assets/protobuf/protometry
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
go test ~/go/src/github.com/The-Tensox/erutan/utils/ -v
```

# Debug

```bash
export GRPC_VERBOSITY=INFO
```

# Entities

# Components

Composed of physical data (position, rotation, scale, shape, collision ...), logic + others ...

# Systems

- Network: for every entity, simply synchronize every added components over network.
- Collision: handle physics (what to do when a movement has been requested, how to handle collisions, gravity ...)
- Herbivorous, Eatable, Vegetation (will probably change name over time): some temporary hard-coded logic
- Render: how it should be rendered on clients

# Roadmap

- [ ] Better visual debugging (octree & others)
- [ ] 2D -> 3D (map procedurally generated for example)
- [ ] More (useful) characteristics (no point in adding characteristics that doesn't help survival)
- [ ] Environment-based evolution (stay near lakes, need more aquatic food, swim better idk, stay near desert, more resistant to sun ...)
- [ ] Other languages libraries (Python, JS ...) allowing either other front-ends either building bots, client-side heavy computation stuff ...
- [x] Octree
- [ ] Deployable (docker, kubernetes)
