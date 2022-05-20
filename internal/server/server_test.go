package server

import (
	"github.com/ArtyomArtamonov/msg/internal/mocks"
	"github.com/ArtyomArtamonov/msg/internal/utils"
)

var jwtManagerMock *mocks.JWTManagerMock
var refreshTokenStoreMock *mocks.RefreshTokenStoreMock
var roomStoreMock *mocks.RoomStoreMock
var userStoreMock *mocks.UserStoreMock
var amqpProducerMock *mocks.AMQPProducerMock
var apiServer *ApiServer
var authServer *AuthServer

func setupTest() {
	utils.MockNow(utils.DefaultMockTime)
	jwtManagerMock = new(mocks.JWTManagerMock)
	refreshTokenStoreMock = new(mocks.RefreshTokenStoreMock)
	roomStoreMock = new(mocks.RoomStoreMock)
	userStoreMock = new(mocks.UserStoreMock)
	amqpProducerMock = new(mocks.AMQPProducerMock)
	apiServer = NewApiServer(jwtManagerMock, roomStoreMock, amqpProducerMock)
	authServer = &AuthServer{
		userStore:         userStoreMock,
		refreshTokenStore: refreshTokenStoreMock,
		jwtManager:        jwtManagerMock,
	}
}
