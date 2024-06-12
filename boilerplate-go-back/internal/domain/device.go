package domain

import "time"

type Device struct {
	Id              uint64
	OrganizationId  uint64
	RoomId          uint64
	InventoryNumber string
	SerialNumber    string
	Characteristics string
	Category        string
	Units           string
	CreatedDate     time.Time
	UpdatedDate     time.Time
	DeletedDate     *time.Time
}
