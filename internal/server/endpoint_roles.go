package server

import "github.com/ArtyomArtamonov/msg/internal/model"

type EndpointRoles = map[string][]string

func NewEndpointRoles(endpoints *Endpoints) EndpointRoles {
	return EndpointRoles{
		endpoints.ApiService.CreateRoom:      {model.ADMIN_ROLE, model.USER_ROLE},
		endpoints.ApiService.ListRooms:       {model.ADMIN_ROLE, model.USER_ROLE},
		endpoints.AuthService.Login:          nil,
		endpoints.AuthService.Register:       nil,
		endpoints.AuthService.Refresh:        nil,
		endpoints.MessageService.SendMessage: {model.ADMIN_ROLE, model.USER_ROLE},
		endpoints.MessageService.GetMessages: {model.ADMIN_ROLE, model.USER_ROLE},
	}
}
