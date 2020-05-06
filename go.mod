module github.com/The-Tensox/Erutan-go

go 1.14

require (
	github.com/The-Tensox/octree v0.0.0-20200502124658-d5eedbdf3820
	github.com/The-Tensox/protometry v0.0.0-20200502124743-c5fd69c974e2
	github.com/aquilax/go-perlin v0.0.0-20191229124216-0af9ce917c28
	github.com/golang/protobuf v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/procfs v0.0.11 // indirect
	golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.21.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace github.com/The-Tensox/protometry => ../protometry
replace github.com/The-Tensox/octree => ../octree
