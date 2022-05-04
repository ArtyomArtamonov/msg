package server

import (
	context "context"
	"errors"
	"testing"

	"github.com/ArtyomArtamonov/msg/internal/mocks"
	"github.com/ArtyomArtamonov/msg/internal/model"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthInterceptor_UnaryFailsWithNoToken(t *testing.T) {
	setupTest()

	expectedRequest := &struct{}{}

	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		// empty metadata
	})

	endpoint := endpoints.MessageService.SendMessage
	authInterceptor := NewAuthInterceptor(jwtManagerMock, endpointRoles)
	res, err := authInterceptor.Unary()(
		contextMock,
		expectedRequest,
		&grpc.UnaryServerInfo{
			Server:     nil,
			FullMethod: endpoint,
		},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.FailNow(t, "Should not have been called")
			return nil, nil
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "authorization token is not provided"), err)
}

func TestAuthInterceptor_UnaryFailsWithInvalidToken(t *testing.T) {
	setupTest()

	expectedRequest := &struct{}{}
	expectedError := errors.New("some_error")

	token := "some_access_token"
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {token},
	})
	jwtManagerMock.On("Verify", token).Return(nil, expectedError)

	endpoint := endpoints.MessageService.SendMessage
	authInterceptor := NewAuthInterceptor(jwtManagerMock, endpointRoles)
	res, err := authInterceptor.Unary()(
		contextMock,
		expectedRequest,
		&grpc.UnaryServerInfo{
			Server:     nil,
			FullMethod: endpoint,
		},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.FailNow(t, "Should not have been called")
			return nil, nil
		},
	)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "authorization token is invalid: some_error"), err)
}

func TestAuthInterceptor_UnarySussess(t *testing.T) {
	setupTest()

	expectedRequest := &pb.MessageRequest{
		To:      "some_user",
		Message: "some message",
	}
	expectedResponse := &pb.MessageRequestStatus{
		Success: true,
		Message: "some message",
	}
	expectedError := errors.New("some_error")

	token := "some_access_token"
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {token},
	})
	jwtManagerMock.On("Verify", token).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)

	endpoint := endpoints.MessageService.SendMessage
	authInterceptor := NewAuthInterceptor(jwtManagerMock, endpointRoles)
	res, err := authInterceptor.Unary()(
		contextMock,
		expectedRequest,
		&grpc.UnaryServerInfo{
			Server:     nil,
			FullMethod: endpoint,
		},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, expectedRequest, req)
			return expectedResponse, expectedError
		},
	)

	assert.Equal(t, expectedResponse, res)
	assert.ErrorIs(t, expectedError, err)
}

func TestAuthInterceptor_StreamFailsWithNoToken(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}

	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		// empty metadata
	})
	expectedStream.On("Context").Return(contextMock)

	endpoint := endpoints.MessageService.SendMessage
	authInterceptor := NewAuthInterceptor(jwtManagerMock, endpointRoles)
	err := authInterceptor.Stream()(
		contextMock,
		expectedStream,
		&grpc.StreamServerInfo{
			FullMethod:     endpoint,
			IsClientStream: false,
			IsServerStream: true,
		},
		func(srv interface{}, stream grpc.ServerStream) error {
			assert.FailNow(t, "Should not have been called")
			return nil
		},
	)

	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "authorization token is not provided"), err)
}

func TestAuthInterceptor_StreamFailsWithInvalidToken(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}
	expectedError := errors.New("some_error")

	token := "some_access_token"
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {token},
	})
	expectedStream.On("Context").Return(contextMock)
	jwtManagerMock.On("Verify", token).Return(nil, expectedError)

	endpoint := endpoints.MessageService.SendMessage
	authInterceptor := NewAuthInterceptor(jwtManagerMock, endpointRoles)
	err := authInterceptor.Stream()(
		contextMock,
		expectedStream,
		&grpc.StreamServerInfo{
			FullMethod:     endpoint,
			IsClientStream: false,
			IsServerStream: true,
		},
		func(srv interface{}, stream grpc.ServerStream) error {
			assert.FailNow(t, "Should not have been called")
			return nil
		},
	)

	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "authorization token is invalid: some_error"), err)
}

func TestAuthInterceptor_StreamSussess(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}
	expectedError := errors.New("some_error")

	token := "some_access_token"
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {token},
	})
	expectedStream.On("Context").Return(contextMock)
	jwtManagerMock.On("Verify", token).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)

	endpoint := endpoints.MessageService.SendMessage
	authInterceptor := NewAuthInterceptor(jwtManagerMock, endpointRoles)
	err := authInterceptor.Stream()(
		contextMock,
		expectedStream,
		&grpc.StreamServerInfo{
			FullMethod:     endpoint,
			IsClientStream: false,
			IsServerStream: true,
		},
		func(srv interface{}, stream grpc.ServerStream) error {
			assert.Equal(t, stream, expectedStream)
			return expectedError
		},
	)

	assert.ErrorIs(t, expectedError, err)
}
