module github.com/The-Tensox/erutan

go 1.14

require (
	github.com/The-Tensox/octree v0.0.0-20200401181246-11a4110a7917
	github.com/The-Tensox/protometry v0.0.0-20200402182624-0a8c69d9271d
	github.com/aquilax/go-perlin v0.0.0-20191229124216-0af9ce917c28
	github.com/golang/protobuf v1.3.5
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	google.golang.org/grpc v1.28.0
)

replace (
	github.com/The-Tensox/octree v0.0.0-20200401181246-11a4110a7917 => ../octree
	github.com/The-Tensox/protometry v0.0.0-20200329160116-a97dbae83b93 => ../protometry
)
