module github.com/The-Tensox/erutan

go 1.14

require (
	github.com/The-Tensox/octree v0.0.0-20200402182846-1f7eeddf526a
	github.com/The-Tensox/protometry v0.0.0-20200402182624-0a8c69d9271d
	github.com/aquilax/go-perlin v0.0.0-20191229124216-0af9ce917c28 // indirect
	github.com/golang/protobuf v1.3.5
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200403120447-c50568487044 // indirect
	google.golang.org/grpc v1.28.0
)

replace (
	github.com/The-Tensox/octree => ../octree
	github.com/The-Tensox/protometry => ../protometry
)
