package entities

import "github.com/google/uuid"

type Gig struct {
	ID            uuid.UUID `gorm:"primaryKey;unique;not null"`
	FreelancerId  uuid.UUID
	Title         string
	CategoryId    int
	SkillId       int
	Description   string
	PackageTypeId int
	PackageType   PackageType `gorm:"foreignKey:PackageTypeId"`
	Price         float64
	DeliveryDays  string
}

type Images struct {
	Id    uuid.UUID `gorm:"primaryKey;unique;not null"`
	GigId uuid.UUID
	Gig   Gig `gorm:"foreignKey:GigId"`
	Image string
}

type PackageType struct {
	Id   int
	Type string
}
