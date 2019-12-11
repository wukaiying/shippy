package main

import (
	"context"
	"log"
	pb "shippy/consignment-service/proto/consignment"
	vesselpb "shippy/vessel-service/proto/vessel"

	"gopkg.in/mgo.v2"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
// 实现方法的传参等可参考生成的 consignment.pb.go
type handler struct {
	session *mgo.Session
	vesselClient vesselpb.VesselService
}

// 从主会话中 Clone() 出新会话处理查询
func (h *handler)GetRepo() Repository  {
	return &ConsignmentRepository{h.session.Clone()}
}

//
// service 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
func (h *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	defer h.GetRepo().Close()

	//检查是否有合适的货轮
	vReq := &vesselpb.Specification{
		Capacity: int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	// 接收承运的货物
	err = h.GetRepo().Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = req
	return nil
}

//查询consignment 服务
func (h *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	allConsignments, err := h.GetRepo().GetAll()
	if err != nil {
		return err
	}
	res.Consignments = allConsignments
	return nil
}
