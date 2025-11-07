package handler

import (
	"context"
	"text/writer/internal/api/application/controller"

	pb "github.com/daeroworld/shared/proto/text"
	"github.com/golang/protobuf/ptypes/empty"
)

type GRPCHandler struct {
	pb.UnimplementedTextWriterServer
	ctrl controller.IController
}

func NewGRPCHandler(ctrl controller.IController) *GRPCHandler {
	return &GRPCHandler{ctrl: ctrl}
}

func (hdlr *GRPCHandler) Generate(ctx context.Context, req *pb.TextGenerateRequest) (*pb.TextGenerateResponse, error) {

	cnt, err := hdlr.ctrl.Create(req.GetId())

	return &pb.TextGenerateResponse{
		Id:    req.GetId(), // return the same id or some generated value
		Count: cnt,
	}, err
}

func (hdlr *GRPCHandler) UpdateContent(ctx context.Context, req *pb.UpdateContentRequest) (*empty.Empty, error) {
	return nil, hdlr.ctrl.UpdateContent(req.GetId(), req.GetContent())
}
