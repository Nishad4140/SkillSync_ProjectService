package usecase

import "github.com/Nishad4140/SkillSync_ProjectService/internal/adapter"

type ProjectUsecase struct {
	projectAdapter adapter.AdapterInterface
}

func NewProjectUsecase(projectAdapter adapter.AdapterInterface) *ProjectUsecase {
	return &ProjectUsecase{
		projectAdapter: projectAdapter,
	}
}
