/**
定义接口并实现，类似于java
**/
package main

import (
    "log"
    "os"

    "github.com/micro/go-micro"
    // 导如 protoc 自动生成的包
    pb "shippy/consignment-service/proto/consignment"
    vesselpb "shippy/vessel-service/proto/vessel"
)

const (
    DEFAULT_HOST = "localhost:27017"
)

func main() {
    // 获取容器设置的数据库地址环境变量的值
    dbHost := os.Getenv("DB_HOST")
    if dbHost == ""{
        dbHost = DEFAULT_HOST
    }
    //创建session，handler中使用clone来使用session.
    session, err := CreateSession(dbHost)
    // 创建于 MongoDB 的主会话，需在退出 main() 时候手动释放连接
    defer session.Close()
    if err != nil {
        log.Fatalf("create session error: %v\n", err)
    }

    // 创建一个新服务，其中可以包括一些可选的配置
    server := micro.NewService(
        // 这个 Name 必须和 consignment.proto 中的 package 一致

        micro.Name("go.micro.srv.consignment"),
        micro.Version("latest"),
        )
    server.Init()
    vesselClient := vesselpb.NewVesselService("go.micro.srv.vessel", server.Client())
    pb.RegisterShippingServiceHandler(server.Server(),  &handler{session,vesselClient})
    if err := server.Run(); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}







