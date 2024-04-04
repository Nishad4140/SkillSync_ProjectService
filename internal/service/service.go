package service

import (
	"context"
	"fmt"
	"io"

	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/helper"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/usecase"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	UserClient pb.UserServiceClient
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
	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: req.FreelancerId,
	})
	if err != nil {
		return nil, err
	}
	if freelancerData.CategoryId != req.CategoryId {
		return nil, fmt.Errorf("this is not your category")
	}
	skills, err := UserClient.FreelancerGetAllSkill(context.Background(), &pb.GetUserById{
		Id: req.FreelancerId,
	})
	if err != nil {
		return nil, err
	}
	skillMap := make(map[int32]bool)
	for {
		skill, err := skills.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		skillMap[skill.Id] = true
	}
	if !skillMap[req.SkillId] {
		return nil, fmt.Errorf("this skill is not in your skill list")
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

func (project *ProjectService) UpdateGig(ctx context.Context, req *pb.GigResponse) (*emptypb.Empty, error) {
	freelancerId, err := uuid.Parse(req.FreelancerId)
	if err != nil {
		return nil, err
	}
	gigId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: req.FreelancerId,
	})
	if err != nil {
		return nil, err
	}
	if freelancerData.CategoryId != req.CategoryId {
		return nil, fmt.Errorf("this is not your category")
	}
	skills, err := UserClient.FreelancerGetAllSkill(context.Background(), &pb.GetUserById{
		Id: req.FreelancerId,
	})
	if err != nil {
		return nil, err
	}
	skillMap := make(map[int32]bool)
	for {
		skill, err := skills.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		skillMap[skill.Id] = true
	}
	if !skillMap[req.SkillId] {
		return nil, fmt.Errorf("this skill is not in your skill list")
	}
	reqEntity := entities.Gig{
		ID:            gigId,
		FreelancerId:  freelancerId,
		Title:         req.Title,
		CategoryId:    int(req.CategoryId),
		SkillId:       int(req.SkillId),
		Description:   req.Description,
		PackageTypeId: int(req.PackageTypeId),
		Price:         float64(req.Price),
		DeliveryDays:  req.DelivaryDays,
	}
	if err := project.adapters.UpdateGig(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) GetGig(ctx context.Context, req *pb.GetById) (*pb.GigResponse, error) {
	gig, err := project.adapters.GetGigById(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.GigResponse{
		Id:            gig.ID.String(),
		Title:         gig.Title,
		FreelancerId:  gig.FreelancerId.String(),
		CategoryId:    int32(gig.CategoryId),
		SkillId:       int32(gig.SkillId),
		Description:   gig.Description,
		PackageTypeId: int32(gig.PackageTypeId),
		Price:         float32(gig.Price),
		DelivaryDays:  gig.DeliveryDays,
	}
	return res, nil
}

func (project *ProjectService) GetAllFreelancerGigs(req *pb.GetByUserId, srv pb.ProjectService_GetAllFreelancerGigsServer) error {
	gigs, err := project.adapters.GetAllFreelancerGigs(req.Id)
	if err != nil {
		return err
	}

	for _, gig := range gigs {
		res := &pb.GigResponse{
			Id:            gig.ID.String(),
			Title:         gig.Title,
			FreelancerId:  gig.FreelancerId.String(),
			CategoryId:    int32(gig.CategoryId),
			SkillId:       int32(gig.SkillId),
			Description:   gig.Description,
			PackageTypeId: int32(gig.PackageTypeId),
			Price:         float32(gig.Price),
			DelivaryDays:  gig.DeliveryDays,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (project *ProjectService) GetAllGigs(e *emptypb.Empty, srv pb.ProjectService_GetAllGigsServer) error {
	gigs, err := project.adapters.GetAllGigs()
	if err != nil {
		return err
	}

	for _, gig := range gigs {
		res := &pb.GigResponse{
			Id:            gig.ID.String(),
			Title:         gig.Title,
			FreelancerId:  gig.FreelancerId.String(),
			CategoryId:    int32(gig.CategoryId),
			SkillId:       int32(gig.SkillId),
			Description:   gig.Description,
			PackageTypeId: int32(gig.PackageTypeId),
			Price:         float32(gig.Price),
			DelivaryDays:  gig.DeliveryDays,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (project *ProjectService) ClientAddRequest(ctx context.Context, req *pb.AddClientGigRequest) (*emptypb.Empty, error) {
	clientId, err := uuid.Parse(req.ClientId)
	if err != nil {
		return nil, err
	}
	delivaryDate, err := helper.ConvertStringToDate(req.DeliveryDate)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.ClientRequest{
		ClientId:     clientId,
		Title:        req.Title,
		CategoryId:   int(req.CategoryId),
		SkillId:      int(req.SkillId),
		Description:  req.Description,
		Price:        float64(req.Price),
		DelivaryDate: delivaryDate,
	}
	if err := project.adapters.ClientAddRequest(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) ClientUpdateRequest(ctx context.Context, req *pb.ClientRequestResponse) (*emptypb.Empty, error) {
	clientId, err := uuid.Parse(req.ClientId)
	if err != nil {
		return nil, err
	}
	reqId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	delivaryDate, err := helper.ConvertStringToDate(req.DeliveryDate)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.ClientRequest{
		ID:           reqId,
		ClientId:     clientId,
		Title:        req.Title,
		CategoryId:   int(req.CategoryId),
		SkillId:      int(req.SkillId),
		Description:  req.Description,
		Price:        float64(req.Price),
		DelivaryDate: delivaryDate,
	}
	if err := project.adapters.ClientUpdateRequest(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) GetClientRequest(ctx context.Context, req *pb.GetById) (*pb.ClientRequestResponse, error) {
	clientReq, err := project.adapters.GetClientRequest(req.Id)
	if err != nil {
		return nil, err
	}
	res := &pb.ClientRequestResponse{
		Id:           clientReq.ID.String(),
		ClientId:     clientReq.ClientId.String(),
		Title:        clientReq.Title,
		CategoryId:   int32(clientReq.CategoryId),
		SkillId:      int32(clientReq.SkillId),
		Price:        float32(clientReq.Price),
		Description:  clientReq.Description,
		DeliveryDate: clientReq.DelivaryDate.String(),
	}
	return res, nil
}

func (project *ProjectService) GetAllClientRequest(req *pb.GetByUserId, srv pb.ProjectService_GetAllClientRequestServer) error {
	clientReqs, err := project.adapters.GetAllClientRequest(req.Id)
	if err != nil {
		return err
	}

	for _, clientReq := range clientReqs {
		res := &pb.ClientRequestResponse{
			Id:           clientReq.ID.String(),
			ClientId:     clientReq.ClientId.String(),
			Title:        clientReq.Title,
			CategoryId:   int32(clientReq.CategoryId),
			SkillId:      int32(clientReq.SkillId),
			Price:        float32(clientReq.Price),
			Description:  clientReq.Description,
			DeliveryDate: clientReq.DelivaryDate.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (project *ProjectService) GetAllClientRequestForFreelancers(req *pb.GetByUserId, srv pb.ProjectService_GetAllClientRequestForFreelancersServer) error {
	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: req.Id,
	})
	if err != nil {
		return err
	}
	clientReqs, err := project.adapters.GetAllClientRequestForFreelancers(int(freelancerData.CategoryId))
	if err != nil {
		return err
	}
	for _, clientReq := range clientReqs {
		res := &pb.ClientRequestResponse{
			Id:           clientReq.ID.String(),
			ClientId:     clientReq.ClientId.String(),
			Title:        clientReq.Title,
			CategoryId:   int32(clientReq.CategoryId),
			SkillId:      int32(clientReq.SkillId),
			Price:        float32(clientReq.Price),
			Description:  clientReq.Description,
			DeliveryDate: clientReq.DelivaryDate.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (project *ProjectService) AddPackageType(ctx context.Context, req *pb.AddPackageTypeRequest) (*emptypb.Empty, error) {
	packageData, err := project.adapters.GetPackgeTypeByName(req.PackageType)
	if err != nil {
		return nil, err
	}
	if packageData.Type != "" {
		return nil, fmt.Errorf("this package type already exists")
	}

	if err := project.adapters.AddPackageType(req.PackageType); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) EditPackageType(ctx context.Context, req *pb.PackageTypeResponse) (*emptypb.Empty, error) {
	packageData, err := project.adapters.GetPackgeTypeById(req.Id)
	if err != nil {
		return nil, err
	}
	if packageData.Type == "" {
		return nil, fmt.Errorf("there is no such id exists to edit")
	}
	nameCheck, err := project.adapters.GetPackgeTypeByName(req.PackageType)
	if err != nil {
		return nil, err
	}
	if nameCheck.Type != "" {
		return nil, fmt.Errorf("this package type already exists")
	}

	reqEntity := entities.PackageType{
		Id:   int(req.Id),
		Type: req.PackageType,
	}

	if err := project.adapters.EditPackgeType(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) GetPackageType(e *emptypb.Empty, srv pb.ProjectService_GetPackageTypeServer) error {
	packageTypes, err := project.adapters.GetAllPAckageType()
	if err != nil {
		return err
	}

	for _, types := range packageTypes {
		res := &pb.PackageTypeResponse{
			Id:          int32(types.Id),
			PackageType: types.Type,
		}
		err := srv.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}
