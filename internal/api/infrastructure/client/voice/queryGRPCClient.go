package voice

import (
	"context"
	"fmt"
	"log"

	pb "github.com/daeroworld/shared/proto/voice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VoiceReaderClient struct {
	clnt pb.VoiceReaderClient
	conn *grpc.ClientConn
}

func NewVoiceReaderClient(host string, port uint16) *VoiceReaderClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	clnt := pb.NewVoiceReaderClient(conn)
	return &VoiceReaderClient{clnt: clnt, conn: conn}
}

func (c *VoiceReaderClient) Shutdown() {
	c.conn.Close()
}

func (c *VoiceReaderClient) Count(jobId string) (response *pb.VoiceCountResponse, err error) {

	ctx := context.Background()

	return c.clnt.Count(ctx, &pb.VoiceCountRequest{Id: jobId})
}

func (c *VoiceReaderClient) Retrieve(jobId string, index int32) (response *pb.VoiceRetrieveResponse, err error) {

	ctx := context.Background()

	return c.clnt.Retrieve(ctx, &pb.VoiceRetrieveRequest{Id: jobId, Index: index})
}
