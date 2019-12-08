# shippy
section 2 通过makefile build docker image 来启动服务
启动命令： make build && make run
mac 在启动cli.go时，需要运行docker-machine ip 来吧客户端中的host换成得到的ip。

#protoc -I proto/consignment --go_out=plugins=micro:proto/consignment proto/consignment/consignment.proto
#protoc -I proto/consignment --go_out=plugins=grpc:proto/consignment proto/consignment/consignment.proto
#protoc -I proto/consignment --go_out=plugins=micro:proto/consignment --micro_out=plugins=micro:proto/consignment proto/consignment/consignment.proto
#protoc -I proto/consignment --micro_out=plugins=micro:proto/consignment proto/consignment/consignment.proto

	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/shippy/consignment-service proto/consignment/consignment.proto


	protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto

应该安装 会生成两个文件
go get github.com/micro/protoc-gen-micro