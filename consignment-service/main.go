/**
定义接口并实现，类似于java
**/
package main

import (
    "context"
    "github.com/micro/go-micro/util/log"

    // 导如 protoc 自动生成的包
    pb "shippy/consignment-service/proto/consignment"
    "github.com/micro/go-micro"
)

const (
    PORT = ":50051"
)

//
// 仓库接口
//
type IRepository interface {
    Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
    GetAll() []*pb.Consignment
}

//
// 我们存放多批货物的仓库，实现了 IRepository 接口
//
type Repository struct {
    consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    repo.consignments = append(repo.consignments, consignment)
    return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
    return repo.consignments
}

//
// 定义微服务
//
type service struct {
    repo Repository
}

//
// service 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
    // 接收承运的货物
    consignment, err := s.repo.Create(req)
    if err != nil {
        return nil, err
    }
    resp := &pb.Response{Created: true, Consignment: consignment}
    return resp, nil
}

//查询consignment 服务
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	allConsignments := s.repo.GetAll()
	resp := &pb.Response{Consignments: allConsignments}
    return resp, nil
}


func main() {
    server := micro.NewService(
        micro.Name("go.micro.srv.consignment"),
        micro.Version("latest"),
        )
    server.Init()
    repo := Repository{}
    pb.RegisterShippingServiceHandler(server.Server(), &service{repo})
    if err := server.Run(); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}







