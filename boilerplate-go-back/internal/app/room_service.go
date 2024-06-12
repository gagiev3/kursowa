package app

import (
	"errors"
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type RoomService interface {
	Save(r domain.Room, uId uint64) (domain.Room, error)
	FindByOrgId(orgId uint64) ([]domain.Room, error)
	Find(id uint64) (interface{}, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
	// Save(o domain.Organization) (domain.Organization, error)
	// FindByOrgId(uId uint64) ([]domain.Organization, error)
	// Find(id uint64) (interface{}, error)
	// Update(r domain.Room) (domain.Room, error)
	// Delete(id uint64) error
}

type roomService struct {
	orgRepo  database.OrganizationRepository
	roomRepo database.RoomRepository
}

func NewRoomService(
	or database.OrganizationRepository,
	rr database.RoomRepository) OrganizationService {
	return organizationService{
		organizationRepo: or,
		roomRepo:         rr,
	}
}

func (s roomService) Save(r domain.Room, uId uint64) (domain.Room, error) {
	org, err := s.orgRepo.FindById(r.OrganizationId)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	if org.UserId != uId {
		err = errors.New("access denied")
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	r, err = s.roomRepo.Save(r)

	return r, err
}

func (s roomService) FindByOrgId(orgId uint64) ([]domain.Room, error) {
	rooms, err := s.roomRepo.FindByOrgId(orgId)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return nil, err
	}

	return rooms, err
}

func (s roomService) Find(id uint64) (interface{}, error) {
	room, err := s.roomRepo.FindById(id)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return room, err
}

func (s roomService) Update(r domain.Room) (domain.Room, error) {
	_, err := s.orgRepo.FindById(r.OrganizationId)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	room, err := s.roomRepo.Update(r)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return domain.Room{}, err
	}

	return room, err
}

func (s roomService) Delete(id uint64) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		log.Printf("RoomService: %s", err)
		return err
	}

	return nil
}
