package server

import (
	"context"
	"errors"
	"testing"

	"github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestApiServer_CreateRoomFailsIfLessThanTwoUsers(t *testing.T) {
	setupTest()

	for i := 0; i < 2; i++ {
		users := []string{}
		for j := 0; j < i; j++ {
			users = append(users, "some_user_id")
		}

		res, err := apiServer.CreateRoom(
			context.TODO(),
			&proto.CreateRoomRequest{
				Name:         "",
				IsDialogRoom: false,
				Users:        users,
			},
		)

		assert.Nil(t, res)
		assert.ErrorIs(t, status.Error(codes.InvalidArgument, "cannot be less than 2 users"), err)
	}
}

func TestApiServer_CreateRoomFailsIfDialogHasMoreThanTwoUsers(t *testing.T) {
	setupTest()

	res, err := apiServer.CreateRoom(
		context.TODO(),
		&proto.CreateRoomRequest{
			Name:         "",
			IsDialogRoom: true,
			Users: []string{
				"some_user_id_1",
				"some_user_id_2",
				"some_user_id_3",
			},
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.InvalidArgument, "dialog room cannot have more than 2 users"), err)
}

func TestApiServer_CreateRoomFailsIfDatabaseQueryFails(t *testing.T) {
	setupTest()

	expectedError := errors.New("some_error")

	roomStoreMock.On("Add", mock.Anything).Return(expectedError)

	res, err := apiServer.CreateRoom(
		context.TODO(),
		&proto.CreateRoomRequest{
			Name:         "",
			IsDialogRoom: false,
			Users: []string{
				"some_user_id_1",
				"some_user_id_2",
			},
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.Internal, "cannot create room: %v", expectedError), err)
}

func TestApiServer_CreateRoomSuccess(t *testing.T) {
	setupTest()

	expectedName := "some_user"
	expectedUsers := []string{
		"some_user_id_1",
		"some_user_id_2",
	}

	roomStoreMock.On("Add", mock.Anything).Return(nil)

	res, err := apiServer.CreateRoom(
		context.TODO(),
		&proto.CreateRoomRequest{
			Name:         expectedName,
			IsDialogRoom: false,
			Users:        expectedUsers,
		},
	)

	assert.NotEmpty(t, res.RoomId)
	assert.Equal(t, res.Name, expectedName)
	assert.Equal(t, res.Users, expectedUsers)
	assert.Nil(t, err)
}
