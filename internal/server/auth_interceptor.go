package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ArtyomArtamonov/msg/internal/service"
)

type AuthInterceptor struct {
	jwtManager    service.JWTManagerProtol
	endpointRoles EndpointRoles
}

func NewAuthInterceptor(jwtManager service.JWTManagerProtol, endpointRoles EndpointRoles) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:    jwtManager,
		endpointRoles: endpointRoles,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logrus.Trace("Unary auth interceptor here: ", info.FullMethod)

		err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logrus.Trace("Streaming auth interceptor here: ", info.FullMethod)

		err := i.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) error {
	endpointRoles, ok := i.endpointRoles[method]
	if endpointRoles == nil || !ok {
		// anyone has access
		return nil
	}

	claims, err := service.GetAndVerifyClaimsFromContext(ctx, i.jwtManager)
	if err != nil {
		return err
	}

	for _, role := range endpointRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Errorf(codes.PermissionDenied, "user does not have permission")
}
