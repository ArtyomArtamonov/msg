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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthInterceptor_UnaryFailsWithInvalidToken(t *testing.T) {
	setupTest()

	expectedRequest := &struct{}{}
	expectedError := errors.New("some_error")

	ctx := context.TODO()
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(nil, expectedError)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		endpoint: {model.USER_ROLE},
	})
	res, err := authInterceptor.Unary()(
		ctx,
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
	assert.ErrorIs(t, expectedError, err)
}

func TestAuthInterceptor_UnaryFailsWithoutEndpointPermissions(t *testing.T) {
	setupTest()

	expectedRequest := &struct{}{}

	ctx := context.TODO()
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		endpoint: {},
	})
	res, err := authInterceptor.Unary()(
		ctx,
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
	assert.ErrorIs(t, status.Errorf(codes.PermissionDenied, "user does not have permission"), err)
}

func TestAuthInterceptor_UnarySussessForAnyone(t *testing.T) {
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

	ctx := context.TODO()
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		// allow all endpoints for anyone
	})
	res, err := authInterceptor.Unary()(
		ctx,
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

func TestAuthInterceptor_UnarySussessWithPermissions(t *testing.T) {
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

	ctx := context.TODO()
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		endpoint: {model.USER_ROLE},
	})
	res, err := authInterceptor.Unary()(
		ctx,
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

func TestAuthInterceptor_StreamFailsWithInvalidToken(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}
	expectedError := errors.New("some_error")

	ctx := context.TODO()
	expectedStream.On("Context").Return(ctx)
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(nil, expectedError)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		endpoint: {model.USER_ROLE},
	})
	err := authInterceptor.Stream()(
		ctx,
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

	assert.ErrorIs(t, expectedError, err)
}

func TestAuthInterceptor_StreamFailsWithoutEndpointPermissions(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}

	ctx := context.TODO()
	expectedStream.On("Context").Return(ctx)
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		endpoint: {},
	})
	err := authInterceptor.Stream()(
		ctx,
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

	assert.ErrorIs(t, status.Errorf(codes.PermissionDenied, "user does not have permission"), err)
}

func TestAuthInterceptor_StreamSussessForAnyone(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}
	expectedError := errors.New("some_error")

	ctx := context.TODO()
	expectedStream.On("Context").Return(ctx)
	expectedStream.On("Context").Return(ctx)
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		// allow all endpoints for anyone
	})
	err := authInterceptor.Stream()(
		ctx,
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

func TestAuthInterceptor_StreamSussessWithPermissions(t *testing.T) {
	setupTest()

	expectedStream := &mocks.ServerStreamMock{}
	expectedError := errors.New("some_error")

	ctx := context.TODO()
	expectedStream.On("Context").Return(ctx)
	expectedStream.On("Context").Return(ctx)
	jwtManagerMock.On("GetAndVerifyClaims", ctx).Return(&model.UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Role:           model.USER_ROLE,
	}, nil)
	endpoint := "some_endpoint"
	authInterceptor := NewAuthInterceptor(jwtManagerMock, EndpointRoles{
		endpoint: {model.USER_ROLE},
	})
	err := authInterceptor.Stream()(
		ctx,
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
