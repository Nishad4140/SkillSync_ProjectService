package usecase

import "github.com/Nishad4140/SkillSync_ProtoFiles/pb"

type UsecaseInterface interface {
	UploadFreelancerFile(req *pb.FileRequest) (string, error)
}
