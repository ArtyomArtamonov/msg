package server

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/repository"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ApiServer struct {
	pb.UnimplementedApiServiceServer

	jwtManager service.JWTManagerProtol
	roomStore  repository.RoomStore
}

func NewApiServer(jwtManager service.JWTManagerProtol, roomStore repository.RoomStore) *ApiServer {
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

	// TODO: remove this, this logic should be hidden from the client - we should create new dialog seamlessly, when user sends his first message in chat room
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

func (s *ApiServer) ListRooms(ctx context.Context, req *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	if req.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "page_size cannot be bigger than 100")
	}

	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot parse uuid: %v", err)
	}

	var rooms []model.Room
	if req.PageToken == nil {
		rooms, err = s.roomStore.ListRoomsFirst(userId, int(req.PageSize))
	} else {
		lastMessageTime, e := decodePageToken(req.PageToken.Value)
		if e != nil {
			return nil, status.Errorf(codes.InvalidArgument, "cannot parse next token: %v", e)
		}
		rooms, err = s.roomStore.ListRooms(userId, *lastMessageTime, int(req.PageSize))
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get rooms: %v", err)
	}

	var nextToken string
	if len(rooms) < int(req.PageSize) {
		nextToken = ""
	} else {
		lastRoom := rooms[len(rooms)-1]
		nextToken = encodePageToken(lastRoom.LastMessageTime)
	}

	pbRooms := []*pb.Room{}
	for _, room := range rooms {
		pbRooms = append(pbRooms, room.PbRoom())
	}

	var listRoomsResponse *pb.ListRoomsResponse
	if len(nextToken) == 0 {
		listRoomsResponse = &pb.ListRoomsResponse{
			NextToken: nil,
			Rooms:     pbRooms,
		}
	} else {
		listRoomsResponse = &pb.ListRoomsResponse{
			NextToken: &wrapperspb.StringValue{Value: nextToken},
			Rooms:     pbRooms,
		}
	}

	return listRoomsResponse, nil

}

func decodePageToken(token string) (*time.Time, error) {
	buf, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, string(buf))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func encodePageToken(lastMessageDate time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(lastMessageDate.Format(time.RFC3339)))
}
