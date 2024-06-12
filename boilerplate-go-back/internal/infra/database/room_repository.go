package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const RoomsTableName = "rooms"

type room struct {
	Id             uint64     `db:"id,omitempty"`
	OrganizationId uint64     `db:"organization_id"`
	Name           string     `db:"name"`
	Description    string     `db:"description"`
	CreatedDate    time.Time  `db:"created_date"`
	UpdatedDate    time.Time  `db:"updated_date"`
	DeletedDate    *time.Time `db:"deleted_date"`
}

type RoomRepository interface {
	Save(dr domain.Room) (domain.Room, error)
	FindByOrgId(oId uint64) ([]domain.Room, error)
	FindById(id uint64) (domain.Room, error)
	Update(o domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

type roomRepository struct {
	coll db.Collection
	sess db.Session
}

func NewRoomRepository(dbSession db.Session) RoomRepository {
	return roomRepository{
		coll: dbSession.Collection(RoomsTableName),
		sess: dbSession,
	}
}

func (r roomRepository) Save(dr domain.Room) (domain.Room, error) {
	rm := r.mapDomainToModel(dr)
	rm.CreatedDate, rm.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&rm)
	if err != nil {
		return domain.Room{}, err
	}
	dr = r.mapModelToDomain(rm)
	return dr, nil
}

func (r roomRepository) FindByOrgId(oId uint64) ([]domain.Room, error) {
	var rs []room
	err := r.coll.Find(db.Cond{"organization_id": oId, "deleted_date": nil}).All(&rs)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(rs)
	return res, nil
}

func (r roomRepository) FindById(id uint64) (domain.Room, error) {
	var dr room
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&dr)
	if err != nil {
		return domain.Room{}, err
	}
	rs := r.mapModelToDomain(dr)
	return rs, nil
}

func (r roomRepository) Update(rs domain.Room) (domain.Room, error) {
	dr := r.mapDomainToModel(rs)
	dr.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": dr.Id, "deleted_date": nil}).Update(&dr)
	if err != nil {
		return domain.Room{}, err
	}
	rs = r.mapModelToDomain(dr)
	return rs, nil
}

func (r roomRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r roomRepository) mapDomainToModel(d domain.Room) room {
	return room{
		Id:             d.Id,
		OrganizationId: d.OrganizationId,
		Name:           d.Name,
		Description:    d.Description,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomain(d room) domain.Room {
	return domain.Room{
		Id:             d.Id,
		OrganizationId: d.OrganizationId,
		Name:           d.Name,
		Description:    d.Description,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomainCollection(rs []room) []domain.Room {
	var rooms []domain.Room
	for _, dr := range rs {
		dr := r.mapModelToDomain(dr)
		rooms = append(rooms, dr)
	}
	return rooms
}
