package service

import (
	"context"

	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/usecase"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProjectService struct {
	adapters adapter.AdapterInterface
	usecase  usecase.UsecaseInterface
	pb.UnimplementedProjectServiceServer
}

func NewProjectService(adapters adapter.AdapterInterface, usecase usecase.UsecaseInterface) *ProjectService {
	return &ProjectService{
		adapters: adapters,
		usecase:  usecase,
	}
}

func (project *ProjectService) CreateGig(ctx context.Context, req *pb.CreateGigRequest) (*emptypb.Empty, error) {
	freelancerId, err := uuid.Parse(req.FreelancerId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Gig{
		FreelancerId:  freelancerId,
		Title:         req.Title,
		CategoryId:    int(req.CategoryId),
		SkillId:       int(req.SkillId),
		Description:   req.Description,
		PackageTypeId: int(req.PackageTypeId),
		Price:         float64(req.Price),
		DeliveryDays:  req.DelivaryDays,
	}
	if err := project.adapters.CreateGig(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}
