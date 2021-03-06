
# Erutan-go

[![Go Report Card](https://goreportcard.com/badge/github.com/The-Tensox/erutan-go?style=flat-square)](https://goreportcard.com/report/github.com/The-Tensox/erutan-go)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/The-Tensox/erutan-go)
[![Release](https://img.shields.io/github/release/The-Tensox/erutan-go.svg?style=flat-square)](https://github.com/The-Tensox/erutan-go/releases/latest)

WIP - huge refactors to be expected, some ugly code at some places :).

Trying to simulate evolution, synchronized over gRPC to clients that render a 3D visualisation.

To be used with [the Unity client](https://github.com/The-Tensox/Erutan-unity)

## Usage

### SSL/TLS configuration

```bash
# Edit your /etc/ssl/openssl.cnf add subjectAltName = IP:127.0.0.1 in [v3_ca] section

# Maybe it will do the trick but not tested :D
sed -i -e 's/#subjectAltName = IP:127.0.0.1/subjectAltName = IP:127.0.0.1/g' /etc/ssl/openssl.cnf
```

```bash
openssl genrsa -out server1.key 2048 &&
openssl req -new -x509 -sha256 -key server1.key -out server1.crt -days 3650
cp server1.crt $UNITY_PROJECT_PATH/Assets/StreamingAssets
```

### Run

```bash
go get github.com/The-Tensox/Erutan-go
cd $GOPATH/src/github.com/The-Tensox/Erutan-go
make run
```

### With Docker

```bash
make dbuild
make drun
```

You can tweak the [base configuration](config.yml.

## Tests

```bash
go test -v
```

## Monitoring

Some metrics are exposed.

![](docs/images/grafana.png)

Install and run [Grafana](https://grafana.com) + [Prometheus](https://prometheus.io/docs/introduction/overview) to monitor erutan-go:

```bash
make dmon
```

## Design

The Entity(well you will usually find them named "object" instead, but same principle) Component System(ECS) design is used.

### Components

Composed of physical data (position, rotation, scale, shape, collision ...), logic + others ...

### Systems

Sorted in execution order:

1. Collision: handle physics (what to do when a movement has been requested, how to handle collisions, gravity ...)
2. Network: for every object, simply synchronize every added components over network.
3. Logic: Herbivorous, Eatable, Vegetation (will probably change name over time): some temporary hard-coded logic

### Lifecycle

All out-going messages are handled asynchronously: throw in the channel then it will be sent to the client.

#### Main loop

1. Process incoming client packets (synchronously)
2. Update each system (by their order of priority)

#### Physics

1. Compute physics and notify of new state
2. All observers update their local state accordingly

## TODO

See issues, feel free to drop any things
