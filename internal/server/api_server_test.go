package server

import (
	"context"
	"errors"
	"testing"

	"github.com/ArtyomArtamonov/msg/internal/model"
	proto "github.com/ArtyomArtamonov/msg/internal/server/msg-proto"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestApiServer_CreateRoomSuccess(t *testing.T) {
	setupTest()

	createRoomRequest := &proto.CreateRoomRequest{
		Name: "room_name",
		UserIds: []string{
			"f13550bb-ea7d-4c5b-b9a4-a352c13f1232",
		},
	}
	
	expectedClaimsResult := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id: "7815f165-6e48-452d-8fe9-b4075b35e194",
		},
	}
	expectedResponse := proto.CreateRoomStatus{
		Name:   "room_name",
		Users:  []string{
			"7815f165-6e48-452d-8fe9-b4075b35e194",
			"f13550bb-ea7d-4c5b-b9a4-a352c13f1232",
		},
	}

	jwtManagerMock.On("GetAndVerifyClaims", mock.Anything).Return(expectedClaimsResult, nil)
	roomStoreMock.On("Add", mock.Anything).Return(nil)

	res, err := apiServer.CreateRoom(
		context.TODO(),
		createRoomRequest,
	)

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse.Name, res.Name)
	assert.Equal(t, expectedResponse.Name, res.Name)
	assert.Contains(t, res.Users, createRoomRequest.UserIds[0])
	assert.Contains(t, res.Users, expectedClaimsResult.Id)
}

func TestApiServer_ListRoomsFailsIfPageSizeExceedsLimit(t *testing.T) {
	setupTest()

	res, err := apiServer.ListRooms(
		context.TODO(),
		&proto.ListRoomsRequest{
			PageSize: 101,
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.InvalidArgument, "page_size cannot be bigger than 100"), err)
}

func TestApiServer_ListRoomsFailsIfAuthFails(t *testing.T) {
	setupTest()

	expectedError := errors.New("some_error")

	ctx := context.TODO()
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(nil, expectedError)

	res, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			PageSize: 100,
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, expectedError, err)
}

func TestApiServer_ListRoomsFailsIfInvalidUserId(t *testing.T) {
	setupTest()

	ctx := context.TODO()
	userClaims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id: "",
		},
	}
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(userClaims, nil)

	res, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			PageSize: 100,
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.Internal, "cannot parse uuid: %v", "invalid UUID length: 0"), err)
}

func TestApiServer_ListRoomsFailsIfDatabaseFailsWithNoPageToken(t *testing.T) {
	setupTest()

	expectedError := errors.New("some_error")

	ctx := context.TODO()
	userId := uuid.New()
	pageSize := 100
	userClaims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id: userId.String(),
		},
		Username: "",
		Role:     model.USER_ROLE,
	}
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(userClaims, nil)
	roomStoreMock.On("ListRoomsFirst", userId, pageSize).Return(nil, expectedError)

	res, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			PageSize: int32(pageSize),
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.Internal, "cannot get rooms: %v", expectedError), err)
}

func TestApiServer_ListRoomsFailsIfDatabaseFailsWithPageTokenPresent(t *testing.T) {
	setupTest()

	token := encodePageToken(utils.Now())
	lastMessageTime, _ := decodePageToken(token)
	expectedError := errors.New("some_error")

	ctx := context.TODO()
	userId := uuid.New()
	pageSize := 100
	userClaims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id: userId.String(),
		},
		Username: "",
		Role:     model.USER_ROLE,
	}
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(userClaims, nil)
	roomStoreMock.On("ListRooms", userId, *lastMessageTime, pageSize).Return(nil, expectedError)

	res, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			NextToken: &wrapperspb.StringValue{
				Value: token,
			},
			PageSize: int32(pageSize),
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.Internal, "cannot get rooms: %v", expectedError), err)
}

func TestApiServer_ListRoomsFailsIfInvalidPageToken(t *testing.T) {
	setupTest()

	token := "some_invalid_token"
	ctx := context.TODO()
	userId := uuid.New()
	pageSize := 100
	userClaims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id: userId.String(),
		},
		Username: "",
		Role:     model.USER_ROLE,
	}
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(userClaims, nil)
	roomStoreMock.On("ListRoomsFirst", userId, pageSize).Return(nil, nil)

	res, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			NextToken: &wrapperspb.StringValue{
				Value: token,
			},
			PageSize: int32(pageSize),
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.InvalidArgument, "cannot parse next token: %v", "illegal base64 data at input byte 4"), err)
}

func TestApiServer_ListRoomsSuccess(t *testing.T) {
	setupTest()

	ctx := context.TODO()
	userId := uuid.New()
	pageSize := 2
	userClaims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id: userId.String(),
		},
		Username: "",
		Role:     model.USER_ROLE,
	}
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(userClaims, nil)
	room := *model.NewRoom("name", true)
	roomStoreMock.On("ListRoomsFirst", userId, pageSize).Return([]model.Room{
		room,
		room,
	}, nil)

	res1, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			PageSize: int32(pageSize),
		},
	)

	nextToken := res1.NextToken.Value
	lastMessageTime, _ := decodePageToken(nextToken)
	roomStoreMock.On("ListRooms", userId, *lastMessageTime, pageSize).Return([]model.Room{
		room,
	}, nil)

	res2, err := apiServer.ListRooms(
		ctx,
		&proto.ListRoomsRequest{
			NextToken: &wrapperspb.StringValue{
				Value: nextToken,
			},
			PageSize: int32(pageSize),
		},
	)

	assert.Equal(
		t,
		&proto.ListRoomsResponse{
			NextToken: &wrapperspb.StringValue{
				Value: nextToken,
			},
			Rooms: []*proto.Room{
				room.PbRoom(),
				room.PbRoom(),
			},
		},
		res1,
	)
	assert.Equal(
		t,
		&proto.ListRoomsResponse{
			NextToken: nil,
			Rooms: []*proto.Room{
				room.PbRoom(),
			},
		},
		res2,
	)
	assert.Nil(t, err)
}
