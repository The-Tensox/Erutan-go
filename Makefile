UNITY_PROJECT_PATH ?= $(HOME)/Documents/unity/Erutan-unity
NS ?= erutan
VERSION ?= 1.0.0
PORTS = -p 34555:34555 -p 50051:50051
IMAGE_NAME ?= erutan
CONTAINER_NAME ?= erutan
CONTAINER_INSTANCE ?= default

.PHONY: help install run proto deploy_proto

help:
	@echo ''
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@echo '  install    	compile protos and go project'
	@echo '  run    	run the go project'
	@echo '  proto    	compile protos'
	@echo '  deploy_proto	copy protos to client project'
	@echo '  docker_build    	build docker image'
	@echo '  docker_run    	run docker container'
	@echo '  docker_shell    	open a shell into the container'
	@echo '  docker_start    	start the container'
	@echo '  docker_stop    	stop the container'
	@echo ''

install: proto
	go install cmd/server
	which erutan

run:
	go run cmd/server/main.go

proto:
	protoc --go_out=plugins=grpc:. protobuf/protometry/*.proto --go_opt=paths=source_relative
	protoc --go_out=plugins=grpc:. protobuf/erutan/*.proto

deploy_proto:
	cp protobuf/*.proto $(UNITY_PROJECT_PATH)/Assets/protobuf
	cp protobuf/protometry/*.proto $(UNITY_PROJECT_PATH)/Assets/protobuf/protometry

.PHONY: docker_build docker_run docker_shell docker_start docker_stop docker_mon
docker_build:
	docker build -t $(NS)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .

docker_run:
	docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION)

docker_shell:
	docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) -i -t $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION) /bin/bash

docker_start:
	docker run -d --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION)

docker_stop:
	docker stop $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

docker_mon:
	docker run -d --rm --name prom -p 9090:9090 -v $(pwd)/mon/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
	docker run -d --rm --name graf -p 3000:3000 grafana/grafana

default: docker_build
