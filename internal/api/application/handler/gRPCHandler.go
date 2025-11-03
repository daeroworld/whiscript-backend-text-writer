package handler

import (
	"context"
	"log"
	"text/writer/internal/api/application/controller"

	pb "github.com/daeroworld/shared/proto/text"
)

type GRPCHandler struct {
	pb.UnimplementedTextWriterServer
	ctrl controller.IController
}

func NewGRPCHandler(ctrl controller.IController) *GRPCHandler {
	return &GRPCHandler{ctrl: ctrl}
}

func (hdlr *GRPCHandler) Generate(ctx context.Context, req *pb.TextGenerateRequest) (*pb.TextGenerateResponse, error) {
	log.Printf("Generate called with ID: %s", req.Id)

	cnt, err := hdlr.ctrl.Create(req.GetId())

	return &pb.TextGenerateResponse{
		Id:    req.GetId(), // return the same id or some generated value
		Count: cnt,
	}, err
}
