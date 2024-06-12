package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type RoomController struct {
	roomService app.RoomService
}

func NewRoomController(os app.RoomService) RoomController {
	return RoomController{
		roomService: os,
	}
}

func (c RoomController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		org := r.Context().Value(OrgKey).(domain.Room)
		room, err := requests.Bind(r, requests.RoomRequest{}, domain.Room{})
		if err != nil {
			log.Printf("RoomController: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		room.OrganizationId = org.Id
		room, err = c.roomService.Save(room, user.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		var roomDto resources.RoomDto
		Success(w, roomDto.DomainToDto(room))
	}
}

func (c OrganizationController) FindByOrgId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		orgs, err := c.organizationService.FindForUser(user.Id)
		if err != nil {
			log.Printf("OrganizationController: %s", err)
			InternalServerError(w, err)
			return
		}

		var orgsDto resources.OrgsDto
		resposnse := orgsDto.DomainToDto(orgs)
		Success(w, resposnse)
	}
}

func (c RoomController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room := r.Context().Value(RoomKey).(domain.Room)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if room.OrganizationId != org.Id {
			err := fmt.Errorf("access denied: organization %d does not own room %d", org.Id, room.Id)
			Forbidden(w, err)
			return
		}

		var roomDto resources.RoomDto
		Success(w, roomDto.DomainToDto(room))
	}
}

func (r RoomController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		org, err := requests.Bind(r, requests.OrganizationRequest{}, domain.Organization{})
		if err != nil {
			log.Printf("RoomController: %s", err)
			BadRequest(w, err)
			return
		}

		organization := r.Context().Value(OrgKey).(domain.Organization)
		if organization.UserId != user.Id {
			err := fmt.Errorf("access deied", org.Id, user.Id)
			Forbidden(w, err)
			return
		}

		organization.Name = org.Name
		organization.Address = org.Address
		organization.City = org.City
		organization.Lat = org.Lat
		organization.Lon = org.Lon

		//org.UserId = user.Id
		//org, err = c.organizationService.Save(org)

		if err != nil {
			log.Printf("OrganizationController: %s", err)
			InternalServerError(w, err)
			return
		}

		var orgDto resources.OrgDto
		Success(w, orgDto.DomainToDto(organization))
	}
}

func (c RoomController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room := r.Context().Value(RoomKey).(domain.Room)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if room.OrganizationId != org.Id {
			err := fmt.Errorf("access denied: organization %d does not own room %d", org.Id, room.Id)
			Forbidden(w, err)
			return
		}

		err := c.roomService.Delete(room.Id)
		if err != nil {
			log.Printf("RoomController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
