UNITY_PROJECT_PATH ?= /home/louis/Documents/unity/Erutan
NS ?= erutan
VERSION ?= 1.0.0

IMAGE_NAME ?= erutan
CONTAINER_NAME ?= erutan
CONTAINER_INSTANCE ?= default

.PHONY: help install run proto deploy_proto

help:
	@echo ''
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@echo '  install    	blabla'
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

.PHONY: dbuild drun dshell dstart dstop drm
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

drm:
	docker rm $(CONTAINER_NAME)-$(CONTAINER_INSTANCE)

default: dbuild