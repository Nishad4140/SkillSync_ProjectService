package adapter

import "github.com/Nishad4140/SkillSync_ProjectService/entities"

type AdapterInterface interface {
	CreateGig(req entities.Gig) error
	UpdateGig(req entities.Gig) error
	GetGigById(id string) (entities.Gig, error)
	GetAllFreelancerGigs(freelancerId string) ([]entities.Gig, error)

	AddPackageType(name string) error
	EditPackgeType(req entities.PackageType) error
	GetAllPAckageType() ([]entities.PackageType, error)
	GetPackgeTypeByName(name string) (entities.PackageType, error)
	GetPackgeTypeById(id int32) (entities.PackageType, error)

	GetAllGigs() ([]entities.Gig, error)
}
