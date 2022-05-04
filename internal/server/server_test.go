package server

import (
	"github.com/ArtyomArtamonov/msg/internal/mocks"
)

var jwtManagerMock *mocks.JWTManagerMock
var refreshTokenStoreMock *mocks.RefreshTokenStoreMock
var roomStoreMock *mocks.RoomStoreMock
var userStoreMock *mocks.UserStoreMock
var apiServer *ApiServer
var authServer *AuthServer

func setupTest() {
	jwtManagerMock = new(mocks.JWTManagerMock)
	refreshTokenStoreMock = new(mocks.RefreshTokenStoreMock)
	roomStoreMock = new(mocks.RoomStoreMock)
	userStoreMock = new(mocks.UserStoreMock)
	apiServer = NewApiServer(jwtManagerMock, roomStoreMock)
	authServer = &AuthServer{
		userStore:         userStoreMock,
		refreshTokenStore: refreshTokenStoreMock,
		jwtManager:        jwtManagerMock,
	}
}
