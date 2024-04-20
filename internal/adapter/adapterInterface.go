package adapter

import (
	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	helperstruct "github.com/Nishad4140/SkillSync_ProjectService/internal/helperStruct"
)

type AdapterInterface interface {
	CreateGig(req entities.Gig) error
	UpdateGig(req entities.Gig) error
	GetGigById(id string) (entities.Gig, error)
	GetAllFreelancerGigs(freelancerId string) ([]entities.Gig, error)
	GetAllClientRequestForFreelancers(categoryId int, queryParams helperstruct.FilterQuery) ([]entities.ClientRequest, error)
	AddIntrestToClientRequest(req entities.Intrest) error

	AddPackageType(name string) error
	EditPackgeType(req entities.PackageType) error
	GetAllPAckageType() ([]entities.PackageType, error)
	GetPackgeTypeByName(name string) (entities.PackageType, error)
	GetPackgeTypeById(id int32) (entities.PackageType, error)

	ClientAddRequest(req entities.ClientRequest) error
	ClientUpdateRequest(req entities.ClientRequest) error
	GetClientRequest(reqId string) (entities.ClientRequest, error)
	GetAllClientRequest(clientId string, queryParams helperstruct.FilterQuery) ([]entities.ClientRequest, error)
	GetAllClientRequestIntrest(reqId string) ([]entities.Intrest, error)
	GetClientIdByRequestId(reqId string) (string, error)
	ClientAddIntrestAcknowledgment(req entities.IntrestAcknowledgment) error
	GetIntrestById(id string) (entities.Intrest, error)

	GetAllGigs(queryParams helperstruct.FilterQuery) ([]entities.Gig, error)
}
