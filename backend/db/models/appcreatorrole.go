package models

import "time"

type AppCreatorRoleID int

const (
	AppCreatorRoleIDPublisher AppCreatorRoleID = 1
	AppCreatorRoleIDDeveloper AppCreatorRoleID = 2
	AppCreatorRoleIDFranchise AppCreatorRoleID = 3
)

type AppCreatorRole struct {
	RoleID    AppCreatorRoleID `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;unique"`
}
