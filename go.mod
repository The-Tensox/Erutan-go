module github.com/The-Tensox/erutan

go 1.14

require (
	github.com/The-Tensox/octree v0.0.0-20200402182846-1f7eeddf526a
	github.com/The-Tensox/protometry v0.0.0-20200402182624-0a8c69d9271d
	github.com/aquilax/go-perlin v0.0.0-20191229124216-0af9ce917c28 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/sys v0.0.0-20200409092240-59c9f1ba88fa // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200410110633-0848e9f44c36 // indirect
	google.golang.org/grpc v1.28.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace (
	github.com/The-Tensox/octree => ../octree
	github.com/The-Tensox/protometry => ../protometry
)
