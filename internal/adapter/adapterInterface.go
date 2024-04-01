package adapter

import "github.com/Nishad4140/SkillSync_ProjectService/entities"

type AdapterInterface interface {
	CreateGig(req entities.Gig) error
}
