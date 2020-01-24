
# Grpc_push

```bash
protoc -I protos/realtime --go_out=plugins=grpc:protos/realtime protos/realtime/realtime.proto
```

```bash
cp protos/realtime/realtime.proto ~/Documents/unity/Erutan/Assets/Protos/Realtime
```

```bash
go build -o ./bin/erutan
bin/erutan -s -v -p "" -h "0.0.0.0:50051"
```
