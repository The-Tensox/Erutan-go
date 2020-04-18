UNITY_PROJECT_PATH ?= /home/louis/Documents/unity/Erutan
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
	@echo '  dbuild    	build docker image'
	@echo '  drun    	run docker container'
	@echo '  dshell    	open a shell into the container'
	@echo '  dstart    	start the container'
	@echo '  dstop    	stop the container'
	@echo ''

install: proto
	go install .
	which erutan

run:
	go run .

proto:
	protoc --go_out=plugins=grpc:. protobuf/protometry/*.proto --go_opt=paths=source_relative
	protoc --go_out=plugins=grpc:. protobuf/*.proto --go_opt=paths=source_relative

deploy_proto:
	cp protobuf/*.proto $(UNITY_PROJECT_PATH)/Assets/protobuf
	cp protobuf/protometry/*.proto $(UNITY_PROJECT_PATH)/Assets/protobuf/protometry

.PHONY: dbuild drun dshell dstart dstop drm dmon
dbuild:
	docker build -t $(NS)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .

drun:
	docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION)

dshell:
	docker run --rm --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) -i -t $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION) /bin/bash

dstart:
	docker run -d --name $(CONTAINER_NAME)-$(CONTAINER_INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(IMAGE_NAME):$(VERSION)

dstop:
	docker stop $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

dmon:
	docker run -d --rm --name prom -p 9090:9090 -v $(pwd)/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
	docker run -d --rm --name graf -p 3000:3000 grafana/grafana

default: dbuild
