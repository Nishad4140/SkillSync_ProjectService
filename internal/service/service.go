package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/helper"
	helperstruct "github.com/Nishad4140/SkillSync_ProjectService/internal/helperStruct"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/usecase"
	"github.com/Nishad4140/SkillSync_ProjectService/kafka"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	UserClient pb.UserServiceClient
	NotiClient pb.NotificationServiceClient
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

func (project *ProjectService) GetAllGigs(req *pb.GigFilterQuery, srv pb.ProjectService_GetAllGigsServer) error {

	reqEntity := helperstruct.FilterQuery{
		Page:     int(req.Page),
		Limit:    int(req.Limit),
		Query:    req.Query,
		Filter:   req.Filter,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}
	gigs, err := project.adapters.GetAllGigs(reqEntity)
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

func (project *ProjectService) GetAllClientRequest(req *pb.RequestFilterQuery, srv pb.ProjectService_GetAllClientRequestServer) error {
	reqEntity := helperstruct.FilterQuery{
		Page:     int(req.Page),
		Limit:    int(req.Limit),
		Query:    req.Query,
		Filter:   req.Filter,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}
	clientReqs, err := project.adapters.GetAllClientRequest(req.UserId, reqEntity)
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

func (project *ProjectService) GetAllClientRequestForFreelancers(req *pb.RequestFilterQuery, srv pb.ProjectService_GetAllClientRequestForFreelancersServer) error {
	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: req.UserId,
	})
	if err != nil {
		return err
	}
	reqEntity := helperstruct.FilterQuery{
		Page:     int(req.Page),
		Limit:    int(req.Limit),
		Query:    req.Query,
		Filter:   req.Filter,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}
	clientReqs, err := project.adapters.GetAllClientRequestForFreelancers(int(freelancerData.CategoryId), reqEntity)
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

func (project *ProjectService) ShowIntrest(ctx context.Context, req *pb.IntrestRequest) (*emptypb.Empty, error) {
	freelancerId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	reqId, err := uuid.Parse(req.RequestId)
	if err != nil {
		return nil, err
	}
	gigId, err := uuid.Parse(req.GigId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Intrest{
		ClientRequestId: reqId,
		FreelancerId:    freelancerId,
		GigId:           gigId,
	}

	if err := project.adapters.AddIntrestToClientRequest(reqEntity); err != nil {
		return nil, err
	}
	clientId, err := project.adapters.GetClientIdByRequestId(req.RequestId)
	if err != nil {
		return nil, err
	}

	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	clientRequest, err := project.adapters.GetClientRequest(req.RequestId)
	if err != nil {
		return nil, err
	}

	if _, err := NotiClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
		UserId:       clientId,
		Notification: fmt.Sprintf(`{"message" : "%s intrested on your %s titled request"}`, freelancerData.Name, clientRequest.Title),
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (project *ProjectService) GetAllIntrest(req *pb.GetAllIntrestRequest, srv pb.ProjectService_GetAllIntrestServer) error {
	reqs, err := project.adapters.GetAllClientRequestIntrest(req.RequestId)
	if err != nil {
		return err
	}

	for _, req := range reqs {
		freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
			Id: req.FreelancerId.String(),
		})
		if err != nil {
			return err
		}

		res := &pb.IntrestResponse{
			Id:        req.Id.String(),
			UserId:    req.FreelancerId.String(),
			Name:      freelancerData.Name,
			RequestId: req.ClientRequestId.String(),
			GigId:     req.GigId.String(),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (project *ProjectService) ClientIntrestAcknowledgment(ctx context.Context, req *pb.IntrestAcknowledgmentRequest) (*emptypb.Empty, error) {
	clientId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	intrestId, err := uuid.Parse(req.IntrestId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.IntrestAcknowledgment{
		ClientId:  clientId,
		IntrestId: intrestId,
	}
	if err := project.adapters.ClientAddIntrestAcknowledgment(reqEntity); err != nil {
		return nil, err
	}
	intrest, err := project.adapters.GetIntrestById(req.IntrestId)
	if err != nil {
		return nil, err
	}

	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: intrest.FreelancerId.String(),
	})
	if err != nil {
		return nil, err
	}

	clientData, err := UserClient.GetClientById(context.Background(), &pb.GetUserById{
		Id: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	Clinetreq, err := project.adapters.GetClientRequest(intrest.ClientRequestId.String())
	if err != nil {
		return nil, err
	}
	if err := kafka.IntrestAcknowledgment(freelancerData.Email, clientData.Name, Clinetreq.Title); err != nil {
		log.Println("error happened in pushing to stream", err)
	}
	return nil, nil
}

func (project *ProjectService) CreateProject(ctx context.Context, req *pb.ProjectRequest) (*pb.ProjectResponse, error) {
	freelancerId, err := uuid.Parse(req.FreelancerId)
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	freelancerData, err := UserClient.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: req.FreelancerId,
	})
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	if freelancerData.Name == "" {
		return &pb.ProjectResponse{}, fmt.Errorf("thre is no such frelancer")
	}

	clientId, err := uuid.Parse(req.ClientId)
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	clientData, err := UserClient.GetClientById(context.Background(), &pb.GetUserById{
		Id: req.ClientId,
	})
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	if clientData.Name == "" {
		return &pb.ProjectResponse{}, fmt.Errorf("thre is no such client")
	}

	gigId, err := uuid.Parse(req.GigId)
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	gigData, err := project.adapters.GetGigById(req.GigId)
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	if gigData.Title == "" {
		return &pb.ProjectResponse{}, fmt.Errorf("thre is no such gig")
	}

	endDate, err := helper.ConvertStringToDate(req.EndDate)
	if err != nil {
		return &pb.ProjectResponse{}, err
	}
	id := uuid.New()
	reqEntity := entities.Project{
		Id:           id,
		ClientId:     clientId,
		FreelancerId: freelancerId,
		GigId:        gigId,
		EndDate:      endDate,
		Price:        gigData.Price,
	}
	proj, err := project.adapters.CreateProject(reqEntity)
	if err != nil {
		return &pb.ProjectResponse{}, err
	}

	if _, err := NotiClient.AddNotification(context.Background(), &pb.AddNotificationRequest{
		UserId:       req.FreelancerId,
		Notification: fmt.Sprintf(`{"message" : "%s invited you %s project"}`, clientData.Name, id),
	}); err != nil {
		return &pb.ProjectResponse{}, err
	}

	status := proj.Status

	_, err = strconv.Atoi(proj.Status)
	if err == nil {
		status += "%"
	}

	res := &pb.ProjectResponse{
		Id:           proj.Id.String(),
		ClientId:     proj.ClientId.String(),
		FreelancerId: proj.FreelancerId.String(),
		GigId:        proj.GigId.String(),
		StartDate:    proj.StartDate.Format("02-01-2006"),
		EndDate:      proj.EndDate.Format("02-01-2006"),
		Status:       status,
		IsPaid:       false,
		Price:        float32(proj.Price),
	}
	return res, nil
}

func (project *ProjectService) UpdateProject(ctx context.Context, req *pb.ProjectResponse) (*emptypb.Empty, error) {
	clientId, err := uuid.Parse(req.ClientId)
	if err != nil {
		return nil, err
	}
	projectId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	projectData, err := project.adapters.GetProject(req.Id)
	if err != nil {
		return nil, err
	}
	if projectData.Status == "" {
		return nil, fmt.Errorf("thre is no such project")
	}
	clientData, err := UserClient.GetClientById(context.Background(), &pb.GetUserById{
		Id: req.ClientId,
	})
	if err != nil {
		return nil, err
	}
	if clientData.Name == "" {
		return nil, fmt.Errorf("thre is no such client")
	}

	gigId, err := uuid.Parse(req.GigId)
	if err != nil {
		return nil, err
	}
	gigData, err := project.adapters.GetGigById(req.GigId)
	if err != nil {
		return nil, err
	}
	if gigData.Title == "" {
		return nil, fmt.Errorf("thre is no such gig")
	}

	endDate, err := helper.ConvertStringToDate(req.EndDate)
	if err != nil {
		return nil, err
	}

	reqEntity := entities.Project{
		Id:       projectId,
		ClientId: clientId,
		GigId:    gigId,
		EndDate:  endDate,
	}
	if err := project.adapters.UpdateProject(reqEntity); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) GetProject(ctx context.Context, req *pb.GetProjectById) (*pb.ProjectResponse, error) {
	proj, err := project.adapters.GetProject(req.Id)
	if err != nil {
		return nil, err
	}
	formattedStartDate := proj.StartDate.Format("2006-01-02")
	formattedEndDate := proj.EndDate.Format("2006-01-02")

	if proj.ClientId.String() != req.UserId && proj.FreelancerId.String() != req.UserId {
		return &pb.ProjectResponse{}, fmt.Errorf("there is no such project for you")
	}

	status := proj.Status

	_, err = strconv.Atoi(proj.Status)
	if err == nil {
		status += "%"
	}

	res := &pb.ProjectResponse{
		Id:                 proj.Id.String(),
		ClientId:           proj.ClientId.String(),
		FreelancerId:       proj.FreelancerId.String(),
		GigId:              proj.GigId.String(),
		StartDate:          formattedStartDate,
		EndDate:            formattedEndDate,
		Status:             status,
		IsManagementNeeded: proj.IsManagement,
		IsPaid:             proj.IsPaid,
		Price:              float32(proj.Price),
	}
	return res, nil
}

func (project *ProjectService) RemoveProject(ctx context.Context, req *pb.GetProjectById) (*emptypb.Empty, error) {
	proj, err := project.adapters.GetProject(req.Id)
	if err != nil {
		return nil, err
	}

	if proj.ClientId.String() != req.UserId {
		return nil, fmt.Errorf("you dont have the authority to delete this project")
	}
	if proj.Status != "not started" {
		return nil, fmt.Errorf("the project is already started, you cant delete")
	}
	if err := project.adapters.RemoveProjects(req.Id); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) GetAllProjects(req *pb.ProjectResponse, srv pb.ProjectService_GetAllProjectsServer) error {

	// reqEntity := helperstruct.FilterQuery{
	// 	Page:     int(req.Page),
	// 	Limit:    int(req.Limit),
	// 	Query:    req.Query,
	// 	Filter:   req.Filter,
	// 	SortBy:   req.SortBy,
	// 	SortDesc: req.SortDesc,
	// }
	var clientId uuid.UUID
	var freelancerId uuid.UUID
	var err error
	if req.ClientId != "" {
		clientId, err = uuid.Parse(req.ClientId)
		if err != nil {
			return err
		}
	} else {
		freelancerId, err = uuid.Parse(req.FreelancerId)
		if err != nil {
			return err
		}
	}

	reqEntity := entities.Project{
		ClientId:     clientId,
		FreelancerId: freelancerId,
	}
	gigs, err := project.adapters.GetAllProjects(reqEntity)
	if err != nil {
		return err
	}

	for _, gig := range gigs {
		res := &pb.ProjectResponse{
			Id:                 gig.Id.String(),
			FreelancerId:       gig.FreelancerId.String(),
			ClientId:           gig.ClientId.String(),
			GigId:              gig.GigId.String(),
			StartDate:          gig.StartDate.String(),
			EndDate:            gig.EndDate.String(),
			Status:             gig.Status,
			Price:              float32(gig.Price),
			IsManagementNeeded: gig.IsManagement,
			IsPaid:             gig.IsPaid,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (project *ProjectService) FreelancerProjectAcknowledgment(ctx context.Context, req *pb.ProjectAcknowledgmentRequest) (*emptypb.Empty, error) {
	freelancerId, err := uuid.Parse(req.FreelancerId)
	if err != nil {
		return nil, err
	}
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}
	reqEntity := entities.Project{
		Id:           projectId,
		FreelancerId: freelancerId,
		Status:       "0",
	}
	if err := project.adapters.FreelancerUpdateStatus(reqEntity); err != nil {
		return nil, err
	}

	return nil, nil
}

func (project *ProjectService) FreelancerProjectManagement(ctx context.Context, req *pb.ManagementRequest) (*emptypb.Empty, error) {
	managementId, err := project.adapters.ProjectManagement(helperstruct.ProjectManagement{
		IsManagementNeeded: req.IsManagementNeeded,
		ModuleNumber:       int(req.ModuleNumber),
		ProjectId:          req.Projectid,
	})
	if err != nil {
		return nil, err
	}
	if managementId != "" {
		for _, module := range req.ModuleDetails {
			moduleDetails := []string{module.ModuleName, module.Description}
			if err := project.adapters.ModuleManagement(helperstruct.ModuleManagement{
				ModuleDetails: moduleDetails,
				ManagementId:  managementId,
			}); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (project *ProjectService) FreelancerModuleUpdation(ctx context.Context, req *pb.ModuleUpdation) (*emptypb.Empty, error) {

	managementData, err := project.adapters.GetManagementByProjectId(req.ProjectId)
	if err != nil {
		return nil, err
	}
	moduleId, err := uuid.Parse(req.ModuleId)
	if err != nil {
		return nil, err
	}
	management, err := uuid.Parse(managementData.Id.String())
	if err != nil {
		return nil, err
	}

	if err := project.adapters.ModuleStatusUpdate(entities.Module{
		Id:           moduleId,
		Status:       int(req.Status),
		ManagementId: management,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (project *ProjectService) FreelancerUpdateProjectStatus(ctx context.Context, req *pb.UpdateProjectStatusRequest) (*emptypb.Empty, error) {
	freelancerId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	proj, err := project.adapters.GetProject(req.ProjectId)
	if err != nil {
		return nil, err
	}

	if proj.IsManagement {
		return nil, fmt.Errorf("you selected the project management, please update the modules")
	}

	reqEntity := entities.Project{
		Id:           projectId,
		FreelancerId: freelancerId,
		Status:       req.Status,
	}
	if err := project.adapters.FreelancerUpdateStatus(reqEntity); err != nil {
		return nil, err
	}

	return nil, nil
}

func (project *ProjectService) GetProjectManagement(ctx context.Context, req *pb.GetProjectById) (*pb.ProjectManagementResponse, error) {
	managementData, err := project.adapters.GetManagementByProjectId(req.Id)
	if err != nil {
		return nil, err
	}

	moduleData, err := project.adapters.GetModuleByManagementId(managementData.Id.String())
	if err != nil {
		return nil, err
	}

	moduleDetails := []*pb.Module{}

	for _, v := range moduleData {
		module := &pb.Module{
			Id:          v.Id.String(),
			ModuleName:  v.ModuleName,
			Description: v.ModuleDescription,
			Status:      int32(v.Status),
		}
		moduleDetails = append(moduleDetails, module)
	}

	res := &pb.ProjectManagementResponse{
		ProjectId:     managementData.ProjectId.String(),
		ManagementId:  managementData.Id.String(),
		ModuleNumber:  int32(managementData.ModulesNumber),
		ModuleDetails: moduleDetails,
	}

	return res, nil
}

func (project *ProjectService) PaymentStatusChange(ctx context.Context, req *pb.PaymentStatusRequest) (*emptypb.Empty, error) {
	if err := project.adapters.UpdatePaymentStatus(req.ProjectId); err != nil {
		return nil, err
	}
	return nil, nil
}

func (project *ProjectService) FreelancerUploadFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	proj, err := project.adapters.GetProject(req.ProjectId)
	if err != nil {
		return &pb.FileResponse{}, err
	}
	if proj.Status != "completed" {
		return nil, fmt.Errorf("complete teh project")
	}
	freelancerId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	if proj.FreelancerId != freelancerId {
		return nil, fmt.Errorf("you cant edit on this")
	}
	url, err := project.usecase.UploadFreelancerFile(req)
	res := &pb.FileResponse{
		Url: url,
	}
	return res, nil
}

func (project *ProjectService) GetFile(ctx context.Context, req *pb.GetProjectById) (*pb.FileResponse, error) {
	proj, err := project.adapters.GetProject(req.Id)
	if err != nil {
		return &pb.FileResponse{}, err
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	if proj.ClientId == userId {
		if !proj.IsPaid {
			return nil, fmt.Errorf("complete the payment")
		}
	}
	file, err := project.adapters.GetFile(proj.Id.String())
	res := &pb.FileResponse{
		Url: file.File,
	}
	return res, nil
}
