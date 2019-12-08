/**
定义接口并实现，类似于java
**/
package main

import (
    "context"
    "log"

    "github.com/micro/go-micro"
    // 导如 protoc 自动生成的包
    pb "shippy/consignment-service/proto/consignment"
    vesselpb "shippy/vessel-service/proto/vessel"
)

const (
    DEFAULT_HOST = "localhost:27017"
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
// Repository - 虚拟存储库，模拟数据存储的使用，以后会用一个真正的实现来替换
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
    vesselClient vesselpb.VesselService  //添加vessel client
}

//
// service 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
    //检查是否有合适的货轮
    vReq := &vesselpb.Specification{
        Capacity: int32(len(req.Containers)),
        MaxWeight: req.Weight,
    }
    vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
    if err != nil {
        return err
    }
    log.Printf("found vessel: %s\n", vResp.Vessel.Name)
    req.VesselId = vResp.Vessel.Id

    // 接收承运的货物
    consignment, err := s.repo.Create(req)
    if err != nil {
        return err
    }
    //resp := &pb.Response{Created: true, Consignment: consignment}

    res.Created = true
    res.Consignment = consignment
    return nil
}

//查询consignment 服务
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	allConsignments := s.repo.GetAll()
	//resp := &pb.Response{Consignments: allConsignments}
	res.Consignments = allConsignments
    return nil
}


func main() {
    // 创建一个新服务，其中可以包括一些可选的配置

    server := micro.NewService(
        // 这个 Name 必须和 consignment.proto 中的 package 一致

        micro.Name("go.micro.srv.consignment"),
        micro.Version("latest"),
        )
    server.Init()
    repo := Repository{}

    vesselClient := vesselpb.NewVesselService("go.micro.srv.vessel", server.Client())
    pb.RegisterShippingServiceHandler(server.Server(),  &service{repo,vesselClient})
    if err := server.Run(); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}







