package voice

import pb "github.com/daeroworld/shared/proto/voice"

type IVoiceReader interface {
	Count(id string) (*pb.VoiceCountResponse, error)
	Retrieve(id string, idx int32) (*pb.VoiceRetrieveResponse, error)
}
