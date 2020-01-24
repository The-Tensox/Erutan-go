.PHONY: install
install: protos/realtime/realtime.pb.go
	go install .
	which erutan

protos/realtime/realtime.pb.go:
	protoc --go_out="plugins=grpc:." protos/realtime/realtime.proto

.PHONY: docker
docker:
	sudo docker build --rm -t louis030195/erutan .