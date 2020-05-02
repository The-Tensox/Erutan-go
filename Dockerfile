FROM golang:1.14-alpine as builder
RUN apk --no-cache add git

RUN go get -d \
    github.com/pkg/errors \
    golang.org/x/net/context \
    google.golang.org/grpc \
    github.com/golang/protobuf/ptypes

WORKDIR /go/src/github.com/The-Tensox/erutan
COPY . .

#RUN go build -o app .
RUN go install .

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/erutan

# Prometheus metrics
EXPOSE 34555

# gRPC
EXPOSE 50051

## --- Execution Stage
#
#FROM alpine:latest
#EXPOSE 34555/tcp
#
#WORKDIR /root/
#COPY --from=builder /go/src/github.com/The-Tensox/Erutan-go/app .
#
#ENTRYPOINT ["./app"]
