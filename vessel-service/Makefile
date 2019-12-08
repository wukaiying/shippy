# 一定要注意 Makefile 中的缩进(要跟一个tab，而不是4个空格)，否则执行 make build 报错
build:
	protoc --proto_path=/Users/wukaiying/go/src/shippy/vessel-service/proto/vessel:. --micro_out=/Users/wukaiying/go/src/shippy/vessel-service/proto/vessel --go_out=/Users/wukaiying/go/src/shippy/vessel-service/proto/vessel vessel.proto
	GOOS=linux GOARCH=amd64 go build
	docker build -t vessel-service .
run:
	docker run -p 50052:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
    -e MICRO_REGISTRY=mdns \
    vessel-service