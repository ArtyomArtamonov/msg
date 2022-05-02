package server

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/repository"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiServer struct {
	pb.UnimplementedApiServiceServer

	jwtManager *service.JWTManager
	roomStore  repository.RoomStore
}

func NewApiServer(jwtManager *service.JWTManager, roomStore repository.RoomStore) *ApiServer {
	return &ApiServer{
		jwtManager: jwtManager,
		roomStore:  roomStore,
	}
}

func (s *ApiServer) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomStatus, error) {
	if len(req.Users) < 2 {
		return nil, status.Error(codes.InvalidArgument, "cannot be less than 2 users")
	}

	if req.IsDialogRoom && len(req.Users) != 2 {
		return nil, status.Error(codes.InvalidArgument, "dialog room cannot have more than 2 users")
	}

	if req.IsDialogRoom && len(req.Users) == 2 {
		id1, err := uuid.Parse(req.Users[0])
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "could not parse uuid")
		}

		id2, err := uuid.Parse(req.Users[1])
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "could not parse uuid")
		}

		room, err := s.roomStore.FindDialogRoom(id1, id2)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.AlreadyExists {
				createRoomStatus := &pb.CreateRoomStatus{
					RoomId: room.Id.String(),
					Name:   room.Name,
					Users:  req.Users,
				}
				return createRoomStatus, err
			} else {
				return nil, status.Errorf(codes.Internal, "cannot create room: %v", err)
			}
		}
	}

	newRoom := model.NewRoom(req.Name, req.IsDialogRoom, req.Users...)
	err := s.roomStore.Add(newRoom)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create room: %v", err)
	}

	createRoomStatus := &pb.CreateRoomStatus{
		RoomId: newRoom.Id.String(),
		Name:   newRoom.Name,
		Users:  req.Users,
	}

	return createRoomStatus, nil
}
