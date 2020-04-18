
# Erutan-go

Trying to simulate evolution, synchronized over gRPC to clients that render a 3D visualisation.

To be used with [the Unity client](https://github.com/The-Tensox/Erutan-unity)

## Usage

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

## SSL/TLS configuration

```bash
# Edit your /etc/ssl/openssl.cnf add subjectAltName = IP:127.0.0.1 in [v3_ca] section

```bash
# Maybe it will do the trick but not tested :D
sed -i -e 's/#subjectAltName = IP:127.0.0.1/subjectAltName = IP:127.0.0.1/g' /etc/ssl/openssl.cnf
```

```bash
openssl genrsa -out server1.key 2048 &&
openssl req -new -x509 -sha256 -key server1.key -out server1.crt -days 3650
cp server1.crt $UNITY_PROJECT_PATH/Assets/StreamingAssets
```

## Tests

```bash
go test -v
```

## Monitoring

Some metrics are exposed.

![](docs/images/grafana.png)

Install and run [Grafana](https://grafana.com) + [Prometheus](https://prometheus.io/docs/introduction/overview) to monitor erutan-go:

```bash
docker run -d --rm --name prom -p 9090:9090 -v $(pwd)/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
docker run -d --rm --name graf -p 3000:3000 grafana/grafana

# Or simply
make dmon
```

## ECS

### Entities

### Components

Composed of physical data (position, rotation, scale, shape, collision ...), logic + others ...

### Systems

- Network: for every entity, simply synchronize every added components over network.
- Collision: handle physics (what to do when a movement has been requested, how to handle collisions, gravity ...)
- Herbivorous, Eatable, Vegetation (will probably change name over time): some temporary hard-coded logic
- Render: how it should be rendered on clients

### TODO

- [x] Better visual debugging (octree & others)
- [ ] More (useful) characteristics (no point in adding characteristics that doesn't help survival)
- [ ] Environment-based evolution (stay near lakes, need more aquatic food, swim better idk, stay near desert, more resistant to sun ...)
- [ ] Other languages libraries (Python, JS ...) allowing either other front-ends either building bots, client-side heavy computation stuff ...
- [x] Octree
