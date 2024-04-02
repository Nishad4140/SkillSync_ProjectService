package adapter

import "github.com/Nishad4140/SkillSync_ProjectService/entities"

type AdapterInterface interface {
	CreateGig(req entities.Gig) error

	AddPackageType(name string) error
	EditPackgeType(req entities.PackageType) error
	GetAllPAckageType() ([]entities.PackageType, error)
	GetPackgeTypeByName(name string) (entities.PackageType, error)
	GetPackgeTypeById(id int32) (entities.PackageType, error)
}
