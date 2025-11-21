package handler

import (
	"context"
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

	cnt, err := hdlr.ctrl.Create(req.GetId())

	return &pb.TextGenerateResponse{
		Id:    req.GetId(), // return the same id or some generated value
		Count: cnt,
	}, err
}

func (hdlr *GRPCHandler) Put(ctx context.Context, req *pb.TextUpdateRequest) (*pb.TextUpdateResponse, error) {
	id, err := hdlr.ctrl.Put(req.GetId(), req.GetFilename(), int(req.GetSentence()), int(req.GetWord()), req.GetStart(), req.GetEnd(), req.GetContent())
	return &pb.TextUpdateResponse{
		Id: id,
	}, err
}
